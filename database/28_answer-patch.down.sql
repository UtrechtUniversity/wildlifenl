ALTER TABLE "answer"
    DROP CONSTRAINT FK_next_question;

ALTER TABLE "answer"
    DROP COLUMN "nextQuestionID";