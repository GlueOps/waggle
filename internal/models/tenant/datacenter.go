package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Datacenter struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name              string    `gorm:"size:255" json:"name"`
	Url               string    `gorm:"size:255" json:"url"`
	EncryptedTokenKey string    `gorm:"type:text" json:"-"`
	TokenKeyIV        string    `gorm:"type:text" json:"-"`
	TokenKeyTag       string    `gorm:"type:text" json:"-"`
	// InsecureSkipVerify disables TLS certificate verification when talking to
	// this datacenter's Proxmox API. Common for self-signed homelab clusters.
	// Prefer a valid cert; this is an explicit per-datacenter opt-out.
	InsecureSkipVerify bool      `gorm:"column:insecure_skip_verify;not null;default:false" json:"insecure_skip_verify"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
