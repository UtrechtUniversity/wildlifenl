ALTER TABLE "detection"
	RENAME COLUMN "timestamp" TO "start";
ALTER TABLE "detection"
	ADD "end" TIMESTAMPTZ NOT NULL DEFAULT now();
ALTER TABLE "detection"
	RENAME COLUMN "sensorID" TO "deploymentID";