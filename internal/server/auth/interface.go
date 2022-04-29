package auth

type Service interface {
	ParseClaims(bearerToken string) (DefaultClaims, error)
	RegisterNewUser(login, password string) (string, error)
	LoginUser(login, password string) (string, error)
}
