ALTER TABLE "experiment"
DROP CONSTRAINT start_must_be_later_than_now;

ALTER TABLE "experiment"
DROP CONSTRAINT end_must_be_later_than_start;