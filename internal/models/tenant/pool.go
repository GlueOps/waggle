package tenant

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Pool struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DatacenterID uuid.UUID      `gorm:"type:uuid;index;not null" json:"datacenter_id"`
	SlotID       uuid.UUID      `gorm:"type:uuid;index;not null" json:"slot_id"`
	Name         string         `gorm:"size:255;index;not null" json:"name"`
	DesiredCount int            `gorm:"not null;default:0" json:"desired_count"`
	Metadata     datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
