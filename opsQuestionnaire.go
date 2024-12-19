package wildlifenl

import (
	"context"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
	"github.com/danielgtaylor/huma/v2"
)

type QuestionnaireHolder struct {
	Body *models.Questionnaire `json:"questionnaire"`
}

type QuestionnairesHolder struct {
	Body []models.Questionnaire `json:"questionnaires"`
}

type QuestionnaireAddInput struct {
	Input
	Body *models.QuestionnaireRecord `json:"questionnaire"`
}

type QuestionnaireUpdateInput struct {
	Input
	ID   string                      `path:"id" format:"uuid" doc:"The ID of the questionnaire to be updated."`
	Body *models.QuestionnaireRecord `json:"questionnaire"`
}

type QuestionnaireDeleteInput struct {
	Input
	ID string `path:"id" format:"uuid" doc:"The ID of the questionnaire to be deleted."`
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
		ID string `path:"id" format:"uuid" doc:"The ID of the questionnaire to retrieve."`
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

func (o *questionnaireOperations) RegisterGetByExperimentDeprecated(api huma.API) {
	name := "Get Questionnaires By Experiment [deprecated]"
	description := "Retrieve all questionnaires by experimentID. DEPRECATED: Use Get Questionnaires By Experiment instead."
	path := "/" + o.Endpoint + "s/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the experiment to retrieve questionnaires for."`
	}) (*QuestionnairesHolder, error) {
		questionnaires, err := stores.NewQuestionnaireStore(relationalDB).GetByExperiment(input.ID)
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
	}, func(ctx context.Context, input *QuestionnaireAddInput) (*QuestionnaireHolder, error) {
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

func (o *questionnaireOperations) RegisterUpdate(api huma.API) {
	name := "Update Questionnaire"
	description := "Update an existing questionnaire for which the experiment has not started yet."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{"researcher"}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *QuestionnaireUpdateInput) (*QuestionnaireHolder, error) {
		questionnaire, err := stores.NewQuestionnaireStore(relationalDB).Update(input.credential.UserID, input.ID, input.Body)
		if err != nil {
			return nil, handleError(err)
		}
		if questionnaire == nil {
			return nil, generateNotFoundForThisUserError("questionnaire", input.ID)
		}
		return &QuestionnaireHolder{Body: questionnaire}, nil
	})
}

func (o *questionnaireOperations) RegisterDelete(api huma.API) {
	name := "Delete Questionnaire"
	description := "Delete a questionnaire."
	path := "/" + o.Endpoint + "/{id}"
	scopes := []string{}
	method := http.MethodDelete
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *QuestionnaireDeleteInput) (*struct{}, error) {
		err := stores.NewQuestionnaireStore(relationalDB).Delete(input.ID, input.credential.UserID)
		if err != nil {
			return nil, handleError(err)
		}
		return nil, nil
	})
}

func (o *questionnaireOperations) RegisterGetByExperiment(api huma.API) {
	name := "Get Questionnaires By Experiment"
	description := "Retrieve all questionnaires for a specific experiment."
	path := "/" + o.Endpoint + "s/experiment/{id}"
	scopes := []string{"researcher"}
	method := http.MethodGet
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" format:"uuid" doc:"The ID of the experiment to retrieve questionnaires for."`
	}) (*QuestionnairesHolder, error) {
		questionnaires, err := stores.NewQuestionnaireStore(relationalDB).GetByExperiment(input.ID)
		if err != nil {
			return nil, handleError(err)
		}
		return &QuestionnairesHolder{Body: questionnaires}, nil
	})
}
