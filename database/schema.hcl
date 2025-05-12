schema "public" {}

table "urls" {
  schema = schema.public

  column "id" {
    type = uuid
    null = false
  }

  column "alias" {
    type = text
    null = false
  }

  column "url" {
    type = text
    null = false
  }

  column "domain" {
    type = text
    null = false
  }

  primary_key {
    columns = [column.id]
  }
  index "idx_urls_alias" {
    columns = [column.alias]
    unique  = true
  }
  foreign_key "fk_urls_alias" {
    columns     = [column.alias]
    ref_columns = [table.aliases.column.alias]
    on_delete   = CASCADE
  }
}

table "aliases" {
  schema = schema.public

  column "alias" {
    type = text
    null = false
  }

  column "used" {
    type    = boolean
    null    = false
    default = false
  }

  primary_key {
    columns = [column.alias]
  }
  index "idx_aliases_used" {
    columns = [column.used]
  }
}
