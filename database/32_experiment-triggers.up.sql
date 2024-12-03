ALTER TABLE "experiment"
DROP CONSTRAINT start_must_be_later_than_now;

ALTER TABLE "experiment"
DROP CONSTRAINT end_must_be_later_than_start;


CREATE OR REPLACE FUNCTION check_insert_constraints()
RETURNS TRIGGER AS '
BEGIN
    IF NEW."start" <= NOW() THEN
        RAISE EXCEPTION ''start must be later than now'';
    END IF;

    IF NEW."end" IS NOT NULL AND NEW."end" <= NEW."start" THEN
        RAISE EXCEPTION ''end must be later than start'';
    END IF;

    RETURN NEW;
END;
' LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION check_update_constraints()
RETURNS TRIGGER AS '
BEGIN
    IF OLD."start" IS NOT NULL AND OLD."end" IS NOT NULL THEN
        RAISE EXCEPTION ''cannot update a finished experiment'';
    END IF;

    IF NEW."start" IS DISTINCT FROM OLD."start" AND NEW."end" IS DISTINCT FROM OLD."end" THEN
        IF OLD."start" <= NOW() THEN
            RAISE EXCEPTION ''cannot update a started experiment 1'';
        END IF;

        IF NEW."start" <= NOW() THEN
            RAISE EXCEPTION ''start must be later than now'';
        END IF;

        IF NEW."end" <= NEW."start" THEN
            RAISE EXCEPTION ''end must be later than start'';
        END IF;

    ELSIF NEW."start" IS DISTINCT FROM OLD."start" THEN
        IF OLD."start" <= NOW() THEN
            RAISE EXCEPTION ''cannot update a started experiment 2'';
        END IF;

        IF NEW."start" <= NOW() THEN
            RAISE EXCEPTION ''start must be later than now'';
        END IF;

        IF OLD."end" IS NOT NULL AND NEW."start" >= OLD."end" THEN
            RAISE EXCEPTION ''start must be earlier than end'';
        END IF;

    ELSIF NEW."end" IS DISTINCT FROM OLD."end" THEN
        IF NEW."end" <= OLD."start" THEN
            RAISE EXCEPTION ''end must be later than start'';
        END IF;
    END IF;

    RETURN NEW;
END;
' LANGUAGE plpgsql;


CREATE TRIGGER check_insert
BEFORE INSERT ON "experiment"
FOR EACH ROW
EXECUTE FUNCTION check_insert_constraints();

CREATE TRIGGER check_update
BEFORE UPDATE ON "experiment"
FOR EACH ROW
EXECUTE FUNCTION check_update_constraints();