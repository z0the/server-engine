package hubber_tmp

type IServer interface {
	Run(handler IHandler)
	Stop()
}

type IHandler interface {
	Register(client IClient) int64
	Unregister(id int64)
	Handle(request IRequest)
}

type IClient interface {
	Kill()
	AsyncSend(resp IResponse)
}

type IResponse interface {
	SetReceiverID(id int64)
	ReceiverID() int64
	GetAction() string
	ParseData(pointer interface{})
	WriteData(pointer interface{})
}

type IRequest interface {
	SenderID() int64
	GetAction() string
	ParseData(pointer interface{})
	WriteData(pointer interface{})
}
