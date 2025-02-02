package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type AlarmStore Store

func NewAlarmStore(db *sql.DB) *AlarmStore {
	s := AlarmStore{
		relationalDB: db,
		query: `
		SELECT a."ID", a."timestamp", z."ID", z."deactivated", z."name", z."description", z."area", s."ID", s."name", s."commonName", u."ID", u."name", COALESCE(i."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(i."timestamp", '2000-01-01'), COALESCE(i."description",''), COALESCE(i."location",'(0,0)'), COALESCE(t."ID",0), COALESCE(t."name",''), COALESCE(t."description",''), COALESCE(x."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(x."name",''), COALESCE(d."ID",0), COALESCE(d."location", '(0,0)'), COALESCE(d."timestamp",'2000-01-01'), COALESCE(d."sensorID",''), COALESCE(n."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(n."name",''), COALESCE(n."location",'(0,0)'), COALESCE(c."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(c."timestamp",'2000-01-01'), COALESCE(m."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(m."name",''), COALESCE(m."severity",0), COALESCE(m."text",''), COALESCE(m."trigger",''), COALESCE(e."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(e."name",''), COALESCE(e."description",''), COALESCE(e."start",'2000-01-01'), COALESCE(e."end",'2000-01-01'), COALESCE(y."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(y."name",'')
		FROM "alarm" a
		INNER JOIN "zone" z ON a."zoneID" = z."ID"
		INNER JOIN "user" u ON u."ID" = z."userID"
		LEFT JOIN "interaction" i ON a."interactionID" = i."ID"
		LEFT JOIN "interactionType" t ON t."ID" = i."typeID"
		LEFT JOIN "user" x ON x."ID" = i."userID"
		LEFT JOIN "detection" d ON a."detectionID" = d."ID"
		LEFT JOIN "animal" n ON a."animalID" = n."ID"
		INNER JOIN "species" s ON s."ID" = COALESCE(i."speciesID", d."speciesID", n."speciesID")
		LEFT JOIN "conveyance" c ON c."alarmID" = a."ID"
		LEFT JOIN "message" m ON m."ID" = c."messageID"
		LEFT JOIN "experiment" e ON e."ID" = m."experimentID"
		LEFT JOIN "user" y ON y."ID" = e."userID"
		`,
	}
	return &s
}

func (s *AlarmStore) process(rows *sql.Rows, err error) ([]models.Alarm, error) {
	if err != nil {
		return nil, err
	}
	alarms := make([]models.Alarm, 0)
	var alarm models.Alarm
	for rows.Next() {
		var a models.Alarm
		var i models.Interaction
		var d models.Detection
		var n models.Animal
		var s models.Species
		var c models.Conveyance
		if err := rows.Scan(&a.ID, &a.Timestamp, &a.Zone.ID, &a.Zone.Deactivated, &a.Zone.Name, &a.Zone.Description, &a.Zone.Area, &s.ID, &s.Name, &s.CommonName, &a.Zone.User.ID, &a.Zone.User.Name, &i.ID, &i.Timestamp, &i.Description, &i.Location, &i.Type.ID, &i.Type.Name, &i.Type.Description, &i.User.ID, &i.User.Name, &d.ID, &d.Location, &d.Timestamp, &d.SensorID, &n.ID, &n.Name, &n.Location, &c.ID, &c.Timestamp, &c.Message.ID, &c.Message.Name, &c.Message.Severity, &c.Message.Text, &c.Message.Trigger, &c.Message.Experiment.ID, &c.Message.Experiment.Name, &c.Message.Experiment.Description, &c.Message.Experiment.Start, &c.Message.Experiment.End, &c.Message.Experiment.User.ID, &c.Message.Experiment.User.Name); err != nil {
			return nil, err
		}
		if alarm.ID != a.ID {
			if alarm.ID != "" {
				alarms = append(alarms, alarm)
			}
			alarm = a
			if i.ID != "00000000-0000-0000-0000-000000000000" {
				i.Species = s
				alarm.Interaction = &i
			}
			if d.ID > 0 {
				d.Species = s
				alarm.Detection = &d
			}
			if n.ID != "00000000-0000-0000-0000-000000000000" {
				n.Species = s
				alarm.Animal = &n
			}
		}
		if c.ID != "00000000-0000-0000-0000-000000000000" {
			c.User = alarm.Zone.User
			alarm.Conveyances = append(alarm.Conveyances, c)
		}
	}
	if alarm.ID != "" {
		alarms = append(alarms, alarm)
	}
	return alarms, nil
}

func (s *AlarmStore) Get(alarmID int) (*models.Alarm, error) {
	query := s.query + `
		WHERE a."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, alarmID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *AlarmStore) GetAll() ([]models.Alarm, error) {
	query := s.query + `
		ORDER BY a."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *AlarmStore) GetByUser(userID string) ([]models.Alarm, error) {
	query := s.query + `
		WHERE u."ID" = $1
		ORDER BY a."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}

func (s *AlarmStore) AddAllFromDetection(detection *models.Detection) ([]string, error) {
	query := `
		INSERT INTO "alarm"("zoneID", "detectionID")
		SELECT z."ID", d."ID"
		FROM "zone" z
		INNER JOIN "user" u ON u."ID" = z."userID"
		LEFT JOIN "zone_species" x ON x."zoneID" = z."ID"
		LEFT JOIN "species" s ON s."ID" = x."speciesID"
		LEFT JOIN "detection" d ON z."area" @> d."location" AND d."speciesID" = s."ID"
		WHERE d."ID" = $1
		AND z."deactivated" IS NULL
		AND z."created" < d."timestamp"
		RETURNING "ID"
	`
	ids := make([]string, 0)
	rows, err := s.relationalDB.Query(query, detection.ID)
	if err != nil {
		return nil, err
	}
	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *AlarmStore) AddAllFromInteraction(interaction *models.Interaction) ([]string, error) {
	query := `
		INSERT INTO "alarm"("zoneID", "interactionID")
		SELECT z."ID", i."ID"
		FROM "zone" z
		INNER JOIN "user" u ON u."ID" = z."userID"
		LEFT JOIN "zone_species" x ON x."zoneID" = z."ID"
		LEFT JOIN "species" s ON s."ID" = x."speciesID"
		LEFT JOIN "interaction" i ON z."area" @> i."location" AND i."speciesID" = s."ID"
		WHERE i."ID" = $1
		AND z."deactivated" IS NULL
		AND z."created" < i."timestamp"
		RETURNING "ID"
	`
	ids := make([]string, 0)
	rows, err := s.relationalDB.Query(query, interaction.ID)
	if err != nil {
		return nil, err
	}
	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *AlarmStore) AddAllFromAnimal(animal *models.Animal) ([]string, error) {
	query := `
		INSERT INTO "alarm"("zoneID", "animalID")
		SELECT z."ID", n."ID"
		FROM "zone" z
		INNER JOIN "user" u ON u."ID" = z."userID"
		LEFT JOIN "zone_species" x ON x."zoneID" = z."ID"
		LEFT JOIN "species" s ON s."ID" = x."speciesID"
		LEFT JOIN "animal" n ON z."area" @> n."location" AND n."speciesID" = s."ID"
		WHERE n."ID" = $1
		AND z."deactivated" IS NULL
		AND z."created" < n."locationTimestamp"
		RETURNING "ID"
	`
	ids := make([]string, 0)
	rows, err := s.relationalDB.Query(query, animal.ID)
	if err != nil {
		return nil, err
	}
	var id string
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
