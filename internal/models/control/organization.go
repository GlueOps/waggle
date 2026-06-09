package control

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Organization struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name               string         `gorm:"size:255;not null" json:"name"`
	Slug               string         `gorm:"size:255;uniqueIndex" json:"slug"`
	Domain             *string        `gorm:"size:255;uniqueIndex" json:"domain,omitempty"`
	Status             string         `gorm:"size:32;not null;default:'pending';index" json:"status"`
	ConnectionString   string         `gorm:"type:text" json:"-"`
	EncryptedTenantKey string         `gorm:"type:text" json:"-"`
	TenantKeyIV        string         `gorm:"type:text" json:"-"`
	TenantKeyTag       string         `gorm:"type:text" json:"-"`
	Metadata           datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
