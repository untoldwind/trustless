package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/routing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIntPart(t *testing.T) {
	Convey("Given an int port matcher", t, func() {
		var subMatcherPath string
		var subMatcherCalled bool
		var subMatcherMatches bool

		subMatcher := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			subMatcherPath = remainingPath
			subMatcherCalled = true
			return subMatcherMatches
		}

		var extractedPart int
		matcher := routing.IntPart(func(part int) routing.Matcher {
			extractedPart = part
			return subMatcher
		})

		request, _ := http.NewRequest("GET", "/matches", nil)
		recorder := httptest.NewRecorder()

		Convey("And sub-matcher always matches", func() {
			subMatcherMatches = true

			Convey("When http request with non-empty remaining path is matched", func() {
				result := matcher("12345/and/more", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/and/more")
				So(extractedPart, ShouldEqual, 12345)
			})

			Convey("When http request with beging / in remaining path is matches", func() {
				result := matcher("/12345/and/more", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/and/more")
				So(extractedPart, ShouldEqual, 12345)
			})

			Convey("When http request with trailing string part", func() {
				result := matcher("12345", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
				So(extractedPart, ShouldEqual, 12345)
			})

			Convey("When http request with trailing string part beging with /", func() {
				result := matcher("/12345", recorder, request)

				So(result, ShouldBeTrue)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "")
				So(extractedPart, ShouldEqual, 12345)
			})

			Convey("When http request with non-number remaining path is matched", func() {
				result := matcher("/abcdef", recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeFalse)
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
				result := matcher("12345/and/more", recorder, request)

				So(result, ShouldBeFalse)
				So(subMatcherCalled, ShouldBeTrue)
				So(subMatcherPath, ShouldEqual, "/and/more")
				So(extractedPart, ShouldEqual, 12345)
			})
		})
	})
}
