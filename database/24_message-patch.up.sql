ALTER TABLE "message"
	ADD COLUMN "trigger" VARCHAR CHECK ("trigger" IN ('encounter', 'answer', 'alarm')) NOT NULL;
ALTER TABLE "message"
    ADD COLUMN "encounterMeters" INT NULL DEFAULT NULL;
ALTER TABLE "message"
    ADD COLUMN "encounterMinutes" INT NULL DEFAULT NULL;