package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type UserStore Store

func NewUserStore(db *sql.DB) *UserStore {
	s := UserStore{
		relationalDB: db,
		query: `
		SELECT u."ID", u."name"
		FROM "user" u
		`,
	}
	return &s
}

func (s *UserStore) process(rows *sql.Rows, err error) ([]models.User, error) {
	if err != nil {
		return nil, err
	}
	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (s *UserStore) Get(userID string) (*models.User, error) {
	query := s.query + `
		WHERE u."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *UserStore) GetAll() ([]models.User, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}
