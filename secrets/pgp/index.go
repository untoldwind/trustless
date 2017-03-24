package pgp

import (
	"sort"
	"sync"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
	"github.com/untoldwind/trustless/store/model"
)

type IndexEntry struct {
	api.SecretEntry
	Blocks secrets.IDSet
}

type Index struct {
	lock          sync.Mutex
	Entries       map[string]*IndexEntry
	Commits       secrets.IDSet
	Tags          secrets.IDSet
	DeletedBlocks secrets.IDSet
}

func (i *Index) list() *api.SecretList {
	i.lock.Lock()
	defer i.lock.Unlock()

	tags := make(sort.StringSlice, 0, len(i.Tags))
	for tag := range i.Tags {
		tags = append(tags, tag)
	}
	tags.Sort()

	entries := make([]*api.SecretEntry, 0, len(i.Entries))
	for _, entry := range i.Entries {
		entries = append(entries, &entry.SecretEntry)
	}
	return &api.SecretList{
		AllTags: tags,
		Entries: entries,
	}
}

func (i *Index) registerCommit(commitID string, changedBlocks map[string]*secrets.SecretBlock) {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.Commits.Contains(commitID) {
		return
	}

	for blockID, secretBlock := range changedBlocks {
		if i.DeletedBlocks.Contains(blockID) {
			continue
		}
		if secretBlock != nil {
			i.Tags.AddAll(secretBlock.Version.Tags)
			entry, ok := i.Entries[secretBlock.ID]
			if !ok {
				entry = &IndexEntry{
					SecretEntry: api.SecretEntry{
						ID:        secretBlock.ID,
						Type:      secretBlock.Type,
						Name:      secretBlock.Version.Name,
						Tags:      secretBlock.Version.Tags,
						Timestamp: secretBlock.Version.Timestamp,
					},
					Blocks: secrets.IDSet{},
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
			i.DeletedBlocks.Add(blockID)
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

func (s *pgpSecrets) buildIndex() error {
	if locked, _ := s.IsLocked(); locked {
		return secrets.ErrSecretsLocked
	}
	s.index = &Index{
		Entries: map[string]*IndexEntry{},
		Commits: secrets.IDSet{},
		Tags:    secrets.IDSet{},
	}
	heads, err := s.store.Heads()
	if err != nil {
		return err
	}
	for _, head := range heads {
		commitID := head.CommitID
		for commitID != "" {
			if s.index.Commits.Contains(commitID) {
				break
			}
			commit, err := s.store.GetCommit(commitID)
			if err != nil {
				return err
			}
			changedBlocks := map[string]*secrets.SecretBlock{}
			for _, change := range commit.Changes {
				switch change.Operation {
				case model.ChangeOpAdd:
					block, err := s.store.GetBlock(change.BlockID)
					if err != nil {
						return err
					}
					secretBlock, err := s.decryptSecret(block)
					if err != nil {
						return err
					}
					changedBlocks[change.BlockID] = secretBlock
				case model.ChangeOpDelete:
					changedBlocks[change.BlockID] = nil
				}
			}
			s.index.registerCommit(commitID, changedBlocks)
			commitID = commit.PrevCommitID
		}
	}

	return nil
}
