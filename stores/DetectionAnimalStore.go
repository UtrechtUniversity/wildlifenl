package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type DetectionAnimalStore Store

func NewDetectionAnimalStore(db *sql.DB) *DetectionAnimalStore {
	s := DetectionAnimalStore{
		relationalDB: db,
		query: `
		SELECT a."confidence", a."sex", a."lifeStage", a."condition", a."behaviour", a."description"
		FROM "detectionAnimal" a
		`,
	}
	return &s
}

func (s *DetectionAnimalStore) process(rows *sql.Rows, err error) ([]models.DetectionAnimal, error) {
	if err != nil {
		return nil, err
	}
	animalInfos := make([]models.DetectionAnimal, 0)
	for rows.Next() {
		var a models.DetectionAnimal
		if err := rows.Scan(&a.Confidence, &a.Sex, &a.LifeStage, &a.Condition, &a.Behaviour, &a.Description); err != nil {
			return nil, err
		}
		animalInfos = append(animalInfos, a)
	}
	return animalInfos, nil
}

func (s *DetectionAnimalStore) GetAllForDetection(detectionID string) ([]models.DetectionAnimal, error) {
	query := s.query + `
		WHERE a."detectionID" = $1
	`
	rows, err := s.relationalDB.Query(query, detectionID)
	return s.process(rows, err)
}

func (s *DetectionAnimalStore) addMany(detectionAnimals []models.DetectionAnimal, detectionID string) error {
	query := `
		INSERT INTO "detectionAnimal"("detectionID", "confidence", "sex", "lifeStage", "condition", "behaviour", "description") VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING "detectionID"
	`
	for _, a := range detectionAnimals {
		var id string
		row := s.relationalDB.QueryRow(query, detectionID, a.Confidence, a.Sex, a.LifeStage, a.Condition, a.Behaviour, a.Description)
		if err := row.Scan(&id); err != nil {
			return err
		}
	}
	return nil
}
