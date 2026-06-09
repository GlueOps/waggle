# Waggle

**A placement oracle and ledger for Proxmox infrastructure.**

Waggle manages VM placement decisions and maintains an authoritative ledger of node pool configurations across Proxmox hypervisors. Named after the "waggle dance" bees use to communicate directions and timing to the rest of the hive.

> **Note**: Waggle does not create or provision VMs itself. It serves as the source of truth for placement state. VM provisioning is handled by companion systems like River-based workers or your Terraform provider.

## Features

- **Declarative Node Pool Management** — Define desired VM counts per pool; Waggle tracks actual placements
- **Multi-Tenant Architecture** — Isolated organizations with encrypted per-tenant databases
- **Placement Oracle** — Allocate VMs to hypervisors based on resource availability and constraints
- **Encryption at Rest** — Two-tier KEK→DEK envelope encryption for sensitive data
- **Job Queue System** — River-based background workers for tenant provisioning/destruction
- **OpenAPI-Driven SDKs** — Automatically generated TypeScript, Go, and REST clients
- **Audit Trail** — Complete auth event logging and session tracking
- **Database Versioning** — Schema management via Goose migrations for control and tenant databases

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.26.3 |
| API Framework | [Huma](https://github.com/danielgtaylor/huma) v2 |
| Routing | Chi v5 |
| Database | PostgreSQL |
| ORM | GORM |
| Job Queue | [River](https://github.com/riverqueue/river) |
| Migrations | Goose v3 + Atlas |
| Auth | JWT (golang-jwt) |
| Encryption | AES-GCM (crypto/cipher) |
| CLI | Cobra |

## Prerequisites

- **Go** 1.26.3 or later
- **PostgreSQL** 13+
- **Docker & Docker Compose** (for local development)
- **Node.js** (for SDK generation)

## Quick Start

### 1. Clone & Setup

```bash
git clone https://github.com/glueops/waggle.git
cd waggle
cp .env.example .env  # Configure as needed
```

### 2. Start Local Database

```bash
just up  # Runs docker-compose up + migrations
```

### 3. Run the Server

```bash
go run . serve
```

Server starts on `http://localhost:8080` with API docs at `/api/docs`.

### 4. Run Background Workers

In a separate terminal:

```bash
go run . worker
```

Workers process River jobs (tenant provisioning, destruction, etc.).

## Installation & Development

### Build from Source

```bash
go build -o waggle ./cmd
./waggle serve
```

### Available Commands

```bash
waggle serve              # Start API server
waggle worker             # Run background job workers
waggle migrate [args]     # Manage database migrations
waggle encrypt            # Key management (generate-master, rotate-master)
waggle generate           # Code generation (migrations, SDKs)
waggle version            # Print version info
```

### Useful Justfile Targets

```bash
just generate <name>      # Create new DB migration
just sdk                  # Generate OpenAPI spec + SDKs
just migrate up           # Run pending migrations
just reset                # Drop & recreate database
just release              # Build release artifacts
```

## Documentation

In-depth reference docs live in [`docs/`](docs/README.md):

- [Architecture](docs/ARCHITECTURE.md) — planes, request flow, jobs, encryption, placement
- [Database](docs/DATABASE.md) — control & tenant schema reference
- [API](docs/API.md) — endpoint and auth reference
- [Configuration](docs/CONFIGURATION.md) — every environment variable
- [Development](docs/DEVELOPMENT.md) — setup, code generation, migrations, contributing

## Architecture Overview

### Data Model

**Control Database** (shared, encrypted at rest)
- Organizations & Tenants
- Users, Accounts, Emails
- Token Sessions & Refresh Tokens
- Auth Audit Events
- Policies & Passpkeys

**Tenant Databases** (per-organization, encrypted connection string)
- Datacenters
- Hypervisors (CPU/RAM/Disk capacity tracking)
- Node Pools (declared desired count)
- VM Placements (actual VMID assignments)
- Slots (VM size definitions: vCPU, RAM, Disk)

### Encryption

Waggle uses envelope encryption for sensitive data:
```
Master Key (KEK) 32 bytes
    ↓
Tenant Key (DEK) — encrypted with Master Key
    ↓
Field-Level Encryption — encrypted with Tenant Key
```

**Intentionally Plaintext**: `Organization.ConnectionString` is stored unencrypted in the control DB for operational simplicity.

### Job Queue (River)

Asynchronous work is handled by River:
- `TenantProvisionerWorker` — Creates tenant databases
- `TenantDestroyerWorker` — Destroys tenant databases

Scale workers independently from the API server.

## API Documentation

### Quick Links

- **OpenAPI Spec**: `GET /api/openapi.json`
- **Interactive Docs**: `GET /api/docs`
- **Base Path**: `/api/v1`

### Key Endpoints

```
POST   /auth/register              — Create account
POST   /auth/login                 — Generate JWT
POST   /auth/refresh               — Refresh access token
POST   /auth/logout                — Revoke session

GET    /organizations              — List orgs (paginated)
GET    /organizations/{id}         — Get org details
POST   /organizations              — Create org

GET    /datacenters/{id}           — Get datacenter
GET    /hypervisors/{id}           — Get hypervisor state
GET    /pools/{id}                 — Get pool placements
POST   /placements                 — Create/update placement
```

See `/docs` for the complete interactive specification.

## Deployment

### Environment Variables

```bash
# Server
BIND_HOST=0.0.0.0
BIND_PORT=8080
BASE_PATH=/api/v1
FRONTEND_MODE=embed  # none|embed|external

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/waggle?sslmode=disable
ADMIN_DATABASE_URL=  # Optional: superuser DSN for tenant creation

# Encryption
ENCRYPTION_MASTER_KEY=<base64-32-bytes>  # Generate via: go run . encrypt generate-master

# JWT
JWT_SECRET=<random-secret>
JWT_ISSUER=waggle
JWT_ACCESS_TTL_MIN=60
JWT_REFRESH_TTL_HOUR=720
JWT_AUDIENCE=waggle-api
```

### Generate Master Key

```bash
go run . encrypt generate-master
# Output: ENCRYPTION_MASTER_KEY=<base64>
```

### Rotate Master Key

```bash
go run . encrypt rotate-master
# Safely re-wraps all tenant keys with the new master key
```

### Migrations

```bash
go run . migrate up              # Apply all pending
go run . migrate status          # Show applied migrations
go run . migrate down            # Rollback one step
go run . migrate down-to 001     # Rollback to specific version
```

Control DB and tenant DBs are migrated separately via the `--scope` flag.

### Docker Deployment

```bash
docker build -t waggle:latest .
docker run -d \
  -e DATABASE_URL="postgres://..." \
  -e ENCRYPTION_MASTER_KEY="..." \
  waggle:latest serve
```

## Contributing

### Code Generation

Waggle uses code generation heavily:

1. **Database Migrations** — Generate via Atlas diff
   ```bash
   just generate add_user_field
   ```

2. **OpenAPI + SDKs** — Auto-generate from Huma router
   ```bash
   just sdk
   # Generates:
   # - docs/openapi.json
   # - ui/src/sdk/          (Hey API, embedded in UI)
   # - sdk/ts/              (npm package)
   # - sdk/go/              (Go module)
   ```

3. **Security** — Generated Go SDK patches request/response logging to exclude sensitive headers/bodies

### Development Workflow

1. Modify GORM models in `internal/models/`
2. Generate migration: `just generate my_change`
3. Review SQL in `internal/migrations/`
4. Run migration: `just migrate up`
5. Modify API handlers in `internal/handlers/`
6. Regenerate SDK if API changes: `just sdk`
7. Commit generated files (migrations, SDK)

### Code Style

- **Go**: Follow `go fmt`, no comments unless the "why" is non-obvious
- **Models**: Use GORM tags for schema; JSON tags for API responses
- **Error Handling**: Return early, wrap with context via `fmt.Errorf`
- **Secrets**: Never log auth headers, tokens, or connection strings — patch generated SDKs

### Testing

```bash
go test ./...
```

Integration tests assume a local Postgres instance (provided by `docker-compose`).

## Troubleshooting

### Database Connection Issues

```bash
# Test control DB connection
DATABASE_URL=postgres://... go run . migrate status

# Test specific tenant
go run . migrate --scope tenants --org <org-id> status
```

### Master Key Rotation

If you rotate the master key, **all River jobs and running processes must be restarted** after the new key is deployed.

### SDK Generation Failures

Ensure Node.js and `npx` are available. Generated SDKs require npm dependencies:
```bash
npm install sdk/ts
npm install sdk/go
```

## License

[Add your license]

## Support

For issues, please open a GitHub issue. For security concerns, email security@glueops.com.

---

**Built with ❤️ by the GlueOps team**
