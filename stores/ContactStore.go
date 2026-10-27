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
			SELECT c."ID", c."start", c."end", c."startLocation", c."endLocation", d."sensorID", d."contactHardwareAddress", d."start", d."end", a."ID", a."name", a."location", s."ID", s."name", s."commonName", u."ID", u."name"
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
		if err := rows.Scan(&c.ID, &c.Start, &c.End, &c.StartLocation, &c.EndLocation, &c.BorneSensorDeployment.SensorID, &c.BorneSensorDeployment.ContactHardwareAddress, &c.BorneSensorDeployment.Start, &c.BorneSensorDeployment.End, &c.BorneSensorDeployment.Animal.ID, &c.BorneSensorDeployment.Animal.Name, &c.BorneSensorDeployment.Animal.Location, &c.BorneSensorDeployment.Animal.Species.ID, &c.BorneSensorDeployment.Animal.Species.Name, &c.BorneSensorDeployment.Animal.Species.CommonName, &c.User.ID, &c.User.Name); err != nil {
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

func (s *ContactStore) GetByExperiment(experimentID string) ([]models.Contact, error) {
	query := s.query + `
		INNER JOIN "experiment" e ON e."ID" = $1
		LEFT JOIN "livingLab" l ON l."ID" = e."livingLabID"
		WHERE (
			(c."start" >= e."start" AND (e."end" IS NULL OR c."start" <= e."end"))
			OR
			(c."end" >= e."start" AND (e."end" IS NULL OR c."end" <= e."end"))
		) AND (
			l."ID" IS NULL
			OR
			c."startLocation" <@ l."definition"
			OR
			c."endLocation" <@ l."definition"
		)
		ORDER BY c."start" DESC
		`
	rows, err := s.relationalDB.Query(query, experimentID)
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

func (s *ContactStore) Add(userID string, contactHardwareAddress string, location *models.Point) (*models.Contact, error) {
	query := `
		INSERT INTO "contact" ("userID", "borneSensorDeploymentID", "startLocation")
		SELECT $1, d."ID", $3
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
	row := s.relationalDB.QueryRow(query, userID, contactHardwareAddress, location)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}

func (s *ContactStore) End(userID string, contactID string, location *models.Point) (*models.Contact, error) {
	query := `
		UPDATE "contact" SET "end" = now(), "endLocation" = $3
		WHERE "userID" = $1
		AND "ID" = $2
		AND "end" IS NULL
		RETURNING "ID";
	`
	var id string
	row := s.relationalDB.QueryRow(query, userID, contactID, location)
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return s.Get(id)
}
