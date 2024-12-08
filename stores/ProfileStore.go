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
		SELECT u."ID", u."name", u."email", u."location", u."locationTimestamp", u."dateOfBirth", u."gender", u."postcode", u."description", COALESCE(r."ID", 0), COALESCE(r."name", '')
		FROM "user" u
		LEFT JOIN "user_role" x ON x."userID" = u."ID"
		LEFT JOIN "role" r ON r."ID" = x."roleID"
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
		var userDateOfBirth *string
		var userGender *string
		var userPostcode *string
		var userDescription *string
		var r models.Role
		if err := rows.Scan(&userID, &userName, &userEmail, &userLocation, &userLocationTimestamp, &userDateOfBirth, &userGender, &userPostcode, &userDescription, &r.ID, &r.Name); err != nil {
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
		user.DateOfBirth = userDateOfBirth
		user.Gender = userGender
		user.Postcode = userPostcode
		user.Description = userDescription
		if r.ID > 0 {
			user.Roles = append(user.Roles, r)
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

func (s *ProfileStore) Update(profileID string, profile *models.ProfileRecord) (*models.Profile, error) {
	query := `
		UPDATE "user" SET "name" = $2, "dateOfBirth" = $3, "gender" = $4, "postcode" = $5, "description" = $6
		WHERE "ID" = $1
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, profileID, profile.Name, profile.DateOfBirth, profile.Gender, profile.Postcode, profile.Description)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}
