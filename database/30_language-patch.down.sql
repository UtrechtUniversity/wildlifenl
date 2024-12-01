ALTER TABLE "species"
	RENAME COLUMN "commonName" TO "commonNameNL";
ALTER TABLE "species"
	ALTER COLUMN "commonNameNL" TYPE VARCHAR,
	ALTER COLUMN "commonNameNL" SET NOT NULL,
	ALTER COLUMN "commonNameNL" DROP DEFAULT;
ALTER TABLE "species"
	ADD COLUMN "commonNameEN" VARCHAR NULL;

ALTER TABLE "interactionType"
	RENAME COLUMN "name" TO "nameNL";
ALTER TABLE "interactionType"
	ALTER COLUMN "nameNL" TYPE VARCHAR,
	ALTER COLUMN "nameNL" SET NOT NULL,
	ALTER COLUMN "nameNL" DROP DEFAULT;
ALTER TABLE "interactionType"
	RENAME COLUMN "description" TO "descriptionNL";
ALTER TABLE "interactionType"
	ALTER COLUMN "descriptionNL" TYPE TEXT,
	ALTER COLUMN "descriptionNL" SET NOT NULL,
	ALTER COLUMN "descriptionNL" SET DEFAULT '';
ALTER TABLE "interactionType"
	ADD COLUMN "nameEN" VARCHAR NULL;
ALTER TABLE "interactionType"
	ADD COLUMN "descriptionEN" VARCHAR NULL;