ALTER TABLE "detection"
	ADD "sensorType" VARCHAR;

ALTER TABLE "detection"
	ADD "uri" VARCHAR;

CREATE TABLE "detectionAnimal" (
	"ID" SERIAL NOT NULL,
	"detectionID" UUID NOT NULL,
	"confidence" INT NOT NULL,
	"sex" VARCHAR,
	"lifeStage" VARCHAR,
	"condition" VARCHAR,
	"behaviour" TEXT,
	"description" TEXT,
	PRIMARY KEY ("ID"),
    CONSTRAINT "FK__detection" FOREIGN KEY ("detectionID") REFERENCES "detection" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);