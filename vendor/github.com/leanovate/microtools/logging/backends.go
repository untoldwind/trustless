package logging

import "sync"

type backendFactory func(options Options) Logger

var backendLock sync.Mutex
var backends map[string]backendFactory = map[string]backendFactory{}

func RegisterBackend(name string, factory backendFactory) {
	backendLock.Lock()
	defer backendLock.Unlock()

	backends[name] = factory
}
