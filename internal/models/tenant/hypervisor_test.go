package tenant

import "testing"

func TestHypervisorEffectiveCPUTotal(t *testing.T) {
	cases := []struct {
		name  string
		cores int
		ratio float64
		want  int
	}{
		{"unset ratio degrades to 1:1", 32, 0, 32},
		{"negative ratio degrades to 1:1", 32, -4, 32},
		{"explicit 1.0 is a no-op", 32, 1.0, 32},
		{"4:1 quadruples the pool", 32, 4.0, 128},
		{"fractional ratio rounds down", 32, 1.5, 48},
		{"undercommit is allowed", 32, 0.5, 16},
		{"rounds down, never up", 7, 1.5, 10}, // 10.5 -> 10
		{"zero cores stay zero", 0, 8.0, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := Hypervisor{CPUTotal: c.cores, CPUOvercommitRatio: c.ratio}
			if got := h.EffectiveCPUTotal(); got != c.want {
				t.Fatalf("EffectiveCPUTotal() = %d, want %d", got, c.want)
			}
		})
	}
}

// A node that predates the column reads as ratio 0 in Go. It must schedule as
// 1:1 rather than reporting zero capacity, which would silently drain it.
func TestZeroRatioDoesNotZeroCapacity(t *testing.T) {
	h := Hypervisor{CPUTotal: 64}
	if got := h.EffectiveCPUOvercommitRatio(); got != DefaultCPUOvercommitRatio {
		t.Fatalf("ratio = %v, want %v", got, DefaultCPUOvercommitRatio)
	}
	if got := h.EffectiveCPUTotal(); got != 64 {
		t.Fatalf("EffectiveCPUTotal() = %d, want 64", got)
	}
}

func TestDatacenterEffectiveCPUOvercommitRatio(t *testing.T) {
	cases := []struct {
		name  string
		ratio float64
		want  float64
	}{
		{"unset degrades to 1.0", 0, DefaultCPUOvercommitRatio},
		{"negative degrades to 1.0", -1, DefaultCPUOvercommitRatio},
		{"explicit value passes through", 4.0, 4.0},
		{"undercommit passes through", 0.5, 0.5},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			d := Datacenter{CPUOvercommitRatio: c.ratio}
			if got := d.EffectiveCPUOvercommitRatio(); got != c.want {
				t.Fatalf("EffectiveCPUOvercommitRatio() = %v, want %v", got, c.want)
			}
		})
	}
}
