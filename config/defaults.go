package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/pkg/errors"
)

func GenerateNodeID() (string, error) {
	jitter := make([]byte, 65536)
	if _, err := rand.Read(jitter); err != nil {
		return "", errors.Wrap(err, "Secure random failed")
	}
	hash := sha256.New()
	if _, err := hash.Write(jitter); err != nil {
		return "", errors.Wrap(err, "Hashing failed")
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
