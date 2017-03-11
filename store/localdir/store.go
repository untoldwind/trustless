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

type LocaldirStore struct {
	lock    sync.Mutex
	baseDir string
	logger  logging.Logger
}

func NewLocaldirStore(dirUrl *url.URL, logger logging.Logger) (*LocaldirStore, error) {
	baseDir := dirUrl.Path
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, errors.Wrap(err, "Create store directory failed")
	}
	return &LocaldirStore{
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
