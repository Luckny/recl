// Package main is the main package for the recl command.
package main

import (
	"os"

	"github.com/Luckny/recl/config"
	"github.com/Luckny/recl/internal/logging"
)

var (
	Success      = 0
	GeneralError = 1
	ConfigError  = 2
)

func main() {
	logConfig := &logging.LoggingConfig{}
	logger := logging.New(logConfig)

	if _, err := config.LoadFile(config.DefaultFileName); err != nil {
		logger.Error("Error loading config", "err", err)
		os.Exit(ConfigError)
	}

	logger.Info("Config loaded successfully")
}
