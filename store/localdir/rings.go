package localdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (s *LocaldirStore) GetRing() ([]byte, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	ringFile := filepath.Join(s.baseDir, "ring")
	raw, err := ioutil.ReadFile(ringFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Read ring file failed")
	}
	return raw, nil
}

func (s *LocaldirStore) StoreRing(raw []byte) error {
	ringFile := filepath.Join(s.baseDir, "ring")

	current, err := s.GetRing()
	if err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if current != nil {
		if err := ioutil.WriteFile(ringFile+".bak", current, 0700); err != nil {
			return errors.Wrap(err, "Writing backup failed")
		}
	}
	if err := ioutil.WriteFile(ringFile, raw, 0700); err != nil {
		return errors.Wrap(err, "Writing ring failed")
	}
	return nil
}
