package daemon

import (
	"net"
	"net/http"

	"github.com/leanovate/microtools/logging"
	"github.com/leanovate/microtools/rest"
	"github.com/leanovate/microtools/routing"
	"github.com/untoldwind/trustless/secrets"
)

// Daemon is the daomnized trustless server
type Daemon struct {
	logger   logging.Logger
	secrets  secrets.Secrets
	listener net.Listener
}

// NewDaemon creates a new trustless daeom
func NewDaemon(secrets secrets.Secrets, logger logging.Logger) *Daemon {
	return &Daemon{
		logger:  logger.WithField("package", "daemon"),
		secrets: secrets,
	}
}

// Start the daemon
func (d *Daemon) Start() error {
	server := &http.Server{
		Handler: routing.NewLoggingHandler(
			d.routeHandler(),
			d.logger.WithContext(map[string]interface{}{"type": "access"}),
		),
	}

	var err error
	d.listener, err = d.createListener()
	if err != nil {
		return err
	}

	go func() {
		d.logger.Infof("Starting daemon: %v", d.listener.Addr())
		if err := server.Serve(d.listener); err != nil {
			d.logger.ErrorErr(err)
		}
	}()

	return nil
}

// Stop the daemon
func (d *Daemon) Stop() error {
	d.logger.Info("Stopping daemon")
	if d.listener != nil {
		return d.listener.Close()
	}
	return nil
}

func (d *Daemon) routeHandler() http.Handler {
	return routing.NewRouteHandler(
		rest.ResourceMatcher(NewRootResource(d.secrets, d.logger)),
	)
}
