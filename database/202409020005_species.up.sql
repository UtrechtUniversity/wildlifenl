CREATE TABLE IF NOT EXISTS "species" (
	"id" UUID NOT NULL DEFAULT gen_random_uuid(),
	"name" VARCHAR NOT NULL,
	"commonNameNL" VARCHAR NOT NULL,
	"commonNameEN" VARCHAR NOT NULL,
	PRIMARY KEY ("id")
);