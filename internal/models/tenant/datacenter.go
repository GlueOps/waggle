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
	InsecureSkipVerify bool `gorm:"column:insecure_skip_verify;not null;default:false" json:"insecure_skip_verify"`
	// CPUOvercommitRatio is this datacenter's default vCPU-per-physical-core
	// ratio, stamped onto each hypervisor the first time discovery sees it.
	// Changing it does NOT retroactively re-rate existing nodes — they keep
	// whatever ratio they were stamped with (or an operator later set), the
	// same way Reserve* defaults behave. Per-tenant by construction: this row
	// lives in the tenant DB, so one org's overcommit policy can never move
	// another's capacity.
	CPUOvercommitRatio float64   `gorm:"column:cpu_overcommit_ratio;type:numeric(5,2);not null;default:1.0" json:"cpu_overcommit_ratio"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// EffectiveCPUOvercommitRatio is the ratio to stamp onto hypervisors in this
// datacenter. A non-positive stored value (a row predating the column, or bad
// data) is read as DefaultCPUOvercommitRatio, so nodes inherit no-overcommit
// rather than a zero that would make them unschedulable.
func (d Datacenter) EffectiveCPUOvercommitRatio() float64 {
	if d.CPUOvercommitRatio <= 0 {
		return DefaultCPUOvercommitRatio
	}
	return d.CPUOvercommitRatio
}
