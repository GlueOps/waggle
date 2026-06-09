package control

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID      uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:ux_user_account_org" json:"account_id"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:ux_user_account_org;index" json:"organization_id"`
	DisplayName    string     `gorm:"size:255" json:"display_name"`
	// Role governs what the member may do in the org: owner (full control incl.
	// delete), admin (manage org + members), member (use the fleet).
	Role           string     `gorm:"size:16;not null;default:'member'" json:"role"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
