DROP TRIGGER check_insert ON "experiment";
DROP TRIGGER check_update ON "experiment";

DROP FUNCTION check_insert_constraints();
DROP FUNCTION check_update_constraints();

ALTER TABLE "experiment"
ADD CONSTRAINT start_must_be_later_than_now CHECK ("start" > NOW());

ALTER TABLE "experiment"
ADD CONSTRAINT end_must_be_later_than_start CHECK ("end" IS NULL OR "end" > "start");
