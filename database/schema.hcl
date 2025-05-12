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
}
