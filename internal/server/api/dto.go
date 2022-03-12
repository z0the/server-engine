package api

import (
	"encoding/json"
	"errors"
)

type ServerMessage struct {
	IsSyncResponse bool   `json:"is_sync_response"`
	Payload        []byte `json:"payload"`
}

func NewBaseRequest(reqType RequestType, payload interface{}) (*BaseRequest, error) {
	baseReq := &BaseRequest{ReqType: reqType}
	err := baseReq.EncodeJSONPayload(payload)
	if err != nil {
		return nil, err
	}
	return baseReq, nil
}

type BaseRequest struct {
	connUID string
	ReqType RequestType `json:"req_type"`
	Payload []byte      `json:"payload"`
}

func (r *BaseRequest) GetConnUID() string {
	return r.connUID
}

func (r *BaseRequest) SetConnUID(connUID string) {
	r.connUID = connUID
}

func (r *BaseRequest) Validate() error {
	err := r.ReqType.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (r *BaseRequest) RequestTypeCode() string {
	return r.ReqType.String()
}

func (r *BaseRequest) EncodeJSONPayload(payload interface{}) error {
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
func (r *BaseRequest) DecodeJSONPayload(targetPayload interface{}) error {
	err := json.Unmarshal(r.Payload, targetPayload)
	if err != nil {
		return err
	}
	return nil
}

type RequestType string

func (rt RequestType) String() string {
	return string(rt)
}

func (rt RequestType) Validate() error {
	for _, reqType := range getAllRequestTypes() {
		if reqType == rt {
			return nil
		}
	}
	return ErrInvalidRequestType
}

const (
	RegistrationReqType RequestType = "REGISTRATION"
	LoginReqType        RequestType = "LOGIN"
)

func getAllRequestTypes() []RequestType {
	return []RequestType{
		RegistrationReqType,
		LoginReqType,
	}
}

type RegistrationPayload struct {
	Login    string
	Password string
}

type LoginPayload struct {
	Login    string
	Password string
}
