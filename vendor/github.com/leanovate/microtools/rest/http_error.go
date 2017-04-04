package rest

import (
	"fmt"
	"net/http"

	"github.com/leanovate/microtools/routing"
)

// HTTPError is an error result of a HTTP/REST operation.
// Implements the Error interface.
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

// Send the HTTPError the a http.ResponseWriter
func (e *HTTPError) Send(response http.ResponseWriter, encoder ResponseEncoder) {
	response.WriteHeader(e.Code)
	if err := encoder.Encode(response, e); err != nil {
		response.Write([]byte(e.Message))
	}
}

// WithDetails creates a new HTTPError with extra detail message
func (e *HTTPError) WithDetails(details string) *HTTPError {
	return &HTTPError{
		Code:    e.Code,
		Type:    e.Type,
		Message: e.Message,
		Details: details,
	}
}

// HTTPErrorMatcher is a routing.Matcher that always response with a given
// HTTPError. Usually useful at the end of a routing chain as catch all
// for MethodNotAllowed or NotFound-
func HTTPErrorMatcher(httpError *HTTPError) routing.Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		encoder := StdResponseEncoderChooser(req)
		httpError.Send(resp, encoder)
		return true
	}
}

// WrapError wrap a generic error as HTTPError.
// If err already is a HTTPError it will be left intact, otherwise the error
// will be mapped to InternalServerError
func WrapError(err error) *HTTPError {
	switch err.(type) {
	case *HTTPError:
		return err.(*HTTPError)
	default:
		return HTTPInternalServerError(err)
	}
}

// HTTPBadRequest is a HTTP bad request 400
var HTTPBadRequest = &HTTPError{
	Code:    400,
	Type:    "https://httpstatus.es/400",
	Message: "Bad request",
}

// HTTPUnauthorized is a HTTP Unauthorized 401
var HTTPUnauthorized = &HTTPError{
	Code:    401,
	Type:    "https://httpstatus.es/401",
	Message: "Unauthorized",
}

// HTTPForbidden is a HTTP forbidden 403
var HTTPForbidden = &HTTPError{
	Code:    403,
	Type:    "https://httpstatus.es/403",
	Message: "Forbidden",
}

// HTTPNotFound is a HTTP not found 404
var HTTPNotFound = &HTTPError{
	Code:    404,
	Type:    "https://httpstatus.es/404",
	Message: "Not found",
}

// HTTPMethodNotAllowed is a HTTP method not allowed 405
var HTTPMethodNotAllowed = &HTTPError{
	Code:    405,
	Type:    "https://httpstatus.es/405",
	Message: "Method not allowed",
}

// HTTPConflict is a HTTP conflict 409
var HTTPConflict = &HTTPError{
	Code:    409,
	Type:    "https://httpstatus.es/409",
	Message: "Conflict",
}

// HTTPInternalServerError wraps an arbitrary Error as
// HTTP internal server error 500
func HTTPInternalServerError(err error) *HTTPError {
	return &HTTPError{
		Code:    500,
		Type:    "https://httpstatus.es/500",
		Message: "InternalServerError",
		Details: err.Error(),
	}
}
