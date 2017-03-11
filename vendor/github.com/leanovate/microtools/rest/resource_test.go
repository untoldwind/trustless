package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/leanovate/microtools/rest"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

type minimalResource struct {
	rest.ResourceBase
}

func (minimalResource) Self() rest.Link {
	return rest.SimpleLink("")
}

func TestResource(t *testing.T) {
	Convey("Given a minimal resource", t, func() {
		matcher := rest.ResourceMatcher(minimalResource{})

		Convey("When GET /", func() {
			request, _ := http.NewRequest("GET", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 405)
			So(received.Type, ShouldEndWith, "/405")
		})

		Convey("When GET /something", func() {
			request, _ := http.NewRequest("GET", "/something", nil)
			response := httptest.NewRecorder()
			matches := matcher("/something", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 404)
			So(received.Type, ShouldEndWith, "/404")
		})

		Convey("When PUT /", func() {
			request, _ := http.NewRequest("PUT", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 405)
			So(received.Type, ShouldEndWith, "/405")
		})

		Convey("When PATCH /", func() {
			request, _ := http.NewRequest("PATCH", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 405)
			So(received.Type, ShouldEndWith, "/405")
		})

		Convey("When DELETE /", func() {
			request, _ := http.NewRequest("DELETE", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 405)
			So(received.Type, ShouldEndWith, "/405")
		})
	})

}
