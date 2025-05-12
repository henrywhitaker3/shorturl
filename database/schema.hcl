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

table "clicks" {
  schema = schema.public

  column "id" {
    type = uuid
    null = false
  }

  column "url_id" {
    type = uuid
    null = false
  }

  column "ip" {
    type = text
    null = false
  }

  column "clicked_at" {
    type = bigint
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "fk_clicks_url_id" {
    columns     = [column.url_id]
    ref_columns = [table.urls.column.id]
    on_delete   = CASCADE
  }
  index "idx_clicks_url_id" {
    columns = [column.url_id]
  }
}
