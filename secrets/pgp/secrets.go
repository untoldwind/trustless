package pgp

import (
	"bytes"
	"sync"
	"time"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/secrets"
	"github.com/untoldwind/trustless/store"
	"golang.org/x/crypto/openpgp"
)

type pgpSecrets struct {
	lock          sync.Mutex
	store         store.Store
	nodeID        string
	logger        logging.Logger
	entities      openpgp.EntityList
	index         *Index
	masterKeyBits int
	autolocker    *secrets.Autolocker
}

// NewPGPSecrets creats a new secrets store based on openpgp
func NewPGPSecrets(storeURL, nodeID string, masterKeyBits int, unlockTimout time.Duration, unlockTimeoutHard bool, logger logging.Logger) (secrets.Secrets, error) {
	store, err := store.NewStore(storeURL, logger)
	if err != nil {
		return nil, err
	}
	pgp := &pgpSecrets{
		store:         store,
		nodeID:        nodeID,
		logger:        logger.WithField("package", "secrets"),
		masterKeyBits: masterKeyBits,
	}
	pgp.autolocker = secrets.NewAutolocker(pgp, unlockTimout, unlockTimeoutHard)
	if err := pgp.readRing(); err != nil {
		return nil, err
	}
	return pgp, nil
}

func (s *pgpSecrets) IsInitialized() bool {
	return len(s.entities) > 0
}

func (s *pgpSecrets) readRing() error {
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
