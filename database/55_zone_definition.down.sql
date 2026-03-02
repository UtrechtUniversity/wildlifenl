ALTER TABLE "zone"
    ADD COLUMN "area" circle;

UPDATE "zone"
    SET "area" = '<(0,0),1>'::circle
    WHERE "area" IS NULL;

ALTER TABLE "zone"
    ALTER COLUMN "area" SET NOT NULL;

ALTER TABLE "zone"
	DROP COLUMN "definition";