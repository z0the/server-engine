package matchmaker

import (
	"fmt"
	"math/rand"
	"sync"

	"go.uber.org/zap"

	"rpg/internal/server/async_controller/asyncapi"
	"rpg/internal/server/eventbus"
)

func NewMatchMakerService(lg *zap.SugaredLogger, bus eventbus.Bus) Service {
	svc := &service{
		lg:  lg,
		bus: bus,
		roomList: map[string]*Room{
			"1": newRoom(bus, lg),
		},
	}
	return svc
}

type service struct {
	sync.RWMutex
	lg       *zap.SugaredLogger
	bus      eventbus.Bus
	roomList map[string]*Room
}

func (s *service) JoinRoom(connUID, clientUID, login string) string {
	s.Lock()
	defer s.Unlock()
	player := &Player{
		connUID: connUID,
		UID:     clientUID,
		Login:   login,
		PosX:    rand.Intn(100),
		PosY:    rand.Intn(100),
	}
	room := s.roomList["1"]
	room.addPlayer(player)

	room.sendUpdateSceneMsg()

	return "1"
}

func newRoom(bus eventbus.Bus, lg *zap.SugaredLogger) *Room {
	return &Room{
		lg:         lg,
		bus:        bus,
		UID:        "1",
		playerList: map[string]*Player{},
	}
}

type Room struct {
	sync.RWMutex
	lg         *zap.SugaredLogger
	bus        eventbus.Bus
	UID        string
	playerList map[string]*Player
}

func (r *Room) addPlayer(player *Player) {
	r.Lock()
	defer r.Unlock()
	r.playerList[player.connUID] = player
	r.bus.SubscribeWithFromPrefix(player.connUID, r.msgHandler)
}

func (r *Room) msgHandler(msg eventbus.BusMsg) {
	msgType := asyncapi.MsgType(msg.MsgType)
	switch msgType {
	case asyncapi.MoveMsgType:
		payload, ok := msg.Payload.(asyncapi.MovePayloadIN)
		if !ok {
			r.lg.Errorw("wrong payload type", "payload", payload)
			return
		}
		fmt.Printf("Move speed:%d direction:%s\n", payload.Speed, payload.Direction)
	default:
		r.lg.Errorw("wrong msg type", "msg_type", msgType)
	}
}

func (r *Room) sendUpdateSceneMsg() {
	r.RLock()
	defer r.RUnlock()

	dtoPlayers := make([]asyncapi.Player, 0, len(r.playerList))
	for _, player := range r.playerList {
		dtoPlayers = append(dtoPlayers, player.makeDTO())
	}
	for _, player := range r.playerList {
		go r.bus.PublishWithToPrefix(
			player.topicKey(), eventbus.BusMsg{
				RecipientConnUID: player.connUID,
				MsgType:          asyncapi.SceneUpdateMsgType.String(),
				Payload: asyncapi.SceneUpdatePayloadOUT{
					PlayersList: dtoPlayers,
				},
			},
		)
	}
}

type Player struct {
	connUID string
	UID     string
	Login   string
	PosX    int
	PosY    int
	Score   uint
}

func (p *Player) makeDTO() asyncapi.Player {
	return asyncapi.Player{
		UID:   p.UID,
		Login: p.Login,
		PosX:  p.PosX,
		PosY:  p.PosY,
	}
}

func (p *Player) topicKey() string {
	return p.connUID
}
