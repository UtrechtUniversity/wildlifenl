ALTER TABLE "damageReport"
	DROP COLUMN "perceivedLoss",
	DROP COLUMN "preventiveMeasures",
	DROP COLUMN "preventiveMeasuresDescription";

ALTER TABLE "damageReport"
	ADD "impactType" VARCHAR NOT NULL;
ALTER TABLE "damageReport"
	ADD "impactValue" INTEGER NOT NULL;
ALTER TABLE "damageReport"
	ADD "estimatedDamage" INTEGER NOT NULL;
ALTER TABLE "damageReport"
	ADD "estimatedLoss" INTEGER NOT NULL;