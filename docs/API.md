# API Reference

The API is built with [Huma](https://github.com/danielgtaylor/huma) on Chi and
is fully described by the generated OpenAPI spec. The summaries below are a map;
the spec is the source of truth.

- **Base path:** `/api/v1` (configurable via `BASE_PATH`)
- **Interactive docs:** `GET /api/v1/docs`
- **OpenAPI spec:** `GET /api/v1/openapi` (also committed at
  [`docs/openapi.json`](./openapi.json))
- **Health:** `GET /api/v1/health` → `{"status":"ok"}`

## Authentication

Two credential types, selected per route:

| Type | Header | Used by | Routes |
|------|--------|---------|--------|
| **User JWT** | `Authorization: Bearer <jwt>` | humans (UI, login flow) | everything |
| **Org API key** | `Authorization: Bearer wgl_<secret>` | automation (Terraform) | fleet/tenant resources only |

- JWTs are HS256, signed with `JWT_SECRET`, carrying the account and current
  organization. Obtain one from `/auth/login` (or `/auth/switch` to change org).
- API keys are minted with `POST /api-keys`, scoped to one organization, and
  identified by the `wgl_` prefix. They **cannot** perform human-only actions:
  organization/member administration and minting further API keys require a JWT.

The middleware attaches a `Principal` of `user` or `api_key` to each request.

## Auth & account — `/auth`

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/auth/signup` | Create account + org + first user; enqueue tenant provisioning |
| POST | `/auth/login` | Exchange credentials for an access + refresh token pair |
| POST | `/auth/refresh` | Rotate the refresh token; issue a new pair |
| POST | `/auth/logout` | Revoke a refresh-token session (idempotent) |
| POST | `/auth/verify-email` | Consume an email verification token (idempotent) |
| GET | `/auth/me` | Current account, emails, memberships, active org |
| POST | `/auth/switch` | Issue a token pair for another org the account belongs to |
| POST | `/auth/accept-invite` | Accept an invite, set password if new, sign in |

## Organizations & members — `/organizations`

User JWT only — these are human administrative actions.

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/organizations` | Create org (caller becomes owner); enqueue provisioning |
| GET | `/organizations` | List orgs the caller belongs to, with their role |
| GET | `/organizations/{id}` | Get one org |
| PATCH | `/organizations/{id}` | Rename (admin+) |
| DELETE | `/organizations/{id}` | Delete + enqueue teardown (owner only) |
| GET | `/organizations/{id}/members` | List members |
| POST | `/organizations/{id}/members` | Add/invite by email (admin+; owner to grant owner) |
| PATCH | `/organizations/{id}/members/{userId}` | Change role (admin+; owner for owners) |
| DELETE | `/organizations/{id}/members/{userId}` | Remove member (never the last owner) |

## API keys — `/api-keys`

User JWT only (an API key cannot bootstrap more keys).

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/api-keys` | Mint a key; plaintext returned **once**; `expires_in_days` optional |
| GET | `/api-keys` | List keys (metadata only, no secrets) |
| DELETE | `/api-keys/{id}` | Revoke a key (idempotent) |

## Fleet / tenant resources

JWT or org API key. Org scope is taken from the credential.

### Datacenters — `/datacenters`

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/datacenters` | Create (Proxmox token encrypted with tenant DEK) |
| GET | `/datacenters` | List |
| GET | `/datacenters/{id}` | Get |
| PUT | `/datacenters/{id}` | Update (token optional on update) |
| DELETE | `/datacenters/{id}` | Delete |
| POST | `/datacenters/{id}/discover` | Discover hypervisors from Proxmox; `?async` queues a job |

### Slots — `/slots`

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/slots` | Create a VM size |
| GET | `/slots` | List |
| GET | `/slots/{id}` | Get |
| PUT | `/slots/{id}` | Update |
| DELETE | `/slots/{id}` | Delete |

### Hypervisors — `/hypervisors`

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/hypervisors` | Create (manual entry) |
| GET | `/hypervisors` | List; optional `datacenter_id` filter |
| GET | `/hypervisors/{id}` | Get |
| PUT | `/hypervisors/{id}` | Update (preserves `reserved` / `schedulable` on re-discovery) |
| DELETE | `/hypervisors/{id}` | Delete |

Responses report derived bookable capacity (`Total − Reserved − Used`).

### Pools & placements

| Method | Path | Purpose |
|--------|------|---------|
| POST | `/pools` | Create a pool and place its VMs (anti-affinity, all-or-nothing) |
| GET | `/pools` | List |
| GET | `/pools/{id}` | Get |
| PATCH | `/pools/{id}` | Resize (grow: place new all-or-nothing; shrink: LIFO release) |
| DELETE | `/pools/{id}` | Delete and release all placements |
| GET | `/pools/{id}/placements` | Placements for one pool |
| GET | `/placements` | Fleet-wide placements with pool/slot/hypervisor context |
| PATCH | `/placements/{id}` | Backfill the Proxmox `vmid` onto a placement |

## Errors

Errors follow Huma's RFC 7807-style problem shape:

```json
{
  "title": "Unprocessable Entity",
  "status": 422,
  "detail": "validation failed",
  "errors": [
    { "message": "expected required property", "location": "body.name" }
  ]
}
```

A pool create/resize that exceeds bookable capacity returns an error and changes
nothing (all-or-nothing).

## Typical automation flow (Terraform)

1. `POST /auth/login` (human) → `POST /api-keys` → store the `wgl_…` secret.
2. With the API key: `POST /datacenters` (with Proxmox token) →
   `POST /datacenters/{id}/discover` to learn hypervisors.
3. `POST /slots` to define sizes, then `POST /pools` with a `desired_count`.
4. Read `GET /placements`, provision the VMs in Proxmox, and
   `PATCH /placements/{id}` to record each `vmid`.

See the [Terraform provider](https://github.com/GlueOps/terraform-provider-waggle)
for the managed version of this flow.
