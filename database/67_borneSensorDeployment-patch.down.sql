ALTER TABLE "borneSensorDeployment"
	DROP CONSTRAINT IF EXISTS "FK_borneSensorDeployment_user";

ALTER TABLE "borneSensorDeployment"
	DROP COLUMN "userID";