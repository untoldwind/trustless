package rest

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type ResponseEncoder interface {
	ContentType() string
	Encode(output io.Writer, data interface{}) error
}

type jsonResponseEncoder struct{}

var JsonResponseEncoder jsonResponseEncoder

func (jsonResponseEncoder) ContentType() string {
	return "application/json"
}

func (jsonResponseEncoder) Encode(output io.Writer, data interface{}) error {
	return json.NewEncoder(output).Encode(data)
}

type xmlResponseEncoder struct{}

var XmlResponseEncoder xmlResponseEncoder

func (xmlResponseEncoder) ContentType() string {
	return "text/xml"
}

func (xmlResponseEncoder) Encode(output io.Writer, data interface{}) error {
	return xml.NewEncoder(output).Encode(data)
}

type ResponseEncoderChooser func(*http.Request) ResponseEncoder

func StdResponseEncoderChooser(request *http.Request) ResponseEncoder {
	accept := request.Header.Get("accept")

	if strings.Contains(accept, "text/xml") || strings.Contains(accept, "application/xml") {
		return XmlResponseEncoder
	}
	return JsonResponseEncoder
}
