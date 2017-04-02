package api

import "time"

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

// SecretList contains a list of all SecretEntries
type SecretList struct {
	AllTags []string       `json:"all_tags"`
	Entries []*SecretEntry `json:"entries"`
}
