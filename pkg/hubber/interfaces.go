package hubber

type ConnectionController interface {
	// HandleClientConnection - when called should register connectionWrapper connection in system,
	// and then call the ConnectionWrapper.Run() method
	HandleClientConnection(client ConnectionWrapper)
	// Handle(msg any)
}

type ServerConnection interface {
	Run(syncChan, asyncChan <-chan []byte)
	Send(msg Message)
}

type ConnectionWrapper interface {
	// StartReading
	// connectionIsDeadChan - is a chan to which the connectionWrapper should sends its uid when connection becomes dead.
	StartReading(connUID string, outChan chan<- Message)
	// Kill - should stops all connectionWrapper serving goroutines and notify the connectionIsDeadChan.
	Kill()
	Send(msg []byte)
	GetLastRequestNumber() uint
}

type Message interface {
	GetConnUID() string
	GetRawData() []byte
}
