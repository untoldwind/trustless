package secrets

import (
	"sync"
	"time"
)

type Autolocker struct {
	lock        sync.Mutex
	secrets     Secrets
	autolockAt  time.Time
	timeout     time.Duration
	hardTimeout bool
	ticker      *time.Ticker
}

func NewAutolocker(secrets Secrets, timeout time.Duration, hardTimeout bool) *Autolocker {
	return &Autolocker{
		secrets:     secrets,
		timeout:     timeout,
		hardTimeout: hardTimeout,
	}
}

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

func (a *Autolocker) Reset() {
	if a.hardTimeout {
		return
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	a.autolockAt = time.Now().Add(a.timeout)
}

func (a *Autolocker) Cancel() {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.ticker != nil {
		a.ticker.Stop()
		a.ticker = nil
	}
}

func (a *Autolocker) GetAutolockAt() time.Time {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.autolockAt
}

func (a *Autolocker) autolocker(ticks <-chan time.Time) {
	for t := range ticks {
		if t.After(a.GetAutolockAt()) {
			a.secrets.Lock()
			a.Cancel()
			return
		}
	}
}
