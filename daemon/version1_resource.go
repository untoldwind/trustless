package daemon

import (
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
)

type Version1Resource struct {
	rest.ResourceBase
	logger logging.Logger
}

func NewVersion1Resource(logger logging.Logger) *Version1Resource {
	return &Version1Resource{
		logger: logger.WithField("resource", "v1"),
	}
}

func (Version1Resource) Self() rest.Link {
	return rest.SimpleLink("/v1")
}

func (r Version1Resource) Get(request *http.Request) (interface{}, error) {
	return &ServiceDocument{
		Links: map[string]rest.Link{
			"self": r.Self(),
		},
	}, nil
}
