package model

import "time"

type ChangeOp string

const (
	ChangeOpAdd    ChangeOp = "add"
	ChangeOpDelete ChangeOp = "delete"
)

type Change struct {
	Operation ChangeOp `json:"op"`
	BlockID   string   `json:"block"`
}

type Commit struct {
	NodeID       string    `json:"node"`
	PrevCommitID string    `json:"prev,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	Changes      []Change  `json:"changes"`
}
