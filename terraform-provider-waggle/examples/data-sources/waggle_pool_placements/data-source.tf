# waggle_pool_placements — read the placements the oracle computed for a pool.
# Returns a typed list of { id, hypervisor_id, hypervisor_name, vmid, created_at },
# the input you feed to a VM provisioner (see ../../proxmox-node-pool).

data "waggle_pool_placements" "workers" {
  pool_id = "00000000-0000-0000-0000-000000000000" # waggle_pools.workers.id
}

output "placement_targets" {
  description = "hypervisor + vmid per placement."
  value = [
    for p in data.waggle_pool_placements.workers.placements : {
      hypervisor = p.hypervisor_name
      vmid       = p.vmid
    }
  ]
}
