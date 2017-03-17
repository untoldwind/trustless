package pgp

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestPadding(t *testing.T) {
	parameters := gopter.DefaultTestParameters()

	if !testing.Short() {
		parameters.MaxSize = 200000
	}

	properties := gopter.NewProperties(parameters)

	properties.Property("any non-zero content can be padded/unpadded", prop.ForAll(
		func(content []byte) (string, error) {
			padded, err := padBlock(content)
			if err != nil {
				return "", err
			}
			if len(padded) < 0x2000 {
				return fmt.Sprintf("Too short: %d", len(padded)), nil
			}
			unpadded, err := unpadBlock(padded)
			if err != nil {
				return "", err
			}

			expected := hex.EncodeToString(content)
			actual := hex.EncodeToString(unpadded)

			if expected != actual {
				return fmt.Sprintf("%s != %s", expected, actual), nil
			}
			return "", nil
		},
		gen.SliceOf(gen.UInt8Range(1, 255)),
	))

	properties.TestingRun(t)
}
