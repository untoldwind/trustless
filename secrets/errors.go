package secrets

import "errors"

var SecretsLockedError = errors.New("Secrets are locked")
var SecretNotFound = errors.New("Secret not found")
