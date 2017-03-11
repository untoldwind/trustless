package localdir

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store/model"
)

func (s *LocaldirStore) Commit(nodeID string, changes []model.Change) (string, error) {
	head, err := s.GetHead(nodeID)
	if err != nil {
		return "", err
	}
	commit := model.Commit{
		NodeID:       nodeID,
		PrevCommitID: head,
		Changes:      changes,
	}
	raw, err := json.Marshal(&commit)
	if err != nil {
		return "", errors.Wrap(err, "Json marshal of commit failed")
	}
	commitID, err := generateID(raw)
	if err != nil {
		return "", err
	}
	commitFile, err := s.commitFileName(commitID)
	if err := os.MkdirAll(filepath.Dir(commitFile), 0700); err != nil {
		return "", errors.Wrap(err, "Create block directory failed")
	}
	if err := ioutil.WriteFile(commitFile, raw, 0600); err != nil {
		return "", errors.Wrap(err, "Writing blockfile failed")
	}
	if err := s.storeHead(nodeID, commitID); err != nil {
		return "", err
	}

	return commitID, nil
}

func (s *LocaldirStore) GetCommit(commitID string) (*model.Commit, error) {
	commitFile, err := s.commitFileName(commitID)
	raw, err := ioutil.ReadFile(commitFile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read commit file")
	}
	var commit model.Commit
	if err := json.Unmarshal(raw, &commit); err != nil {
		return nil, errors.Wrap(err, "Parsing commit file failed")
	}
	return &commit, nil
}

func (s *LocaldirStore) commitFileName(commitID string) (string, error) {
	if len(commitID) < 3 {
		return "", errors.New("CommitID too short")
	}

	return filepath.Join(s.baseDir, "commits", commitID[0:2], commitID), nil
}
