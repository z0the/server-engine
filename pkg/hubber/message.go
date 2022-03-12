package hubber

func NewMessage(connUID string, rawData []byte) Message {
	return &message{
		connUID: connUID,
		rawData: rawData,
	}
}

type message struct {
	connUID string
	rawData []byte
}

func (m *message) GetConnUID() string {
	return m.connUID
}

func (m *message) GetRawData() []byte {
	return append([]byte{}, m.rawData...)
}
