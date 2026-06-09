terraform {
  required_version = ">= 1.6"

  required_providers {
    # The placement oracle. Declares the pool; computes placements.
    waggle = {
      source = "registry.terraform.io/glueops/waggle"
    }

    # Creates the actual VMs on Proxmox VE.
    proxmox = {
      source  = "bpg/proxmox"
      version = ">= 0.66"
    }
  }
}
