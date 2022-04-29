package game

import (
	"math/rand"
	"time"

	"go.uber.org/zap"
)

func NewGame(logger *zap.SugaredLogger, characters []*Player, msgPipe chan<- hubber_tmp.IResponse) *Game {
	return &Game{
		logger:     logger,
		characters: characters,
		events:     make(chan hubber_tmp.IRequest, 10),
		msgPipe:    msgPipe,
	}
}

type Game struct {
	logger     *zap.SugaredLogger
	characters []*Player
	events     chan hubber_tmp.IRequest
	msgPipe    chan<- hubber_tmp.IResponse
}

func (g *Game) HandleRequest(event hubber_tmp.IRequest) {
	g.events <- event
}

func (g *Game) Start() {
	// g.init()
	go g.updater()
	go g.mainLoop()
	go g.observer()
}

func (g *Game) init() {
	rand.Seed(time.Now().Unix())
	for _, player := range g.characters {
		player.PosY = float64(rand.Intn(40) + 10)
		player.PosX = float64(rand.Intn(40) + 10)
		player.ColorRGBA = [4]uint8{
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			255,
		}
	}
}

func (g *Game) mainLoop() {
	for event := range g.events {
		switch event.GetAction() {
		case "gameMove":
			g.handleMove(event)
		}
	}
}

const (
	Left Direction = iota
	Right
	Top
	Bottom
)

type Direction int

func (s *Direction) GetSide() string {
	return []string{"Left", "Right", "Top", "Bottom"}[*s]
}

type MoveData struct {
	Direction `json:"side"`
}

func (g *Game) handleMove(event hubber_tmp.IRequest) {
	var character *Player
	for _, checkCharacter := range g.characters {
		if checkCharacter.OwnerID == event.SenderID() {
			character = checkCharacter
		}
	}
	if character == nil {
		g.logger.Warn("no such character in game")
		return
	}
	data := &MoveData{}
	event.ParseData(data)
	switch data.Direction {
	case Left:
		character.PosX -= 2
	case Right:
		character.PosX += 2
	case Top:
		character.PosY -= 2
	case Bottom:
		character.PosY += 2
	default:
		g.logger.Warn("wrong move side")
	}
}

func (g *Game) InitPlayer(player *Player) {
	g.characters = append(g.characters, player)
}

func (g *Game) observer() {
	for {
		g.logger.Info(*g)
		time.Sleep(time.Second)
	}
}

func (g *Game) updater() {
	ticker := time.NewTicker(16 * time.Millisecond)
	for range ticker.C {
		g.sendGameState()
	}
}

type Event struct {
	Type string
	Data any
}

type StateData struct {
	Characters []*Player `json:"characters"`
}

func (g *Game) sendGameState() {
	for _, character := range g.characters {
		data := &StateData{
			Characters: g.characters,
		}
		resp := &hubber_tmp.Response{}
		resp.SetReceiverID(character.OwnerID)
		resp.Action = "gameState"
		resp.WriteData(data)
		g.msgPipe <- resp
	}
}

// PlayerData - is used for send to client
type PlayerData struct {
	OwnerID   int64
	ColorRGBA [4]uint8
	Score     int     `json:"score"`
	PosX      float64 `json:"posX"`
	PosY      float64 `json:"posY"`
}

type Player struct {
	isLoaded bool
	PlayerData
}

func NewPlayer(ownerID int64) *Player {
	rand.Seed(time.Now().Unix())
	player := &Player{}
	player.OwnerID = ownerID
	player.PosY = float64(rand.Intn(10)+10) * (-1)
	player.PosX = float64(rand.Intn(10) + 10)
	player.ColorRGBA = [4]uint8{
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		uint8(rand.Intn(256)),
		255,
	}
	return player
}

func (p *Player) IsLoaded() bool {
	return p.isLoaded
}

func (p *Player) SetLoaded() {
	p.isLoaded = true
}

func (p *Player) SetNotLoaded() {
	p.isLoaded = false
}
