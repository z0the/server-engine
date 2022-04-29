package matchmaker

type Service interface {
	JoinRoom(connUID, clientUID, login string) (roomUID string)
}

// type Service interface {
// 	JoinRoom(ownerID int64, data RoomData) (int64, error)
// 	SendToGame(roomID int64, req hubber_tmp.IRequest) error
// }
//
// type Services struct {
// 	Auth auth.service
// 	Service
// }
//
// func NewMatchMakerService(logger *zap.SugaredLogger) *Services {
// 	return &Services{
// 		Auth:       auth.NewAuthService(logger),
// 		Service: NewMatchMakerService(logger),
// 	}
// }
