package hubber_tmp

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// Request implement with json communication
type Request struct {
	senderID int64
	Action   string          `json:"action"`
	Data     json.RawMessage `json:"data,omitempty"`
}

func (r *Request) SenderID() int64 {
	return r.senderID
}

func (r *Request) GetAction() string {
	return r.Action
}

func (r *Request) ParseData(pointer interface{}) {
	if r.Data != nil {
		if err := json.Unmarshal(r.Data, pointer); err != nil {
			logrus.Error("Parsing error!")
			panic(err)
		}
	} else {
		logrus.Warn("No data in request!")
	}
}

func (r *Request) WriteData(pointer interface{}) {
	raw, err := json.Marshal(pointer)
	r.Data = raw
	if err != nil {
		logrus.Error("Writing error!")
		panic(err)
	}
}
