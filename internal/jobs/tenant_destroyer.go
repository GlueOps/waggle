package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/glueops/waggle/internal/database"
	"github.com/glueops/waggle/internal/models/control"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type TenantDestroyerService struct {
	controlDB    *gorm.DB
	controlDBURL string
	adminDBURL   string
	tenants      *database.TenantManager
}

func NewTenantDestroyerService(controlDB *gorm.DB, controlDBURL, adminOverride string, tenants *database.TenantManager) (*TenantDestroyerService, error) {
	adminURL, err := deriveAdminDBURL(controlDBURL, adminOverride)
	if err != nil {
		return nil, fmt.Errorf("resolve admin db url: %w", err)
	}
	return &TenantDestroyerService{
		controlDB:    controlDB,
		controlDBURL: controlDBURL,
		adminDBURL:   adminURL,
		tenants:      tenants,
	}, nil
}

func (d *TenantDestroyerService) DestroyOrganization(ctx context.Context, organizationID string) error {
	log.Printf("tenant_destroyer: loading organization %s", organizationID)

	var org control.Organization
	if err := d.controlDB.WithContext(ctx).First(&org, "id = ?", organizationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // idempotent
		}
		return fmt.Errorf("load org: %w", err)
	}

	switch org.Status {
	case "destroying":
		// proceed
	case "deleted":
		return nil // idempotent — already destroyed
	default:
		return river.JobCancel(fmt.Errorf("org %s status=%q, refusing to destroy (expected 'destroying')", organizationID, org.Status))
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return river.JobCancel(fmt.Errorf("invalid org id %q: %w", organizationID, err))
	}

	log.Printf("tenant_destroyer: forgetting cached tenant connection for %s", organizationID)
	d.tenants.Forget(orgUUID)

	// Drop the tenant database.
	dbName, err := tenantDBName(organizationID)
	if err != nil {
		return fmt.Errorf("derive tenant db name: %w", err)
	}

	adminDB, err := sql.Open("pgx", d.adminDBURL)
	if err != nil {
		return fmt.Errorf("open admin db: %w", err)
	}
	defer adminDB.Close()

	log.Printf("tenant_destroyer: dropping database %s", dbName)
	if _, err := adminDB.ExecContext(ctx, `DROP DATABASE IF EXISTS `+quoteIdentifier(dbName)+` WITH (FORCE)`); err != nil {
		return fmt.Errorf("drop database %s: %w", dbName, err)
	}

	// Clean up all org-scoped control DB records.
	log.Printf("tenant_destroyer: cleaning up control records for %s", organizationID)
	if err := d.controlDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		orgID := org.ID.String()
		deletes := []struct {
			table string
			where string
		}{
			{"auth_audit_events", "organization_id = ?"},
			{"token_sessions", "organization_id = ?"},
			{"users", "organization_id = ?"},
			{"organizations", "id = ?"},
		}

		for _, d := range deletes {
			if err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", d.table, d.where), orgID).Error; err != nil {
				return fmt.Errorf("delete from %s: %w", d.table, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	log.Printf("tenant_destroyer: organization %s destroyed", organizationID)
	return nil
}

type TenantDestroyerArgs struct {
	OrganizationID string `json:"organization_id"`
}

func (TenantDestroyerArgs) Kind() string { return "tenant_destroyer" }

func (TenantDestroyerArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		UniqueOpts: river.UniqueOpts{ByArgs: true},
	}
}

type TenantDestroyerWorker struct {
	river.WorkerDefaults[TenantDestroyerArgs]
	svc *TenantDestroyerService
}

func NewTenantDestroyerWorker(svc *TenantDestroyerService) *TenantDestroyerWorker {
	return &TenantDestroyerWorker{svc: svc}
}

func (w *TenantDestroyerWorker) Work(ctx context.Context, job *river.Job[TenantDestroyerArgs]) error {
	return w.svc.DestroyOrganization(ctx, job.Args.OrganizationID)
}
