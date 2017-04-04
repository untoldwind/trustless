package remote

import (
	"context"
	"net"
	"os"
	"path/filepath"
)

func dialDaemon(ctx context.Context, network, address string) (net.Conn, error) {
	var location string

	location = filepath.Join(os.Getenv("HOME"), ".trustless", "daemon.sock")

	return net.DialUnix("unix", nil, &net.UnixAddr{
		Net:  "unix",
		Name: location,
	})
}
