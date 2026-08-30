package proxmox

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

// newTestClient points a Client at a stub Proxmox serving the given handler.
func newTestClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c, err := New(Config{BaseURL: srv.URL, Token: "user@pve!id=secret"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c
}

const gib = int64(1024 * 1024 * 1024)

// TestNodeUsageReportsManagedResidency is the hv23 regression: a Waggle-managed
// guest that migrated onto this node must be reported in Managed even though it
// is excluded from the capacity sums. Without it discovery cannot tell that the
// guest's placement row names the wrong hypervisor, and the node is oversold by
// exactly that guest's size.
func TestNodeUsageReportsManagedResidency(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api2/json/nodes/hv23/qemu":
			_, _ = w.Write([]byte(`{"data":[
				{"vmid":103009401,"cpus":32,"maxmem":68719476736,"maxdisk":107374182400},
				{"vmid":102723902,"cpus":16,"maxmem":34359738368,"maxdisk":107374182400},
				{"vmid":999,"cpus":2,"maxmem":8589934592,"maxdisk":10737418240}
			]}`))
		case "/api2/json/nodes/hv23/lxc":
			_, _ = w.Write([]byte(`{"data":[]}`))
		default:
			http.NotFound(w, r)
		}
	})

	managed := map[int]struct{}{103009401: {}, 102723902: {}}
	usage, err := c.NodeUsage(context.Background(), "hv23", managed)
	if err != nil {
		t.Fatalf("NodeUsage: %v", err)
	}

	// Only the unmanaged guest counts toward usage.
	if usage.Guests != 1 || usage.VCPU != 2 || usage.MemBytes != 8*gib {
		t.Errorf("unmanaged sums = guests %d vcpu %d mem %d GiB; want 1/2/8",
			usage.Guests, usage.VCPU, usage.MemBytes/gib)
	}

	got := append([]int(nil), usage.Managed...)
	sort.Ints(got)
	want := []int{102723902, 103009401}
	if len(got) != len(want) {
		t.Fatalf("Managed = %v; want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("Managed = %v; want %v", got, want)
		}
	}
}

// TestNodeUsageErrorsWhenNoListingSucceeds guards the silent-zero footgun: if
// every guest endpoint fails, usage is UNKNOWN. Returning a zero NodeUsage with
// a nil error lets discovery overwrite a correct ram_gb_used with 0, which
// makes a full node look empty to the scheduler.
func TestNodeUsageErrorsWhenNoListingSucceeds(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "permission denied", http.StatusForbidden)
	})

	if _, err := c.NodeUsage(context.Background(), "hv23", nil); err == nil {
		t.Fatal("NodeUsage returned nil error when every guest listing failed")
	}
}

// TestNodeUsageToleratesOneKindFailing keeps the original best-effort
// behaviour: an LXC endpoint that 403s must not discard the QEMU sums.
func TestNodeUsageToleratesOneKindFailing(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api2/json/nodes/hv23/qemu" {
			_, _ = w.Write([]byte(`{"data":[{"vmid":101,"cpus":4,"maxmem":8589934592,"maxdisk":10737418240}]}`))
			return
		}
		http.Error(w, "nope", http.StatusForbidden)
	})

	usage, err := c.NodeUsage(context.Background(), "hv23", nil)
	if err != nil {
		t.Fatalf("NodeUsage: %v", err)
	}
	if usage.Guests != 1 || usage.MemBytes != 8*gib {
		t.Errorf("usage = guests %d mem %d GiB; want 1/8", usage.Guests, usage.MemBytes/gib)
	}
}
