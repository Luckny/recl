// Package logging is a package that provides a logger for the application.
package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Luckny/recl/config"
)

const (
	reservedKeyPrefix = "logged_"
)

var defaultWriter = os.Stderr

func New(config *config.LogConfig) *slog.Logger {
	logLevel := &Level{}
	logLevel.Set(config.Level)

	logHandlerOpts := &slog.HandlerOptions{
		Level:     logLevel.Level(),
		AddSource: true,
	}

	logHandlerOpts.ReplaceAttr = replaceAttrFunc(logLevel)

	if config.Format != "" && config.Format == "json" {
		return slog.New(slog.NewJSONHandler(defaultWriter, logHandlerOpts))
	}

	return slog.New(slog.NewTextHandler(defaultWriter, logHandlerOpts))
}

func replaceAttrFunc(level *Level) func(groups []string, a slog.Attr) slog.Attr {
	return func(_ []string, a slog.Attr) slog.Attr {
		key := a.Key
		switch key {
		case slog.TimeKey, "ts":
			if t, ok := a.Value.Any().(time.Time); ok {
				a.Key = "ts"

				// This timestamp format differs from RFC3339Nano by using .000 instead
				// of .999999999 which changes the timestamp from 9 variable to 3 fixed
				// decimals (.130 instead of .130987456).
				a.Value = slog.StringValue(t.UTC().Format("2006-01-02T15:04:05.000Z07:00"))
			} else {
				// If we can't cast the any from the value to a
				// time.Time, it means the caller logged
				// another attribute with a key of `ts`.
				// Prevent duplicate keys (necessary for proper
				// JSON) by renaming the key to `logged_ts`.
				a.Key = reservedKeyPrefix + key
			}
		case slog.SourceKey, "caller":
			if src, ok := a.Value.Any().(*slog.Source); ok {
				a.Key = "caller"
				switch level.String() {
				case "debug":
					a.Value = slog.StringValue(filepath.Base(src.File) + "(" + filepath.Base(src.Function) + "):" + strconv.Itoa(src.Line))
				default:
					a.Value = slog.StringValue(filepath.Base(src.File) + ":" + strconv.Itoa(src.Line))
				}
			} else {
				// If we can't cast the any from the value to an *slog.Source, it means the caller logged
				// another attribute with a key of `caller`.
				// Prevent duplicate keys (necessary for proper
				// JSON) by renaming the key to
				// `logged_caller`.
				a.Key = reservedKeyPrefix + key
			}
		case slog.LevelKey:
			if lvl, ok := a.Value.Any().(slog.Level); ok {
				a.Value = slog.StringValue(strings.ToLower(lvl.String()))
			} else {
				// If we can't cast the any from the value to
				// an slog.Level, it means the caller logged
				// another attribute with a key of `level`.
				// Prevent duplicate keys (necessary for proper
				// JSON) by renaming the key to `logged_level`.
				a.Key = reservedKeyPrefix + key
			}
		default:
		}

		if v, ok := a.Value.Any().(time.Duration); ok {
			a.Value = slog.StringValue(v.String())
		}

		return a
	}
}
