# Configuration

Waggle is configured entirely through environment variables (loaded via Viper;
a local `.env` is read if present). Copy [`.env.example`](../.env.example) to
`.env` to get started.

**Required:** `DATABASE_URL`, `JWT_SECRET`, `BASE_URL`.
`FRONTEND_MODE` must be one of `embed | proxy | none`.

## Server

| Variable | Default | Purpose |
|----------|---------|---------|
| `ENV` | — | `dev` / `staging` / `prod`; affects HSTS behavior |
| `BIND_HOST` | `127.0.0.1` | HTTP bind address |
| `BIND_PORT` | `8080` | HTTP bind port |
| `BASE_URL` | `http://localhost:8080` | Public base URL used in email links |
| `BASE_PATH` | `/api/v1` | API route prefix |
| `FRONTEND_MODE` | `embed` | `embed` (built SPA), `proxy` (Vite dev), `none` (API only) |
| `VITE_DEV_URL` | `http://localhost:5173` | Vite dev server target when `FRONTEND_MODE=proxy` |

## Databases

| Variable | Default | Purpose |
|----------|---------|---------|
| `DATABASE_URL` | **required** | Control DB DSN (`postgres://…`) |
| `ADMIN_DATABASE_URL` | derived (`/postgres`) | Superuser DSN used to create/drop tenant DBs |
| `RIVER_DATABASE_URL` | `DATABASE_URL` | DSN for the River job queue |

## Encryption

| Variable | Default | Purpose |
|----------|---------|---------|
| `ENCRYPTION_MASTER_KEY` | — | Base64 of 32 random bytes (KEK). Required to use the tenant manager. Generate with `go run . encrypt generate-master` |

Without a master key the server runs but cannot resolve tenant databases or
decrypt tenant secrets. See [ARCHITECTURE.md](./ARCHITECTURE.md#encryption-envelope-kek--dek--field).

## JWT

| Variable | Default | Purpose |
|----------|---------|---------|
| `JWT_SECRET` | **required** | HMAC-SHA256 signing key |
| `JWT_ISSUER` | — | `iss` claim |
| `JWT_AUDIENCE` | `waggle` | `aud` claim |
| `JWT_ACCESS_TTL_MIN` | — | Access-token lifetime (minutes) |
| `JWT_REFRESH_TTL_HOUR` | — | Refresh-token lifetime (hours) |

## SMTP

Transactional email (verification, invites). If `SMTP_SERVER` is unset, Waggle
logs emails to stdout instead of sending (dev mode). The bundled
`docker-compose` runs [Mailpit](https://github.com/axllent/mailpit) on
`localhost:1025` (web UI at `:8025`).

| Variable | Default | Purpose |
|----------|---------|---------|
| `SMTP_SERVER` | — | SMTP host; unset → log-only dev sender |
| `SMTP_PORT` | `1025` | SMTP port |
| `SMTP_USER` | — | Username |
| `SMTP_PASSWORD` | — | Password |
| `SMTP_FROM` | `SMTP_USER` | `From:` header |

## Capacity reservation

Headroom held back the first time a hypervisor is discovered. Operator changes
to a hypervisor's reserved values are preserved on later re-discovery.

| Variable | Default | Purpose |
|----------|---------|---------|
| `RESERVE_CPU` | `0` | Cores reserved on newly discovered nodes |
| `RESERVE_RAM_GB` | `2` | RAM (GB) reserved on newly discovered nodes |
| `RESERVE_DISK_GB` | `10` | Disk (GB) reserved on newly discovered nodes |

## Reserved for future use

These are read but not yet wired into active features:

| Variable | Purpose |
|----------|---------|
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` / `OAUTH_REDIRECT_BASE` | Google OAuth2 sign-in |
| `WEBAUTHN_RP_ID` / `WEBAUTHN_RP_ORIGIN` / `WEBAUTHN_RP_NAME` | WebAuthn / passkeys |

## Docker Compose variables

Used only by `docker-compose.yml` for the local Postgres container:

| Variable | Purpose |
|----------|---------|
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | Local Postgres credentials/database |
| `LOCAL_POSTGRES_PORT` / `LOCAL_BOUNCER_PORT` | Host port mappings |
