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

	"github.com/sirupsen/logrus"

	"rpg/internal/server/api"
	"rpg/internal/server/matchmaker"
	"rpg/internal/server/utils"
)

func NewController(logger *logrus.Logger, services *matchmaker.Services) *Controller {
	hdl := &Controller{
		lg:       logger,
		services: services,
	}

	hdl.handlers = map[string]EpWrapper{
		registrationURL: {
			Endpoint:    hdl.makeRegistrationEndpoint(),
			RequestType: reflect.TypeOf(api.RegistrationRequest{}),
			Description: "registration",
		},
	}

	return hdl
}

type Controller struct {
	sync.Mutex
	lg       *logrus.Logger
	services *matchmaker.Services
	handlers map[string]EpWrapper
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	ctx := r.Context()

	handleError := makeHandleErrorFunc(w, enc, c.lg)

	bearerToken := r.Header.Get("Authorization")
	if bearerToken != utils.EmptyString {
		claims, err := c.services.Auth.ParseClaims(bearerToken)
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
	requestPointer := reflect.New(epWrapper.RequestType).Interface()

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
		api.BaseResponse{
			Success: true,
			Data:    res,
		},
	)
	if err != nil {
		c.lg.WithError(err).Warn("Failed to encode response")
		return
	}
}

func makeHandleErrorFunc(w http.ResponseWriter, enc *json.Encoder, lg *logrus.Logger) func(errResp error) {
	return func(errResp error) {
		resp := &api.BaseResponse{
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
				lg.WithError(err).Error("failed to encode error")
			}
		}()

		typedErr, ok := errResp.(api.ErrorResponse)
		if !ok {
			resp.Data = errResp.Error()
			return
		}

		status = typedErr.HTTPStatusCode
		resp.Data = typedErr.Err.Error()
	}
}
