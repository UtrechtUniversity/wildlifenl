CREATE TABLE IF NOT EXISTS "answer" (
	"ID" UUID NOT NULL DEFAULT gen_random_uuid(),
	"questionID" UUID NOT NULL,
	"text" VARCHAR NOT NULL,
    "index" INT NOT NULL,
    PRIMARY KEY ("ID"),
	CONSTRAINT "FK__question" FOREIGN KEY ("questionID") REFERENCES "question" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);