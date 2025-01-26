package types

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Logging struct {
	Folder     string `koanf:"folder" json:"folder,omitempty" validate:"required,folder_exists"`
	Filename   string `koanf:"filename" json:"filename,omitempty" validate:"required,endswith=.log"`
	MaxSize    int    `koanf:"maxSize" yaml:"maxSize" json:"maxSize,omitempty" validate:"required,min=5"`
	MaxBackups int    `koanf:"maxBackups" yaml:"maxBackups" json:"maxBackups,omitempty" validate:"required,min=1"`
	MaxAge     int    `koanf:"maxAge" yaml:"maxAge" json:"maxAge,omitempty" validate:"required,min=1"`
	Compress   bool   `koanf:"compress" json:"compress,omitempty"`
	Level      string `koanf:"level" yaml:"level" json:"level,omitempty" validate:"oneof=debug info warn error"`
}

func NewLogging() Logging {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine user home directory")
	}
	return Logging{
		Folder:     filepath.Join(homeDir, ".khedra", "logs"),
		Filename:   "khedra.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     10,
		Compress:   true,
		Level:      "info",
	}
}

const LevelProgress slog.Level = slog.LevelInfo + 1

type multiHandler struct {
	writeBoth     bool
	screenHandler slog.Handler
	fileHandler   slog.Handler
}

// Enabled determines whether the log level should be processed by this handler.
// Returns true if the screen handler is enabled for the given level or, if writeBoth
// is true, if the file handler is enabled for the level. This ensures logs are
// processed if at least one of the handlers supports the level.
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

// NewLogger creates a logger with optional file logging
func NewLogger(logging Logging) *slog.Logger {
	screenHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: convertLevel(logging.Level),
	})

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

	return slog.New(handler)
}

// convertLevel converts a string log level to a slog.Level.
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

/*
type customHandler struct {
	writer      io.Writer
	baseHandler slog.Handler
	attrs       []slog.Attr // Stores attributes added via WithAttrs
	groups      []string    // Stores group names added via WithGroup
}

func newCustomHandler(w io.Writer, l string) *customHandler {
	return &customHandler{
		writer: w,
		baseHandler: slog.NewTextHandler(
			w,
			&slog.HandlerOptions{Level: convertLevel(l)},
		),
	}
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.baseHandler.Enabled(ctx, r.Level) {
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

	if len(h.groups) > 0 {
		logMsg += fmt.Sprintf("groups=[%s] ", strings.Join(h.groups, "."))
	}

	for _, attr := range h.attrs {
		logMsg += fmt.Sprintf("%s=%v ", attr.Key, attr.Value)
	}

	r.Attrs(func(attr slog.Attr) bool {
		logMsg += fmt.Sprintf("%s=%v ", attr.Key, attr.Value)
		return true
	})

	fmt.Fprintln(h.writer, logMsg)
	return nil
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.baseHandler.Enabled(ctx, level)
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a copy of the handler
	newHandler := *h
	// Merge new attributes with existing ones
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	// Create a copy of the handler
	newHandler := *h
	// Add the new group to the list of groups
	newHandler.groups = append(newHandler.groups, name)
	return &newHandler
}

func NewCustomLogger() (*slog.Logger, slog.Level) {
	logger.SetLoggerWriter(io.Discard)
	level := slog.LevelInfo
	if ll, ok := os.LookupEnv("TB_KHEDRA_LOGGING_LEVEL"); ok {
		switch strings.ToLower(ll) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		}
	}
	customHandler := &customHandler{
		writer: os.Stderr,
		level:  level,
	}
	return slog.New(customHandler), level
}
*/
