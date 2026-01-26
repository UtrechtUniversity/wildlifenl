
ALTER TABLE "detection"
ADD COLUMN "id_tmp" UUID NOT NULL DEFAULT gen_random_uuid();

ALTER TABLE "alarm"
ADD COLUMN "detectionid_tmp" UUID;

UPDATE "alarm" a
SET "detectionid_tmp" = d."id_tmp"
FROM "detection" d
WHERE a."detectionID" = d."ID";

ALTER TABLE "alarm"
DROP CONSTRAINT "FK__detection";

ALTER TABLE "detection"
DROP CONSTRAINT detection_pkey;

ALTER TABLE "alarm"
DROP COLUMN "detectionID";

ALTER TABLE "detection"
DROP COLUMN "ID";

ALTER TABLE "detection"
RENAME COLUMN "id_tmp" TO "ID";

ALTER TABLE "alarm"
RENAME COLUMN "detectionid_tmp" TO "detectionID";

ALTER TABLE "detection"
ADD CONSTRAINT detection_pkey PRIMARY KEY ("ID");

ALTER TABLE "alarm"
ADD CONSTRAINT FK_detection
FOREIGN KEY ("detectionID")
REFERENCES "detection"("ID")
ON UPDATE RESTRICT
ON DELETE RESTRICT;