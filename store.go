package wildlifenl

import "github.com/google/uuid"

type Store struct {
}

func newStore() (*Store, error) {
	store := Store{}
	return &store, nil
}

func (d *Store) getAccount(token string) *Account {
	scopes := make([]string, 0)
	if token == "bas" {
		scopes = append(scopes, "researcher")
	}
	return &Account{
		Token:  token,
		Scopes: scopes,
	}
}

func (d *Store) createAccount() *Account {
	token := uuid.New().String()
	return &Account{
		Token:  token,
		Scopes: []string{},
	}
}

func (s *Store) getAnimal(id int) *Animal {
	if id == 42 {
		return &Animal{
			ID:      42,
			Name:    "Flupke",
			Species: Species{Name: "Canis familiaris", CommonName: "Dog"},
		}
	}
	return nil
}
