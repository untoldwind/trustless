package remote

import (
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/untoldwind/trustless/secrets"
)

// Client is a trustless client (communicating with a daemon)
type remoteSecrets struct {
	logger     logging.Logger
	httpClient *http.Client
}

// NewRemoteSecrets creates a new remote secrets store. This is the client-side
// counterpart of the daemon api
func NewRemoteSecrets(logger logging.Logger) secrets.Secrets {
	return &remoteSecrets{
		logger: logger.WithField("package", "client"),
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: dialDaemon,
			},
		},
	}
}
