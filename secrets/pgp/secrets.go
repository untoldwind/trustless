package pgp

import (
	"bytes"
	"sync"
	"time"

	"github.com/awnumar/memguard"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
	"github.com/untoldwind/trustless/store"
	"golang.org/x/crypto/openpgp"
)

type pgpSecrets struct {
	lock          sync.Mutex
	store         store.Store
	scrypted      bool
	nodeID        string
	logger        logging.Logger
	identities    []api.Identity
	entities      openpgp.EntityList
	index         *Index
	masterKeyBits int
	autolocker    *secrets.Autolocker
	buffers       []*memguard.LockedBuffer
}

// NewPGPSecrets creats a new secrets store based on openpgp
func NewPGPSecrets(storeURL string, scrypted bool, nodeID string, masterKeyBits int, unlockTimout time.Duration, unlockTimeoutHard bool, logger logging.Logger) (secrets.Secrets, error) {
	store, err := store.NewStore(storeURL, logger)
	if err != nil {
		return nil, err
	}
	pgp := &pgpSecrets{
		store:         store,
		scrypted:      scrypted,
		nodeID:        nodeID,
		logger:        logger.WithField("package", "secrets"),
		masterKeyBits: masterKeyBits,
	}
	pgp.autolocker = secrets.NewAutolocker(pgp, unlockTimout, unlockTimeoutHard)
	if err := pgp.readIdentities(); err != nil {
		return nil, err
	}
	return pgp, nil
}

func (s *pgpSecrets) IsInitialized() bool {
	return len(s.identities) > 0
}

func (s *pgpSecrets) readIdentities() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	ring, err := s.store.GetPublicRing()
	if err != nil {
		return nil
	}
	if ring != nil {
		entities, err := openpgp.ReadKeyRing(bytes.NewBuffer(ring))
		if err != nil {
			return errors.Wrap(err, "Failed to read ring")
		}
		s.identities = identitiesFromEntities(entities)
		return nil
	}
	ring, err = s.store.GetRing()
	if err != nil {
		return err
	}
	if ring == nil {
		return nil
	}
	entities, err := openpgp.ReadKeyRing(bytes.NewBuffer(ring))
	if err != nil {
		return errors.Wrap(err, "Failed to read ring")
	}
	s.identities = identitiesFromEntities(entities)
	buf := bytes.NewBuffer(nil)

	for _, entity := range entities {
		if err := entity.Serialize(buf); err != nil {
			return err
		}
	}
	if err := s.store.StorePublicRing(buf.Bytes()); err != nil {
		return err
	}
	return nil
}
