package rest

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
)

type restHandler struct {
	before  func(resp http.ResponseWriter, req *http.Request) bool
	handler func(request *http.Request) (interface{}, error)
}

func (h restHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	encoder := StdResponseEncoderChooser(req)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, string(debug.Stack()))
			HTTPInternalServerError(fmt.Errorf("Paniced: %v", r)).Send(resp, encoder)
		}
	}()
	var err error
	if h.before != nil && !h.before(resp, req) {
		return
	}
	result, err := h.handler(req)
	if err == nil {
		switch result.(type) {
		case *Result:
			err = result.(*Result).Send(resp, encoder)
		default:
			err = Ok().WithBody(result).Send(resp, encoder)
		}
	}
	if err != nil {
		WrapError(err).Send(resp, encoder)
	}
}

type createHandler struct {
	before  func(resp http.ResponseWriter, req *http.Request) bool
	handler func(*http.Request) (Resource, error)
}

func (h createHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	encoder := StdResponseEncoderChooser(req)
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, string(debug.Stack()))
			HTTPInternalServerError(fmt.Errorf("Paniced: %v", r)).Send(resp, encoder)
		}
	}()
	var err error
	if h.before != nil && !h.before(resp, req) {
		return
	}
	var resource Resource
	resource, err = h.handler(req)
	if err == nil {
		if resource != nil {
			var result interface{}
			result, err = resource.Get(req)
			if err == nil {
				switch result.(type) {
				case *Result:
					err = result.(*Result).
						AddHeader("location", resource.Self().Href).
						WithStatus(201).
						Send(resp, encoder)
				default:
					err = Created().
						AddHeader("location", resource.Self().Href).
						WithBody(result).
						Send(resp, encoder)
				}
			}
		}
	}
	if err != nil {
		WrapError(err).Send(resp, encoder)
	}
}
