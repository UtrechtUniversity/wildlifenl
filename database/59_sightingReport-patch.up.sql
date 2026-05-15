ALTER TABLE "sightingReport"
	ADD "humanActivity" VARCHAR NOT NULL DEFAULT 'unknown';
ALTER TABLE "sightingReport"
	ADD "humanActivityOther" VARCHAR NULL;
ALTER TABLE "sightingReport"
	ADD "perceivedAnimalActivity" VARCHAR NOT NULL DEFAULT 'unknown';
ALTER TABLE "sightingReport"
	ADD "perceivedAnimalActivityOther" VARCHAR NULL;