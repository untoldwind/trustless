package routing

import "net/http"

// Matcher matches a HTTP request. If the matcher returns true, the HTTP
// request is considered to be handled and no other matcher should be tried.
type Matcher func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool
