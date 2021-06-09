package service

import (
	"errors"
	"rpg/pkg/hubber"
)

type User struct {
	//OwnerID Dynamic assignment for every session
	OwnerID  int64
	Name     string
	Password string
	Coins    int
}

var bdUsers = []User{
	{Name: "Pavel", Password: "123", Coins: 100},
	{Name: "Vova", Password: "321", Coins: 1000},
}

type RoomService struct {
	lastID int64
}

func NewRoomService(msgPipe chan<- hubber.IResponse) *RoomService {
	return &RoomService{}
}

func (s *RoomService) Login(user User) (*User, error) {
	for _, bdUser := range bdUsers {
		if bdUser.Name == user.Name && bdUser.Password == user.Password {
			s.lastID++
			user = bdUser
			user.OwnerID = s.lastID
			return &user, nil
		}
	}
	return nil, errors.New("no such registered user")
}
