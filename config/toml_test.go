package config

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteConfigToTomlFile(t *testing.T) {
	dir, err := os.Getwd()
	require.Nil(t, err)
	fileName := filepath.Join(dir, "logger.toml")
	cfg := DefaultLoggerConfig()
	err = WriteConfigToTomlFile(fileName, &cfg)
	require.Nil(t, err)

	_ = os.Remove(fileName)
}

func TestGetLoggerConfigFromPath(t *testing.T) {
	dir, err := os.Getwd()
	require.Nil(t, err)
	fileName := filepath.Join(dir, "logger.toml")
	cfg := DefaultLoggerConfig()
	err = WriteConfigToTomlFile(fileName, &cfg)
	require.Nil(t, err)

	defer os.Remove(fileName)

	cfgR, err := GetLoggerConfigFromPath(dir, nil)
	require.Nil(t, err)

	require.Equal(t, cfg, *cfgR)
}

func TestGetLoggerConfigFromFile(t *testing.T) {
	dir, err := os.Getwd()
	require.Nil(t, err)
	fileName := filepath.Join(dir, "logger.toml")
	cfg := DefaultLoggerConfig()
	err = WriteConfigToTomlFile(fileName, &cfg)
	require.Nil(t, err)

	defer os.Remove(fileName)

	cfgR, err := GetLoggerConfigFromFile(fileName, nil)
	require.Nil(t, err)

	require.Equal(t, cfg, *cfgR)
}
