package secrets

import "errors"

// ErrSecretsLocked occurs if the secret store is locked
var ErrSecretsLocked = errors.New("Secrets are locked")

// ErrSecretNotFound occurs if a secret could not be found
var ErrSecretNotFound = errors.New("Secret not found")
