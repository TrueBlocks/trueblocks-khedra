package types

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logging struct {
	Folder     string `koanf:"folder" validate:"required,folder_exists"`
	Filename   string `koanf:"filename" validate:"required,endswith=.log"`
	MaxSizeMb  int    `koanf:"max_size_mb" yaml:"max_size_mb" validate:"required,min=5"`
	MaxBackups int    `koanf:"max_backups" yaml:"max_backups" validate:"required,min=1"`
	MaxAgeDays int    `koanf:"max_age_days" yaml:"max_age_days" validate:"required,min=1"`
	Compress   bool   `koanf:"compress"`
	LogLevel   string `koanf:"log_level" yaml:"log_level" validate:"oneof=debug info warn error"`
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
	fileHandler := &customHandler{
		writer: &lumberjack.Logger{
			Filename:   filepath.Join(cfg.Folder, cfg.Filename),
			MaxSize:    cfg.MaxSizeMb,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		},
		level: convertLogLevel(cfg.LogLevel),
	}
	fileLogger := slog.New(fileHandler)

	progressHandler := &customHandler{
		writer: os.Stderr,
		level:  convertLogLevel(cfg.LogLevel),
	}
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

type customHandler struct {
	writer io.Writer
	level  slog.Level
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < h.level {
		return nil
	}
	levels := map[string]string{
		"DEBU": "DEBG",
		"INFO": "INFO",
		"WARN": "WARN",
		"ERRO": "EROR",
	}
	lev := r.Level.String()[:4]
	timeFormat := r.Time.Format("02-01|15:04:05.000")
	formattedMessage := r.Message
	if r.NumAttrs() > 0 {
		formattedMessage = fmt.Sprintf("%-25.25s", r.Message)
	}
	logMsg := fmt.Sprintf("%4.4s[%s] %s ", levels[lev], timeFormat, formattedMessage)
	r.Attrs(func(attr slog.Attr) bool {
		logMsg += fmt.Sprintf(" %s=%v", colors.Green+attr.Key+colors.Off, attr.Value)
		return true
	})
	fmt.Fprintln(h.writer, logMsg)
	return nil
}

func (h *customHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	return h
}

/*
func NewCustomLogger() (*slog.Logger, slog.Level) {
	logger.SetLoggerWriter(io.Discard)
	logLevel := slog.LevelInfo
	if ll, ok := os.LookupEnv("TB_LOGLEVEL"); ok {
		switch strings.ToLower(ll) {
		case "debug":
			logLevel = slog.LevelDebug
		case "info":
			logLevel = slog.LevelInfo
		case "warn":
			logLevel = slog.LevelWarn
		case "error":
			logLevel = slog.LevelError
		}
	}
	customHandler := &customHandler{
		writer: os.Stderr,
		level:  logLevel,
	}
	return slog.New(customHandler), logLevel
}
*/
