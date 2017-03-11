package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/routing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSequence(t *testing.T) {
	Convey("Given a sequence of matchers", t, func() {
		var matcher1Matches bool
		var matcher1Called bool
		matcher1 := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			matcher1Called = true
			return matcher1Matches
		}
		var matcher2Matches bool
		var matcher2Called bool
		matcher2 := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			matcher2Called = true
			return matcher2Matches
		}
		var matcher3Matches bool
		var matcher3Called bool
		matcher3 := func(remainingPath string, resp http.ResponseWriter, req *http.Request) bool {
			matcher3Called = true
			return matcher3Matches
		}

		matcher := routing.Sequence(matcher1, matcher2, matcher3)

		request, _ := http.NewRequest("GET", "/something", nil)
		recorder := httptest.NewRecorder()

		Convey("When matcher1 matches", func() {
			matcher1Matches = true

			result := matcher(request.URL.Path, recorder, request)

			So(result, ShouldBeTrue)
			So(matcher1Called, ShouldBeTrue)
			So(matcher2Called, ShouldBeFalse)
			So(matcher3Called, ShouldBeFalse)
		})

		Convey("WHen matcher2 matches", func() {
			matcher2Matches = true

			result := matcher(request.URL.Path, recorder, request)

			So(result, ShouldBeTrue)
			So(matcher1Called, ShouldBeTrue)
			So(matcher2Called, ShouldBeTrue)
			So(matcher3Called, ShouldBeFalse)
		})

		Convey("WHen matcher3 matches", func() {
			matcher3Matches = true

			result := matcher(request.URL.Path, recorder, request)

			So(result, ShouldBeTrue)
			So(matcher1Called, ShouldBeTrue)
			So(matcher2Called, ShouldBeTrue)
			So(matcher3Called, ShouldBeTrue)
		})

		Convey("WHen no matcher matches", func() {
			result := matcher(request.URL.Path, recorder, request)

			So(result, ShouldBeFalse)
			So(matcher1Called, ShouldBeTrue)
			So(matcher2Called, ShouldBeTrue)
			So(matcher3Called, ShouldBeTrue)
		})
	})
}
