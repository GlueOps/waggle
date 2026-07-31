package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/glueops/waggle/internal/utils"
)

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name    string
		pw      string
		wantErr bool
	}{
		{"too short", strings.Repeat("a", minPasswordLen-1), true},
		{"empty", "", true},
		{"at minimum", strings.Repeat("a", minPasswordLen), false},
		{"typical", "correct horse battery", false},
		{"at bcrypt maximum", strings.Repeat("a", maxPasswordLen), false},
		{"one over bcrypt maximum", strings.Repeat("a", maxPasswordLen+1), true},
		{"far over bcrypt maximum", strings.Repeat("a", 128), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validatePassword(c.pw)
			if c.wantErr && err == nil {
				t.Fatalf("validatePassword(%d chars) = nil, want an error", len(c.pw))
			}
			if !c.wantErr && err != nil {
				t.Fatalf("validatePassword(%d chars) = %v, want nil", len(c.pw), err)
			}
			if c.wantErr && !errors.Is(err, ErrInvalidInput) {
				t.Fatalf("validatePassword(%d chars) returned %v, which does not wrap ErrInvalidInput — "+
					"the API would map it to a 500 instead of a 4xx", len(c.pw), err)
			}
		})
	}
}

// Regression: the password bounds must match what bcrypt actually accepts.
// Previously only a minimum was enforced while the API allowed up to 128
// characters, so a 73-128 character password passed validation and then failed
// inside utils.HashPassword — surfacing to the user as a 500 rather than a 422.
func TestPasswordPolicyMatchesBcryptLimit(t *testing.T) {
	// Anything validatePassword accepts must be hashable.
	for _, n := range []int{minPasswordLen, 20, 71, maxPasswordLen} {
		pw := strings.Repeat("a", n)
		if err := validatePassword(pw); err != nil {
			t.Fatalf("validatePassword(%d chars) unexpectedly rejected: %v", n, err)
		}
		if _, err := utils.HashPassword(pw); err != nil {
			t.Fatalf("validatePassword accepted a %d-char password that HashPassword rejects: %v", n, err)
		}
	}

	// And the first length bcrypt rejects must already be rejected by policy,
	// so the error is ErrInvalidInput (4xx) rather than a hash failure (500).
	over := strings.Repeat("a", maxPasswordLen+1)
	if _, err := utils.HashPassword(over); err == nil {
		t.Fatalf("expected bcrypt to reject %d chars; maxPasswordLen may be stale", maxPasswordLen+1)
	}
	if err := validatePassword(over); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("validatePassword(%d chars) = %v, want ErrInvalidInput", maxPasswordLen+1, err)
	}
}
