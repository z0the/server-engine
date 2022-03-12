package auth

import "errors"

var (
	ErrUserDoesNotExist                   = errors.New("user does not exist")
	ErrWrongPassword                      = errors.New("wrong password")
	ErrUserWithSameLoginAlreadyRegistered = errors.New("user with same login already registered")
)
