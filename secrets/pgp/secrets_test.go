package pgp_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/leanovate/microtools/logging"
	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets/pgp"
)

func TestSecrets(t *testing.T) {
	require := require.New(t)
	logger := logging.NewSimpleLoggerNull()

	tempDir, err := ioutil.TempDir(os.TempDir(), "secrets_test")
	require.Nil(err)

	secrets, err := pgp.NewPGPSecrets("file://"+tempDir, "test-client", 1024, 5*time.Minute, false, logger)
	require.Nil(err)

	status, err := secrets.Status(context.Background())
	require.Nil(err)
	require.False(status.Initialized)
	require.True(status.Locked)

	err = secrets.Unlock(context.Background(), "Tester", "tester@mail.com", "12345678")
	require.Nil(err)

	status, err = secrets.Status(context.Background())
	require.Nil(err)
	require.False(status.Locked)
	require.True(status.Initialized)

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
	err = secrets.Add(context.Background(), "secret1", api.SecretTypeLogin, version1)
	require.Nil(err)

	list, err := secrets.List(context.Background(), api.SecretListFilter{})
	require.Nil(err)
	require.Len(list.AllTags, 2)
	require.Len(list.Entries, 1)
	require.Equal("secret1", list.Entries[0].ID)
	require.Equal(version1.Name, list.Entries[0].Name)
	require.Equal(version1.Tags, list.Entries[0].Tags)

	version2 := api.SecretVersion{
		Timestamp: now.Add(1 * time.Minute),
		Name:      "my-login",
		Tags:      []string{"private"},
		Properties: map[string]string{
			"url":      "https://site.com",
			"username": "tester",
			"password": "supersecret2",
		},
	}
	err = secrets.Add(context.Background(), "secret1", api.SecretTypeLogin, version2)
	require.Nil(err)

	list, err = secrets.List(context.Background(), api.SecretListFilter{})
	require.Nil(err)
	require.Len(list.AllTags, 2)
	require.Len(list.Entries, 1)
	require.Equal("secret1", list.Entries[0].ID)
	require.Equal(version2.Name, list.Entries[0].Name)
	require.Equal(version2.Tags, list.Entries[0].Tags)

	actualSecret, err := secrets.Get(context.Background(), "secret1")

	require.Nil(err)
	require.Equal("secret1", actualSecret.ID)
	require.Len(actualSecret.Versions, 2)
	require.Equal(version2.Name, actualSecret.Current.Name)
	require.Equal(version2.Tags, actualSecret.Current.Tags)
	require.Equal(version2.Properties, actualSecret.Current.Properties)
	require.Equal(version1.Tags, actualSecret.Versions[1].Tags)
	require.Equal(version1.Properties, actualSecret.Versions[1].Properties)
}
