ALTER TABLE "contact"
	ADD "startLocation" POINT NOT NULL DEFAULT '(0,90)',
	ADD "endLocation" POINT NULL;