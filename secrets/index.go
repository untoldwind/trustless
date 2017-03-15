package secrets

import (
	"sync"

	"github.com/untoldwind/trustless/api"
)

type IndexEntry struct {
	api.SecretEntry
	Blocks IDSet
}

type Index struct {
	lock    sync.Mutex
	Entries map[string]*IndexEntry
	Commits IDSet
}

func (i *Index) registerCommit(commitID string, changes map[string]*SecretBlock) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.Commits.Contains(commitID) {
		return
	}

	for blockID, secretBlock := range changes {
		if secretBlock != nil {
			entry, ok := i.Entries[secretBlock.ID]
			if !ok {
				entry := &IndexEntry{
					SecretEntry: api.SecretEntry{
						ID:        secretBlock.ID,
						Type:      secretBlock.Type,
						Name:      secretBlock.Version.Name,
						Tags:      secretBlock.Version.Tags,
						Timestamp: secretBlock.Version.Timestamp,
					},
					Blocks: IDSet{},
				}
				i.Entries[secretBlock.ID] = entry
			}
			entry.Blocks.Add(blockID)
			if secretBlock.Version.Timestamp.After(entry.Timestamp) {
				entry.Name = secretBlock.Version.Name
				entry.Tags = secretBlock.Version.Tags
				entry.Timestamp = secretBlock.Version.Timestamp
			}
		} else {
			for entryID, entry := range i.Entries {
				if entry.Blocks.Contains(blockID) {
					entry.Blocks.Remove(blockID)
					if len(entry.Blocks) == 0 {
						delete(i.Entries, entryID)
					}
					break
				}
			}
		}
	}
}

func (s *Secrets) buildIndex() error {
	if s.IsLocked() {
		return SecretsLockedError
	}
	blocks, err := newBlockSet(s.store)
	if err != nil {
		return err
	}

	for blockID, blockRef := range blocks {
		if blockRef.delete {
			continue
		}
		block, err := s.store.GetBlock(blockID)
		if err != nil {
			return err
		}
		if len(block) == 0 {
			continue
		}
	}
	return nil
}
