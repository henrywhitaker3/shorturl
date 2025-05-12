-- create "clicks" table
CREATE TABLE "public"."clicks" (
  "id" uuid NOT NULL,
  "url_id" uuid NOT NULL,
  "ip" text NOT NULL,
  "clicked_at" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_clicks_url_id" FOREIGN KEY ("url_id") REFERENCES "public"."urls" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_clicks_url_id" to table: "clicks"
CREATE INDEX "idx_clicks_url_id" ON "public"."clicks" ("url_id");
