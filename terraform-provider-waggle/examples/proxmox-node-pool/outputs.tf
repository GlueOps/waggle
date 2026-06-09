output "pool_id" {
  description = "ID of the Waggle pool."
  value       = waggle_pools.this.id
}

output "placements" {
  description = "Placements the oracle computed for the pool."
  value       = local.placements
}

output "nodes" {
  description = "Created VMs keyed by placement id."
  value = {
    for id, vm in proxmox_virtual_environment_vm.node : id => {
      name       = vm.name
      vmid       = vm.vm_id
      hypervisor = vm.node_name
      ipv4       = try(vm.ipv4_addresses, null)
    }
  }
}
