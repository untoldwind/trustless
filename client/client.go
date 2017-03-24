package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

// Client is a trustless client (communicating with a daemon)
type Client struct {
	logger     logging.Logger
	httpClient *http.Client
}

// NewClient creates a new trustless client
func NewClient(logger logging.Logger) *Client {
	return &Client{
		logger: logger.WithField("package", "client"),
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: dialDaemon,
			},
		},
	}
}

// Status gets the current daemon status (and implicitly tests if the daemon is available)
func (c *Client) Status(ctx context.Context) (*api.Status, error) {
	result, err := c.get(ctx, "/status")
	if err != nil {
		return nil, err
	}
	var status api.Status
	if err := json.Unmarshal(result, &status); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize status")
	}
	return &status, nil
}
