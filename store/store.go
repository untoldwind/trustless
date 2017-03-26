package store

import (
	"net/url"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store/localdir"
	"github.com/untoldwind/trustless/store/model"
)

// Store is the common interface for a backend store.
// Secrets are stored in encrypted blocks, its the Store's responsibility to
// do the necessary I/O stuff.
type Store interface {
	GetRing() ([]byte, error)
	StoreRing(raw []byte) error

	ChangeLogs() ([]model.ChangeLog, error)

	AddBlock(block []byte) (string, error)
	GetBlock(blockID string) ([]byte, error)

	Commit(nodeID string, changes []model.Change) error
}

// NewStore creates a new backend store from a URL.
func NewStore(storeURLStr string, logger logging.Logger) (Store, error) {
	storeURL, err := url.Parse(storeURLStr)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid StoreURL")
	}
	switch storeURL.Scheme {
	case "file":
		return localdir.NewLocaldirStore(storeURL, logger)
	}
	return nil, errors.Errorf("Invalid store url: %s", storeURLStr)
}
