package auth

import (
	"sync"

	"github.com/google/uuid"
)

//nolint:gochecknoglobals

func NewAuthService() Service {
	return &authService{
		loginData: make(map[string]*loginUser),
		takenUIDs: make(map[string]struct{}),
	}
}

type authService struct {
	sync.Mutex
	loginData map[string]*loginUser
	takenUIDs map[string]struct{}
}

func (s *authService) RegisterNewUser(login, password string) (*User, error) {
	s.Lock()
	defer s.Unlock()

	_, alreadyRegistered := s.loginData[login]
	if alreadyRegistered {
		return nil, ErrUserWithSameLoginAlreadyRegistered
	}

	lu := &loginUser{
		UID:      s.getNewUID(),
		Login:    login,
		Password: password,
		Nick:     login,
		Coins:    0,
	}

	s.loginData[lu.Login] = lu
	s.takenUIDs[lu.UID] = struct{}{}

	return lu.User(), nil
}

func (s *authService) LoginUser(login, password string) (*User, error) {
	lu, userExists := s.loginData[login]
	if !userExists {
		return nil, ErrUserDoesNotExist
	}
	if lu.Password != password {
		return nil, ErrWrongPassword
	}
	return lu.User(), nil
}

// getNewUID returns first free uid for loginData
func (s *authService) getNewUID() string {
	newUID := uuid.NewString()
	_, isAlreadyExists := s.takenUIDs[newUID]
	if isAlreadyExists {
		return s.getNewUID()
	}
	return newUID
}
