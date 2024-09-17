package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type AnswerStore Store

func NewAnswerStore(db *sql.DB) *AnswerStore {
	s := AnswerStore{
		relationalDB: db,
		query: `
		SELECT a."ID", a."text", a."index"
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
		var answer models.Answer
		if err := rows.Scan(&answer.ID, &answer.Text, &answer.Index); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
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
		INSERT INTO "answer"("text", "index", "questionID") VALUES($1, $2, $3)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, answer.Text, answer.Index, answer.QuestionID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
