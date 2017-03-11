package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/routing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPrefix(t *testing.T) {
	Convey("Given a prefix matcher", t, func() {
		var subMatcherPath string
		var subMatcherCalled bool
		var subMatcherMatches bool

		subMatcher := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			subMatcherPath = remainingPath
			subMatcherCalled = true
			return subMatcherMatches
		}
		matcher := routing.Prefix("/matches", subMatcher)

		Convey("And sub-matcher always matches", func() {
			subMatcherMatches = true

			Convey("When http request with correct prefix is matches", func() {
				request, _ := http.NewRequest("GET", "/matches/some_more", nil)
				recorder := httptest.NewRecorder()

				result := matcher(request.URL.Path, recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/some_more")
			})

			Convey("When http request with incorrect prefix is matches", func() {
				request, _ := http.NewRequest("GET", "/notmatch/some_more", nil)
				recorder := httptest.NewRecorder()

				result := matcher(request.URL.Path, recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeFalse)
			})

			Convey("When http request with to short prefix is matches", func() {
				request, _ := http.NewRequest("GET", "/tos", nil)
				recorder := httptest.NewRecorder()

				result := matcher(request.URL.Path, recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeFalse)
			})
		})

		Convey("And sub-matcher never matches", func() {
			subMatcherMatches = false

			Convey("When http request with correct prefix is matches", func() {
				request, _ := http.NewRequest("GET", "/matches/some_more", nil)
				recorder := httptest.NewRecorder()

				result := matcher(request.URL.Path, recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/some_more")
			})

			Convey("When http request with incorrect prefix is matches", func() {
				request, _ := http.NewRequest("GET", "/notmatch/some_more", nil)
				recorder := httptest.NewRecorder()

				result := matcher(request.URL.Path, recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeFalse)
			})

		})
	})
}
