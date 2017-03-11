package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/routing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStringPart(t *testing.T) {
	Convey("Given a string part matcher", t, func() {
		var subMatcherPath string
		var subMatcherCalled bool
		var subMatcherMatches bool

		subMatcher := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			subMatcherPath = remainingPath
			subMatcherCalled = true
			return subMatcherMatches
		}

		var extractedPart string
		matcher := routing.StringPart(func(part string) routing.Matcher {
			extractedPart = part
			return subMatcher
		})

		request, _ := http.NewRequest("GET", "/matches", nil)
		recorder := httptest.NewRecorder()

		Convey("And sub-matcher always matches", func() {
			subMatcherMatches = true

			Convey("When http request with non-empty remaining path is matched", func() {
				result := matcher("something/and/more", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/and/more")
				So(extractedPart, ShouldEqual, "something")
			})

			Convey("When http request with beging / in remaining path is matches", func() {
				result := matcher("/something/and/more", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/and/more")
				So(extractedPart, ShouldEqual, "something")
			})

			Convey("When http request with trailing string part", func() {
				result := matcher("something", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
				So(extractedPart, ShouldEqual, "something")
			})

			Convey("When http request with trailing string part beging with /", func() {
				result := matcher("/something", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
				So(extractedPart, ShouldEqual, "something")
			})

			Convey("When http request with empty remaining path is matched", func() {
				result := matcher("", recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeFalse)
			})
		})

		Convey("And sub-matcher never matches", func() {
			subMatcherMatches = false

			Convey("When http request with non-empty remaining path is matched", func() {
				result := matcher("something/and/more", recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/and/more")
				So(extractedPart, ShouldEqual, "something")
			})
		})
	})
}
