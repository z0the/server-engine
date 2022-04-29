package async_controller

import (
	"errors"

	"rpg/internal/server/async_controller/asyncapi"
)

func (h *Controller) loginHandler(reqMsg asyncapi.BaseMessage) (any, error) {
	payload := asyncapi.LoginPayload{}
	err := reqMsg.DecodeJSONPayload(payload)
	if err != nil {
		return nil, err
	}

	claims, err := h.auth.ParseClaims(payload.Token)
	if err != nil {
		return nil, err
	}

	clientConn := h.conns[reqMsg.GetConnUID()]
	clientConn.userUID = claims.UserUID
	clientConn.login = claims.Login
	clientConn.authorized = true

	result := asyncapi.LoginResPayload{
		Success: true,
	}

	return result, nil
}

func (h *Controller) ConnectToRoomHandler(reqMsg asyncapi.BaseMessage) (any, error) {
	if !h.conns[reqMsg.GetConnUID()].authorized {
		return nil, errors.New("not authorized")
	}

	payload := asyncapi.ConnectToRoomPayload{}
	err := reqMsg.DecodeJSONPayload(payload)
	if err != nil {
		return nil, err
	}

	result := asyncapi.LoginResPayload{
		Success: true,
	}

	return result, nil
}
