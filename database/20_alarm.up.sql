CREATE TABLE "alarm" (
    "ID" UUID NOT NULL DEFAULT gen_random_uuid(),
	"timestamp" TIMESTAMPTZ NOT NULL DEFAULT now(),
	"zoneID" UUID NOT NULL,
	"interactionID" UUID NULL,
    "detectionID" INT NULL,
    "animalID" UUID NULL,
	PRIMARY KEY ("ID"),
	CONSTRAINT "FK__zone" FOREIGN KEY ("zoneID") REFERENCES "zone" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT "FK__interaction" FOREIGN KEY ("interactionID") REFERENCES "interaction" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT "FK__detection" FOREIGN KEY ("detectionID") REFERENCES "detection" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
    CONSTRAINT "FK__animal" FOREIGN KEY ("animalID") REFERENCES "animal" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);