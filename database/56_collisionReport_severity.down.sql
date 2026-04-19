ALTER TABLE "collisionReport"
	RENAME COLUMN "severity" TO "intensity";
ALTER TABLE "collisionReport"
	ADD "urgency" VARCHAR NOT NULL;