package api

// Status contains relevant status information of the daemon
type Status struct {
	Initialized bool   `json:"initialized"`
	Locked      bool   `json:"locked"`
	Version     string `json:"version"`
}
