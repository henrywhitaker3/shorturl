-- create "urls" table
CREATE TABLE "public"."urls" (
  "id" uuid NOT NULL,
  "alias" text NOT NULL,
  "url" text NOT NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_urls_alias" to table: "urls"
CREATE UNIQUE INDEX "idx_urls_alias" ON "public"."urls" ("alias");
