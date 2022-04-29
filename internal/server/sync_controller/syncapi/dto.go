package syncapi

type BaseResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type RegistrationPayloadIN struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegistrationPayloadOUT struct {
	Token string `json:"token"`
}

type LoginPayloadIN struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginPayloadOUT struct {
	Token string `json:"token"`
}

type JoinRoomPayloadIN struct {
	ConnectionUID string `json:"connection_uid"`
}

type JoinRoomPayloadOUT struct {
	RoomUID string `json:"room_uid"`
}

type ServerConnectPayloadIn struct {
	ConnectionUID string `json:"connection_uid"`
}

type ServerConnectPayloadOut struct{}
