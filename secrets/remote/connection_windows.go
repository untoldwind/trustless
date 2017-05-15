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
	conn, err := npipe.Dial("\\\\.\\pipe\\trustless")
	if err != nil {
		return false
	}
	defer conn.Close()
	return false
}
