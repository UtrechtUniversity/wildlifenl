package stores

import (
	"database/sql"
	"errors"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type QuestionnaireStore Store

func NewQuestionnaireStore(db *sql.DB) *QuestionnaireStore {
	s := QuestionnaireStore{
		relationalDB: db,
		query: `
		SELECT q."ID" AS "questionnaireID", q."name", q."identifier", e."ID", e."name", e."start", e."end", t."ID", t."name", t."description", u."ID", u."name"
		FROM "questionnaire" q
		INNER JOIN "experiment" e ON e."ID" = q."experimentID"
		INNER JOIN "interactionType" t ON t."ID" = q."interactionTypeID"
		INNER JOIN "user" u ON u."ID" = e."userID"
		`,
	}
	return &s
}

func (s *QuestionnaireStore) process(rows *sql.Rows, err error) ([]models.Questionnaire, error) {
	if err != nil {
		return nil, err
	}
	questionnaires := make([]models.Questionnaire, 0)
	for rows.Next() {
		var questionnaire models.Questionnaire
		var experiment models.Experiment
		var interactionType models.InteractionType
		var user models.User
		if err := rows.Scan(&questionnaire.ID, &questionnaire.Name, &questionnaire.Identifier, &experiment.ID, &experiment.Name, &experiment.Start, &experiment.End, &interactionType.ID, &interactionType.Name, &interactionType.Description, &user.ID, &user.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		experiment.User = user
		questionnaire.Experiment = experiment
		questionnaire.InteractionType = interactionType
		questionnaires = append(questionnaires, questionnaire)
	}
	return questionnaires, nil
}

func (s *QuestionnaireStore) Get(questionnaireID string) (*models.Questionnaire, error) {
	query := s.query + `
		WHERE q."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, questionnaireID)
	questionnaires, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(questionnaires) != 1 {
		return nil, nil
	}
	return s.addQuestions(&questionnaires[0])
}

func (s *QuestionnaireStore) GetAll() ([]models.Questionnaire, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *QuestionnaireStore) Add(questionnaire *models.QuestionnaireRecord) (*models.Questionnaire, error) {
	query := `
		INSERT INTO "questionnaire"("name", "identifier", "experimentID", "interactionTypeID") VALUES($1, $2, $3, $4)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, questionnaire.Name, questionnaire.Identifier, questionnaire.ExperimentID, questionnaire.InteractionTypeID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *QuestionnaireStore) GetByUser(userID string) ([]models.Questionnaire, error) {
	query := s.query + `
		WHERE u."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	questionnaires, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return s.addQuestionsAll(questionnaires)
}

func (s *QuestionnaireStore) GetByExperiment(experimentID string) ([]models.Questionnaire, error) {
	query := s.query + `
		WHERE e."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, experimentID)
	questionnaires, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return s.addQuestionsAll(questionnaires)
}

func (s *QuestionnaireStore) GetRandomByInteraction(interaction *models.Interaction) (*models.Questionnaire, error) {
	query := `
		WITH selected AS (
	` + s.query + `
			LEFT JOIN "livingLab" l ON l."ID" = e."livingLabID"
			WHERE q."interactionTypeID" = $1
			AND e."start" < $2
			AND (e."end" IS NULL OR e."end" > $2)
			AND (l."ID" IS NULL OR l."definition" @> $3)
			ORDER BY random()
			LIMIT 1
		)
		INSERT INTO "assignment" ("userID", "questionnaireID", "interactionID")
		SELECT $4, "questionnaireID", $5
		FROM selected
		WHERE "questionnaireID" IS NOT NULL
		RETURNING "questionnaireID"
		`
	var id string
	row := s.relationalDB.QueryRow(query, interaction.Type.ID, interaction.Timestamp, interaction.Location, interaction.User.ID, interaction.ID)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	questionnaire, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	return s.addQuestions(questionnaire)
}

func (s *QuestionnaireStore) Update(userID string, questionnaireID string, questionnaire *models.QuestionnaireRecord) (*models.Questionnaire, error) {
	query := `
		WITH update_query AS (
			UPDATE "questionnaire" q
			SET "name" = $1, "identifier" = $2, "interactionTypeID" = $4
			FROM "experiment" e
			WHERE q."ID" = $5
			AND e."userID" = $6
			AND e."ID" = $3
			AND e."start" > NOW()
			RETURNING q."ID", q."experimentID"
		)
		SELECT 
			u."ID", 
			CASE 
				WHEN c."start" <= NOW() THEN 'STARTED'
				WHEN u."ID" IS NULL THEN 'WRONG'	
				WHEN u."ID" IS NOT NULL THEN 'OK'
			END AS status
		FROM (
			SELECT "ID", "start"
			FROM "experiment"
			WHERE "ID" = $3 
		) c
		LEFT JOIN update_query u ON c."ID" = u."experimentID"
	`
	var id *string
	var status string
	row := s.relationalDB.QueryRow(query, questionnaire.Name, questionnaire.Identifier, questionnaire.ExperimentID, questionnaire.InteractionTypeID, questionnaireID, userID)
	if err := row.Scan(&id, &status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	switch status {
	case "OK":
		return s.Get(*id)
	case "WRONG": // It is technically possible to allow moving a questionnaire from a non-started experiment to another non-started expirement and that is why we accept experimentID in the input body of the end-point. However, for now we do not allow this as there is no user story for it.
		return nil, &ErrRecordImmutable{message: "experimentID cannot be changed"}
	case "STARTED":
		return nil, &ErrRecordImmutable{message: "experiment already started"}
	}
	return nil, nil
}

func (s *QuestionnaireStore) addQuestions(questionnaire *models.Questionnaire) (*models.Questionnaire, error) {
	questions, err := NewQuestionStore(s.relationalDB).GetByQuestionnaire(questionnaire.ID)
	if err != nil {
		return nil, err
	}
	questionnaire.Questions = questions
	return questionnaire, nil
}

func (s *QuestionnaireStore) addQuestionsAll(questionnaires []models.Questionnaire) ([]models.Questionnaire, error) {
	// This is a potential performance issue because it calls the DB in a loop, but let's cross that bridge when we get there.
	for i := 0; i < len(questionnaires); i++ {
		questionnaire, err := s.addQuestions(&questionnaires[i])
		if err != nil {
			return nil, err
		}
		questionnaires[i] = *questionnaire
	}
	return questionnaires, nil
}

func (s *QuestionnaireStore) Delete(questionID string, userID string) error {
	query := `
		WITH deleted AS (
			DELETE FROM "questionnaire" qq
			USING "experiment" e 
			WHERE qq."ID" = $1
			AND e."userID" = $2
			AND e.start > NOW()
			RETURNING qq."ID"
		)
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM "questionnaire" WHERE "ID" = $1) THEN 'INVALID'
				WHEN NOT EXISTS (SELECT 1 FROM "questionnaire" qq JOIN "experiment" e ON qq."experimentID" = e."ID" WHERE qq."ID" = $1 AND e."userID" = $2) THEN 'USER'
				WHEN EXISTS (SELECT 1 FROM "questionnaire" WHERE "ID" = $1) AND NOT EXISTS (SELECT 1 FROM deleted) THEN 'STARTED'
				WHEN EXISTS (SELECT 1 FROM deleted) THEN 'OK'
			END AS result;
	`
	var state string
	row := s.relationalDB.QueryRow(query, questionID, userID)
	if err := row.Scan(&state); err != nil {
		return err
	}
	switch state {
	case "INVALID":
		return &ErrRecordInattainable{message: "questionnaire was not found"}
	case "USER":
		return &ErrRecordInattainable{message: "questionnaire does not exist for the current user"}
	case "STARTED":
		return &ErrRecordImmutable{message: "cannot delete questionnaire for an experiment that has started"}
	case "OK":
		return nil
	}
	return errors.New("unknown error")
}
