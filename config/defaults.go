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

func DefaultConfigFile() string {
	return filepath.Join(findConfigPath(), "config.yaml")
}

func DefaultStoreURL() string {
	return "file://" + filepath.Join(findConfigPath(), "store")
}
