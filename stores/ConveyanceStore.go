package stores

import (
	"database/sql"
	"strings"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/lib/pq"
)

type ConveyanceStore Store

func NewConveyanceStore(db *sql.DB) *ConveyanceStore {
	s := ConveyanceStore{
		relationalDB: db,
		query: `
		SELECT c."ID", c."timestamp", u."ID", u."name", m."ID", m."name", m."severity", m."text", m."trigger", m."encounterMeters", m."encounterMinutes", e."ID", e."name", e."description", e."start", e."end", uu."ID", uu."name", COALESCE(ll."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(ll."name", ''), COALESCE(n."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(n."name",''), COALESCE(s."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(s."name",''), COALESCE(s."commonName",''), COALESCE(r."ID", '00000000-0000-0000-0000-000000000000'), r."text", COALESCE(q."ID", '00000000-0000-0000-0000-000000000000'), q."text", q."description", q."index", q."allowMultipleResponse", q."allowOpenResponse", COALESCE(i."ID", '00000000-0000-0000-0000-000000000000'), i."timestamp", i."description", i."location", COALESCE(t."ID", '00000000-0000-0000-0000-000000000000'), t."name", t."description", COALESCE(ts."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(ts."name",''), COALESCE(ts."commonName",''), COALESCE(a."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(a."text",''), COALESCE(a."index",0), COALESCE(l."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(l."timestamp",'2000-01-01'), COALESCE(z."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(z."deactivated",'200-01-01'), COALESCE(z."name",''), COALESCE(z."description",''), COALESCE(z."area",'<(0,0),1>'), COALESCE(ms."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(ms."name",''), COALESCE(ms."commonName",'') 
		FROM "conveyance" c
		INNER JOIN "user" u ON u."ID" = c."userID"
		INNER JOIN "message" m ON m."ID" = c."messageID"
		INNER JOIN "experiment" e ON e."ID" = m."experimentID"
		INNER JOIN "user" uu ON uu."ID" = e."userID"
		LEFT JOIN "livingLab" ll ON ll."ID" = e."livingLabID"
		LEFT JOIN "species" ms ON ms."ID" = m."speciesID"
		LEFT JOIN "animal" n ON n."ID" = c."animalID"
		LEFT JOIN "species" s ON s."ID" = n."speciesID"
		LEFT JOIN "response" r ON r."ID" = c."responseID"
		LEFT JOIN "question" q ON q."ID" = r."questionID"
		LEFT JOIN "interaction" i ON i."ID" = r."interactionID"
		LEFT JOIN "interactionType" t ON t."ID" = i."typeID"
		LEFT JOIN "species" ts ON ts."ID" = i."speciesID"
		LEFT JOIN "answer" a ON a."ID" = r."answerID"
		LEFT JOIN "alarm" l ON l."ID" = c."alarmID"
		LEFT JOIN "zone" z ON z."ID" = l."zoneID"
		`,
	}
	return &s
}

func (s *ConveyanceStore) process(rows *sql.Rows, err error) ([]models.Conveyance, error) {
	if err != nil {
		return nil, err
	}
	conveyances := make([]models.Conveyance, 0)
	for rows.Next() {
		var c models.Conveyance
		var n models.Animal
		var r models.Response
		var a models.Answer
		var l models.Alarm
		var ll models.LivingLab
		var s models.Species
		if err := rows.Scan(&c.ID, &c.Timestamp, &c.User.ID, &c.User.Name, &c.Message.ID, &c.Message.Name, &c.Message.Severity, &c.Message.Text, &c.Message.Trigger, &c.Message.EncounterMeters, &c.Message.EncounterMinutes, &c.Message.Experiment.ID, &c.Message.Experiment.Name, &c.Message.Experiment.Description, &c.Message.Experiment.Start, &c.Message.Experiment.End, &c.Message.Experiment.User.ID, &c.Message.Experiment.User.Name, &ll.ID, &ll.Name, &n.ID, &n.Name, &n.Species.ID, &n.Species.Name, &n.Species.CommonName, &r.ID, &r.Text, &r.Question.ID, &r.Question.Text, &r.Question.Description, &r.Question.Index, &r.Question.AllowMultipleResponse, &r.Question.AllowOpenResponse, &r.Interaction.ID, &r.Interaction.Timestamp, &r.Interaction.Description, &r.Interaction.Location, &r.Interaction.Type.ID, &r.Interaction.Type.Name, &r.Interaction.Type.Description, &r.Interaction.Species.ID, &r.Interaction.Species.Name, &r.Interaction.Species.CommonName, &a.ID, &a.Text, &a.Index, &l.ID, &l.Timestamp, &l.Zone.ID, &l.Zone.Deactivated, &l.Zone.Name, &l.Zone.Description, &l.Zone.Area, &s.ID, &s.Name, &s.CommonName); err != nil {
			return nil, err
		}
		if n.ID != "00000000-0000-0000-0000-000000000000" {
			c.Animal = &n
		}
		if r.ID != "00000000-0000-0000-0000-000000000000" {
			r.Interaction.User = c.User
			c.Response = &r
			if a.ID != "00000000-0000-0000-0000-000000000000" {
				r.Answer = &a
			}
		}
		if l.ID != "00000000-0000-0000-0000-000000000000" {
			c.Alarm = &l
		}
		if ll.ID != "00000000-0000-0000-0000-000000000000" {
			c.Message.Experiment.LivingLab = &ll
		}
		if s.ID != "00000000-0000-0000-0000-000000000000" {
			c.Message.Species = &s
		}
		conveyances = append(conveyances, c)
	}
	return conveyances, nil
}

func (s *ConveyanceStore) Get(conveyanceID string) (*models.Conveyance, error) {
	query := s.query + `
		WHERE c."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, conveyanceID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ConveyanceStore) GetAll() ([]models.Conveyance, error) {
	query := s.query + `
		ORDER BY c."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *ConveyanceStore) GetByUser(userID string) ([]models.Conveyance, error) {
	query := s.query + `
		WHERE u."ID" = $1
		ORDER BY c."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}

func (s *ConveyanceStore) AddForResponse(response *models.Response) (*models.Conveyance, error) {
	query := `
		WITH inserted AS (
			INSERT INTO "conveyance"("userID", "messageID", "responseID")
			SELECT u."ID", m."ID", r."ID"
			FROM "response" r
			INNER JOIN "interaction" i ON r."interactionID" = i."ID"
			INNER JOIN "user" u ON u."ID" = i."userID"
			LEFT JOIN "answer" a ON r."answerID" = a."ID"
			INNER JOIN "message" m ON m."answerID" = a."ID"
			LEFT JOIN "experiment" e ON e."ID" = m."experimentID"
			LEFT JOIN "livingLab" l ON l."ID" = e."livingLabID"
			WHERE r."ID" = $1
			AND e."start" < i."timestamp"
			AND (e."end" IS NULL OR e."end" > i."timestamp")
			AND (l."ID" IS NULL OR l."definition" @> i."location")
			ORDER BY RANDOM()
			LIMIT 1
			RETURNING "ID", "timestamp", "userID", "messageID", "animalID", "responseID", "alarmID"
		)
	` + strings.Replace(s.query, "FROM \"conveyance\"", "FROM inserted", 1)
	rows, err := s.relationalDB.Query(query, response.ID)
	if err != nil {
		return nil, err
	}
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ConveyanceStore) AddForTrackingReading(trackingReading *models.TrackingReading) (*models.Conveyance, error) {
	query := `
		WITH inserted AS (
			INSERT INTO "conveyance"("userID", "messageID", "animalID")
			SELECT u."ID", m."ID", n."ID"
			FROM "animal" n
			INNER JOIN "species" s ON s."ID" = n."speciesID"
			INNER JOIN "message" m ON m."speciesID" = s."ID"
			INNER JOIN "experiment" e ON e."ID" = m."experimentID"
			INNER JOIN "user" u ON u."ID" = $1
			LEFT JOIN "livingLab" l ON l."ID" = e."livingLabID"
			WHERE ABS(extract(epoch FROM u."locationTimestamp" - n."locationTimestamp")) / 60 < m."encounterMinutes" 
			AND CIRCLE(n."location", CAST(m."encounterMeters" AS FLOAT) / 10000) @> u."location"
			AND e."start" < u."locationTimestamp"
			AND (e."end" IS NULL OR e."end" > u."locationTimestamp")
			AND (l."ID" IS NULL OR l."definition" @> $2)
			ORDER BY RANDOM()
			LIMIT 1
			RETURNING "ID", "timestamp", "userID", "messageID", "animalID", "responseID", "alarmID"
		)
	` + strings.Replace(s.query, "FROM \"conveyance\"", "FROM inserted", 1)
	rows, err := s.relationalDB.Query(query, trackingReading.UserID, trackingReading.Location)
	if err != nil {
		return nil, err
	}
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ConveyanceStore) AddForAlarmIDs(alarmIDs []string) error {
	if len(alarmIDs) < 1 {
		return nil
	}
	query := `
		INSERT INTO "conveyance"("userID", "messageID", "alarmID")
		SELECT u."ID", m."ID", a."ID"
		FROM "alarm" a
		INNER JOIN "zone" z ON a."zoneID" = z."ID"
		INNER JOIN "user" u ON u."ID" = z."userID"
		LEFT JOIN "interaction" i ON i."ID" = a."interactionID"
		LEFT JOIN "detection" d ON d."ID" = a."detectionID"
		LEFT JOIN "animal" n ON n."ID" = a."animalID"
		INNER JOIN "species" s ON s."ID" = COALESCE(i."speciesID", d."speciesID", n."speciesID")
		INNER JOIN "message" m ON m."speciesID" = s."ID" AND m."trigger" = 'alarm'
		INNER JOIN "experiment" e ON e."ID" = m."experimentID"
		LEFT JOIN "livingLab" l ON l."ID" = e."livingLabID"
		WHERE a."ID" = ANY($1)
		AND e."start" < a."timestamp"
		AND (e."end" IS NULL OR e."end" > a."timestamp")
		AND (l."ID" IS NULL OR l."definition" @> @@z."area")
	`
	if _, err := s.relationalDB.Exec(query, pq.Array(alarmIDs)); err != nil {
		return err
	}
	return nil
}

func (s *ConveyanceStore) GetByExperiment(experimentID string) ([]models.Conveyance, error) {
	query := s.query + `
		WHERE e."ID" = $1
		ORDER BY c."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query, experimentID)
	return s.process(rows, err)
}
