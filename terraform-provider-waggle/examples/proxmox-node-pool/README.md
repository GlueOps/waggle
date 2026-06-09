# Waggle pool → Proxmox node pool

End-to-end: **Waggle decides where each VM goes, `bpg/proxmox` builds it.**

1. `waggle_pools.this` declares the pool. Creating/resizing it makes the oracle
   place `desired_count` VMs across hypervisors (anti-affinity, all-or-nothing)
   and reserve a `vmid` for each.
2. `data.waggle_pool_placements.this` reads the pool's placements natively —
   each carries `hypervisor_name` and `vmid`.
3. `proxmox_virtual_environment_vm.node` is created `for_each` placement, with
   `node_name = hypervisor_name` and `vm_id = vmid`.

```
waggle_pools ──▶ data.waggle_pool_placements ──▶ [{hypervisor_name, vmid}, ...]
                                                        │ for_each
                                                        ▼
                                            proxmox_virtual_environment_vm
```

## The `waggle_pool_placements` data source

Placements are produced by a pool — there is no create or read-by-id for an
individual placement — so they aren't a managed resource. They're exposed
read-only via the `waggle_pool_placements` data source (takes `pool_id`,
returns a typed `placements` list of `{id, hypervisor_id, hypervisor_name,
vmid, created_at}`). It depends on `waggle_pools.this.id`, so it re-reads
whenever the pool changes. No third-party providers are required — only
`waggle` and `bpg/proxmox`.

> This data source is a hand-authored overlay (the OpenAPI Generator can't
> model a list of nested objects). It's injected into the provider on every
> regen from `cmd/overlays/pool_placements_data_source.go.tmpl` in the waggle
> repo, so it persists across `generate`.

## Usage

```sh
cp terraform.tfvars.example terraform.tfvars
# edit terraform.tfvars (datacenter_id, slot_id, proxmox creds, template id)
terraform init
terraform plan
terraform apply
```

## Requirements

- A reachable Waggle API with an existing datacenter (`datacenter_id`) and slot
  (`slot_id`).
- A Proxmox VE cluster, an API token, and a clonable template (`template_vm_id`)
  whose VMID space does not collide with the VMIDs Waggle reserves.

See `../README.md` for how attribute roles (Required vs Computed) are derived.
