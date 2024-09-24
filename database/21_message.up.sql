CREATE TABLE "message" (
    "ID" UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" VARCHAR NOT NULL,
    "severity" INT NOT NULL,
    "text" TEXT NOT NULL,
    "experimentID" UUID NOT NULL,
    "speciesID" UUID NULL,
    "answerID" UUID NULL,
	PRIMARY KEY ("ID"),
    CONSTRAINT "FK__experiment" FOREIGN KEY ("experimentID") REFERENCES "experiment" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK__species" FOREIGN KEY ("speciesID") REFERENCES "species" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT "FK__answer" FOREIGN KEY ("answerID") REFERENCES "answer" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);