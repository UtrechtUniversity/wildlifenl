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
		SELECT l."id", l."name", l."definition"
		FROM "livingLab" l
		`,
	}
	return &s
}

func (s *LivingLabStore) process(rows *sql.Rows, err error) ([]models.LivingLab, error) {
	if err != nil {
		return nil, err
	}
	livinglabs := make([]models.LivingLab, 0)
	for rows.Next() {
		var livinglab models.LivingLab
		if err := rows.Scan(&livinglab.ID, &livinglab.Name, &livinglab.Definition); err != nil {
			return nil, err
		}
		livinglabs = append(livinglabs, livinglab)
	}
	return livinglabs, nil
}

func (s *LivingLabStore) Get(livinglabID string) (*models.LivingLab, error) {
	query := s.query + `
		WHERE l."id" = $1
		`
	rows, err := s.relationalDB.Query(query, livinglabID)
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

func (s *LivingLabStore) Add(livinglab *models.LivingLab) (*models.LivingLab, error) {
	query := `
		INSERT INTO livingLab("name", "definition") VALUES($1, $2)
		RETURNING "id"
	`
	var id string
	row := s.relationalDB.QueryRow(query, livinglab.Name, livinglab.Definition)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
