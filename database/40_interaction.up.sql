ALTER TABLE "interaction"
	ALTER COLUMN "timestamp" TYPE TIMESTAMPTZ,
	ALTER COLUMN "timestamp" SET NOT NULL,
	ALTER COLUMN "timestamp" SET DEFAULT now();

ALTER TABLE "interaction"
	ADD "moment" TIMESTAMPTZ NOT NULL DEFAULT now();
