package matchmaker

import (
	"errors"
	"sync"

	"rpg/internal/server/game"
	"rpg/pkg/hubber_tmp"

	"github.com/sirupsen/logrus"
)

const (
	defaultType = "DEFAULT"
)

type Service struct {
	sync.Mutex
	logger    *logrus.Logger
	msgPipe   chan<- hubber_tmp.IResponse
	roomCount int
	lastID    int64
	rooms     map[int64]*Room
	// freeRooms   map[int64]*Room
}

func NewMatchMakerService(logger *logrus.Logger) *Service {
	return &Service{
		logger: logger,
		// msgPipe: msgPipe,
		rooms: make(map[int64]*Room),
		// freeRooms:   make(map[int64]*Room),
	}
}

func (s *Service) JoinRoom(ownerID int64, data RoomData) (int64, error) {
	s.Lock()
	defer s.Unlock()

	newPlayer := game.NewPlayer(ownerID)

	data = s.validateRoomData(data)
	foundRoom := s.findAppropriateRoom(data)
	if foundRoom != nil {
		if err := foundRoom.addPlayer(newPlayer); err != nil {
			return 0, err
		}
		// if foundRoom.isFull() {
		//	delete(s.freeRooms, foundRoom.id)
		//	s.rooms[foundRoom.id] = foundRoom
		// }
		return foundRoom.id, nil
	}
	newRoom := s.createRoom(data)
	if err := newRoom.addPlayer(newPlayer); err != nil {
		return 0, err
	}
	return newRoom.id, nil
}

func (s *Service) SendToGame(roomID int64, req hubber_tmp.IRequest) error {
	if room, ok := s.rooms[roomID]; ok {
		if err := room.handleRequest(req); err != nil {
			return err
		}
		return nil
	}
	return ErrNoSuchRoom
}

func (s *Service) validateRoomData(data RoomData) RoomData {
	if data.GameType != defaultType {
		data.GameType = defaultType
	}
	if data.MaxPlayers < 2 {
		data.MaxPlayers = 2
	}
	if data.MaxPlayers > 10 {
		data.MaxPlayers = 10
	}
	return data
}

func (s *Service) createRoom(data RoomData) *Room {
	s.lastID++
	newRoom := &Room{
		msgPipe:  s.msgPipe,
		logger:   s.logger,
		id:       s.lastID,
		RoomData: data,
	}
	s.rooms[s.lastID] = newRoom
	return newRoom
}

func (s *Service) findAppropriateRoom(data RoomData) *Room {
	for _, room := range s.rooms {
		if room.isAppropriate(data) {
			return room
		}
	}
	return nil
}

type RoomData struct {
	MaxPlayers int
	GameType   string
	IsPrivate  bool
}

type Room struct {
	sync.Mutex
	RoomData
	*game.Game
	logger   *logrus.Logger
	msgPipe  chan<- hubber_tmp.IResponse
	isActive bool
	id       int64
	players  []*game.Player
}

func (r *Room) addPlayer(player *game.Player) error {
	r.Lock()
	defer r.Unlock()
	for _, checkPlayer := range r.players {
		if checkPlayer.OwnerID == player.OwnerID {
			return ErrPlayerIsAlreadyInRoom
		}
	}
	r.players = append(r.players, player)
	if r.Game != nil {
		r.Game.InitPlayer(player)
	}
	r.startGame()
	return nil
}

func (r *Room) handleRequest(req hubber_tmp.IRequest) error {
	if r.Game != nil {
		r.Game.HandleRequest(req)
		return nil
	}
	return ErrRoomIsNotRunning
}

func (r *Room) startGame() {
	if r.Game != nil {
		return
	}
	r.Game = game.NewGame(r.logger, r.players, r.msgPipe)
	r.Start()
}

func (r *Room) removePlayer(ownerID int64) error {
	r.Lock()
	defer r.Unlock()
	for index, player := range r.players {
		if player.OwnerID == ownerID {
			r.players = append(r.players[:index], r.players[index+1:]...)
			return nil
		}
	}
	return ErrNoSuchPlayerInRoom
}

func (r *Room) isAppropriate(data RoomData) bool {
	return r.GameType == data.GameType &&
		r.MaxPlayers == data.MaxPlayers &&
		r.IsPrivate == data.IsPrivate
}

func (r *Room) isFull() bool {
	return len(r.players) >= r.MaxPlayers
}

var (
	ErrPlayerIsAlreadyInRoom = errors.New("player is already in this room")
	ErrNoSuchPlayerInRoom    = errors.New("no such player in this room")
	ErrNoSuchRoom            = errors.New("no such room")
	ErrRoomIsNotRunning      = errors.New("room is not running")
)
