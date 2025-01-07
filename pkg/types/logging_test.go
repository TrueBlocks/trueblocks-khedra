package types

import (
	"bytes"
	"log/slog"
	"regexp"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestLoggingNew(t *testing.T) {
	logging := NewLogging()
	assert.Equal(t, "~/.khedra/logs", logging.Folder)
	assert.Equal(t, "khedra.log", logging.Filename)
	assert.Equal(t, 10, logging.MaxSizeMb)
	assert.Equal(t, 3, logging.MaxBackups)
	assert.Equal(t, 10, logging.MaxAgeDays)
	assert.Equal(t, "info", logging.LogLevel)
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
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: false,
		},
		{
			name: "Valid Logging level warn",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "warn",
			},
			wantErr: false,
		},
		{
			name: "Invalid Logging level",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "bogus",
			},
			wantErr: true,
		},
		{
			name: "Missing Folder",
			logging: Logging{
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "Non-existent Folder",
			logging: Logging{
				Folder:     "/non/existent/path",
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "Missing Filename",
			logging: Logging{
				Folder:     tempDir,
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "Filename without .log extension",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.txt",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "MaxSizeMb is zero",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  0,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "MaxBackups is negative",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: -1,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "MaxAgeDays is negative",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: -1,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.logging)
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
  log_level: debug
  max_size_mb: 10
  max_backups: 3
  max_age_days: 10
  compress: true
`

	assertions := func(t *testing.T, logging *Logging) {
		assert.Equal(t, "~/.khedra/logs", logging.Folder, "Folder should match the expected value")
		assert.Equal(t, "khedra.log", logging.Filename, "Filename should match the expected value")
		assert.Equal(t, "debug", logging.LogLevel, "LogLevel should match the expected value")
		assert.Equal(t, 10, logging.MaxSizeMb, "MaxSizeMb should match the expected value")
		assert.Equal(t, 3, logging.MaxBackups, "MaxBackups should match the expected value")
		assert.Equal(t, 10, logging.MaxAgeDays, "MaxAgeDays should match the expected value")
		assert.True(t, logging.Compress, "Compress should be true")
	}

	ReadAndWriteWithAssertions[Logging](t, tempFilePath, content, assertions)
}

func TestCustomHandlerLogFormatting(t *testing.T) {
	var output bytes.Buffer
	handler := &customHandler{
		writer: &output,
		level:  slog.LevelInfo,
	}

	logger := slog.New(handler)
	logger.Info("Test message", slog.Int("key", 42))

	logOutput := stripAnsiCodes(output.String())

	expectedSubstring := "INFO"
	assert.Contains(t, logOutput, expectedSubstring, "Expected log output to contain log level")
	assert.Contains(t, logOutput, "Test message", "Expected log output to contain the message")
	assert.Contains(t, logOutput, "key=42", "Expected log output to contain attributes")
}

func TestCustomHandlerLogLevels(t *testing.T) {
	var output bytes.Buffer
	handler := &customHandler{
		writer: &output,
		level:  slog.LevelWarn,
	}

	logger := slog.New(handler)
	logger.Info("This should not be logged")
	logger.Warn("This should be logged")

	logOutput := stripAnsiCodes(output.String())

	assert.NotContains(t, logOutput, "This should not be logged", "Logs below the configured level should not be written")
	assert.Contains(t, logOutput, "This should be logged", "Logs at or above the configured level should be written")
}

func TestNewLoggersIntegration(t *testing.T) {
	var fileOutput, progOutput bytes.Buffer

	fileHandler := &customHandler{writer: &fileOutput, level: slog.LevelInfo}
	progHandler := &customHandler{writer: &progOutput, level: slog.LevelInfo}

	fileLogger := slog.New(fileHandler)
	progLogger := slog.New(progHandler)

	fileLogger.Info("File logger message")
	progLogger.Info("Prog logger message")

	fileLogOutput := stripAnsiCodes(fileOutput.String())
	progLogOutput := stripAnsiCodes(progOutput.String())

	assert.Contains(t, fileLogOutput, "File logger message", "Expected log to appear in file logger output")
	assert.Contains(t, progLogOutput, "Prog logger message", "Expected log to appear in progress logger output")
}

func stripAnsiCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}
