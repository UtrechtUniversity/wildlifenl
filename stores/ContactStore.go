package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ContactStore Store

func NewContactStore(relationalDB *sql.DB) *ContactStore {
	s := ContactStore{
		relationalDB: relationalDB,
		query: `
			SELECT c."ID", c."start", c."end", d."sensorID", d."contactHardwareAddress", d."start", d."end", a."ID", a."name", a."location", s."ID", s."name", s."commonName", u."ID", u."name"
			FROM "contact" c
			INNER JOIN "borneSensorDeployment" d ON d."ID" = c."borneSensorDeploymentID"
			INNER JOIN "animal" a ON a."ID" = d."animalID"
			INNER JOIN "species" s ON s."ID" = a."speciesID"
			INNER JOIN "user" u ON u."ID" = c."userID"
		`,
	}
	return &s
}

func (s *ContactStore) process(rows *sql.Rows, err error) ([]models.Contact, error) {
	if err != nil {
		return nil, err
	}
	contacts := make([]models.Contact, 0)
	for rows.Next() {
		var c models.Contact
		if err := rows.Scan(&c.ID, &c.Start, &c.End, &c.BorneSensorDeployment.SensorID, &c.BorneSensorDeployment.ContactHardwareAddress, &c.BorneSensorDeployment.Start, &c.BorneSensorDeployment.End, &c.BorneSensorDeployment.Animal.ID, &c.BorneSensorDeployment.Animal.Name, &c.BorneSensorDeployment.Animal.Location, &c.BorneSensorDeployment.Animal.Species.ID, &c.BorneSensorDeployment.Animal.Species.Name, &c.BorneSensorDeployment.Animal.Species.CommonName, &c.User.ID, &c.User.Name); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (s *ContactStore) Get(contactID string) (*models.Contact, error) {
	query := s.query + `
		WHERE c."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, contactID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ContactStore) GetForUser(userID string) ([]models.Contact, error) {
	query := s.query + `
		WHERE c."userID" = $1
		ORDER BY c."start" DESC
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}

func (s *ContactStore) NotExists(userID string, contactHardwareAddress string) (bool, error) {
	query := `
		SELECT c."ID"
		FROM "contact" c
		INNER JOIN "borneSensorDeployment" d ON d."ID" = c."borneSensorDeploymentID"
		WHERE c."userID" = $1
		AND d."contactHardwareAddress" = $2
		AND c."end" IS NULL
	`
	var id string
	row := s.relationalDB.QueryRow(query, userID, contactHardwareAddress)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (s *ContactStore) Exists(userID string, contactID string) (bool, error) {
	query := `
		SELECT c."ID"
		FROM "contact" c
		WHERE c."userID" = $1
		AND c."ID" = $2
		AND c."end" IS NULL
	`
	var id string
	row := s.relationalDB.QueryRow(query, userID, contactID)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *ContactStore) Add(userID string, contactHardwareAddress string) (*models.Contact, error) {
	query := `
		INSERT INTO "contact" ("userID", "borneSensorDeploymentID")
		SELECT $1, d."ID"
		FROM "borneSensorDeployment" d
		WHERE d."contactHardwareAddress" = $2
		AND d."start" <= now()
		AND d."end" IS NULL
		AND NOT EXISTS (
			SELECT 1 FROM "contact" c
			WHERE c."userID" = $1
			AND c."borneSensorDeploymentID" = d."ID"
			AND c."end" IS NULL
		)
		ORDER BY d."start" DESC
		LIMIT 1
		RETURNING "ID";
	`
	var id string
	row := s.relationalDB.QueryRow(query, userID, contactHardwareAddress)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}

func (s *ContactStore) End(userID string, contactID string) (*models.Contact, error) {
	query := `
		UPDATE "contact" SET "end" = now()
		WHERE "userID" = $1
		AND "ID" = $2 
		RETURNING "ID";
	`
	var id string
	row := s.relationalDB.QueryRow(query, userID, contactID)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}
