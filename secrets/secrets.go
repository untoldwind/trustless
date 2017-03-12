package secrets

import (
	"sync"

	"golang.org/x/crypto/openpgp"

	"github.com/leanovate/microtools/logging"
	"github.com/untoldwind/trustless/store"
)

type Secrets struct {
	lock          sync.Mutex
	store         store.Store
	logger        logging.Logger
	entities      openpgp.EntityList
	MasterKeyBits int
}

func NewSecrets(storeURL string, logger logging.Logger) (*Secrets, error) {
	store, err := store.NewStore(storeURL, logger)
	if err != nil {
		return nil, err
	}
	return &Secrets{
		store:         store,
		logger:        logger.WithField("package", "secrets"),
		MasterKeyBits: 4096,
	}, nil
}
