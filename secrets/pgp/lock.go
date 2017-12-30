package pgp

import (
	"bytes"
	"context"
	"crypto"

	"github.com/awnumar/memguard"

	"github.com/untoldwind/scrypt-go/scryptlib"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/config"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
)

func (s *pgpSecrets) Status(ctx context.Context) (*api.Status, error) {
	if s.isLocked() {
		return &api.Status{
			Initialized: len(s.identities) > 0,
			Locked:      true,
			Version:     config.Version(),
		}, nil
	}
	autolockAt := s.autolocker.GetAutolockAt()

	return &api.Status{
		Initialized: true,
		Locked:      false,
		AutolockAt:  &autolockAt,
		Version:     config.Version(),
	}, nil
}

func (s *pgpSecrets) isLocked() bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, entity := range s.entities {
		if entity.PrivateKey != nil && !entity.PrivateKey.Encrypted {
			return false
		}
	}
	return true
}

func (s *pgpSecrets) Lock(ctx context.Context) error {
	s.logger.Info("Locking secrets")

	s.lock.Lock()
	defer s.lock.Unlock()

	s.preparePurge()
	for _, entity := range s.entities {
		s.purgePrivateKey(entity.PrivateKey)
		for _, subKey := range entity.Subkeys {
			s.purgePrivateKey(subKey.PrivateKey)
		}
	}
	s.destroyBuffers()
	s.autolocker.Cancel()
	return nil
}

func (s *pgpSecrets) Unlock(ctx context.Context, name, email, passphrase string) error {
	s.logger.Info("Unlocking secrets")
	s.lock.Lock()
	defer s.lock.Unlock()

	rawRing, err := s.store.GetRing()
	if err != nil {
		return err
	}
	if rawRing == nil {
		rawRing, err = s.initializeRing(name, email, passphrase)
		if err != nil {
			return err
		}
	}

	ring := rawRing
	if s.scrypted {
		out := bytes.NewBuffer(nil)
		if err := scryptlib.Decrypt([]byte(passphrase), bytes.NewBuffer(rawRing), out); err != nil {
			return err
		}
		ring = out.Bytes()
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
		if err := s.protectPrivateKey(entity.PrivateKey); err != nil {
			return err
		}
		for _, subKey := range entity.Subkeys {
			if err := subKey.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
				return errors.Wrap(err, "Unable to decrypt key")
			}
			if err := s.protectPrivateKey(subKey.PrivateKey); err != nil {
				return err
			}
		}
	}
	s.fetchIndex()
	s.autolocker.Start()
	return nil
}

func (s *pgpSecrets) initializeRing(name, email, passphrase string) ([]byte, error) {
	sha256Id, ok := s2k.HashToHashId(crypto.SHA256)
	if !ok {
		return nil, errors.New("SHA256 id not found")
	}
	sha512Id, ok := s2k.HashToHashId(crypto.SHA512)
	if !ok {
		return nil, errors.New("SHA512 id not found")
	}

	config := &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
		RSABits:       s.masterKeyBits,
	}
	entity, err := openpgp.NewEntity(name, "", email, config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate entity")
	}
	for _, id := range entity.Identities {
		id.SelfSignature.PreferredSymmetric = []uint8{uint8(packet.CipherAES256)}
		id.SelfSignature.PreferredHash = []uint8{sha512Id, sha256Id}

		if err := id.SelfSignature.SignUserId(id.UserId.Id, entity.PrimaryKey, entity.PrivateKey, config); err != nil {
			return nil, errors.Wrap(err, "Failed to sign identity")
		}
	}
	for _, subKey := range entity.Subkeys {
		subKey.Sig.PreferredSymmetric = []uint8{uint8(packet.CipherAES256)}
		subKey.Sig.PreferredHash = []uint8{sha512Id, sha256Id}

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
	publicBuf := bytes.NewBuffer(nil)
	if err := entity.Serialize(publicBuf); err != nil {
		return nil, errors.Wrap(err, "Failed to serialize entity")
	}
	ring := buf.Bytes()
	rawRing := ring

	if s.scrypted {
		buf, err := memguard.NewMutable(len(ring))
		if err != nil {
			return nil, err
		}
		defer buf.Destroy()
		out := &writeTo{buffer: buf}
		if err := scryptlib.Encrypt([]byte(passphrase), bytes.NewBuffer(ring), out); err != nil {
			return nil, err
		}
		rawRing = out.Result()
	}

	s.purgePrivateKey(entity.PrivateKey)

	if err := s.store.StorePublicRing(publicBuf.Bytes()); err != nil {
		return nil, err
	}
	if err := s.store.StoreRing(rawRing); err != nil {
		return nil, err
	}

	return rawRing, nil
}
