package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type MessageStore Store

func NewMessageStore(db *sql.DB) *MessageStore {
	s := MessageStore{
		relationalDB: db,
		query: `
		SELECT m."ID", m."name", m."severity", m."text", e."ID", e."name", e."start", e."end", u."ID", u."name", COALESCE(s."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(s."name",''), COALESCE(s."commonNameNL",''), COALESCE(s."commonNameEN",''), COALESCE(a."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(a."text",''), COALESCE(a."index",0)
		FROM "message" m
		INNER JOIN "experiment" e ON e."ID" = m."experimentID"
		INNER JOIN "user" u ON u."ID" = e."userID"
		LEFT JOIN "species" s ON s."ID" = m."speciesID"
		LEFT JOIN "answer" a ON a."ID" = m."answerID"
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
		if err := rows.Scan(&m.ID, &m.Name, &m.Severity, &m.Text, m.Experiment.ID, &m.Experiment.Start, &m.Experiment.End, &m.Experiment.User.ID, &m.Experiment.User.Name, &s.ID, &s.Name, &s.CommonNameNL, &s.CommonNameEN, &a.ID, &a.Text, &a.Index); err != nil {
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
		INSERT INTO "message"("name", "severity", "text", "experimentID", "speciesID", "answerID") VALUES($1, $2, $3, $4, $5, $6)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, message.Name, message.Severity, message.Text, message.ExperimentID, message.SpeciesID, message.AnswerID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *MessageStore) GetByUser(userID string) ([]models.Message, error) {
	query := s.query + `
		WHERE u."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MessageStore) GetAllForEncounter(encounter *models.Encounter) ([]models.Message, error) {
	query := s.query + `
		WHERE m."speciesID" = $1
		AND e."start" < $2
		AND (e."end" IS NULL OR e."end" > $2)
	`
	rows, err := s.relationalDB.Query(query, encounter.Animal.SpeciesID, encounter.Timestamp)
	return s.process(rows, err)
}
