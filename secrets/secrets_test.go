package secrets_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/leanovate/microtools/logging"
	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func TestSecrets(t *testing.T) {
	require := require.New(t)
	logger := logging.NewSimpleLoggerNull()

	tempDir, err := ioutil.TempDir(os.TempDir(), "secrets_test")
	require.Nil(err)

	secrets, err := secrets.NewSecrets("file://"+tempDir, "test-client", logger)
	require.Nil(err)
	secrets.MasterKeyBits = 1024

	require.True(secrets.IsLocked())

	err = secrets.Unlock("Tester", "tester@mail.com", "12345678")
	require.Nil(err)

	require.False(secrets.IsLocked())

	version1 := api.SecretVersion{
		Timestamp: time.Now(),
		Name:      "my-login",
		Tags:      []string{"web", "private"},
		Properties: map[string]string{
			"url":      "https://site.com",
			"username": "tester",
			"password": "supersecret",
		},
	}
	err = secrets.Add("secret1", api.SecretTypeLogin, version1)
	require.Nil(err)
}
