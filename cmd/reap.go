package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/glueops/waggle/internal/config"
	"github.com/spf13/cobra"
)

const tenantDBPrefix = "tenant_"

var (
	reapConfirm     bool
	reapAdminDBURL  string
	reapAllowNoOrgs bool
)

var reapCmd = &cobra.Command{
	Use:   "reap",
	Short: "Reap zombie tenant databases that have no matching organization",
	Long: `Reap finds PostgreSQL databases prefixed "tenant_" that do not correspond to
any organization in the control database and drops them.

A tenant database is considered LIVE if its name is tenant_<org-id-without-dashes>
for some row in organizations, or if it is referenced by an organization's
connection_string. Anything else with the tenant_ prefix is a ZOMBIE — typically
left behind by a failed or rolled-back provisioning.

Reap runs as a DRY RUN by default and only prints what it would drop. Pass
--confirm to actually DROP DATABASE (irreversible).`,
	// Runtime failures print their own message via Execute's log.Fatal; don't
	// bury them under a flag-usage dump.
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runReap(cmd.Context())
	},
}

func runReap(ctx context.Context) error {
	controlURL, err := resolveControlDBURL()
	if err != nil {
		return err
	}
	adminURL, err := resolveAdminDBURL(controlURL)
	if err != nil {
		return err
	}

	// 1) Build the set of LIVE tenant DB names from the control DB. If this read
	// fails we must NOT reap — an empty/failed list would mark every tenant DB a
	// zombie and delete them all.
	controlSQL, err := sql.Open("pgx", controlURL)
	if err != nil {
		return fmt.Errorf("open control db: %w", err)
	}
	defer controlSQL.Close()

	live, orgCount, err := liveTenantDBs(ctx, controlSQL)
	if err != nil {
		return fmt.Errorf("list organizations (refusing to reap): %w", err)
	}

	// 2) List the actual tenant_* databases via the admin connection.
	adminSQL, err := sql.Open("pgx", adminURL)
	if err != nil {
		return fmt.Errorf("open admin db: %w", err)
	}
	defer adminSQL.Close()

	all, err := tenantDatabases(ctx, adminSQL)
	if err != nil {
		return fmt.Errorf("list tenant databases: %w", err)
	}

	// 3) Anything present on disk but not LIVE is a zombie.
	var zombies []string
	for _, name := range all {
		if !live[name] {
			zombies = append(zombies, name)
		}
	}
	sort.Strings(zombies)

	log.Printf("reap: %d organization(s), %d tenant_* database(s), %d zombie(s)", orgCount, len(all), len(zombies))
	for _, name := range all {
		state := "live"
		if !live[name] {
			state = "ZOMBIE"
		}
		log.Printf("reap:   %-45s %s", name, state)
	}

	if len(zombies) == 0 {
		log.Printf("reap: nothing to do")
		return nil
	}

	// Safety: zero orgs but tenant DBs present almost always means the wrong or
	// an empty control DB, not that every tenant is a zombie. Refuse unless forced.
	if orgCount == 0 && !reapAllowNoOrgs {
		return fmt.Errorf("found %d tenant database(s) but 0 organizations; refusing to reap "+
			"(pass --allow-no-orgs only if the control DB really is empty)", len(all))
	}

	if !reapConfirm {
		log.Printf("reap: DRY RUN — would drop %d database(s); pass --confirm to drop", len(zombies))
		return nil
	}

	var failed int
	for _, name := range zombies {
		log.Printf("reap: dropping %s", name)
		// WITH (FORCE) terminates other sessions so the drop doesn't block on
		// stray connections (PostgreSQL 13+).
		if _, err := adminSQL.ExecContext(ctx, `DROP DATABASE IF EXISTS `+quoteIdent(name)+` WITH (FORCE)`); err != nil {
			log.Printf("reap: FAILED to drop %s: %v", name, err)
			failed++
		}
	}
	if failed > 0 {
		return fmt.Errorf("reap: %d of %d drop(s) failed", failed, len(zombies))
	}
	log.Printf("reap: dropped %d zombie database(s)", len(zombies))
	return nil
}

// liveTenantDBs returns the set of tenant database names that belong to an
// existing organization (by derived name and by connection_string), plus the
// number of organizations seen.
func liveTenantDBs(ctx context.Context, controlSQL *sql.DB) (map[string]bool, int, error) {
	rows, err := controlSQL.QueryContext(ctx, `SELECT id::text, connection_string FROM organizations`)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	live := make(map[string]bool)
	count := 0
	for rows.Next() {
		var id string
		var conn sql.NullString
		if err := rows.Scan(&id, &conn); err != nil {
			return nil, 0, err
		}
		count++
		if name := tenantDBNameFromID(id); name != "" {
			live[name] = true
		}
		if conn.Valid {
			if name := dbNameFromURL(conn.String); name != "" {
				live[name] = true
			}
		}
	}
	return live, count, rows.Err()
}

// tenantDatabases lists every database with the tenant_ prefix. The underscore
// is escaped so it is matched literally rather than as a LIKE wildcard.
func tenantDatabases(ctx context.Context, adminSQL *sql.DB) ([]string, error) {
	rows, err := adminSQL.QueryContext(ctx,
		`SELECT datname FROM pg_database WHERE datname LIKE 'tenant\_%' ESCAPE '\' ORDER BY datname`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}

func tenantDBNameFromID(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	return tenantDBPrefix + strings.ReplaceAll(id, "-", "")
}

func dbNameFromURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(u.Path, "/")
}

func resolveAdminDBURL(controlURL string) (string, error) {
	if reapAdminDBURL != "" {
		return reapAdminDBURL, nil
	}
	if cfg, err := config.Load(); err == nil && strings.TrimSpace(cfg.AdminDatabaseURL) != "" {
		return cfg.AdminDatabaseURL, nil
	}
	// Fall back to the control credentials against the maintenance database.
	u, err := url.Parse(strings.TrimSpace(controlURL))
	if err != nil {
		return "", fmt.Errorf("parse control db url: %w", err)
	}
	u.Path = "/postgres"
	u.RawPath = ""
	return u.String(), nil
}

func init() {
	rootCmd.AddCommand(reapCmd)
	reapCmd.Flags().BoolVar(&reapConfirm, "confirm", false, "Actually DROP the zombie databases (default: dry-run)")
	reapCmd.Flags().StringVar(&reapAdminDBURL, "admin-db-url", "", "Admin DB URL used for DROP DATABASE (default: ADMIN_DATABASE_URL, else control creds @ /postgres)")
	reapCmd.Flags().BoolVar(&reapAllowNoOrgs, "allow-no-orgs", false, "Permit reaping even when the control DB reports zero organizations")
}
