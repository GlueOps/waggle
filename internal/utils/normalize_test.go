package utils

import (
	"strings"
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"lowercases", "User@Example.COM", "user@example.com"},
		{"trims surrounding space", "  user@example.com \t", "user@example.com"},
		{"preserves +tag", "User+Tag@Example.com", "user+tag@example.com"},
		{"already normal", "user@example.com", "user@example.com"},
		{"empty", "", ""},
		{"only whitespace", "   ", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NormalizeEmail(c.in); got != c.want {
				t.Fatalf("NormalizeEmail(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestNormalizeDomain(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"lowercases and trims", "  Example.COM ", "example.com"},
		{"strips leading @", "@example.com", "example.com"},
		{"strips leading dot", ".example.com", "example.com"},
		{"strips trailing dot (FQDN root)", "example.com.", "example.com"},
		{"strips both ends", "  @Example.com.  ", "example.com"},
		{"empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NormalizeDomain(c.in); got != c.want {
				t.Fatalf("NormalizeDomain(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// NormalizeDomain only strips ONE leading "@" / "." because it uses TrimPrefix
// rather than TrimLeft. Pinning that so a future refactor to TrimLeft is a
// deliberate decision rather than a silent behaviour change.
func TestNormalizeDomainStripsOnlyOneLeadingMarker(t *testing.T) {
	if got := NormalizeDomain("..example.com"); got != ".example.com" {
		t.Fatalf("NormalizeDomain(\"..example.com\") = %q, want %q", got, ".example.com")
	}
}

func TestExtractDomain(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"simple", "user@example.com", "example.com"},
		{"normalizes without pre-normalized input", "  User@Example.COM ", "example.com"},
		{"uses last @ so local-part @ is tolerated", "we@ird@example.com", "example.com"},
		{"strips trailing dot", "user@example.com.", "example.com"},
		{"subdomain preserved", "user@mail.example.com", "mail.example.com"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ExtractDomain(c.in)
			if err != nil {
				t.Fatalf("ExtractDomain(%q) returned error: %v", c.in, err)
			}
			if got != c.want {
				t.Fatalf("ExtractDomain(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestExtractDomainErrors(t *testing.T) {
	cases := []struct {
		name string
		in   string
	}{
		{"no at sign", "not-an-email"},
		{"empty", ""},
		{"empty domain", "user@"},
		{"domain is only a dot", "user@."},
		{"domain is only an at", "user@@"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := ExtractDomain(c.in)
			if err == nil {
				t.Fatalf("ExtractDomain(%q) = %q, want an error", c.in, got)
			}
			if got != "" {
				t.Fatalf("ExtractDomain(%q) returned %q alongside an error, want empty", c.in, got)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	cases := []struct {
		name       string
		in         string
		wantPrefix string
	}{
		{"lowercases and dashes spaces", "My Great Org", "my-great-org-"},
		{"collapses runs of separators", "My   Great!!! Org", "my-great-org-"},
		{"strips leading/trailing separators", "  ...Org...  ", "org-"},
		{"keeps digits", "Org 42", "org-42-"},
		{"drops non-ascii", "Café Org", "caf-org-"},
		{"falls back to 'org' when nothing survives", "!!!", "org-"},
		{"empty falls back to 'org'", "", "org-"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Slugify(c.in)
			if !strings.HasPrefix(got, c.wantPrefix) {
				t.Fatalf("Slugify(%q) = %q, want prefix %q", c.in, got, c.wantPrefix)
			}
			suffix := strings.TrimPrefix(got, c.wantPrefix)
			if suffix == "" {
				t.Fatalf("Slugify(%q) = %q, expected a random suffix after %q", c.in, got, c.wantPrefix)
			}
			for _, r := range suffix {
				if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyz234567", r) {
					t.Fatalf("Slugify(%q) suffix %q contains non-base32 rune %q", c.in, suffix, r)
				}
			}
		})
	}
}

// The random suffix exists specifically so two signups with the same org name
// don't collide.
func TestSlugifyIsCollisionResistant(t *testing.T) {
	const n = 200
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		s := Slugify("Acme Corp")
		if _, dup := seen[s]; dup {
			t.Fatalf("Slugify produced duplicate slug %q within %d calls", s, n)
		}
		seen[s] = struct{}{}
	}
}

func TestSlugifyProducesURLSafeOutput(t *testing.T) {
	const allowed = "abcdefghijklmnopqrstuvwxyz0123456789-"
	for _, in := range []string{
		"My Great Org", "Ünïcødé Nåme", "org/with/slashes", "a?b#c&d=e",
		"trailing---dashes---", "	tabs	and	newlines\n", "!!!", "",
	} {
		got := Slugify(in)
		for _, r := range got {
			if !strings.ContainsRune(allowed, r) {
				t.Fatalf("Slugify(%q) = %q contains URL-unsafe rune %q", in, got, r)
			}
		}
		if strings.HasPrefix(got, "-") || strings.HasSuffix(got, "-") {
			t.Fatalf("Slugify(%q) = %q should not start or end with a dash", in, got)
		}
		if strings.Contains(got, "--") {
			t.Fatalf("Slugify(%q) = %q should not contain a doubled dash", in, got)
		}
	}
}
