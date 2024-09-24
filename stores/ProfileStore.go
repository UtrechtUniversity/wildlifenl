package stores

import (
	"database/sql"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ProfileStore Store

func NewProfileStore(db *sql.DB) *ProfileStore {
	s := ProfileStore{
		relationalDB: db,
		query: `
		SELECT u."ID", u."name", u."email", u."location", u."locationTimestamp", COALESCE(r."ID", 0), COALESCE(r."name", ''), COALESCE(l."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(l."name", '')
		FROM "user" u
		LEFT JOIN "user_role" x ON x."userID" = u."ID"
		LEFT JOIN "role" r ON r."ID" = x."roleID"
		LEFT JOIN "livingLab" l ON l."ID" = u."livingLabID"
		`,
	}
	return &s
}

func (s *ProfileStore) process(rows *sql.Rows, err error) ([]models.Profile, error) {
	if err != nil {
		return nil, err
	}
	users := make([]models.Profile, 0)
	var user models.Profile
	for rows.Next() {
		var userID string
		var userName string
		var userEmail string
		var userLocation *models.Point
		var userLocationTimestamp *time.Time
		var r models.Role
		var l models.LivingLab
		if err := rows.Scan(&userID, &userName, &userEmail, &userLocation, &userLocationTimestamp, &r.ID, &r.Name, &l.ID, &l.Name); err != nil {
			return nil, err
		}
		if user.ID != "" && user.ID != userID {
			users = append(users, user)
			user = models.Profile{}
		}
		user.ID = userID
		user.Name = userName
		user.Email = userEmail
		user.Location = userLocation
		user.LocationTimestamp = userLocationTimestamp
		if r.ID > 0 {
			user.Roles = append(user.Roles, r)
		}
		if l.ID != "00000000-0000-0000-0000-000000000000" {
			user.LivingLab = &l
		}
	}
	if user.ID != "" {
		users = append(users, user)
	}
	return users, nil
}

func (s *ProfileStore) Get(userID string) (*models.Profile, error) {
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

func (s *ProfileStore) GetAll() ([]models.Profile, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *ProfileStore) GetByCredentialToken(token string) (*models.Profile, error) {
	query := `
		SELECT u."ID", u."email"
		FROM "user" u
		INNER JOIN "credential" c ON c."email" = u."email"
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
	if userID != "" {
		return s.Get(userID)
	}
	return nil, nil
}
