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
		SELECT q."ID", q."text", q."description", q."index", q."allowMultipleResponse", q."allowOpenResponse", q."openResponseFormat", COALESCE(a."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(a."text", ''), COALESCE(a."index", 0)
		FROM "question" q
		LEFT JOIN "answer" a ON a."questionID" = q."ID"
		`,
	}
	return &s
}

func (s *QuestionStore) process(rows *sql.Rows, err error) ([]models.Question, error) {
	if err != nil {
		return nil, err
	}
	questions := make([]models.Question, 0)
	var question models.Question
	for rows.Next() {
		var q models.Question
		var a models.Answer
		if err := rows.Scan(&q.ID, &q.Text, &q.Description, &q.Index, &q.AllowMultipleResponse, &q.AllowOpenResponse, &q.OpenResponseFormat, &a.ID, &a.Text, &a.Index); err != nil {
			return nil, err
		}
		if question.ID == "" {
			question = q
		} else if question.ID != q.ID {
			questions = append(questions, question)
			question = q
		}
		if a.ID != "00000000-0000-0000-0000-000000000000" {
			question.Answers = append(question.Answers, a)
		}
	}
	if question.ID != "" {
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
		INSERT INTO "question"("text", "description", "index", "allowMultipleResponse", "allowOpenResponse", "openResponseFormat", "questionnaireID") VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, question.Text, question.Description, question.Index, question.AllowMultipleResponse, question.AllowOpenResponse, question.OpenResponseFormat, question.QuestionnaireID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *QuestionStore) GetByQuestionnaire(questionnaireID string) ([]models.Question, error) {
	query := s.query + `
		WHERE q."questionnaireID" = $1
		ORDER BY q."index" ASC
	`
	rows, err := s.relationalDB.Query(query, questionnaireID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}
