package secrets_test

import (
	"bytes"
	"crypto"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/secrets"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

func TestSerializeKeyRing(t *testing.T) {
	require := require.New(t)

	expected1, err := openpgp.NewEntity("Tester", "", "tester@mail.org", &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
		RSABits:       1024,
	})
	require.Nil(err)
	expected2, err := openpgp.NewEntity("Tester", "", "tester@mail.org", &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
		RSABits:       1024,
	})
	require.Nil(err)
	encoded, err := secrets.SerializeKeyRing(openpgp.EntityList{expected1, expected2})
	require.Nil(err)

	actual, err := openpgp.ReadKeyRing(bytes.NewBuffer(encoded))
	require.Nil(err)
	require.Len(actual, 2)

	require.Equal(expected1.PrimaryKey.PublicKey, actual[0].PrimaryKey.PublicKey)
	require.Equal(expected1.PrivateKey.PublicKey.PublicKey, actual[0].PrivateKey.PublicKey.PublicKey)
	require.Equal(expected2.PrimaryKey.PublicKey, actual[1].PrimaryKey.PublicKey)
	require.Equal(expected2.PrivateKey.PublicKey.PublicKey, actual[1].PrivateKey.PublicKey.PublicKey)
}
