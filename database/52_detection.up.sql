ALTER TABLE "detection"
	ADD "sensorType" VARCHAR;

ALTER TABLE "detection"
	ADD "uri" VARCHAR;

ALTER TABLE "detection"
	ADD "userID" UUID NULL;

UPDATE "detection" SET "userID" = (
	SELECT u."ID" 
	FROM "user" u
	INNER JOIN "user_role" ur ON ur."userID" = u."ID"
	INNER JOIN "role" r ON r."ID" = ur."roleID"
	WHERE r."name" = 'administrator'
	LIMIT 1
);

ALTER TABLE "detection"
	ALTER COLUMN "userID" SET NOT NULL;

ALTER TABLE "detection"
	ADD CONSTRAINT "FK__user" FOREIGN KEY ("userID") REFERENCES "user" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT;

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