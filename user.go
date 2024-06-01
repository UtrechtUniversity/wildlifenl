package wildlifenl

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/UtrechtUniversity/wildlifenl/stores"
)

type UserHolder struct {
	Body *models.User `json:"user"`
}

type UsersHolder struct {
	Body []models.User `json:"users"`
}

type userOperations struct{}

func (s *userOperations) RegisterGetUserByID(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-user-by-id",
		Tags:        []string{"user"},
		Method:      http.MethodGet,
		Path:        "/user/{id}",
		Summary:     "Get User By ID",
		Security:    []map[string][]string{{"auth": []string{}}},
		Description: "Retrieve a specific user by ID.",
	}, func(ctx context.Context, input *struct {
		ID string `path:"id" doc:"The ID of the user." format:"uuid"`
	}) (*UserHolder, error) {
		user, err := stores.NewUserStore(database).Get(input.ID)
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		if user == nil {
			return nil, huma.Error404NotFound("No user with ID " + input.ID + " was found")
		}
		return &UserHolder{Body: user}, nil
	})
}

func (s *userOperations) RegisterGetAllUsers(api huma.API) {
	scopes := []string{"researcher"}
	huma.Register(api, huma.Operation{
		OperationID: "get-all-users",
		Tags:        []string{"user"},
		Method:      http.MethodGet,
		Path:        "/user/",
		Summary:     "Get all Users",
		Security:    []map[string][]string{{"auth": scopes}},
		Description: "Retrieve all users. <br/><br/>**Scopes**<br/>" + scopesAsMarkdown(scopes),
	}, func(ctx context.Context, input *struct{}) (*UsersHolder, error) {
		users, err := stores.NewUserStore(database).GetAll()
		if err != nil {
			log.Println("ERROR", logTrace(), err)
			return nil, huma.Error500InternalServerError("")
		}
		return &UsersHolder{Body: users}, nil
	})
}
