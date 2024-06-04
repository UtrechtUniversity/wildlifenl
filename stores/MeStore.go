package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type MeStore Store

func NewMeStore(db *sql.DB) *MeStore {
	s := MeStore{
		db: db,
		query: `
		SELECT u."id", u."name", u."email"
		FROM "user" u
		INNER JOIN credential c	ON c."email" = u."email"
		`,
	}
	return &s
}

func (s *MeStore) Get(token string) (*models.Me, error) {
	query := s.query + `
		WHERE c."token" = $1
	`
	var me models.Me
	row := s.db.QueryRow(query, token)
	if err := row.Scan(&me.ID, &me.Name, &me.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &me, nil
}
