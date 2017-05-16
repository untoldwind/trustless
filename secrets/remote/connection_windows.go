package remote

import (
	"context"
	"net"

	"github.com/leanovate/microtools/logging"

	npipe "gopkg.in/natefinch/npipe.v2"
)

func dialDaemon(ctx context.Context, network, address string) (net.Conn, error) {
	return npipe.Dial("\\\\.\\pipe\\trustless")
}

func remoteAvailable(logger logging.Logger) bool {
	listener, err := npipe.Listen("\\\\.\\pipe\\trustless")
	if err != nil {
		return true
	}
	defer listener.Close()

	return false
}
