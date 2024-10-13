package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type ExperimentStore Store

func NewExperimentStore(relationalDB *sql.DB) *ExperimentStore {
	s := ExperimentStore{
		relationalDB: relationalDB,
		query: `
		SELECT e."ID", e."name", e."description", e."start", e."end", u."ID", u."name", COALESCE(l."ID", '00000000-0000-0000-0000-000000000000'), COALESCE(l."name", '')
		FROM "experiment" e
		INNER JOIN "user" u ON u."ID" = e."userID"
		LEFT JOIN "livingLab" l ON l."ID" = e."livingLabID"
		`,
	}
	return &s
}

func (s *ExperimentStore) process(rows *sql.Rows, err error) ([]models.Experiment, error) {
	if err != nil {
		return nil, err
	}
	experiments := make([]models.Experiment, 0)
	for rows.Next() {
		var experiment models.Experiment
		var user models.User
		var livingLab models.LivingLab
		if err := rows.Scan(&experiment.ID, &experiment.Name, &experiment.Description, &experiment.Start, &experiment.End, &user.ID, &user.Name, &livingLab.ID, &livingLab.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		experiment.User = user
		if livingLab.ID != "00000000-0000-0000-0000-000000000000" {
			experiment.LivingLab = &livingLab
		}
		experiments = append(experiments, experiment)
	}
	return experiments, nil
}

func (s *ExperimentStore) Get(experimentID string) (*models.Experiment, error) {
	query := s.query + `
		WHERE e."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, experimentID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *ExperimentStore) GetAll() ([]models.Experiment, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *ExperimentStore) Add(userID string, experiment *models.ExperimentRecord) (*models.Experiment, error) {
	query := `
		INSERT INTO "experiment"("name", "description", "start", "end", "userID", "livingLabID") VALUES($1, $2, $3, $4, $5, $6)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, experiment.Name, experiment.Description, experiment.Start, experiment.End, userID, experiment.LivingLabID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *ExperimentStore) GetByUser(userID string) ([]models.Experiment, error) {
	query := s.query + `
		WHERE u."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}
