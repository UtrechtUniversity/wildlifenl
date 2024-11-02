package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type NewQuestionnaireInput struct {
	Input
	Body *models.QuestionnaireRecord `json:"questionnaire"`
}

type QuestionnaireHolder struct {
	Body *models.Questionnaire `json:"questionnaire"`
}

type QuestionnairesHolder struct {
	Body []models.Questionnaire `json:"questionnaires"`
}

type questionnaireOperations Operations

func newQuestionnaireOperations() *questionnaireOperations {
	return &questionnaireOperations{Endpoint: "questionnaire"}
}

func (o *questionnaireOperations) RegisterGet(api huma.API) {
	name := "Get Questionnaire By ID"
	description := "Retrieve a specific questionnaire by ID."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of this questionnaire." format:"uuid"`
	}) (*QuestionnaireHolder, error) {
		questionnaire, err := stores.NewQuestionnaireStore(relationalDB).Get(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		if questionnaire == nil {
			return nil, generateNotFoundByIDError(o.Endpoint, input.ID)
		}
		return &QuestionnaireHolder{Body: questionnaire}, nil
	})
}

func (o *questionnaireOperations) RegisterGetAll(api huma.API) {
	name := "Get All Questionnaires"
	description := "Retrieve all questionnaires."
	path := "/" + o.Endpoint + "s/"
	scopes := []string{"administrator"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct{}) (*QuestionnairesHolder, error) {
		questionnaires, err := stores.NewQuestionnaireStore(relationalDB).GetAll()
		if err != nil {
			return nil, handleError(err)
		}
		return &QuestionnairesHolder{Body: questionnaires}, nil
	})
}

func (o *questionnaireOperations) RegisterAdd(api huma.API) {
	name := "Add Questionnaire"
	description := "Add a new questionnaire."
	path := "/" + o.Endpoint + "/"
	scopes := []string{"researcher"}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *NewQuestionnaireInput) (*QuestionnaireHolder, error) {
		experiments, err := stores.NewExperimentStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		var experiment models.Experiment
		for _, e := range experiments {
			if e.ID == input.Body.ExperimentID {
				experiment = e
				break
			}
		}
		if experiment.ID == "" {
			return nil, generateNotFoundForThisUserError("experiment", input.Body.ExperimentID)
		}
		questionnaire, err := stores.NewQuestionnaireStore(relationalDB).Add(input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		return &QuestionnaireHolder{Body: questionnaire}, nil
	})
}

func (o *questionnaireOperations) RegisterGetByExperiment(api huma.API) {
	name := "Get Questionnaires By Experiment"
	description := "Retrieve all questionnaires by experimentID."
	path := "/" + o.Endpoint + "s/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the experiment to retrieve questionnaires for." format:"uuid"`
	}) (*QuestionnairesHolder, error) {
		questionnaires, err := stores.NewQuestionnaireStore(relationalDB).GetByExperiment(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &QuestionnairesHolder{Body: questionnaires}, nil
	})
}

/*
func (o *questionnaireOperations) RegisterGetMine(api huma.API) {
	name := "Get My Questionnaires"
	description := "Retrieve my questionnaires."
	path := "/" + o.Endpoint + "s/me/"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *Input) (*QuestionnairesHolder, error) {
		questionnaires, err := stores.NewQuestionnaireStore(relationalDB).GetByUser(input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return &QuestionnairesHolder{Body: questionnaires}, nil
	})
}
*/
