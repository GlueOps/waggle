# Architecture

Waggle is a multi-tenant control plane for Proxmox VM placement. It is a single
Go binary (`waggle`) that runs in two roles — an HTTP API server and a
background job worker — over a PostgreSQL control database and one PostgreSQL
database per tenant organization.

```
                        ┌────────────────────────────────────────────┐
                        │                  waggle                     │
   Terraform provider   │   serve  ──HTTP/JSON──►  Huma + Chi router  │
   UI (embedded SPA)  ──┼──►          │                               │
   curl / SDKs          │            ▼                                │
                        │      services (auth, fleet, org, apikey)    │
                        │       │            │            │           │
                        │       ▼            ▼            ▼           │
                        │  ControlDB    TenantManager   River queue   │
                        └───────┼────────────┼─────────────┼─────────┘
                                │            │             │
                        ┌───────▼───┐  ┌─────▼──────┐  ┌───▼────────┐
                        │ control   │  │ tenant_<id>│  │  worker    │
                        │ Postgres  │  │ Postgres   │  │ (River)    │
                        └───────────┘  └────────────┘  └─────┬──────┘
                                                             │
                                                       ┌─────▼──────┐
                                                       │  Proxmox   │
                                                       │  (PVE API) │
                                                       └────────────┘
```

## Two planes

**Control plane** — one shared, always-present PostgreSQL database. Holds
platform identity and tenancy: accounts, organizations, users/memberships,
token sessions, organization API keys, and the auth audit log. It also stores,
per organization, the **wrapped tenant data-encryption key** and the connection
string for that org's tenant database.

**Tenant plane** — one PostgreSQL database per organization, named
`tenant_<org-uuid-without-dashes>`, created on demand when the org is
provisioned. Holds the infrastructure ledger: datacenters, hypervisors, slots
(VM sizes), pools (desired VM groups), and placements (actual VM assignments).

This split keeps tenant infrastructure data physically isolated while a single
control database manages identity and lifecycle.

## Process roles

Both roles are subcommands of the same binary and can be scaled independently.

- **`waggle serve`** — boots config, opens the control DB, constructs the
  `TenantManager`, the River client, and the service layer, then mounts the
  Huma/Chi API plus the frontend (per `FRONTEND_MODE`). Listens on
  `BIND_HOST:BIND_PORT`, graceful shutdown on SIGINT/SIGTERM.
- **`waggle worker`** — registers River workers (tenant provisioner, tenant
  destroyer, hypervisor discovery) and drains the job queue. Decouples slow or
  failure-prone background work from API availability.

See [DEVELOPMENT.md](./DEVELOPMENT.md) for the full CLI surface
(`migrate`, `encrypt`, `generate`, `version`).

## Layered design

```
cmd/                 Cobra commands: serve, worker, migrate, encryption, generate, version
internal/
  app/               Deps container — every wired dependency passed to commands
  config/            Env/Viper-backed config + validation
  database/          Control DB open + TenantManager (per-org connection + DEK unwrap)
  api/               Huma routes, middleware (auth), docs, frontend mounting
  service/           Business logic: auth, token, org, fleet, placement, apikey, email, policy
  repo/              GORM repositories
  jobs/              River workers + enqueuer + tenant DSN helpers
  proxmox/           Proxmox (PVE) API client
  models/control/    Control-plane GORM models
  models/tenant/     Tenant-plane GORM models
  migrations/        Goose migrations (control/ and tenant/), embedded
  utils/             crypto (AES-GCM envelope), password hashing, normalization
ui/                  Vite/React SPA, embedded into the binary via go:embed
sdk/                 Generated Go + TypeScript clients
terraform-provider-waggle/  Generated + hand-overlaid Terraform provider
```

Dependencies are assembled once in `internal/app/deps.go` (the `Deps` struct)
and threaded through commands. Routes are registered conditionally on which
services are present, so the same wiring supports both the live server and the
offline spec-generation path used by `waggle generate sdk`.

## Request flow

1. A client calls the API with a **Bearer JWT** (human session) or an
   **organization API key** (`wgl_` prefix, for automation).
2. Middleware (`internal/api/middleware.go`) authenticates the request:
   - `RequireAuth` — JWT only; used for human-only actions (org/member admin,
     API key minting).
   - `RequireBearer` — JWT *or* API key; used for fleet/tenant resources. The
     API-key path synthesizes claims scoped to the key's organization.
   Both attach the resolved claims and a `Principal` (`user` | `api_key`) to the
   request context, and capture client IP / User-Agent for sessions and audit.
3. The handler calls a service. Tenant-scoped services resolve the caller's
   tenant DB through `TenantManager.For(ctx, orgID)`.
4. Mutations that must happen out-of-band (tenant DB create/drop, hypervisor
   discovery) are enqueued onto River rather than run inline.

## Multi-tenancy & the TenantManager

`internal/database/database.go` owns tenant resolution:

- `For(ctx, orgID)` returns a cached GORM handle to the org's tenant DB,
  reopening on miss. It errors with `ErrTenantNotProvisioned` (no connection
  string yet) or `ErrTenantNotActive` (org not `active`).
- `TenantDEK(ctx, orgID)` loads the org, then unwraps its data-encryption key
  with the master key (AES-GCM). The plaintext DEK is returned to the caller and
  never cached.
- `Forget(orgID)` evicts and closes a cached handle — called before a tenant DB
  is dropped.

## Encryption (envelope: KEK → DEK → field)

`internal/utils/crypto.go` implements AES-256-GCM envelope encryption.

```
ENCRYPTION_MASTER_KEY  (KEK, 32 bytes, base64 in env)
        │  wraps
        ▼
Organization.EncryptedTenantKey (+ IV + tag)   ← per-org DEK, 32 bytes, in control DB
        │  decrypts to plaintext DEK at use time
        ▼
Field-level ciphertext in the tenant DB         ← e.g. Datacenter Proxmox token
```

- Each org gets a random 32-byte **DEK** at provisioning; it is stored only in
  wrapped form on the `Organization` row.
- Field secrets (currently the Proxmox API token on `Datacenter`) are encrypted
  with the DEK. The API never returns the token — only a `HasToken` boolean.
- `EncryptAESGCM` returns `(ciphertext, iv, tag)` stored as separate columns;
  `DecryptAESGCM` fails closed if the tag does not verify.
- `waggle encrypt rotate-master` re-wraps every org's DEK under a new KEK
  without touching field ciphertext. Processes must restart with the new key.

> **Intentionally plaintext:** `Organization.ConnectionString` is stored
> unencrypted in the control DB for operational simplicity.

## Background jobs (River)

River is a PostgreSQL-backed queue. Jobs are inserted via the `Enqueuer` facade
(`internal/jobs/enqueuer.go`) and processed by `waggle worker`:

| Worker | Trigger | Work |
|--------|---------|------|
| `TenantProvisioner` | signup / org create | create tenant DB, run tenant migrations, generate + wrap DEK, set org `active` |
| `TenantDestroyer` | org delete | forget cached handle, drop tenant DB, set org `deleted` |
| `HypervisorDiscovery` | `POST /datacenters/{id}/discover?async` | query Proxmox, upsert hypervisors, refresh used capacity |

Provisioner and destroyer are idempotent and bail out if the org status has
moved past the expected state. Discovery coalesces by `(org, datacenter)` to
collapse rapid repeat requests.

## Placement algorithm

A **pool** declares a `DesiredCount` of identical VMs (one `Slot` size) within a
datacenter. `PlacementService` (`internal/service/placement.go`) plans where
those VMs go:

- **Anti-affinity spread** — greedily assign each VM to the hypervisor currently
  holding the fewest VMs from this pool, tie-breaking on most remaining bookable
  CPU. This spreads a pool across nodes for fault tolerance.
- **All-or-nothing** — if the full requested count does not fit the fleet's
  bookable capacity, the operation fails with a capacity error and places
  nothing.
- **Resize** — growing a pool places only the new VMs (all-or-nothing); shrinking
  removes placements LIFO.

Bookable capacity per hypervisor is `Total − Reserved − Used`, where `Used` is
refreshed by discovery and `Reserved` is operator-controlled headroom preserved
across discoveries (see [DATABASE.md](./DATABASE.md)).

Placements carry an optional Proxmox `VMID`, **backfilled** by the downstream
provisioner via `PATCH /placements/{id}` once the VM actually exists.
