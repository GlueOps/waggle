package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Slot struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"size:255;uniqueIndex;not null" json:"name"`
	VCPU      int       `gorm:"column:vcpu;not null" json:"vcpu"`
	RAMGB     int       `gorm:"column:ram_gb;not null" json:"ram_gb"`
	DiskGB    int       `gorm:"column:disk_gb;not null" json:"disk_gb"`
	CreatedAt time.Time `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt time.Time `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
