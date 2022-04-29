package asyncapi

const (
	SceneUpdateMsgType   MsgType = "SCENE_UPDATE"
	ErrorMsgType         MsgType = "ERROR"
	ServerConnectMsgType MsgType = "SERVER_CONNECT"
	MoveMsgType          MsgType = "MOVE"
)

type ServerMessage struct {
	MsgType          MsgType
	RecipientConnUID string
	Payload          any
}

type ErrorPayloadOut struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

const (
	emptyErrorCode = 0
)

type SceneUpdatePayloadOUT struct {
	PlayersList []Player `json:"players_list"`
}

type MovePayloadIN struct {
	Direction string `json:"direction"`
	Speed     int    `json:"speed"`
}

type MovePayloadOUT struct{}

type Player struct {
	UID   string `json:"uid"`
	PosX  int    `json:"posX"`
	PosY  int    `json:"posY"`
	Login string `json:"login"`
}
