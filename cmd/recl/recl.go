// Package main is the main package for the recl command.
package main

import (
	"log/slog"
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
	bootstrapLogger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	cfg, err := config.LoadFile(config.DefaultFileName)
	if err != nil {
		bootstrapLogger.Error("Error loading config", "err", err)
		os.Exit(ConfigError)
	}

	logger := logging.New(&cfg.App.Logging)

	server, err := server.NewServer(&cfg.Recl, logger)
	if err != nil {
		logger.Error("Error creating server state", "err", err)
		os.Exit(GeneralError)
	}

	err = server.Run()
	if err != nil {
		logger.Error("Error running server", "err", err)
		os.Exit(GeneralError)
	}
}
