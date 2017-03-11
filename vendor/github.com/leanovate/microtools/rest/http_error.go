package rest

import (
	"fmt"
	"net/http"

	"github.com/leanovate/microtools/routing"
)

type HTTPError struct {
	Code    int    `json:"code" xml:"code"`
	Type    string `json:"type" xml:"type"`
	Message string `json:"message" xml:"message"`
	Details string `json:"details,omitempty" xml:"details,omitempty"`
}

func (e *HTTPError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s (%d): %s", e.Message, e.Code, e.Details)
	}
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}

func (e *HTTPError) Send(response http.ResponseWriter, encoder ResponseEncoder) {
	response.WriteHeader(e.Code)
	if err := encoder.Encode(response, e); err != nil {
		response.Write([]byte(e.Message))
	}
}

func (e *HTTPError) WithDetails(details string) *HTTPError {
	return &HTTPError{
		Code:    e.Code,
		Type:    e.Type,
		Message: e.Message,
		Details: details,
	}
}

func HttpErrorMatcher(httpError *HTTPError) routing.Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		encoder := StdResponseEncoderChooser(req)
		httpError.Send(resp, encoder)
		return true
	}
}

func WrapError(err error) *HTTPError {
	switch err.(type) {
	case *HTTPError:
		return err.(*HTTPError)
	default:
		return InternalServerError(err)
	}
}

var BadRequest = &HTTPError{
	Code:    400,
	Type:    "https://httpstatus.es/400",
	Message: "Bad request",
}

var UnauthorizedError = &HTTPError{
	Code:    401,
	Type:    "https://httpstatus.es/401",
	Message: "Unauthorized",
}

var Forbidden = &HTTPError{
	Code:    403,
	Type:    "https://httpstatus.es/403",
	Message: "Forbidden",
}

var NotFound = &HTTPError{
	Code:    404,
	Type:    "https://httpstatus.es/404",
	Message: "Not found",
}

var MethodNotAllowed = &HTTPError{
	Code:    405,
	Type:    "https://httpstatus.es/405",
	Message: "Method not allowed",
}

var Conflict = &HTTPError{
	Code:    409,
	Type:    "https://httpstatus.es/409",
	Message: "Conflict",
}

func InternalServerError(err error) *HTTPError {
	return &HTTPError{
		Code:    500,
		Type:    "https://httpstatus.es/500",
		Message: "InternalServerError",
		Details: err.Error(),
	}
}
