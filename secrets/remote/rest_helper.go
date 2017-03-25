package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/leanovate/microtools/rest"
	"github.com/pkg/errors"
)

func (c *remoteSecrets) get(ctx context.Context, uri string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, "http://daemon"+uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Create GET request failed")
	}
	return c.doRequest(req.WithContext(ctx))
}

func (c *remoteSecrets) post(ctx context.Context, uri string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, "http://daemon"+uri, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "Create POST request failed")
	}
	return c.doRequest(req.WithContext(ctx))
}

func (c *remoteSecrets) put(ctx context.Context, uri string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPut, "http://daemon"+uri, bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrap(err, "Create PUT request failed")
	}
	return c.doRequest(req.WithContext(ctx))
}

func (c *remoteSecrets) delete(ctx context.Context, uri string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, "http://daemon"+uri, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Create DELETE request failed")
	}
	return c.doRequest(req.WithContext(ctx))
}

func (c *remoteSecrets) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "perform http request failed")
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read http response body failed")
	}

	if resp.StatusCode >= 300 {
		return nil, c.decodeError(req, resp.StatusCode, data)
	}
	return data, nil
}

func (c *remoteSecrets) decodeError(req *http.Request, status int, data []byte) error {
	var httpError rest.HTTPError

	if err := json.Unmarshal(data, &httpError); err == nil {
		return errors.Wrap(&httpError, "Server request failed")
	}
	return errors.Errorf("Request %s %v failed with %d: %s", req.Method, req.URL, status, string(data))
}
