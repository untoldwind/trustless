package rest

import (
	"errors"
	"net/http"

	"github.com/leanovate/microtools/routing"
)

type Resources interface {
	BeforeFilter(resp http.ResponseWriter, req *http.Request) bool
	Self() Link
	Create(request *http.Request) (Resource, error)
	List(request *http.Request) (interface{}, error)

	FindById(id string) (interface{}, error)
}

type ResourcesBase struct{}

func (ResourcesBase) BeforeFilter(resp http.ResponseWriter, req *http.Request) bool {
	return true
}

func (ResourcesBase) Create(*http.Request) (Resource, error) {
	return nil, HTTPMethodNotAllowed
}

func (ResourcesBase) List(*http.Request) (interface{}, error) {
	return nil, HTTPMethodNotAllowed
}

func (ResourcesBase) FindById(id string) (interface{}, error) {
	return nil, HTTPNotFound
}

type LimitedResources struct {
	ResourcesBase
	RequestSizeLimit int64
}

func (r LimitedResources) BeforeFilter(resp http.ResponseWriter, req *http.Request) bool {
	if req.Body != nil {
		req.Body = http.MaxBytesReader(resp, req.Body, r.RequestSizeLimit)
	}
	return true
}

func ResourcesMatcher(prefix string, collection Resources) routing.Matcher {
	return routing.PrefixSeq(prefix,
		routing.StringPart(func(id string) routing.Matcher {
			result, err := collection.FindById(id)
			if err != nil {
				return HTTPErrorMatcher(WrapError(err))
			}
			switch resource := (result).(type) {
			case Resource:
				return ResourceMatcher(resource)
			case Resources:
				return ResourcesMatcher("", resource)
			default:
				return HTTPErrorMatcher(HTTPInternalServerError(errors.New("Invalid result")))
			}
		}),
		routing.EndSeq(
			routing.GET(restHandler{before: collection.BeforeFilter, handler: collection.List}),
			routing.POST(createHandler{before: collection.BeforeFilter, handler: collection.Create}),
			HTTPErrorMatcher(HTTPMethodNotAllowed),
		),
	)
}
