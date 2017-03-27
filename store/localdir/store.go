package localdir

import (
	"crypto/sha512"
	"encoding/hex"
	"net/url"
	"os"
	"sync"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
)

// Store is an implementation backed by a local directory
type Store struct {
	lock    sync.Mutex
	baseDir string
	logger  logging.Logger
}

// NewLocaldirStore creates a new store backed by a local directory
// Note its save to distribute this directory among several machine via
// Dropbox, opencloud or similiar
func NewLocaldirStore(dirURL *url.URL, logger logging.Logger) (*Store, error) {
	baseDir := dirURL.Path
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, errors.Wrap(err, "Create store directory failed")
	}
	return &Store{
		baseDir: baseDir,
		logger:  logger.WithField("package", "store/localdir"),
	}, nil
}

func generateID(data []byte) (string, error) {
	hash := sha512.New()
	if _, err := hash.Write(data); err != nil {
		return "", errors.Wrap(err, "Hashing failed")
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
