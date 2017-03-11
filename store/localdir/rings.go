package localdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store/model"
)

func (s LocaldirStore) GetRing(ringType model.RingType) ([]byte, error) {
	ringFile, err := s.ringFileName(ringType)
	if err != nil {
		return nil, err
	}
	raw, err := ioutil.ReadFile(ringFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Read ring file failed")
	}
	return raw, nil
}

func (s LocaldirStore) StoreRing(ringType model.RingType, raw []byte) error {
	ringFile, err := s.ringFileName(ringType)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(ringFile), 0700); err != nil {
		return errors.Wrap(err, "Create rings directory failed")
	}

	current, err := s.GetRing(ringType)
	if err != nil {
		return err
	}

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

func (s *LocaldirStore) ringFileName(ringType model.RingType) (string, error) {
	return filepath.Join(s.baseDir, "rings", string(ringType)), nil
}
