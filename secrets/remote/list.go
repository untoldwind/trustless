package remote

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func (c *remoteSecrets) List(ctx context.Context, filter api.SecretListFilter) (*api.SecretList, error) {
	listURL := &url.URL{
		Path: "/v1/secrets",
	}
	if filter.Name != "" {
		listURL.Query().Add("name", filter.Name)
	}
	if filter.Tag != "" {
		listURL.Query().Add("tag", filter.Tag)
	}
	if filter.Type != "" {
		listURL.Query().Add("type", string(filter.Type))
	}
	if filter.URL != "" {
		listURL.Query().Add("url", filter.URL)
	}

	result, err := c.get(ctx, listURL.String())
	if err != nil {
		return nil, err
	}
	var list api.SecretList
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, errors.Wrap(err, "Failed to deserialize status")
	}
	return &list, nil
}
