package routing

import "net/http"

// Method create a matcher that matches a specific HTTP method
func Method(method string, handler http.Handler) Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		if req.Method == method {
			handler.ServeHTTP(resp, req)
			return true
		}
		return false
	}
}

// GET creates a matcher that matches HTTP GET
func GET(handler http.Handler) Matcher {
	return Method("GET", handler)
}

// GETFunc create a matcher that delegates to a function on HTTP GET
func GETFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return GET(http.HandlerFunc(handler))
}

// POST creates a matcher that matches HTTP POST
func POST(handler http.Handler) Matcher {
	return Method("POST", handler)
}

// POSTFunc create a matcher that delegates to a function on HTTP POST
func POSTFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return POST(http.HandlerFunc(handler))
}

// PUT creates a matcher that matches HTTP PUT
func PUT(handler http.Handler) Matcher {
	return Method("PUT", handler)
}

// PUTFunc create a matcher that delegates to a function on HTTP PUT
func PUTFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return PUT(http.HandlerFunc(handler))
}

// PATCH creates a matcher that matches HTTP PATCH
func PATCH(handler http.Handler) Matcher {
	return Method("PATCH", handler)
}

// PATCHFunc create a matcher that delegates to a function on HTTP PATCH
func PATCHFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return PATCH(http.HandlerFunc(handler))
}

// DELETE creates a matcher that matches HTTP DELETE
func DELETE(handler http.Handler) Matcher {
	return Method("DELETE", handler)
}

// DELETEFunc create a matcher that delegates to a function on HTTP DELETE
func DELETEFunc(handler func(http.ResponseWriter, *http.Request)) Matcher {
	return DELETE(http.HandlerFunc(handler))
}

// MethodNotAllowed creates a matcher that matches everything and always
// replies with a 405. Useful at the end of a sequence of method matchers.
func MethodNotAllowed(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
	resp.WriteHeader(405)
	return true
}
