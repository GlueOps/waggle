###############################################################################
# Waggle decides WHERE; bpg/proxmox makes it real.
#
#   1. waggle_pools declares the pool. Creating/resizing it makes the oracle
#      place `desired_count` VMs across hypervisors (anti-affinity, all-or-
#      nothing) and reserve a vmid for each.
#   2. The placements (hypervisor_name + vmid) are read from the API and one
#      VM is cloned per placement, pinned to that hypervisor and vmid.
#
# Placements are a product of the pool, read via the native
# waggle_pool_placements data source. It depends on waggle_pools.this.id, so it
# re-reads whenever the pool changes. No third-party (http) provider needed.
###############################################################################

provider "waggle" {
  endpoint = var.waggle_endpoint
  token    = var.waggle_token
  api_key  = var.waggle_api_key
}

provider "proxmox" {
  endpoint  = var.proxmox_endpoint
  api_token = var.proxmox_api_token
  insecure  = var.proxmox_insecure
}

# --- 1. Declare the pool (oracle places its VMs) -----------------------------

resource "waggle_pools" "this" {
  name          = var.pool_name
  datacenter_id = var.datacenter_id
  slot_id       = var.slot_id
  desired_count = var.desired_count
}

# --- 2. Read the placements the oracle computed for this pool ----------------

data "waggle_pool_placements" "this" {
  pool_id = waggle_pools.this.id
}

locals {
  # Key placements by id so adds/removes map to stable VM addresses.
  placements = { for p in data.waggle_pool_placements.this.placements : p.id => p }
}

# --- 3. Realize one VM per placement on its assigned hypervisor --------------

resource "proxmox_virtual_environment_vm" "node" {
  for_each = local.placements

  node_name = each.value.hypervisor_name
  vm_id     = each.value.vmid

  name        = "${var.node_name_prefix}-${each.value.vmid}"
  description = "Managed by Terraform. Placement ${each.key} from Waggle pool ${waggle_pools.this.id}."
  tags        = ["waggle", "pool-${waggle_pools.this.name}"]

  clone {
    vm_id = var.template_vm_id
  }

  cpu {
    cores = var.cores
    type  = "host"
  }

  memory {
    dedicated = var.memory_mb
  }

  disk {
    datastore_id = var.datastore_id
    interface    = "scsi0"
    size         = var.disk_gb
  }

  network_device {
    bridge = var.network_bridge
  }

  initialization {
    ip_config {
      ipv4 {
        address = "dhcp"
      }
    }

    dynamic "user_account" {
      for_each = length(var.ssh_public_keys) > 0 ? [1] : []
      content {
        username = "ubuntu"
        keys     = var.ssh_public_keys
      }
    }
  }

  operating_system {
    type = "l26"
  }

  lifecycle {
    ignore_changes = [clone]
  }
}
