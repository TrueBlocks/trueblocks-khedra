package types

import (
	"bytes"
	"log/slog"
	"regexp"
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
	"github.com/alecthomas/assert/v2"
)

// Testing status: reviewed

func TestLoggingNew(t *testing.T) {
	logging := NewLogging()
	assert.Equal(t, "~/.khedra/logs", logging.Folder)
	assert.Equal(t, "khedra.log", logging.Filename)
	assert.Equal(t, 10, logging.MaxSize)
	assert.Equal(t, 3, logging.MaxBackups)
	assert.Equal(t, 10, logging.MaxAge)
	assert.Equal(t, "info", logging.Level)
	assert.True(t, logging.Compress)
}

func TestLoggingValidation(t *testing.T) {
	tempDir := createTempDir(t, true) // Helper function to create a temp writable directory

	tests := []struct {
		name    string
		logging Logging
		wantErr bool
	}{
		{
			name: "Valid Logging struct",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: false,
		},
		{
			name: "Valid Logging level warn",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "warn",
			},
			wantErr: false,
		},
		{
			name: "Invalid Logging level",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "bogus",
			},
			wantErr: true,
		},
		{
			name: "Missing Folder",
			logging: Logging{
				Filename:   "app.log",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
		{
			name: "Non-existent Folder",
			logging: Logging{
				Folder:     "/non/existent/path",
				Filename:   "app.log",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
		{
			name: "Missing Filename",
			logging: Logging{
				Folder:     tempDir,
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
		{
			name: "Filename without .log extension",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.txt",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
		{
			name: "MaxSize is zero",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSize:    0,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
		{
			name: "MaxBackups is negative",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSize:    10,
				MaxBackups: -1,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
		{
			name: "MaxAge is negative",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     -1,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Validate4(&tt.logging)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for test case '%s'", tt.name)
			} else {
				assert.NoError(t, err, "Did not expect error for test case '%s'", tt.name)
			}
		})
	}
}

func TestLoggingReadAndWrite(t *testing.T) {
	tempFilePath := "temp_config.yaml"
	content := `
  folder: ~/.khedra/logs
  filename: khedra.log
  level: debug
  maxSize: 10
  maxBackups: 3
  maxAge: 10
  compress: true
`

	assertions := func(t *testing.T, logging *Logging) {
		assert.Equal(t, "~/.khedra/logs", logging.Folder, "Folder should match the expected value")
		assert.Equal(t, "khedra.log", logging.Filename, "Filename should match the expected value")
		assert.Equal(t, "debug", logging.Level, "Level should match the expected value")
		assert.Equal(t, 10, logging.MaxSize, "MaxSize should match the expected value")
		assert.Equal(t, 3, logging.MaxBackups, "MaxBackups should match the expected value")
		assert.Equal(t, 10, logging.MaxAge, "MaxAge should match the expected value")
		assert.True(t, logging.Compress, "Compress should be true")
	}

	ReadAndWriteWithAssertions[Logging](t, tempFilePath, content, assertions)
}

func TestCustomHandlerLogFormatting(t *testing.T) {
	var output bytes.Buffer
	handler := newCustomHandler(&output, "info")
	logger := slog.New(handler)
	logger.Info("Test message", slog.Int("key", 42))

	logOutput := stripAnsiCodes(t, output.String())

	expectedSubstring := "INFO"
	assert.Contains(t, logOutput, expectedSubstring, "Expected log output to contain log level")
	assert.Contains(t, logOutput, "Test message", "Expected log output to contain the message")
	assert.Contains(t, logOutput, "key=42", "Expected log output to contain attributes")
}

func TestCustomHandlerLevels(t *testing.T) {
	var output bytes.Buffer
	handler := newCustomHandler(&output, "warn")
	logger := slog.New(handler)
	logger.Info("This should not be logged")
	logger.Warn("This should be logged")

	logOutput := stripAnsiCodes(t, output.String())

	assert.NotContains(t, logOutput, "This should not be logged", "Logs below the configured level should not be written")
	assert.Contains(t, logOutput, "This should be logged", "Logs at or above the configured level should be written")
}

func TestNewLoggersIntegration(t *testing.T) {
	var fileOutput, progOutput bytes.Buffer

	fileHandler := newCustomHandler(&fileOutput, "info")
	progHandler := newCustomHandler(&progOutput, "info")

	fileLogger := slog.New(fileHandler)
	progLogger := slog.New(progHandler)

	fileLogger.Info("File logger message")
	progLogger.Info("Prog logger message")

	fileLogOutput := stripAnsiCodes(t, fileOutput.String())
	progLogOutput := stripAnsiCodes(t, progOutput.String())

	assert.Contains(t, fileLogOutput, "File logger message", "Expected log to appear in file logger output")
	assert.Contains(t, progLogOutput, "Prog logger message", "Expected log to appear in progress logger output")
}

func TestConvertLevelUnsupported(t *testing.T) {
	level := convertLevel("unsupported")
	if level != slog.LevelInfo {
		t.Errorf("expected fallback to DefaultLevel, got %v", level)
	}
}

func TestCustomHandlerWithAttrs(t *testing.T) {
	var output bytes.Buffer
	handler := newCustomHandler(&output, "info")
	updatedHandler := handler.WithAttrs([]slog.Attr{slog.Int("globalKey", 2)})
	if updatedHandler == nil {
		t.Errorf("expected non-nil handler after WithAttrs")
	}
	logger := slog.New(updatedHandler)
	logger.Info("Test message", slog.Int("localKey", 1))

	logOutput := stripAnsiCodes(t, output.String())

	assert.Contains(t, logOutput, "INFO", "Expected log output to contain log level")
	assert.Contains(t, logOutput, "Test message", "Expected log output to contain the message")
	assert.Contains(t, logOutput, "globalKey=2", "Expected log output to contain attributes")
	assert.Contains(t, logOutput, "localKey=1", "Expected log output to contain attributes")
}

func TestCustomHandlerWithGroup(t *testing.T) {
	var output bytes.Buffer
	handler := newCustomHandler(&output, "debug")
	updatedHandler := handler.WithGroup("group")
	if updatedHandler == nil {
		t.Errorf("expected non-nil handler after WithGroup")
	}
	logger := slog.New(updatedHandler)
	logger.Debug("Test message", slog.Int("localKey", 1))

	logOutput := stripAnsiCodes(t, output.String())

	assert.Contains(t, logOutput, "DEBG", "Expected log output to contain log level")
	assert.Contains(t, logOutput, "Test message", "Expected log output to contain the message")
	assert.Contains(t, logOutput, "groups=[group]", "Expected log output to contain the group")
	assert.Contains(t, logOutput, "localKey=1", "Expected log output to contain the attribute")
}

func TestCustomHandlerWithBoth(t *testing.T) {
	var output bytes.Buffer
	handler := newCustomHandler(&output, "debug")
	updatedHandler := handler.WithGroup("group").WithAttrs([]slog.Attr{slog.Int("globalKey", 2)})
	if updatedHandler == nil {
		t.Errorf("expected non-nil handler after WithGroup")
	}
	logger := slog.New(updatedHandler)
	logger.Debug("Test message", slog.Int("localKey", 1))

	logOutput := stripAnsiCodes(t, output.String())

	assert.Contains(t, logOutput, "DEBG", "Expected log output to contain log level")
	assert.Contains(t, logOutput, "Test message", "Expected log output to contain the message")
	assert.Contains(t, logOutput, "groups=[group]", "Expected log output to contain the group")
	assert.Contains(t, logOutput, "globalKey=2", "Expected log output to contain attributes")
	assert.Contains(t, logOutput, "localKey=1", "Expected log output to contain the attribute")
}

func stripAnsiCodes(t *testing.T, input string) string {
	t.Helper()
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}
