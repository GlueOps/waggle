package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Hypervisor struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	DatacenterID   uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:ux_hypervisor_datacenter_name" json:"datacenter_id"`
	Name           string     `gorm:"size:255;not null;uniqueIndex:ux_hypervisor_datacenter_name" json:"name"`
	CPUTotal       int        `gorm:"column:cpu_total;not null;default:0" json:"cpu_total"`
	CPUReserved    int        `gorm:"column:cpu_reserved;not null;default:0" json:"cpu_reserved"`
	RAMGBTotal     int        `gorm:"column:ram_gb_total;not null;default:0" json:"ram_gb_total"`
	RAMGBReserved  int        `gorm:"column:ram_gb_reserved;not null;default:0" json:"ram_gb_reserved"`
	DiskGBTotal    int        `gorm:"column:disk_gb_total;not null;default:0" json:"disk_gb_total"`
	DiskGBReserved int        `gorm:"column:disk_gb_reserved;not null;default:0" json:"disk_gb_reserved"`
	// *Used is capacity already committed to existing guests on the node,
	// populated by discovery (sum of each VM/container's allocated resources).
	// Distinct from *Reserved (operator headroom); both reduce bookable space.
	CPUUsed    int `gorm:"column:cpu_used;not null;default:0" json:"cpu_used"`
	RAMGBUsed  int `gorm:"column:ram_gb_used;not null;default:0" json:"ram_gb_used"`
	DiskGBUsed int `gorm:"column:disk_gb_used;not null;default:0" json:"disk_gb_used"`
	// Schedulable gates whether placement may book VMs onto this hypervisor.
	// Defaults true; operators flip it false to drain/exclude a node (e.g. for
	// maintenance) without deleting it. Discovery preserves the existing value.
	Schedulable    bool       `gorm:"column:schedulable;not null;default:true" json:"schedulable"`
	LastSyncedAt   *time.Time `gorm:"type:timestamptz;column:last_synced_at" json:"last_synced_at,omitempty"`
	CreatedAt      time.Time  `gorm:"type:timestamptz;column:created_at;not null;default:now()" json:"created_at" format:"date-time"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;autoUpdateTime;column:updated_at;not null;default:now()" json:"updated_at" format:"date-time"`
}
