package scryptlib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/scrypt"
)

func Encrypt(password []byte, in io.Reader, out io.Writer) error {
	params, err := pickparams(1024*1024*1024, 1.0)
	if err != nil {
		return err
	}
	return EncryptWithParams(password, in, out, params)
}

func EncryptWithParams(password []byte, in io.Reader, out io.Writer, params *Params) error {
	header := Header{
		Magic:   Magic,
		Version: 0,
		Params:  *params,
	}
	if n, err := rand.Read(header.Salt[:]); err != nil || n != 32 {
		return errors.New("Failed to create salt")
	}
	headerOut := bytes.NewBuffer(nil)
	if err := binary.Write(headerOut, binary.BigEndian, &header); err != nil {
		return err
	}
	headerBytes := headerOut.Bytes()
	headerhash := sha256.Sum256(headerBytes[0:48])
	copy(headerBytes[48:64], headerhash[:16])

	N := 1 << header.Params.LogN
	dk, err := scrypt.Key(password, header.Salt[:], N, int(header.Params.R), int(header.Params.P), 64)
	if err != nil {
		return err
	}

	hmac256 := hmac.New(sha256.New, dk[32:])
	hmac256.Write(headerBytes[:64])
	copy(headerBytes[64:96], hmac256.Sum(nil))

	if _, err := out.Write(headerBytes); err != nil {
		return err
	}
	hmac256.Write(headerBytes[64:96])

	return encryptStream(dk[0:32], hmac256, in, out)
}

func encryptStream(aesKey []byte, hmac256 hash.Hash, in io.Reader, out io.Writer) error {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return err
	}
	iv := make([]byte, 16)
	stream := cipher.NewCTR(block, iv)

	buffer := make([]byte, 8192)

	for {
		n, err := in.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		stream.XORKeyStream(buffer[:n], buffer[:n])
		hmac256.Write(buffer[:n])
		if _, outErr := out.Write(buffer[:n]); outErr != nil {
			return outErr
		}
	}
	if _, err := out.Write(hmac256.Sum(nil)); err != nil {
		return err
	}

	return nil
}
