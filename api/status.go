package api

type Status struct {
	Initialized bool   `json:"initialized"`
	Locked      bool   `json:"locked"`
	Version     string `json:"version"`
}
