package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/glueops/waggle/internal/database"
	"github.com/glueops/waggle/internal/models/tenant"
	"github.com/glueops/waggle/internal/proxmox"
	"github.com/glueops/waggle/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrNotFound is returned by FleetService when a resource doesn't exist in
// the caller's tenant database.
var ErrNotFound = errors.New("not found")

// ErrDiscovery wraps failures talking to the upstream Proxmox cluster during
// hypervisor discovery (network, auth, or unexpected API responses). Distinct
// from ErrNotFound/ErrInvalidInput so the API can return a gateway error.
var ErrDiscovery = errors.New("hypervisor discovery failed")

const bytesPerGB = 1024 * 1024 * 1024

// ReservationDefaults is the capacity held back from placement (OS/host
// overhead) applied to newly-discovered hypervisors.
//
// The CPU overcommit ratio deliberately does NOT live here: this struct comes
// from process config, which is shared by every tenant, and capacity policy
// must not cross tenant boundaries. It is carried per-datacenter in the tenant
// DB instead (Datacenter.CPUOvercommitRatio).
type ReservationDefaults struct {
	CPU    int
	RAMGB  int
	DiskGB int
}

// FleetService manages tenant-scoped placement resources (datacenters,
// slots, and — in later chunks — hypervisors, pools, placements). Each
// method resolves the caller's tenant DB via the TenantManager.
type FleetService struct {
	tenants *database.TenantManager
	reserve ReservationDefaults
}

func NewFleetService(tenants *database.TenantManager, reserve ReservationDefaults) *FleetService {
	return &FleetService{tenants: tenants, reserve: reserve}
}

func (s *FleetService) db(ctx context.Context, orgID uuid.UUID) (*gorm.DB, error) {
	if s == nil || s.tenants == nil {
		return nil, errors.New("fleet service: tenant manager not configured")
	}
	return s.tenants.For(ctx, orgID)
}

// ---- Datacenter ----

type DatacenterInput struct {
	Name string
	URL  string
	// Token is the Proxmox API token used for discovery. nil means "leave
	// unchanged" (on update) or "none" (on create); a non-empty value is
	// encrypted with the tenant DEK before storage and never returned.
	Token *string
	// InsecureSkipVerify: nil means "default false" on create / "leave
	// unchanged" on update; non-nil sets it explicitly.
	InsecureSkipVerify *bool
	// CPUOvercommitRatio is the default stamped onto hypervisors discovered in
	// this datacenter. nil means "default 1.0" on create / "leave unchanged"
	// on update; non-nil sets it explicitly.
	CPUOvercommitRatio *float64
}

func (in DatacenterInput) validate() error {
	if in.Name == "" || in.URL == "" {
		return ErrInvalidInput
	}
	// Same bounds as the per-hypervisor override — see MaxCPUOvercommitRatio.
	if in.CPUOvercommitRatio != nil &&
		(*in.CPUOvercommitRatio <= 0 || *in.CPUOvercommitRatio > MaxCPUOvercommitRatio) {
		return ErrInvalidInput
	}
	return nil
}

// encryptToken seals a plaintext PVE token with the org's tenant DEK, returning
// the base64 ciphertext/iv/tag for storage on the Datacenter row.
func (s *FleetService) encryptToken(ctx context.Context, orgID uuid.UUID, token string) (ct, iv, tag string, err error) {
	dek, err := s.tenants.TenantDEK(ctx, orgID)
	if err != nil {
		return "", "", "", err
	}
	ctb, ivb, tagb, err := utils.EncryptAESGCM([]byte(token), dek)
	if err != nil {
		return "", "", "", fmt.Errorf("encrypt datacenter token: %w", err)
	}
	return utils.EncodeB64(ctb), utils.EncodeB64(ivb), utils.EncodeB64(tagb), nil
}

// DatacenterToken decrypts and returns the stored PVE token for a datacenter.
// Returns ErrInvalidInput when no token has been configured.
func (s *FleetService) DatacenterToken(ctx context.Context, orgID uuid.UUID, dc *tenant.Datacenter) (string, error) {
	if dc.EncryptedTokenKey == "" || dc.TokenKeyIV == "" || dc.TokenKeyTag == "" {
		return "", fmt.Errorf("%w: datacenter has no API token configured", ErrInvalidInput)
	}
	dek, err := s.tenants.TenantDEK(ctx, orgID)
	if err != nil {
		return "", err
	}
	ct, err := utils.DecodeB64(dc.EncryptedTokenKey)
	if err != nil {
		return "", fmt.Errorf("decode datacenter token: %w", err)
	}
	iv, err := utils.DecodeB64(dc.TokenKeyIV)
	if err != nil {
		return "", fmt.Errorf("decode datacenter token iv: %w", err)
	}
	tag, err := utils.DecodeB64(dc.TokenKeyTag)
	if err != nil {
		return "", fmt.Errorf("decode datacenter token tag: %w", err)
	}
	pt, err := utils.DecryptAESGCM(ct, dek, iv, tag)
	if err != nil {
		return "", fmt.Errorf("decrypt datacenter token: %w", err)
	}
	return string(pt), nil
}

func (s *FleetService) CreateDatacenter(ctx context.Context, orgID uuid.UUID, in DatacenterInput) (*tenant.Datacenter, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	dc := &tenant.Datacenter{
		Name: in.Name,
		Url:  in.URL,
		// Sell cores 1:1 unless asked otherwise, so adding the feature does
		// not change the capacity of any fleet that never opts in.
		CPUOvercommitRatio: tenant.DefaultCPUOvercommitRatio,
	}
	if in.InsecureSkipVerify != nil {
		dc.InsecureSkipVerify = *in.InsecureSkipVerify
	}
	if in.CPUOvercommitRatio != nil {
		dc.CPUOvercommitRatio = *in.CPUOvercommitRatio
	}
	if in.Token != nil && *in.Token != "" {
		ct, iv, tag, err := s.encryptToken(ctx, orgID, *in.Token)
		if err != nil {
			return nil, err
		}
		dc.EncryptedTokenKey, dc.TokenKeyIV, dc.TokenKeyTag = ct, iv, tag
	}
	if err := db.WithContext(ctx).Create(dc).Error; err != nil {
		return nil, fmt.Errorf("create datacenter: %w", err)
	}
	return dc, nil
}

func (s *FleetService) ListDatacenters(ctx context.Context, orgID uuid.UUID) ([]tenant.Datacenter, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var dcs []tenant.Datacenter
	if err := db.WithContext(ctx).Order("created_at").Find(&dcs).Error; err != nil {
		return nil, fmt.Errorf("list datacenters: %w", err)
	}
	return dcs, nil
}

func (s *FleetService) GetDatacenter(ctx context.Context, orgID, id uuid.UUID) (*tenant.Datacenter, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var dc tenant.Datacenter
	if err := db.WithContext(ctx).First(&dc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get datacenter: %w", err)
	}
	return &dc, nil
}

func (s *FleetService) UpdateDatacenter(ctx context.Context, orgID, id uuid.UUID, in DatacenterInput) (*tenant.Datacenter, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var dc tenant.Datacenter
	if err := db.WithContext(ctx).First(&dc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get datacenter: %w", err)
	}
	dc.Name = in.Name
	dc.Url = in.URL
	if in.InsecureSkipVerify != nil {
		dc.InsecureSkipVerify = *in.InsecureSkipVerify
	}
	// Affects nodes discovered from here on. Existing hypervisors keep the
	// ratio they were stamped with; re-rating a live fleet has to be an
	// explicit per-node act, not a side effect of editing the datacenter.
	if in.CPUOvercommitRatio != nil {
		dc.CPUOvercommitRatio = *in.CPUOvercommitRatio
	}
	if in.Token != nil && *in.Token != "" {
		ct, iv, tag, err := s.encryptToken(ctx, orgID, *in.Token)
		if err != nil {
			return nil, err
		}
		dc.EncryptedTokenKey, dc.TokenKeyIV, dc.TokenKeyTag = ct, iv, tag
	}
	if err := db.WithContext(ctx).Save(&dc).Error; err != nil {
		return nil, fmt.Errorf("update datacenter: %w", err)
	}
	return &dc, nil
}

func (s *FleetService) DeleteDatacenter(ctx context.Context, orgID, id uuid.UUID) error {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return err
	}
	res := db.WithContext(ctx).Delete(&tenant.Datacenter{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete datacenter: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// ---- Slot ----

type SlotInput struct {
	Name   string
	VCPU   int
	RAMGB  int
	DiskGB int
}

func (in SlotInput) validate() error {
	if in.Name == "" || in.VCPU <= 0 || in.RAMGB <= 0 || in.DiskGB <= 0 {
		return ErrInvalidInput
	}
	return nil
}

func (s *FleetService) CreateSlot(ctx context.Context, orgID uuid.UUID, in SlotInput) (*tenant.Slot, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	slot := &tenant.Slot{Name: in.Name, VCPU: in.VCPU, RAMGB: in.RAMGB, DiskGB: in.DiskGB}
	if err := db.WithContext(ctx).Create(slot).Error; err != nil {
		return nil, fmt.Errorf("create slot: %w", err)
	}
	return slot, nil
}

func (s *FleetService) ListSlots(ctx context.Context, orgID uuid.UUID, name string) ([]tenant.Slot, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	q := db.WithContext(ctx).Order("name")
	if name != "" {
		q = q.Where("name = ?", name)
	}
	var slots []tenant.Slot
	if err := q.Find(&slots).Error; err != nil {
		return nil, fmt.Errorf("list slots: %w", err)
	}
	return slots, nil
}

func (s *FleetService) GetSlot(ctx context.Context, orgID, id uuid.UUID) (*tenant.Slot, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var slot tenant.Slot
	if err := db.WithContext(ctx).First(&slot, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get slot: %w", err)
	}
	return &slot, nil
}

func (s *FleetService) UpdateSlot(ctx context.Context, orgID, id uuid.UUID, in SlotInput) (*tenant.Slot, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var slot tenant.Slot
	if err := db.WithContext(ctx).First(&slot, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get slot: %w", err)
	}
	slot.Name = in.Name
	slot.VCPU = in.VCPU
	slot.RAMGB = in.RAMGB
	slot.DiskGB = in.DiskGB
	if err := db.WithContext(ctx).Save(&slot).Error; err != nil {
		return nil, fmt.Errorf("update slot: %w", err)
	}
	return &slot, nil
}

func (s *FleetService) DeleteSlot(ctx context.Context, orgID, id uuid.UUID) error {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return err
	}
	res := db.WithContext(ctx).Delete(&tenant.Slot{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete slot: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// ---- Hypervisor ----

type HypervisorInput struct {
	DatacenterID   uuid.UUID
	Name           string
	CPUTotal       int
	CPUReserved    int
	RAMGBTotal     int
	RAMGBReserved  int
	DiskGBTotal    int
	DiskGBReserved int
	// Schedulable: nil means "default true" on create / "leave unchanged" on
	// update; non-nil sets it explicitly.
	Schedulable *bool
	// CPUOvercommitRatio: nil means "use the configured default" on create /
	// "leave unchanged" on update; non-nil sets it explicitly. Same
	// omit-means-don't-touch contract as Schedulable.
	CPUOvercommitRatio *float64
}

// MaxCPUOvercommitRatio caps how far a node may be oversold. Not a technical
// limit — a guard against a fat-fingered 40 where 4.0 was meant, which would
// silently book ten times the intended VMs. Also keeps the value inside the
// column's numeric(5,2) range.
const MaxCPUOvercommitRatio = 64.0

func (in HypervisorInput) validate() error {
	if in.Name == "" {
		return ErrInvalidInput
	}
	if in.CPUTotal < 0 || in.CPUReserved < 0 ||
		in.RAMGBTotal < 0 || in.RAMGBReserved < 0 ||
		in.DiskGBTotal < 0 || in.DiskGBReserved < 0 {
		return ErrInvalidInput
	}
	// Reserved is checked against PHYSICAL cores, not the overcommitted pool:
	// headroom means "cores the host keeps", so its ceiling shouldn't move
	// when the operator changes how aggressively the rest is sold.
	if in.CPUReserved > in.CPUTotal ||
		in.RAMGBReserved > in.RAMGBTotal ||
		in.DiskGBReserved > in.DiskGBTotal {
		return ErrInvalidInput
	}
	// Ratios below 1 (deliberate underselling) are allowed; zero or negative
	// is not — it would zero out the node's capacity rather than express a
	// policy, and is indistinguishable from an unset field.
	if in.CPUOvercommitRatio != nil &&
		(*in.CPUOvercommitRatio <= 0 || *in.CPUOvercommitRatio > MaxCPUOvercommitRatio) {
		return ErrInvalidInput
	}
	return nil
}

// loadDatacenter fetches a datacenter by ID from the tenant DB, doubling as
// the existence check for inbound references. A missing row is treated as a
// bad reference in the request (ErrInvalidInput) rather than ErrNotFound.
// Callers that only need validity can ignore the returned row; hypervisor
// creation uses it to inherit the datacenter's overcommit default.
func loadDatacenter(ctx context.Context, db *gorm.DB, id uuid.UUID) (*tenant.Datacenter, error) {
	var dc tenant.Datacenter
	if err := db.WithContext(ctx).First(&dc, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidInput
		}
		return nil, fmt.Errorf("verify datacenter: %w", err)
	}
	return &dc, nil
}

func (s *FleetService) CreateHypervisor(ctx context.Context, orgID uuid.UUID, in HypervisorInput) (*tenant.Hypervisor, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	dc, err := loadDatacenter(ctx, db, in.DatacenterID)
	if err != nil {
		return nil, err
	}
	hv := &tenant.Hypervisor{
		DatacenterID:   in.DatacenterID,
		Name:           in.Name,
		CPUTotal:       in.CPUTotal,
		CPUReserved:    in.CPUReserved,
		RAMGBTotal:     in.RAMGBTotal,
		RAMGBReserved:  in.RAMGBReserved,
		DiskGBTotal:    in.DiskGBTotal,
		DiskGBReserved: in.DiskGBReserved,
		Schedulable:    in.Schedulable == nil || *in.Schedulable,
		// An omitted ratio inherits the datacenter's default, so a manually
		// registered node lands on the same footing as a discovered one.
		CPUOvercommitRatio: dc.EffectiveCPUOvercommitRatio(),
	}
	if in.CPUOvercommitRatio != nil {
		hv.CPUOvercommitRatio = *in.CPUOvercommitRatio
	}
	if err := db.WithContext(ctx).Create(hv).Error; err != nil {
		return nil, fmt.Errorf("create hypervisor: %w", err)
	}
	return hv, nil
}

func (s *FleetService) ListHypervisors(ctx context.Context, orgID uuid.UUID, datacenterID *uuid.UUID) ([]tenant.Hypervisor, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	q := db.WithContext(ctx)
	if datacenterID != nil {
		q = q.Where("datacenter_id = ?", *datacenterID)
	}
	var hvs []tenant.Hypervisor
	if err := q.Order("name").Find(&hvs).Error; err != nil {
		return nil, fmt.Errorf("list hypervisors: %w", err)
	}
	if err := attachConsumption(db.WithContext(ctx), hvs); err != nil {
		return nil, err
	}
	return hvs, nil
}

func (s *FleetService) GetHypervisor(ctx context.Context, orgID, id uuid.UUID) (*tenant.Hypervisor, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	hvs := make([]tenant.Hypervisor, 1)
	if err := db.WithContext(ctx).First(&hvs[0], "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get hypervisor: %w", err)
	}
	if err := attachConsumption(db.WithContext(ctx), hvs); err != nil {
		return nil, err
	}
	return &hvs[0], nil
}

func (s *FleetService) UpdateHypervisor(ctx context.Context, orgID, id uuid.UUID, in HypervisorInput) (*tenant.Hypervisor, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var hv tenant.Hypervisor
	if err := db.WithContext(ctx).First(&hv, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get hypervisor: %w", err)
	}
	if in.DatacenterID != hv.DatacenterID {
		if _, err := loadDatacenter(ctx, db, in.DatacenterID); err != nil {
			return nil, err
		}
	}
	hv.DatacenterID = in.DatacenterID
	hv.Name = in.Name
	hv.CPUTotal = in.CPUTotal
	hv.CPUReserved = in.CPUReserved
	hv.RAMGBTotal = in.RAMGBTotal
	hv.RAMGBReserved = in.RAMGBReserved
	hv.DiskGBTotal = in.DiskGBTotal
	hv.DiskGBReserved = in.DiskGBReserved
	if in.Schedulable != nil {
		hv.Schedulable = *in.Schedulable
	}
	// Omitted means "leave unchanged": moving a node between datacenters must
	// not silently re-rate it to the new datacenter's default.
	if in.CPUOvercommitRatio != nil {
		hv.CPUOvercommitRatio = *in.CPUOvercommitRatio
	}
	if err := db.WithContext(ctx).Save(&hv).Error; err != nil {
		return nil, fmt.Errorf("update hypervisor: %w", err)
	}
	one := []tenant.Hypervisor{hv}
	if err := attachConsumption(db.WithContext(ctx), one); err != nil {
		return nil, err
	}
	return &one[0], nil
}

func (s *FleetService) DeleteHypervisor(ctx context.Context, orgID, id uuid.UUID) error {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return err
	}
	res := db.WithContext(ctx).Delete(&tenant.Hypervisor{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete hypervisor: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// DiscoverHypervisors queries the datacenter's Proxmox cluster and upserts a
// hypervisor row per node, keyed on (datacenter_id, name). Discovered totals
// (CPU/RAM/disk) and last_synced_at are refreshed; operator-managed fields
// (reserved capacity, schedulable) are preserved on existing rows. New rows
// default to schedulable=true. Returns the datacenter's current hypervisor set.
func (s *FleetService) DiscoverHypervisors(ctx context.Context, orgID, datacenterID uuid.UUID) ([]tenant.Hypervisor, error) {
	db, err := s.db(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var dc tenant.Datacenter
	if err := db.WithContext(ctx).First(&dc, "id = ?", datacenterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get datacenter: %w", err)
	}

	token, err := s.DatacenterToken(ctx, orgID, &dc)
	if err != nil {
		return nil, err
	}

	client, err := proxmox.New(proxmox.Config{
		BaseURL:            dc.Url,
		Token:              token,
		InsecureSkipVerify: dc.InsecureSkipVerify,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDiscovery, err)
	}
	nodes, err := client.ListNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDiscovery, err)
	}

	// vmids of guests Waggle already tracks in its placement ledger for this
	// datacenter. NodeUsage excludes them so discovered CPUUsed reflects only
	// non-Waggle guests; the ledger accounts for Waggle's own VMs at scheduling
	// time. Counting both would double-subtract their capacity. vmid is
	// cluster-unique in Proxmox, so a flat set (not keyed per hypervisor) is
	// both sufficient and drift-safe: a migrated VM is still excluded wherever
	// discovery finds it, staying counted exactly once via the ledger.
	var managedVMIDs []int
	if err := db.WithContext(ctx).
		Table("placements").
		Joins("JOIN hypervisors ON hypervisors.id = placements.hypervisor_id").
		Where("hypervisors.datacenter_id = ?", dc.ID).
		Where("placements.vmid IS NOT NULL").
		Distinct().
		Pluck("placements.vmid", &managedVMIDs).Error; err != nil {
		return nil, fmt.Errorf("load managed vmids: %w", err)
	}
	managed := make(map[int]struct{}, len(managedVMIDs))
	for _, v := range managedVMIDs {
		managed[v] = struct{}{}
	}

	now := time.Now().UTC()
	for _, n := range nodes {
		// Capacity already committed to existing guests, so bookable space
		// reflects reality rather than assuming an empty node. Best-effort:
		// usage errors don't abort discovery of the rest of the node's data.
		usage, uerr := client.NodeUsage(ctx, n.Name, managed)
		if uerr != nil {
			usage = proxmox.NodeUsage{}
		}
		hv := tenant.Hypervisor{
			DatacenterID: dc.ID,
			Name:         n.Name,
			CPUTotal:     n.MaxCPU,
			RAMGBTotal:   int(n.MaxMem / bytesPerGB),
			DiskGBTotal:  int(n.MaxDisk / bytesPerGB),
			CPUUsed:      usage.VCPU,
			RAMGBUsed:    int(usage.MemBytes / bytesPerGB),
			DiskGBUsed:   int(usage.DiskBytes / bytesPerGB),
			// Reserved + Schedulable + overcommit apply to INSERT only
			// (excluded from DoUpdates), so operator overrides survive
			// re-discovery. The overcommit ratio comes from the datacenter
			// rather than process config so the policy stays tenant-scoped.
			CPUReserved:        s.reserve.CPU,
			RAMGBReserved:      s.reserve.RAMGB,
			DiskGBReserved:     s.reserve.DiskGB,
			Schedulable:        true,
			CPUOvercommitRatio: dc.EffectiveCPUOvercommitRatio(),
			LastSyncedAt:       &now,
		}
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "datacenter_id"}, {Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"cpu_total", "ram_gb_total", "disk_gb_total",
				"cpu_used", "ram_gb_used", "disk_gb_used",
				"last_synced_at", "updated_at",
			}),
		}).Create(&hv).Error; err != nil {
			return nil, fmt.Errorf("upsert discovered hypervisor %q: %w", n.Name, err)
		}
	}

	return s.ListHypervisors(ctx, orgID, &datacenterID)
}
