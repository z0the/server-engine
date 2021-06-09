package service

import (
	"rpg/pkg/hubber"

	"github.com/sirupsen/logrus"
)

type Auth interface {
	Login(user User) (*User, error)
}

type MatchMaker interface {
	JoinRoom(ownerID int64, data RoomData) (int64, error)
	SendToGame(roomID int64, req hubber.IRequest) error
}

type Service struct {
	Auth
	MatchMaker
}

func NewService(logger *logrus.Logger, msgPipe chan<- hubber.IResponse) *Service {
	return &Service{
		Auth:       NewRoomService(msgPipe),
		MatchMaker: NewMatchMakerService(logger, msgPipe),
	}
}
