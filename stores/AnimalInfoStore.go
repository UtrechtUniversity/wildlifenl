package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type AnimalInfoStore Store

func NewAnimalInfoStore(db *sql.DB) *AnimalInfoStore {
	s := AnimalInfoStore{
		relationalDB: db,
		query: `
		SELECT a."sex", a."lifeStage", a."condition"
		FROM "animalInfo" a
		`,
	}
	return &s
}

func (s *AnimalInfoStore) process(rows *sql.Rows, err error) ([]models.AnimalInfo, error) {
	if err != nil {
		return nil, err
	}
	animalInfos := make([]models.AnimalInfo, 0)
	for rows.Next() {
		var a models.AnimalInfo
		if err := rows.Scan(&a.Sex, &a.LifeStage, &a.Condition); err != nil {
			return nil, err
		}
		animalInfos = append(animalInfos, a)
	}
	return animalInfos, nil
}

func (s *AnimalInfoStore) GetAllForInteraction(interactionID string) ([]models.AnimalInfo, error) {
	query := s.query + `
		WHERE a."interactionID" = $1
	`
	rows, err := s.relationalDB.Query(query, interactionID)
	return s.process(rows, err)
}

func (s *AnimalInfoStore) addMany(animalInfos []models.AnimalInfo, interactionID string) error {
	query := `
		INSERT INTO "animalInfo"("interactionID", "sex", "lifeStage", "condition") VALUES ($1, $2, $3, $4)
		RETURNING "interactionID"
	`
	for _, animalInfo := range animalInfos {
		var id string
		row := s.relationalDB.QueryRow(query, interactionID, animalInfo.Sex, animalInfo.LifeStage, animalInfo.Condition)
		if err := row.Scan(&id); err != nil {
			return err
		}
	}
	return nil
}
