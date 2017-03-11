package routing

import "net/http"

func Sequence(matchers ...Matcher) Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		for _, matcher := range matchers {
			if matcher(remainingPath, resp, req) {
				return true
			}
		}
		return false
	}
}
