ALTER TABLE "damageReport"
	ALTER COLUMN "belongingID" TYPE UUID,
	ALTER COLUMN "belongingID" SET NOT NULL,
	ALTER COLUMN "belongingID" DROP DEFAULT;
ALTER TABLE "damageReport"
	DROP COLUMN "belonging";