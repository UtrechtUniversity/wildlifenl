ALTER TABLE "detection"
ADD COLUMN "id_tmp_int" SERIAL; 

ALTER TABLE "alarm"
ADD COLUMN "detectionid_tmp_int" INTEGER;

UPDATE "alarm" a
SET "detectionid_tmp_int" = d."id_tmp_int"
FROM "detection" d
WHERE a."detectionID" = d."ID";   

ALTER TABLE "alarm"
DROP CONSTRAINT IF EXISTS "FK_detection";

ALTER TABLE "detection"
DROP CONSTRAINT IF EXISTS "detection_pkey";

ALTER TABLE "alarm"
DROP COLUMN "detectionID";

ALTER TABLE "detection"
DROP COLUMN "ID";

ALTER TABLE "detection"
RENAME COLUMN "id_tmp_int" TO "ID";

ALTER TABLE "alarm"
RENAME COLUMN "detectionid_tmp_int" TO "detectionID";

ALTER TABLE "detection"
ADD CONSTRAINT detection_pkey PRIMARY KEY ("ID");

ALTER TABLE "alarm"
ADD CONSTRAINT "FK_detection"
FOREIGN KEY ("detectionID")
REFERENCES "detection"("ID")
ON UPDATE RESTRICT
ON DELETE RESTRICT;