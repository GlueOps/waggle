package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

type Enqueuer struct {
	river *river.Client[pgx.Tx]
}

func NewEnqueuer(rc *river.Client[pgx.Tx]) (*Enqueuer, error) {
	if rc == nil {
		return nil, errors.New("jobs: nil river client")
	}
	return &Enqueuer{river: rc}, nil
}

func (e *Enqueuer) EnqueueProvisionTenant(ctx context.Context, orgID uuid.UUID) error {
	if e == nil || e.river == nil {
		return errors.New("jobs: enqueuer not initialized")
	}
	if _, err := e.river.Insert(ctx, TenantProvisionerArgs{OrganizationID: orgID.String()}, nil); err != nil {
		return fmt.Errorf("enqueue tenant provisioner: %w", err)
	}
	return nil
}

func (e *Enqueuer) EnqueueDestroyTenant(ctx context.Context, orgID uuid.UUID) error {
	if e == nil || e.river == nil {
		return errors.New("jobs: enqueuer not initialized")
	}
	if _, err := e.river.Insert(ctx, TenantDestroyerArgs{OrganizationID: orgID.String()}, nil); err != nil {
		return fmt.Errorf("enqueue tenant destroyer: %w", err)
	}
	return nil
}

func (e *Enqueuer) EnqueueDiscoverHypervisors(ctx context.Context, orgID, datacenterID uuid.UUID) error {
	if e == nil || e.river == nil {
		return errors.New("jobs: enqueuer not initialized")
	}
	if _, err := e.river.Insert(ctx, HypervisorDiscoveryArgs{
		OrganizationID: orgID.String(),
		DatacenterID:   datacenterID.String(),
	}, nil); err != nil {
		return fmt.Errorf("enqueue hypervisor discovery: %w", err)
	}
	return nil
}
