package custom_errors

import (
	"fmt"
	"net/http"
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

func MakeWrongRequestTypeError(wrongReqType string) error {
	return fmt.Errorf("invalid request type %s", wrongReqType)
}
