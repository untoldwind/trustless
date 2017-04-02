package localdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// AddBlock adds a block (of encrypted data) to the store and
// return its id
func (s *Store) AddBlock(block []byte) (string, error) {
	blockID, err := generateID(block)
	if err != nil {
		return "", err
	}
	blockFile, err := s.blockFileName(blockID)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(blockFile), 0700); err != nil {
		return "", errors.Wrap(err, "Create block directory failed")
	}
	if err := ioutil.WriteFile(blockFile, block, 0600); err != nil {
		return "", errors.Wrap(err, "Writing blockfile failed")
	}

	return blockID, nil
}

// GetBlock retrieves a block by its id
func (s *Store) GetBlock(blockID string) ([]byte, error) {
	blockFile, err := s.blockFileName(blockID)
	if err != nil {
		return nil, err
	}
	block, err := ioutil.ReadFile(blockFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Read block file failed")
	}
	return block, nil
}

func (s *Store) blockFileName(blockID string) (string, error) {
	if len(blockID) < 2 {
		return "", errors.New("BlockID too short")
	}

	return filepath.Join(s.baseDir, "blocks", blockID[0:1], blockID), nil
}
