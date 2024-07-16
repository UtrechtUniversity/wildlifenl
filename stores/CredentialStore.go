package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type CredentialStore Store

func NewCredentialStore(db *sql.DB) *CredentialStore {
	s := CredentialStore{
		relationalDB: db,
		query: `
		SELECT u."id", c."token", c."email", c."lastLogon", r."name"
		FROM "credential" c
		INNER JOIN "user" u ON u."email" = C."email"
		LEFT JOIN "user_role" x ON x."userID" = u."id"
		LEFT JOIN "role" r ON r."id" = x."roleID"
		`,
	}
	return &s
}

func (s *CredentialStore) Get(token string) (*models.Credential, error) {
	query := s.query + `
		WHERE c."token" = $1
	`
	rows, err := s.relationalDB.Query(query, token)
	if err != nil {
		return nil, err
	}
	credential := models.Credential{
		Scopes: make([]string, 0),
	}
	for rows.Next() {
		var role string
		rows.Scan(&credential.UserID, &credential.Token, &credential.Email, &credential.LastLogon, &role)
		credential.Scopes = append(credential.Scopes, role)
	}
	return &credential, nil
}

func (s *CredentialStore) Create(appName, userName, email string) (*models.Credential, error) {
	query := `
		INSERT INTO "user"("name", "email") VALUES($1, $2) 
		ON CONFLICT("email") DO UPDATE SET "name" = $3 
		RETURNING "id"
	`
	row := s.relationalDB.QueryRow(query, userName, email, userName)
	var userID string
	if err := row.Scan(&userID); err != nil {
		return nil, err
	}
	query = `
		SELECT r."name"
		FROM "user_role" x
		INNER JOIN "role" r ON r."id" = x."roleID"
		INNER JOIN "user" u ON u."id" = x."userID"
		WHERE u."id" = $1;
	`
	rows, err := s.relationalDB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	scopes := make([]string, 0)
	for rows.Next() {
		var role string
		rows.Scan(&role)
		scopes = append(scopes, role)
	}
	query = `
		INSERT INTO credential("email", "appName") VALUES($1, $2)
		RETURNING "token", "email", "lastLogon";
	`
	row = s.relationalDB.QueryRow(query, email, appName)
	credential := models.Credential{UserID: userID}
	if err := row.Scan(&credential.Token, &credential.Email, &credential.LastLogon); err != nil {
		return nil, err
	}
	credential.Scopes = scopes
	return &credential, nil
}
