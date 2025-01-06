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
		SELECT r."ID", r."text", q."ID", q."text", q."description", q."index", q."allowMultipleResponse", q."allowOpenResponse", n."ID", n."name", n."identifier", e."ID", e."name", e."description", e."start", e."end", uu."ID", uu."name", COALESCE(ll."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(ll."name", ''), i."ID", i."timestamp", i."description", i."location", u."ID", u."name", t."ID", t."name", t."description", s."name", s."commonName", s."category", s."advice", s."didYouKnow", COALESCE(a."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(a."text", ''), COALESCE(a."index", 0)
		FROM "response" r
		INNER JOIN "question" q ON q."ID" = r."questionID"
		INNER JOIN "questionnaire" n ON q."questionnaireID" = n."ID"
		INNER JOIN "experiment" e ON e."ID" = n."experimentID"
		INNER JOIN "user" uu ON uu."ID" = e."userID"
		INNER JOIN "interaction" i ON i."ID" = r."interactionID"
		INNER JOIN "user" u ON u."ID" = i."userID"
		INNER JOIN "interactionType" t ON t."ID" = i."typeID"
		INNER JOIN "species" s ON s."ID" = i."speciesID"
		LEFT JOIN "answer" a ON a."ID" = r."answerID"
		LEFT JOIN "livingLab" ll ON ll."ID" = e."livingLabID"
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
		r.Question.Questionnaire = new(models.Questionnaire)
		var a models.Answer
		var ll models.LivingLab
		if err := rows.Scan(&r.ID, &r.Text, &r.Question.ID, &r.Question.Text, &r.Question.Description, &r.Question.Index, &r.Question.AllowMultipleResponse, &r.Question.AllowOpenResponse, &r.Question.Questionnaire.ID, &r.Question.Questionnaire.Name, &r.Question.Questionnaire.Identifier, &r.Question.Questionnaire.Experiment.ID, &r.Question.Questionnaire.Experiment.Name, &r.Question.Questionnaire.Experiment.Description, &r.Question.Questionnaire.Experiment.Start, &r.Question.Questionnaire.Experiment.End, &r.Question.Questionnaire.Experiment.User.ID, &r.Question.Questionnaire.Experiment.User.Name, &ll.ID, &ll.Name, &r.Interaction.ID, &r.Interaction.Timestamp, &r.Interaction.Description, &r.Interaction.Location, &r.Interaction.User.ID, &r.Interaction.User.Name, &r.Interaction.Type.ID, &r.Interaction.Type.Name, &r.Interaction.Type.Description, &r.Interaction.Species.Name, &r.Interaction.Species.CommonName, &r.Interaction.Species.Category, &r.Interaction.Species.Advice, &r.Interaction.Species.DidYouKnow, &a.ID, &a.Text, &a.Index); err != nil {
			return nil, err
		}
		if a.ID != "00000000-0000-0000-0000-000000000000" {
			r.Answer = &a
		}
		if ll.ID != "00000000-0000-0000-0000-000000000000" {
			r.Question.Questionnaire.Experiment.LivingLab = &ll
		}
		r.Question.Questionnaire.InteractionType = r.Interaction.Type
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

func (s *ResponseStore) GetByQuestionnaire(questionnaireID string) ([]models.Response, error) {
	query := s.query + `
		WHERE n."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, questionnaireID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return result, nil
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

func (s *ResponseStore) GetByExperiment(experimentID string) ([]models.Response, error) {
	query := s.query + `
		WHERE e."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, experimentID)
	return s.process(rows, err)
}
