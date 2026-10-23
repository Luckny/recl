// Package config provides the configuration for the application.
package config

import (
	"fmt"
	"log/slog"
	"os"

	"go.yaml.in/yaml/v2"
)

var DefaultFileName = "recl.yaml"

type AOFSyncMode string

const (
	AOFSyncModeEverysec AOFSyncMode = "everysec"
	AOFSyncModeAlways   AOFSyncMode = "always"
	AOFSyncModeNo       AOFSyncMode = "no"
)

type AOFConfig struct {
	Enabled  bool        `yaml:"enabled"`
	Filename string      `yaml:"filename"`
	Sync     AOFSyncMode `yaml:"sync"`
}

type Config struct {
	DataDir string    `yaml:"dataDir"`
	AOF     AOFConfig `yaml:"aof"`
}

// LoadFile parses and validates the given YAML file into a read-only Config.
func LoadFile(filename string, logger *slog.Logger) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := load(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing YAML file %s: %w", filename, err)
	}

	// Force DataDir
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}

	return cfg, nil
}

// Load loads the configuration from the environment.
func load(content string) (*Config, error) {
	cfg := &Config{}

	err := yaml.UnmarshalStrict([]byte(content), cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
