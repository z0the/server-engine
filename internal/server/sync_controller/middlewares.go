package sync_controller

import "context"

func MakeParseJWTMiddleware() Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request any) (any, error) {
			return next(ctx, request)
		}
	}
}
