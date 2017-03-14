package secrets

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/leanovate/microtools/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/api"
)

func TestCrypt(t *testing.T) {
	require := require.New(t)
	logger := logging.NewSimpleLoggerNull()

	tempDir, err := ioutil.TempDir(os.TempDir(), "secrets_test")
	require.Nil(err)

	parameters := gopter.DefaultTestParameters()

	if !testing.Short() {
		parameters.MaxSize = 200000
	}

	secrets, err := NewSecrets("file://"+tempDir, logger)
	require.Nil(err)
	secrets.MasterKeyBits = 1024

	err = secrets.Unlock("Tester", "tester@mail.com", "12345678")
	require.Nil(err)

	properties := gopter.NewProperties(parameters)

	properties.Property("Any SecretBlock can be encrypted/decrypted", prop.ForAll(
		func(original *SecretBlock) (string, error) {
			encrypted, err := secrets.encryptSecret(original)
			if err != nil {
				return "", err
			}
			decrypted, err := secrets.decryptSecret(encrypted)
			if err != nil {
				return "", err
			}
			if !assert.Equal(t, original, decrypted) {
				return "Original != decrypted", nil
			}
			return "", nil
		},
		gen.StructPtr(reflect.TypeOf(&SecretBlock{}), map[string]gopter.Gen{
			"ID": gen.Identifier(),
			"Version": gen.StructPtr(reflect.TypeOf(&api.SecretVersion{}), map[string]gopter.Gen{
				"Name": gen.AnyString(),
			}),
		}),
	))

	properties.TestingRun(t)
}
