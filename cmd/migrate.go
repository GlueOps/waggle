package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/internal/migrations"
	"github.com/glueops/waggle/internal/models/control"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	migrateDBUrl string
	migrateScope string
	migrateOrg   string
	migrateSteps int
	migrateRiver bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	// Runtime failures (e.g. the schema-privilege preflight) print their own
	// guidance via Execute's log.Fatal; don't bury it under a flag-usage dump.
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMigrations(cmd.Context(), gooseAction{kind: "up"})
	},
}

type gooseAction struct {
	kind    string
	version int64
}

func newMigrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply migrations (all by default, or --steps via up-by-one)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations(cmd.Context(), gooseAction{kind: "up"})
		},
	}
}

func newMigrateUpToCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "up-to VERSION",
		Short: "Migrate up to a specific VERSION",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := parseVersion(args[0])
			if err != nil {
				return err
			}
			return runMigrations(cmd.Context(), gooseAction{kind: "up-to", version: v})
		},
	}
	return c
}

func newMigrateDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Rollback migrations (one by default, or --steps N)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations(cmd.Context(), gooseAction{kind: "down"})
		},
	}
}

func newMigrateDownToCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down-to VERSION",
		Short: "Rollback migrations to a specific VERSION",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := parseVersion(args[0])
			if err != nil {
				return err
			}
			return runMigrations(cmd.Context(), gooseAction{kind: "down-to", version: v})
		},
	}
}

func newMigrateRedoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "redo",
		Short: "Re-run the latest migration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations(cmd.Context(), gooseAction{kind: "redo"})
		},
	}
}

func newMigrateResetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reset",
		Short: "Rollback ALL migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations(cmd.Context(), gooseAction{kind: "reset"})
		},
	}
}

func newMigrateStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Print migration status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations(cmd.Context(), gooseAction{kind: "status"})
		},
	}
}

func newMigrateVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print current migration version",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrations(cmd.Context(), gooseAction{kind: "version"})
		},
	}
}

func runMigrations(ctx context.Context, action gooseAction) error {
	scope := strings.ToLower(strings.TrimSpace(migrateScope))
	if scope == "" {
		scope = "all"
	}
	if scope != "control" && scope != "tenants" && scope != "all" {
		return fmt.Errorf("invalid --scope %q (must be control|tenants|all)", migrateScope)
	}

	dbURL, err := resolveControlDBURL()
	if err != nil {
		return err
	}

	// Goose setup (dialect is global)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Connect control DB
	controlSQL, err := sql.Open("pgx", dbURL)
	if err != nil {
		return err
	}
	defer controlSQL.Close()

	// For down/reset, fetch tenants FIRST (control schema might disappear)
	var orgs []control.Organization
	if scope != "control" && (action.kind == "down" || action.kind == "down-to" || action.kind == "reset") {
		orgs = fetchOrgsBestEffort(controlSQL, migrateOrg)
	}

	if scope == "control" || scope == "all" {
		if actionWritesSchema(action.kind) {
			if err := preflightSchemaCreate(ctx, controlSQL, "control"); err != nil {
				return err
			}
		}
		log.Printf("Control DB: goose %s", action.kind)
		if err := runGooseAction(controlSQL, migrations.ControlFS, "control", action); err != nil {
			return fmt.Errorf("control migrations failed: %w", err)
		}

		// River migrations generally live in control DB
		if migrateRiver {
			if err := runRiver(ctx, controlSQL, action); err != nil {
				return fmt.Errorf("river migrations failed: %w", err)
			}
		}
	}

	if scope == "tenants" || scope == "all" {
		// For up/status/etc, fetch orgs after control is ensured
		if orgs == nil {
			orgs = fetchOrgsBestEffort(controlSQL, migrateOrg)
		}
		if len(orgs) == 0 {
			log.Println("No tenants found (or control schema not available).")
			return nil
		}

		for _, org := range orgs {
			if org.ConnectionString == "" {
				log.Printf("Skipping org %s (%s): empty connection_string", org.Name, org.ID)
				continue
			}

			log.Printf("Tenant DB [%s / %s]: goose %s", org.Name, org.ID, action.kind)
			tenantSQL, err := sql.Open("pgx", org.ConnectionString)
			if err != nil {
				log.Printf("Failed to connect tenant DB %s: %v", org.ID, err)
				continue
			}

			if actionWritesSchema(action.kind) {
				if err := preflightSchemaCreate(ctx, tenantSQL, fmt.Sprintf("tenant %s (%s)", org.Name, org.ID)); err != nil {
					log.Printf("%v", err)
					_ = tenantSQL.Close()
					continue
				}
			}

			if err := runGooseAction(tenantSQL, migrations.TenantFS, "tenant", action); err != nil {
				log.Printf("Tenant migrations failed for %s: %v", org.ID, err)
			}
			_ = tenantSQL.Close()
		}
	}
	return nil
}

func runGooseAction(db *sql.DB, base fs.FS, dir string, action gooseAction) error {
	goose.SetBaseFS(base)

	switch action.kind {
	case "up":
		if migrateSteps > 0 {
			for i := 0; i < migrateSteps; i++ {
				if err := goose.UpByOne(db, dir); err != nil {
					return err
				}
			}
			return nil
		}
		return goose.Up(db, dir)

	case "up-to":
		return goose.UpTo(db, dir, action.version)

	case "down":
		steps := migrateSteps
		if steps <= 0 {
			steps = 1
		}
		for i := 0; i < steps; i++ {
			if err := goose.Down(db, dir); err != nil {
				return err
			}
		}
		return nil

	case "down-to":
		return goose.DownTo(db, dir, action.version)

	case "redo":
		return goose.Redo(db, dir)

	case "reset":
		return goose.Reset(db, dir)

	case "status":
		return goose.Status(db, dir)

	case "version":
		return goose.Version(db, dir)

	default:
		return fmt.Errorf("unknown goose action: %s", action.kind)
	}
}

// actionWritesSchema reports whether a goose action creates or drops objects
// (and therefore needs the goose_db_version table, which requires CREATE on the
// public schema). Pure read actions are excluded so they don't trip the
// preflight on a role that only has read access.
func actionWritesSchema(kind string) bool {
	switch kind {
	case "status", "version":
		return false
	default:
		return true
	}
}

// preflightSchemaCreate verifies the connected role can create objects in the
// "public" schema before goose tries to create goose_db_version. PostgreSQL 15+
// no longer grants CREATE on "public" to non-owner roles, and
// GRANT ALL PRIVILEGES ON DATABASE does not cover schema-level privileges — so
// this is the most common cause of "permission denied for schema public" on a
// freshly provisioned database. When it can't run the check itself, it logs and
// lets goose surface the real error rather than blocking.
func preflightSchemaCreate(ctx context.Context, db *sql.DB, label string) error {
	const q = `
SELECT current_user,
       current_database(),
       EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'public'),
       CASE WHEN EXISTS (SELECT 1 FROM pg_namespace WHERE nspname = 'public')
            THEN has_schema_privilege(current_user, 'public', 'CREATE')
            ELSE NULL END`

	var (
		user     string
		database string
		schemaOK bool
		createOK sql.NullBool
	)
	if err := db.QueryRowContext(ctx, q).Scan(&user, &database, &schemaOK, &createOK); err != nil {
		log.Printf("Preflight (%s): could not check schema privileges (%v); continuing", label, err)
		return nil
	}

	if !schemaOK {
		return fmt.Errorf("%s database %s has no \"public\" schema.\n\n"+
			"Create it (as a superuser or the database owner) and hand ownership to the app role:\n\n"+
			"    \\c %s\n"+
			"    CREATE SCHEMA public AUTHORIZATION %s;\n\n"+
			"Then re-run: waggle migrate",
			label, quoteIdent(database), quoteIdent(database), quoteIdent(user))
	}
	if createOK.Valid && createOK.Bool {
		return nil
	}

	return fmt.Errorf("role %s cannot CREATE in schema \"public\" of %s database %s.\n\n"+
		"PostgreSQL 15+ does not grant CREATE on the \"public\" schema to non-owner roles, and\n"+
		"GRANT ALL PRIVILEGES ON DATABASE does not cover schema-level privileges.\n\n"+
		"Fix it with ONE of the following, connected as a superuser or the database owner:\n\n"+
		"    -- Option A (recommended): make the role own the database, so it owns \"public\":\n"+
		"    ALTER DATABASE %s OWNER TO %s;\n\n"+
		"    -- Option B: grant CREATE on the public schema (run while connected to %s):\n"+
		"    \\c %s\n"+
		"    GRANT ALL ON SCHEMA public TO %s;\n\n"+
		"Per-tenant provisioning also needs the role to create databases:\n"+
		"    ALTER ROLE %s CREATEDB;\n\n"+
		"Then re-run: waggle migrate",
		quoteIdent(user), label, quoteIdent(database),
		quoteIdent(database), quoteIdent(user),
		quoteIdent(database), quoteIdent(database), quoteIdent(user),
		quoteIdent(user))
}

// quoteIdent double-quotes a SQL identifier so the remediation snippets are
// safe to copy-paste verbatim.
func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

func runRiver(ctx context.Context, controlSQL *sql.DB, action gooseAction) error {
	// Only do River schema changes for actions that change schema.
	switch action.kind {
	case "status", "version":
		return nil
	}

	log.Printf("Control DB: river migrate (%s)", riverDirectionFor(action))
	m, err := rivermigrate.New(riverdatabasesql.New(controlSQL), nil)
	if err != nil {
		return err
	}

	dir := riverDirectionFor(action)

	// Reasonable defaults:
	// - Up: apply all outstanding
	// - Down: apply 1 step unless --steps specified
	// - Reset: remove all River tables (TargetVersion = -1)
	var opts *rivermigrate.MigrateOpts
	if dir == rivermigrate.DirectionDown {
		opts = &rivermigrate.MigrateOpts{}
		if action.kind == "reset" {
			opts.TargetVersion = -1
		} else if migrateSteps > 0 {
			opts.MaxSteps = migrateSteps
		} // else defaults to 1 down step per River docs
	}

	_, err = m.Migrate(ctx, dir, opts)
	return err
}

func riverDirectionFor(action gooseAction) rivermigrate.Direction {
	switch action.kind {
	case "down", "down-to", "reset":
		return rivermigrate.DirectionDown
	default:
		return rivermigrate.DirectionUp
	}
}

func resolveControlDBURL() (string, error) {
	if migrateDBUrl != "" {
		return migrateDBUrl, nil
	}
	cfg, err := config.Load()
	if err == nil && cfg.DatabaseURL != "" {
		return cfg.DatabaseURL, nil
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v, nil
	}
	return "", errors.New("DATABASE_URL is required (via --db-url, config, or env var)")
}

func fetchOrgsBestEffort(controlSQL *sql.DB, orgFilter string) []control.Organization {
	controlGorm, err := gorm.Open(postgres.New(postgres.Config{Conn: controlSQL}), &gorm.Config{})
	if err != nil {
		log.Printf("Could not init gorm for control DB: %v", err)
		return nil
	}

	var orgs []control.Organization
	q := controlGorm.Model(&control.Organization{}).Select("id", "name", "connection_string")
	if orgFilter != "" {
		q = q.Where("id::text = ? OR name = ?", orgFilter, orgFilter)
	}
	if err := q.Find(&orgs).Error; err != nil {
		log.Printf("Could not fetch orgs: %v", err)
		return nil
	}
	return orgs
}

func parseVersion(s string) (int64, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("invalid VERSION %q", s)
	}
	return v, nil
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.PersistentFlags().StringVar(&migrateDBUrl, "db-url", "", "Control DB URL (or DATABASE_URL / config fallback)")
	migrateCmd.PersistentFlags().StringVar(&migrateScope, "scope", "all", "Migration scope: control|tenants|all")
	migrateCmd.PersistentFlags().StringVar(&migrateOrg, "org", "", "Limit tenant migrations to a single org (id or name)")
	migrateCmd.PersistentFlags().IntVar(&migrateSteps, "steps", 0, "Number of steps (for up-by-one/down). 0 means default behavior.")
	migrateCmd.PersistentFlags().BoolVar(&migrateRiver, "river", true, "Also migrate River tables in the control DB")

	// Subcommands
	migrateCmd.AddCommand(newMigrateUpCmd())
	migrateCmd.AddCommand(newMigrateUpToCmd())
	migrateCmd.AddCommand(newMigrateDownCmd())
	migrateCmd.AddCommand(newMigrateDownToCmd())
	migrateCmd.AddCommand(newMigrateRedoCmd())
	migrateCmd.AddCommand(newMigrateResetCmd())
	migrateCmd.AddCommand(newMigrateStatusCmd())
	migrateCmd.AddCommand(newMigrateVersionCmd())

	// Cobra only consults the executed command (and the root) for these flags,
	// so propagate them to every subcommand — otherwise a runtime failure like
	// the schema-privilege preflight gets buried under a flag-usage dump.
	for _, c := range migrateCmd.Commands() {
		c.SilenceUsage = true
		c.SilenceErrors = true
	}
}
