package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func ReadConfig(configFile string, logger logging.Logger) (*CommonConfig, error) {
	raw, err := ioutil.ReadFile(configFile)

	if os.IsNotExist(err) {
		logger.Warnf("Configuration %s does not exists, create default")
		defaultConfig, err := DefaultCommonConfig()
		if err != nil {
			return nil, err
		}
		if err := WriteClientConfig(configFile, defaultConfig); err != nil {
			return nil, err
		}
		return defaultConfig, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "read config file failed")
	}

	var config CommonConfig
	if strings.HasSuffix(configFile, ".json") {
		if err := json.Unmarshal(raw, &config); err != nil {
			return nil, errors.Wrap(err, "config.ClientConfig json unmarshal failed")
		}
	} else {
		if err := yaml.Unmarshal(raw, &config); err != nil {
			return nil, errors.Wrap(err, "config.ClientConfig yaml unmarshal failed")
		}
	}
	return &config, nil
}

func WriteClientConfig(configFile string, config *CommonConfig) error {
	var raw []byte
	var err error
	if strings.HasSuffix(configFile, ".json") {
		raw, err = json.Marshal(config)
		if err != nil {
			return errors.Wrap(err, "config.CommonConfig json marshal failed")
		}
	} else {
		raw, err = yaml.Marshal(config)
		if err != nil {
			return errors.Wrap(err, "config.CommonConfig yaml marshal failed")
		}
	}
	if err := os.MkdirAll(filepath.Dir(configFile), 0700); err != nil {
		return errors.Wrap(err, "creating config file directory failed")
	}
	if err := ioutil.WriteFile(configFile, raw, 0600); err != nil {
		return errors.Wrap(err, "write config file failed")
	}
	return nil
}
