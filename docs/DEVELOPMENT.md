# Development

## Prerequisites

- **Go** 1.26+
- **Docker** & Docker Compose (local Postgres + Mailpit)
- **Node.js** + Yarn (frontend and SDK generation)
- Optional: [`atlas`](https://atlasgo.io) (migration diffing),
  [`just`](https://github.com/casey/just), [`mprocs`](https://github.com/pvolok/mprocs),
  [`goreleaser`](https://goreleaser.com)

## First-time setup

```bash
cp .env.example .env
# Generate a real master key and paste it into .env:
go run . encrypt generate-master

just up            # docker-compose up -d  +  migrate up
go run . serve     # API on http://localhost:8080, docs at /api/v1/docs
go run . worker    # in another terminal: River background jobs
```

`mprocs` runs the API (in `proxy` frontend mode), the worker, and the Vite dev
server together:

```bash
mprocs   # reads mprocs.yaml
```

## CLI commands

```
waggle serve       Start the HTTP API server
waggle worker      Run River background workers
waggle migrate     Manage database migrations (see below)
waggle encrypt     Key management: generate-master, rotate-master
waggle generate    Code generation: migrations, sdk, terraform
waggle version     Print version / commit / build date
```

### `migrate`

Migrations are powered by Goose and are **scoped** — the control DB and tenant
DBs migrate separately.

```bash
go run . migrate up                         # apply all (control scope by default)
go run . migrate up --steps 1
go run . migrate down --steps 1
go run . migrate up-to <version>
go run . migrate down-to <version>
go run . migrate redo
go run . migrate status
go run . migrate create <name>

# Tenant-scoped:
go run . migrate up --scope tenant --org-id <org-uuid>

# Flags: --db-url, --scope {control|tenant}, --org-id, --steps, --river
```

### `encrypt`

```bash
go run . encrypt generate-master   # print a new ENCRYPTION_MASTER_KEY
go run . encrypt rotate-master     # re-wrap every org DEK under a new KEK
```

After rotating the master key, **restart all serve/worker processes** with the
new value.

## Code generation

Waggle leans on generated code; commit the outputs.

### Migrations (Atlas diff → Goose files)

```bash
just generate add_some_field       # = go run . generate migrations add_some_field
```

This diffs the current GORM models against the DB and writes new Goose
migrations under `internal/migrations/{control,tenant}/`. Review the SQL before
applying.

### OpenAPI spec + SDKs

```bash
just sdk                           # = go run . generate sdk
```

Produces, from the live Huma router:

- `docs/openapi.json` — the spec
- `ui/src/sdk/` — Hey API client embedded in the UI
- `sdk/ts/` — TypeScript npm package
- `sdk/go/` — Go module

The generated Go SDK is patched so request/response logging excludes sensitive
headers and bodies.

### Terraform provider

```bash
just terraform                     # = go run . generate terraform openapi-generator
```

Regenerates `terraform-provider-waggle/` from the spec, then re-applies the
hand-authored overlays in `cmd/overlays/` (e.g. the pool-placements data
source).

## Typical change workflow

1. Edit GORM models in `internal/models/{control,tenant}/`.
2. `just generate <name>` and review the SQL under `internal/migrations/`.
3. `just migrate up`.
4. Edit handlers in `internal/api/` and logic in `internal/service/`.
5. If the API surface changed: `just sdk` (and `just terraform` if needed).
6. Commit code **and** the generated migrations/SDKs together.

## Building

```bash
just build         # builds ui/dist, then `go build -o bin/waggle .`
just release       # goreleaser build --snapshot --clean
```

The frontend is embedded via `//go:embed all:dist` in `ui/embed.go`, so the UI
must be built (`just ui` / `yarn build`) before `go build` for `FRONTEND_MODE=embed`
to serve real assets. Without a build, the embed falls back to "none".

> **Note:** `ui/dist/` is git-ignored (see `ui/.gitignore`); build it locally or
> in CI before producing an embedded binary.

## Testing

```bash
go test ./...
```

Integration tests expect a local Postgres — `just up` provides one.

## Code style

- **Go:** `go fmt`; comments only where the "why" isn't obvious.
- **Models:** GORM tags for schema, JSON tags for API responses.
- **Errors:** return early; wrap with context via `fmt.Errorf`.
- **Secrets:** never log auth headers, tokens, or connection strings — keep the
  generated-SDK logging patch intact.

## Repository layout

See [ARCHITECTURE.md](./ARCHITECTURE.md#layered-design) for the annotated tree.
