// Package server implements the server.
package server

import (
	"log/slog"
	"os"

	"github.com/Luckny/recl/config"
	"github.com/Luckny/recl/internal/aof"
)

type ServerState struct {
	AOF *aof.AOF
}

func NewServerState(conf *config.Config, logger *slog.Logger) (*ServerState, error) {
	// create data dir if not created
	err := os.MkdirAll(conf.DataDir, 0o755)
	if err != nil {
		return nil, err
	}

	var state ServerState
	if conf.AOF.Enabled {
		aof, err := aof.NewAOF(conf.DataDir, &conf.AOF)
		if err != nil {
			return nil, err
		}

		aof.SyncFromDisk(logger)
		state.AOF = aof
	}
	return &state, nil
}
