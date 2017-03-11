package routing

import "net/http"

func Prefix(prefix string, subMatcher Matcher) Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		if len(remainingPath) >= len(prefix) && remainingPath[:len(prefix)] == prefix {
			return subMatcher(remainingPath[len(prefix):], resp, req)
		}
		return false
	}
}

func PrefixSeq(prefix string, subMatchers ...Matcher) Matcher {
	if prefix == "" {
		return Sequence(subMatchers...)
	}
	return Prefix(prefix, Sequence(subMatchers...))
}
