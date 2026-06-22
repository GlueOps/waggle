package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glueops/waggle/internal/models/tenant"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CapacityError is returned when a pool's requested VM count cannot be fully
// placed. Placement is all-or-nothing, so nothing is committed; Fit reports
// how many WOULD have fit, for the caller's benefit.
type CapacityError struct {
	Requested int
	Fit       int
}

func (e *CapacityError) Error() string {
	return fmt.Sprintf("insufficient capacity: requested %d, can place %d", e.Requested, e.Fit)
}

// slotCost is the per-VM resource footprint of a slot.
type slotCost struct {
	vcpu, ramGB, diskGB int
}

// hvCapacity is a hypervisor's remaining bookable capacity plus how many VMs
// of the pool being placed already sit on it (for anti-affinity spread).
type hvCapacity struct {
	id            uuid.UUID
	cpuRemaining  int
	ramRemaining  int
	diskRemaining int
	poolCount     int
}

func (h hvCapacity) fits(c slotCost) bool {
	return h.cpuRemaining >= c.vcpu && h.ramRemaining >= c.ramGB && h.diskRemaining >= c.diskGB
}

// planPlacements greedily assigns `count` VMs of cost `c` across `hvs`,
// preferring the hypervisor with the FEWEST VMs already from this pool
// (anti-affinity / spread), tie-broken by the most remaining CPU (balance).
// Returns the chosen hypervisor IDs. If it can't place all `count`, ok is
// false and `fit` is how many it managed before getting stuck.
func planPlacements(hvs []hvCapacity, c slotCost, count int) (ids []uuid.UUID, fit int, ok bool) {
	work := make([]hvCapacity, len(hvs))
	copy(work, hvs)

	ids = make([]uuid.UUID, 0, count)
	for placed := 0; placed < count; placed++ {
		best := -1
		for i := range work {
			if !work[i].fits(c) {
				continue
			}
			if best == -1 {
				best = i
				continue
			}
			if work[i].poolCount < work[best].poolCount ||
				(work[i].poolCount == work[best].poolCount && work[i].cpuRemaining > work[best].cpuRemaining) {
				best = i
			}
		}
		if best == -1 {
			return ids, placed, false
		}
		ids = append(ids, work[best].id)
		work[best].cpuRemaining -= c.vcpu
		work[best].ramRemaining -= c.ramGB
		work[best].diskRemaining -= c.diskGB
		work[best].poolCount++
	}
	return ids, count, true
}

// plan computes a placement plan for `count` VMs of the given slot within the
// datacenter, accounting for capacity already consumed by ALL existing
// placements (across every pool) and biasing the spread by how many VMs of
// `poolID` already sit on each hypervisor. Must run inside the locking tx.
func (s *FleetService) plan(tx *gorm.DB, datacenterID, slotID, poolID uuid.UUID, count int) ([]uuid.UUID, error) {
	var slot tenant.Slot
	if err := tx.First(&slot, "id = ?", slotID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidInput
		}
		return nil, fmt.Errorf("load slot: %w", err)
	}

	var hvs []tenant.Hypervisor
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("datacenter_id = ?", datacenterID).
		Where("schedulable IS TRUE").
		Order("name").Find(&hvs).Error; err != nil {
		return nil, fmt.Errorf("lock hypervisors: %w", err)
	}
	if len(hvs) == 0 {
		return nil, &CapacityError{Requested: count, Fit: 0}
	}

	hvIDs := make([]uuid.UUID, len(hvs))
	for i, h := range hvs {
		hvIDs[i] = h.ID
	}

	// Capacity consumed by every existing placement, summed per hypervisor.
	type consumedRow struct {
		HypervisorID uuid.UUID
		CPU          int
		RAM          int
		Disk         int
	}
	var crows []consumedRow
	if err := tx.Table("placements").
		Select("placements.hypervisor_id AS hypervisor_id, "+
			"COALESCE(SUM(slots.vcpu),0) AS cpu, "+
			"COALESCE(SUM(slots.ram_gb),0) AS ram, "+
			"COALESCE(SUM(slots.disk_gb),0) AS disk").
		Joins("JOIN pools ON pools.id = placements.pool_id").
		Joins("JOIN slots ON slots.id = pools.slot_id").
		Where("placements.hypervisor_id IN ?", hvIDs).
		Group("placements.hypervisor_id").
		Scan(&crows).Error; err != nil {
		return nil, fmt.Errorf("compute consumed capacity: %w", err)
	}
	consumed := make(map[uuid.UUID]consumedRow, len(crows))
	for _, r := range crows {
		consumed[r.HypervisorID] = r
	}

	// How many of THIS pool's VMs already sit on each hypervisor (spread bias).
	poolCounts := make(map[uuid.UUID]int)
	if poolID != uuid.Nil {
		type pcRow struct {
			HypervisorID uuid.UUID
			N            int
		}
		var pcs []pcRow
		if err := tx.Table("placements").
			Select("hypervisor_id, COUNT(*) AS n").
			Where("pool_id = ?", poolID).
			Group("hypervisor_id").
			Scan(&pcs).Error; err != nil {
			return nil, fmt.Errorf("count pool placements: %w", err)
		}
		for _, r := range pcs {
			poolCounts[r.HypervisorID] = r.N
		}
	}

	caps := make([]hvCapacity, 0, len(hvs))
	for _, h := range hvs {
		c := consumed[h.ID]
		caps = append(caps, hvCapacity{
			id: h.ID,
			// Bookable = total − operator headroom − existing-guest allocation
			// (from discovery) − capacity already consumed by Waggle placements.
			cpuRemaining:  (h.CPUTotal - h.CPUReserved - h.CPUUsed) - c.CPU,
			ramRemaining:  (h.RAMGBTotal - h.RAMGBReserved - h.RAMGBUsed) - c.RAM,
			diskRemaining: (h.DiskGBTotal - h.DiskGBReserved - h.DiskGBUsed) - c.Disk,
			poolCount:     poolCounts[h.ID],
		})
	}

	ids, fit, ok := planPlacements(caps, slotCost{slot.VCPU, slot.RAMGB, slot.DiskGB}, count)
	if !ok {
		return nil, &CapacityError{Requested: count, Fit: fit}
	}
	return ids, nil
}

type PoolInput struct {
	DatacenterID uuid.UUID
	SlotID       uuid.UUID
	Name         string
	DesiredCount int
	Metadata     []byte
}

func (in PoolInput) validate() error {
	if in.Name == "" || in.DatacenterID == uuid.Nil || in.SlotID == uuid.Nil || in.DesiredCount < 0 {
		return ErrInvalidInput
	}
	return nil
}

type PlacementView struct {
	ID             uuid.UUID
	PoolID         uuid.UUID
	HypervisorID   uuid.UUID
	HypervisorName string
	VMID           *int
	CreatedAt      time.Time
}

type PoolResult struct {
	Pool       tenant.Pool
	Placements []PlacementView
}

func (s *FleetService) CreatePool(ctx context.Context, orgID uuid.UUID, in PoolInput) (*PoolResult, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var pool tenant.Pool
	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var dc tenant.Datacenter
		if err := tx.First(&dc, "id = ?", in.DatacenterID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrInvalidInput
			}
			return fmt.Errorf("load datacenter: %w", err)
		}

		hvIDs, err := s.plan(tx, in.DatacenterID, in.SlotID, uuid.Nil, in.DesiredCount)
		if err != nil {
			return err
		}

		pool = tenant.Pool{
			DatacenterID: in.DatacenterID,
			SlotID:       in.SlotID,
			Name:         in.Name,
			DesiredCount: in.DesiredCount,
		}
		if len(in.Metadata) > 0 {
			pool.Metadata = datatypes.JSON(in.Metadata)
		}
		if err := tx.Create(&pool).Error; err != nil {
			return fmt.Errorf("create pool: %w", err)
		}

		return createPlacements(tx, pool.ID, hvIDs)
	})
	if err != nil {
		return nil, err
	}
	return s.loadPoolResult(ctx, db, pool.ID)
}

func (s *FleetService) ResizePool(ctx context.Context, orgID, poolID uuid.UUID, desiredCount int) (*PoolResult, error) {
	if desiredCount < 0 {
		return nil, ErrInvalidInput
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var pool tenant.Pool
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&pool, "id = ?", poolID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("load pool: %w", err)
		}

		delta := desiredCount - pool.DesiredCount
		switch {
		case delta > 0:
			hvIDs, err := s.plan(tx, pool.DatacenterID, pool.SlotID, pool.ID, delta)
			if err != nil {
				return err
			}
			if err := createPlacements(tx, pool.ID, hvIDs); err != nil {
				return err
			}
		case delta < 0:
			// LIFO: remove the newest placements first.
			var victims []tenant.Placement
			if err := tx.Where("pool_id = ?", pool.ID).
				Order("created_at DESC").
				Limit(-delta).
				Find(&victims).Error; err != nil {
				return fmt.Errorf("select placements to remove: %w", err)
			}
			ids := make([]uuid.UUID, len(victims))
			for i, v := range victims {
				ids[i] = v.ID
			}
			if len(ids) > 0 {
				if err := tx.Delete(&tenant.Placement{}, "id IN ?", ids).Error; err != nil {
					return fmt.Errorf("delete placements: %w", err)
				}
			}
		}

		return tx.Model(&pool).Update("desired_count", desiredCount).Error
	})
	if err != nil {
		return nil, err
	}
	return s.loadPoolResult(ctx, db, poolID)
}

func (s *FleetService) GetPool(ctx context.Context, orgID, poolID uuid.UUID) (*PoolResult, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return s.loadPoolResult(ctx, db, poolID)
}

func (s *FleetService) ListPools(ctx context.Context, orgID uuid.UUID) ([]tenant.Pool, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var pools []tenant.Pool
	if err := db.WithContext(ctx).Order("created_at").Find(&pools).Error; err != nil {
		return nil, fmt.Errorf("list pools: %w", err)
	}
	return pools, nil
}

func (s *FleetService) DeletePool(ctx context.Context, orgID, poolID uuid.UUID) error {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.First(&tenant.Pool{}, "id = ?", poolID)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return res.Error
		}
		if err := tx.Delete(&tenant.Placement{}, "pool_id = ?", poolID).Error; err != nil {
			return fmt.Errorf("delete placements: %w", err)
		}
		if err := tx.Delete(&tenant.Pool{}, "id = ?", poolID).Error; err != nil {
			return fmt.Errorf("delete pool: %w", err)
		}
		return nil
	})
}

func (s *FleetService) ListPlacements(ctx context.Context, orgID, poolID uuid.UUID) ([]PlacementView, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if err := db.WithContext(ctx).First(&tenant.Pool{}, "id = ?", poolID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("load pool: %w", err)
	}
	return placementsForPool(ctx, db, poolID)
}

// GetPlacement returns a single placement by ID, joined with its hypervisor name.
func (s *FleetService) GetPlacement(ctx context.Context, orgID, placementID uuid.UUID) (*PlacementView, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	return placementByID(ctx, db, placementID)
}

// DeletePlacement removes a single placement from its pool. The pool's
// desired_count is NOT adjusted; callers should resize the pool separately
// if they want waggle to re-fill the vacancy.
func (s *FleetService) DeletePlacement(ctx context.Context, orgID, placementID uuid.UUID) error {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return err
	}
	res := db.WithContext(ctx).Delete(&tenant.Placement{}, "id = ?", placementID)
	if res.Error != nil {
		return fmt.Errorf("delete placement: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// BackfillVMID attaches the externally-assigned Proxmox vmid to a placement
// (the BGP/Proxmox pipeline calls this after the VM is created).
func (s *FleetService) BackfillVMID(ctx context.Context, orgID, placementID uuid.UUID, vmid int) (*PlacementView, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var pl tenant.Placement
	if err := db.WithContext(ctx).First(&pl, "id = ?", placementID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("load placement: %w", err)
	}
	if err := db.WithContext(ctx).Model(&pl).Update("vmid", &vmid).Error; err != nil {
		return nil, fmt.Errorf("update vmid: %w", err)
	}
	return placementByID(ctx, db, placementID)
}

// placementByID fetches a single placement row joined with its hypervisor name.
func placementByID(ctx context.Context, db *gorm.DB, placementID uuid.UUID) (*PlacementView, error) {
	type row struct {
		ID             uuid.UUID
		PoolID         uuid.UUID
		HypervisorID   uuid.UUID
		HypervisorName string
		VMID           *int
		CreatedAt      time.Time
	}
	var r row
	if err := db.WithContext(ctx).
		Table("placements").
		Select("placements.id, placements.pool_id, placements.hypervisor_id, "+
			"hypervisors.name AS hypervisor_name, placements.vmid, placements.created_at").
		Joins("JOIN hypervisors ON hypervisors.id = placements.hypervisor_id").
		Where("placements.id = ?", placementID).
		First(&r).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get placement: %w", err)
	}
	return &PlacementView{
		ID:             r.ID,
		PoolID:         r.PoolID,
		HypervisorID:   r.HypervisorID,
		HypervisorName: r.HypervisorName,
		VMID:           r.VMID,
		CreatedAt:      r.CreatedAt,
	}, nil
}

func createPlacements(tx *gorm.DB, poolID uuid.UUID, hvIDs []uuid.UUID) error {
	if len(hvIDs) == 0 {
		return nil
	}
	placements := make([]tenant.Placement, 0, len(hvIDs))
	for _, hv := range hvIDs {
		placements = append(placements, tenant.Placement{PoolID: poolID, HypervisorID: hv})
	}
	if err := tx.Create(&placements).Error; err != nil {
		return fmt.Errorf("create placements: %w", err)
	}
	return nil
}

func (s *FleetService) loadPoolResult(ctx context.Context, db *gorm.DB, poolID uuid.UUID) (*PoolResult, error) {
	var pool tenant.Pool
	if err := db.WithContext(ctx).First(&pool, "id = ?", poolID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("load pool: %w", err)
	}
	placements, err := placementsForPool(ctx, db, poolID)
	if err != nil {
		return nil, err
	}
	return &PoolResult{Pool: pool, Placements: placements}, nil
}

// FleetPlacementView is a tenant-wide placement enriched with its pool, slot,
// and hypervisor context for overview displays.
type FleetPlacementView struct {
	ID             uuid.UUID
	PoolID         uuid.UUID
	PoolName       string
	HypervisorID   uuid.UUID
	HypervisorName string
	SlotName       string
	VCPU           int
	RAMGB          int
	DiskGB         int
	VMID           *int
	CreatedAt      time.Time
}

// ListAllPlacements returns every placement in the caller's tenant, newest
// first, joined with pool/slot/hypervisor context (for fleet-wide overviews).
func (s *FleetService) ListAllPlacements(ctx context.Context, orgID uuid.UUID) ([]FleetPlacementView, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var rows []FleetPlacementView
	if err := db.WithContext(ctx).
		Table("placements").
		Select("placements.id AS id, placements.pool_id AS pool_id, pools.name AS pool_name, "+
			"placements.hypervisor_id AS hypervisor_id, hypervisors.name AS hypervisor_name, "+
			"slots.name AS slot_name, slots.vcpu AS vcpu, slots.ram_gb AS ram_gb, slots.disk_gb AS disk_gb, "+
			"placements.vmid AS vmid, placements.created_at AS created_at").
		Joins("JOIN pools ON pools.id = placements.pool_id").
		Joins("JOIN hypervisors ON hypervisors.id = placements.hypervisor_id").
		Joins("JOIN slots ON slots.id = pools.slot_id").
		Order("placements.created_at DESC").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list placements: %w", err)
	}
	return rows, nil
}

func placementsForPool(ctx context.Context, db *gorm.DB, poolID uuid.UUID) ([]PlacementView, error) {
	type row struct {
		ID             uuid.UUID
		HypervisorID   uuid.UUID
		HypervisorName string
		VMID           *int
		CreatedAt      time.Time
	}
	var rows []row
	if err := db.WithContext(ctx).
		Table("placements").
		Select("placements.id AS id, placements.hypervisor_id AS hypervisor_id, "+
			"hypervisors.name AS hypervisor_name, placements.vmid AS vmid, placements.created_at AS created_at").
		Joins("JOIN hypervisors ON hypervisors.id = placements.hypervisor_id").
		Where("placements.pool_id = ?", poolID).
		Order("placements.created_at").
		Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list placements: %w", err)
	}
	out := make([]PlacementView, 0, len(rows))
	for _, r := range rows {
		out = append(out, PlacementView{
			ID:             r.ID,
			PoolID:         poolID,
			HypervisorID:   r.HypervisorID,
			HypervisorName: r.HypervisorName,
			VMID:           r.VMID,
			CreatedAt:      r.CreatedAt,
		})
	}
	return out, nil
}
