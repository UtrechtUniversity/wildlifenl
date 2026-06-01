ALTER TABLE "interaction"
	RENAME COLUMN "description" TO "notes";
ALTER TABLE "interaction"
	ALTER COLUMN "notes" TYPE TEXT,
	ALTER COLUMN "notes" DROP NOT NULL,
	ALTER COLUMN "notes" DROP DEFAULT;