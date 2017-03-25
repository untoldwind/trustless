package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

func readPassphrase() (string, error) {
	rl, err := readline.New("")
	if err != nil {
		return "", errors.Wrap(err, "Failed to create readline")
	}
	defer rl.Close()

	config := rl.GenPasswordConfig()
	config.MaskRune = '*'
	config.Prompt = boldRed("Master Passphrase: ")

	pwd, err := rl.ReadPasswordWithConfig(config)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read passphrase")
	}
	return string(pwd), nil
}

func readInitialUnlock() (*api.MasterKeyUnlock, error) {
	fmt.Println("Store has not been initialized yet.")
	fmt.Println("Provide initial identity")
	fmt.Println()

	rl, err := readline.New("")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create readline")
	}
	defer rl.Close()

	rl.SetPrompt("Name : ")
	name, err := rl.Readline()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read name")
	}

	rl.SetPrompt("Email: ")
	email, err := rl.Readline()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read email")
	}

	config := rl.GenPasswordConfig()
	config.MaskRune = '*'
	config.Prompt = boldRed("Master Passphrase          : ")

	pwd1, err := rl.ReadPasswordWithConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read passphrase")
	}

	config.Prompt = boldRed("Master Passphrase (confirm): ")

	pwd2, err := rl.ReadPasswordWithConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read passphrase")
	}

	if string(pwd1) != string(pwd2) {
		return nil, errors.New("Passphrase do not match")
	}

	return &api.MasterKeyUnlock{
		Identity: api.Identity{
			Name:  name,
			Email: email,
		},
		Passphrase: string(pwd1),
	}, nil
}

func confirm(prompt string) bool {
	fmt.Println()
	rl, err := readline.New(prompt + " [y/N]: ")
	if err != nil {
		return false
	}
	defer rl.Close()

	answer, err := rl.Readline()
	if err != nil {
		return false
	}
	answer = strings.ToLower(answer)
	return answer == "y" || answer == "yes"
}

func reportStatus(status *api.Status) {
	fmt.Println()
	if status.Initialized {
		fmt.Printf("Store is        : %s\n", green("Initialized"))
	} else {
		fmt.Printf("Store is        : %s\n", boldRed("Not initialized"))
	}
	if status.Locked {
		fmt.Printf("Store is        : %s\n", green("Locked"))
	} else {
		fmt.Printf("Store is        : %s\n", yellow("Unlocked"))
		if status.AutolockAt != nil {
			timeout := status.AutolockAt.Sub(time.Now())
			timeout = (timeout / time.Second) * time.Second

			fmt.Printf("Will autolock in: %v\n", timeout)
		}
	}
}
