package rest_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/leanovate/microtools/rest"
	. "github.com/smartystreets/goconvey/convey"
)

type writableString string

func (s writableString) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return int64(n), nil
}

func TestResult(t *testing.T) {
	Convey("Given an OK result", t, func() {
		result := rest.Ok()

		Convey("When result is send", func() {
			response := httptest.NewRecorder()

			result.Send(response, rest.JsonResponseEncoder)

			So(response.Code, ShouldEqual, 200)
		})

		Convey("When the status is changed to ACCEPTED", func() {
			result = result.WithStatus(202)
			response := httptest.NewRecorder()

			result.Send(response, rest.JsonResponseEncoder)

			So(response.Code, ShouldEqual, 202)
		})

		Convey("When the body is set to nil", func() {
			result = result.WithBody(nil)
			response := httptest.NewRecorder()

			result.Send(response, rest.JsonResponseEncoder)

			So(response.Code, ShouldEqual, 204)
		})

		Convey("When the body is set to json encodable", func() {
			body := &struct {
				Data string `json:"data"`
			}{
				Data: "Data content",
			}
			result = result.WithBody(body)
			response := httptest.NewRecorder()

			result.Send(response, rest.JsonResponseEncoder)

			So(response.Code, ShouldEqual, 200)
			So(response.Header().Get("Content-Type"), ShouldEqual, "application/json")
			So(response.Body.String(), ShouldEqual, "{\"data\":\"Data content\"}\n")
		})

		Convey("When the body is set to a reader", func() {
			body := bytes.NewBufferString("Data content")
			result = result.WithBody(body)
			response := httptest.NewRecorder()

			result.Send(response, rest.JsonResponseEncoder)

			So(response.Code, ShouldEqual, 200)
			So(response.Body.String(), ShouldEqual, "Data content")
		})

		Convey("When the body is set to a WriterTo", func() {
			result = result.WithBody(writableString("Data content"))
			response := httptest.NewRecorder()

			result.Send(response, rest.JsonResponseEncoder)

			So(response.Code, ShouldEqual, 200)
			So(response.Body.String(), ShouldEqual, "Data content")
		})
	})

	Convey("Given an CREATED result", t, func() {
		result := rest.Created()

		Convey("When result is send", func() {
			response := httptest.NewRecorder()

			result.Send(response, rest.JsonResponseEncoder)

			So(response.Code, ShouldEqual, 201)
		})
	})
}
