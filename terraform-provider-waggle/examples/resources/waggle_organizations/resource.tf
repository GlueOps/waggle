# waggle_organizations — an organization (tenant) in the control plane.
# Only `name` is an input; the server assigns id/slug/status/role and (when
# applicable) domain, all exposed as computed attributes.

resource "waggle_organizations" "acme" {
  name = "ACME Corp"
}

output "acme" {
  value = {
    id     = waggle_organizations.acme.id
    slug   = waggle_organizations.acme.slug
    status = waggle_organizations.acme.status
    role   = waggle_organizations.acme.role
  }
}
