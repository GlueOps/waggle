# Waggle provider examples

> **Durability:** everything under `examples/` is hand-authored and is
> **preserved across `waggle generate terraform openapi-generator`**. The
> generator's clean step (`RemoveAll`) snapshots this directory and overlays it
> back after regeneration (see `snapshotDir`/`copyTreeOverwrite` in
> `cmd/generate.go`), so these files survive a regen even though the rest of the
> provider is regenerated from scratch. The generator still writes its own
> `examples/provider/provider.tf`, which the restore then overwrites with the
> richer version here.

## Layout

| Path | What it shows |
|---|---|
| `provider/provider.tf` | Provider config: endpoint + auth (token or api_key) |
| `resources/waggle_slots/` | A VM size (vcpu/ram/disk) |
| `resources/waggle_datacenters/` | A Proxmox cluster (read/track; token is write-only) |
| `resources/waggle_pools/` | A node pool — the core placement request |
| `resources/waggle_organizations/` | A tenant/org |
| `data-sources/waggle_pool_placements/` | Read a pool's placements (hypervisor + vmid) |
| `proxmox-node-pool/` | **End-to-end: `waggle_pools` → `waggle_pool_placements` → bpg/proxmox VMs** |

## Attribute roles

The OpenAPI Generator marks *every* attribute `Required`. A post-generation
pass (`patchResourceSchemaRoles` in `cmd/generate.go`) corrects this from the
spec's create/update bodies, so each resource now has sensible roles:

- **Required** — real inputs (`name`, `vcpu`, `datacenter_id`, `desired_count`, …).
- **Optional** / **Optional+Computed** — optional inputs (`metadata`,
  `insecure_skip_verify`, `schedulable`).
- **Computed** — server-assigned, read-only fields (`id`, `created_at`,
  `updated_at`, `has_token`, the hypervisor `*_used`/`*_bookable` capacity, …).

So you only set inputs; computed fields are read back as attributes (e.g.
`waggle_pools.workers.id`) with no placeholder values and no perpetual diff.

> Caveat: `datacenters` cannot set the Proxmox token (it's write-only on the
> API and absent from `DatacenterView`); configure it out-of-band. And on an
> in-place update, computed fields show `(known after apply)` in the plan —
> cosmetic, not drift.

## Usage

The provider must be built and installed locally first
(`make build && make install` in the provider root), with a `dev_overrides`
block pointing Terraform at the local binary. Then, in an example directory:

```sh
terraform plan
```

(`proxmox-node-pool/` is a full root module — copy `terraform.tfvars.example`
to `terraform.tfvars` first.)
