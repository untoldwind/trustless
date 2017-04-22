package api

import (
	"net/url"
	"strings"
	"time"
)

// SecretListFilter filter options to query secret list
type SecretListFilter struct {
	URL     string     `json:"url,omitempty"`
	Tag     string     `json:"tag,omitempty"`
	Type    SecretType `json:"type,omitempty"`
	Name    string     `json:"name,omitempty"`
	Deleted bool       `json:"deleted,omitempty"`
}

// SecretEntry is a reference to a stored secret
type SecretEntry struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Type      SecretType `json:"type"`
	Tags      []string   `json:"tags"`
	URLs      []string   `json:"urls"`
	Timestamp time.Time  `json:"timestamp"`
	Deleted   bool       `json:"deleted"`
}

func (e *SecretEntry) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (e *SecretEntry) MatchesURL(urlStr string) bool {
	lookupURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	for _, entryURLStr := range e.URLs {
		entryURL, err := url.Parse(entryURLStr)
		if err != nil {
			continue
		}
		if urlMatches(entryURL, lookupURL) {
			return true
		}
	}
	return false
}

func (e *SecretEntry) Matches(filter SecretListFilter) bool {
	if filter.Name != "" && !strings.Contains(e.Name, filter.Name) {
		return false
	}
	if filter.Tag != "" && !e.HasTag(filter.Tag) {
		return false
	}
	if filter.URL != "" && !e.MatchesURL(filter.URL) {
		return false
	}
	if filter.Type != "" && e.Type != filter.Type {
		return false
	}
	if !filter.Deleted && e.Deleted {
		return false
	}
	return true
}

// SecretList contains a list of all SecretEntries
type SecretList struct {
	AllTags []string       `json:"all_tags"`
	Entries []*SecretEntry `json:"entries"`
}

func urlMatches(url1, url2 *url.URL) bool {
	if url1.Host == url2.Host {
		return true
	}
	return false
}
