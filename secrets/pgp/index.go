package pgp

import (
	"encoding/json"
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
	KnownBlocks   secrets.IDSet
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

func (i *Index) registerChanges(changedBlocks map[string]*secrets.SecretBlock) {
	i.lock.Lock()
	defer i.lock.Unlock()

	for blockID, secretBlock := range changedBlocks {
		if i.KnownBlocks.Contains(blockID) || i.DeletedBlocks.Contains(blockID) {
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
						URLs:      secretBlock.Version.URLs,
						Timestamp: secretBlock.Version.Timestamp,
						Deleted:   secretBlock.Version.Deleted,
					},
					Blocks: secrets.IDSet{},
				}
				i.Entries[secretBlock.ID] = entry
			}
			entry.Blocks.Add(blockID)
			if secretBlock.Version.Timestamp.After(entry.Timestamp) {
				entry.Name = secretBlock.Version.Name
				entry.Tags = secretBlock.Version.Tags
				entry.URLs = secretBlock.Version.URLs
				entry.Timestamp = secretBlock.Version.Timestamp
				entry.Deleted = secretBlock.Version.Deleted
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
		i.KnownBlocks.Add(blockID)
	}
}

func (i *Index) serialize() ([]byte, error) {
	i.lock.Lock()
	defer i.lock.Unlock()
	return json.Marshal(i)
}

func (s *pgpSecrets) buildIndex() error {
	s.ensureIndex()
	changeLogs, err := s.store.ChangeLogs()
	if err != nil {
		return err
	}
	changed := false
	for _, changeLog := range changeLogs {
		changedBlocks := map[string]*secrets.SecretBlock{}
		for _, change := range changeLog.Changes {
			if s.index.KnownBlocks.Contains(change.BlockID) {
				continue
			}
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
		if len(changedBlocks) > 0 {
			s.index.registerChanges(changedBlocks)
			changed = true
		}
	}
	if changed {
		s.storeIndex()
	}

	return nil
}

func (s *pgpSecrets) ensureIndex() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.index == nil {
		s.index = &Index{
			Entries:       map[string]*IndexEntry{},
			KnownBlocks:   secrets.IDSet{},
			Tags:          secrets.IDSet{},
			DeletedBlocks: secrets.IDSet{},
		}
	}
}

func (s *pgpSecrets) fetchIndex() {
	indexBlock, err := s.store.GetIndex(s.nodeID)
	if err != nil {
		s.logger.Warnf("Failed to retrieve index block: %v", err)
		return
	}
	if indexBlock == nil {
		return
	}
	indexData, err := s.decryptData(indexBlock)
	if err != nil {
		s.logger.Warnf("Failed to decrypt index block: %v", err)
		return
	}
	var index Index
	if err := json.Unmarshal(indexData, &index); err != nil {
		s.logger.Warnf("Failed to unmarshal index block: %v", err)
		return
	}
	s.index = &index
}

func (s *pgpSecrets) storeIndex() {
	if s.index == nil {
		return
	}
	indexData, err := s.index.serialize()
	if err != nil {
		s.logger.Warnf("Failed to marshal index: %v", err)
		return
	}
	indexBlock, err := s.encryptData(indexData)
	if err != nil {
		s.logger.Warnf("Failed to encrypt index: %v", err)
		return
	}
	s.store.StoreIndex(s.nodeID, indexBlock)
}
