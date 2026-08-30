package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Placement struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PoolID uuid.UUID `gorm:"type:uuid;index;not null" json:"pool_id"`
	// HypervisorID is the ASSIGNMENT: the hypervisor Waggle booked for this VM
	// and the one the provisioning pipeline is obliged to build on. It is part
	// of that pipeline's Terraform state, so it is never rewritten server-side
	// -- not even when the guest turns up somewhere else.
	HypervisorID uuid.UUID `gorm:"type:uuid;index;not null" json:"hypervisor_id"`
	// ObservedHypervisorID is where discovery last actually found the guest.
	// nil means "on its assignment, or not yet observed". A non-nil value that
	// differs from HypervisorID is a violated booking, and capacity is charged
	// against it rather than the assignment so the scheduler stops overselling
	// the host really carrying the guest. The assignment stays untouched as the
	// record of what was promised.
	ObservedHypervisorID *uuid.UUID `gorm:"type:uuid;index;column:observed_hypervisor_id" json:"observed_hypervisor_id,omitempty"`
	VMID                 *int       `gorm:"column:vmid;index" json:"vmid,omitempty"`
	CreatedAt            time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt            time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
