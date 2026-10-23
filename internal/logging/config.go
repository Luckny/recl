package logging

import (
	"fmt"
	"log/slog"
	"strings"
)

type Level struct {
	lvl slog.LevelVar
}

// Level returns the value of the logging level as an slog.Level.
func (l *Level) Level() slog.Level {
	return l.lvl.Level()
}

// String returns the current level.
func (l *Level) String() string {
	switch l.lvl.Level() {
	case slog.LevelDebug:
		return "debug"
	case slog.LevelInfo:
		return "info"
	case slog.LevelWarn:
		return "warn"
	case slog.LevelError:
		return "error"
	default:
		return ""
	}
}

func (l *Level) Set(s string) error {
	switch strings.ToLower(s) {
	case "debug":
		l.lvl.Set(slog.LevelDebug)
	case "info":
		l.lvl.Set(slog.LevelInfo)
	case "warn":
		l.lvl.Set(slog.LevelWarn)
	case "error":
		l.lvl.Set(slog.LevelError)
	default:
		return fmt.Errorf("unrecognized log level %s", s)
	}
	return nil
}
