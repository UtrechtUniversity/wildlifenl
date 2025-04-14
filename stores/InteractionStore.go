package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type InteractionStore Store

func NewInteractionStore(db *sql.DB) *InteractionStore {
	s := InteractionStore{
		relationalDB: db,
		query: `
		SELECT i."ID", i."timestamp", i."description", i."location", i."moment", i."place", s."ID", s."name", s."commonName", u."ID", u."name", t."ID", t."name", t."description", COALESCE(dr."impactType",''), COALESCE(dr."impactValue",0), COALESCE(dr."estimatedDamage",0), COALESCE(dr."estimatedLoss",0), COALESCE(cr."estimatedDamage",0), COALESCE(cr."intensity",''), COALESCE(cr."urgency",''), COALESCE(b."ID",'00000000-0000-0000-0000-000000000000'), COALESCE(b."name",''), COALESCE(b."category",'')FROM "interaction" i
		INNER JOIN "species" s ON s."ID" = i."speciesID"
		INNER JOIN "user" u ON u."ID" = i."userID"
		LEFT JOIN "interactionType" t ON t."ID" = i."typeID"
		LEFT JOIN "sightingReport" sr ON i."ID" = sr."interactionID" AND i."typeID" = 1
		LEFT JOIN "damageReport" dr ON i."ID" = dr."interactionID" AND i."typeID" = 2
		LEFT JOIN "collisionReport" cr ON i."ID" = cr."interactionID" AND i."typeID" = 3
		LEFT JOIN "belonging" b ON b."ID" = dr."belongingID" AND i."typeID" = 2
		`,
	}
	return &s
}

func (s *InteractionStore) process(rows *sql.Rows, err error) ([]models.Interaction, error) {
	if err != nil {
		return nil, err
	}
	interactions := make([]models.Interaction, 0)
	for rows.Next() {
		var i models.Interaction
		var sr models.SightingReport
		var dr models.DamageReport
		var cr models.CollisionReport
		if err := rows.Scan(&i.ID, &i.Timestamp, &i.Description, &i.Location, &i.Moment, &i.Place, &i.Species.ID, &i.Species.Name, &i.Species.CommonName, &i.User.ID, &i.User.Name, &i.Type.ID, &i.Type.Name, &i.Type.Description, &dr.ImpactType, &dr.ImpactValue, &dr.EstimatedDamage, &dr.EstimatedLoss, &cr.EstimatedDamage, &cr.Intensity, &cr.Urgency, &dr.Belonging.ID, &dr.Belonging.Name, &dr.Belonging.Category); err != nil {
			return nil, err
		}
		if i.Type.ID == 1 {
			involvedAnimals, err := NewAnimalInfoStore(s.relationalDB).GetAllForInteraction(i.ID) // possible performance issue because it is being called in a loop.
			if err != nil {
				return nil, err
			}
			sr.InvolvedAnimals = involvedAnimals
			i.ReportOfSighting = &sr
		}
		if i.Type.ID == 2 {
			i.ReportOfDamage = &dr
		}
		if i.Type.ID == 3 {
			involvedAnimals, err := NewAnimalInfoStore(s.relationalDB).GetAllForInteraction(i.ID) // possible performance issue because it is being called in a loop.
			if err != nil {
				return nil, err
			}
			cr.InvolvedAnimals = involvedAnimals
			i.ReportOfCollision = &cr
		}
		interactions = append(interactions, i)
	}
	return interactions, nil
}

func (s *InteractionStore) Get(interactionID string) (*models.Interaction, error) {
	query := s.query + `
		WHERE i."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, interactionID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *InteractionStore) GetAll() ([]models.Interaction, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *InteractionStore) Add(userID string, interaction *models.InteractionRecord) (*models.Interaction, error) {
	// This query works if sent directly to postgres, but apparently does not when doing so via the
	// Go connector: pq: cannot insert multiple commands into a prepared statement
	// So let's do them one by one, which is not so nice....
	/*
		query := `
			DROP TABLE IF EXISTS inserted;
			CREATE TEMP TABLE inserted ("ID" UUID, "typeID" INT);
			WITH "insert" as (
				INSERT INTO "interaction"("description", "location", "moment", "speciesID", "userID", "typeID") VALUES($1, $2, $3, $4, $5, $6)
				RETURNING "ID", "typeID"
			)
			INSERT INTO inserted
			SELECT "ID", "typeID"
			FROM "insert";
			INSERT INTO "sightingReport" ("interactionID")
			SELECT "ID" FROM "inserted" WHERE "typeID" = 1;
			INSERT INTO "damageReport" ("interactionID", "belongingID", "impactType", "impactValue", "estimatedDamage", "estimatedLoss")
			SELECT "ID", $7, $8, $9, $10, $11 FROM inserted WHERE "typeID" = 2;
			INSERT INTO "collisionReport" ("interactionID", "estimatedDamage", "intensity", "urgency")
			SELECT "ID", $12, $13, $14 FROM inserted WHERE "typeID" = 3;
			SELECT "ID" FROM inserted;
			DROP TABLE inserted;
		`
		if interaction.ReportOfSighting == nil {
			interaction.ReportOfSighting = &models.SightingReport{}
		}
		if interaction.ReportOfDamage == nil {
			interaction.ReportOfDamage = &models.DamageReport{}
		}
		if interaction.ReportOfCollision == nil {
			interaction.ReportOfCollision = &models.CollisionReport{}
		}
		var id string
		row := s.relationalDB.QueryRow(query, interaction.Description, interaction.Location, interaction.Moment, interaction.SpeciesID, userID, interaction.TypeID, interaction.ReportOfDamage.Belonging.ID, interaction.ReportOfDamage.ImpactType, interaction.ReportOfDamage.ImpactValue, interaction.ReportOfDamage.EstimatedDamage, interaction.ReportOfDamage.EstimatedLoss, interaction.ReportOfCollision.EstimatedDamage, interaction.ReportOfCollision.Intensity, interaction.ReportOfCollision.Urgency)
		if err := row.Scan(&id); err != nil {
			return nil, err
		}
	*/
	query := `
		INSERT INTO "interaction"("description", "location", "moment", "place", "speciesID", "userID", "typeID") VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, interaction.Description, interaction.Location, interaction.Moment, interaction.Place, interaction.SpeciesID, userID, interaction.TypeID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	switch interaction.TypeID {
	case 1:
		query = `
			INSERT INTO "sightingReport" ("interactionID") VALUES ($1)
		`
		if _, err := s.relationalDB.Exec(query, id); err != nil {
			return nil, err
		}
		if err := NewAnimalInfoStore(s.relationalDB).addMany(interaction.ReportOfSighting.InvolvedAnimals, id); err != nil {
			return nil, err
		}
	case 2:
		query = `
			INSERT INTO "damageReport" ("interactionID", "belongingID", "impactType", "impactValue", "estimatedDamage", "estimatedLoss") VALUES ($1, $2, $3, $4, $5, $6)
		`
		if _, err := s.relationalDB.Exec(query, id, interaction.ReportOfDamage.Belonging.ID, interaction.ReportOfDamage.ImpactType, interaction.ReportOfDamage.ImpactValue, interaction.ReportOfDamage.EstimatedDamage, interaction.ReportOfDamage.EstimatedLoss); err != nil {
			return nil, err
		}
	case 3:
		query = `
			INSERT INTO "collisionReport" ("interactionID", "estimatedDamage", "intensity", "urgency") VALUES ($1, $2, $3, $4)
		`
		if _, err := s.relationalDB.Exec(query, id, interaction.ReportOfCollision.EstimatedDamage, interaction.ReportOfCollision.Intensity, interaction.ReportOfCollision.Urgency); err != nil {
			return nil, err
		}
		if err := NewAnimalInfoStore(s.relationalDB).addMany(interaction.ReportOfCollision.InvolvedAnimals, id); err != nil {
			return nil, err
		}
	}
	return s.Get(id)
}

func (s *InteractionStore) GetByUser(userID string) ([]models.Interaction, error) {
	query := s.query + `
		WHERE u."ID" = $1
		ORDER BY i."timestamp" DESC
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}
