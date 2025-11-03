package types

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/colors"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/logger"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logging struct {
	Folder     string `koanf:"folder" json:"folder,omitempty" validate:"required,folder_exists"`
	Filename   string `koanf:"filename" json:"filename,omitempty" validate:"required,endswith=.log"`
	ToFile     bool   `koanf:"toFile" json:"toFile,omitempty"`
	MaxSize    int    `koanf:"maxSize" yaml:"maxSize" json:"maxSize,omitempty" validate:"required,min=5"`
	MaxBackups int    `koanf:"maxBackups" yaml:"maxBackups" json:"maxBackups,omitempty" validate:"required,min=1"`
	MaxAge     int    `koanf:"maxAge" yaml:"maxAge" json:"maxAge,omitempty" validate:"required,min=1"`
	Compress   bool   `koanf:"compress" json:"compress,omitempty"`
	Level      string `koanf:"level" yaml:"level" json:"level,omitempty" validate:"oneof=debug info warn error"`
}

func NewLogging() Logging {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Panic("could not determine user home directory")
	}
	return Logging{
		Folder:     filepath.Join(homeDir, ".khedra", "logs"),
		Filename:   "khedra.log",
		ToFile:     false,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     10,
		Compress:   true,
		Level:      "info",
	}
}

func convertLevel(level string) slog.Level {
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

const LevelProgress slog.Level = slog.LevelInfo + 1

func levelToString(level slog.Level) string {
	switch level {
	case LevelProgress:
		return "PROG"
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		return level.String()
	}
}

type multiHandler struct {
	writeBoth     bool
	screenHandler slog.Handler
	fileHandler   slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return m.screenHandler.Enabled(ctx, level) || (m.writeBoth && m.fileHandler.Enabled(ctx, level))
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level == LevelProgress {
		if m.screenHandler.Enabled(ctx, r.Level) {
			return m.screenHandler.Handle(ctx, r)
		}
		return nil
	}

	if m.screenHandler.Enabled(ctx, r.Level) {
		if err := m.screenHandler.Handle(ctx, r); err != nil {
			return err
		}
	}
	if m.writeBoth && m.fileHandler.Enabled(ctx, r.Level) {
		if err := m.fileHandler.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &multiHandler{
		screenHandler: m.screenHandler.WithAttrs(attrs),
		fileHandler:   m.fileHandler.WithAttrs(attrs),
		writeBoth:     m.writeBoth,
	}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	return &multiHandler{
		screenHandler: m.screenHandler.WithGroup(name),
		fileHandler:   m.fileHandler.WithGroup(name),
		writeBoth:     m.writeBoth,
	}
}

type CustomLogger struct {
	*slog.Logger
	screenHandler slog.Handler
}

func (c *CustomLogger) Panic(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	panic(s)
}

func (c *CustomLogger) Fatal(msg string, args ...any) {
	c.Error(msg, args...)
	os.Exit(1)
}

func (c *CustomLogger) Progress(msg string, args ...any) {
	if c.screenHandler.Enabled(context.Background(), LevelProgress) {
		c.Logger.Log(context.Background(), LevelProgress, msg, args...)
	}
}

type ColorTextHandler struct {
	Writer io.Writer
	Level  slog.Level
}

func (h *ColorTextHandler) Handle(ctx context.Context, r slog.Record) error {
	_ = ctx
	levelColors := map[slog.Level]string{
		slog.LevelDebug: colors.Cyan,
		slog.LevelInfo:  colors.Green,
		slog.LevelWarn:  colors.BrightYellow,
		slog.LevelError: colors.Red,
		LevelProgress:   colors.BrightBlue,
	}

	timestamp := r.Time.Format(time.RFC3339)
	levelColor, exists := levelColors[r.Level]
	if !exists {
		levelColor = colors.Off
	}
	levelStr := fmt.Sprintf("%s%s%s", levelColor, levelToString(r.Level), colors.Off)

	fixedMsg := ""
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		if a.Key == "time" {
			return true
		}
		attrs += fmt.Sprintf("%s%s%s=%v ", colors.Green, a.Key, colors.Off, a.Value)
		return true
	})

	if attrs != "" {
		fixedMsg = fmt.Sprintf("%-25.25s", r.Message)
	} else {
		fixedMsg = r.Message
	}

	var msg string
	if attrs != "" {
		msg = fmt.Sprintf("%s %s\t%s\t%s\n", timestamp, levelStr, fixedMsg, attrs)
	} else {
		msg = fmt.Sprintf("%s %s\t%s\n", timestamp, levelStr, fixedMsg)
	}

	_, err := h.Writer.Write([]byte(msg))
	return err
}

func (h *ColorTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	_ = ctx
	return level >= h.Level
}

func (h *ColorTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	_ = attrs
	return h
}

func (h *ColorTextHandler) WithGroup(name string) slog.Handler {
	_ = name
	return h
}

func NewLogger(logging Logging) *CustomLogger {
	screenHandler := &ColorTextHandler{
		Writer: os.Stderr,
		Level:  convertLevel(logging.Level),
	}

	var fileHandler slog.Handler
	if logging.Filename != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logging.Folder, logging.Filename),
			MaxSize:    logging.MaxSize,
			MaxBackups: logging.MaxBackups,
			MaxAge:     logging.MaxAge,
			Compress:   logging.Compress,
		}
		fileHandler = slog.NewTextHandler(fileWriter, &slog.HandlerOptions{
			Level: convertLevel(logging.Level),
		})
	}

	handler := &multiHandler{
		screenHandler: screenHandler,
		fileHandler:   fileHandler,
		writeBoth:     logging.Filename != "",
	}

	return &CustomLogger{
		Logger:        slog.New(handler),
		screenHandler: screenHandler,
	}
}

func (c *CustomLogger) GetLogger() *slog.Logger {
	return c.Logger
}
