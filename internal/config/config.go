package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env                 string `mapstructure:"env" yaml:"env"`
	BindHost            string `mapstructure:"bind_host" yaml:"bind_host"`
	BindPort            string `mapstructure:"bind_port" yaml:"bind_port"`
	BaseURL             string `mapstructure:"base_url" yaml:"base_url"`
	BasePath            string `mapstructure:"base_path" yaml:"base_path"`
	FrontendMode        string `mapstructure:"frontend_mode" yaml:"frontend_mode"`
	ViteDevURL          string `mapstructure:"vite_dev_url" yaml:"vite_dev_url"`
	DatabaseURL         string `mapstructure:"database_url" yaml:"database_url"`
	AdminDatabaseURL    string `mapstructure:"admin_database_url" yaml:"admin_database_url"`
	RiverDatabaseURL    string `mapstructure:"river_database_url" yaml:"river_database_url"`
	EncryptionMasterKey string `mapstructure:"encryption_master_key" yaml:"encryption_master_key"`
	JWTSecret           string `mapstructure:"jwt_secret" yaml:"jwt_secret"`
	JWTIssuer           string `mapstructure:"jwt_issuer" yaml:"jwt_issuer"`
	JWTAudience         string `mapstructure:"jwt_audience" yaml:"jwt_audience"`
	JWTAccessTTLMin     int    `mapstructure:"jwt_access_ttl_min" yaml:"jwt_access_ttl_min"`
	JWTRefreshTTLHour   int    `mapstructure:"jwt_refresh_ttl_hour" yaml:"jwt_refresh_ttl_hour"`
	GoogleClientID      string `mapstructure:"google_client_id" yaml:"google_client_id"`
	GoogleClientSecret  string `mapstructure:"google_client_secret" yaml:"google_client_secret"`
	OAuthRedirectBase   string `mapstructure:"oauth_redirect_base" yaml:"oauth_redirect_base"`
	WebAuthnRPID        string `mapstructure:"webauthn_rp_id" yaml:"webauthn_rp_id"`
	WebAuthnRPOrigin    string `mapstructure:"webauthn_rp_origin" yaml:"webauthn_rp_origin"`
	WebAuthnRPName      string `mapstructure:"webauthn_rp_name" yaml:"webauthn_rp_name"`

	// SMTP transport for transactional email (verification, invites). When
	// SMTPServer is empty the app falls back to logging tokens (LogSender).
	SMTPServer   string `mapstructure:"smtp_server" yaml:"smtp_server"`
	SMTPPort     int    `mapstructure:"smtp_port" yaml:"smtp_port"`
	SMTPUser     string `mapstructure:"smtp_user" yaml:"smtp_user"`
	SMTPPassword string `mapstructure:"smtp_password" yaml:"smtp_password"`
	SMTPFrom     string `mapstructure:"smtp_from" yaml:"smtp_from"`

	// Default capacity held back from placement on newly-discovered hypervisors
	// (OS/host overhead). Applied only on first discovery of a node; operator
	// overrides via the API are preserved on re-discovery.
	ReserveCPU    int `mapstructure:"reserve_cpu" yaml:"reserve_cpu"`
	ReserveRAMGB  int `mapstructure:"reserve_ram_gb" yaml:"reserve_ram_gb"`
	ReserveDiskGB int `mapstructure:"reserve_disk_gb" yaml:"reserve_disk_gb"`
}

const (
	FrontendModeProxy = "proxy"
	FrontendModeEmbed = "embed"
	FrontendModeNone  = "none"
)

var envKeys = []string{
	"env",
	"bind_host",
	"bind_port",
	"base_url",
	"base_path",
	"frontend_mode",
	"vite_dev_url",
	"database_url",
	"admin_database_url",
	"river_database_url",
	"encryption_master_key",
	"jwt_secret",
	"jwt_issuer",
	"jwt_audience",
	"jwt_access_ttl_min",
	"jwt_refresh_ttl_hour",
	"google_client_id",
	"google_client_secret",
	"oauth_redirect_base",
	"webauthn_rp_id",
	"webauthn_rp_origin",
	"webauthn_rp_name",
	"smtp_server",
	"smtp_port",
	"smtp_user",
	"smtp_password",
	"smtp_from",
	"reserve_cpu",
	"reserve_ram_gb",
	"reserve_disk_gb",
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env")

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	for _, key := range envKeys {
		_ = v.BindEnv(key)
	}

	v.SetDefault("bind_host", "127.0.0.1")
	v.SetDefault("bind_port", "8080")
	v.SetDefault("base_path", "/api/v1")
	v.SetDefault("base_url", "http://localhost:8080")
	v.SetDefault("smtp_port", 1025)
	// Conservative default OS reservation: keep 2 GB RAM and 10 GB disk per
	// node for the hypervisor itself. CPU is overcommittable, so 0 by default.
	v.SetDefault("reserve_cpu", 0)
	v.SetDefault("reserve_ram_gb", 2)
	v.SetDefault("reserve_disk_gb", 10)

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var missing []string
	if strings.TrimSpace(c.DatabaseURL) == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if strings.TrimSpace(c.JWTSecret) == "" {
		missing = append(missing, "JWT_SECRET")
	}
	if strings.TrimSpace(c.BaseURL) == "" {
		missing = append(missing, "BASE_URL")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required config: %s", strings.Join(missing, ", "))
	}

	switch c.FrontendMode {
	case FrontendModeProxy, FrontendModeEmbed, FrontendModeNone, "":
	default:
		return fmt.Errorf("invalid FRONTEND_MODE %q (expected proxy|embed|none)", c.FrontendMode)
	}
	return nil
}
