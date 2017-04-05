package config

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/pkg/errors"
)

// CommonConfig contains the common client/daemon configuration
type CommonConfig struct {
	StoreURL          string        `json:"store_url" yaml:"store_url"`
	NodeID            string        `json:"node_id" yaml:"node_id"`
	UnlockTimeout     time.Duration `json:"unlock_timeout" yaml:"unlock_timeout"`
	UnlockTimeoutHard bool          `json:"unlock_timeout_hard" yaml:"unlock_timeout_hard"`
}

// DefaultCommonConfig create a CommonConfig with reasonable defaults
func DefaultCommonConfig() (*CommonConfig, error) {
	nodeID, err := generateNodeID()
	if err != nil {
		return nil, err
	}
	return &CommonConfig{
		StoreURL:          DefaultStoreURL(),
		NodeID:            nodeID,
		UnlockTimeout:     5 * time.Minute,
		UnlockTimeoutHard: false,
	}, nil
}

func generateNodeID() (string, error) {
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
