package routing

import (
	"net/http"
	"net/url"
	"strings"
)

func StringPart(matcherForPart func(string) Matcher) Matcher {
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
		strPart, err := url.QueryUnescape(remainingPath[startIndex : endIndex+startIndex])
		if err != nil {
			return false
		}
		matcher := matcherForPart(strPart)
		return matcher(remainingPath[endIndex+startIndex:], resp, req)
	}
}
