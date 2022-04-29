package eventbus

type BusMsg struct {
	RecipientConnUID string
	MsgType          string
	Payload          any
}

type MsgHandler func(msg BusMsg)
