# Waggle Documentation

Reference docs for **Waggle** — a placement oracle and ledger for Proxmox
infrastructure. Start with the [root README](../README.md) for the project
overview and quick start, then dig into the topics below.

| Doc | What's inside |
|-----|---------------|
| [Architecture](./ARCHITECTURE.md) | Control vs. tenant planes, request flow, jobs, encryption, placement algorithm |
| [Database](./DATABASE.md) | Control & tenant schema reference, capacity model, lifecycle states |
| [API](./API.md) | Endpoint reference, auth model, error format |
| [Configuration](./CONFIGURATION.md) | Every environment variable with defaults |
| [Development](./DEVELOPMENT.md) | Local setup, code generation, migrations, testing, contributing |
| [`openapi.json`](./openapi.json) | Generated OpenAPI 3.1 spec (served live at `/api/v1/docs`) |

## What Waggle is (and isn't)

Waggle decides **where** VMs should live across a fleet of Proxmox hypervisors
and records that decision as an authoritative ledger of pools and placements. It
is named after the "waggle dance" bees use to communicate direction and timing.

Waggle **does not create or boot VMs**. It is the source of truth for placement
state; actual provisioning is performed by downstream consumers — primarily the
[Terraform provider](https://github.com/GlueOps/terraform-provider-waggle) — that read placements from
Waggle and apply them to Proxmox.
