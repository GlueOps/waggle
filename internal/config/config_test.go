package config

import (
	"strings"
	"testing"
	"time"
)

// validConfig is the minimum that Validate accepts; tests mutate a copy.
//
// Validate only checks these fields for non-emptiness, so the values are
// deliberately non-credential-shaped: no user:password in the database URL and
// no literal assigned to JWTSecret, both of which trip secret scanning.
func validConfig() *Config {
	return &Config{
		DatabaseURL: "postgres://localhost:5432/waggle_test?sslmode=disable",
		JWTSecret:   strings.Repeat("x", 32),
		BaseURL:     "http://localhost:8080",
	}
}

func TestValidateAcceptsMinimalConfig(t *testing.T) {
	if err := validConfig().Validate(); err != nil {
		t.Fatalf("Validate rejected a minimal valid config: %v", err)
	}
}

func TestValidateReportsMissingRequired(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*Config)
		wantKey string
	}{
		{"no database url", func(c *Config) { c.DatabaseURL = "" }, "DATABASE_URL"},
		{"blank database url", func(c *Config) { c.DatabaseURL = "   " }, "DATABASE_URL"},
		{"no jwt secret", func(c *Config) { c.JWTSecret = "" }, "JWT_SECRET"},
		{"blank jwt secret", func(c *Config) { c.JWTSecret = "\t" }, "JWT_SECRET"},
		{"no base url", func(c *Config) { c.BaseURL = "" }, "BASE_URL"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := validConfig()
			c.mutate(cfg)
			err := cfg.Validate()
			if err == nil {
				t.Fatal("Validate() = nil, want an error")
			}
			if !strings.Contains(err.Error(), c.wantKey) {
				t.Fatalf("error %q does not name the missing key %q", err, c.wantKey)
			}
		})
	}
}

// All missing keys should be reported at once so an operator fixes the whole
// config in one pass rather than one variable per restart.
func TestValidateReportsAllMissingKeysAtOnce(t *testing.T) {
	err := (&Config{}).Validate()
	if err == nil {
		t.Fatal("Validate() on an empty config = nil, want an error")
	}
	for _, key := range []string{"DATABASE_URL", "JWT_SECRET", "BASE_URL"} {
		if !strings.Contains(err.Error(), key) {
			t.Fatalf("error %q omits missing key %q", err, key)
		}
	}
}

func TestValidateFrontendMode(t *testing.T) {
	valid := []string{FrontendModeProxy, FrontendModeEmbed, FrontendModeNone, ""}
	for _, mode := range valid {
		t.Run("accepts "+mode, func(t *testing.T) {
			cfg := validConfig()
			cfg.FrontendMode = mode
			if err := cfg.Validate(); err != nil {
				t.Fatalf("Validate rejected FRONTEND_MODE %q: %v", mode, err)
			}
		})
	}

	for _, mode := range []string{"nope", "PROXY", "embedded", " proxy"} {
		t.Run("rejects "+mode, func(t *testing.T) {
			cfg := validConfig()
			cfg.FrontendMode = mode
			err := cfg.Validate()
			if err == nil {
				t.Fatalf("Validate accepted invalid FRONTEND_MODE %q", mode)
			}
			if !strings.Contains(err.Error(), "FRONTEND_MODE") {
				t.Fatalf("error %q does not mention FRONTEND_MODE", err)
			}
		})
	}
}

func TestDiscoveryIntervalDuration(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want time.Duration
	}{
		{"empty disables", "", 0},
		{"zero disables", "0", 0},
		{"whitespace disables", "   ", 0},
		{"seconds", "30s", 30 * time.Second},
		{"minutes", "5m", 5 * time.Minute},
		{"hours", "2h", 2 * time.Hour},
		{"compound", "1h30m", 90 * time.Minute},
		{"surrounding space tolerated", "  15m  ", 15 * time.Minute},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := validConfig()
			cfg.DiscoveryInterval = c.in
			got, err := cfg.DiscoveryIntervalDuration()
			if err != nil {
				t.Fatalf("DiscoveryIntervalDuration(%q): %v", c.in, err)
			}
			if got != c.want {
				t.Fatalf("DiscoveryIntervalDuration(%q) = %s, want %s", c.in, got, c.want)
			}
		})
	}
}

func TestDiscoveryIntervalDurationErrors(t *testing.T) {
	for _, in := range []string{"soon", "30", "5 minutes", "-1h", "-30s"} {
		t.Run(in, func(t *testing.T) {
			cfg := validConfig()
			cfg.DiscoveryInterval = in
			got, err := cfg.DiscoveryIntervalDuration()
			if err == nil {
				t.Fatalf("DiscoveryIntervalDuration(%q) = %s, want an error", in, got)
			}
			if !strings.Contains(err.Error(), "DISCOVERY_INTERVAL") {
				t.Fatalf("error %q does not name DISCOVERY_INTERVAL", err)
			}
		})
	}
}

// A bad interval must fail startup validation, not just the accessor —
// otherwise the daemon boots and silently never runs discovery.
func TestValidateRejectsBadDiscoveryInterval(t *testing.T) {
	cfg := validConfig()
	cfg.DiscoveryInterval = "not-a-duration"
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate accepted an unparseable DISCOVERY_INTERVAL")
	}

	cfg.DiscoveryInterval = "-5m"
	if err := cfg.Validate(); err == nil {
		t.Fatal("Validate accepted a negative DISCOVERY_INTERVAL")
	}
}
