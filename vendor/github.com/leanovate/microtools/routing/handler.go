package routing

import "net/http"

// RouteHandler is a http.Handler delegating the request to a Matcher.
type RouteHandler struct {
	Matcher Matcher
}

// NewRouteHandler creates a RouteHandler for a sequence of Matcher.
func NewRouteHandler(matchers ...Matcher) *RouteHandler {
	return &RouteHandler{
		Matcher: Sequence(matchers...),
	}
}

func (r *RouteHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if !r.Matcher(req.URL.Path, resp, req) {
		resp.WriteHeader(404)
	}
}

func RouteHandlerFunc(matcher Matcher) func(http.ResponseWriter, *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		if !matcher(req.URL.Path, resp, req) {
			resp.WriteHeader(404)
		}
	}
}
