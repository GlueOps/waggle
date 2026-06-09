terraform {
  required_providers {
    waggle = {
      source = "registry.terraform.io/glueops/waggle"
    }
  }
}

# Waggle is the placement oracle + ledger for a Proxmox fleet. It records
# datacenters, slots (t-shirt VM sizes), and hypervisors, and computes pool
# placements (which hypervisor + which vmid). It does NOT create VMs.
#
# Authentication (Configure() uses the first that is set):
#   - token    -> sent as "Authorization: Bearer <token>"
#   - api_key  -> sent as "Authorization: <api_key>" (raw, no prefix)
#   - username + password are also accepted.
#
# Pass secrets via TF_VAR_/env, not literals.
provider "waggle" {
  # Full base URL of the Waggle API. The provider's built-in default is the
  # bare "/api/v1" path, so set the real host here.
  endpoint = "https://waggle.example.com/api/v1"

  # token   = var.waggle_token
  api_key = var.waggle_api_key
}

variable "waggle_api_key" {
  type        = string
  sensitive   = true
  description = "Waggle API key (sent verbatim in the Authorization header)."
  default     = null
}

variable "waggle_token" {
  type        = string
  sensitive   = true
  description = "Waggle bearer token (sent as 'Authorization: Bearer <token>')."
  default     = null
}
