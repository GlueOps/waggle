package service

import (
	"errors"
	"testing"

	"github.com/glueops/waggle/internal/models/tenant"

	"github.com/google/uuid"
)

// capacityOf mirrors the CPU arithmetic plan() builds each hvCapacity from, so
// these tests exercise the same expression the scheduler uses without needing
// a database.
func capacityOf(h tenant.Hypervisor, consumedCPU int) int {
	return (h.EffectiveCPUTotal() - h.CPUReserved - h.CPUUsed) - consumedCPU
}

func TestOvercommitExpandsSchedulableCPU(t *testing.T) {
	// 32 physical cores, 2 held back for the host, 8 vCPU already allocated to
	// non-Waggle guests, 10 vCPU committed by Waggle's own ledger.
	hv := tenant.Hypervisor{CPUTotal: 32, CPUReserved: 2, CPUUsed: 8}

	cases := []struct {
		name  string
		ratio float64
		want  int
	}{
		// 32 - 2 - 8 - 10
		{"no overcommit", 1.0, 12},
		// floor(32*4) - 2 - 8 - 10 -- only the physical total scales; the
		// vCPU subtrahends come off unscaled.
		{"4:1", 4.0, 108},
		{"1.5:1", 1.5, 28},
		{"undercommit", 0.5, -4},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			hv.CPUOvercommitRatio = c.ratio
			if got := capacityOf(hv, 10); got != c.want {
				t.Fatalf("cpuRemaining = %d, want %d", got, c.want)
			}
		})
	}
}

// The point of the feature: a slot that does not fit at 1:1 fits once the node
// is oversold, and RAM/disk remain hard limits that overcommit cannot relax.
func TestPlanPlacementsRespectsOvercommittedCPU(t *testing.T) {
	id := uuid.New()
	// 8 cores, nothing reserved or used. Ten 2-vCPU VMs need 20 vCPU.
	hv := tenant.Hypervisor{CPUTotal: 8, RAMGBTotal: 1024, DiskGBTotal: 10240}
	cost := slotCost{vcpu: 2, ramGB: 1, diskGB: 1}

	mk := func(ratio float64) []hvCapacity {
		hv.CPUOvercommitRatio = ratio
		return []hvCapacity{{
			id:            id,
			cpuRemaining:  capacityOf(hv, 0),
			ramRemaining:  hv.RAMGBTotal,
			diskRemaining: hv.DiskGBTotal,
		}}
	}

	if _, fit, ok := planPlacements(mk(1.0), cost, 10); ok || fit != 4 {
		t.Fatalf("at 1:1 got (fit=%d, ok=%v), want (4, false) — 8 vCPU fits four 2-vCPU VMs", fit, ok)
	}
	if ids, fit, ok := planPlacements(mk(4.0), cost, 10); !ok || fit != 10 || len(ids) != 10 {
		t.Fatalf("at 4:1 got (fit=%d, ok=%v, ids=%d), want all 10 placed", fit, ok, len(ids))
	}

	// RAM is not overcommittable: a node oversold on CPU still refuses VMs it
	// cannot back with real memory.
	starved := mk(64.0)
	starved[0].ramRemaining = 3
	if _, fit, ok := planPlacements(starved, cost, 10); ok || fit != 3 {
		t.Fatalf("got (fit=%d, ok=%v), want (3, false) — RAM must still bind", fit, ok)
	}
}

func TestHypervisorInputValidatesOvercommitRatio(t *testing.T) {
	ratio := func(f float64) *float64 { return &f }
	base := func() HypervisorInput {
		return HypervisorInput{Name: "n1", CPUTotal: 8, RAMGBTotal: 16, DiskGBTotal: 100}
	}

	cases := []struct {
		name    string
		ratio   *float64
		wantErr bool
	}{
		{"omitted inherits the datacenter default", nil, false},
		{"1.0 accepted", ratio(1.0), false},
		{"4.0 accepted", ratio(4.0), false},
		{"undercommit accepted", ratio(0.25), false},
		{"at the cap accepted", ratio(MaxCPUOvercommitRatio), false},
		{"zero rejected", ratio(0), true},
		{"negative rejected", ratio(-2), true},
		{"above the cap rejected", ratio(MaxCPUOvercommitRatio + 0.01), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := base()
			in.CPUOvercommitRatio = c.ratio
			err := in.validate()
			if c.wantErr && !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("validate() = %v, want ErrInvalidInput", err)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("validate() = %v, want nil", err)
			}
		})
	}
}

// Headroom is measured against physical cores, so raising the ratio must not
// change whether a given cpu_reserved is legal.
func TestReservedIsValidatedAgainstPhysicalCores(t *testing.T) {
	ratio := 8.0
	in := HypervisorInput{
		Name: "n1", CPUTotal: 8, CPUReserved: 16,
		RAMGBTotal: 16, DiskGBTotal: 100,
		CPUOvercommitRatio: &ratio,
	}
	// 16 reserved > 8 physical cores, even though 8*8=64 vCPU would "fit" it.
	if err := in.validate(); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("validate() = %v, want ErrInvalidInput", err)
	}
}

func TestDatacenterInputValidatesOvercommitRatio(t *testing.T) {
	ratio := func(f float64) *float64 { return &f }
	cases := []struct {
		name    string
		ratio   *float64
		wantErr bool
	}{
		{"omitted defaults to 1.0", nil, false},
		{"4.0 accepted", ratio(4.0), false},
		{"zero rejected", ratio(0), true},
		{"above the cap rejected", ratio(MaxCPUOvercommitRatio + 1), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in := DatacenterInput{Name: "dc", URL: "https://pve.example", CPUOvercommitRatio: c.ratio}
			err := in.validate()
			if c.wantErr && !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("validate() = %v, want ErrInvalidInput", err)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("validate() = %v, want nil", err)
			}
		})
	}
}
