package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type QuestionStore Store

func NewQuestionStore(db *sql.DB) *QuestionStore {
	s := QuestionStore{
		relationalDB: db,
		query: `
		SELECT q."ID", q."text", q."description", q."index", q."allowMultipleResponse", q."allowOpenResponse"
		FROM "question" q
		`,
	}
	return &s
}

func (s *QuestionStore) process(rows *sql.Rows, err error) ([]models.Question, error) {
	if err != nil {
		return nil, err
	}
	questions := make([]models.Question, 0)
	for rows.Next() {
		var question models.Question
		if err := rows.Scan(&question.ID, &question.Text, &question.Description, &question.Index, &question.AllowMultipleResponse, &question.AllowOpenResponse); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, nil
}

func (s *QuestionStore) Get(questionID string) (*models.Question, error) {
	query := s.query + `
		WHERE q."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, questionID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *QuestionStore) Add(question *models.QuestionRecord) (*models.Question, error) {
	query := `
		INSERT INTO "question"("text", "description", "index", "allowMultipleResponse", "allowOpenResponse", "questionnaireID") VALUES($1, $2, $3, $4, $5)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, question.Text, question.Description, question.Index, question.AllowMultipleResponse, question.AllowOpenResponse, question.QuestionnaireID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
