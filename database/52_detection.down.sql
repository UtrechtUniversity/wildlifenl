ALTER TABLE "detection"
	DROP COLUMN "sensorType";

ALTER TABLE "detection"
	ADD COLUMN "uri";

DROP TABLE "detectionAnimal";