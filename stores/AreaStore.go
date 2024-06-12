package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type AreaStore Store

func NewAreaStore(db *sql.DB) *AreaStore {
	s := AreaStore{
		db: db,
		query: `
		SELECT a."id", a."description", a."definition", u."id", u."name"
		FROM area a
		INNER JOIN "user" u ON u."id" = a."userID"
		`,
	}
	return &s
}

func (s *AreaStore) process(rows *sql.Rows, err error) ([]models.Area, error) {
	if err != nil {
		return nil, err
	}
	notices := make([]models.Area, 0)
	for rows.Next() {
		var area models.Area
		var user models.User
		if err := rows.Scan(&area.ID, &area.Description, &area.Definition, &user.ID, &user.Name); err != nil {
			return nil, err
		}
		area.User = user
		notices = append(notices, area)
	}
	return notices, nil
}

func (s *AreaStore) Get(areaID string) (*models.Area, error) {
	query := s.query + `
		WHERE a."id" = $1
		`
	rows, err := s.db.Query(query, areaID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *AreaStore) GetAll() ([]models.Area, error) {
	rows, err := s.db.Query(s.query)
	return s.process(rows, err)
}

func (s *AreaStore) GetByUser(userID string) ([]models.Area, error) {
	query := s.query + `
		WHERE u."id" = $1
		`
	rows, err := s.db.Query(query, userID)
	return s.process(rows, err)
}

func (s *AreaStore) Add(userID string, notice *models.AreaRecord) (*models.Area, error) {
	query := `
		INSERT INTO area("description", "definition", "userID") VALUES($1, $2, $3)
		RETURNING "id"
	`
	var id string
	row := s.db.QueryRow(query, notice.Description, notice.Definition, userID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
