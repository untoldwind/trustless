package cmds

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/davecgh/go-spew/spew"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/untoldwind/trustless/config"
)

var cfgFile string

var rootCommand = &cobra.Command{
	Use:     "trustless",
	Short:   "Password storage",
	Long:    "Password storage (if you have less trust in all the others)",
	Version: config.Version(),
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		showError(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	exportCmd.Flags().StringVar(&exportFlags.format, "format", "json", "Export format (json, keepass)")

	rootCommand.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default ./trustless.toml, $HOME/.trustless/trustless.toml")
	rootCommand.PersistentFlags().String("node-id", "", "ID of this node")
	rootCommand.PersistentFlags().String("store-url", "", "URL of the store to open (only file:// URLs are supported)")
	rootCommand.PersistentFlags().String("store-scheme", "openpgp", "Storage scheme (openpgp, openpgp+scrypt)")
	rootCommand.PersistentFlags().String("log-file", "", "Log to file instead stdout")
	rootCommand.PersistentFlags().String("log-format", "text", "Log format to use (test, json, logstash)")
	rootCommand.PersistentFlags().Bool("debug", false, "Enable debug information")
	rootCommand.PersistentFlags().Duration("unlock-timeout", 5*time.Minute, "AUtomatic lock timeout")
	rootCommand.PersistentFlags().Bool("unlock-timeout-hard", false, "Enable hard timeout")

	generateCmd.PersistentFlags().IntVar(&GenerateFlags.Count, "count", 10, "Number of passwords to generate")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.Words, "words", false, "Generate password based on words")
	generateCmd.PersistentFlags().IntVar(&GenerateFlags.Length, "length", 16, "Desired password length (or number of words)")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.IncludeUpper, "include-upper", true, "Include upper chars in password")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.RequireUpper, "require-upper", false, "Require at least one upper chars in password")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.IncludeNumbers, "include-number", true, "Include numbers in password")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.RequireNumber, "require-number", false, "Require at least one number in password")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.IncludeSymbols, "include-symbols", false, "Include symbols in password")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.RequireSymbol, "require-symbol", false, "Require at least one symbol in password")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.ExcludeSimilar, "exclude-similar", false, "Exclude similar chars")
	generateCmd.PersistentFlags().BoolVar(&GenerateFlags.CharsParameter.ExcludeSimilar, "exclude-ambigous", true, "Exclude ambigous chars")
	generateCmd.PersistentFlags().StringVar(&GenerateFlags.WordsParameter.Delim, "delim", ".", "Delimiter for words")

	rootCommand.AddCommand(infoCmd, lockCmd, unlockCmd, listCmd, estimateCmd, generateCmd, exportCmd, daemonCmd)

	viper.BindPFlags(rootCommand.PersistentFlags())
}

func initConfig() {
	viper.SetEnvPrefix("TRUSTLESS")
	viper.AutomaticEnv()

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			showError(err)
		}
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".config"))
		viper.SetConfigName("trustless")
	}

	if err := viper.ReadInConfig(); err != nil {
		if err := writeDefaultConfig(); err != nil {
			showError(err)
		}
	}

	if err := viper.Unmarshal(&cmdSettings); err != nil {
		showError(err)
	}

	if cmdSettings.Debug {
		fmt.Fprintln(os.Stderr, "---- Settings [snip]----")
		spew.Fdump(os.Stderr, cmdSettings)
		fmt.Fprintln(os.Stderr, "---- Settings [snap] ----")
	}
}
