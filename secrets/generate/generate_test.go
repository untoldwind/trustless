package generate_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets/generate"
)

func TestGenerateDefault(t *testing.T) {
	require := require.New(t)

	tries := 100

	if testing.Short() {
		tries = 10
	}

	for i := 0; i < tries; i++ {
		pwd, err := generate.Password(api.GenerateParameter{})
		require.Nil(err)
		require.Len(pwd, 14)
	}
}

func TestCharWithRequired(t *testing.T) {
	require := require.New(t)

	parameters := api.GenerateParameter{
		Chars: &api.CharsParameter{
			NumChars:      20,
			RequireNumber: true,
			RequireSymbol: true,
			RequireUpper:  true,
		},
	}

	tries := 100

	if testing.Short() {
		tries = 10
	}

	for i := 0; i < tries; i++ {
		pwd, err := generate.Password(parameters)
		require.Nil(err)
		require.Len(pwd, 20)

		require.True(strings.IndexFunc(pwd, unicode.IsUpper) >= 0)
		require.True(strings.IndexFunc(pwd, unicode.IsDigit) >= 0)
		require.True(strings.IndexAny(pwd, "!-+*#_$%&/()=?{}[]()/\\'\"`-,;:.<>") >= 0)
	}
}

func TestWords(t *testing.T) {
	require := require.New(t)

	parameters := api.GenerateParameter{
		Words: &api.WordsParameter{
			NumWords: 4,
			Delim:    ".",
		},
	}
	tries := 100

	if testing.Short() {
		tries = 10
	}

	for i := 0; i < tries; i++ {
		pwd, err := generate.Password(parameters)
		require.Nil(err)

		require.True(len(pwd) > 4*3)
		require.Equal(3, strings.Count(pwd, "."))
	}
}
