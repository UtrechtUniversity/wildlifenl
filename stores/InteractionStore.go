package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type InteractionStore Store

func NewInteractionStore(db *sql.DB) *InteractionStore {
	s := InteractionStore{
		relationalDB: db,
		query: `
		SELECT i."ID", i."timestamp", i."description", i."location", i."moment", s."ID", s."name", s."commonName", u."ID", u."name", t."ID", t."name", t."description"
		FROM "interaction" i
		INNER JOIN "species" s ON s."ID" = i."speciesID"
		INNER JOIN "user" u ON u."ID" = i."userID"
		LEFT JOIN "interactionType" t ON t."ID" = i."typeID"
		`,
	}
	return &s
}

func (s *InteractionStore) process(rows *sql.Rows, err error) ([]models.Interaction, error) {
	if err != nil {
		return nil, err
	}
	interactions := make([]models.Interaction, 0)
	for rows.Next() {
		var i models.Interaction
		var s models.Species
		var u models.User
		var t models.InteractionType
		if err := rows.Scan(&i.ID, &i.Timestamp, &i.Description, &i.Location, &i.Moment, &s.ID, &s.Name, &s.CommonName, &u.ID, &u.Name, &t.ID, &t.Name, &t.Description); err != nil {
			return nil, err
		}
		i.Species = s
		i.User = u
		i.Type = t
		interactions = append(interactions, i)
	}
	return interactions, nil
}

func (s *InteractionStore) Get(interactionID string) (*models.Interaction, error) {
	query := s.query + `
		WHERE i."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, interactionID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *InteractionStore) GetAll() ([]models.Interaction, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *InteractionStore) Add(userID string, interaction *models.InteractionRecord) (*models.Interaction, error) {
	query := `
		INSERT INTO "interaction"("description", "location", "moment", "speciesID", "userID", "typeID") VALUES($1, $2, $3, $4, $5, $6)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, interaction.Description, interaction.Location, interaction.Moment, interaction.SpeciesID, userID, interaction.TypeID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *InteractionStore) GetByUser(userID string) ([]models.Interaction, error) {
	query := s.query + `
		WHERE u."ID" = $1
		ORDER BY i."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}
