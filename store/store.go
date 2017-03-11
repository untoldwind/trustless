package store

import (
	"net/url"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store/localdir"
	"github.com/untoldwind/trustless/store/model"
)

type Store interface {
	GetRing(ringType model.RingType) ([]byte, error)
	StoreRing(ringType model.RingType, raw []byte) error

	Heads() ([]model.Head, error)
	GetHead(nodeID string) (string, error)

	AddBlock(block []byte) (string, error)
	GetBlock(blockID string) ([]byte, error)

	Commit(nodeID string, changes []model.Change) (string, error)
	GetCommit(commitID string) (*model.Commit, error)
}

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
