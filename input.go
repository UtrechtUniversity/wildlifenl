package wildlifenl

import (
	"strings"

	"github.com/UtrechtUniversity/wildlifenl/models"
	"github.com/danielgtaylor/huma/v2"
)

type Input struct {
	credential *models.Credential
}

func (m *Input) Resolve(ctx huma.Context) []error {
	token := strings.TrimPrefix(ctx.Header("Authorization"), "Bearer ")
	credential, err := getCredential(token)
	if err != nil {
		return []error{err}
	}
	m.credential = credential
	return nil
}
