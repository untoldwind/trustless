package cmds

import (
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/untoldwind/trustless/config"
)

func writeDefaultConfig() error {
	if viper.GetString("store-url") == "" {
		viper.Set("store-url", config.DefaultStoreURL())
	}
	if viper.GetString("node-id") == "" {
		nodeId, err := config.GenerateNodeID()
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
