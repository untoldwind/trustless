package pgp

import (
	"context"
	"strings"

	"github.com/junegunn/fzf/src/algo"

	"github.com/junegunn/fzf/src/util"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

const (
	slab16Size int = 100 * 1024 // 200KB * 32 = 12.8MB
	slab32Size int = 2048       // 8KB * 32 = 256KB
)

func (s *pgpSecrets) List(ctx context.Context, filter api.SecretListFilter) (*api.SecretList, error) {
	if s.isLocked() {
		return nil, secrets.ErrSecretsLocked
	}
	s.logger.Info("List secrets")

	if err := s.buildIndex(); err != nil {
		return nil, err
	}

	s.autolocker.Reset()

	list := s.index.list()

	slab := util.MakeSlab(slab16Size, slab32Size)
	pattern := algo.NormalizeRunes([]rune(strings.ToLower(filter.Name)))

	filtered := make([]*api.SecretEntry, 0, len(list.Entries))
	for _, entry := range list.Entries {
		if match := filterMatch(entry, filter, pattern, slab); match != nil {
			filtered = append(filtered, match)
		}
	}
	secrets.EntrySortNameAsc(filtered)
	return &api.SecretList{
		AllTags: list.AllTags,
		Entries: filtered,
	}, nil
}

func filterMatch(entry *api.SecretEntry, filter api.SecretListFilter, pattern []rune, slab *util.Slab) *api.SecretEntry {
	var nameHighlights []int
	if len(pattern) > 0 {
		input := util.ToChars([]byte(entry.Name))
		result, pos := algo.FuzzyMatchV2(false, true, true, &input, pattern, true, slab)
		if result.Start < 0 {
			return nil
		}
		if pos != nil {
			nameHighlights = *pos
		}
	}
	if filter.Tag != "" && !entry.HasTag(filter.Tag) {
		return nil
	}
	if filter.URL != "" && !entry.MatchesURL(filter.URL) {
		return nil
	}
	if filter.Type != "" && entry.Type != filter.Type {
		return nil
	}
	if filter.Deleted != entry.Deleted {
		return nil
	}
	if len(nameHighlights) > 0 {
		cloned := *entry
		cloned.NameHighlights = nameHighlights
		return &cloned
	}
	return entry
}
