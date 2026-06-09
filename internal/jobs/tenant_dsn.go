package jobs

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

func tenantDBName(orgID string) (string, error) {
	id, err := uuid.Parse(strings.TrimSpace(orgID))
	if err != nil {
		return "", fmt.Errorf("invalid org id: %w", err)
	}
	return "tenant_" + strings.ReplaceAll(id.String(), "-", ""), nil
}

func deriveAdminDBURL(controlURL, override string) (string, error) {
	if s := strings.TrimSpace(override); s != "" {
		return s, nil
	}

	u, err := url.Parse(strings.TrimSpace(controlURL))
	if err != nil {
		return "", fmt.Errorf("parse control db url: %w", err)
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return "", fmt.Errorf("unsupported db scheme %q", u.Scheme)
	}

	u.Path = "/postgres"
	u.RawPath = ""
	return u.String(), nil
}

// deriveAdminTenantDBURL points the admin DSN at a specific tenant database, so
// privileged per-database statements (e.g. GRANT ... ON SCHEMA public) run as
// the owner/superuser rather than the app role.
func deriveAdminTenantDBURL(adminURL, dbName string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(adminURL))
	if err != nil {
		return "", fmt.Errorf("parse admin db url: %w", err)
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return "", fmt.Errorf("unsupported db scheme %q", u.Scheme)
	}
	u.Path = "/" + dbName
	u.RawPath = ""
	return u.String(), nil
}

// dsnRole returns the role (username) encoded in a postgres URL, or "" if none.
func dsnRole(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.User == nil {
		return ""
	}
	return u.User.Username()
}

func deriveTenantDBURL(controlURL, orgID string) (dbName string, tenantURL string, err error) {
	u, err := url.Parse(strings.TrimSpace(controlURL))
	if err != nil {
		return "", "", fmt.Errorf("parse control db url: %w", err)
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return "", "", fmt.Errorf("unsupported db scheme %q", u.Scheme)
	}

	dbName, err = tenantDBName(orgID)
	if err != nil {
		return "", "", err
	}

	u.Path = "/" + dbName
	u.RawPath = ""

	return dbName, u.String(), nil
}
