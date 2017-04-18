package commands

import (
	"context"
	"fmt"

	"github.com/untoldwind/trustless/api"
	cli "gopkg.in/urfave/cli.v2"
)

var GenerateFlags = struct {
	Count          int
	Words          bool
	Length         int
	CharsParameter api.CharsParameter
	WordsParameter api.WordsParameter
}{}

var GenerateCommand = &cli.Command{
	Name:   "generate",
	Usage:  "Generate a password",
	Action: withDetailedErrors(generatePassword),
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:        "count",
			Usage:       "Number of passwords to generate (to pick the nicest from)",
			Value:       10,
			Destination: &GenerateFlags.Count,
		},
		&cli.BoolFlag{
			Name:        "words",
			Usage:       "Generate password based on words",
			Destination: &GenerateFlags.Words,
		},
		&cli.IntFlag{
			Name:        "length",
			Usage:       "Password length or word count",
			Value:       14,
			Destination: &GenerateFlags.Length,
		},
		&cli.BoolFlag{
			Name:        "include-upper",
			Usage:       "Include upper chars in password",
			Value:       true,
			Destination: &GenerateFlags.CharsParameter.IncludeUpper,
		},
		&cli.BoolFlag{
			Name:        "require-upper",
			Usage:       "Require at least one upper chars in password",
			Destination: &GenerateFlags.CharsParameter.RequireUpper,
		},
		&cli.BoolFlag{
			Name:        "include-number",
			Usage:       "Include numbers in password",
			Value:       true,
			Destination: &GenerateFlags.CharsParameter.IncludeNumbers,
		},
		&cli.BoolFlag{
			Name:        "require-number",
			Usage:       "Require at least one number in password",
			Destination: &GenerateFlags.CharsParameter.RequireNumber,
		},
		&cli.BoolFlag{
			Name:        "include-symbols",
			Usage:       "Include symbols in password",
			Value:       true,
			Destination: &GenerateFlags.CharsParameter.IncludeSymbols,
		},
		&cli.BoolFlag{
			Name:        "require-symbol",
			Usage:       "Require at least one symbol in password",
			Destination: &GenerateFlags.CharsParameter.RequireSymbol,
		},
		&cli.BoolFlag{
			Name:        "exclude-similar",
			Usage:       "Exclude similar chars",
			Destination: &GenerateFlags.CharsParameter.ExcludeSimilar,
		},
		&cli.BoolFlag{
			Name:        "exclude-ambigous",
			Usage:       "Exclude ambigous chars",
			Value:       true,
			Destination: &GenerateFlags.CharsParameter.ExcludeAmbiguous,
		},
		&cli.StringFlag{
			Name:        "delim",
			Usage:       "Delimeter for words",
			Value:       ".",
			Destination: &GenerateFlags.WordsParameter.Delim,
		},
	},
}

func generatePassword(ctx *cli.Context) error {
	logger := createLogger()
	client := createRemote(logger)

	var parameters api.GenerateParameter
	if GenerateFlags.Words {
		parameters.Words = &GenerateFlags.WordsParameter
		parameters.Words.NumWords = GenerateFlags.Length
	} else {
		parameters.Chars = &GenerateFlags.CharsParameter
		parameters.Chars.NumChars = GenerateFlags.Length
	}

	for i := 0; i < GenerateFlags.Count; i++ {
		pwd, err := client.GeneratePassword(context.Background(), parameters)
		if err != nil {
			return err
		}
		fmt.Println(pwd)
	}

	return nil
}
