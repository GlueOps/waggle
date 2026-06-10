package jobs

import (
	"context"
	"fmt"
	"log"

	"github.com/glueops/waggle/internal/models/control"
	"github.com/glueops/waggle/internal/models/tenant"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

// DatacenterLister is the subset of FleetService the sweep needs to enumerate a
// tenant's datacenters. Declared here (rather than importing service) because
// service imports jobs; FleetService satisfies it structurally.
type DatacenterLister interface {
	ListDatacenters(ctx context.Context, orgID uuid.UUID) ([]tenant.Datacenter, error)
}

// HypervisorDiscoverySweepService fans a single periodic tick out into one
// per-datacenter discovery job for every active tenant. It does the cheap
// enumeration (active orgs -> datacenters) and lets the existing
// hypervisor_discovery worker do the expensive Proxmox calls, so a slow or
// unreachable cluster can't stall the sweep or the other datacenters.
type HypervisorDiscoverySweepService struct {
	control *gorm.DB
	fleet   DatacenterLister
}

func NewHypervisorDiscoverySweepService(control *gorm.DB, fleet DatacenterLister) *HypervisorDiscoverySweepService {
	return &HypervisorDiscoverySweepService{control: control, fleet: fleet}
}

// Sweep enqueues a hypervisor_discovery job for every token-configured
// datacenter across all active tenants. Per-org failures are logged and
// skipped so one broken tenant doesn't abort the whole sweep. Returns the
// number of discovery jobs enqueued.
func (s *HypervisorDiscoverySweepService) Sweep(ctx context.Context) (int, error) {
	if s == nil || s.control == nil || s.fleet == nil {
		return 0, fmt.Errorf("hypervisor discovery sweep service not configured")
	}

	// River injects the client into the worker context; we reuse it to enqueue
	// the per-datacenter jobs.
	client, err := river.ClientFromContextSafely[pgx.Tx](ctx)
	if err != nil {
		return 0, fmt.Errorf("hypervisor discovery sweep: %w", err)
	}

	var orgs []control.Organization
	if err := s.control.WithContext(ctx).
		Where("status = ?", "active").
		Find(&orgs).Error; err != nil {
		return 0, fmt.Errorf("list active organizations: %w", err)
	}

	enqueued := 0
	for _, org := range orgs {
		dcs, err := s.fleet.ListDatacenters(ctx, org.ID)
		if err != nil {
			log.Printf("hypervisor_discovery_sweep: skip org=%s: %v", org.ID, err)
			continue
		}
		for _, dc := range dcs {
			// Datacenters without an API token can't be discovered; skip them
			// rather than enqueue jobs that are guaranteed to fail.
			if dc.EncryptedTokenKey == "" {
				continue
			}
			// ByArgs uniqueness on HypervisorDiscoveryArgs coalesces this with
			// any still-pending discovery for the same datacenter (e.g. a manual
			// trigger), so overlapping sweeps don't pile up duplicates.
			if _, err := client.Insert(ctx, HypervisorDiscoveryArgs{
				OrganizationID: org.ID.String(),
				DatacenterID:   dc.ID.String(),
			}, nil); err != nil {
				log.Printf("hypervisor_discovery_sweep: enqueue org=%s dc=%s: %v", org.ID, dc.ID, err)
				continue
			}
			enqueued++
		}
	}
	return enqueued, nil
}

type HypervisorDiscoverySweepArgs struct{}

func (HypervisorDiscoverySweepArgs) Kind() string { return "hypervisor_discovery_sweep" }

func (HypervisorDiscoverySweepArgs) InsertOpts() river.InsertOpts {
	// Coalesce overlapping sweeps: if a sweep is still pending/running when the
	// next interval fires, don't stack another on top of it.
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{ByArgs: true},
	}
}

type HypervisorDiscoverySweepWorker struct {
	river.WorkerDefaults[HypervisorDiscoverySweepArgs]
	svc *HypervisorDiscoverySweepService
}

func NewHypervisorDiscoverySweepWorker(svc *HypervisorDiscoverySweepService) *HypervisorDiscoverySweepWorker {
	return &HypervisorDiscoverySweepWorker{svc: svc}
}

func (w *HypervisorDiscoverySweepWorker) Work(ctx context.Context, job *river.Job[HypervisorDiscoverySweepArgs]) error {
	n, err := w.svc.Sweep(ctx)
	if err != nil {
		return err
	}
	log.Printf("hypervisor_discovery_sweep: enqueued %d discovery job(s)", n)
	return nil
}
