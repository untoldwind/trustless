package config

import (
	"os"
	"path/filepath"
)

// DefaultStoreURL gets the default location of the store
func DefaultStoreURL() string {
	return "file://" + filepath.Join(os.Getenv("HOME"), ".trustless_store")
}
