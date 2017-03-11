package secrets

import (
	"bytes"
	"crypto"
	"io"

	"github.com/pkg/errors"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

// TODO. This might become obsolete with https://github.com/golang/go/issues/16664
func SerializeKeyRing(entities openpgp.EntityList) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	for _, entity := range entities {
		if err := SerializeEntity(entity, buffer); err != nil {
			return nil, errors.Wrap(err, "Serialize entity failed")
		}
	}
	return buffer.Bytes(), nil
}

func SerializeEntity(e *openpgp.Entity, w io.Writer) (err error) {
	config := &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
	}
	err = e.PrivateKey.Serialize(w)
	if err != nil {
		return
	}
	for _, ident := range e.Identities {
		err = ident.UserId.Serialize(w)
		if err != nil {
			return
		}
		err = ident.SelfSignature.SignUserId(ident.UserId.Id, e.PrimaryKey, e.PrivateKey, config)
		if err != nil {
			return
		}
		err = ident.SelfSignature.Serialize(w)
		if err != nil {
			return
		}
	}
	for _, subkey := range e.Subkeys {
		err = subkey.PrivateKey.Serialize(w)
		if err != nil {
			return
		}
		err = subkey.Sig.SignKey(subkey.PublicKey, e.PrivateKey, config)
		if err != nil {
			return
		}
		err = subkey.Sig.Serialize(w)
		if err != nil {
			return
		}
	}
	return nil
}
