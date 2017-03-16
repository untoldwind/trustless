package daemon

import (
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/leanovate/microtools/routing"
	"github.com/untoldwind/trustless/secrets"
)

type ServiceDocument struct {
	Links map[string]rest.Link `json:"_links"`
}

type RootResource struct {
	rest.ResourceBase
	logger logging.Logger
	status *StatusResource
}

func NewRootResource(secrets *secrets.Secrets, logger logging.Logger) *RootResource {
	return &RootResource{
		logger: logger.WithField("resource", "service"),
		status: NewStatusResource(secrets),
	}
}

func (RootResource) Self() rest.Link {
	return rest.SimpleLink("/")
}

func (r RootResource) Get(request *http.Request) (interface{}, error) {
	return &ServiceDocument{
		Links: map[string]rest.Link{
			"self":   r.Self(),
			"status": r.status.Self(),
		},
	}, nil
}

func (r *RootResource) SubResources() routing.Matcher {
	return routing.Sequence(
		routing.Prefix("/status", rest.ResourceMatcher(r.status)),
	)
}
