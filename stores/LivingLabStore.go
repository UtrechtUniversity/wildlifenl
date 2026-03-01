package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type LivingLabStore Store

func NewLivingLabStore(db *sql.DB) *LivingLabStore {
	s := LivingLabStore{
		relationalDB: db,
		query: `
		SELECT l."ID", l."name", l."definition"
		FROM "livingLab" l
		`,
	}
	return &s
}

func (s *LivingLabStore) process(rows *sql.Rows, err error) ([]models.LivingLab, error) {
	if err != nil {
		return nil, err
	}
	livingLabs := make([]models.LivingLab, 0)
	for rows.Next() {
		var livingLab models.LivingLab
		if err := rows.Scan(&livingLab.ID, &livingLab.Name, &livingLab.Definition); err != nil {
			return nil, err
		}
		livingLabs = append(livingLabs, livingLab)
	}
	return livingLabs, nil
}

func (s *LivingLabStore) Get(livingLabID string) (*models.LivingLab, error) {
	query := s.query + `
		WHERE l."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, livingLabID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *LivingLabStore) GetAll() ([]models.LivingLab, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *LivingLabStore) Add(livingLab *models.LivingLab) (*models.LivingLab, error) {
	query := `
		INSERT INTO "livingLab"("name", "definition") VALUES($1, $2)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, livingLab.Name, livingLab.Definition)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *LivingLabStore) Update(livingLabID string, livingLab *models.LivingLab) (*models.LivingLab, error) {
	query := `
		UPDATE "livingLab" SET "name" = $2, "definition" = $3
		WHERE "ID" = $1
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, livingLabID, livingLab.Name, livingLab.Definition)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}
