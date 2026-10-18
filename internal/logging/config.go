package logging

import (
	"io"
	"log/slog"
)

type LoggingConfig struct {
	Level  *slog.LevelVar
	Writer io.Writer
	Format string
}
