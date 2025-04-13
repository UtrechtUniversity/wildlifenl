CREATE TABLE "damageReport" (
	"interactionID" UUID NOT NULL,
	"belongingID" UUID NOT NULL,
	"impactType" VARCHAR NOT NULL,
	"impactValue" INTEGER NOT NULL,
	"estimatedDamage" INTEGER NOT NULL,
	"estimatedLoss" INTEGER NOT NULL,
	PRIMARY KEY ("interactionID"),
	CONSTRAINT "FK__interaction" FOREIGN KEY ("interactionID") REFERENCES "interaction" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK__belonging" FOREIGN KEY ("belongingID") REFERENCES "belonging" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);