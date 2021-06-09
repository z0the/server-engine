package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"net"
	"rpg/internal/game"
	"rpg/internal/service"
	"rpg/pkg/hubber"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320
	screenHeight = 240
)

func init() {

}

type Tile struct {
	*ebiten.Image
	id   int64
	x, y float64
}

type Game struct {
	clientID   int64
	conn       net.Conn
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
	req := &hubber.Request{}
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
	//time.Sleep(5 * time.Second)
	resp := &hubber.Response{}
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
		//var data []byte
		//_, err := g.conn.Read(data)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println("data: ", data)
		resp := &hubber.Response{}
		if err := dec.Decode(resp); err != nil {
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			continue
		}
		//fmt.Println("Received in game: ", resp)
		g.handleResponse(resp)
	}
}

func (g *Game) sendToServer(req hubber.IRequest) {
	enc := json.NewEncoder(g.conn)
	if err := enc.Encode(req); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) handleResponse(resp hubber.IResponse) {
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
				g.initChar(char)
				return
			}
		}
	}
}

func (g *Game) initChar(data *game.Player) {

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

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Keyboard (Ebiten Demo)")
	newGame := &Game{
		conn: conn,
	}
	roomData := &service.RoomData{
		GameType:  "default",
		IsPrivate: false,
	}
	req := &hubber.Request{Action: "joinRoom"}
	req.WriteData(roomData)
	newGame.sendToServer(req)
	newGame.Init()
	if err := ebiten.RunGame(newGame); err != nil {
		log.Fatal(err)
	}
}
