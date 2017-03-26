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
	changeLogs, err := store.ChangeLogs()
	require.Nil(err)
	require.Len(changeLogs, 0)

	// Commit first block
	blockID := commitRandomBlock(t, store, nodeID)

	changeLogs, err = store.ChangeLogs()
	require.Nil(err)
	require.Len(changeLogs, 1)
	myChangeLog := changeLogs[0]
	require.Equal(nodeID, myChangeLog.NodeID)
	require.Len(myChangeLog.Changes, 1)
	require.Equal(model.ChangeOpAdd, myChangeLog.Changes[0].Operation)
	require.Equal(blockID, myChangeLog.Changes[0].BlockID)

	// Commit second block
	block2ID := commitRandomBlock(t, store, nodeID)

	changeLogs, err = store.ChangeLogs()
	require.Nil(err)
	require.Len(changeLogs, 1)
	myChangeLog = changeLogs[0]
	require.Equal(nodeID, myChangeLog.NodeID)
	require.Len(myChangeLog.Changes, 2)
	require.Equal(model.ChangeOpAdd, myChangeLog.Changes[0].Operation)
	require.Equal(blockID, myChangeLog.Changes[0].BlockID)
	require.Equal(model.ChangeOpAdd, myChangeLog.Changes[1].Operation)
	require.Equal(block2ID, myChangeLog.Changes[1].BlockID)
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

	err = store.Commit(nodeID, []model.Change{
		{Operation: model.ChangeOpAdd, BlockID: blockID},
	})
	require.Nil(err)

	return blockID
}
