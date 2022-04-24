package hubber_tmp

import (
	"context"
	"fmt"
	"net"
	"runtime"

	"github.com/sirupsen/logrus"
)

type ConnectionHandler interface {
	HandleConnection(ctx context.Context, conn net.Conn)
	CloneWithServices() ConnectionHandler
}

type Server struct {
	port       string
	listener   net.Listener
	log        *logrus.Logger
	shouldStop bool
}

func NewServer(port string, logger *logrus.Logger) IServer {
	return &Server{
		port: port,
		log:  logger,
	}
}

func (s *Server) Run(handler IHandler) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Warn("Recovered in server: ", r)
		}
		s.log.Info("Server has stopped...")
	}()
	s.log.Infof("Starting server on port %s...", s.port)
	s.loadListener()
	for {
		if s.shouldStop {
			return
		}
		conn, err := s.listener.Accept()

		if err != nil {
			switch typedErr := err.(type) {
			case *net.OpError:
				if typedErr.Timeout() {
					continue
				}
			default:
				s.log.Fatal("Error during client conn attempt: ", err)
			}
		}
		// runClient(conn, s.log, async_controller)
		fmt.Println(conn)
		s.log.Info("start new client...")
		s.log.Info("Num of running gorutines: ", runtime.NumGoroutine())
	}
}

func (s *Server) loadListener() {
	var err error
	s.listener, err = net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.log.Fatal(err)
	}
}

func (s *Server) Stop() {
	s.shouldStop = true
}
