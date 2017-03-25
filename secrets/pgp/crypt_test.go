package pgp

import (
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/leanovate/microtools/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func TestCrypt(t *testing.T) {
	require := require.New(t)
	logger := logging.NewSimpleLoggerNull()

	tempDir, err := ioutil.TempDir(os.TempDir(), "secrets_test")
	require.Nil(err)

	parameters := gopter.DefaultTestParameters()

	_secrets, err := NewPGPSecrets("file://"+tempDir, "testNode", 1024, 5*time.Minute, false, logger)
	require.Nil(err)
	pgpSecrets := _secrets.(*pgpSecrets)

	err = pgpSecrets.Unlock(context.Background(), "Tester", "tester@mail.com", "12345678")
	require.Nil(err)

	properties := gopter.NewProperties(parameters)

	properties.Property("Any SecretBlock can be encrypted/decrypted", prop.ForAll(
		func(original *secrets.SecretBlock) (string, error) {
			encrypted, err := pgpSecrets.encryptSecret(original)
			if err != nil {
				return "", err
			}
			decrypted, err := pgpSecrets.decryptSecret(encrypted)
			if err != nil {
				return "", err
			}
			if !assert.Equal(t, original, decrypted) {
				return "Original != decrypted", nil
			}
			return "", nil
		},
		gen.StructPtr(reflect.TypeOf(&secrets.SecretBlock{}), map[string]gopter.Gen{
			"ID": gen.Identifier(),
			"Version": gen.Struct(reflect.TypeOf(api.SecretVersion{}), map[string]gopter.Gen{
				"Name":      gen.AnyString(),
				"Timestamp": gen.Time(),
				"Tags":      gen.SliceOf(gen.AnyString()),
				"Attachments": gen.SliceOf(gen.Struct(reflect.TypeOf(api.SecretAttachment{}), map[string]gopter.Gen{
					"Name": gen.AnyString(),
				})),
			}),
		}),
	))

	properties.TestingRun(t)
}
