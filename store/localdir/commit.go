package localdir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store/model"
)

// Commit changes made to the store (i.e. write them the the change log)
func (s *Store) Commit(nodeID string, changes []model.Change) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	changeLogFileName := filepath.Join(s.baseDir, "logs", nodeID)
	if err := os.MkdirAll(filepath.Dir(changeLogFileName), 0700); err != nil {
		return errors.Wrapf(err, "Failed to create changelog dir: %s", filepath.Dir(changeLogFileName))
	}
	changeLogFile, err := os.OpenFile(changeLogFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return errors.Wrapf(err, "Failed to open changelog: %s", changeLogFileName)
	}
	defer changeLogFile.Close()

	for _, change := range changes {
		switch change.Operation {
		case model.ChangeOpAdd:
			if _, err := changeLogFile.WriteString(fmt.Sprintf("A %s\n", change.BlockID)); err != nil {
				return errors.Wrap(err, "Failed writing change")
			}
		case model.ChangeOpDelete:
			if _, err := changeLogFile.WriteString(fmt.Sprintf("D %s\n", change.BlockID)); err != nil {
				return errors.Wrap(err, "Failed writing change")
			}
		default:
			return errors.Errorf("Invalid change operation: %s", change.Operation)
		}
	}
	return nil
}
