# waggle_pools — a node pool. Creating it asks the oracle to place
# `desired_count` VMs across hypervisors (anti-affinity spread, all-or-nothing).
# The pool entity round-trips cleanly (create/read/update all return PoolView);
# the resulting placements (hypervisor + vmid) are read via the
# waggle_pool_placements data source — see ../../proxmox-node-pool.
#
# id/created_at/updated_at are computed (server-assigned).

resource "waggle_pools" "workers" {
  name          = "workers"
  datacenter_id = waggle_datacenters.primary.id
  slot_id       = waggle_slots.small.id
  desired_count = 3

  # Optional free-form JSON metadata (string-encoded in the schema).
  metadata = jsonencode({ team = "platform", env = "prod" })
}

output "workers_pool_id" {
  value = waggle_pools.workers.id
}
