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
	// GetRing retrieves the key ring of the store
	GetRing() ([]byte, error)
	// StoreRing stores the key ring of the store
	StoreRing(raw []byte) error

	// ChangeLogs retrieves the change logs of all nodes
	ChangeLogs() ([]model.ChangeLog, error)

	// GetIndex retrieves the index block of a node
	GetIndex(nodeID string) ([]byte, error)
	// StoreIndex stores the index block of a node
	StoreIndex(nodeID string, indexBlock []byte) error

	// AddBlock adds a block (of encrypted data) to the store and
	// return its id
	AddBlock(block []byte) (string, error)
	// GetBlock retrieves a block by its id
	GetBlock(blockID string) ([]byte, error)

	// Commit changes made to the store (i.e. write them the the change log)
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
