package control

import (
	"time"

	"github.com/google/uuid"
)

type AccountEmail struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AccountID  uuid.UUID  `gorm:"type:uuid;not null;index;uniqueIndex:ux_account_emails_primary,where:is_primary" json:"account_id"`
	Email      string     `gorm:"size:255;not null;uniqueIndex" json:"email"`
	IsPrimary  bool       `gorm:"not null;default:false" json:"is_primary"`
	VerifiedAt *time.Time `gorm:"type:timestamptz" json:"verified_at,omitempty"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt  time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
