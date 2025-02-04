ALTER TABLE "species"
	DROP COLUMN "behaviour";
ALTER TABLE "species"
	RENAME COLUMN "roleInNature" TO "didYouKnow";
ALTER TABLE "species"
	DROP COLUMN "description";