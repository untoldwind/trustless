package scryptlib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash"
	"io"

	"golang.org/x/crypto/scrypt"
)

func Decrypt(password []byte, in io.Reader, out io.Writer) error {
	headerBytes := make([]byte, 96)
	if _, err := io.ReadFull(in, headerBytes); err != nil {
		return err
	}
	dk, err := decodeHeader(password, headerBytes)
	if err != nil {
		return err
	}
	hmac256 := hmac.New(sha256.New, dk[32:])
	hmac256.Write(headerBytes)

	contentHMAC, err := decryptStream(dk[0:32], hmac256, in, out)
	if err != nil {
		return err
	}

	if !bytes.Equal(hmac256.Sum(nil), contentHMAC) {
		return errors.New("Content HMAC does not match")
	}

	return nil
}

func decodeHeader(password, headerBytes []byte) ([]byte, error) {
	var header Header

	if err := binary.Read(bytes.NewReader(headerBytes), binary.BigEndian, &header); err != nil {
		return nil, err
	}
	if header.Magic != Magic {
		return nil, errors.New("Invalid file magic")
	}
	if header.Version != 0 {
		return nil, errors.New("Invalid file version")
	}
	headerhash := sha256.Sum256(headerBytes[0:48])

	if !bytes.Equal(headerhash[:16], header.HeaderHash[:]) {
		return nil, errors.New("Header checksum does not match")
	}

	N := 1 << header.Params.LogN
	dk, err := scrypt.Key(password, header.Salt[:], N, int(header.Params.R), int(header.Params.P), 64)
	if err != nil {
		return nil, err
	}

	hmac256 := hmac.New(sha256.New, dk[32:])
	hmac256.Write(headerBytes[:64])

	if !bytes.Equal(hmac256.Sum(nil), header.HeaderHMAC[:]) {
		return nil, errors.New("Header HMAC does not match")
	}

	return dk, nil
}

func decryptStream(aesKey []byte, hmac256 hash.Hash, in io.Reader, out io.Writer) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, 16)
	stream := cipher.NewCTR(block, iv)

	buffer := make([]byte, 8192)
	bufferLen := 0

	for {
		n, err := in.Read(buffer)
		bufferLen += n
		if bufferLen > 32 {
			hmac256.Write(buffer[:bufferLen-32])
			stream.XORKeyStream(buffer[:bufferLen-32], buffer[:bufferLen-32])
			if _, outErr := out.Write(buffer[:bufferLen-32]); outErr != nil {
				return nil, outErr
			}
			copy(buffer[0:32], buffer[bufferLen-32:bufferLen])
			bufferLen = 32
		}
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
	}
	if bufferLen < 32 {
		return nil, errors.New("Content too short (no final hmac)")
	}

	return buffer[bufferLen-32 : bufferLen], nil
}
