package types

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func TestLoggingNew(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	expectedFolder := filepath.Join(homeDir, ".khedra", "logs")

	logging := NewLogging()
	assert.Equal(t, expectedFolder, logging.Folder)
	assert.Equal(t, "khedra.log", logging.Filename)
	assert.False(t, logging.ToFile)
	assert.Equal(t, 10, logging.MaxSize)
	assert.Equal(t, 3, logging.MaxBackups)
	assert.Equal(t, 10, logging.MaxAge)
	assert.Equal(t, "info", logging.Level)
	assert.True(t, logging.Compress)
}

func TestLoggingValidation(t *testing.T) {
	tempDir := createTempDir(t, true)

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
				ToFile:     true,
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
				ToFile:     false,
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
				ToFile:     false,
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
				ToFile:     false,
				MaxSize:    10,
				MaxBackups: 3,
				MaxAge:     7,
				Compress:   true,
				Level:      "info",
			},
			wantErr: true,
		},
		// {
		// 	name: "Non-existent Folder",
		// 	logging: Logging{
		// 		Folder:     "/non/existent/path",
		// 		Filename:   "app.log",
		// 		ToFile:     true,
		// 		MaxSize:    10,
		// 		MaxBackups: 3,
		// 		MaxAge:     7,
		// 		Compress:   true,
		// 		Level:      "info",
		// 	},
		// 	wantErr: true,
		// },
		{
			name: "Missing Filename",
			logging: Logging{
				Folder:     tempDir,
				MaxSize:    10,
				ToFile:     false,
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
				ToFile:     false,
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
				ToFile:     false,
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
				ToFile:     false,
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
				ToFile:     true,
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
			err := Validate(&tt.logging)
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
  toFile: false
  maxSize: 10
  maxBackups: 3
  maxAge: 10
  compress: true
  level: debug
`

	assertions := func(t *testing.T, logging *Logging) {
		assert.Equal(t, "~/.khedra/logs", logging.Folder, "Folder should match the expected value")
		assert.Equal(t, "khedra.log", logging.Filename, "Filename should match the expected value")
		assert.False(t, logging.ToFile, "ToFile should be true")
		assert.Equal(t, "debug", logging.Level, "Level should match the expected value")
		assert.Equal(t, 10, logging.MaxSize, "MaxSize should match the expected value")
		assert.Equal(t, 3, logging.MaxBackups, "MaxBackups should match the expected value")
		assert.Equal(t, 10, logging.MaxAge, "MaxAge should match the expected value")
		assert.True(t, logging.Compress, "Compress should be true")
	}

	ReadAndWriteWithAssertions(t, tempFilePath, content, assertions)
}

func TestConvertLevelUnsupported(t *testing.T) {
	level := convertLevel("unsupported")
	if level != slog.LevelInfo {
		t.Errorf("expected fallback to DefaultLevel, got %v", level)
	}
}

// Added tests per ai/TestDesign_logging.go.md (safe: all file writes confined to temp dirs)
func TestConvertLevel_AllMappings(t *testing.T) {
	assert.Equal(t, slog.LevelDebug, convertLevel("debug"))
	assert.Equal(t, slog.LevelInfo, convertLevel("info"))
	assert.Equal(t, slog.LevelWarn, convertLevel("warn"))
	assert.Equal(t, slog.LevelError, convertLevel("error"))
}

func TestLevelToString_CustomAndStandard(t *testing.T) {
	assert.Equal(t, "PROG", levelToString(LevelProgress))
	assert.Equal(t, "DEBUG", levelToString(slog.LevelDebug))
	assert.Equal(t, "INFO", levelToString(slog.LevelInfo))
	assert.Equal(t, "WARN", levelToString(slog.LevelWarn))
	assert.Equal(t, "ERROR", levelToString(slog.LevelError))
}

func TestNewLogger_ScreenOnly_NoFile(t *testing.T) {
	tempDir := t.TempDir()
	// Screen-only: empty Filename prevents file handler creation
	logging := Logging{Folder: tempDir, Filename: "", Level: "info", MaxSize: 5, MaxBackups: 1, MaxAge: 1, Compress: false}

	// Capture stderr
	origStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	logger := NewLogger(logging)
	logger.Info("screen only message")
	w.Close()
	os.Stderr = origStderr
	data, _ := io.ReadAll(r)
	output := string(data)
	assert.Contains(t, output, "screen only message")
	// Ensure no file created
	entries, _ := os.ReadDir(tempDir)
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".log") {
			t.Fatalf("unexpected log file created: %s", e.Name())
		}
	}
}

func TestNewLogger_FileLogging_Writes(t *testing.T) {
	tempDir := t.TempDir()
	logging := Logging{Folder: tempDir, Filename: "test.log", Level: "info", MaxSize: 5, MaxBackups: 1, MaxAge: 1, Compress: false}
	logger := NewLogger(logging)
	msg := "file target message"
	logger.Info(msg)
	logPath := filepath.Join(tempDir, "test.log")
	// give a moment for synchronous write (not really needed, but safe)
	time.Sleep(10 * time.Millisecond)
	content, err := os.ReadFile(logPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), msg)
}

func TestCustomLogger_ProgressShownAndSuppressed(t *testing.T) {
	// Shown (level=info)
	tempDir := t.TempDir()
	loggingInfo := Logging{Folder: tempDir, Filename: "", Level: "info", MaxSize: 5, MaxBackups: 1, MaxAge: 1}
	// capture stderr
	origStderr := os.Stderr
	r1, w1, _ := os.Pipe()
	os.Stderr = w1
	loggerInfo := NewLogger(loggingInfo)
	loggerInfo.Progress("progress message one")
	w1.Close()
	os.Stderr = origStderr
	data1, _ := io.ReadAll(r1)
	assert.Contains(t, string(data1), "progress message one")

	// Suppressed (level=error)
	r2, w2, _ := os.Pipe()
	os.Stderr = w2
	loggingErr := Logging{Folder: tempDir, Filename: "", Level: "error", MaxSize: 5, MaxBackups: 1, MaxAge: 1}
	loggerErr := NewLogger(loggingErr)
	loggerErr.Progress("hidden progress")
	w2.Close()
	os.Stderr = origStderr
	data2, _ := io.ReadAll(r2)
	assert.NotContains(t, string(data2), "hidden progress")
}

func TestColorTextHandler_Format(t *testing.T) {
	buf := &bytes.Buffer{}
	h := &ColorTextHandler{Writer: buf, Level: slog.LevelDebug}
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "message body", 0)
	rec.AddAttrs(slog.String("key", "value"))
	_ = h.Handle(context.Background(), rec)
	raw := buf.String()
	// strip ANSI color codes
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	stripped := re.ReplaceAllString(raw, "")
	assert.Contains(t, stripped, "INFO")
	assert.Contains(t, stripped, "message body")
	assert.Contains(t, stripped, "key=value")
}
