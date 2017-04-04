package routing_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/routing"
	. "github.com/smartystreets/goconvey/convey"
)

type mockHandler struct {
	request *http.Request
	status  int
}

func (h *mockHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.request = req
	resp.WriteHeader(h.status)
}

func TestLoggingHander(t *testing.T) {
	Convey("Given a logging handler", t, func() {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		buffer := bytes.NewBuffer([]byte{})
		logger := logging.NewSimpleLogger(logging.Options{Output: buffer, Level: logging.Debug})
		handler := &mockHandler{
			status: 200,
		}
		loggingHandler := routing.NewLoggingHandler(handler, logger)

		Convey("When GET /therequest is send with OK result", func() {
			request, err := http.NewRequest("GET", "/therequest", nil)
			request.RequestURI = "/therequest"
			request.Header.Add("X-Flow-Id", "flow")

			So(err, ShouldBeNil)

			recorder := httptest.NewRecorder()

			loggingHandler.ServeHTTP(recorder, request)

			So(recorder.Code, ShouldEqual, 200)
			So(buffer.String(), ShouldContainSubstring, "uri=/therequest")
			So(buffer.String(), ShouldContainSubstring, "status=200")
			So(buffer.String(), ShouldContainSubstring, "Request: Success")
			So(buffer.String(), ShouldContainSubstring, "flow_id=flow")
		})

		Convey("When GET /therequest is send with SEE OTHER result", func() {
			handler.status = 303
			request, err := http.NewRequest("GET", "/therequest", nil)
			request.RequestURI = "/therequest"

			So(err, ShouldBeNil)

			recorder := httptest.NewRecorder()

			loggingHandler.ServeHTTP(recorder, request)

			So(recorder.Code, ShouldEqual, 303)
			So(buffer.String(), ShouldContainSubstring, "uri=/therequest")
			So(buffer.String(), ShouldContainSubstring, "status=303")
			So(buffer.String(), ShouldContainSubstring, "Request: Redirect")
		})

		Convey("When GET /therequest is send with NOT FOUND result", func() {
			handler.status = 404
			request, err := http.NewRequest("GET", "/therequest", nil)
			request.RequestURI = "/therequest"

			So(err, ShouldBeNil)

			recorder := httptest.NewRecorder()

			loggingHandler.ServeHTTP(recorder, request)

			So(recorder.Code, ShouldEqual, 404)
			So(buffer.String(), ShouldContainSubstring, "uri=/therequest")
			So(buffer.String(), ShouldContainSubstring, "status=404")
			So(buffer.String(), ShouldContainSubstring, "Request: Client error")
		})
	})
}
