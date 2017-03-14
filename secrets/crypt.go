package secrets

import (
	"bytes"
	"crypto"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"

	"github.com/pkg/errors"
)

func (s *Secrets) encryptSecret(secretBlock *SecretBlock) ([]byte, error) {
	if s.IsLocked() { // NOTE: Strickly speaking, we can encrypt even if store is locked, future enhancement maybe
		return nil, SecretsLockedError
	}
	content, err := json.Marshal(secretBlock)
	if err != nil {
		return nil, errors.Wrap(err, "Json marshal of secret block failed")
	}
	block, err := padBlock(content)
	if err != nil {
		return nil, err
	}
	out := bytes.NewBuffer(nil)
	writer, err := openpgp.Encrypt(out, s.entities, nil, nil, &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Openpgp create writer failed")
	}
	if _, err := writer.Write(block); err != nil {
		return nil, errors.Wrap(err, "Openpgp encryption failed")
	}
	if err := writer.Close(); err != nil {
		return nil, errors.Wrap(err, "Openpgp encryption failed")
	}
	return out.Bytes(), nil
}

func (s *Secrets) decryptSecret(encrypted []byte) (*SecretBlock, error) {
	if s.IsLocked() {
		return nil, SecretsLockedError
	}

	message, err := openpgp.ReadMessage(bytes.NewBuffer(encrypted), s.entities, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Read openpgp message failed")
	}
	block, err := ioutil.ReadAll(message.UnverifiedBody)
	if err != nil {
		return nil, errors.Wrap(err, "Openpgp decrypt failed")
	}
	content, err := unpadBlock(block)
	if err != nil {
		return nil, err
	}

	var secert SecretBlock
	if err := json.Unmarshal(content, &secert); err != nil {
		return nil, errors.Wrap(err, "Json unmarshal failed")
	}
	return &secert, nil
}
