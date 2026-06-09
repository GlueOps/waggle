package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/glueops/waggle/internal/database"
	"github.com/glueops/waggle/internal/service"
)

func TestMapFleetError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"org not found -> re-auth", fmt.Errorf("%w: some-id", database.ErrOrgNotFound), 401},
		{"tenant not provisioned", database.ErrTenantNotProvisioned, 503},
		{"tenant not active", fmt.Errorf("%w: org x status=pending", database.ErrTenantNotActive), 503},
		{"not found", service.ErrNotFound, 404},
		{"invalid input", service.ErrInvalidInput, 422},
		{"discovery", service.ErrDiscovery, 502},
		{"unmapped -> 500", errors.New("boom"), 500},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var se huma.StatusError
			if !errors.As(mapFleetError(c.err), &se) {
				t.Fatalf("mapFleetError did not return a huma.StatusError")
			}
			if se.GetStatus() != c.want {
				t.Fatalf("got status %d, want %d", se.GetStatus(), c.want)
			}
		})
	}
}
