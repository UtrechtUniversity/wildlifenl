package stores

import (
	"database/sql"
	"time"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ZoneStore Store

func NewZoneStore(db *sql.DB) *ZoneStore {
	s := ZoneStore{
		relationalDB: db,
		query: `
		SELECT z."ID", z."deactivated", z."name", z."description", z."area", u."ID", u."name", COALESCE(s."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(s."name",''), COALESCE(s."commonNameNL",''), COALESCE(s."commonNameEN", '')
		FROM "zone" z
		INNER JOIN "user" u ON u."ID" = z."userID"
		LEFT JOIN "zone_species" x ON x."zoneID" = z."ID"
		LEFT JOIN "species" s ON s."ID" = x."speciesID"
		`,
	}
	return &s
}

func (s *ZoneStore) process(rows *sql.Rows, err error) ([]models.Zone, error) {
	if err != nil {
		return nil, err
	}
	zones := make([]models.Zone, 0)
	zone := models.Zone{}
	var zoneID string
	var zoneDeactivated *time.Time
	var zoneName string
	var zoneDescription string
	var zoneArea models.Circle
	for rows.Next() {
		var user models.User
		var species models.Species
		if err := rows.Scan(&zoneID, &zoneDeactivated, &zoneName, &zoneDescription, &zoneArea, &user.ID, &user.Name, &species.ID, &species.Name, &species.CommonNameNL, &species.CommonNameEN); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		if zone.ID != "" && zone.ID != zoneID {
			zones = append(zones, zone)
			zone = models.Zone{}
		}
		zone.ID = zoneID
		zone.Deactivated = zoneDeactivated
		zone.Name = zoneName
		zone.Description = zoneDescription
		zone.Area = zoneArea
		zone.User = user
		if species.ID != "00000000-0000-0000-0000-000000000000" {
			zone.Species = append(zone.Species, species)
		}
	}
	if zone.ID != "" {
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
		AND z."deactivated" IS NULL
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}

func (s *ZoneStore) SetZoneSpecies(zoneID string, speciesIDs []string) (*models.Zone, error) {
	query := `
		DELETE FROM "zone_species" 
		WHERE "zoneID" = $1
		`
	if _, err := s.relationalDB.Exec(query, zoneID); err != nil {
		return nil, err
	}
	query = `
		INSERT INTO "zone_species"("zoneID", "speciesID") VALUES($1, $2)
		`
	for _, speciesID := range speciesIDs {
		if _, err := s.relationalDB.Exec(query, zoneID, speciesID); err != nil {
			return nil, err
		}
	}
	return s.Get(zoneID)
}

func (s *ZoneStore) Deactivate(zoneID string) (*models.Zone, error) {
	query := `
		UPDATE "zone" SET "deactivated" = now()
		WHERE "ID" = $1
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, zoneID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}
