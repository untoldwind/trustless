package commands

import (
	"os"

	"github.com/chzyer/readline"
)

func readPassphrase(prompt string) (string, error) {
	rl, err := readline.NewEx(&readline.Config{Prompt: "", Stdout: os.Stderr})
	if err != nil {
		return "", err
	}
	defer rl.Close()

	config := rl.GenPasswordConfig()
	config.MaskRune = '*'
	config.Prompt = prompt

	pwd, err := rl.ReadPasswordWithConfig(config)
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}
