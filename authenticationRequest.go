package wildlifenl

type AuthenticationRequest struct {
	appName  string
	userName string
	email    string
	code     string
}
