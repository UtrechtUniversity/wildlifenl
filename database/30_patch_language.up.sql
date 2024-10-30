ALTER TABLE "species"
	RENAME COLUMN "commonNameNL" TO "commonName";
ALTER TABLE "species"
	ALTER COLUMN "commonName" TYPE VARCHAR,
	ALTER COLUMN "commonName" SET NOT NULL,
	ALTER COLUMN "commonName" DROP DEFAULT;
ALTER TABLE "species"
	DROP COLUMN "commonNameEN";

ALTER TABLE "interactionType"
	RENAME COLUMN "nameNL" TO "name";
ALTER TABLE "interactionType"
	ALTER COLUMN "name" TYPE VARCHAR,
	ALTER COLUMN "name" SET NOT NULL,
	ALTER COLUMN "name" DROP DEFAULT;
ALTER TABLE "interactionType"
	RENAME COLUMN "descriptionNL" TO "description";
ALTER TABLE "interactionType"
	ALTER COLUMN "description" TYPE TEXT,
	ALTER COLUMN "description" SET NOT NULL,
	ALTER COLUMN "description" SET DEFAULT '';
ALTER TABLE "interactionType"
	DROP COLUMN "nameEN";
ALTER TABLE "interactionType"
	DROP COLUMN "descriptionEN";
