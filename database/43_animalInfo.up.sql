CREATE TABLE "animalInfo" (
	"ID" SERIAL NOT NULL,
	"interactionID" UUID NOT NULL,
	"sex" VARCHAR NOT NULL,
	"lifeStage" VARCHAR NOT NULL,
	"condition" VARCHAR NOT NULL,
	PRIMARY KEY ("ID"),
    CONSTRAINT "FK__interaction" FOREIGN KEY ("interactionID") REFERENCES "interaction" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);