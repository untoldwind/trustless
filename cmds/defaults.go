package cmds

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/untoldwind/trustless/config"

	"github.com/pkg/errors"
)

func writeDefaultConfig() error {
	if viper.GetString("store-url") == "" {
		viper.Set("store-url", config.DefaultStoreURL())
	}
	if viper.GetString("node-id") == "" {
		nodeId, err := generateNodeID()
		if err != nil {
			return err
		}
		viper.Set("node-id", nodeId)
	}
	home, err := homedir.Dir()
	if err != nil {
		showError(err)
	}
	os.MkdirAll(filepath.Join(home, ".config"), 0700)

	return viper.WriteConfigAs(filepath.Join(home, ".config", "trustless.toml"))
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
