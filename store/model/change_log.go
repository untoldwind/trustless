package model

// ChangeLog contains all changes performed by a specific node
type ChangeLog struct {
	NodeID  string
	Changes []Change
}
