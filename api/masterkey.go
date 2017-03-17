package api

type MasterKey struct {
	Locked bool `json:"locked"`
}

type MasterKeyUnlock struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	// NODE: Tempoary measure until there is proper pinentry
	Passphrase string `json:"passphrase"`
}
