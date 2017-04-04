package rest

import (
	"net/http"

	"github.com/leanovate/microtools/routing"
)

type Resource interface {
	BeforeFilter(resp http.ResponseWriter, req *http.Request) bool
	Self() Link
	Get(request *http.Request) (interface{}, error)
	Post(request *http.Request) (interface{}, error)
	Patch(request *http.Request) (interface{}, error)
	Update(request *http.Request) (interface{}, error)
	Delete(request *http.Request) (interface{}, error)

	SubResources() routing.Matcher
}

type ResourceBase struct{}

func (ResourceBase) BeforeFilter(resp http.ResponseWriter, req *http.Request) bool {
	return true
}

func (ResourceBase) Get(request *http.Request) (interface{}, error) {
	return nil, HTTPMethodNotAllowed
}

func (ResourceBase) Post(request *http.Request) (interface{}, error) {
	return nil, HTTPMethodNotAllowed
}

func (ResourceBase) Patch(request *http.Request) (interface{}, error) {
	return nil, HTTPMethodNotAllowed
}

func (ResourceBase) Update(request *http.Request) (interface{}, error) {
	return nil, HTTPMethodNotAllowed
}

func (ResourceBase) Delete(request *http.Request) (interface{}, error) {
	return nil, HTTPMethodNotAllowed
}

func (ResourceBase) SubResources() routing.Matcher {
	return HTTPErrorMatcher(HTTPNotFound)
}

type LimitedResource struct {
	ResourceBase
	RequestSizeLimit int64
}

func (r LimitedResource) BeforeFilter(resp http.ResponseWriter, req *http.Request) bool {
	if req.Body != nil {
		req.Body = http.MaxBytesReader(resp, req.Body, r.RequestSizeLimit)
	}
	return true
}

func ResourceMatcher(resource Resource) routing.Matcher {
	return routing.Sequence(
		routing.EndSeq(
			routing.GET(restHandler{before: resource.BeforeFilter, handler: resource.Get}),
			routing.POST(restHandler{before: resource.BeforeFilter, handler: resource.Post}),
			routing.PUT(restHandler{before: resource.BeforeFilter, handler: resource.Update}),
			routing.PATCH(restHandler{before: resource.BeforeFilter, handler: resource.Patch}),
			routing.DELETE(restHandler{before: resource.BeforeFilter, handler: resource.Delete}),
			HTTPErrorMatcher(HTTPMethodNotAllowed),
		),
		resource.SubResources(),
	)
}
