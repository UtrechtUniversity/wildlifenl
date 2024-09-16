CREATE TABLE IF NOT EXISTS "question" (
	"ID" UUID NOT NULL DEFAULT gen_random_uuid(),
    "questionnaireID" UUID NOT NULL,
	"text" VARCHAR NOT NULL,
	"description" TEXT NOT NULL,
	"index" INT NOT NULL,
	"allowMultipleResponse" BOOLEAN NOT NULL,
	"allowOpenResponse" BOOLEAN NOT NULL,
    PRIMARY KEY ("ID"),
	CONSTRAINT "FK__questionnaire" FOREIGN KEY ("questionnaireID") REFERENCES "questionnaire" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);