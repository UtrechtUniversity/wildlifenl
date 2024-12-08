package stores

import (
	"database/sql"
	"errors"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type MessageStore Store

func NewMessageStore(db *sql.DB) *MessageStore {
	s := MessageStore{
		relationalDB: db,
		query: `
		SELECT m."ID", m."name", m."severity", m."text", m."trigger", m."encounterMeters", m."encounterMinutes", COALESCE(c."x", 0), e."ID", e."name", e."description", e."start", e."end", u."ID", u."name", COALESCE(s."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(s."name",''), COALESCE(s."commonName",''), COALESCE(a."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(a."text",''), COALESCE(a."index",0)
		FROM "message" m
		INNER JOIN "experiment" e ON e."ID" = m."experimentID"
		INNER JOIN "user" u ON u."ID" = e."userID"
		LEFT JOIN "species" s ON s."ID" = m."speciesID"
		LEFT JOIN "answer" a ON a."ID" = m."answerID"
		LEFT JOIN (SELECT "messageID", COUNT("ID") AS x FROM "conveyance" GROUP BY "messageID") c ON c."messageID" = m."ID"
		`,
	}
	return &s
}

func (s *MessageStore) process(rows *sql.Rows, err error) ([]models.Message, error) {
	if err != nil {
		return nil, err
	}
	messages := make([]models.Message, 0)
	for rows.Next() {
		var m models.Message
		var s models.Species
		var a models.Answer
		if err := rows.Scan(&m.ID, &m.Name, &m.Severity, &m.Text, &m.Trigger, &m.EncounterMeters, &m.EncounterMinutes, &m.Activity, &m.Experiment.ID, &m.Experiment.Name, &m.Experiment.Description, &m.Experiment.Start, &m.Experiment.End, &m.Experiment.User.ID, &m.Experiment.User.Name, &s.ID, &s.Name, &s.CommonName, &a.ID, &a.Text, &a.Index); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		if s.ID != "00000000-0000-0000-0000-000000000000" {
			m.Species = &s
		}
		if a.ID != "00000000-0000-0000-0000-000000000000" {
			m.Answer = &a
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (s *MessageStore) Get(messageID string) (*models.Message, error) {
	query := s.query + `
		WHERE m."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, messageID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *MessageStore) GetAll() ([]models.Message, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *MessageStore) Add(message *models.MessageRecord) (*models.Message, error) {
	query := `
		INSERT INTO "message"("name", "severity", "text", "trigger", "encounterMeters", "encounterMinutes", "experimentID", "speciesID", "answerID") VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, message.Name, message.Severity, message.Text, message.Trigger, message.EncounterMeters, message.EncounterMinutes, message.ExperimentID, message.SpeciesID, message.AnswerID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *MessageStore) GetByExperiment(experimentID string) ([]models.Message, error) {
	query := s.query + `
		WHERE e."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, experimentID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MessageStore) Delete(messageID string, userID string) error {
	query := `
		WITH deleted AS (
			DELETE FROM "message" m
			USING "experiment" e 
			WHERE m."ID" = $1
			AND e."userID" = $2
			AND e.start > NOW()
			RETURNING m."ID"
		)
		SELECT 
			CASE 
				WHEN NOT EXISTS (SELECT 1 FROM "message" WHERE "ID" = $1) THEN 'INVALID'
				WHEN NOT EXISTS (SELECT 1 FROM "message" m JOIN "experiment" e ON m."experimentID" = e."ID" WHERE m."ID" = $1 AND e."userID" = $2) THEN 'USER'
				WHEN EXISTS (SELECT 1 FROM "message" WHERE "ID" = $1) AND NOT EXISTS (SELECT 1 FROM deleted) THEN 'STARTED'
				WHEN EXISTS (SELECT 1 FROM deleted) THEN 'OK'
			END AS result;
	`
	var state string
	row := s.relationalDB.QueryRow(query, messageID, userID)
	if err := row.Scan(&state); err != nil {
		return err
	}
	switch state {
	case "INVALID":
		return &ErrRecordInattainable{message: "message was not found"}
	case "USER":
		return &ErrRecordInattainable{message: "message does not exist for the current user"}
	case "STARTED":
		return &ErrRecordImmutable{message: "cannot delete message for an experiment that has started"}
	case "OK":
		return nil
	}
	return errors.New("unknown error")
}
