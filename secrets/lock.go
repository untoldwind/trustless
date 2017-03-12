package secrets

import (
	"bytes"
	"crypto"

	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
)

func (s *Secrets) IsLocked() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, entity := range s.entities {
		if entity.PrivateKey != nil && !entity.PrivateKey.Encrypted {
			return false
		}
	}
	return true
}

func (s *Secrets) Lock() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, entity := range s.entities {
		s.purgePrivateKey(entity.PrivateKey)
		for _, subKey := range entity.Subkeys {
			s.purgePrivateKey(subKey.PrivateKey)
		}
	}
}

func (s *Secrets) Unlock(name, email, passphrase string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	ring, err := s.store.GetRing()
	if err != nil {
		return err
	}
	if ring == nil {
		ring, err = s.initializeRing(name, email, passphrase)
		if err != nil {
			return err
		}
	}
	entities, err := openpgp.ReadKeyRing(bytes.NewBuffer(ring))
	if err != nil {
		return errors.Wrap(err, "Failed to read ring")
	}
	s.entities = entities
	for _, entity := range s.entities {
		if err := entity.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
			return errors.Wrap(err, "Unable to decrypt key")
		}
		for _, subKey := range entity.Subkeys {
			if err := subKey.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
				return errors.Wrap(err, "Unable to decrypt key")
			}
		}
	}
	return nil
}

func (s *Secrets) initializeRing(name, email, passphrase string) ([]byte, error) {
	config := &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
		RSABits:       s.MasterKeyBits,
	}
	entity, err := openpgp.NewEntity(name, "", email, config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate entity")
	}
	for _, id := range entity.Identities {
		if err := id.SelfSignature.SignUserId(id.UserId.Id, entity.PrimaryKey, entity.PrivateKey, config); err != nil {
			return nil, errors.Wrap(err, "Failed to sign identity")
		}
	}
	for _, subKey := range entity.Subkeys {
		if err := subKey.Sig.SignKey(subKey.PublicKey, entity.PrivateKey, config); err != nil {
			return nil, errors.Wrap(err, "Failed to sign identity")
		}
	}

	if err := entity.PrivateKey.EncryptWithParameters([]byte(passphrase), packet.CipherAES256, s2k.ModeIterated, s2k.Config{
		S2KCount: 65536,
		Hash:     crypto.SHA512,
	}); err != nil {
		return nil, errors.Wrap(err, "Failed to encrypt key")
	}
	for _, subKey := range entity.Subkeys {
		if err := subKey.PrivateKey.EncryptWithParameters([]byte(passphrase), packet.CipherAES256, s2k.ModeIterated, s2k.Config{
			S2KCount: 65536,
			Hash:     crypto.SHA512,
		}); err != nil {
			return nil, errors.Wrap(err, "Failed to encrypt key")
		}
	}

	buf := bytes.NewBuffer(nil)
	if err := entity.SerializePrivate(buf, config); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize entity")
	}
	ring := buf.Bytes()
	s.purgePrivateKey(entity.PrivateKey)

	if err := s.store.StoreRing(ring); err != nil {
		return nil, err
	}

	return ring, nil
}
