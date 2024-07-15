package wildlifenl

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/danielgtaylor/huma/v2"
)

type authOperations Operations

func newAuthOperations(database *sql.DB) *authOperations {
	o := authOperations{
		Database: database,
		Endpoint: "auth",
	}
	return &o
}

type AuthenticationInput struct {
	Body struct {
		DisplayNameApp  string `json:"displayNameApp" doc:"The display name of the requesting app, will be used in the email message." example:"MyApp"`
		DisplayNameUser string `json:"displayNameUser" doc:"The display name of the user, will be used in the email message." example:"Jane Smith"`
		Email           string `json:"email" doc:"The email address that the authentication code should be send to." format:"email"`
		//Language        string `json:"language,omitempty" doc:"The two digit code [nl,en] for the language that the email message should be written in. Default:nl" minLength:"2" maxLength:"2" example:"nl"`
	}
}

type AuthenticationResult struct {
	Body string `json:"message"`
}

func (o *authOperations) RegisterAuthentication(api huma.API) {
	name := "Authenticate"
	description := "Start the log on process and request a code by email, then call Authorize with this code."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPost
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *AuthenticationInput) (*AuthenticationResult, error) {
		if err := authenticate(input.Body.DisplayNameApp, input.Body.DisplayNameUser, input.Body.Email); err != nil {
			log.Println("ERROR authentication:", err)
			return nil, huma.Error500InternalServerError("An email could not be sent to the provided email address.")
		}

		// TODO add mailing code, then remove code below.
		code := ""
		if value, ok := authRequests.Get(input.Body.Email); ok {
			code = value.(AuthenticationRequest).code
		}
		return &AuthenticationResult{Body: "The authentication code: " + code}, nil
		// ----

		// return &AuthenticationResult{Body: "The authentication code has been sent to: " + input.Body.Email}, nil
	})
}

type AuthorizationInput struct {
	Body struct {
		Email string `json:"email" doc:"The email address that the authentication process was started with." format:"email"`
		Code  string `json:"code" doc:"The code as received by email." example:"123456"`
	}
}

type AuthorizationResult struct {
	Body *models.Credential `json:"credential"`
}

func (o *authOperations) RegisterAuthorisation(api huma.API) {
	name := "Authorize"
	description := "Finalize the log on process by providing the code as received by email and get a bearer token."
	path := "/" + o.Endpoint + "/"
	scopes := []string{}
	method := http.MethodPut
	huma.Register(api, huma.Operation{
		OperationID: name, Summary: name, Path: path, Method: method, Tags: []string{o.Endpoint}, Description: generateDescription(description, scopes), Security: []map[string][]string{{"auth": scopes}},
	}, func(ctx context.Context, input *AuthorizationInput) (*AuthorizationResult, error) {
		credential, err := authorize(input.Body.Email, input.Body.Code)
		if err != nil {
			return nil, handleError(err)
		}
		if credential == nil {
			return nil, huma.Error403Forbidden("The combination of email and code does not match a previous authentication")
		}
		return &AuthorizationResult{Body: credential}, nil
	})
}
