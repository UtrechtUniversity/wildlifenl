ALTER TABLE "user"
	RENAME COLUMN "description" TO "notes";
ALTER TABLE "user"
	ADD "natureVisitAvgWeeklyFrequency" INTEGER NOT NULL DEFAULT 0;