package config

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"

	"github.com/pkg/errors"
)

type CommonConfig struct {
	StoreURL string `json:"store_url" yaml:"store_url"`
	NodeID   string `json:"node_id" yaml:"node_id"`
}

func DefaultCommonConfig() (*CommonConfig, error) {
	nodeID, err := generateNodeID()
	if err != nil {
		return nil, err
	}
	return &CommonConfig{
		StoreURL: DefaultStoreURL(),
		NodeID:   nodeID,
	}, nil
}

func generateNodeID() (string, error) {
	jitter := make([]byte, 65536)
	if _, err := rand.Read(jitter); err != nil {
		return "", errors.Wrap(err, "Secure random failed")
	}
	hash := sha512.New()
	if _, err := hash.Write(jitter); err != nil {
		return "", errors.Wrap(err, "Hashing failed")
	}

	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil)), nil
}
