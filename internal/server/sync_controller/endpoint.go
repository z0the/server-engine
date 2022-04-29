package sync_controller

import (
	"context"
	"reflect"
)

type Endpoint func(ctx context.Context, request any) (any, error)

type EpWrapper struct {
	Endpoint
	PayloadType reflect.Type
	Description string
}

type Middleware func(ep Endpoint) Endpoint
