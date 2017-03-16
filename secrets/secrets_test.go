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

	require.True(secrets.IsInitialized())
	require.True(secrets.IsLocked())

	err = secrets.Unlock("Tester", "tester@mail.com", "12345678")
	require.Nil(err)

	require.False(secrets.IsLocked())
	require.False(secrets.IsInitialized())

	now := time.Now().Add(-1 * time.Minute)
	version1 := api.SecretVersion{
		Timestamp: now,
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

	list, err := secrets.List()
	require.Nil(err)
	require.Len(list.AllTags, 2)
	require.Len(list.Entries, 1)
	require.Equal("secret1", list.Entries[0].ID)
	require.Equal(version1.Name, list.Entries[0].Name)
	require.Equal(version1.Tags, list.Entries[0].Tags)

	now.Add(1 * time.Minute)
	version2 := api.SecretVersion{
		Timestamp: now,
		Name:      "my-login",
		Tags:      []string{"private"},
		Properties: map[string]string{
			"url":      "https://site.com",
			"username": "tester",
			"password": "supersecret2",
		},
	}
	err = secrets.Add("secret1", api.SecretTypeLogin, version2)
	require.Nil(err)

	list, err = secrets.List()
	require.Nil(err)
	require.Len(list.AllTags, 2)
	require.Len(list.Entries, 1)
	require.Equal("secret1", list.Entries[0].ID)
	require.Equal(version2.Name, list.Entries[0].Name)
	require.Equal(version2.Tags, list.Entries[0].Tags)
}
