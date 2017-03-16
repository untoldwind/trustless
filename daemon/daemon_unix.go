package daemon

import (
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (d *Daemon) createListener() (net.Listener, error) {
	var location string

	if xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR"); xdgRuntimeDir != "" {
		location = filepath.Join(xdgRuntimeDir, "trustless", "daemon.sock")
	} else {
		location = filepath.Join(os.Getenv("HOME"), ".trustless", "daemon.sock")
	}

	if err := os.MkdirAll(filepath.Dir(location), 0700); err != nil {
		return nil, errors.Wrap(err, "Failed to create socket dir")
	}

	if err := os.Remove(location); err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "Failed to cleanup unix socket")
	}
	listener, err := net.Listen("unix", location)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create listener")
	}
	if err := os.Chmod(location, 0700); err != nil {
		return nil, errors.Wrap(err, "Failed to chmod")
	}

	return listener, nil
}
