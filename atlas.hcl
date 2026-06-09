data "external_schema" "control" {
  program = ["go", "run", "./tools/schema/control"]
}

data "external_schema" "tenant" {
  program = ["go", "run", "./tools/schema/tenant"]
}

env "control" {
  src = data.external_schema.control.url
  # An ephemeral dev database used by Atlas to calculate diffs safely
  dev = "docker://postgres/18/dev?search_path=public"
  migration {
    dir    = "file://internal/migrations/control"
    format = goose
  }
}

env "tenant" {
  src = data.external_schema.tenant.url
  dev = "docker://postgres/18/dev?search_path=public"
  migration {
    dir    = "file://internal/migrations/tenant"
    format = goose
  }
}