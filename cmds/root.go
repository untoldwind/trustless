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
	rootCommand.PersistentFlags().String("log-file", "", "Log to file instead stdout")
	rootCommand.PersistentFlags().String("log-format", "text", "Log format to use (test, json, logstash)")
	rootCommand.PersistentFlags().Bool("debug", false, "Enable debug information")
	rootCommand.PersistentFlags().Duration("unlock-timeout", 5*time.Minute, "AUtomatic lock timeout")
	rootCommand.PersistentFlags().Bool("unlock-timeout-hard", false, "Enable hard timeout")

	rootCommand.AddCommand(infoCmd, lockCmd, unlockCmd, listCmd, estimateCmd, exportCmd, daemonCmd)

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
