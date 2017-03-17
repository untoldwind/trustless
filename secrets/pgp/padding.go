package pgp

import (
	"bytes"
	"math/rand"

	"github.com/pkg/errors"
)

func padBlock(content []byte) ([]byte, error) {
	contentSize := len(content)
	paddedSize := ((contentSize >> 13) + 1) << 13
	fillSize := paddedSize - contentSize
	padBytes := make([]byte, fillSize-1)

	if _, err := rand.Read(padBytes); err != nil {
		return nil, errors.Wrap(err, "Read random failed")
	}
	idx := bytes.IndexByte(padBytes, 0)
	var prefix []byte
	var suffix []byte
	if idx < 0 {
		padBytes[0] = 0
		prefix = padBytes[0:1]
		suffix = padBytes[0:]
	} else {
		prefix = padBytes[0 : idx+1]
		suffix = padBytes[idx:]
	}

	result := make([]byte, 0, paddedSize)
	result = append(result, prefix...)
	result = append(result, content...)
	result = append(result, suffix...)

	return result, nil
}

func unpadBlock(block []byte) ([]byte, error) {
	idx1 := bytes.IndexByte(block, 0)
	if idx1 < 0 {
		return nil, errors.New("Invalid prefix padding")
	}
	idx2 := bytes.IndexByte(block[idx1+1:], 0)
	if idx2 < 0 {
		return nil, errors.New("Invalid suffix padding")
	}

	return block[idx1+1 : idx1+1+idx2], nil
}
