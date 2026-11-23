ALTER TABLE "borneSensorDeployment"
  ADD COLUMN "userID" UUID NULL;

UPDATE "borneSensorDeployment" d
SET "userID" = (
	SELECT u."ID"
	FROM "user" u
	INNER JOIN "user_role" ur ON ur."userID" = u."ID"
	INNER JOIN "role" r ON r."ID" = ur."roleID"
	WHERE r."name" = 'herd-manager'
	LIMIT 1
)
WHERE d."userID" IS NULL;

ALTER TABLE "borneSensorDeployment"
  ALTER COLUMN "userID" SET NOT NULL;

ALTER TABLE "borneSensorDeployment"
  ADD CONSTRAINT "FK_borneSensorDeployment_user"
  FOREIGN KEY ("userID") REFERENCES "user" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT;