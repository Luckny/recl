// Package main is the main package for the recl command.
package main

import (
	"os"

	"github.com/Luckny/recl/config"
	"github.com/Luckny/recl/internal/logging"
	"github.com/Luckny/recl/internal/server"
)

var (
	Success      = 0
	GeneralError = 1
	ConfigError  = 2
)

func main() {
	logConfig := &logging.LoggingConfig{}
	logger := logging.New(logConfig)

	cfg, err := config.LoadFile(config.DefaultFileName, logger)
	if err != nil {
		logger.Error("Error loading config", "err", err)
		os.Exit(ConfigError)
	}

	_, err = server.NewServerState(cfg, logger)
	if err != nil {
		logger.Error("Error creating server state", "err", err)
		os.Exit(GeneralError)
	}

	// server.Run(state)
}
