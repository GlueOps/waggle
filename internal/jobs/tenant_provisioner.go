package jobs

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/glueops/waggle/internal/migrations"
	"github.com/glueops/waggle/internal/models/control"
	"github.com/glueops/waggle/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TenantProvisionerService struct {
	controlDB    *gorm.DB
	controlDBURL string
	adminDBURL   string
	masterKey    []byte
}

func NewTenantProvisionerService(controlDB *gorm.DB, controlDBURL, adminOverride string, masterKey []byte) (*TenantProvisionerService, error) {
	adminURL, err := deriveAdminDBURL(controlDBURL, adminOverride)
	if err != nil {
		return nil, fmt.Errorf("resolve admin db url: %w", err)
	}
	return &TenantProvisionerService{
		controlDB:    controlDB,
		controlDBURL: controlDBURL,
		adminDBURL:   adminURL,
		masterKey:    masterKey,
	}, nil
}

func (p *TenantProvisionerService) ProvisionOrganization(ctx context.Context, organizationID string) error {
	if p.controlDB == nil {
		return fmt.Errorf("control db is nil")
	}
	if strings.TrimSpace(p.controlDBURL) == "" {
		return fmt.Errorf("control db url is empty")
	}
	if len(p.masterKey) != 32 {
		return fmt.Errorf("invalid master key")
	}

	if organizationID == "" {
		return nil
	}

	log.Printf("tenant_provisioner: loading organization %s", organizationID)

	var org control.Organization
	if err := p.controlDB.WithContext(ctx).First(&org, "id = ?", organizationID).Error; err != nil {
		return fmt.Errorf("load org: %w", err)
	}

	envelopeSet := org.ConnectionString != "" &&
		org.EncryptedTenantKey != "" &&
		org.TenantKeyIV != "" &&
		org.TenantKeyTag != ""

	switch {
	case org.Status == "destroying" || org.Status == "deleted":
		return river.JobCancel(fmt.Errorf("org %s status=%q, cannot provision", organizationID, org.Status))
	case org.Status == "active" && envelopeSet:
		return nil
	case org.Status == "" || org.Status == "pending":
		// proceed
	default:
		return river.JobCancel(fmt.Errorf("org %s has unrecognized status %q", organizationID, org.Status))
	}

	dbName, tenantURL, err := deriveTenantDBURL(p.controlDBURL, org.ID.String())
	if err != nil {
		return err
	}

	adminSQL, err := sql.Open("pgx", p.adminDBURL)
	if err != nil {
		return fmt.Errorf("open admin db: %w", err)
	}

	appRole := dsnRole(tenantURL)

	log.Printf("tenant_provisioner: creating database %s (owner=%q)", dbName, appRole)
	if err := createDatabaseIfNotExists(ctx, adminSQL, dbName, appRole); err != nil {
		adminSQL.Close()
		return err
	}

	// Make the app role own the tenant database (and thus its public schema, via
	// pg_database_owner) so it can run migrations and own everything it creates.
	// Also heals databases created before ownership was set. Requires the admin
	// role to be able to SET ROLE to the app role (membership) or be a superuser.
	if appRole != "" {
		log.Printf("tenant_provisioner: setting owner of %s to %q", dbName, appRole)
		if err := setDatabaseOwner(ctx, adminSQL, dbName, appRole); err != nil {
			adminSQL.Close()
			return err
		}
	}
	adminSQL.Close()

	tenantSQL, err := sql.Open("pgx", tenantURL)
	if err != nil {
		return fmt.Errorf("open tenant db: %w", err)
	}
	defer tenantSQL.Close()

	log.Printf("tenant_provisioner: running tenant migrations on %s", dbName)
	if err := runTenantMigrations(tenantSQL); err != nil {
		return err
	}

	log.Printf("tenant_provisioner: wrapping tenant DEK for %s", organizationID)
	dek, err := utils.RandomBytes(32)
	if err != nil {
		return fmt.Errorf("generate tenant dek: %w", err)
	}

	ct, iv, tag, err := utils.EncryptAESGCM(dek, p.masterKey)
	if err != nil {
		return fmt.Errorf("wrap tenant dek: %w", err)
	}

	log.Printf("tenant_provisioner: marking organization %s active", organizationID)
	if err := p.controlDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fresh control.Organization
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&fresh, "id = ?", organizationID).Error; err != nil {
			return err
		}

		// Re-check for idempotency in case another retry already finished.
		if fresh.ConnectionString != "" &&
			fresh.EncryptedTenantKey != "" &&
			fresh.TenantKeyIV != "" &&
			fresh.TenantKeyTag != "" {
			return nil
		}

		return tx.Model(&fresh).Updates(map[string]any{
			"connection_string":    tenantURL,
			"encrypted_tenant_key": utils.EncodeB64(ct),
			"tenant_key_iv":        utils.EncodeB64(iv),
			"tenant_key_tag":       utils.EncodeB64(tag),
			"status":               "active",
		}).Error
	}); err != nil {
		return fmt.Errorf("update org after provisioning: %w", err)
	}

	return nil
}

func runTenantMigrations(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}
	goose.SetBaseFS(migrations.TenantFS)
	if err := goose.Up(db, "tenant"); err != nil {
		return fmt.Errorf("tenant goose up: %w", err)
	}
	return nil
}

func createDatabaseIfNotExists(ctx context.Context, adminDB *sql.DB, dbName, owner string) error {
	var exists bool
	if err := adminDB.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`,
		dbName,
	).Scan(&exists); err != nil {
		return fmt.Errorf("check database existence: %w", err)
	}

	if exists {
		return nil
	}

	stmt := `CREATE DATABASE ` + quoteIdentifier(dbName)
	if owner != "" {
		stmt += ` OWNER ` + quoteIdentifier(owner)
	}
	if _, err := adminDB.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("create database %s: %w", dbName, err)
	}
	return nil
}

func quoteIdentifier(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// setDatabaseOwner reassigns ownership of the tenant database to the app role,
// so it owns the database and (via pg_database_owner) the public schema. This
// lets the app role run migrations and own every object it creates, and heals
// databases created before ownership was set. It is idempotent. The admin role
// must be able to SET ROLE to the app role (membership) or be a superuser.
func setDatabaseOwner(ctx context.Context, adminDB *sql.DB, dbName, owner string) error {
	stmt := fmt.Sprintf("ALTER DATABASE %s OWNER TO %s", quoteIdentifier(dbName), quoteIdentifier(owner))
	if _, err := adminDB.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("set owner of %s to %s: %w", dbName, owner, err)
	}
	return nil
}

type TenantProvisionerArgs struct {
	OrganizationID string `json:"organization_id"`
}

func (TenantProvisionerArgs) Kind() string { return "tenant_provisioner" }

func (TenantProvisionerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{ByArgs: true},
	}
}

type TenantProvisionerWorker struct {
	river.WorkerDefaults[TenantProvisionerArgs]
	svc *TenantProvisionerService
}

func NewTenantProvisionerWorker(svc *TenantProvisionerService) *TenantProvisionerWorker {
	return &TenantProvisionerWorker{svc: svc}
}

func (w *TenantProvisionerWorker) Work(ctx context.Context, job *river.Job[TenantProvisionerArgs]) error {
	return w.svc.ProvisionOrganization(ctx, job.Args.OrganizationID)
}
