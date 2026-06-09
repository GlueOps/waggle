package control

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuthAuditEvent struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationID *uuid.UUID     `gorm:"type:uuid;index" json:"organization_id,omitempty"`
	UserID         *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Event          string         `gorm:"size:64;index;not null" json:"event"`
	Outcome        string         `gorm:"size:32;not null" json:"outcome"`
	IPAddress      string         `gorm:"size:64" json:"ip_address"`
	UserAgent      string         `gorm:"type:text" json:"user_agent"`
	Metadata       datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time      `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt      time.Time      `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
