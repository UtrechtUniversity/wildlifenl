CREATE TABLE "response" (
    "ID" UUID NOT NULL DEFAULT gen_random_uuid(),
    "text" TEXT NOT NULL,
    "questionID" UUID NOT NULL,
    "interactionID" UUID NOT NULL,
    "answerID" UUID NULL,
	PRIMARY KEY ("ID"),
    CONSTRAINT "FK__question" FOREIGN KEY ("questionID") REFERENCES "question" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK__interaction" FOREIGN KEY ("interactionID") REFERENCES "interaction" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT "FK__answer" FOREIGN KEY ("answerID") REFERENCES "answer" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);