package api

// Identity contains information about anyone allowed to access a trust store
type Identity struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
