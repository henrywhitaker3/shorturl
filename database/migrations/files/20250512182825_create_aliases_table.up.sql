-- create "aliases" table
CREATE TABLE "public"."aliases" (
  "alias" text NOT NULL,
  "used" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("alias")
);
-- modify "urls" table
ALTER TABLE "public"."urls" ADD CONSTRAINT "fk_urls_alias" FOREIGN KEY ("alias") REFERENCES "public"."aliases" ("alias") ON UPDATE NO ACTION ON DELETE CASCADE;
-- drop "alias_buffer" table
DROP TABLE "public"."alias_buffer";
