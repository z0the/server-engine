package service

import (
	"rpg/internal/server/auth"
	"rpg/pkg/hubber_tmp"

	"github.com/sirupsen/logrus"
)

type MatchMaker interface {
	JoinRoom(ownerID int64, data RoomData) (int64, error)
	SendToGame(roomID int64, req hubber_tmp.IRequest) error
}

type Services struct {
	Auth auth.Service
	MatchMaker
}

func NewService(logger *logrus.Logger) *Services {
	return &Services{
		Auth:       auth.NewAuthService(),
		MatchMaker: NewMatchMakerService(logger),
	}
}
