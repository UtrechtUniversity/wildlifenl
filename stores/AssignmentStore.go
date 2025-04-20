package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type AssignmentStore Store

func NewAssignmentStore(db *sql.DB) *AssignmentStore {
	s := AssignmentStore{
		relationalDB: db,
		query: `
		SELECT i."ID", i."timestamp", i."description", i."location", i."moment", i."place", t."ID", t."name", t."description", s."ID", s."name", s."commonName", u."ID", u."name", q."ID", q."name", q."identifier", e."ID", e."name", e."start", e."end", qt."ID", qt."name", qt."description", eu."ID", eu."name"
		FROM "assignment" a
		INNER JOIN "interaction" i ON i."ID" = a."interactionID"
		INNER JOIN "interactionType" t ON t."ID" = i."typeID"
		INNER JOIN "species" s ON s."ID" = i."speciesID"
		INNER JOIN "user" u ON u."ID" = a."userID"
		INNER JOIN "questionnaire" q ON q."ID" = a."questionnaireID"
		INNER JOIN "experiment" e ON e."ID" = q."experimentID"
		INNER JOIN "interactionType" qt ON qt."ID" = q."interactionTypeID"
		INNER JOIN "user" eu ON eu."ID" = e."userID"
		`,
	}
	return &s
}

func (s *AssignmentStore) process(rows *sql.Rows, err error) ([]models.Assignment, error) {
	if err != nil {
		return nil, err
	}
	assignments := make([]models.Assignment, 0)
	for rows.Next() {
		var a models.Assignment
		if err := rows.Scan(&a.Interaction.ID, &a.Interaction.Timestamp, &a.Interaction.Description, &a.Interaction.Location, &a.Interaction.Moment, &a.Interaction.Place, &a.Interaction.Type.ID, &a.Interaction.Type.Name, &a.Interaction.Type.Description, &a.Interaction.Species.ID, &a.Interaction.Species.Name, &a.Interaction.Species.CommonName, &a.Interaction.User.ID, &a.Interaction.User.Name, &a.Questionnaire.ID, &a.Questionnaire.Name, &a.Questionnaire.Identifier, &a.Questionnaire.Experiment.ID, &a.Questionnaire.Experiment.Name, &a.Questionnaire.Experiment.Start, &a.Questionnaire.Experiment.End, &a.Questionnaire.InteractionType.ID, &a.Questionnaire.InteractionType.Name, &a.Questionnaire.InteractionType.Description, &a.Questionnaire.Experiment.User.ID, &a.Questionnaire.Experiment.User.Name); err != nil {
			return nil, err
		}
		a.User = a.Interaction.User
		assignments = append(assignments, a)
	}
	return assignments, nil
}

func (s *AssignmentStore) GetByUser(userID string) ([]models.Assignment, error) {
	query := s.query + `
		WHERE a."userID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	assignments, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	return assignments, nil
}
