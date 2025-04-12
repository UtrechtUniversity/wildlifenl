ALTER TABLE "interaction"
	ALTER COLUMN "timestamp" TYPE TIMESTAMP,
	ALTER COLUMN "timestamp" SET NOT NULL,
	ALTER COLUMN "timestamp" SET DEFAULT now();

ALTER TABLE "interaction"
	DROP COLUMN "moment";