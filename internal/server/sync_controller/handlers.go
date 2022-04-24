package sync_controller

import (
	"context"

	"rpg/internal/server/api"
)

func (c *Controller) makeRegistrationEndpoint() func(ctx context.Context, rawReq any) (any, error) {
	return func(ctx context.Context, rawReq any) (any, error) {
		req, ok := rawReq.(api.RegistrationRequest)
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

func (c *Controller) registrationHandler(_ context.Context, req api.RegistrationRequest) (
	api.RegistrationResponse,
	error,
) {

	signedToken, err := c.services.Auth.RegisterNewUser(req.Login, req.Password)
	if err != nil {
		return api.RegistrationResponse{}, err
	}

	return api.RegistrationResponse{
		Token: signedToken,
	}, nil
}

func (c *Controller) makeLoginEndpoint() func(ctx context.Context, rawReq any) (any, error) {
	return func(ctx context.Context, rawReq any) (any, error) {
		req, ok := rawReq.(api.LoginRequest)
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

func (c *Controller) loginHandler(_ context.Context, req api.LoginRequest) (
	api.LoginResponse,
	error,
) {

	signedToken, err := c.services.Auth.RegisterNewUser(req.Login, req.Password)
	if err != nil {
		return api.LoginResponse{}, err
	}

	return api.LoginResponse{
		Token: signedToken,
	}, nil
}

func (c *Controller) makeJoinRoomEndpoint() func(ctx context.Context, rawReq any) (any, error) {
	return func(ctx context.Context, rawReq any) (any, error) {
		payload, ok := rawReq.(api.LoginRequest)
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

func (c *Controller) joinRoomHandler(_ context.Context, req api.LoginRequest) (
	api.LoginResponse,
	error,
) {

	signedToken, err := c.services.Auth.RegisterNewUser(req.Login, req.Password)
	if err != nil {
		return api.LoginResponse{}, err
	}

	return api.LoginResponse{
		Token: signedToken,
	}, nil
}
