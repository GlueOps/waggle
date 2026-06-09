package control

import (
	"time"

	"github.com/google/uuid"
)

// OrgAPIKey is a long-lived, organization-scoped credential used by automation
// (notably the Terraform provider) to authenticate as a tenant without a user
// login. The plaintext token is shown once at creation; only its SHA-256 hash
// is stored, so the column is a safe-to-leak lookup key rather than a secret.
//
// A key is valid when RevokedAt IS NULL and (ExpiresAt IS NULL OR ExpiresAt is
// in the future). CreatedByAccountID records the human who minted it (nullable
// so platform-issued keys can omit it).
type OrgAPIKey struct {
	ID                 uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OrganizationID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"organization_id"`
	Name               string     `gorm:"size:255;not null" json:"name"`
	TokenHash          string     `gorm:"size:64;uniqueIndex;not null" json:"-"`
	Prefix             string     `gorm:"size:16;not null" json:"prefix"`
	CreatedByAccountID *uuid.UUID `gorm:"type:uuid" json:"created_by_account_id,omitempty"`
	LastUsedAt         *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt          *time.Time `gorm:"index" json:"expires_at,omitempty"`
	RevokedAt          *time.Time `json:"revoked_at,omitempty"`
	CreatedAt          time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt          time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
