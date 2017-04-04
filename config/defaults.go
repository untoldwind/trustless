package config

import (
	"os"
	"path/filepath"

	"github.com/shibukawa/configdir"
)

var configDirs = configdir.New("untoldwind", "trustless")

func findConfigPath() string {
	localSettings := configDirs.QueryFolders(configdir.Global)
	if len(localSettings) > 0 {
		return localSettings[0].Path
	}
	return filepath.Join(os.Getenv("HOME"), ".trustless")
}

// DefaultConfigFile gets the default location of the configuration file
func DefaultConfigFile() string {
	return filepath.Join(findConfigPath(), "config.yaml")
}
