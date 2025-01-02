package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logging struct {
	Folder     string `koanf:"folder" validate:"required,dirpath"`
	Filename   string `koanf:"filename" validate:"required,endswith=.log"`
	MaxSizeMb  int    `koanf:"max_size_mb" validate:"required,min=5"`
	MaxBackups int    `koanf:"max_backups" validate:"required,min=1"`
	MaxAgeDays int    `koanf:"max_age_days" validate:"required,min=1"`
	Compress   bool   `koanf:"compress"`
	LogLevel   string `koanf:"log_level" validate:"oneof=debug info warn error"`
}

func NewLogging() Logging {
	return Logging{
		Folder:     "~/.khedra/logs",
		Filename:   "khedra.log",
		MaxSizeMb:  10,
		MaxBackups: 3,
		MaxAgeDays: 10,
		Compress:   true,
		LogLevel:   "info",
	}
}

// NewLoggers creates and returns two loggers: one (fileLogger) for
// logging to a file and another (progressLogger) for logging to stderr.
func NewLoggers(cfg Logging) (*slog.Logger, *slog.Logger) {
	fileLog := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Folder, cfg.Filename),
		MaxSize:    cfg.MaxSizeMb,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	fileHandler := slog.NewTextHandler(fileLog, &slog.HandlerOptions{
		Level: convertLogLevel(cfg.LogLevel),
	})
	fileLogger := slog.New(fileHandler)

	progressHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: convertLogLevel(cfg.LogLevel),
	})
	progressLogger := slog.New(progressHandler)

	return fileLogger, progressLogger
}

// convertLogLevel converts a string log level to a slog.Level.
func convertLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
