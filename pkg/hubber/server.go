package hubber

import (
	"net"
	"runtime"

	"go.uber.org/zap"
)

type server struct {
	log            *zap.SugaredLogger
	port           string
	listener       net.Listener
	connectionCtrl ConnectionController
	shouldStop     bool
}

func NewServer(logger *zap.SugaredLogger, port string, connCtrl ConnectionController) *server {
	return &server{
		log:            logger,
		port:           port,
		connectionCtrl: connCtrl,
	}
}

func (s *server) Run() error {
	defer func() {
		if r := recover(); r != nil {
			s.log.Warn("Recovered in server: ", r)
		}
		s.log.Info("server has stopped...")
	}()

	s.log.Infof("Starting server on port %s...", s.port)

	err := s.loadListener()
	if err != nil {
		return err
	}

	for {
		if s.shouldStop {
			break
		}
		conn, err := s.listener.Accept()

		if err != nil {
			switch typedErr := err.(type) {
			case *net.OpError:
				if typedErr.Timeout() {
					continue
				}
			default:
				s.log.Fatal("Error during connectionWrapper conn attempt: ", err)
			}
		}
		s.connectionCtrl.HandleClientConnection(NewConnectionWrapper(conn, s.log, DefaultReadConnDelimiter))
		s.log.Info("Num of running goroutines: ", runtime.NumGoroutine())
	}
	return nil
}

func (s *server) loadListener() error {
	var err error
	s.listener, err = net.Listen("tcp", ":"+s.port)
	if err != nil {
		s.log.Fatal(err)
	}
	return err
}

func (s *server) Stop() {
	s.shouldStop = true
}
