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
		SELECT i."id", i."createdAt", i."description", i."latitude", i."longitude", s."id", s."name", s."commonNameNL", s."commonNameEN", u."id", u."name"
		FROM interaction i
		INNER JOIN "species" s ON s."id" = i."speciesID"
		INNER JOIN "user" u ON u."id" = i."userID"
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
		var interaction models.Interaction
		var species models.Species
		var user models.User
		if err := rows.Scan(&interaction.ID, &interaction.CreatedAt, &interaction.Description, &interaction.Latitude, &interaction.Longitude, &species.ID, &species.Name, &species.CommonNameNL, &species.CommonNameEN, &user.ID, &user.Name); err != nil {
			return nil, err
		}
		interaction.Species = species
		interaction.User = user
		interactions = append(interactions, interaction)
	}
	return interactions, nil
}

func (s *InteractionStore) Get(interactionID string) (*models.Interaction, error) {
	query := s.query + `
		WHERE i."id" = $1
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

func (s *InteractionStore) GetByUser(userID string) ([]models.Interaction, error) {
	query := s.query + `
		WHERE u."id" = $1
		ORDER BY i."createdAt" DESC
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}

func (s *InteractionStore) Add(userID string, interaction *models.InteractionRecord) (*models.Interaction, error) {
	query := `
		INSERT INTO interaction("description", "latitude", "longitude", "speciesID", "userID") VALUES($1, $2, $3, $4, $5)
		RETURNING "id"
	`
	var id string
	row := s.relationalDB.QueryRow(query, interaction.Description, interaction.Latitude, interaction.Longitude, interaction.SpeciesID, userID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
