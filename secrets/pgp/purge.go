package pgp

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"math/big"

	"golang.org/x/crypto/openpgp/elgamal"
	"golang.org/x/crypto/openpgp/packet"
)

// overwrite private part of key with zeros
func (s *pgpSecrets) purgePrivateKey(key *packet.PrivateKey) {
	if key.PrivateKey == nil || key.Encrypted {
		return
	}
	switch privKey := key.PrivateKey.(type) {
	case *rsa.PrivateKey:
		clearBigInt(privKey.D)
		for _, p := range privKey.Primes {
			clearBigInt(p)
		}
		clearBigInt(privKey.Precomputed.Dp)
		clearBigInt(privKey.Precomputed.Dq)
		clearBigInt(privKey.Precomputed.Qinv)
		for _, crt := range privKey.Precomputed.CRTValues {
			clearBigInt(crt.Coeff)
			clearBigInt(crt.Exp)
			clearBigInt(crt.R)
		}
	case *dsa.PrivateKey:
		clearBigInt(privKey.X)
	case *elgamal.PrivateKey:
		clearBigInt(privKey.X)
	case *ecdsa.PrivateKey:
		clearBigInt(privKey.D)
	default:
		s.logger.Warn("Unable to purge key from memory")
	}
	key.PrivateKey = nil
	key.Encrypted = true
}

func clearBigInt(i *big.Int) {
	if i == nil {
		return
	}
	words := i.Bits()
	for i := range i.Bits() {
		words[i] = 0
	}
}
