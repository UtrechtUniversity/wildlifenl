CREATE TABLE IF NOT EXISTS "credential" (
	"token" UUID NOT NULL DEFAULT gen_random_uuid(),
	"email" VARCHAR NOT NULL,
	"lastLogon" TIMESTAMP NOT NULL DEFAULT now(),
	"appName" VARCHAR NOT NULL DEFAULT '',
	UNIQUE ("email", "appName")
);