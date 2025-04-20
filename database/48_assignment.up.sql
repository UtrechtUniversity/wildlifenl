CREATE TABLE "assignment" (
	"userID" UUID NOT NULL,
	"questionnaireID" UUID NOT NULL,
	"interactionID" UUID NOT NULL,
	PRIMARY KEY ("userID", "questionnaireID", "interactionID"),
	CONSTRAINT "FK_questionnaireAssignment_user" FOREIGN KEY ("userID") REFERENCES "user" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK_questionnaireAssignment_questionnaire" FOREIGN KEY ("questionnaireID") REFERENCES "questionnaire" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT,
	CONSTRAINT "FK_questionnaireAssignment_interaction" FOREIGN KEY ("interactionID") REFERENCES "interaction" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT
);