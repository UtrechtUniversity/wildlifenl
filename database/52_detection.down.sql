ALTER TABLE "detection"
	DROP COLUMN "sensorType";

ALTER TABLE "detection"
	ADD COLUMN "uri";

ALTER TABLE "detection"
	DROP CONSTRAINT "FK__user";

ALTER TABLE "detection"
	ADD COLUMN "userID";

DROP TABLE "detectionAnimal";