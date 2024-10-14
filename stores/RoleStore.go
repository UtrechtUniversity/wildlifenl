package stores

import (
	"database/sql"

	"github.com/UtrechtUniversity/wildlifenl/models"
)

type RoleStore Store

func NewRoleStore(db *sql.DB) *RoleStore {
	s := RoleStore{
		relationalDB: db,
		query: `
		SELECT r."ID", r."name"
		FROM "role" r
		`,
	}
	return &s
}

func (s *RoleStore) process(rows *sql.Rows, err error) ([]models.Role, error) {
	if err != nil {
		return nil, err
	}
	roles := make([]models.Role, 0)
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (s *RoleStore) GetAll() ([]models.Role, error) {
	query := s.query + `
		ORDER BY r."ID"
		`
	rows, err := s.relationalDB.Query(query)
	return s.process(rows, err)
}

func (s *RoleStore) AddRoleToUser(userID string, roleID int) error {
	query := `
		INSERT INTO "user_role"("userID", "roleID") VALUES($1, $2)
	`
	_, err := s.relationalDB.Exec(query, userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func (s *RoleStore) RemoveRoleFromUser(userID string, roleID int) error {
	query := `
		DELETE FROM "user_role"
		WHERE "userID" = $1
		AND "roleID" = $2
	`
	_, err := s.relationalDB.Exec(query, userID, roleID)
	if err != nil {
		return err
	}
	return nil
}
