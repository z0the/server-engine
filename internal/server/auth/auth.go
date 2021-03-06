package auth

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	"go.uber.org/zap"

	"rpg/internal/server/utils"
)

func NewAuthService(lg *zap.SugaredLogger) Service {
	return &service{
		lg:        lg,
		loginData: make(map[string]*loginUser),
		takenUIDs: make(map[string]struct{}),
	}
}

type service struct {
	sync.Mutex
	lg        *zap.SugaredLogger
	loginData map[string]*loginUser
	takenUIDs map[string]struct{}
}

const jwtSecret = "123456"

func (s *service) ParseClaims(bearerToken string) (DefaultClaims, error) {
	jwtToken := strings.TrimPrefix(bearerToken, "Bearer ")
	parsedClaims := DefaultClaims{}
	_, err := jwt.ParseWithClaims(
		jwtToken, &parsedClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		},
	)
	if err != nil {
		return DefaultClaims{}, err
	}
	fmt.Println(parsedClaims.Login)

	return parsedClaims, nil
}

func (s *service) RegisterNewUser(login, password string) (string, error) {
	s.Lock()
	defer s.Unlock()

	_, alreadyRegistered := s.loginData[login]
	if alreadyRegistered {
		return utils.EmptyString, ErrUserWithSameLoginAlreadyRegistered
	}

	lu := loginUser{
		UID:      s.getNewUID(),
		Login:    login,
		Password: password,
		Nick:     login,
		Coins:    0,
	}

	s.loginData[lu.Login] = &lu
	s.takenUIDs[lu.UID] = struct{}{}

	return s.makeJWTToken(lu)
}

func (s *service) LoginUser(login, password string) (string, error) {
	lu, userExists := s.loginData[login]
	if !userExists {
		return utils.EmptyString, ErrUserDoesNotExist
	}
	if lu.Password != password {
		return utils.EmptyString, ErrWrongPassword
	}

	return s.makeJWTToken(*lu)
}

// getNewUID returns first free uid for loginData
func (s *service) getNewUID() string {
	newUID := uuid.NewString()
	_, isAlreadyExists := s.takenUIDs[newUID]
	if isAlreadyExists {
		return s.getNewUID()
	}
	return newUID
}

func (s *service) makeJWTToken(lu loginUser) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		DefaultClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			Login:   lu.Login,
			UserUID: lu.UID,
		},
	)

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
