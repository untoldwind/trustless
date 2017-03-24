package api

import "time"

// Status contains relevant status information of the daemon
type Status struct {
	Initialized bool       `json:"initialized"`
	Locked      bool       `json:"locked"`
	AutolockAt  *time.Time `json:"autolock_at,omitempty"`
	Version     string     `json:"version"`
}
