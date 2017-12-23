package localdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetRing retrieves the key ring of the store
func (s *Store) GetRing() ([]byte, error) {
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

// StoreRing stores the key ring of the store
func (s *Store) StoreRing(raw []byte) error {
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

// GetPublicRing receives the public key ring of the store
func (s *Store) GetPublicRing() ([]byte, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	ringFile := filepath.Join(s.baseDir, "ring.pub")
	raw, err := ioutil.ReadFile(ringFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Read ring file failed")
	}
	return raw, nil
}

// StorePublicRing stores the public key ring of the store
func (s *Store) StorePublicRing(raw []byte) error {
	ringFile := filepath.Join(s.baseDir, "ring.pub")

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
