CREATE TABLE IF NOT EXISTS "zone" (
	"ID" UUID NOT NULL DEFAULT gen_random_uuid(),
	"name" VARCHAR NOT NULL,
	"description" TEXT NOT NULL DEFAULT '',
	"area" CIRCLE NOT NULL,
    "userID" UUID NOT NULL,
    PRIMARY KEY ("ID"),
    CONSTRAINT "FK__user" FOREIGN KEY ("userID") REFERENCES "user" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);