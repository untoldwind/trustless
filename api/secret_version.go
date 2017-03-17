package api

import (
	"sort"
	"time"
)

type SecretAttachment struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

type SecretVersion struct {
	Timestamp   time.Time          `json:"timestamp"`
	Name        string             `json:"name"`
	Tags        []string           `json:"tags"`
	Properties  map[string]string  `json:"properties"`
	Attachments []SecretAttachment `json:"attachment"`
	Deleted     bool               `json:"deleted"`
}

type SecretVersions []*SecretVersion

func (s SecretVersions) Len() int {
	return len(s)
}

func (s SecretVersions) Less(i, j int) bool {
	return s[i].Timestamp.Before(s[j].Timestamp)
}

func (s SecretVersions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SecretVersions) Sort() {
	sort.Sort(s)
}
