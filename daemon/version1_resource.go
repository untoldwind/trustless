package daemon

import (
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/leanovate/microtools/routing"
	"github.com/untoldwind/trustless/secrets"
)

// Version1Resource is a REST resource containing all resources that are part of
// the v1 API.
type Version1Resource struct {
	rest.ResourceBase
	logger             logging.Logger
	masterKeyResource  *MasterKeyResource
	identitiesResource *IdentitiesResource
	secretsResource    *SecretsResource
	estimateResource   *EstimateResource
}

// NewVersion1Resource creates a new Version1Resource
func NewVersion1Resource(secrets secrets.Secrets, logger logging.Logger) *Version1Resource {
	return &Version1Resource{
		logger:             logger.WithField("resource", "v1"),
		masterKeyResource:  NewMasterKeyResource(secrets, logger),
		identitiesResource: NewIdentitiesResource(secrets, logger),
		secretsResource:    NewSecretsResource(secrets, logger),
		estimateResource:   NewEstimateResource(secrets),
	}
}

// Self link to the resource
func (Version1Resource) Self() rest.Link {
	return rest.SimpleLink("/v1")
}

// Get all hypermedia links the the v1 resources
func (r Version1Resource) Get(request *http.Request) (interface{}, error) {
	return &ServiceDocument{
		Links: map[string]rest.Link{
			"self":       r.Self(),
			"masterkey":  r.masterKeyResource.Self(),
			"identities": r.identitiesResource.Self(),
			"secrets":    r.secretsResource.Self(),
			"estimate":   r.estimateResource.Self(),
		},
	}, nil
}

// SubResources creates routes to all sub-resources
func (r Version1Resource) SubResources() routing.Matcher {
	return routing.Sequence(
		routing.Prefix("/masterkey", rest.ResourceMatcher(r.masterKeyResource)),
		rest.ResourcesMatcher("/identities", r.identitiesResource),
		rest.ResourcesMatcher("/secrets", r.secretsResource),
		routing.Prefix("/estimate", rest.ResourceMatcher(r.estimateResource)),
	)
}
