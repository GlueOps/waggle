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

	log.Printf("tenant_provisioner: creating database %s", dbName)
	if err := createDatabaseIfNotExists(ctx, adminSQL, dbName); err != nil {
		adminSQL.Close()
		return err
	}
	adminSQL.Close()

	// Ensure the app role can create objects in the new tenant DB's public
	// schema before goose runs. When the admin/owner role differs from the app
	// role (e.g. a superuser ADMIN_DATABASE_URL), PostgreSQL 15+ would otherwise
	// deny tenant migrations with "permission denied for schema public".
	if appRole := dsnRole(tenantURL); appRole != "" {
		log.Printf("tenant_provisioner: granting public schema on %s to %q", dbName, appRole)
		if err := grantTenantPublicSchema(ctx, p.adminDBURL, dbName, appRole); err != nil {
			return err
		}
	}

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

func createDatabaseIfNotExists(ctx context.Context, adminDB *sql.DB, dbName string) error {
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

	_, err := adminDB.ExecContext(ctx, `CREATE DATABASE `+quoteIdentifier(dbName))
	if err != nil {
		return fmt.Errorf("create database %s: %w", dbName, err)
	}
	return nil
}

func quoteIdentifier(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// grantTenantPublicSchema connects to the freshly created tenant database as the
// admin/owner role and grants the app role full access to the public schema
// (USAGE + CREATE), which is what goose needs to create goose_db_version and the
// tenant tables. It is idempotent and a no-op cost when the admin and app roles
// are the same (the app already owns the database it created).
func grantTenantPublicSchema(ctx context.Context, adminURL, dbName, role string) error {
	adminTenantURL, err := deriveAdminTenantDBURL(adminURL, dbName)
	if err != nil {
		return fmt.Errorf("derive admin url for %s: %w", dbName, err)
	}

	db, err := sql.Open("pgx", adminTenantURL)
	if err != nil {
		return fmt.Errorf("open admin connection to %s: %w", dbName, err)
	}
	defer db.Close()

	stmt := "GRANT ALL ON SCHEMA public TO " + quoteIdentifier(role)
	if _, err := db.ExecContext(ctx, stmt); err != nil {
		return fmt.Errorf("grant public schema on %s to %s: %w", dbName, role, err)
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
