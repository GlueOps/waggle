package control

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DisplayName  string     `gorm:"size:255" json:"display_name"`
	PasswordHash string     `gorm:"type:text" json:"-"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
