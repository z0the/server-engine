package hubber_tmp

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

const Dead = "DEAD"

func NewClient(conn net.Conn, log *logrus.Logger) *Client {
	return &Client{
		conn:     conn,
		log:      log,
		sendChan: make(chan IResponse, 4),
	}
}

type Client struct {
	id         int64
	authorized bool
	log        *logrus.Logger
	sendChan   chan IResponse
	conn       net.Conn
}

func (c *Client) Run() {
}

func (c *Client) Kill() {
	c.sendChan = nil
	if err := c.conn.Close(); err != nil {
		c.log.Error(err)
	}
}

func (c *Client) Send(resp IResponse) {
	if c.sendChan != nil {
		c.sendChan <- resp
	}
}

func (c *Client) sendClientIsDead(handler IHandler) {
	req := &Request{}
	req.senderID = c.id
	req.Action = Dead
	handler.Handle(req)
}

func (c *Client) readPump(handler IHandler) {
	defer func() {
		c.log.Warn("Stop read pump...")
		if err := c.conn.Close(); err != nil {
			c.log.Error(err)
		}
		c.sendChan = nil
		c.sendClientIsDead(handler)
	}()
	decoder := json.NewDecoder(c.conn)
	for {
		req := &Request{}
		//if err := c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); err != nil {
		//	c.log.Error(err)
		//	break
		//}
		fmt.Println("start decoding")
		if err := decoder.Decode(&req); err != nil {
			//c.log.Error(err)
			break
		}
		if req.Action != "ping" {
			c.log.Infof("received from %s : %s", c.conn.RemoteAddr(), req)
		}
		req.senderID = c.id
		handler.Handle(req)
	}
}

func (c *Client) writePump() {
	defer func() {
		c.log.Warn("Stop write pump...")
		if err := c.conn.Close(); err != nil {
			c.log.Error(err)
		}
	}()
	encoder := json.NewEncoder(c.conn)
	for message := range c.sendChan {
		//if err := c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		//	c.log.Error(err)
		//	break
		//}
		if message.GetAction() != "pong" {
			c.log.Infof("send to    %s : %s", c.conn.RemoteAddr(), message)
		}
		if err := encoder.Encode(message); err != nil {
			c.log.Error(err)
			break
		}
	}
}
