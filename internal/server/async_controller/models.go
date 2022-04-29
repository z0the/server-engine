package async_controller

import (
	"rpg/internal/server/async_controller/asyncapi"
	"rpg/pkg/hubber"
)

type ClientConnWrapper struct {
	hubber.ConnectionWrapper
	authorized bool
	userUID    string
	login      string
}

type Handler func(req asyncapi.BaseMessage) (any, error)
