package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Placement struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PoolID       uuid.UUID `gorm:"type:uuid;index;not null" json:"pool_id"`
	HypervisorID uuid.UUID `gorm:"type:uuid;index;not null" json:"hypervisor_id"`
	VMID         *int      `gorm:"column:vmid;index" json:"vmid,omitempty"`
	CreatedAt    time.Time `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
