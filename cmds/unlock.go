package cmds

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock the store",
	Run:   withDetailedErrors(unlockStoreCmd),
}

func unlockStoreCmd(cmd *cobra.Command, args []string) error {
	logger := createLogger()
	client := createRemote(logger)

	status, err := unlockStore(client)
	if err != nil {
		return err
	}

	reportStatus(status)
	return nil
}

func unlockStore(client secrets.Secrets) (*api.Status, error) {
	status, err := client.Status(createClientContext())
	if err != nil {
		return nil, err
	}

	if !status.Locked {
		return status, err
	}
	if !status.Initialized {
		initialUnlock, err := readInitialUnlock()
		if err != nil {
			return nil, err
		}
		if err := client.Unlock(createClientContext(), initialUnlock.Name, initialUnlock.Email, initialUnlock.Passphrase); err != nil {
			return nil, err
		}
	} else {
		identities, err := client.Identities(createClientContext())
		if err != nil {
			return nil, err
		}
		if len(identities) == 0 {
			return nil, errors.New("There are no identities")
		}
		identity := identities[0]
		fmt.Fprintf(os.Stderr, "Name : %s\n", identity.Name)
		fmt.Fprintf(os.Stderr, "Email: %s\n", identity.Email)
		passphrase, err := readPassphrase("Master Passphrase: ")
		defer passphrase.Destroy()
		if err != nil {
			return nil, err
		}
		if err := client.Unlock(createClientContext(), identity.Name, identity.Email, string(passphrase.Buffer())); err != nil {
			return nil, err
		}
	}

	status, err = client.Status(createClientContext())
	if err != nil {
		return nil, err
	}

	return status, nil
}
