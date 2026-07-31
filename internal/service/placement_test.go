package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

// hv builds a hypervisor capacity entry with a stable, readable id so failures
// name the host rather than a random UUID.
func hv(t *testing.T, name string, cpu, ram, disk, poolCount int) hvCapacity {
	t.Helper()
	return hvCapacity{
		id:            uuid.NewSHA1(uuid.Nil, []byte(name)),
		cpuRemaining:  cpu,
		ramRemaining:  ram,
		diskRemaining: disk,
		poolCount:     poolCount,
	}
}

// names maps chosen ids back to readable host names for assertions.
func names(t *testing.T, ids []uuid.UUID, hvs []hvCapacity, labels []string) []string {
	t.Helper()
	byID := make(map[uuid.UUID]string, len(hvs))
	for i, h := range hvs {
		byID[h.id] = labels[i]
	}
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		n, ok := byID[id]
		if !ok {
			t.Fatalf("planPlacements returned unknown hypervisor id %s", id)
		}
		out = append(out, n)
	}
	return out
}

func TestHVCapacityFits(t *testing.T) {
	h := hvCapacity{cpuRemaining: 4, ramRemaining: 8, diskRemaining: 100}
	cases := []struct {
		name string
		cost slotCost
		want bool
	}{
		{"exact fit", slotCost{4, 8, 100}, true},
		{"comfortably under", slotCost{1, 1, 1}, true},
		{"zero cost", slotCost{0, 0, 0}, true},
		{"cpu over", slotCost{5, 8, 100}, false},
		{"ram over", slotCost{4, 9, 100}, false},
		{"disk over", slotCost{4, 8, 101}, false},
		{"all over", slotCost{99, 99, 999}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := h.fits(c.cost); got != c.want {
				t.Fatalf("fits(%+v) = %v, want %v", c.cost, got, c.want)
			}
		})
	}
}

func TestPlanPlacementsSpreadsAcrossHypervisors(t *testing.T) {
	labels := []string{"a", "b", "c"}
	hvs := []hvCapacity{
		hv(t, "a", 16, 64, 1000, 0),
		hv(t, "b", 16, 64, 1000, 0),
		hv(t, "c", 16, 64, 1000, 0),
	}
	cost := slotCost{vcpu: 2, ramGB: 4, diskGB: 50}

	ids, fit, ok := planPlacements(hvs, cost, 3)
	if !ok || fit != 3 {
		t.Fatalf("planPlacements = (fit %d, ok %v), want (3, true)", fit, ok)
	}

	got := names(t, ids, hvs, labels)
	seen := map[string]int{}
	for _, n := range got {
		seen[n]++
	}
	if len(seen) != 3 {
		t.Fatalf("3 VMs across 3 idle equal hosts landed on %d distinct hosts (%v); anti-affinity should spread them", len(seen), got)
	}
}

// The primary sort key is "fewest VMs already from this pool", so an
// already-loaded host must be skipped while an emptier one has room.
func TestPlanPlacementsPrefersFewestPoolVMs(t *testing.T) {
	labels := []string{"loaded", "empty"}
	hvs := []hvCapacity{
		hv(t, "loaded", 64, 256, 4000, 3), // far more spare capacity...
		hv(t, "empty", 8, 16, 200, 0),     // ...but this one has none of this pool
	}
	cost := slotCost{vcpu: 2, ramGB: 4, diskGB: 50}

	ids, fit, ok := planPlacements(hvs, cost, 1)
	if !ok || fit != 1 {
		t.Fatalf("planPlacements = (fit %d, ok %v), want (1, true)", fit, ok)
	}
	if got := names(t, ids, hvs, labels)[0]; got != "empty" {
		t.Fatalf("placed on %q; anti-affinity should win over spare capacity and pick %q", got, "empty")
	}
}

// Ties on poolCount are broken by most remaining CPU, to keep hosts balanced.
func TestPlanPlacementsTieBreaksOnMostFreeCPU(t *testing.T) {
	labels := []string{"small", "big"}
	hvs := []hvCapacity{
		hv(t, "small", 8, 64, 1000, 0),
		hv(t, "big", 32, 64, 1000, 0),
	}
	cost := slotCost{vcpu: 2, ramGB: 4, diskGB: 50}

	ids, _, ok := planPlacements(hvs, cost, 1)
	if !ok {
		t.Fatal("planPlacements failed with ample capacity")
	}
	if got := names(t, ids, hvs, labels)[0]; got != "big" {
		t.Fatalf("tie on poolCount placed on %q, want %q (most remaining CPU)", got, "big")
	}
}

func TestPlanPlacementsReportsPartialFit(t *testing.T) {
	hvs := []hvCapacity{
		hv(t, "a", 4, 64, 1000, 0), // room for exactly 2 VMs at 2 vcpu
		hv(t, "b", 2, 64, 1000, 0), // room for exactly 1
	}
	cost := slotCost{vcpu: 2, ramGB: 4, diskGB: 50}

	ids, fit, ok := planPlacements(hvs, cost, 10)
	if ok {
		t.Fatal("planPlacements reported success for an unsatisfiable request")
	}
	if fit != 3 {
		t.Fatalf("fit = %d, want 3 (2 on 'a' + 1 on 'b')", fit)
	}
	if len(ids) != fit {
		t.Fatalf("returned %d ids but reported fit %d; they must agree", len(ids), fit)
	}
}

// Placement is all-or-nothing at the caller level, so a failed plan must not
// have mutated the caller's slice — otherwise a retry would see phantom usage.
func TestPlanPlacementsDoesNotMutateInput(t *testing.T) {
	hvs := []hvCapacity{
		hv(t, "a", 8, 32, 500, 1),
		hv(t, "b", 8, 32, 500, 2),
	}
	before := make([]hvCapacity, len(hvs))
	copy(before, hvs)

	planPlacements(hvs, slotCost{2, 4, 50}, 100)

	for i := range hvs {
		if hvs[i] != before[i] {
			t.Fatalf("planPlacements mutated caller input at %d:\n got %+v\nwant %+v", i, hvs[i], before[i])
		}
	}
}

func TestPlanPlacementsEdgeCases(t *testing.T) {
	cost := slotCost{vcpu: 2, ramGB: 4, diskGB: 50}

	t.Run("zero count succeeds with no placements", func(t *testing.T) {
		ids, fit, ok := planPlacements([]hvCapacity{hv(t, "a", 8, 32, 500, 0)}, cost, 0)
		if !ok || fit != 0 || len(ids) != 0 {
			t.Fatalf("planPlacements(count=0) = (%d ids, fit %d, ok %v), want (0, 0, true)", len(ids), fit, ok)
		}
	})

	t.Run("no hypervisors fails immediately", func(t *testing.T) {
		ids, fit, ok := planPlacements(nil, cost, 1)
		if ok || fit != 0 || len(ids) != 0 {
			t.Fatalf("planPlacements(no hosts) = (%d ids, fit %d, ok %v), want (0, 0, false)", len(ids), fit, ok)
		}
	})

	t.Run("single host absorbs many VMs", func(t *testing.T) {
		hvs := []hvCapacity{hv(t, "solo", 100, 400, 5000, 0)}
		ids, fit, ok := planPlacements(hvs, cost, 5)
		if !ok || fit != 5 || len(ids) != 5 {
			t.Fatalf("planPlacements = (%d ids, fit %d, ok %v), want (5, 5, true)", len(ids), fit, ok)
		}
		for _, id := range ids {
			if id != hvs[0].id {
				t.Fatal("placements landed somewhere other than the only host")
			}
		}
	})

	t.Run("host that fits nothing is skipped", func(t *testing.T) {
		labels := []string{"tiny", "usable"}
		hvs := []hvCapacity{
			hv(t, "tiny", 1, 1, 1, 0),
			hv(t, "usable", 8, 32, 500, 5),
		}
		ids, fit, ok := planPlacements(hvs, cost, 1)
		if !ok || fit != 1 {
			t.Fatalf("planPlacements = (fit %d, ok %v), want (1, true)", fit, ok)
		}
		if got := names(t, ids, hvs, labels)[0]; got != "usable" {
			t.Fatalf("placed on %q, want %q — the undersized host must be skipped", got, "usable")
		}
	})
}

// Capacity must be decremented as VMs are assigned, or a host would be
// oversubscribed within a single plan.
func TestPlanPlacementsDecrementsCapacityAsItGoes(t *testing.T) {
	hvs := []hvCapacity{hv(t, "solo", 4, 8, 100, 0)}
	cost := slotCost{vcpu: 2, ramGB: 4, diskGB: 50}

	// Exactly 2 fit; a 3rd must not.
	if _, fit, ok := planPlacements(hvs, cost, 2); !ok || fit != 2 {
		t.Fatalf("planPlacements(2) = (fit %d, ok %v), want (2, true)", fit, ok)
	}
	if _, fit, ok := planPlacements(hvs, cost, 3); ok || fit != 2 {
		t.Fatalf("planPlacements(3) = (fit %d, ok %v), want (2, false) — capacity is not being consumed", fit, ok)
	}
}

// Whichever dimension is scarcest must bound the plan.
func TestPlanPlacementsRespectsEachDimension(t *testing.T) {
	cases := []struct {
		name    string
		host    hvCapacity
		wantFit int
	}{
		{"cpu bound", hvCapacity{cpuRemaining: 4, ramRemaining: 999, diskRemaining: 9999}, 2},
		{"ram bound", hvCapacity{cpuRemaining: 999, ramRemaining: 12, diskRemaining: 9999}, 3},
		{"disk bound", hvCapacity{cpuRemaining: 999, ramRemaining: 999, diskRemaining: 200}, 4},
	}
	cost := slotCost{vcpu: 2, ramGB: 4, diskGB: 50}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.host.id = uuid.New()
			_, fit, ok := planPlacements([]hvCapacity{c.host}, cost, 100)
			if ok {
				t.Fatal("expected the request to exceed capacity")
			}
			if fit != c.wantFit {
				t.Fatalf("fit = %d, want %d", fit, c.wantFit)
			}
		})
	}
}

func TestCapacityErrorMessage(t *testing.T) {
	err := &CapacityError{Requested: 10, Fit: 3}
	want := "insufficient capacity: requested 10, can place 3"
	if err.Error() != want {
		t.Fatalf("Error() = %q, want %q", err.Error(), want)
	}
	// Callers match on the concrete type to read Fit back out.
	var ce *CapacityError
	if !errors.As(fmt.Errorf("wrapped: %w", err), &ce) {
		t.Fatal("CapacityError does not survive errors.As through a wrap")
	}
	if ce.Fit != 3 || ce.Requested != 10 {
		t.Fatalf("unwrapped CapacityError = %+v, want Requested 10 / Fit 3", ce)
	}
}

func TestPoolInputValidate(t *testing.T) {
	valid := PoolInput{
		DatacenterID: uuid.New(),
		SlotID:       uuid.New(),
		Name:         "web",
		DesiredCount: 3,
	}
	if err := valid.validate(); err != nil {
		t.Fatalf("valid input rejected: %v", err)
	}

	t.Run("zero desired count is allowed", func(t *testing.T) {
		in := valid
		in.DesiredCount = 0
		if err := in.validate(); err != nil {
			t.Fatalf("DesiredCount 0 should be allowed (scale to zero): %v", err)
		}
	})

	bad := []struct {
		name   string
		mutate func(*PoolInput)
	}{
		{"empty name", func(p *PoolInput) { p.Name = "" }},
		{"nil datacenter", func(p *PoolInput) { p.DatacenterID = uuid.Nil }},
		{"nil slot", func(p *PoolInput) { p.SlotID = uuid.Nil }},
		{"negative count", func(p *PoolInput) { p.DesiredCount = -1 }},
	}
	for _, c := range bad {
		t.Run(c.name, func(t *testing.T) {
			in := valid
			c.mutate(&in)
			err := in.validate()
			if err == nil {
				t.Fatal("validate() = nil, want an error")
			}
			if !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("validate() = %v, want ErrInvalidInput so the API maps it to 422", err)
			}
		})
	}
}
