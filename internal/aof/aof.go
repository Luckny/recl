// Package aof implements the AOF.
package aof

import (
	"log/slog"
	"os"
	"path"

	"github.com/Luckny/recl/config"
)

type AOF struct {
	file *os.File
}

func NewAOF(dataDir string, conf *config.AOFConfig) (*AOF, error) {
	filePath := path.Join(dataDir, conf.Filename)

	// create if not exists | write only | append if exists
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &AOF{file}, nil
}

func (a *AOF) SyncFromDisk(logger *slog.Logger) error {
	// not yet implemented
	logger.Info("Syncing AOF from disk")
	logger.Debug("AOF file path", "path", a.file.Name())
	logger.Debug("Sync Not Yet Implemented")
	return nil
}
