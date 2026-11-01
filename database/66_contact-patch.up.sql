ALTER TABLE "contact"
	ADD "startLocation" POINT NOT NULL DEFAULT '(0,90)',
	ADD "endLocation" POINT NOT NULL DEFAULT '(0,90)';
ALTER TABLE "contact"
	ALTER COLUMN "startLocation" DROP DEFAULT;