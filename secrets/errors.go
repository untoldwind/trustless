package secrets

import "errors"

var SecretsLockedError = errors.New("Secrets are locked")
