# waggle_datacenters — a Proxmox cluster Waggle places into.
#
# NOTE: the Proxmox API token is write-only on the API and is not part of
# DatacenterView, so this generated resource cannot set it. Configure the
# datacenter's token out-of-band (API/UI); the computed `has_token` reflects
# whether one is stored. id/created_at/updated_at/has_token are computed.

resource "waggle_datacenters" "primary" {
  name = "primary"
  url  = "https://pve.example.com:8006/api2/json"

  insecure_skip_verify = false # optional
}

output "primary_datacenter_id" {
  value = waggle_datacenters.primary.id
}

output "primary_has_token" {
  value = waggle_datacenters.primary.has_token
}
