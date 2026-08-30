//go:build dbtest

package service

import (
	"context"
	"os"
	"testing"

	"github.com/glueops/waggle/internal/models/tenant"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// openTestDB connects to the throwaway Postgres named by WAGGLE_TEST_DSN. These
// tests exist because recordObservedResidency and consumedByHypervisor are SQL
// against a live schema: v0.2.6 shipped a Pluck into a bare uuid.UUID that
// compiled, passed vet, and failed every discovery sweep in production.
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("WAGGLE_TEST_DSN")
	if dsn == "" {
		t.Skip("WAGGLE_TEST_DSN not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	return db
}

// seed builds a datacenter with two hypervisors, a slot, a pool and one
// placement assigned to hvA, and returns the ids.
func seed(t *testing.T, db *gorm.DB) (dcID, hvA, hvB, poolID uuid.UUID) {
	t.Helper()
	dcID, hvA, hvB, poolID = uuid.New(), uuid.New(), uuid.New(), uuid.New()
	slotID := uuid.New()
	exec := func(q string, args ...any) {
		if err := db.Exec(q, args...).Error; err != nil {
			t.Fatalf("seed %q: %v", q, err)
		}
	}
	exec(`INSERT INTO datacenters (id,name,url) VALUES (?,?,?)`, dcID, "dc-"+dcID.String()[:8], "https://pve.example")
	exec(`INSERT INTO hypervisors (id,datacenter_id,name,cpu_total,ram_gb_total,disk_gb_total) VALUES (?,?,?,?,?,?)`,
		hvA, dcID, "hv-a-"+hvA.String()[:8], 80, 125, 1737)
	exec(`INSERT INTO hypervisors (id,datacenter_id,name,cpu_total,ram_gb_total,disk_gb_total) VALUES (?,?,?,?,?,?)`,
		hvB, dcID, "hv-b-"+hvB.String()[:8], 80, 125, 1737)
	exec(`INSERT INTO slots (id,name,vcpu,ram_gb,disk_gb) VALUES (?,?,?,?,?)`, slotID, "3xlarge-"+slotID.String()[:8], 32, 64, 100)
	exec(`INSERT INTO pools (id,datacenter_id,slot_id,name,desired_count) VALUES (?,?,?,?,?)`, poolID, dcID, slotID, "pool-"+poolID.String()[:8], 1)
	return
}

func hvName(t *testing.T, db *gorm.DB, id uuid.UUID) string {
	t.Helper()
	var names []string
	if err := db.Model(&tenant.Hypervisor{}).Where("id = ?", id).Pluck("name", &names).Error; err != nil || len(names) == 0 {
		t.Fatalf("hvName: %v", err)
	}
	return names[0]
}

// TestRecordObservedResidency is the v0.2.6 regression: a guest assigned to hvA
// but found on hvB must get observed_hypervisor_id = hvB, and the capacity must
// follow it.
func TestRecordObservedResidency(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	dcID, hvA, hvB, poolID := seed(t, db)

	vmid := 103009401
	plID := uuid.New()
	if err := db.Exec(`INSERT INTO placements (id,pool_id,hypervisor_id,vmid) VALUES (?,?,?,?)`,
		plID, poolID, hvA, vmid).Error; err != nil {
		t.Fatalf("insert placement: %v", err)
	}

	// Before: the whole 64 GB is charged to hvA, nothing to hvB.
	got, err := consumedByHypervisor(db, []uuid.UUID{hvA, hvB})
	if err != nil {
		t.Fatalf("consumedByHypervisor: %v", err)
	}
	if got[hvA].RAM != 64 || got[hvB].RAM != 0 {
		t.Fatalf("before: hvA=%d hvB=%d; want 64/0", got[hvA].RAM, got[hvB].RAM)
	}

	// Discovery finds the guest on hvB instead.
	if err := recordObservedResidency(ctx, db, dcID, map[string][]int{hvName(t, db, hvB): {vmid}}); err != nil {
		t.Fatalf("recordObservedResidency: %v", err)
	}

	var observed []uuid.UUID
	if err := db.Model(&tenant.Placement{}).Where("id = ?", plID).
		Pluck("observed_hypervisor_id", &observed).Error; err != nil {
		t.Fatalf("read observed: %v", err)
	}
	if len(observed) != 1 || observed[0] != hvB {
		t.Fatalf("observed_hypervisor_id = %v; want %s", observed, hvB)
	}

	// The assignment must be untouched -- it is the pipeline's Terraform state.
	var assigned []uuid.UUID
	if err := db.Model(&tenant.Placement{}).Where("id = ?", plID).
		Pluck("hypervisor_id", &assigned).Error; err != nil {
		t.Fatalf("read assigned: %v", err)
	}
	if assigned[0] != hvA {
		t.Fatalf("hypervisor_id was rewritten to %s; must stay %s", assigned[0], hvA)
	}

	// After: capacity follows the guest.
	got, err = consumedByHypervisor(db, []uuid.UUID{hvA, hvB})
	if err != nil {
		t.Fatalf("consumedByHypervisor: %v", err)
	}
	if got[hvA].RAM != 0 || got[hvB].RAM != 64 {
		t.Fatalf("after: hvA=%d hvB=%d; want 0/64", got[hvA].RAM, got[hvB].RAM)
	}
	if got[hvB].CPU != 32 || got[hvB].Disk != 100 {
		t.Errorf("after: hvB cpu=%d disk=%d; want 32/100", got[hvB].CPU, got[hvB].Disk)
	}
}

// TestRecordObservedResidencyClearsWhenHome verifies a stale violation is
// dropped once the guest is seen on its assignment again.
func TestRecordObservedResidencyClearsWhenHome(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	dcID, hvA, hvB, poolID := seed(t, db)

	vmid := 103009400
	plID := uuid.New()
	if err := db.Exec(`INSERT INTO placements (id,pool_id,hypervisor_id,vmid,observed_hypervisor_id) VALUES (?,?,?,?,?)`,
		plID, poolID, hvA, vmid, hvB).Error; err != nil {
		t.Fatalf("insert placement: %v", err)
	}

	if err := recordObservedResidency(ctx, db, dcID, map[string][]int{hvName(t, db, hvA): {vmid}}); err != nil {
		t.Fatalf("recordObservedResidency: %v", err)
	}

	// Assert on SQL NULL directly: plucking a NULL uuid into *uuid.UUID yields
	// a pointer to the zero UUID, not nil, so a Go-side nil check would pass
	// even against a row holding all-zeroes.
	var nulls int64
	if err := db.Model(&tenant.Placement{}).
		Where("id = ? AND observed_hypervisor_id IS NULL", plID).
		Count(&nulls).Error; err != nil {
		t.Fatalf("count nulls: %v", err)
	}
	if nulls != 1 {
		var raw []string
		_ = db.Raw(`SELECT COALESCE(observed_hypervisor_id::text,'<null>') FROM placements WHERE id = ?`, plID).Scan(&raw).Error
		t.Fatalf("observed_hypervisor_id = %v; want SQL NULL once the guest is home", raw)
	}
}
