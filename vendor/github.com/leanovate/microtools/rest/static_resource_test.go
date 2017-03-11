package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/rest"
	. "github.com/smartystreets/goconvey/convey"
)

func TestStaticResource(t *testing.T) {
	Convey("Given a static resource", t, func() {
		staticResource := rest.StaticContent([]byte("Hello world"), "text/plain", "/")
		matcher := rest.ResourceMatcher(staticResource)

		Convey("Then self link should be passed through", func() {
			So(staticResource.Self().Href, ShouldEqual, "/")
		})

		Convey("When GET /", func() {
			request, _ := http.NewRequest("GET", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			So(response.Code, ShouldEqual, 200)
			So(response.Body.String(), ShouldEqual, "Hello world")
			So(response.Header().Get("Content-Type"), ShouldEqual, "text/plain")
		})
		Convey("When POST /", func() {
			request, _ := http.NewRequest("POST", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			So(response.Code, ShouldEqual, 405)
		})

		Convey("When PATCH /", func() {
			request, _ := http.NewRequest("PATCH", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			So(response.Code, ShouldEqual, 405)
		})

		Convey("When PUT /", func() {
			request, _ := http.NewRequest("PUT", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			So(response.Code, ShouldEqual, 405)
		})

		Convey("When DELETE /", func() {
			request, _ := http.NewRequest("DELETE", "/", nil)
			response := httptest.NewRecorder()
			matches := matcher("/", response, request)

			So(matches, ShouldBeTrue)

			So(response.Code, ShouldEqual, 405)
		})
	})

}
