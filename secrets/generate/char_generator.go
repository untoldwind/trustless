package generate

import (
	"crypto/rand"
	"encoding/binary"
	"strings"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

const lowers = "abcdefghijklmnopqrstuvwxyz"
const uppers = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers = "0123456789"
const symbols = "!-+*#_$%&/()=?{}[]()/\\'\"`-,;:.<>"
const ambigousChars = "{}[]()/\\'\"`-,;:.<>"
const similarChars = "QO01lIB8S5G62ZUV"

var defaultCharsParameter = &api.CharsParameter{
	NumChars:         14,
	IncludeUpper:     true,
	IncludeNumbers:   true,
	IncludeSymbols:   true,
	RequireNumber:    false,
	RequireSymbol:    false,
	RequireUpper:     false,
	ExcludeSimilar:   false,
	ExcludeAmbiguous: true,
}

func generateChars(parameters *api.CharsParameter) (string, error) {
	pool, err := createRandomPool(2 * parameters.NumChars)

	if err != nil {
		return "", err
	}

	picks := make([]rune, parameters.NumChars)
	idx := 0

	if parameters.RequireUpper {
		picks[idx] = pickCharFrom(parameters, uppers, pool[idx])
		idx++
	}
	if parameters.RequireNumber {
		picks[idx] = pickCharFrom(parameters, numbers, pool[idx])
		idx++
	}
	if parameters.RequireSymbol {
		picks[idx] = pickCharFrom(parameters, symbols, pool[idx])
		idx++
	}

	if idx < parameters.NumChars {
		candidates := createBaseSet(parameters)
		for ; idx < parameters.NumChars; idx++ {
			picks[idx] = candidates[pool[idx]%uint64(len(candidates))]
		}
	}

	// shuffle the picks with second half of the pool
	pool = pool[parameters.NumChars:]
	result := make([]rune, parameters.NumChars)
	for i := 0; i < parameters.NumChars; i++ {
		pickIdx := pool[i] % uint64(len(picks))
		result[i] = picks[pickIdx]
		picks = append(picks[:pickIdx], picks[pickIdx+1:]...)
	}
	return string(result), nil
}

func createBaseSet(parameters *api.CharsParameter) []rune {
	candidates := make([]rune, 0, len(lowers)+len(uppers)+len(numbers)+len(symbols))
	candidates = filterSet(candidates, parameters, lowers)
	if parameters.IncludeUpper {
		candidates = filterSet(candidates, parameters, uppers)
	}
	if parameters.IncludeNumbers {
		candidates = filterSet(candidates, parameters, numbers)
	}
	if parameters.IncludeSymbols {
		candidates = filterSet(candidates, parameters, symbols)
	}
	return candidates
}

func filterSet(candidates []rune, parameters *api.CharsParameter, set string) []rune {
	for _, ch := range set {
		if parameters.ExcludeSimilar && strings.ContainsRune(similarChars, ch) {
			continue
		}
		if parameters.ExcludeAmbiguous && strings.ContainsRune(ambigousChars, ch) {
			continue
		}
		candidates = append(candidates, ch)
	}
	return candidates
}

func pickCharFrom(parameters *api.CharsParameter, set string, pick uint64) rune {
	candidateSet := make([]rune, 0, len(set))
	candidateSet = filterSet(candidateSet, parameters, set)

	return candidateSet[pick%uint64(len(candidateSet))]
}

// We use uint64 to ensure a (mostly) equal distribution on modulo operations.
// E.g. when picking via from a set of 240 candidates via
// (random uint8) % 240
// results in the first 16 candidates to be picked twice as often
// as the others
func createRandomPool(size int) ([]uint64, error) {
	pool := make([]uint64, size)
	if err := binary.Read(rand.Reader, binary.BigEndian, &pool); err != nil {
		return nil, errors.Wrap(err, "Failed to create random pool")
	}
	return pool, nil
}
