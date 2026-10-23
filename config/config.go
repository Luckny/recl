// Package config provides the configuration for the application.
package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v2"
)

var (
	DefaultFileName  = "recl.yaml"
	DefaultDataDir   = "./data"
	DefaultPort      = 6379
	DefaultLogLevel  = "info"
	DefaultLogFormat = "json"
)

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

type ReclConfig struct {
	DataDir string    `yaml:"dataDir"`
	AOF     AOFConfig `yaml:"aof"`
}

type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type AppConfig struct {
	Port    int       `yaml:"port"`
	Logging LogConfig `yaml:"logging"`
}

type Config struct {
	App  AppConfig  `yaml:"app"`
	Recl ReclConfig `yaml:"recl"`
}

func (c *Config) ensureDefaults() error {
	// ensure data dir exists
	err := os.MkdirAll(c.Recl.DataDir, 0o755)
	if err != nil {
		return err
	}

	// ensure port is set
	if c.App.Port == 0 {
		c.App.Port = 6379
	}

	// ensure default log
	if c.App.Logging.Level == "" {
		c.App.Logging.Level = "info"
	}
	if c.App.Logging.Format == "" {
		c.App.Logging.Format = "json"
	}

	return nil
}

// LoadFile parses and validates the given YAML file into a read-only Config.
func LoadFile(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg, err := load(string(content))
	if err != nil {
		return nil, fmt.Errorf("parsing YAML file %s: %w", filename, err)
	}

	if err := cfg.ensureDefaults(); err != nil {
		return nil, err
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
