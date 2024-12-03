package stores

import (
	"database/sql"
	"errors"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ExperimentStore Store

func NewExperimentStore(relationalDB *sql.DB) *ExperimentStore {
	s := ExperimentStore{
		relationalDB: relationalDB,
		query: `
		SELECT e."ID", e."name", e."description", e."start", e."end", COALESCE(qc."x", 0), COALESCE(mc."x", 0), COALESCE(qa."x", 0), COALESCE(ca."x", 0), u."ID", u."name", COALESCE(l."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(l."name", '')
		FROM "experiment" e
		INNER JOIN "user" u ON u."ID" = e."userID"
		LEFT JOIN "livingLab" l ON l."ID" = e."livingLabID"
		LEFT JOIN (SELECT "experimentID", COUNT("ID") AS x FROM "questionnaire" GROUP BY "experimentID") qc ON qc."experimentID" = e."ID"
		LEFT JOIN (SELECT "experimentID", COUNT("ID") AS x FROM "message" GROUP BY "experimentID") mc ON mc."experimentID" = e."ID"
		LEFT JOIN (SELECT "experimentID", COUNT(*) AS x FROM (SELECT a."experimentID" AS "experimentID", a."ID" FROM "response" r INNER JOIN "question" q ON q."ID" = r."questionID" INNER JOIN "questionnaire" a ON a."ID" = q."questionnaireID" GROUP BY a."experimentID", a."ID") AS z GROUP BY "experimentID") qa ON qa."experimentID" = e."ID"
		LEFT JOIN (SELECT m."experimentID", COUNT(c."ID") AS x FROM "conveyance" c INNER JOIN "message" m ON m."ID" = c."messageID" GROUP BY m."experimentID") ca ON ca."experimentID" = e."ID"
		`,
	}
	return &s
}

func (s *ExperimentStore) process(rows *sql.Rows, err error) ([]models.Experiment, error) {
	if err != nil {
		return nil, err
	}
	experiments := make([]models.Experiment, 0)
	for rows.Next() {
		var e models.Experiment
		var u models.User
		var l models.LivingLab
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.Start, &e.End, &e.NumberOfQuestionnaires, &e.NumberOfMessages, &e.QuestionnaireActivity, &e.MessageActivity, &u.ID, &u.Name, &l.ID, &l.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		e.User = u
		if l.ID != "00000000-0000-0000-0000-000000000000" {
			e.LivingLab = &l
		}
		experiments = append(experiments, e)
	}
	return experiments, nil
}

func (s *ExperimentStore) Get(experimentID string) (*models.Experiment, error) {
	query := s.query + `
		WHERE e."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, experimentID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ExperimentStore) GetAll() ([]models.Experiment, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *ExperimentStore) Add(userID string, experiment *models.ExperimentRecord) (*models.Experiment, error) {
	query := `
		INSERT INTO "experiment"("name", "description", "start", "end", "userID", "livingLabID") VALUES($1, $2, $3, $4, $5, $6)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, experiment.Name, experiment.Description, experiment.Start, experiment.End, userID, experiment.LivingLabID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *ExperimentStore) GetByUser(userID string) ([]models.Experiment, error) {
	query := s.query + `
		WHERE u."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}

func (s *ExperimentStore) Update(userID string, experimentID string, experiment *models.ExperimentRecord) (*models.Experiment, error) {
	query := `
		WITH update_query AS (
			UPDATE "experiment"
			SET "name" = $1, "description" = $2, "start" = $3, "end" = $4, "livingLabID" = $5
			WHERE "ID" = $6 AND "userID" = $7 AND "start" > NOW()
			RETURNING "ID"
		)
		SELECT 
			c."ID", 
			CASE 
				WHEN u."ID" IS NOT NULL THEN 'OK'
				WHEN c."start" <= NOW() THEN 'STARTED'
			END AS status
		FROM (
			SELECT "ID", "start"
			FROM "experiment"
			WHERE "ID" = $6 
		) c
		LEFT JOIN update_query u ON c."ID" = u."ID"
	`
	var id string
	var status string
	row := s.relationalDB.QueryRow(query, experiment.Name, experiment.Description, experiment.Start, experiment.End, experiment.LivingLabID, experimentID, userID)
	if err := row.Scan(&id, &status); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	switch status {
	case "STARTED":
		return nil, &CannotUpdateError{message: "experiment already started"}
	case "OK":
		return s.Get(id)
	}
	return nil, nil
}

func (s *ExperimentStore) EndNow(userID string, experimentID string) (*models.Experiment, error) {
	query := `
		UPDATE "experiment"
		SET "end" = NOW()
		WHERE "ID" = $1 AND "userID" = $2
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, experimentID, userID)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}

func (s *ExperimentStore) Delete(questionID string, userID string) error {
	query := `
		WITH deleted AS (
			DELETE FROM "experiment" 
			WHERE "ID" = $1
			AND "userID" = $2
			AND start > NOW()
			RETURNING "ID"
		)
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM "experiment" WHERE "ID" = $1) THEN 'INVALID'
				WHEN NOT EXISTS (SELECT 1 FROM "experiment" WHERE "ID" = $1 AND "userID" = $2) THEN 'USER'
				WHEN EXISTS (SELECT 1 FROM "experiment" WHERE "ID" = $1) AND NOT EXISTS (SELECT 1 FROM deleted) THEN 'STARTED'
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
		return &CannotUpdateError{message: "experiment was not found"}
	case "USER":
		return &CannotUpdateError{message: "experiment does not exist for the current user"}
	case "STARTED":
		return &CannotUpdateError{message: "cannot delete an experiment that has started"}
	case "OK":
		return nil
	}
	return errors.New("unknown error")
}
