ALTER TABLE "message"
	ADD COLUMN "trigger" VARCHAR CHECK ("trigger" IN ('encounter', 'answer', 'alarm')) NOT NULL;
ALTER TABLE "message"
    ADD COLUMN "encounterMeters" INT NOT NULL DEFAULT 0;
ALTER TABLE "message"
    ADD COLUMN "encounterMinutes" INT NOT NULL DEFAULT 0;