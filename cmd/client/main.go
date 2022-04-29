package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "image/png"
	"log"
	"net"
	"net/http"

	"go.uber.org/zap"

	"rpg/internal/server/async_controller/asyncapi"
	"rpg/internal/server/sync_controller/syncapi"
	"rpg/pkg/hubber"
)

// const (
// 	screenWidth  = 320
// 	screenHeight = 240
// )
//
// type Tile struct {
// 	*ebiten.Image
// 	id   int64
// 	x, y float64
// }
//
// func NewGame(connCtrl *connectionController) *Game {
// 	g := &Game{
// 		connCtrl: connCtrl,
// 	}
// 	// g.connCtrl.
// }
//
// type Game struct {
// 	connCtrl   *connectionController
// 	keys       []ebiten.Key
// 	characters []*Tile
// }
//
// func (g *Game) Init() {
// 	go g.listenServer()
// }
//
// func (g *Game) Update() error {
// 	g.keys = inpututil.PressedKeys()
// 	for _, char := range g.characters {
// 		if char.id == g.clientID {
// 			for _, key := range g.keys {
// 				if key == ebiten.KeyD {
// 					char.x += 2
// 					g.sendMove(game.Right)
// 				}
// 				if key == ebiten.KeyA {
// 					char.x -= 2
// 					g.sendMove(game.Left)
// 				}
// 				if key == ebiten.KeyW {
// 					char.y -= 2
// 					g.sendMove(game.Top)
// 				}
// 				if key == ebiten.KeyS {
// 					char.y += 2
// 					g.sendMove(game.Bottom)
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }
//
// func (g *Game) sendMove(direction game.Direction) {
// 	req := &hubber_tmp.Request{}
// 	req.Action = "gameMove"
// 	data := &game.MoveData{
// 		Direction: direction,
// 	}
// 	req.WriteData(data)
// 	fmt.Println("Handle move: ", req)
// 	g.sendToServer(req)
// }
//
// func (g *Game) Draw(screen *ebiten.Image) {
// 	for _, tile := range g.characters {
// 		op := &ebiten.DrawImageOptions{}
// 		op.GeoM.Translate(tile.x, tile.y)
// 		screen.DrawImage(tile.Image, op)
// 	}
// }
//
// func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
// 	return outsideWidth, outsideHeight
// }
//
// func (g *Game) listenServer() {
// 	resp := &hubber_tmp.Response{}
// 	dec := json.NewDecoder(g.conn)
// 	if err := dec.Decode(resp); err != nil {
// 		if e, ok := err.(*json.SyntaxError); ok {
// 			log.Printf("syntax error at byte offset %d", e.Offset)
// 		}
// 	}
// 	ID := &struct {
// 		ID int64
// 	}{}
// 	err := json.Unmarshal(resp.Data, ID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Got clientID: ", ID.ID)
// 	g.clientID = ID.ID
// 	for {
// 		resp := &hubber_tmp.Response{}
// 		if err := dec.Decode(resp); err != nil {
// 			if e, ok := err.(*json.SyntaxError); ok {
// 				log.Printf("syntax error at byte offset %d", e.Offset)
// 			}
// 			continue
// 		}
// 		g.handleResponse(resp)
// 	}
// }
//
// func (g *Game) sendToServer(req hubber_tmp.IRequest) {
// 	enc := json.NewEncoder(g.conn)
// 	if err := enc.Encode(req); err != nil {
// 		log.Fatal(err)
// 	}
// }
//
// func (g *Game) handleResponse(resp hubber_tmp.IResponse) {
// 	switch resp.GetAction() {
// 	case "gameState":
// 		data := &game.StateData{}
// 		resp.ParseData(data)
// 		for _, char := range data.Characters {
// 			found := false
// 			for _, oldChar := range g.characters {
// 				if char.OwnerID == oldChar.id {
// 					oldChar.x = char.PosX
// 					oldChar.y = char.PosY
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				g.initCharacter(char)
// 				return
// 			}
// 		}
// 	}
// }
//
// func (g *Game) initCharacter(data *game.Player) {
// 	img := ebiten.NewImage(20, 20)
// 	img.Fill(color.RGBA{R: data.ColorRGBA[0], G: data.ColorRGBA[1], B: data.ColorRGBA[2], A: data.ColorRGBA[3]})
// 	newPlayer := &Tile{
// 		id:    data.OwnerID,
// 		Image: img,
// 		x:     data.PosX,
// 		y:     data.PosY,
// 	}
// 	g.characters = append(g.characters, newPlayer)
// }

const address = "localhost:3000"

func main() {
	w := serveConn()

	fmt.Println("waiting conn uid")
	var connUID string
	_, err := fmt.Scan(&connUID)
	if err != nil {
		panic(err)
	}

	rawData, err := json.Marshal(
		&syncapi.RegistrationPayloadIN{
			Login:    "test1",
			Password: "123456",
		},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("registration")
	resp, err := http.Post("http://localhost:5000/registration", "application/json", bytes.NewReader(rawData))
	if err != nil {
		panic(err)
	}
	baseResp := syncapi.BaseResponse{}
	err = json.NewDecoder(resp.Body).Decode(&baseResp)
	if err != nil {
		panic(err)
	}
	var data syncapi.RegistrationPayloadOUT
	mapToStruct(baseResp.Data.(map[string]any), &data)

	rawData, err = json.Marshal(
		&syncapi.JoinRoomPayloadIN{
			ConnectionUID: connUID,
		},
	)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:5000/joinRoom", bytes.NewReader(rawData))
	if err != nil {
		panic(err)
	}

	fmt.Println("join room")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", data.Token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	moveMsg, err := asyncapi.NewBaseMessage(
		asyncapi.MoveMsgType, asyncapi.MovePayloadIN{
			Direction: "right",
			Speed:     2,
		},
	)
	if err != nil {
		panic(err)
	}

	rawMsg, err := json.Marshal(moveMsg)
	if err != nil {
		panic(err)
	}
	fmt.Println("send move")

	w.Send(rawMsg)

	fmt.Println("waiting stop")
	fmt.Scan(&connUID)
}

func serveConn() hubber.ConnectionWrapper {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	sugar := logger.Sugar()

	msgPipe := make(chan hubber.Message)
	go func() {
		for msg := range msgPipe {
			var payload asyncapi.BaseMessage
			msg.GetRawData()
			err := json.Unmarshal(msg.GetRawData(), &payload)
			if err != nil {
				panic(err)
			}
			fmt.Printf("MsgType:%s\nPayload:%s\n", payload.MsgType, string(payload.Payload))
		}
	}()

	// go listenConn(conn, msgPipe)
	wrapper := hubber.NewConnectionWrapper(conn, sugar, hubber.DefaultReadConnDelimiter)
	wrapper.StartReading("", msgPipe)
	return wrapper
}

const defaultReadConnDelimiter byte = '\n'

func listenConn(conn net.Conn, msgPipe chan<- asyncapi.BaseMessage) {
	defer func() {
		fmt.Println("Stop listen pump...")
		err := conn.Close()
		if err != nil {
			fmt.Println("Error: ", err)
		}
	}()

	// bufReader := bufio.NewReader(conn)
	dec := json.NewDecoder(conn)
	for {
		fmt.Println("listening...")
		// rawMsg, err := bufReader.ReadBytes(defaultReadConnDelimiter)
		// if err != nil {
		// 	fmt.Println("Error: ", err)
		// 	break
		// }
		// fmt.Println("string Data:", string(rawMsg))
		msg := asyncapi.BaseMessage{}
		err := dec.Decode(&msg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		// err = json.Unmarshal(rawMsg, &msg)
		// if err != nil {
		// 	fmt.Println("Error: ", err)
		// 	continue
		msgPipe <- msg
	}
}

func mapToStruct(data map[string]any, targetPointer any) {

	raw, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, targetPointer)
	if err != nil {
		panic(err)
	}
}
