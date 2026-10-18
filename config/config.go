// Package config provides the configuration for the application.
package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v2"
)

var DefaultFileName = "recl.yaml"

type Config struct{}

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
