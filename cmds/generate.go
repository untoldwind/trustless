package cmds

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/untoldwind/trustless/api"
)

var GenerateFlags = struct {
	Count          int
	Words          bool
	Length         int
	CharsParameter api.CharsParameter
	WordsParameter api.WordsParameter
}{}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new password",
	Run:   withDetailedErrors(generate),
}

func generate(cmd *cobra.Command, args []string) error {
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
