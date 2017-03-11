package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/routing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEnd(t *testing.T) {
	Convey("Given an end matcher", t, func() {
		var subMatcherPath string
		var subMatcherCalled bool
		var subMatcherMatches bool

		subMatcher := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			subMatcherPath = remainingPath
			subMatcherCalled = true
			return subMatcherMatches
		}
		matcher := routing.End(subMatcher)

		Convey("And sub-matcher always matches", func() {
			subMatcherMatches = true

			Convey("When http request with empty remaining is matchess", func() {
				request, _ := http.NewRequest("GET", "/matches", nil)
				recorder := httptest.NewRecorder()

				result := matcher("", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
			})

			Convey("When http request with trailing / is matchess", func() {
				request, _ := http.NewRequest("GET", "/matches", nil)
				recorder := httptest.NewRecorder()

				result := matcher("/", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
			})

			Convey("When http request with raminging patch", func() {
				request, _ := http.NewRequest("GET", "/notmatch/some_more", nil)
				recorder := httptest.NewRecorder()

				result := matcher(request.URL.Path, recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeFalse)
			})
		})

		Convey("And sub-matcher never matches", func() {
			subMatcherMatches = false

			Convey("When http request with empty remaining is matchess", func() {
				request, _ := http.NewRequest("GET", "/matches", nil)
				recorder := httptest.NewRecorder()

				result := matcher("", recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
			})

			Convey("When http request with trailing / is matchess", func() {
				request, _ := http.NewRequest("GET", "/matches", nil)
				recorder := httptest.NewRecorder()

				result := matcher("/", recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
			})

			Convey("When http request with raminging patch", func() {
				request, _ := http.NewRequest("GET", "/notmatch/some_more", nil)
				recorder := httptest.NewRecorder()

				result := matcher(request.URL.Path, recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeFalse)
			})
		})
	})
}
