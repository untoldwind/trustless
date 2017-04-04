package remote

import (
	"context"
	"net"

	npipe "gopkg.in/natefinch/npipe.v2"
)

func dialDaemon(ctx context.Context, network, address string) (net.Conn, error) {
	return npipe.Dial("\\\\.\\pipe\\trustless")
}
