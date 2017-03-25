package secrets

import (
	"context"
	"sync"
	"time"
)

// Autolocker is a helper to automatically lock the secrets store after a given
// timeout. The timeout can be set to hard-mode so that the the store will be
// locked no matter what.
type Autolocker struct {
	lock        sync.Mutex
	secrets     Secrets
	autolockAt  time.Time
	timeout     time.Duration
	hardTimeout bool
	ticker      *time.Ticker
}

// NewAutolocker creates a new Autolocker
func NewAutolocker(secrets Secrets, timeout time.Duration, hardTimeout bool) *Autolocker {
	return &Autolocker{
		secrets:     secrets,
		timeout:     timeout,
		hardTimeout: hardTimeout,
	}
}

// Start the autolock timeout (usually after an unlock)
func (a *Autolocker) Start() {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.ticker != nil {
		a.ticker.Stop()
		a.ticker = nil
	}
	a.autolockAt = time.Now().Add(a.timeout)
	a.ticker = time.NewTicker(1 * time.Second)
	go a.autolocker(a.ticker.C)
}

// Reset the timeout (will be ignored if timeout is hard-mode)
func (a *Autolocker) Reset() {
	if a.hardTimeout {
		return
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	a.autolockAt = time.Now().Add(a.timeout)
}

// Cancel the autolock timeout (usually because the store has been manually locked)
func (a *Autolocker) Cancel() {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.ticker != nil {
		a.ticker.Stop()
		a.ticker = nil
	}
}

// GetAutolockAt gets the current autolock timestamp
func (a *Autolocker) GetAutolockAt() time.Time {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.autolockAt
}

func (a *Autolocker) autolocker(ticks <-chan time.Time) {
	for t := range ticks {
		if t.After(a.GetAutolockAt()) {
			a.secrets.Lock(context.Background())
			a.Cancel()
			return
		}
	}
}
