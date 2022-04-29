package asyncapi

import (
	"encoding/json"
	"errors"

	"rpg/internal/server/custom_errors"
)

func NewBaseMessage(reqType MsgType, payload any) (BaseMessage, error) {
	baseReq := BaseMessage{MsgType: reqType}
	err := baseReq.EncodeJSONPayload(payload)
	if err != nil {
		return BaseMessage{}, err
	}
	return baseReq, nil
}

type BaseMessage struct {
	connUID string
	MsgType MsgType `json:"req_type"`
	Payload []byte  `json:"payload"`
}

func (r *BaseMessage) GetConnUID() string {
	return r.connUID
}

func (r *BaseMessage) SetConnUID(connUID string) {
	r.connUID = connUID
}

func (r *BaseMessage) Validate() error {
	err := r.MsgType.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *BaseMessage) RequestTypeCode() string {
	return r.MsgType.String()
}

func (r *BaseMessage) EncodeJSONPayload(payload any) error {
	if len(r.Payload) != 0 {
		return errors.New("Payload is not empty")
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r.Payload = rawPayload
	return nil
}

// DecodeJSONPayload - parse byte data from internal to targetPayload
// targetPayload should be a pointer to a struct
func (r *BaseMessage) DecodeJSONPayload(targetPayload any) error {
	err := json.Unmarshal(r.Payload, targetPayload)
	if err != nil {
		return err
	}
	return nil
}

func (r *BaseMessage) CloneWithoutPayload() BaseMessage {
	return BaseMessage{
		connUID: r.connUID,
		MsgType: r.MsgType,
	}
}

type MsgType string

const (
	LoginMsgType = "LOGIN"
)

func (rt MsgType) String() string {
	return string(rt)
}

func (rt MsgType) Validate() error {
	_, ok := getAllRequestTypes()[rt]
	if !ok {
		return custom_errors.MakeWrongRequestTypeError(rt.String())
	}
	return nil
}

func getAllRequestTypes() map[MsgType]struct{} {
	return map[MsgType]struct{}{
		LoginMsgType:         {},
		SceneUpdateMsgType:   {},
		ServerConnectMsgType: {},
		MoveMsgType:          {},
	}
}

type LoginPayload struct {
	Token string `json:"token"`
}

type LoginResPayload struct {
	Success bool `json:"success"`
}

type ConnectToRoomPayload struct {
	RoomUID string `json:"room_uid"`
}

type ConnectToRoomPayloadOut struct {
	Success bool `json:"success"`
}

type ServerConnectPayloadOUT struct {
	ConnectionUID string `json:"connection_uid"`
}
