package commands

import (
	"errors"
	"fmt"

	"github.com/untoldwind/trustless/api"

	cli "gopkg.in/urfave/cli.v2"
)

// UnlockCommand is the commandline command to unlock the store on the daemon
var UnlockCommand = &cli.Command{
	Name:   "unlock",
	Usage:  "Unlock the store",
	Action: withDetailedErrors(unlockStore),
}

func unlockStore(ctx *cli.Context) error {
	logger := createLogger()
	client := createClient(logger)

	status, err := client.Status(createClientContext())
	if err != nil {
		return err
	}

	if !status.Initialized {
		initialUnlock, err := readInitialUnlock()
		if err != nil {
			return err
		}
		if err := client.Unlock(createClientContext(), *initialUnlock); err != nil {
			return err
		}
	} else {
		identities, err := client.Identities(createClientContext())
		if err != nil {
			return err
		}
		if len(identities) == 0 {
			return errors.New("There are no identities")
		}
		identity := identities[0]
		fmt.Printf("Name : %s\n", identity.Name)
		fmt.Printf("Email: %s\n", identity.Email)
		passphrase, err := readPassphrase()
		if err != nil {
			return err
		}
		if err := client.Unlock(createClientContext(), api.MasterKeyUnlock{
			Identity:   identity,
			Passphrase: passphrase,
		}); err != nil {
			return err
		}
	}

	status, err = client.Status(createClientContext())
	if err != nil {
		return err
	}
	reportStatus(status)
	return nil
}
