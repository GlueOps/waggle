package control

import (
	"time"

	"github.com/google/uuid"
)

type TokenSession struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID        uuid.UUID  `gorm:"type:uuid;index;not null" json:"account_id"`
	OrganizationID   uuid.UUID  `gorm:"type:uuid;index" json:"organization_id"`
	RefreshTokenHash string     `gorm:"size:128;uniqueIndex;not null" json:"-"`
	UserAgent        string     `gorm:"type:text" json:"user_agent"`
	IPAddress        string     `gorm:"size:64" json:"ip_address"`
	ExpiresAt        time.Time  `gorm:"index" json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	CreatedAt        time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt        time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
