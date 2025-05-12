-- reverse: drop "alias_buffer" table
CREATE TABLE "public"."alias_buffer" (
  "alias" text NOT NULL,
  PRIMARY KEY ("alias")
);
-- reverse: modify "urls" table
ALTER TABLE "public"."urls" DROP CONSTRAINT "fk_urls_alias";
-- reverse: create "aliases" table
DROP TABLE "public"."aliases";
