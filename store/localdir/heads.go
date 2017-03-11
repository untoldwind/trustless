package localdir

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store/model"
)

func (s *LocaldirStore) Heads() ([]model.Head, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	headsDir := filepath.Join(s.baseDir, "heads")
	files, err := ioutil.ReadDir(headsDir)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Open heads dir failed")
	}
	result := make([]model.Head, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		headFile := filepath.Join(headsDir, file.Name())
		commitID, err := ioutil.ReadFile(headFile)
		if err != nil {
			return nil, errors.Wrap(err, "Read head file failed")
		}
		result = append(result, model.Head{
			NodeID:   file.Name(),
			CommitID: string(commitID),
		})
	}

	return result, nil
}

func (s *LocaldirStore) GetHead(nodeID string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	headFile, err := s.headFileName(nodeID)
	if err != nil {
		return "", err
	}
	commitID, err := ioutil.ReadFile(headFile)
	if os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", errors.Wrap(err, "Read head file failed")
	}
	return string(commitID), nil
}

func (s *LocaldirStore) storeHead(nodeID, commitID string) error {
	headFile, err := s.headFileName(nodeID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(headFile), 0700); err != nil {
		return errors.Wrap(err, "Create heads directory failed")
	}
	if err := ioutil.WriteFile(headFile, []byte(commitID), 0600); err != nil {
		return errors.Wrap(err, "Writing head failed")
	}
	return nil
}

func (s *LocaldirStore) headFileName(nodeID string) (string, error) {
	return filepath.Join(s.baseDir, "heads", nodeID), nil
}
