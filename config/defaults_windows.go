package config

import (
	"os"
	"path/filepath"
)

// DefaultStoreURL gets the default location of the store
func DefaultStoreURL() string {
	return "file:///" + filepath.ToSlash(filepath.Join(os.Getenv("USERPROFILE"), ".trustless_store"))
}
