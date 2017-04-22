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
	query := listURL.Query()
	if filter.Name != "" {
		query.Add("name", filter.Name)
	}
	if filter.Tag != "" {
		query.Add("tag", filter.Tag)
	}
	if filter.Type != "" {
		query.Add("type", string(filter.Type))
	}
	if filter.URL != "" {
		query.Add("url", filter.URL)
	}
	if filter.Deleted {
		query.Add("deleted", "true")
	}
	listURL.RawQuery = query.Encode()

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
