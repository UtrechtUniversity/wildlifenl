package stores

import (
	"database/sql"
	"math/rand"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ConveyanceStore Store

func NewConveyanceStore(db *sql.DB) *ConveyanceStore {
	s := ConveyanceStore{
		relationalDB: db,
		query: `
		SELECT c."ID", c."timestamp", m."ID", m."name", m."severity", m."text", COALESCE(e."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(e."timestamp",'2000-01-01'), COALESCE(e."userLocation",'(0,0)'), COALESCE(e."animalLocation",'(0,0)'), COALESCE(n."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(n."name",''), COALESCE(s."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(s."name",''), COALESCE(s."commonNameNL",''), COALESCE(s."commonNameEN",''), COALESCE(s."encounterMeters",0), COALESCE(s."encounterMinutes",1), r."ID", r."text", q."ID", q."text", q."description", q."index", q."allowMultipleResponse", q."allowOpenResponse", i."ID", i."timestamp", i."description", i."location", t."ID", t."nameNL", t."nameEN", t."descriptionNL", t."descriptionEN", COALESCE(a."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(a."text", ''), COALESCE(a."index", 0), u."ID", u."name"
		FROM "conveyance" c
		INNER JOIN "message" m ON m."ID" = c."messageID"
		LEFT JOIN "encounter" e ON e."ID" = c."encounterID"	
		LEFT JOIN "animal" n ON n."ID" = e."animalID"
		LEFT JOIN "species" s ON s."ID" = n."speciesID"
		LEFT JOIN "response" r ON r."ID" = c."responseID"
		LEFT JOIN "question" q ON q."ID" = r."questionID"
		LEFT JOIN "interaction" i ON i."ID" = r."interactionID"
		LEFT JOIN "interactionType" t ON t."ID" = i."typeID"
		LEFT JOIN "answer" a ON a."ID" = r."answerID"
		INNER JOIN "user" u ON u."ID" = COALESCE(e."userID", i."userID")
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
		var e models.Encounter
		var r models.Response
		var a models.Answer
		var u models.User
		if err := rows.Scan(&c.ID, &c.Timestamp, &c.Message.ID, &c.Message.Severity, &c.Message.Text, &e.ID, &e.Timestamp, &e.UserLocation, &e.AnimalLocation, &e.Animal.ID, &e.Animal.Name, &e.Animal.Species.ID, &e.Animal.Species.Name, &e.Animal.Species.CommonNameNL, &e.Animal.Species.CommonNameEN, &e.Animal.Species.EncounterMeters, &e.Animal.Species.EncounterMinutes, &r.ID, &r.Text, &r.Question.ID, &r.Question.Text, &r.Question.Description, &r.Question.Index, &r.Question.AllowMultipleResponse, &r.Question.AllowOpenResponse, &r.Interaction.ID, &r.Interaction.Timestamp, &r.Interaction.Description, &r.Interaction.Location, &r.Interaction.Type.NameNL, &r.Interaction.Type.NameEN, &r.Interaction.Type.DescriptionNL, &r.Interaction.Type.DescriptionEN, &a.ID, &a.Text, &a.Index, &u.ID, &u.Name); err != nil {
			return nil, err
		}
		if e.ID != "00000000-0000-0000-0000-000000000000" {
			e.User = u
			c.Encounter = &e
		}
		if r.ID != "00000000-0000-0000-0000-000000000000" {
			r.Interaction.User = u
			c.Response = &r
			if a.ID != "00000000-0000-0000-0000-000000000000" {
				r.Answer = &a
			}
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

func (s ConveyanceStore) AddForEncounters(encounters []models.Encounter) (*models.Conveyance, error) {
	messages := make(map[*models.Encounter][]models.Message)
	messageStore := NewMessageStore(s.relationalDB)
	for _, encounter := range encounters {
		m, err := messageStore.GetAllForEncounter(&encounter)
		if err != nil {
			return nil, err
		}
		messages[&encounter] = m
	}
	if len(messages) == 0 {
		return nil, nil
	}
	keys := make([]*models.Encounter, 0, len(messages))
	for k := range messages {
		keys = append(keys, k)
	}
	encounter := keys[rand.Intn(len(keys))]
	message := messages[encounter][rand.Intn(len(messages[encounter]))]
	query := `
		INSERT INTO "conveyance"("messageID", "encounterID") VALUES($1, $2, $3)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, message.ID, encounter.ID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}