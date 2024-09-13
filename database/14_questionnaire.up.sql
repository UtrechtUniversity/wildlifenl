CREATE TABLE IF NOT EXISTS "questionnaire" (
	"ID" UUID NOT NULL DEFAULT gen_random_uuid(),
	"name" VARCHAR NOT NULL,
	"experimentID" UUID NOT NULL,
	"interactionTypeID" INTEGER NOT NULL,
	PRIMARY KEY ("ID"),
	CONSTRAINT "FK__experiment" FOREIGN KEY ("experimentID") REFERENCES "experiment" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK__interactionType" FOREIGN KEY ("interactionTypeID") REFERENCES "interactionType" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);