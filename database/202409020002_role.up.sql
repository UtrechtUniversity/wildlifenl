CREATE TABLE IF NOT EXISTS "role" (
	"id" SERIAL NOT NULL,
	"name" VARCHAR NOT NULL,
	PRIMARY KEY ("id")
);

INSERT INTO "role"("name") VALUES ('administrator');
INSERT INTO "role"("name") VALUES ('data-system');
INSERT INTO "role"("name") VALUES ('researcher');
INSERT INTO "role"("name") VALUES ('land-user');
INSERT INTO "role"("name") VALUES ('nature-area-manager');
INSERT INTO "role"("name") VALUES ('wildlife-manager');
INSERT INTO "role"("name") VALUES ('herd-manager');