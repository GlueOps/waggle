package jobs

import (
	"context"
	"fmt"
	"log"

	"github.com/glueops/waggle/internal/models/tenant"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// activeUniqueStates restricts ByArgs uniqueness to jobs that are still in
// flight. River's default ByState also includes JobStateCompleted, which means
// a once-completed unique job keeps blocking duplicate inserts until the job
// cleaner removes it (~24h) — that silently stops a periodic/recurring job from
// ever running again. Limiting uniqueness to the in-flight states (the four
// required by River, plus retryable) coalesces overlapping enqueues without
// blocking the next scheduled run.
var activeUniqueStates = []rivertype.JobState{
	rivertype.JobStateAvailable,
	rivertype.JobStatePending,
	rivertype.JobStateRunning,
	rivertype.JobStateScheduled,
	rivertype.JobStateRetryable,
}

// HypervisorDiscoverer is the subset of FleetService the discovery worker needs.
// Declared here (rather than importing service) because service imports jobs;
// FleetService satisfies it structurally.
type HypervisorDiscoverer interface {
	DiscoverHypervisors(ctx context.Context, orgID, datacenterID uuid.UUID) ([]tenant.Hypervisor, error)
}

type HypervisorDiscoveryService struct {
	fleet HypervisorDiscoverer
}

func NewHypervisorDiscoveryService(fleet HypervisorDiscoverer) *HypervisorDiscoveryService {
	return &HypervisorDiscoveryService{fleet: fleet}
}

// Discover runs discovery for one datacenter. orgID/datacenterID arrive as
// strings from the job payload; an unparseable id cancels the job (retrying
// won't fix bad input).
func (s *HypervisorDiscoveryService) Discover(ctx context.Context, orgID, datacenterID string) error {
	if s == nil || s.fleet == nil {
		return fmt.Errorf("hypervisor discovery service not configured")
	}
	oid, err := uuid.Parse(orgID)
	if err != nil {
		return river.JobCancel(fmt.Errorf("invalid organization_id %q: %w", orgID, err))
	}
	did, err := uuid.Parse(datacenterID)
	if err != nil {
		return river.JobCancel(fmt.Errorf("invalid datacenter_id %q: %w", datacenterID, err))
	}
	hvs, err := s.fleet.DiscoverHypervisors(ctx, oid, did)
	if err != nil {
		return err
	}
	log.Printf("hypervisor_discovery: org=%s dc=%s discovered %d hypervisor(s)", oid, did, len(hvs))
	return nil
}

type HypervisorDiscoveryArgs struct {
	OrganizationID string `json:"organization_id"`
	DatacenterID   string `json:"datacenter_id"`
}

func (HypervisorDiscoveryArgs) Kind() string { return "hypervisor_discovery" }

func (HypervisorDiscoveryArgs) InsertOpts() river.InsertOpts {
	// Coalesce duplicate discovery requests for the same datacenter that are
	// still in flight, but allow a fresh discovery once the previous one has
	// completed (see activeUniqueStates).
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{ByArgs: true, ByState: activeUniqueStates},
	}
}

type HypervisorDiscoveryWorker struct {
	river.WorkerDefaults[HypervisorDiscoveryArgs]
	svc *HypervisorDiscoveryService
}

func NewHypervisorDiscoveryWorker(svc *HypervisorDiscoveryService) *HypervisorDiscoveryWorker {
	return &HypervisorDiscoveryWorker{svc: svc}
}

func (w *HypervisorDiscoveryWorker) Work(ctx context.Context, job *river.Job[HypervisorDiscoveryArgs]) error {
	return w.svc.Discover(ctx, job.Args.OrganizationID, job.Args.DatacenterID)
}
