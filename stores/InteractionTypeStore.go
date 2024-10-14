package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type InteractionTypeStore Store

func NewInteractionTypeStore(db *sql.DB) *InteractionTypeStore {
	s := InteractionTypeStore{
		relationalDB: db,
		query: `
		SELECT t."ID", t."nameNL", t."nameEN", t."descriptionNL", t."descriptionEN"
		FROM "interactionType" t
		`,
	}
	return &s
}

func (s *InteractionTypeStore) process(rows *sql.Rows, err error) ([]models.InteractionType, error) {
	if err != nil {
		return nil, err
	}
	interactionTypes := make([]models.InteractionType, 0)
	for rows.Next() {
		var interactionType models.InteractionType
		if err := rows.Scan(&interactionType.ID, &interactionType.NameNL, &interactionType.NameEN, &interactionType.DescriptionNL, &interactionType.DescriptionEN); err != nil {
			return nil, err
		}
		interactionTypes = append(interactionTypes, interactionType)
	}
	return interactionTypes, nil
}

func (s *InteractionTypeStore) Get(interactionTypeID int) (*models.InteractionType, error) {
	query := s.query + `
		WHERE t."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, interactionTypeID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *InteractionTypeStore) GetAll() ([]models.InteractionType, error) {
	query := s.query + `
		ORDER BY "ID"
	`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *InteractionTypeStore) Add(interactionType *models.InteractionType) (*models.InteractionType, error) {
	query := `
		INSERT INTO "interactionType"("nameNL", "nameEN", "descriptionNL", "descriptionEN") VALUES($1, $2, $3, $4)
		RETURNING "ID"
	`
	var id int
	row := s.relationalDB.QueryRow(query, interactionType.NameNL, interactionType.NameEN, interactionType.DescriptionNL, interactionType.DescriptionEN)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *InteractionTypeStore) Update(interactionTypeID int, interactionType *models.InteractionType) (*models.InteractionType, error) {
	query := `
		UPDATE "interactionType" SET "nameNL" = $2, "nameEN" = $3, "descriptionNL" = $4, "descriptionEN" = $5
		WHERE "ID" = $1
		RETURNING "ID"
	`
	var id int
	row := s.relationalDB.QueryRow(query, interactionTypeID, interactionType.NameNL, interactionType.NameEN, interactionType.DescriptionNL, interactionType.DescriptionEN)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
