# waggle_slots — look up a slot (t-shirt-size VM template) by name or by id.
# Exactly one of name or id must be set; name is unique within the tenant.

# Look up by name — handy for referencing a slot created elsewhere.
data "waggle_slots" "small" {
  name = "small"
}

# Or look up by id.
data "waggle_slots" "by_id" {
  id = "00000000-0000-0000-0000-000000000000"
}

output "small_slot" {
  description = "Resolved slot id and sizing."
  value = {
    id      = data.waggle_slots.small.id
    vcpu    = data.waggle_slots.small.vcpu
    ram_gb  = data.waggle_slots.small.ram_gb
    disk_gb = data.waggle_slots.small.disk_gb
  }
}
