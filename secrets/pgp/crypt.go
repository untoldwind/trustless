package pgp

import (
	"bytes"
	"crypto"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/secrets"
)

func (s *pgpSecrets) encryptData(content []byte) ([]byte, error) {
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

func (s *pgpSecrets) encryptSecret(secretBlock *secrets.SecretBlock) ([]byte, error) {
	content, err := json.Marshal(secretBlock)
	if err != nil {
		return nil, errors.Wrap(err, "Json marshal of secret block failed")
	}
	return s.encryptData(content)
}

func (s *pgpSecrets) decryptData(encrypted []byte) ([]byte, error) {
	message, err := openpgp.ReadMessage(bytes.NewBuffer(encrypted), s.entities, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Read openpgp message failed")
	}
	block, err := ioutil.ReadAll(message.UnverifiedBody)
	if err != nil {
		return nil, errors.Wrap(err, "Openpgp decrypt failed")
	}
	return unpadBlock(block)
}

func (s *pgpSecrets) decryptSecret(encrypted []byte) (*secrets.SecretBlock, error) {
	content, err := s.decryptData(encrypted)
	if err != nil {
		return nil, err
	}
	var secert secrets.SecretBlock
	if err := json.Unmarshal(content, &secert); err != nil {
		return nil, errors.Wrap(err, "Json unmarshal failed")
	}
	return &secert, nil
}
