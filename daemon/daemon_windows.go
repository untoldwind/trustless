package daemon

import (
	"net"

	"github.com/pkg/errors"
	npipe "gopkg.in/natefinch/npipe.v2"
)

func (d *Daemon) createListener() (net.Listener, error) {
	listener, err := npipe.Listen("\\\\.\\pipe\\trustless")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create listener")
	}

	return listener, nil
}
