package hubber_tmp

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// Response implement with json communication
type Response struct {
	receiverID int64
	Action     string          `json:"action"`
	Data       json.RawMessage `json:"data,omitempty"`
}

func (r *Response) SetReceiverID(id int64) {
	r.receiverID = id
}

func (r *Response) ReceiverID() int64 {
	return r.receiverID
}

func (r *Response) GetAction() string {
	return r.Action
}

func (r *Response) ParseData(pointer interface{}) {
	if r.Data != nil {
		if err := json.Unmarshal(r.Data, pointer); err != nil {
			logrus.Error("Parsing error!")
			panic(err)
		}
	} else {
		logrus.Warn("No data in request!")
	}
}

func (r *Response) WriteData(pointer interface{}) {
	raw, err := json.Marshal(pointer)
	r.Data = raw
	if err != nil {
		logrus.Error("Writing error!")
		panic(err)
	}
}
