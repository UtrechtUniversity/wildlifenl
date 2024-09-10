package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ZoneStore Store

func NewZoneStore(db *sql.DB) *ZoneStore {
	s := ZoneStore{
		relationalDB: db,
		query: `
		SELECT z."ID", z."name", z."description", z."area", u."ID", u."name"
		FROM "zone" z
		INNER JOIN "user" u ON u."ID" = z."userID"
		`,
	}
	return &s
}

func (s *ZoneStore) process(rows *sql.Rows, err error) ([]models.Zone, error) {
	if err != nil {
		return nil, err
	}
	zones := make([]models.Zone, 0)
	for rows.Next() {
		var zone models.Zone
		var user models.User
		if err := rows.Scan(&zone.ID, &zone.Name, &zone.Description, &zone.Area, &user.ID, &user.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		zone.User = user
		zones = append(zones, zone)
	}
	return zones, nil
}

func (s *ZoneStore) Get(zoneID string) (*models.Zone, error) {
	query := s.query + `
		WHERE z."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, zoneID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ZoneStore) GetAll() ([]models.Zone, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *ZoneStore) Add(userID string, zone *models.ZoneRecord) (*models.Zone, error) {
	query := `
		INSERT INTO "zone"("name", "description", "area", "userID") VALUES($1, $2, $3, $4)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, zone.Name, zone.Description, zone.Area, userID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *ZoneStore) GetByUser(userID string) ([]models.Zone, error) {
	query := s.query + `
		WHERE u."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}
