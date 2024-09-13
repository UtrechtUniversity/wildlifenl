package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type QuestionnaireStore Store

func NewQuestionnaireStore(db *sql.DB) *QuestionnaireStore {
	s := QuestionnaireStore{
		relationalDB: db,
		query: `
		SELECT q."ID", q."name", e."ID", e."name", e."start", e."end", t."ID", t."nameNL", t."nameEN", t."descriptionNL", t."descriptionEN", u."ID", u."name"
		FROM "questionnaire" q
		INNER JOIN "experiment" e ON e."ID" = q."experimentID"
		INNER JOIN "interactionType" t ON t."ID" = q."interactionTypeID"
		INNER JOIN "user" u ON u."ID" = e."userID"
		`,
	}
	return &s
}

func (s *QuestionnaireStore) process(rows *sql.Rows, err error) ([]models.Questionnaire, error) {
	if err != nil {
		return nil, err
	}
	questionnaires := make([]models.Questionnaire, 0)
	for rows.Next() {
		var questionnaire models.Questionnaire
		var experiment models.Experiment
		var interactionType models.InteractionType
		var user models.User
		if err := rows.Scan(&questionnaire.ID, &questionnaire.Name, &experiment.ID, &experiment.Name, &experiment.Start, &experiment.End, &interactionType.ID, &interactionType.NameNL, &interactionType.NameEN, &interactionType.DescriptionNL, &interactionType.DescriptionEN, &user.ID, &user.Name); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		experiment.User = user
		questionnaire.Experiment = experiment
		questionnaire.InteractionType = interactionType
		questionnaires = append(questionnaires, questionnaire)
	}
	return questionnaires, nil
}

func (s *QuestionnaireStore) Get(questionnaireID string) (*models.Questionnaire, error) {
	query := s.query + `
		WHERE q."ID" = $1
	`
	rows, err := s.relationalDB.Query(query, questionnaireID)
	result, err := s.process(rows, err)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *QuestionnaireStore) GetAll() ([]models.Questionnaire, error) {
	rows, err := s.relationalDB.Query(s.query)
	return s.process(rows, err)
}

func (s *QuestionnaireStore) Add(userID string, questionnaire *models.QuestionnaireRecord) (*models.Questionnaire, error) {
	query := `
		INSERT INTO "questionnaire"("name", "experimentID", "interactionTypeID") VALUES($1, $2, $3)
		RETURNING "ID"
	`
	var id string
	row := s.relationalDB.QueryRow(query, questionnaire.Name, questionnaire.ExperimentID, questionnaire.InteractionTypeID)
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *QuestionnaireStore) GetByUser(userID string) ([]models.Questionnaire, error) {
	query := s.query + `
		WHERE u."ID" = $1
		`
	rows, err := s.relationalDB.Query(query, userID)
	return s.process(rows, err)
}
