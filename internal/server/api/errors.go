package api

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidRequestType = errors.New("invalid request type")
)

type ErrorResponse struct {
	Err            error
	HTTPStatusCode int
	RPCStatusCode  int
}

func (e ErrorResponse) Error() string {
	return e.Err.Error()
}

func MakeBadRequestResponseError(err error) ErrorResponse {
	return ErrorResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func MakeResponseErrorWithHTTPStatusCode(err error, code int) ErrorResponse {
	return ErrorResponse{
		Err:            err,
		HTTPStatusCode: code,
	}
}
