package tenant

import (
	"math"
	"time"

	"github.com/google/uuid"
)

// DefaultCPUOvercommitRatio sells physical cores 1:1 — the no-overcommit
// behaviour Waggle had before the ratio existed.
const DefaultCPUOvercommitRatio = 1.0

type Hypervisor struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DatacenterID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:ux_hypervisor_datacenter_name" json:"datacenter_id"`
	Name           string    `gorm:"size:255;not null;uniqueIndex:ux_hypervisor_datacenter_name" json:"name"`
	CPUTotal       int       `gorm:"column:cpu_total;not null;default:0" json:"cpu_total"`
	CPUReserved    int       `gorm:"column:cpu_reserved;not null;default:0" json:"cpu_reserved"`
	RAMGBTotal     int       `gorm:"column:ram_gb_total;not null;default:0" json:"ram_gb_total"`
	RAMGBReserved  int       `gorm:"column:ram_gb_reserved;not null;default:0" json:"ram_gb_reserved"`
	DiskGBTotal    int       `gorm:"column:disk_gb_total;not null;default:0" json:"disk_gb_total"`
	DiskGBReserved int       `gorm:"column:disk_gb_reserved;not null;default:0" json:"disk_gb_reserved"`
	// *Used is capacity already committed to existing guests on the node,
	// populated by discovery (sum of each VM/container's allocated resources).
	// Distinct from *Reserved (operator headroom); both reduce bookable space.
	CPUUsed    int `gorm:"column:cpu_used;not null;default:0" json:"cpu_used"`
	RAMGBUsed  int `gorm:"column:ram_gb_used;not null;default:0" json:"ram_gb_used"`
	DiskGBUsed int `gorm:"column:disk_gb_used;not null;default:0" json:"disk_gb_used"`
	// CPUOvercommitRatio is how many vCPU the operator is willing to sell per
	// physical core on this node: 1.0 sells cores 1:1, 4.0 sells four vCPU per
	// core, 0.5 deliberately undersells. It scales CPUTotal only (see
	// EffectiveCPUTotal) — *Reserved/*Used/*Consumed are already vCPU counts.
	// Seeded from the parent Datacenter's ratio the first time discovery sees
	// the node (never from process config — that is shared across tenants);
	// operator overrides survive re-discovery. RAM and disk are not
	// overcommittable and have no analogue.
	CPUOvercommitRatio float64 `gorm:"column:cpu_overcommit_ratio;type:numeric(5,2);not null;default:1.0" json:"cpu_overcommit_ratio"`
	// Schedulable gates whether placement may book VMs onto this hypervisor.
	// Defaults true; operators flip it false to drain/exclude a node (e.g. for
	// maintenance) without deleting it. Discovery preserves the existing value.
	Schedulable  bool       `gorm:"column:schedulable;not null;default:true" json:"schedulable"`
	LastSyncedAt *time.Time `gorm:"type:timestamptz;column:last_synced_at" json:"last_synced_at,omitempty"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
	// *Consumed is the placement ledger's committed footprint on this
	// hypervisor (summed slot cost of its placements). Populated at read time,
	// never persisted (gorm:"-"). The API view subtracts it from bookable so
	// the figure matches what the scheduler will actually book. Distinct from
	// *Used, which is non-Waggle guest allocation discovered on the node.
	CPUConsumed    int `gorm:"-" json:"-"`
	RAMGBConsumed  int `gorm:"-" json:"-"`
	DiskGBConsumed int `gorm:"-" json:"-"`
}

// EffectiveCPUOvercommitRatio is the ratio to actually schedule against. A
// non-positive stored value (a row predating the column, or bad data) is read
// as DefaultCPUOvercommitRatio rather than honoured literally, so a missing or
// zeroed value degrades to no-overcommit instead of collapsing the node's
// capacity to zero and silently making it unschedulable.
func (h Hypervisor) EffectiveCPUOvercommitRatio() float64 {
	if h.CPUOvercommitRatio <= 0 {
		return DefaultCPUOvercommitRatio
	}
	return h.CPUOvercommitRatio
}

// EffectiveCPUTotal is the vCPU pool this node is willing to sell: physical
// cores scaled by the overcommit ratio, rounded down. It replaces CPUTotal as
// the starting point for bookable CPU; CPUReserved/CPUUsed/CPUConsumed are
// already vCPU counts and subtract from it unscaled.
//
// Shared by the scheduler (plan) and the hypervisor view so the two never
// disagree on how much CPU exists — the same reason consumedByHypervisor is
// shared for the committed side.
func (h Hypervisor) EffectiveCPUTotal() int {
	return int(math.Floor(float64(h.CPUTotal) * h.EffectiveCPUOvercommitRatio()))
}
