package migrations

import "embed"

//go:embed control/*.sql
var ControlFS embed.FS

//go:embed tenant/*.sql
var TenantFS embed.FS
