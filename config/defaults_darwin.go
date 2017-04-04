package config

import (
	"net/url"
	"os"
	"path/filepath"
)

// DefaultStoreURL gets the default location of the store
func DefaultStoreURL() string {
	return "file://" + url.PathEscape(filepath.Join(os.Getenv("HOME"), ".trustless_store"))
}
