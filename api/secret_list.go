package api

import "time"

type SecretEntry struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Type      SecretType `json:"type"`
	Tags      []string   `json:"tags"`
	Timestamp time.Time  `json:"timestamp"`
}

type SecretList struct {
	AllTags []string       `json:"all_tags"`
	Entries []*SecretEntry `json:"entries"`
}
