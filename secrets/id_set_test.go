package secrets_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/untoldwind/trustless/secrets"
)

type CommitSetHolder struct {
	Commits secrets.IDSet
}

func TestCommitSet(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("any set of commits can be added", prop.ForAll(
		func(commitIDs []string) string {
			commitSet := &secrets.IDSet{}

			for _, commitID := range commitIDs {
				commitSet.Add(commitID)
			}
			for _, commitID := range commitIDs {
				if !commitSet.Contains(commitID) {
					return fmt.Sprintf("%s missing", commitID)
				}
			}
			return ""
		},
		gen.SliceOf(gen.Identifier()),
	))
	properties.Property("any set of commits can be json serialized", prop.ForAll(
		func(commitIDs []string) (string, error) {
			holder := &CommitSetHolder{
				Commits: secrets.IDSet{},
			}

			for _, commitID := range commitIDs {
				holder.Commits.Add(commitID)
			}

			serialized, err := json.Marshal(&holder)
			if err != nil {
				return "", err
			}

			var actual CommitSetHolder
			err = json.Unmarshal(serialized, &actual)
			if err != nil {
				return "", err
			}

			if !actual.Commits.Equals(holder.Commits) {
				return "Sets do not match", nil
			}
			return "", nil
		},
		gen.SliceOf(gen.Identifier()),
	))

	properties.TestingRun(t)
}
