package auth

import "github.com/dgrijalva/jwt-go"

type User struct {
	UID   string
	Nick  string
	Coins uint
}

type loginUser struct {
	UID      string
	Login    string
	Password string
	Nick     string
	Coins    uint
}

func (u *loginUser) User() *User {
	return &User{
		UID:   u.UID,
		Nick:  u.Nick,
		Coins: u.Coins,
	}
}

type DefaultClaims struct {
	jwt.StandardClaims
	UserUID string `json:"user_uid"`
	Login   string `json:"login"`
}
