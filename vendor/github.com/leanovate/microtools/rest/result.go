package rest

import (
	"io"
	"net/http"
)

// Result allows a better control of the response of a REST operation.
// This should be used if a REST handler needs to modify HTTP headers or
// has a specific status.
type Result struct {
	Status int
	Header http.Header
	Body   interface{}
}

// Ok creates an OK (200) response.
func Ok() *Result {
	return Status(200)
}

// Created creates a CREATED (201) response.
func Created() *Result {
	return Status(201)
}

// Status creates a response with a specific HTTP status
func Status(status int) *Result {
	return &Result{Status: status, Header: make(http.Header)}
}

// WithStatus modifies the HTTP status of a response
func (r *Result) WithStatus(status int) *Result {
	r.Status = status
	return r
}

func (r *Result) WithBody(body interface{}) *Result {
	r.Body = body
	if r.Body == nil && r.Status == 200 {
		r.Status = 204
	}
	return r
}

func (r *Result) AddHeader(key, value string) *Result {
	r.Header.Add(key, value)
	return r
}

func (r *Result) Send(resp http.ResponseWriter, encoder ResponseEncoder) error {
	for key, values := range r.Header {
		for _, value := range values {
			resp.Header().Add(key, value)
		}
	}
	switch body := r.Body.(type) {
	case nil:
		resp.WriteHeader(r.Status)
		return nil
	case io.ReadCloser:
		resp.WriteHeader(r.Status)
		defer body.Close()
		_, err := io.Copy(resp, body)
		return err
	case io.Reader:
		resp.WriteHeader(r.Status)
		_, err := io.Copy(resp, body)
		return err
	case io.WriterTo:
		resp.WriteHeader(r.Status)
		_, err := body.WriteTo(resp)
		return err
	case []byte:
		resp.WriteHeader(r.Status)
		_, err := resp.Write(body)
		return err
	default:
		if resp.Header().Get("Content-Type") == "" {
			resp.Header().Set("Content-Type", encoder.ContentType())
		}
		resp.WriteHeader(r.Status)
		return encoder.Encode(resp, r.Body)
	}
}
