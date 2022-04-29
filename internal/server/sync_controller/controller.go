package sync_controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"

	"go.uber.org/zap"

	"rpg/internal/server/auth"
	cerrors "rpg/internal/server/custom_errors"
	"rpg/internal/server/matchmaker"
	"rpg/internal/server/sync_controller/syncapi"
	"rpg/internal/server/utils"
)

func NewController(
	logger *zap.SugaredLogger,
	auth auth.Service,
	matchMaker matchmaker.Service,
) *Controller {
	hdl := &Controller{
		lg:         logger,
		auth:       auth,
		matchMaker: matchMaker,
	}

	hdl.handlers = map[string]EpWrapper{
		registrationURL: {
			Endpoint:    hdl.makeRegistrationEndpoint(),
			PayloadType: reflect.TypeOf(syncapi.RegistrationPayloadIN{}),
			Description: "registration",
		},
		joinURL: {
			Endpoint:    hdl.makeJoinRoomEndpoint(),
			PayloadType: reflect.TypeOf(syncapi.JoinRoomPayloadIN{}),
			Description: "join room",
		},
	}

	return hdl
}

type Controller struct {
	sync.Mutex
	lg         *zap.SugaredLogger
	auth       auth.Service
	matchMaker matchmaker.Service
	handlers   map[string]EpWrapper
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	ctx := r.Context()

	c.lg.Info("Handle request...")
	handleError := makeHandleErrorFunc(w, enc, c.lg)

	bearerToken := r.Header.Get("Authorization")
	if bearerToken != utils.EmptyString {
		claims, err := c.auth.ParseClaims(bearerToken)
		if err != nil {
			handleError(err)
			return
		}
		ctx = context.WithValue(ctx, utils.UIDKey, claims.UserUID)
		ctx = context.WithValue(ctx, utils.LoginKey, claims.Login)
		fmt.Println("context uid: ", ctx.Value(utils.UIDKey))
		fmt.Println("context login: ", ctx.Value(utils.LoginKey))
	}

	epWrapper, ok := c.handlers[r.URL.Path]
	if !ok {
		handleError(errors.New("unknown method " + r.URL.Path))
		return
	}
	requestPointer := reflect.New(epWrapper.PayloadType).Interface()

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := dec.Decode(requestPointer)
	if err != nil {
		if err != io.EOF {
			handleError(err)
			return
		}
	}

	request := reflect.Indirect(reflect.ValueOf(requestPointer)).Interface()

	res, err := epWrapper.Endpoint(ctx, request)
	if err != nil {
		handleError(err)
		return
	}

	err = enc.Encode(
		syncapi.BaseResponse{
			Success: true,
			Data:    res,
		},
	)
	if err != nil {
		c.lg.Warnw("failed to encode response", "err", err)
		return
	}
}

func makeHandleErrorFunc(w http.ResponseWriter, enc *json.Encoder, lg *zap.SugaredLogger) func(errResp error) {
	return func(errResp error) {
		resp := &syncapi.BaseResponse{
			Success: false,
		}
		var status int

		defer func() {
			if status == 0 {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)

			err := enc.Encode(resp)
			if err != nil {
				lg.Warnw("failed to encode error", "err", err)
			}
		}()

		typedErr, ok := errResp.(cerrors.ErrorResponse)
		if !ok {
			resp.Data = errResp.Error()
			return
		}

		status = typedErr.HTTPStatusCode
		resp.Data = typedErr.Err.Error()
	}
}
