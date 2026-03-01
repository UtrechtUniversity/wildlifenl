
ALTER TABLE "user"
    ADD COLUMN "livingLabID" UUID NULL;

ALTER TABLE "user"
    ADD CONSTRAINT "user_livingLabID_fkey" FOREIGN KEY ("livingLabID") REFERENCES "livingLab" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT;

CREATE TABLE "belonging" (
	"ID" UUID NOT NULL DEFAULT gen_random_uuid(),
	"name" VARCHAR NOT NULL,
	"category" VARCHAR NOT NULL,
	PRIMARY KEY ("ID")
);

ALTER TABLE "damageReport"
    ADD COLUMN "belongingID" UUID NULL;

ALTER TABLE "damageReport"
    ADD CONSTRAINT "belonging_fkey" FOREIGN KEY ("belongingID") REFERENCES "belonging" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT;


