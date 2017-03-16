package secrets

import (
	"bytes"
	"sync"

	"golang.org/x/crypto/openpgp"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store"
)

type Secrets struct {
	lock          sync.Mutex
	store         store.Store
	nodeID        string
	logger        logging.Logger
	entities      openpgp.EntityList
	index         *Index
	MasterKeyBits int
}

func NewSecrets(storeURL, nodeID string, logger logging.Logger) (*Secrets, error) {
	store, err := store.NewStore(storeURL, logger)
	if err != nil {
		return nil, err
	}
	secrets := &Secrets{
		store:         store,
		nodeID:        nodeID,
		logger:        logger.WithField("package", "secrets"),
		MasterKeyBits: 4096,
	}
	if err := secrets.readRing(); err != nil {
		return nil, err
	}
	return secrets, nil
}

func (s *Secrets) IsInitialized() bool {
	return len(s.entities) > 0
}

func (s *Secrets) readRing() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	ring, err := s.store.GetRing()
	if err != nil {
		return err
	}
	if ring == nil {
		return nil
	}
	s.entities, err = openpgp.ReadKeyRing(bytes.NewBuffer(ring))
	if err != nil {
		return errors.Wrap(err, "Failed to read ring")
	}
	return nil
}
