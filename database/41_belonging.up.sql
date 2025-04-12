CREATE TABLE "belonging" (
	"ID" UUID NOT NULL DEFAULT gen_random_uuid(),
	"name" VARCHAR NOT NULL,
	"category" VARCHAR NOT NULL,
	PRIMARY KEY ("ID")
);