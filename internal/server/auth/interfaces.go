package auth

type Service interface {
	RegisterNewUser(login, password string) (*User, error)
	LoginUser(login, password string) (*User, error)
}
