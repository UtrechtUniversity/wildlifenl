ALTER TABLE "experiment"
ADD CONSTRAINT start_must_be_later_than_now CHECK ("start" > NOW());

ALTER TABLE "experiment"
ADD CONSTRAINT end_must_be_later_than_start CHECK ("end" IS NULL OR "end" > "start");