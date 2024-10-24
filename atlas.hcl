data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "ariga.io/atlas-provider-gorm",
    "load",
    "--path", "./models",
    "--dialect", "postgres"
  ]
}

data "external" "dot_env" {
  program = [
    "go",
    "run",
    "./scripts/db_url.go",
    "-env=${atlas.env}"
  ]
}

locals {
  dot_env = jsondecode(data.external.dot_env)
  db_url = "${local.dot_env.db_url}"
}

env {
  name = atlas.env // must be either: test, dev, prod
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev?search_path=public"
  url = "${local.db_url}"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}