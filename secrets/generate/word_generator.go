package generate

import (
	"strings"

	"github.com/untoldwind/trustless/api"
)

func generateWords(parameters *api.WordsParameter) (string, error) {
	pool, err := createRandomPool(parameters.NumWords)

	if err != nil {
		return "", err
	}

	picks := make([]string, parameters.NumWords)
	for i := 0; i < parameters.NumWords; i++ {
		picks[i] = defaultWords[pool[i]%uint64(len(defaultWords))]
	}

	return strings.Join(picks, parameters.Delim), nil
}
