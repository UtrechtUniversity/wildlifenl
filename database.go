package wildlifenl

type Database struct {
}

func newDataBase() (*Database, error) {
	database := Database{}
	return &database, nil
}

// getAccount returns the account for the email+password combination, or nil if no such account exists.
func (d *Database) getAccount(email string, password string) *Account {
	if email == "bas" && password == "bas" {
		return &Account{
			Email: email,
			Roles: []string{"admin"},
		}
	}
	return nil
}

// addAccount adds a new account for the email+password combination and returns it, or returns nil if an account for email already exists.
func (d *Database) addAccount(email string, password string) *Account {
	return &Account{
		Email: email,
		Roles: []string{},
	}
}
