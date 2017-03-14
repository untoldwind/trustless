package secrets

import "github.com/untoldwind/trustless/api"

type Index struct {
	Entries     map[string]*api.SecretEntry
	EntryBlocks map[string][]string
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
