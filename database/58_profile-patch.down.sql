ALTER TABLE "user"
	RENAME COLUMN "notes" TO "description";
ALTER TABLE "user"
	DROP COLUMN "natureVisitAvgWeeklyFrequency";