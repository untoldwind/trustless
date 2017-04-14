package commands

import (
	"fmt"

	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets/generate"
	cli "gopkg.in/urfave/cli.v2"
)

var GenerateFlags = struct {
	Words         bool
	Length        int
	CharParameter api.CharsParameter
}{}

var GenerateCommand = &cli.Command{
	Name:   "generate",
	Usage:  "Generate a password",
	Action: withDetailedErrors(generatePassword),
	Flags: []cli.Flag{
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
			Destination: &GenerateFlags.CharParameter.IncludeUpper,
		},
		&cli.BoolFlag{
			Name:        "require-upper",
			Usage:       "Require at least one upper chars in password",
			Destination: &GenerateFlags.CharParameter.RequireUpper,
		},
		&cli.BoolFlag{
			Name:        "include-number",
			Usage:       "Include numbers in password",
			Value:       true,
			Destination: &GenerateFlags.CharParameter.IncludeNumbers,
		},
		&cli.BoolFlag{
			Name:        "require-number",
			Usage:       "Require at least one number in password",
			Destination: &GenerateFlags.CharParameter.RequireNumber,
		},
		&cli.BoolFlag{
			Name:        "include-symbols",
			Usage:       "Include symbols in password",
			Value:       true,
			Destination: &GenerateFlags.CharParameter.IncludeSymbols,
		},
		&cli.BoolFlag{
			Name:        "require-symbol",
			Usage:       "Require at least one symbol in password",
			Destination: &GenerateFlags.CharParameter.RequireSymbol,
		},
		&cli.BoolFlag{
			Name:        "exclude-similar",
			Usage:       "Exclude similar chars",
			Destination: &GenerateFlags.CharParameter.ExcludeSimilar,
		},
		&cli.BoolFlag{
			Name:        "exclude-ambigous",
			Usage:       "Exclude ambigous chars",
			Value:       true,
			Destination: &GenerateFlags.CharParameter.ExcludeAmbiguous,
		},
	},
}

func generatePassword(ctx *cli.Context) error {
	var parameters api.GenerateParameter
	if GenerateFlags.Words {

	} else {
		parameters.Chars = &GenerateFlags.CharParameter
		parameters.Chars.NumChars = GenerateFlags.Length
	}

	pwd, err := generate.Password(parameters)
	if err != nil {
		return err
	}
	fmt.Println(pwd)

	return nil
}
