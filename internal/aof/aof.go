// Package aof implements the AOF.
package aof

import (
	"errors"
	"log/slog"
	"os"
	"path"

	"github.com/Luckny/recl/config"
	"github.com/Luckny/recl/internal/parser"
	"github.com/Luckny/recl/internal/writer"
)

type AOF struct {
	file   *os.File
	writer *writer.Writer
}

func NewAOF(dataDir string, conf *config.AOFConfig) (*AOF, error) {
	filePath := path.Join(dataDir, conf.Filename)

	// create if not exists | write only | append if exists
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &AOF{file: file, writer: writer.NewWriter(file)}, nil
}

func (a *AOF) GetSyncValueFromDisk(logger *slog.Logger) ([]parser.Value, error) {
	content, err := os.ReadFile(a.file.Name())
	if err != nil {
		return nil, err
	}
	ps := parser.NewParser(logger)
	ps.Append(content)

	var commands []parser.Value

	for {
		value, err := ps.ParseNext()
		if errors.Is(err, parser.ErrIncomplete) {
			return commands, nil
		}
		if err != nil {
			logger.Error("error parsing command", "error", err)
			return nil, err
		}
		commands = append(commands, value)
	}
}

func (a *AOF) Write(v parser.Value) {
	a.writer.Write(v)
}

func (a *AOF) Flush() {
	a.writer.Flush()
}
