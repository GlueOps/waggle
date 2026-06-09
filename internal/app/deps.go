package app

import (
	"github.com/glueops/waggle/internal/config"
	"github.com/glueops/waggle/internal/database"
	"github.com/glueops/waggle/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type Deps struct {
	Config        *config.Config
	ControlDB     *gorm.DB
	DBPool        *pgxpool.Pool
	Tenants       *database.TenantManager
	River         *river.Client[pgx.Tx]
	Jobs          any
	Auth          *service.AuthService
	Fleet         *service.FleetService
	Tokens        *service.TokenService
	TokenSessions *service.TokenSessionService
	Policies      *service.PolicyService
	APIKeys       *service.APIKeyService
	Orgs          *service.OrgService
}
