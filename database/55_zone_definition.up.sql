ALTER TABLE "zone"
    ADD COLUMN "definition" polygon;

UPDATE "zone"
    SET "definition" = '((0,0),(1,0),(1,1),(0,1))'::polygon
    WHERE "definition" IS NULL;

ALTER TABLE "zone"
    ALTER COLUMN "definition" SET NOT NULL;

ALTER TABLE "zone"
	DROP COLUMN "area";