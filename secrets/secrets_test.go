package secrets_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/leanovate/microtools/logging"
	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/secrets"
)

func TestSecrets(t *testing.T) {
	require := require.New(t)
	logger := logging.NewSimpleLoggerNull()

	tempDir, err := ioutil.TempDir(os.TempDir(), "secrets_test")
	require.Nil(err)

	secrets, err := secrets.NewSecrets("file://"+tempDir, logger)
	require.Nil(err)
	secrets.MasterKeyBits = 1024

	require.True(secrets.IsLocked())

	err = secrets.Unlock("Tester", "tester@mail.com", "12345678")
	require.Nil(err)

	require.False(secrets.IsLocked())
}
