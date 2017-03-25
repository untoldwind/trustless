package api

import (
	"sort"
	"time"
)

// SecretAttachment is an attachment to a secret
type SecretAttachment struct {
	Name     string `json:"name"`
	MimeType string `json:"mimeType"`
	Content  []byte `json:"content"`
}

// SecretVersion is a concret version of a secret
type SecretVersion struct {
	Timestamp   time.Time          `json:"timestamp"`
	Name        string             `json:"name"`
	Tags        []string           `json:"tags"`
	URLs        []string           `json:"urls"`
	Properties  map[string]string  `json:"properties"`
	Attachments []SecretAttachment `json:"attachment"`
	Deleted     bool               `json:"deleted"`
}

// SecretVersions is a sequence of SecretVersion
type SecretVersions []*SecretVersion

func (s SecretVersions) Len() int {
	return len(s)
}

func (s SecretVersions) Less(i, j int) bool {
	return s[i].Timestamp.After(s[j].Timestamp)
}

func (s SecretVersions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Sort the versions, newest will be first
func (s SecretVersions) Sort() {
	sort.Sort(s)
}
