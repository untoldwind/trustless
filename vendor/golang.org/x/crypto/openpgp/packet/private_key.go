// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

import (
	"bytes"
	"crypto"
	"crypto/cipher"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"math/big"
	"strconv"
	"time"

	"golang.org/x/crypto/openpgp/elgamal"
	"golang.org/x/crypto/openpgp/errors"
	"golang.org/x/crypto/openpgp/s2k"
)

// PrivateKey represents a possibly encrypted private key. See RFC 4880,
// section 5.5.3.
type PrivateKey struct {
	PublicKey
	Encrypted     bool // if true then the private key is unavailable until Decrypt has been called.
	encryptedData []byte
	cipher        CipherFunction
	s2k           func(out, in []byte)
	PrivateKey    interface{} // An *{rsa|dsa|ecdsa}.PrivateKey or a crypto.Signer.
	sha1Checksum  bool
	iv            []byte
	s2kSalt       []byte // s2k parameters
	s2kMode       s2k.Mode
	s2kConfig     s2k.Config
	s2kType       s2kType
}

type s2kType uint8

const (
	s2kTypeNone     s2kType = 0
	s2kTypeSHA1     s2kType = 254
	s2kTypeChecksum s2kType = 255
)

func NewRSAPrivateKey(currentTime time.Time, priv *rsa.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewRSAPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewDSAPrivateKey(currentTime time.Time, priv *dsa.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewDSAPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewElGamalPrivateKey(currentTime time.Time, priv *elgamal.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewElGamalPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

func NewECDSAPrivateKey(currentTime time.Time, priv *ecdsa.PrivateKey) *PrivateKey {
	pk := new(PrivateKey)
	pk.PublicKey = *NewECDSAPublicKey(currentTime, &priv.PublicKey)
	pk.PrivateKey = priv
	return pk
}

// NewSignerPrivateKey creates a sign-only PrivateKey from a crypto.Signer that
// implements RSA or ECDSA.
func NewSignerPrivateKey(currentTime time.Time, signer crypto.Signer) *PrivateKey {
	pk := new(PrivateKey)
	switch pubkey := signer.Public().(type) {
	case rsa.PublicKey:
		pk.PublicKey = *NewRSAPublicKey(currentTime, &pubkey)
		pk.PubKeyAlgo = PubKeyAlgoRSASignOnly
	case ecdsa.PublicKey:
		pk.PublicKey = *NewECDSAPublicKey(currentTime, &pubkey)
	default:
		panic("openpgp: unknown crypto.Signer type in NewSignerPrivateKey")
	}
	pk.PrivateKey = signer
	return pk
}

func (pk *PrivateKey) parse(r io.Reader) (err error) {
	err = (&pk.PublicKey).parse(r)
	if err != nil {
		return
	}
	var buf [1]byte
	_, err = readFull(r, buf[:])
	if err != nil {
		return
	}

	pk.s2kType = s2kType(buf[0])

	switch pk.s2kType {
	case s2kTypeNone:
		pk.s2k = nil
		pk.Encrypted = false
	case s2kTypeSHA1, s2kTypeChecksum:
		_, err = readFull(r, buf[:])
		if err != nil {
			return
		}
		pk.cipher = CipherFunction(buf[0])
		pk.Encrypted = true
		pk.s2k, pk.s2kMode, pk.s2kConfig.Hash, pk.s2kSalt, pk.s2kConfig.S2KCount, err = s2k.ParseWithParameters(r)
		if err != nil {
			return
		}
		if pk.s2kType == s2kTypeSHA1 {
			pk.sha1Checksum = true
		}
	default:
		return errors.UnsupportedError("deprecated s2k function in private key")
	}

	if pk.Encrypted {
		blockSize := pk.cipher.blockSize()
		if blockSize == 0 {
			return errors.UnsupportedError("unsupported cipher in private key: " + strconv.Itoa(int(pk.cipher)))
		}
		pk.iv = make([]byte, blockSize)
		_, err = readFull(r, pk.iv)
		if err != nil {
			return
		}
	}

	pk.encryptedData, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}

	if !pk.Encrypted {
		return pk.parsePrivateKey(pk.encryptedData)
	}

	return
}

func mod64kHash(d []byte) uint16 {
	var h uint16
	for _, b := range d {
		h += uint16(b)
	}
	return h
}

func (pk *PrivateKey) serializeEncrypted(w io.Writer) error {
	hashID, ok := s2k.HashToHashId(pk.s2kConfig.Hash)
	if !ok {
		return errors.UnsupportedError("no such hash algorithm")
	}

	encodedKeyBuf := bytes.NewBuffer(nil)
	encodedKeyBuf.WriteByte(uint8(pk.s2kType))
	encodedKeyBuf.WriteByte(uint8(pk.cipher))
	encodedKeyBuf.WriteByte(uint8(pk.s2kMode))
	encodedKeyBuf.WriteByte(hashID)
	encodedKeyBuf.Write(pk.s2kSalt)
	encodedKeyBuf.WriteByte(pk.s2kConfig.EncodedCount())

	encodedKey := encodedKeyBuf.Bytes()

	w.Write(encodedKey)
	w.Write(pk.iv)
	w.Write(pk.encryptedData)

	return nil
}

func (pk *PrivateKey) serializeUnencrypted(w io.Writer) (err error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(uint8(s2kTypeNone))

	privateKeyBuf := bytes.NewBuffer(nil)

	err = pk.serializePGPPrivate(privateKeyBuf)
	if err != nil {
		return
	}

	privateKeyBytes := privateKeyBuf.Bytes()

	_, err = w.Write([]byte{uint8(s2kTypeNone)})
	if err != nil {
		return
	}
	_, err = w.Write(privateKeyBytes)
	if err != nil {
		return
	}
	if pk.sha1Checksum {
		hash := sha1.New()
		_, err = hash.Write(privateKeyBytes)
		if err != nil {
			return
		}
		sum := hash.Sum(nil)
		_, err = w.Write(sum)
	} else {
		checksum := mod64kHash(privateKeyBytes)
		var checksumBytes [2]byte
		checksumBytes[0] = byte(checksum >> 8)
		checksumBytes[1] = byte(checksum)
		_, err = w.Write(checksumBytes[:])
	}

	return
}

func (pk *PrivateKey) Serialize(w io.Writer) (err error) {
	publicKeyBuf := bytes.NewBuffer(nil)
	err = pk.PublicKey.serializeWithoutHeaders(publicKeyBuf)
	if err != nil {
		return
	}

	privateKeyBuf := bytes.NewBuffer(nil)
	if pk.Encrypted {
		err = pk.serializeEncrypted(privateKeyBuf)
	} else {
		err = pk.serializeUnencrypted(privateKeyBuf)
	}
	if err != nil {
		return
	}

	ptype := packetTypePrivateKey
	publicKeyBytes := publicKeyBuf.Bytes()
	privateKeyBytes := privateKeyBuf.Bytes()
	if pk.IsSubkey {
		ptype = packetTypePrivateSubkey
	}
	err = serializeHeader(w, ptype, len(publicKeyBytes)+len(privateKeyBytes))
	if err != nil {
		return
	}
	_, err = w.Write(publicKeyBytes)
	if err != nil {
		return
	}
	_, err = w.Write(privateKeyBytes)

	return
}

func (pk *PrivateKey) serializePGPPrivate(privateKeyBuf io.Writer) (err error) {
	switch priv := pk.PrivateKey.(type) {
	case *rsa.PrivateKey:
		err = serializeRSAPrivateKey(privateKeyBuf, priv)
	case *dsa.PrivateKey:
		err = serializeDSAPrivateKey(privateKeyBuf, priv)
	case *elgamal.PrivateKey:
		err = serializeElGamalPrivateKey(privateKeyBuf, priv)
	case *ecdsa.PrivateKey:
		err = serializeECDSAPrivateKey(privateKeyBuf, priv)
	default:
		err = errors.InvalidArgumentError("unknown private key type")
	}
	return
}

func serializeRSAPrivateKey(w io.Writer, priv *rsa.PrivateKey) error {
	err := writeBig(w, priv.D)
	if err != nil {
		return err
	}
	err = writeBig(w, priv.Primes[1])
	if err != nil {
		return err
	}
	err = writeBig(w, priv.Primes[0])
	if err != nil {
		return err
	}
	return writeBig(w, priv.Precomputed.Qinv)
}

func serializeDSAPrivateKey(w io.Writer, priv *dsa.PrivateKey) error {
	return writeBig(w, priv.X)
}

func serializeElGamalPrivateKey(w io.Writer, priv *elgamal.PrivateKey) error {
	return writeBig(w, priv.X)
}

func serializeECDSAPrivateKey(w io.Writer, priv *ecdsa.PrivateKey) error {
	return writeBig(w, priv.D)
}

// Encrypt encrypts the private key with default parameters
func (pk *PrivateKey) Encrypt(passphrase []byte) error {
	return pk.EncryptWithParameters(passphrase, CipherAES128, s2k.ModeIterated, s2k.Config{
		S2KCount: 65536,
		Hash:     crypto.SHA1,
	})
}

// EncryptWithParameters encrypts the private key with given parameters
func (pk *PrivateKey) EncryptWithParameters(passphrase []byte, cipherFunc CipherFunction, s2kMode s2k.Mode, s2kConfig s2k.Config) (err error) {
	if pk.Encrypted {
		err = errors.InvalidArgumentError("Key already encrypted")
		return
	}

	privateKeyBuf := bytes.NewBuffer(nil)
	err = pk.serializePGPPrivate(privateKeyBuf)
	if err != nil {
		return
	}
	privateKeyBytes := privateKeyBuf.Bytes()

	pk.cipher = cipherFunc
	pk.s2kMode = s2kMode
	pk.s2kConfig = s2kConfig

	pk.s2kSalt = make([]byte, 8)
	_, err = rand.Read(pk.s2kSalt)
	if err != nil {
		return
	}
	pk.s2k, err = s2k.New(pk.s2kMode, pk.s2kSalt, pk.s2kConfig)
	if err != nil {
		return
	}

	key := make([]byte, pk.cipher.KeySize())
	pk.s2k(key, passphrase)

	block := pk.cipher.new(key)
	pk.iv = make([]byte, pk.cipher.blockSize())
	_, err = rand.Read(pk.iv)
	if err != nil {
		return
	}
	cfb := cipher.NewCFBEncrypter(block, pk.iv)

	if pk.sha1Checksum {
		pk.s2kType = s2kTypeSHA1
		h := sha1.New()
		_, err = h.Write(privateKeyBytes)
		if err != nil {
			return
		}
		sum := h.Sum(nil)
		privateKeyBytes = append(privateKeyBytes, sum...)
	} else {
		pk.s2kType = s2kTypeChecksum
		checksum := mod64kHash(privateKeyBytes)
		privateKeyBytes = append(privateKeyBytes, byte(checksum>>8), byte(checksum))
	}

	pk.encryptedData = make([]byte, len(privateKeyBytes))

	cfb.XORKeyStream(pk.encryptedData, privateKeyBytes)

	pk.Encrypted = true

	return
}

// Decrypt decrypts an encrypted private key using a passphrase.
func (pk *PrivateKey) Decrypt(passphrase []byte) error {
	if !pk.Encrypted {
		return nil
	}

	key := make([]byte, pk.cipher.KeySize())
	pk.s2k(key, passphrase)
	block := pk.cipher.new(key)
	cfb := cipher.NewCFBDecrypter(block, pk.iv)

	data := make([]byte, len(pk.encryptedData))
	cfb.XORKeyStream(data, pk.encryptedData)

	if pk.sha1Checksum {
		if len(data) < sha1.Size {
			return errors.StructuralError("truncated private key data")
		}
		h := sha1.New()
		h.Write(data[:len(data)-sha1.Size])
		sum := h.Sum(nil)
		if !bytes.Equal(sum, data[len(data)-sha1.Size:]) {
			return errors.StructuralError("private key checksum failure")
		}
		data = data[:len(data)-sha1.Size]
	} else {
		if len(data) < 2 {
			return errors.StructuralError("truncated private key data")
		}
		var sum uint16
		for i := 0; i < len(data)-2; i++ {
			sum += uint16(data[i])
		}
		if data[len(data)-2] != uint8(sum>>8) ||
			data[len(data)-1] != uint8(sum) {
			return errors.StructuralError("private key checksum failure")
		}
		data = data[:len(data)-2]
	}

	return pk.parsePrivateKey(data)
}

func (pk *PrivateKey) parsePrivateKey(data []byte) (err error) {
	switch pk.PublicKey.PubKeyAlgo {
	case PubKeyAlgoRSA, PubKeyAlgoRSASignOnly, PubKeyAlgoRSAEncryptOnly:
		return pk.parseRSAPrivateKey(data)
	case PubKeyAlgoDSA:
		return pk.parseDSAPrivateKey(data)
	case PubKeyAlgoElGamal:
		return pk.parseElGamalPrivateKey(data)
	case PubKeyAlgoECDSA:
		return pk.parseECDSAPrivateKey(data)
	}
	panic("impossible")
}

func (pk *PrivateKey) parseRSAPrivateKey(data []byte) (err error) {
	rsaPub := pk.PublicKey.PublicKey.(*rsa.PublicKey)
	rsaPriv := new(rsa.PrivateKey)
	rsaPriv.PublicKey = *rsaPub

	buf := bytes.NewBuffer(data)
	d, _, err := readMPI(buf)
	if err != nil {
		return
	}
	p, _, err := readMPI(buf)
	if err != nil {
		return
	}
	q, _, err := readMPI(buf)
	if err != nil {
		return
	}

	rsaPriv.D = new(big.Int).SetBytes(d)
	rsaPriv.Primes = make([]*big.Int, 2)
	rsaPriv.Primes[0] = new(big.Int).SetBytes(p)
	rsaPriv.Primes[1] = new(big.Int).SetBytes(q)
	if err := rsaPriv.Validate(); err != nil {
		return err
	}
	rsaPriv.Precompute()
	pk.PrivateKey = rsaPriv
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
}

func (pk *PrivateKey) parseDSAPrivateKey(data []byte) (err error) {
	dsaPub := pk.PublicKey.PublicKey.(*dsa.PublicKey)
	dsaPriv := new(dsa.PrivateKey)
	dsaPriv.PublicKey = *dsaPub

	buf := bytes.NewBuffer(data)
	x, _, err := readMPI(buf)
	if err != nil {
		return
	}

	dsaPriv.X = new(big.Int).SetBytes(x)
	pk.PrivateKey = dsaPriv
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
}

func (pk *PrivateKey) parseElGamalPrivateKey(data []byte) (err error) {
	pub := pk.PublicKey.PublicKey.(*elgamal.PublicKey)
	priv := new(elgamal.PrivateKey)
	priv.PublicKey = *pub

	buf := bytes.NewBuffer(data)
	x, _, err := readMPI(buf)
	if err != nil {
		return
	}

	priv.X = new(big.Int).SetBytes(x)
	pk.PrivateKey = priv
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
}

func (pk *PrivateKey) parseECDSAPrivateKey(data []byte) (err error) {
	ecdsaPub := pk.PublicKey.PublicKey.(*ecdsa.PublicKey)

	buf := bytes.NewBuffer(data)
	d, _, err := readMPI(buf)
	if err != nil {
		return
	}

	pk.PrivateKey = &ecdsa.PrivateKey{
		PublicKey: *ecdsaPub,
		D:         new(big.Int).SetBytes(d),
	}
	pk.Encrypted = false
	pk.encryptedData = nil

	return nil
}
