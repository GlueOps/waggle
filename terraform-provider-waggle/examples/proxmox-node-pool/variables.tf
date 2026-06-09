# ---- Waggle (placement oracle) ----

variable "waggle_endpoint" {
  type        = string
  description = "Base URL of the Waggle API, e.g. https://waggle.example.com/api/v1"
}

variable "waggle_token" {
  type        = string
  sensitive   = true
  description = "Waggle bearer token. Sent as 'Authorization: Bearer <token>'."
  default     = null
}

variable "waggle_api_key" {
  type        = string
  sensitive   = true
  description = "Waggle API key. Sent verbatim in the Authorization header. Used if waggle_token is null."
  default     = null
}

variable "datacenter_id" {
  type        = string
  description = "Existing Waggle datacenter UUID to place the pool into."
}

variable "slot_id" {
  type        = string
  description = "Existing Waggle slot (VM size) UUID for the pool."
}

variable "pool_name" {
  type    = string
  default = "workers"
}

variable "desired_count" {
  type    = number
  default = 3
}

# ---- Proxmox (bpg/proxmox) ----

variable "proxmox_endpoint" {
  type        = string
  description = "Proxmox VE API endpoint, e.g. https://pve.example.com:8006/"
}

variable "proxmox_api_token" {
  type        = string
  sensitive   = true
  description = "Proxmox API token in the form 'user@realm!tokenid=uuid'."
}

variable "proxmox_insecure" {
  type    = bool
  default = false
}

# ---- VM shape (applied to every node) ----

variable "template_vm_id" {
  type        = number
  description = "VMID of the Proxmox template to clone for each node."
}

variable "node_name_prefix" {
  type    = string
  default = "waggle-node"
}

variable "cores" {
  type    = number
  default = 4
}

variable "memory_mb" {
  type    = number
  default = 8192
}

variable "disk_gb" {
  type    = number
  default = 40
}

variable "datastore_id" {
  type    = string
  default = "local-lvm"
}

variable "network_bridge" {
  type    = string
  default = "vmbr0"
}

variable "ssh_public_keys" {
  type    = list(string)
  default = []
}
