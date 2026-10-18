// Package logging is a package that provides a logger for the application.
package logging

import (
	"log/slog"
	"os"
)

var defaultWriter = os.Stderr

func New(config *LoggingConfig) *slog.Logger {
	if config.Level == nil {
		config.Level = &slog.LevelVar{}
	}

	if config.Writer == nil {
		config.Writer = defaultWriter
	}

	logHandlerOpts := &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			key := a.key
			switch key {
			case slog.SourceKey, "caller":
						if src, ok := a.Value.Any().(*slog.Source); ok {
							a.Key = "caller"
							switch lvl.String() {
							case "debug":
								a.Value = slog.StringValue(filepath.Base(src.File) + "(" + filepath.Base(src.Function) + "):" + strconv.Itoa(src.Line))
							default:
								a.Value = slog.StringValue(filepath.Base(src.File) + ":" + strconv.Itoa(src.Line))
							}
						} else {
							// If we can't cast the any from the value to
							// an *slog.Source, it means the caller logged
							// another attribute with a key of `caller`.
							// Prevent duplicate keys (necessary for proper
							// JSON) by renaming the key to
							// `logged_caller`.
							a.Key = reservedKeyPrefix + key
						}

			}
			return a

		}
	}
}

	if config.Format != "" && config.Format == "json" {
		return slog.New(slog.NewJSONHandler(config.Writer, logHandlerOpts))
	}

	return slog.New(slog.NewTextHandler(config.Writer, logHandlerOpts))
}
