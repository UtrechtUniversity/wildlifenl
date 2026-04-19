ALTER TABLE "damageReport"
	DROP COLUMN "impactType",
	DROP COLUMN "impactValue",
	DROP COLUMN "estimatedDamage",
	DROP COLUMN "estimatedLoss";

ALTER TABLE "damageReport"
	ADD "perceivedLoss" VARCHAR NOT NULL DEFAULT 'unknown';
ALTER TABLE "damageReport"
	ADD "preventiveMeasures" BOOLEAN NOT NULL DEFAULT 0::boolean;
ALTER TABLE "damageReport"
	ADD "preventiveMeasuresDescription" TEXT NULL;

