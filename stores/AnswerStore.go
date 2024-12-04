package stores

import (
	"database/sql"
	"errors"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type AnswerStore Store

func NewAnswerStore(db *sql.DB) *AnswerStore {
	s := AnswerStore{
		relationalDB: db,
		query: `
		SELECT a."ID", a."text", a."index", a."nextQuestionID"
		FROM "answer" a
		`,
	}
	return &s
}

func (s *AnswerStore) process(rows *sql.Rows, err error) ([]models.Answer, error) {
	if err != nil {
		return nil, err
	}
	answers := make([]models.Answer, 0)
	for rows.Next() {
		var a models.Answer
		if err := rows.Scan(&a.ID, &a.Text, &a.Index, &a.NextQuestionID); err != nil {
			return nil, err
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (s *AnswerStore) Get(answerID string) (*models.Answer, error) {
	query := s.query + `
		WHERE a."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, answerID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *AnswerStore) Add(answer *models.AnswerRecord) (*models.Answer, error) {
	query := `
		INSERT INTO "answer"("text", "index", "questionID", "nextQuestionID") VALUES($1, $2, $3, $4)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, answer.Text, answer.Index, answer.QuestionID, answer.NextQuestionID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *AnswerStore) Delete(answerID string, userID string) error {
	query := `
		WITH deleted AS (
			DELETE FROM "answer" a
			USING "question" q
			JOIN "questionnaire" qq ON q."questionnaireID" = qq."ID"
			JOIN "experiment" e ON qq."experimentID" = e."ID"
			WHERE a."ID" = $1
			AND e."userID" = $2
			AND e.start > NOW()
			RETURNING a."ID"
		)
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM "answer" WHERE "ID" = $1) THEN 'INVALID'
				WHEN NOT EXISTS (SELECT 1 FROM "answer" a JOIN "question" q ON a."questionID" = q."ID" JOIN "questionnaire" qq ON q."questionnaireID" = qq."ID" JOIN "experiment" e ON qq."experimentID" = e."ID" WHERE a."ID" = $1 AND e."userID" = $2) THEN 'USER'
				WHEN EXISTS (SELECT 1 FROM "answer" WHERE "ID" = $1) AND NOT EXISTS (SELECT 1 FROM deleted) THEN 'STARTED'
				WHEN EXISTS (SELECT 1 FROM deleted) THEN 'OK'
			END AS result;
	`
	var state string
	row := s.relationalDB.QueryRow(query, answerID, userID)
	if err := row.Scan(&state); err != nil {
		return err
	}
	switch state {
	case "INVALID":
		return &ErrRecordInattainable{message: "answer was not found"}
	case "USER":
		return &ErrRecordInattainable{message: "answer does not exist for the current user"}
	case "STARTED":
		return &ErrRecordImmutable{message: "cannot delete answer for an experiment that has started"}
	case "OK":
		return nil
	}
	return errors.New("unknown error")
}
