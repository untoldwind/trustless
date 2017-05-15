package remote

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"github.com/leanovate/microtools/logging"
)

func dialDaemon(ctx context.Context, network, address string) (net.Conn, error) {
	location := remoteLocation()

	return net.DialUnix("unix", nil, &net.UnixAddr{
		Net:  "unix",
		Name: location,
	})
}

func remoteLocation() string {
	if xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR"); xdgRuntimeDir != "" {
		return filepath.Join(xdgRuntimeDir, "trustless", "daemon.sock")
	}
	return filepath.Join(os.Getenv("HOME"), ".trustless", "daemon.sock")
}

func remoteAvailable(logger logging.Logger) bool {
	location := remoteLocation()

	logger.Debugf("Check unix socket at: %s", location)

	if _, err := os.Stat(location); os.IsNotExist(err) {
		return false
	} else if err != nil {
		logger.ErrorErr(err)
		return false
	}
	return true
}
