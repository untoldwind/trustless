package secrets_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/secrets"
)

type fakeSecrets struct {
	lock sync.Mutex
	secrets.Secrets
	locked bool
}

func (f *fakeSecrets) isLocked() bool {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.locked
}

func (f *fakeSecrets) Lock(ctx context.Context) error {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.locked = true
	return nil
}

func TestAutolocker(t *testing.T) {
	require := require.New(t)

	fakeSecrets := &fakeSecrets{}
	require.False(fakeSecrets.isLocked())

	autolocker := secrets.NewAutolocker(fakeSecrets, 2*time.Second, false)

	autolocker.Start()
	first := autolocker.GetAutolockAt()
	require.True(first.After(time.Now()))

	time.Sleep(1 * time.Second)
	autolocker.Reset()
	second := autolocker.GetAutolockAt()
	require.True(second.After(first))

	if testing.Short() {
		t.SkipNow()
	}

	time.Sleep(4 * time.Second)
	require.True(fakeSecrets.isLocked())

	fakeSecrets.locked = false
	autolocker.Start()
	autolocker.Cancel()

	time.Sleep(4 * time.Second)
	require.False(fakeSecrets.isLocked())
}
