variable "db_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

env "local" {
  src = "ent://internal/ent/schema"
  dev = "sqlite://ent?mode=memory&cache=shared&_pragma=foreign_keys(1)"
  migration {
    dir = "file://migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "prod" {
  url = var.db_url
  migration {
    dir = "file://migrations"
  }
}
