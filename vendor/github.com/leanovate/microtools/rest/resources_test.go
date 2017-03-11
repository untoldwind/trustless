package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/leanovate/microtools/rest"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

type minimalResources struct {
	rest.ResourcesBase
}

func (minimalResources) Self() rest.Link {
	return rest.SimpleLink("")
}

func TestResources(t *testing.T) {
	Convey("Given base resources", t, func() {
		matcher := rest.ResourcesMatcher("/v1", minimalResources{})

		Convey("When GET /v1", func() {
			request, _ := http.NewRequest("GET", "/v1", nil)
			response := httptest.NewRecorder()
			matches := matcher("/v1", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 405)
			So(received.Type, ShouldEndWith, "/405")
		})

		Convey("When GET /v1/something", func() {
			request, _ := http.NewRequest("GET", "/v1/something", nil)
			response := httptest.NewRecorder()
			matches := matcher("/v1/something", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 404)
			So(received.Type, ShouldEndWith, "/404")
		})

		Convey("When POST /v1", func() {
			request, _ := http.NewRequest("POST", "/v1", nil)
			response := httptest.NewRecorder()
			matches := matcher("/v1", response, request)

			So(matches, ShouldBeTrue)

			var received rest.HTTPError
			err := json.Unmarshal(response.Body.Bytes(), &received)

			So(err, ShouldBeNil)
			So(received.Code, ShouldEqual, 405)
			So(received.Type, ShouldEndWith, "/405")
		})
	})
}
