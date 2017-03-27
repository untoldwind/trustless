package localdir

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/store/model"
)

// ChangeLogs retrieves the change logs of all nodes
func (s *Store) ChangeLogs() ([]model.ChangeLog, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	commitLogDir := filepath.Join(s.baseDir, "logs")
	files, err := ioutil.ReadDir(commitLogDir)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Open heads dir failed")
	}
	result := make([]model.ChangeLog, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		commitLogFile := filepath.Join(commitLogDir, file.Name())
		changes, err := s.parseChangeLog(commitLogFile)
		if err != nil {
			return nil, err
		}
		result = append(result, model.ChangeLog{
			NodeID:  file.Name(),
			Changes: changes,
		})
	}

	return result, nil
}

func (s *Store) parseChangeLog(filename string) ([]model.Change, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open changelog: %s", filename)
	}
	defer file.Close()

	changes := make([]model.Change, 0, 1000)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) != 2 {
			s.logger.Warnf("Invalid entry in %s", filename)
			continue
		}
		switch parts[0] {
		case "A":
			changes = append(changes, model.Change{
				Operation: model.ChangeOpAdd,
				BlockID:   parts[1],
			})
		case "D":
			changes = append(changes, model.Change{
				Operation: model.ChangeOpDelete,
				BlockID:   parts[1],
			})
		default:
			s.logger.Warnf("Invalid operation in %s: %s", filename, parts[0])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "Failed to parse changelog: %s", filename)
	}
	return changes, nil
}
