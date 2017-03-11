package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/untoldwind/trustless/config"
)

func TestDefaultCommonConfig(t *testing.T) {
	require := require.New(t)

	config, err := config.DefaultCommonConfig()
	require.Nil(err)
	require.NotEmpty(config.StoreURL)
	require.NotEmpty(config.NodeID)
}
