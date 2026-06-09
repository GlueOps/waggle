# waggle_slots — a "t-shirt size" VM template (vcpu / ram / disk) used when
# sizing pools. id/created_at/updated_at are computed (server-assigned).

resource "waggle_slots" "small" {
  name    = "small"
  vcpu    = 2
  ram_gb  = 8
  disk_gb = 40
}

output "small_slot_id" {
  value = waggle_slots.small.id
}
