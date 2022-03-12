package hubber

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

func NewClient(conn net.Conn, log *logrus.Logger) ClientConnection {
	return &client{
		conn:     conn,
		log:      log,
		sendChan: make(chan []byte, 4),
	}
}

type client struct {
	conn     net.Conn
	isDead   bool
	log      *logrus.Logger
	sendChan chan []byte
	uid      string
	once     sync.Once
}

func (c *client) Run(uid string, requestChan chan<- Message) {
	c.uid = uid
	go c.listenPump(requestChan)
	go c.sendPump()
}

func (c *client) Kill() {
	c.sendChan = nil
	if err := c.conn.Close(); err != nil {
		c.log.Error(err)
	}
}

func (c *client) GetLastRequestNumber() uint {
	return 0
}

func (c *client) Send(msg []byte) {
	panic("not implemented")
}

func (c *client) AsyncSend(msg []byte) {
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				c.log.Error("recovered in async send: ", r)
			}
		}()
		if c.sendChan != nil {
			c.sendChan <- msg
		}
	}()
}

func (c *client) sendClientIsDead() {
	c.once.Do(
		func() {
			// c.connectionIsDeadChan <- c.uid
		},
	)
}

const defaultReadConnDelimiter byte = '\n'

// nolint
// TODO: Make gracefull shutdown for listen and send pump
func (c *client) listenPump(requestChan chan<- Message) {
	defer func() {
		c.log.Warn("Stop listen pump...")
		err := c.conn.Close()
		if err != nil {
			c.log.Error(err)
		}
		// Closing the channel to stop the sending pump
		// close(c.sendChan)
		c.sendClientIsDead()
	}()
	bufReader := bufio.NewReader(c.conn)
	for {
		if c.isDead {
			break
		}
		c.log.Info("listening...")
		rawReq, err := bufReader.ReadBytes(defaultReadConnDelimiter)
		if err != nil {
			c.log.Error(err)
			break
		}
		fmt.Println("string Data:", string(rawReq))
		requestChan <- NewMessage(c.uid, rawReq)
	}
}

func (c *client) sendPump() {
	defer func() {
		c.log.Warn("Stop send pump...")
		if err := c.conn.Close(); err != nil {
			c.log.WithError(err).Error("failed to close connection")
		}
		c.sendClientIsDead()
	}()

	for msg := range c.sendChan {
		if c.isDead {
			break
		}
		c.log.Info("sending...")
		_, err := c.conn.Write(msg)
		if err != nil {
			c.log.WithError(err).Error("failed to write message to conn")
		}
	}
}
