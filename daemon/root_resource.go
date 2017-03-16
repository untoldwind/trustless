package daemon

import (
	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
)

type RootResource struct {
	rest.ResourceBase
	logger logging.Logger
}

func NewRootResource(logger logging.Logger) *RootResource {
	return &RootResource{
		logger: logger.WithField("resource", "service"),
	}
}

func (RootResource) Self() rest.Link {
	return rest.SimpleLink("/")
}
