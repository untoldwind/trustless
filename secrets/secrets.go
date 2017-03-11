package secrets

import (
	"bytes"
	"crypto"
	"fmt"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store"
)

type Secrets struct {
	store   store.Store
	logger  logging.Logger
	entries openpgp.EntityList
}

func NewSecrets(storeURL string, logger logging.Logger) (*Secrets, error) {
	store, err := store.NewStore(storeURL, logger)
	if err != nil {
		return nil, err
	}
	return &Secrets{
		store:  store,
		logger: logger.WithField("package", "secrets"),
	}, nil
}

func (s *Secrets) Unlock(name, email, passphrase string) error {
	ring, err := s.store.GetRing()
	if err != nil {
		return err
	}
	if ring == nil {
		return s.initializeRing(name, email, passphrase)
	}
	entries, err := openpgp.ReadKeyRing(bytes.NewBuffer(ring))
	if err != nil {
		return errors.Wrap(err, "Failed to read ring")
	}
	s.entries = entries
	return nil
}

func (s *Secrets) initializeRing(name, email, passphrase string) error {
	entity, _ := openpgp.NewEntity(name, "", email, &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
		RSABits:       4096,
	})
	fmt.Printf("%#v\n", entity)
	return nil
}
