package api

// MasterKey contains information about an available masterkey (and its status)
type MasterKey struct {
	Locked bool `json:"locked"`
}

// MasterKeyUnlock is required to unlock a masterkey (might become obsolete)
type MasterKeyUnlock struct {
	Identity
	// NODE: Tempoary measure until there is proper pinentry
	Passphrase string `json:"passphrase"`
}
