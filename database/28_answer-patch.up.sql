ALTER TABLE "answer"
	ADD "nextQuestionID" UUID NULL;
ALTER TABLE "answer"
	ADD CONSTRAINT "FK_next_question" FOREIGN KEY ("nextQuestionID") REFERENCES "question" ("ID") ON UPDATE RESTRICT ON DELETE RESTRICT;