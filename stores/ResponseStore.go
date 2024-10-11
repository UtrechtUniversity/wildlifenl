package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ResponseStore Store

func NewResponseStore(db *sql.DB) *ResponseStore {
	s := ResponseStore{
		relationalDB: db,
		query: `
		SELECT r."ID", r."text", q."ID", q."text", q."description", q."index", q."allowMultipleResponse", q."allowOpenResponse", i."ID", i."timestamp", i."description", i."location", u."ID", u."name", t."ID", t."nameNL", t."nameEN", t."descriptionNL", t."descriptionEN", COALESCE(a."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(a."text", ''), COALESCE(a."index", 0)
		FROM "response" r
		INNER JOIN "question" q ON q."ID" = r."questionID"
		INNER JOIN "interaction" i ON i."ID" = r."interactionID"
		INNER JOIN "user" u ON u."ID" = i."userID"
		INNER JOIN "interactionType" t ON t."ID" = i."typeID"
		LEFT JOIN "answer" a ON a."ID" = r."answerID"
		`,
	}
	return &s
}

func (s *ResponseStore) process(rows *sql.Rows, err error) ([]models.Response, error) {
	if err != nil {
		return nil, err
	}
	responses := make([]models.Response, 0)
	for rows.Next() {
		var r models.Response
		var a models.Answer
		if err := rows.Scan(&r.ID, &r.Text, &r.Question.ID, &r.Question.Text, &r.Question.Description, &r.Question.Index, &r.Question.AllowMultipleResponse, &r.Question.AllowOpenResponse, &r.Interaction.ID, &r.Interaction.Timestamp, &r.Interaction.Description, &r.Interaction.Location, &r.Interaction.User.ID, &r.Interaction.User.Name, &r.Interaction.Type.ID, &r.Interaction.Type.NameNL, &r.Interaction.Type.NameEN, &r.Interaction.Type.DescriptionNL, &r.Interaction.Type.DescriptionEN, &a.ID, &a.Text, &a.Index); err != nil {
			return nil, err
		}
		if a.ID != "00000000-0000-0000-0000-000000000000" {
			r.Answer = &a
		}
		responses = append(responses, r)
	}
	return responses, nil
}

func (s *ResponseStore) Get(responseID string) (*models.Response, error) {
	query := s.query + `
		WHERE r."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, responseID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ResponseStore) Add(userID string, response *models.ResponseRecord) (*models.Response, error) {
	var id string
	var args []any
	query := `
		WITH sanity_check AS (
			SELECT i."ID" 
			FROM "interaction" i
			INNER JOIN "user" u ON u."ID" = i."userID"
			INNER JOIN "interactionType" t ON t."ID" = i."typeID"
			INNER JOIN "questionnaire" n ON n."interactionTypeID" = t."ID"
			INNER JOIN "question" q ON q."questionnaireID" = n."ID"
			LEFT JOIN "answer" a ON a."questionID" = q."ID"
			WHERE u."ID" = $1
			AND q."ID" = $3
			AND i."ID" = $4
	`
	if response.AnswerID == nil {
		query += `
			)
			INSERT INTO "response"("text", "questionID", "interactionID")
			SELECT $2, $3, $4
		`
		args = []any{userID, response.Text, response.QuestionID, response.InteractionID}
	} else {
		query += `
				AND a."ID" = $5
			)
			INSERT INTO "response"("text", "questionID", "interactionID", "answerID")
			SELECT $2, $3, $4, $5
		`
		args = []any{userID, response.Text, response.QuestionID, response.InteractionID, response.AnswerID}
	}
	query += `
		WHERE (SELECT COUNT(*) FROM sanity_check) = 1
		RETURNING "ID";
	`
	row := s.relationalDB.QueryRow(query, args...)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}

func (s *ResponseStore) GetForInteractionByQuestion(interactionID string, questionID string) ([]models.Response, error) {
	query := s.query + `
		WHERE i."ID" = $1
		AND q."ID" = $2
		`
	rows, err := s.relationalDB.Query(query, interactionID, questionID)
	return s.process(rows, err)
}
