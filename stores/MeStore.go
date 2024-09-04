package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type MeStore Store

func NewMeStore(db *sql.DB) *MeStore {
	s := MeStore{
		relationalDB: db,
		query: `
		SELECT u."ID", u."email"
		FROM "user" u
		INNER JOIN "credential" c ON c."email" = u."email"
		`,
	}
	return &s
}

func (s *MeStore) Get(token string) (*models.Me, error) {
	query := s.query + `
		WHERE c."token" = $1
	`
	var userID string
	var email string
	row := s.relationalDB.QueryRow(query, token)
	if err := row.Scan(&userID, &email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	user, err := NewUserStore(s.relationalDB).Get(userID)
	if err != nil {
		return nil, err
	}
	me := models.Me{
		User:  *user,
		Email: email,
	}
	return &me, nil
}
