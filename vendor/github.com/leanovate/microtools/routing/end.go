package routing

import "net/http"

func End(subMatcher Matcher) Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		if len(remainingPath) == 0 || remainingPath == "/" {
			return subMatcher("", resp, req)
		}
		return false
	}
}

func EndSeq(subMatchers ...Matcher) Matcher {
	return End(Sequence(subMatchers...))
}
