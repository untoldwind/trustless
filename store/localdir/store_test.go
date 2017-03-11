package localdir_test

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"testing"

	"github.com/leanovate/microtools/logging"
	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/store"
	"github.com/untoldwind/trustless/store/model"
)

func TestStore(t *testing.T) {
	require := require.New(t)
	logger := logging.NewSimpleLoggerNull()

	tempDir, err := ioutil.TempDir(os.TempDir(), "store_test")
	require.Nil(err)
	store, err := store.NewStore("file://"+tempDir, logger)
	require.Nil(err)
	commonRingFeatures(t, store)
	commonFeatures(t, store)
}

func commonRingFeatures(t *testing.T, store store.Store) {
	require := require.New(t)

	initialRing, err := store.GetRing()
	require.Nil(err)
	require.Nil(initialRing)

	expectedRing := make([]byte, 8192)
	_, err = rand.Read(expectedRing)
	require.Nil(err)

	err = store.StoreRing(expectedRing)
	require.Nil(err)

	actualRing, err := store.GetRing()
	require.Nil(err)
	require.Equal(expectedRing, actualRing)
}

func commonFeatures(t *testing.T, store store.Store) {
	require := require.New(t)

	nodeID := "test-client"

	// Initial heads is empty
	heads, err := store.Heads()
	require.Nil(err)
	require.Len(heads, 0)

	head, err := store.GetHead(nodeID)
	require.Nil(err)
	require.Empty(head)

	// Commit first block
	commit1ID := commitRandomBlock(t, store, nodeID)

	head, err = store.GetHead(nodeID)
	require.Nil(err)
	require.Equal(commit1ID, head)

	// Commit second block
	commit2ID := commitRandomBlock(t, store, nodeID)

	head, err = store.GetHead(nodeID)
	require.Nil(err)
	require.Equal(commit2ID, head)

	commit1, err := store.GetCommit(commit1ID)
	require.Nil(err)
	require.Empty(commit1.PrevCommitID)

	commit2, err := store.GetCommit(commit2ID)
	require.Nil(err)
	require.Equal(commit2.PrevCommitID, commit1ID)

	heads, err = store.Heads()
	require.Nil(err)
	require.Len(heads, 1)
	require.Equal(nodeID, heads[0].NodeID)
	require.Equal(commit2ID, heads[0].CommitID)
}

func commitRandomBlock(t *testing.T, store store.Store, nodeID string) string {
	require := require.New(t)

	expectedBlock := make([]byte, 1024)
	_, err := rand.Read(expectedBlock)
	require.Nil(err)

	blockID, err := store.AddBlock(expectedBlock)
	require.Nil(err)
	require.NotEmpty(blockID)

	actualBlock, err := store.GetBlock(blockID)
	require.Nil(err)
	require.Equal(expectedBlock, actualBlock)

	commitID, err := store.Commit(nodeID, []model.Change{
		{Operation: model.ChangeOpAdd, BlockID: blockID},
	})
	require.Nil(err)
	require.NotEmpty(commitID)

	commit, err := store.GetCommit(commitID)
	require.Nil(err)
	require.Equal(nodeID, commit.NodeID)
	require.Len(commit.Changes, 1)
	require.Equal(model.ChangeOpAdd, commit.Changes[0].Operation)
	require.Equal(blockID, commit.Changes[0].BlockID)

	return commitID
}
