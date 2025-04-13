CREATE TABLE "sightingReport" (
	"interactionID" UUID NOT NULL,
	PRIMARY KEY ("interactionID"),
	CONSTRAINT "FK__interaction" FOREIGN KEY ("interactionID") REFERENCES "interaction" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);