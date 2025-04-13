package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type BelongingStore Store

func NewBelongingStore(db *sql.DB) *BelongingStore {
	s := BelongingStore{
		relationalDB: db,
		query: `
		SELECT b."ID", b."name", b."category"
		FROM "belonging" b
		`,
	}
	return &s
}

func (s *BelongingStore) process(rows *sql.Rows, err error) ([]models.Belonging, error) {
	if err != nil {
		return nil, err
	}
	belongings := make([]models.Belonging, 0)
	for rows.Next() {
		var b models.Belonging
		if err := rows.Scan(&b.ID, &b.Name, &b.Category); err != nil {
			return nil, err
		}
		belongings = append(belongings, b)
	}
	return belongings, nil
}

func (s *BelongingStore) Get(belongingID string) (*models.Belonging, error) {
	query := s.query + `
		WHERE b."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, belongingID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *BelongingStore) GetAll() ([]models.Belonging, error) {
	query := s.query + `
		ORDER BY b."name"
	`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *BelongingStore) Add(belonging *models.Belonging) (*models.Belonging, error) {
	query := `
		INSERT INTO "belonging"("name", "category") VALUES($1, $2)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, belonging.Name, belonging.Category)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *BelongingStore) Update(belongingID string, belonging *models.Belonging) (*models.Belonging, error) {
	query := `
		UPDATE "belonging" SET "name" = $2, "category" = $3
		WHERE "ID" = $1
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, belongingID, belonging.Name, belonging.Category)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}
