package hubber

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
)

func NewConnectionWrapper(
	conn net.Conn,
	log *zap.SugaredLogger,
	readDelimiter rune,
) ConnectionWrapper {
	return &connectionWrapper{
		conn:          conn,
		readDelimiter: readDelimiter,
		log:           log,
	}
}

type connectionWrapper struct {
	conn          net.Conn
	readDelimiter rune
	isDead        bool
	log           *zap.SugaredLogger
	uid           string
	once          sync.Once
}

func (c *connectionWrapper) StartReading(uid string, outChan chan<- Message) {
	c.once.Do(
		func() {
			c.uid = uid
			go c.connReader(outChan)
		},
	)
}

func (c *connectionWrapper) Kill() {
	if err := c.conn.Close(); err != nil {
		c.log.Error(err)
	}
}

func (c *connectionWrapper) GetLastRequestNumber() uint {
	return 0
}

func (c *connectionWrapper) Send(msg []byte) {
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				c.log.Error("recovered in async send: ", r)
			}
		}()

		c.log.Info("sending...")
		c.log.Info("msg: ", string(msg))
		_, err := c.conn.Write(c.addDelimiter(msg))
		if err != nil {
			c.log.Errorw("failed to write message to conn", "err", err)
		}
	}()
}

const DefaultReadConnDelimiter rune = '\n'

func (c *connectionWrapper) addDelimiter(rawMsg []byte) []byte {
	return append(rawMsg, byte(c.readDelimiter))
}

// nolint
// TODO: Make gracefull shutdown for listen and send pump
func (c *connectionWrapper) connReader(outChan chan<- Message) {
	defer func() {
		c.log.Warn("Stop listen pump...")
		err := c.conn.Close()
		if err != nil {
			c.log.Error(err)
		}
	}()
	c.log.Info("start reader...")

	bufReader := bufio.NewReader(c.conn)
	for {
		if c.isDead {
			break
		}
		c.log.Info("listening...")
		rawReq, err := bufReader.ReadBytes(byte(c.readDelimiter))
		if err != nil {
			c.log.Error(err)
			break
		}
		fmt.Println("string Data:", string(rawReq))
		outChan <- NewMessage(c.uid, rawReq)
	}
}
