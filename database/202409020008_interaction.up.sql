CREATE TABLE IF NOT EXISTS "interaction" (
	"id" UUID NOT NULL DEFAULT gen_random_uuid(),
	"description" TEXT NULL DEFAULT NULL,
	"speciesID" UUID NOT NULL,
	"timestamp" TIMESTAMP NOT NULL DEFAULT now(),
	"userID" UUID NOT NULL,
	"location" POINT NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "FK_interaction_species" FOREIGN KEY ("speciesID") REFERENCES "species" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK_interaction_user" FOREIGN KEY ("userID") REFERENCES "user" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT
);