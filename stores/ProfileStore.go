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
		SELECT u."ID", u."name", u."email", u."location", u."locationTimestamp", u."dateOfBirth", u."gender", u."postcode", u."notes", u."natureVisitAvgWeeklyFrequency", u."recreationAppTandC", u."reportAppTandC", COALESCE(u."firebaseCloudMessagingToken", ''), COALESCE(r."ID", 0), COALESCE(r."name", '')
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
	profiles := make([]models.Profile, 0)
	var p models.Profile
	for rows.Next() {
		var id string
		var name string
		var email string
		var location *models.Point
		var locationTimestamp *time.Time
		var dateOfBirth *string
		var gender *string
		var postcode *string
		var notes *string
		var natureVisitAvgWeeklyFrequency int
		var recreationAppTandC bool
		var reportAppTandC bool
		var firebaseCloudMessagingToken string
		var r models.Role
		if err := rows.Scan(&id, &name, &email, &location, &locationTimestamp, &dateOfBirth, &gender, &postcode, &notes, &natureVisitAvgWeeklyFrequency, &recreationAppTandC, &reportAppTandC, &firebaseCloudMessagingToken, &r.ID, &r.Name); err != nil {
			return nil, err
		}
		if p.ID != id {
			if p.ID != "" {
				profiles = append(profiles, p)
				p = models.Profile{}
			}
			p.ID = id
			p.Name = name
			p.Email = email
			p.Location = location
			p.LocationTimestamp = locationTimestamp
			p.DateOfBirth = dateOfBirth
			p.Gender = gender
			p.Postcode = postcode
			p.Notes = notes
			p.NatureVisitAvgWeeklyFrequency = natureVisitAvgWeeklyFrequency
			p.RecreationAppTandC = recreationAppTandC
			p.ReportAppTandC = reportAppTandC
			if firebaseCloudMessagingToken != "" {
				p.FirebaseCloudMessagingToken = &firebaseCloudMessagingToken
			}
		}
		if r.ID > 0 {
			p.Roles = append(p.Roles, r)
		}
	}
	if p.ID != "" {
		profiles = append(profiles, p)
	}
	return profiles, nil
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
	query := s.query + `
		WHERE u."ID"::TEXT || '@wildlifenl.nl' IS DISTINCT FROM u."email"
	` // Closed accounts have an email equal to their ID + @wildlifenl.nl, this excludes closed accounts.
	query += `
		ORDER BY u."ID"
	`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *ProfileStore) Update(profileID string, profile *models.ProfileRecord) (*models.Profile, error) {
	query := `
		UPDATE "user" SET "name" = $2, "dateOfBirth" = $3, "gender" = $4, "postcode" = $5, "notes" = $6, "natureVisitAvgWeeklyFrequency" = $7, "recreationAppTandC" = $8, "reportAppTandC" = $9, "firebaseCloudMessagingToken" = $10
		WHERE "ID" = $1
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, profileID, profile.Name, profile.DateOfBirth, profile.Gender, profile.Postcode, profile.Notes, profile.NatureVisitAvgWeeklyFrequency, profile.RecreationAppTandC, profile.ReportAppTandC, profile.FirebaseCloudMessagingToken)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}

func (s *ProfileStore) Delete(profileID string) error {
	query := `
		UPDATE "user" SET "name" = $1, "email" = $2
		WHERE "ID" = $1
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, profileID, profileID+"@wildlifenl.nl")
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}
