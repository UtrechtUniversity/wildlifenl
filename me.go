package wildlifenl

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
)

type MeInput struct {
	Input
}

type MeHolder struct {
	Body *models.Me `json:"me"`
}

type meOperations struct{}

func (s *meOperations) RegisterGet(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-me",
		Tags:        []string{"me"},
		Method:      http.MethodGet,
		Path:        "/me/",
		Summary:     "Get my profile",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve the logged in user.",
	}, func(ctx context.Context, input *MeInput) (*MeHolder, error) {
		me, err := stores.NewMeStore(database).Get(input.credential.Token)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		if me == nil {
			return nil, huma.Error404NotFound("The logged in user was not found.")
		}
		return &MeHolder{Body: me}, nil
	})
}

func (s *meOperations) RegisterPut(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "put-me",
		Tags:        []string{"me"},
		Method:      http.MethodPut,
		Path:        "/me/",
		Summary:     "Update my profile",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Update the logged in user. Note that if you change your email address you will automatically invalidate all your credentials and need to log in again. If you provide an email address that you do not have access to, you lose access to your account.",
	}, func(ctx context.Context, input *MeHolder) (*MeHolder, error) {
		return nil, huma.Error501NotImplemented("")
	})
}
