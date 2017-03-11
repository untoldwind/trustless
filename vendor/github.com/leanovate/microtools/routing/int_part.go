package routing

import (
	"net/http"
	"strconv"
	"strings"
)

func IntPart(matcherForPart func(int) Matcher) Matcher {
	return func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
		if len(remainingPath) == 0 || remainingPath == "/" {
			return false
		}
		startIndex := 0
		if remainingPath[0] == '/' {
			startIndex = 1
		}
		endIndex := strings.IndexRune(remainingPath[startIndex:], '/')
		if endIndex < 0 {
			endIndex = len(remainingPath) - startIndex
		}
		intPart, err := strconv.ParseInt(remainingPath[startIndex:endIndex+startIndex], 10, 32)
		if err != nil {
			return false
		}
		matcher := matcherForPart(int(intPart))
		return matcher(remainingPath[endIndex+startIndex:], resp, req)
	}
}
