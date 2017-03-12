package secrets

import (
	"time"

	"github.com/untoldwind/trustless/store"
	"github.com/untoldwind/trustless/store/model"
)

type blockRef struct {
	delete    bool
	timestamp time.Time
}

type blockSet map[string]blockRef

func newBlockSet(store store.Store) (blockSet, error) {
	result := blockSet{}
	commits := map[string]bool{}

	heads, err := store.Heads()
	if err != nil {
		return nil, err
	}
	for _, head := range heads {
		commitID := head.CommitID
		for commitID != "" {
			if _, ok := commits[commitID]; ok {
				continue
			}
			commit, err := store.GetCommit(commitID)
			if err != nil {
				return nil, err
			}
			commits[commitID] = true

			for _, change := range commit.Changes {
				ref, ok := result[change.BlockID]
				switch change.Operation {
				case model.ChangeOpAdd:
					if !ok || commit.Timestamp.After(ref.timestamp) {
						result[change.BlockID] = blockRef{delete: false, timestamp: commit.Timestamp}
					}
				case model.ChangeOpDelete:
					if !ok || commit.Timestamp.After(ref.timestamp) {
						result[change.BlockID] = blockRef{delete: true, timestamp: commit.Timestamp}
					}
				}
			}
		}
	}

	return result, err
}
