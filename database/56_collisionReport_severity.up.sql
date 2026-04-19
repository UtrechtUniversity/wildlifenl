ALTER TABLE "collisionReport"
	RENAME COLUMN "intensity" TO "severity";
ALTER TABLE "collisionReport"
	DROP COLUMN "urgency";