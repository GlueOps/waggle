# Database Schema

Waggle uses two schemas. The **control** schema lives in one shared database;
the **tenant** schema is instantiated once per organization in a database named
`tenant_<org-uuid-without-dashes>`. Models are defined as GORM structs under
`internal/models/{control,tenant}/`; the SQL is managed by Goose migrations
under `internal/migrations/{control,tenant}/` (generated via Atlas diff).

Migrations are scoped — see [DEVELOPMENT.md](./DEVELOPMENT.md#migrations).

## Control plane

Identity, tenancy, sessions, and audit. Shared across all organizations.

### `accounts`
Root login identity, independent of any organization.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `display_name` | |
| `password_hash` | bcrypt; never returned |
| `is_active` | |
| `last_login_at` | nullable |
| `created_at` / `updated_at` | |

### `account_emails`
One account may have several emails. A partial unique index enforces a single
primary per account.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `account_id` | uuid, FK → accounts |
| `email` | unique |
| `is_primary` | unique per account where `is_primary` is true |
| `verified_at` | nullable |

### `users`
Membership of an account in an organization, with a role. One account can be a
member of many orgs.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `account_id` | uuid, FK → accounts |
| `organization_id` | uuid, FK → organizations |
| `role` | `owner` \| `admin` \| `member` (default `member`) |
| `is_active` | |
| `last_login_at` | nullable |

Unique on `(account_id, organization_id)`.

### `organizations`
The tenant container and the anchor for tenant-DB encryption.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `name` | |
| `slug` | unique |
| `domain` | unique, optional |
| `status` | `pending` → `active` → `destroying` → `deleted` (indexed) |
| `connection_string` | **plaintext** DSN of the tenant DB (intentional) |
| `encrypted_tenant_key` / `tenant_key_iv` / `tenant_key_tag` | wrapped per-org DEK |
| `metadata` | JSONB |

Lifecycle is driven by River workers; `connection_string` and the wrapped DEK
are populated when provisioning completes.

### `token_sessions`
One row per authenticated refresh-token session.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `account_id` / `organization_id` | uuid |
| `refresh_token_hash` | unique; the plaintext refresh token is never stored |
| `user_agent` / `ip_address` | captured for security |
| `expires_at` | indexed |
| `revoked_at` | nullable; set on logout |

### `org_api_keys`
Long-lived org credentials for automation (e.g. Terraform).

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `organization_id` | uuid |
| `name` | |
| `token_hash` | unique; SHA-256 of the `wgl_…` secret |
| `prefix` | first chars, shown in UI to identify a key |
| `created_by_account_id` | nullable |
| `last_used_at` / `expires_at` / `revoked_at` | nullable |

A key is valid when `revoked_at IS NULL` and (`expires_at IS NULL` or in the
future). The plaintext secret is shown once at creation and never recoverable.

### `auth_audit_events`
Append-only security log.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `organization_id` / `user_id` | nullable, indexed |
| `event` | indexed, e.g. `signup`, `login`, `password_reset` |
| `outcome` | `success` \| `failure` |
| `ip_address` / `user_agent` | |
| `metadata` | JSONB |

## Tenant plane

The Proxmox placement ledger, isolated per organization.

### `datacenters`
A Proxmox cluster endpoint.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `name` | |
| `url` | PVE API URL |
| `encrypted_token_key` / `token_key_iv` / `token_key_tag` | Proxmox token encrypted with the tenant DEK |
| `insecure_skip_verify` | bool, default false |

The token is write-only over the API; responses expose only `HasToken`.

### `hypervisors`
A single Proxmox node and its capacity. Unique on `(datacenter_id, name)`.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `datacenter_id` | uuid, FK → datacenters |
| `name` | |
| `cpu_total` / `cpu_reserved` / `cpu_used` | cores |
| `ram_gb_total` / `ram_gb_reserved` / `ram_gb_used` | GB |
| `disk_gb_total` / `disk_gb_reserved` / `disk_gb_used` | GB |
| `schedulable` | bool, default true |
| `last_synced_at` | set by discovery |

**Capacity model**

```
Bookable = Total − Reserved − Used
```

- `Used` — sum of existing guest allocations; refreshed by discovery.
- `Reserved` — operator headroom; **preserved** across re-discovery and seeded
  from `RESERVE_CPU` / `RESERVE_RAM_GB` / `RESERVE_DISK_GB` on first sight.
- `Schedulable` — gate for placement; flip to false to drain a node for
  maintenance without deleting it.

The API reports the derived `*Bookable` values to clients.

### `slots`
A reusable VM "t-shirt size". Unique on `name`.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `name` | unique, e.g. `small` |
| `vcpu` / `ram_gb` / `disk_gb` | |

### `pools`
A managed group of identical VMs in a datacenter.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `datacenter_id` | uuid, FK → datacenters |
| `slot_id` | uuid, FK → slots |
| `name` | indexed (not unique — see invariants) |
| `desired_count` | target number of VMs |
| `metadata` | JSONB, optional |

### `placements`
One row per VM: the realized assignment of a pool VM to a hypervisor.

| Field | Notes |
|-------|-------|
| `id` | uuid, PK |
| `pool_id` | uuid, FK → pools |
| `hypervisor_id` | uuid, FK → hypervisors |
| `vmid` | nullable Proxmox VM ID, backfilled after provisioning |

## Schema invariants worth knowing

A few rules look like oversights but are deliberate:

- `pools.name` is **not unique** — multiple pools may share a name.
- `placements.vmid` is **nullable** — Waggle records intent before a VM exists;
  the downstream provisioner backfills the real Proxmox ID.
- The Proxmox token on `datacenters` is effectively **write-only** via the API.
- `organizations.connection_string` is **plaintext by design**.
