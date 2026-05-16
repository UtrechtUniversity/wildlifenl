ALTER TABLE "borneSensorDeployment"
	ADD "contactHardwareAddress" VARCHAR NULL;
ALTER TABLE "borneSensorDeployment"
	ADD "ID" UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE "borneSensorDeployment"
	ADD PRIMARY KEY ("ID");