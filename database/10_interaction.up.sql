CREATE TABLE IF NOT EXISTS "interaction" (
	"ID" UUID NOT NULL DEFAULT gen_random_uuid(),
	"description" TEXT NOT NULL DEFAULT '',
	"speciesID" UUID NOT NULL,
	"timestamp" TIMESTAMP NOT NULL DEFAULT now(),
	"userID" UUID NOT NULL,
	"location" POINT NOT NULL,
	"typeID" INT NOT NULL,
	PRIMARY KEY ("ID"),
	CONSTRAINT "FK_interaction_species" FOREIGN KEY ("speciesID") REFERENCES "species" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK_interaction_user" FOREIGN KEY ("userID") REFERENCES "user" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK_interaction_interactionType" FOREIGN KEY ("typeID") REFERENCES "interactionType" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);