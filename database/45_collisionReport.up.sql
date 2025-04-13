CREATE TABLE "collisionReport" (
	"interactionID" UUID NOT NULL,
	"estimatedDamage" INTEGER NOT NULL,
	"intensity" VARCHAR NOT NULL,
	"urgency" VARCHAR NOT NULL,
	PRIMARY KEY ("interactionID"),
	CONSTRAINT "FK__interaction" FOREIGN KEY ("interactionID") REFERENCES "interaction" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);