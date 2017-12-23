package config

import "time"

type Settings struct {
	Debug             bool
	LogFormat         string        `mapstructure:"log-format"`
	LogFile           string        `mapstructure:"log-file"`
	NodeID            string        `mapstructure:"node-id"`
	UnlockTimeout     time.Duration `mapstructure:"unlock-timeout"`
	UnlockTimeoutHard bool          `mapstructure:"unlock-timeout-hard"`
	StoreURL          string        `mapstructure:"store-url"`
	StoreScheme       string        `mapstructure:"store-scheme"`
}
