package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"net"

	"rpg/internal/server/api"
	"rpg/internal/server/game"
	"rpg/pkg/hubber_tmp"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

type Tile struct {
	*ebiten.Image
	id   int64
	x, y float64
}

func NewGame(connCtrl *connectionController) *Game {
	g := &Game{
		connCtrl: connCtrl,
	}
	// g.connCtrl.
}

type Game struct {
	connCtrl   *connectionController
	keys       []ebiten.Key
	characters []*Tile
}

func (g *Game) Init() {
	go g.listenServer()
}

func (g *Game) Update() error {
	g.keys = inpututil.PressedKeys()
	for _, char := range g.characters {
		if char.id == g.clientID {
			for _, key := range g.keys {
				if key == ebiten.KeyD {
					char.x += 2
					g.sendMove(game.Right)
				}
				if key == ebiten.KeyA {
					char.x -= 2
					g.sendMove(game.Left)
				}
				if key == ebiten.KeyW {
					char.y -= 2
					g.sendMove(game.Top)
				}
				if key == ebiten.KeyS {
					char.y += 2
					g.sendMove(game.Bottom)
				}
			}
		}
	}
	return nil
}

func (g *Game) sendMove(direction game.Direction) {
	req := &hubber_tmp.Request{}
	req.Action = "gameMove"
	data := &game.MoveData{
		Direction: direction,
	}
	req.WriteData(data)
	fmt.Println("Handle move: ", req)
	g.sendToServer(req)
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, tile := range g.characters {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(tile.x, tile.y)
		screen.DrawImage(tile.Image, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) listenServer() {
	resp := &hubber_tmp.Response{}
	dec := json.NewDecoder(g.conn)
	if err := dec.Decode(resp); err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
	}
	ID := &struct {
		ID int64
	}{}
	err := json.Unmarshal(resp.Data, ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Got clientID: ", ID.ID)
	g.clientID = ID.ID
	for {
		resp := &hubber_tmp.Response{}
		if err := dec.Decode(resp); err != nil {
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			continue
		}
		g.handleResponse(resp)
	}
}

func (g *Game) sendToServer(req hubber_tmp.IRequest) {
	enc := json.NewEncoder(g.conn)
	if err := enc.Encode(req); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) handleResponse(resp hubber_tmp.IResponse) {
	switch resp.GetAction() {
	case "gameState":
		data := &game.StateData{}
		resp.ParseData(data)
		for _, char := range data.Characters {
			found := false
			for _, oldChar := range g.characters {
				if char.OwnerID == oldChar.id {
					oldChar.x = char.PosX
					oldChar.y = char.PosY
					found = true
					break
				}
			}
			if !found {
				g.initCharacter(char)
				return
			}
		}
	}
}

func (g *Game) initCharacter(data *game.Player) {
	img := ebiten.NewImage(20, 20)
	img.Fill(color.RGBA{R: data.ColorRGBA[0], G: data.ColorRGBA[1], B: data.ColorRGBA[2], A: data.ColorRGBA[3]})
	newPlayer := &Tile{
		id:    data.OwnerID,
		Image: img,
		x:     data.PosX,
		y:     data.PosY,
	}
	g.characters = append(g.characters, newPlayer)
}

const address = "localhost:3000"

func main() {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	go listenConn(conn)

	enc := json.NewEncoder(conn)

	req, err := api.NewBaseRequest(
		api.RegistrationReqType,
		api.RegistrationPayload{
			Login:    "zothe",
			Password: "123456",
		},
	)
	if err != nil {
		panic(err)
	}

	// time.Sleep(time.Second * 4)
	err = enc.Encode(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("Send request")

	// ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	// ebiten.SetWindowTitle("Simple Game")
	// newGame := &Game{
	// 	conn: conn,
	// }

	// newGame.sendToServer(req)
	// newGame.Init()
	// if err := ebiten.RunGame(newGame); err != nil {
	// 	log.Fatal(err)
	// }
	<-context.Background().Done()
}
