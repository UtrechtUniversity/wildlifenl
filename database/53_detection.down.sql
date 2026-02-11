ALTER TABLE "detection"
	RENAME COLUMN "start" TO "timestamp";
ALTER TABLE "detection"
    DROP COLUMN "end";
ALTER TABLE "detection"
	RENAME COLUMN "deploymentID" TO "sensorID";