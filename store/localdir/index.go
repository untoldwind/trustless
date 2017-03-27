package localdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetIndex retrieves the index block of a node
func (s *Store) GetIndex(nodeID string) ([]byte, error) {
	indexFile := s.indexFilename(nodeID)
	indexBlock, err := ioutil.ReadFile(indexFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Read index file failed")
	}
	return indexBlock, nil
}

// StoreIndex stores the index block of a node
func (s *Store) StoreIndex(nodeID string, indexBlock []byte) error {
	indexFile := s.indexFilename(nodeID)
	if err := os.MkdirAll(filepath.Dir(indexFile), 0700); err != nil {
		return errors.Wrap(err, "Create indexes directory failed")
	}
	if err := ioutil.WriteFile(indexFile, indexBlock, 0600); err != nil {
		return errors.Wrap(err, "Writing index file failed")
	}
	return nil
}

func (s *Store) indexFilename(nodeID string) string {
	return filepath.Join(s.baseDir, "indexes", nodeID)
}
