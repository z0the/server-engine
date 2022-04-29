package sync_controller

import (
	"context"
	"errors"

	"rpg/internal/server/sync_controller/syncapi"
	"rpg/internal/server/utils"
)

func (c *Controller) makeRegistrationEndpoint() func(ctx context.Context, rawReq any) (any, error) {
	return func(ctx context.Context, rawReq any) (any, error) {
		req, ok := rawReq.(syncapi.RegistrationPayloadIN)
		if !ok {
			return nil, wrongReqType
		}

		res, err := c.registrationHandler(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func (c *Controller) registrationHandler(_ context.Context, req syncapi.RegistrationPayloadIN) (
	syncapi.RegistrationPayloadOUT,
	error,
) {

	signedToken, err := c.auth.RegisterNewUser(req.Login, req.Password)
	if err != nil {
		return syncapi.RegistrationPayloadOUT{}, err
	}

	return syncapi.RegistrationPayloadOUT{
		Token: signedToken,
	}, nil
}

func (c *Controller) makeLoginEndpoint() func(ctx context.Context, rawReq any) (any, error) {
	return func(ctx context.Context, rawReq any) (any, error) {
		req, ok := rawReq.(syncapi.LoginPayloadIN)
		if !ok {
			return nil, wrongReqType
		}

		res, err := c.loginHandler(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func (c *Controller) loginHandler(_ context.Context, req syncapi.LoginPayloadIN) (
	syncapi.LoginPayloadOUT,
	error,
) {

	signedToken, err := c.auth.RegisterNewUser(req.Login, req.Password)
	if err != nil {
		return syncapi.LoginPayloadOUT{}, err
	}

	return syncapi.LoginPayloadOUT{
		Token: signedToken,
	}, nil
}

func (c *Controller) makeJoinRoomEndpoint() func(ctx context.Context, rawReq any) (any, error) {
	return func(ctx context.Context, rawReq any) (any, error) {
		payload, ok := rawReq.(syncapi.JoinRoomPayloadIN)
		if !ok {
			return nil, wrongReqType
		}

		res, err := c.joinRoomHandler(ctx, payload)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

func (c *Controller) joinRoomHandler(ctx context.Context, req syncapi.JoinRoomPayloadIN) (
	syncapi.JoinRoomPayloadOUT,
	error,
) {
	userUID, ok := ctx.Value(utils.UIDKey).(string)
	if !ok {
		return syncapi.JoinRoomPayloadOUT{}, errors.New("no uid in ctx")
	}
	login, ok := ctx.Value(utils.LoginKey).(string)
	if !ok {
		return syncapi.JoinRoomPayloadOUT{}, errors.New("no login in ctx")
	}

	roomUID := c.matchMaker.JoinRoom(req.ConnectionUID, userUID, login)

	return syncapi.JoinRoomPayloadOUT{
		RoomUID: roomUID,
	}, nil
}
