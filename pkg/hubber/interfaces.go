package hubber

type GameController interface {
	// HandleClientConnection - when called should register client connection in system,
	// and then call the ClientConnection.Run() method
	HandleClientConnection(client ClientConnection)
	// Handle(msg interface{})
}

type ServerConnection interface {
	Run(syncChan, asyncChan <-chan []byte)
	Send(msg Message)
}

type ClientConnection interface {
	// Run - should start listen and send pumps
	// requestChan - is a chan to which client should sends msgs received from low level connection
	// connectionIsDeadChan - is a chan to which the client should sends its uid when connection becomes dead.
	Run(connUID string, requestChan chan<- Message)
	// Kill - should stops all client serving goroutines and notify the connectionIsDeadChan.
	Kill()
	AsyncSend(msg []byte)
	GetLastRequestNumber() uint
	// Send(msg []byte)
}

type Message interface {
	GetConnUID() string
	GetRawData() []byte
}
