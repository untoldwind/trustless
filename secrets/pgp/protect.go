package pgp

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"math/big"

	"github.com/awnumar/memguard"
	"github.com/pkg/errors"

	"golang.org/x/crypto/openpgp/elgamal"
	"golang.org/x/crypto/openpgp/packet"
)

func (s *pgpSecrets) preparePurge() {
	for _, buffer := range s.buffers {
		buffer.MakeMutable()
	}
}

func (s *pgpSecrets) destroyBuffers() {
	for _, buffer := range s.buffers {
		buffer.Destroy()
	}
	s.buffers = nil
}

func (s *pgpSecrets) protectPrivateKey(key *packet.PrivateKey) error {
	if key.PrivateKey == nil || key.Encrypted {
		return nil
	}
	switch privKey := key.PrivateKey.(type) {
	case *rsa.PrivateKey:
		privKey.Precompute()
		s.protectBigInt(privKey.D)
		for _, p := range privKey.Primes {
			s.protectBigInt(p)
		}
		s.protectBigInt(privKey.Precomputed.Dp)
		s.protectBigInt(privKey.Precomputed.Dq)
		s.protectBigInt(privKey.Precomputed.Qinv)
		for _, crt := range privKey.Precomputed.CRTValues {
			s.protectBigInt(crt.Coeff)
			s.protectBigInt(crt.Exp)
			s.protectBigInt(crt.R)
		}
	case *dsa.PrivateKey:
		s.protectBigInt(privKey.X)
	case *elgamal.PrivateKey:
		s.protectBigInt(privKey.X)
	case *ecdsa.PrivateKey:
		s.protectBigInt(privKey.D)
	default:
		s.logger.Warn("Unable to purge key from memory")
	}
	return nil
}

func (s *pgpSecrets) protectBigInt(i *big.Int) error {
	words := i.Bits()
	locked, err := memguard.NewImmutableFromBytes(wordsToBytes(words))
	if err != nil {
		return errors.Wrap(err, "Memguard failed")
	}
	s.buffers = append(s.buffers, locked)
	i.SetBits(bytesToWords(locked.Buffer()))

	return nil
}
