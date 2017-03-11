package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/routing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHandler(t *testing.T) {
	Convey("Given a route matcher", t, func() {
		matcher := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			return req.URL.Path == "/matches"
		}

		handler := routing.RouteHandlerFunc(matcher)

		Convey("When a http request matches", func() {
			request, _ := http.NewRequest("GET", "/matches", nil)
			recorder := httptest.NewRecorder()

			handler(recorder, request)

			So(recorder.Code, ShouldEqual, 200)
		})

		Convey("When a http does not match", func() {
			request, _ := http.NewRequest("GET", "/notmatch", nil)
			recorder := httptest.NewRecorder()

			handler(recorder, request)

			So(recorder.Code, ShouldEqual, 404)
		})
	})
}
