CREATE TABLE IF NOT EXISTS "user" (
	"id" UUID NOT NULL DEFAULT gen_random_uuid(),
	"name" VARCHAR NULL DEFAULT NULL,
	"email" VARCHAR NOT NULL,
	PRIMARY KEY ("id"),
	UNIQUE ("email")
);