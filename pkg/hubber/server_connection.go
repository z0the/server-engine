package hubber

import (
	"encoding/json"
	"fmt"
	"net"

	"rpg/internal/server/async_controller/asyncapi"
)

func NewServerConnection(conn net.Conn) *serverConnection {
	return &serverConnection{
		conn:      conn,
		syncChan:  make(chan []byte),
		asyncChan: make(chan []byte, 10),
	}
}

type serverConnection struct {
	conn      net.Conn
	syncChan  chan []byte
	asyncChan chan []byte
}

func (sc *serverConnection) listenConn() {
	decoder := json.NewDecoder(sc.conn)

	for {
		msg := new(asyncapi.ServerMessage)
		err := decoder.Decode(msg)
		if err != nil {
			panic(err)
		}
		fmt.Println("response: ", msg)
	}
}
