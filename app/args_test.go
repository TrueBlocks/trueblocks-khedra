package app

import (
	"os"
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	"github.com/google/go-cmp/cmp"
)

// Testing status: reviewed

func TestArgsParseArgsInternal(t *testing.T) {
	defer types.SetupTest([]string{})()
	tests := []struct {
		name        string
		args        []string
		expHelp     bool
		expVersion  bool
		expCmds     []string
		expCmdCount int
	}{
		{
			name:        "No args",
			args:        []string{"khedra"},
			expHelp:     true,
			expVersion:  false,
			expCmds:     []string{},
			expCmdCount: 0,
		},
		{
			name:        "Only help",
			args:        []string{"khedra", "--help"},
			expHelp:     true,
			expVersion:  false,
			expCmds:     []string{},
			expCmdCount: 0,
		},
		{
			name:        "Only version",
			args:        []string{"khedra", "--version"},
			expHelp:     false,
			expVersion:  true,
			expCmds:     []string{},
			expCmdCount: 0,
		},
		{
			name:        "Help and version",
			args:        []string{"khedra", "--help", "--version"},
			expHelp:     true,
			expVersion:  true,
			expCmds:     []string{},
			expCmdCount: 0,
		},
		{
			name:        "Commands only",
			args:        []string{"khedra", "config", "edit"},
			expHelp:     false,
			expVersion:  false,
			expCmds:     []string{"config", "edit"},
			expCmdCount: 2,
		},
		{
			name:        "Commands with help",
			args:        []string{"khedra", "config", "--help"},
			expHelp:     true,
			expVersion:  false,
			expCmds:     []string{"config"},
			expCmdCount: 1,
		},
		{
			name:        "Commands with version",
			args:        []string{"khedra", "config", "--version"},
			expHelp:     false,
			expVersion:  true,
			expCmds:     []string{"config"},
			expCmdCount: 1,
		},
		{
			name:        "Commands with help and version",
			args:        []string{"khedra", "--help", "config", "--version"},
			expHelp:     true,
			expVersion:  true,
			expCmds:     []string{"config"},
			expCmdCount: 1,
		},
		{
			name:        "Duplicate commands",
			args:        []string{"khedra", "config", "config"},
			expHelp:     false,
			expVersion:  false,
			expCmds:     []string{"config"},
			expCmdCount: 1,
		},
		{
			name:        "Non-standard flags",
			args:        []string{"khedra", "-unknown", "--flag"},
			expHelp:     false,
			expVersion:  false,
			expCmds:     []string{"-unknown", "--flag"},
			expCmdCount: 0,
		},
		{
			name:        "Complex order of commands and flags",
			args:        []string{"khedra", "--help", "config", "--version"},
			expHelp:     true,
			expVersion:  true,
			expCmds:     []string{"config"},
			expCmdCount: 1,
		},
		{
			name:        "No commands, multiple flags",
			args:        []string{"khedra", "--help", "--version", "--flag"},
			expHelp:     true,
			expVersion:  true,
			expCmds:     []string{"--flag"},
			expCmdCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			hasHelp, hasVersion, commands, commandCount := parseArgsInternal()
			if hasHelp != tt.expHelp {
				t.Errorf("expected hasHelp=%v, got %v", tt.expHelp, hasHelp)
			}
			if hasVersion != tt.expVersion {
				t.Errorf("expected hasVersion=%v, got %v", tt.expVersion, hasVersion)
			}
			if diff := cmp.Diff(tt.expCmds, commands); diff != "" {
				t.Errorf("commands mismatch (-want +got):\n%s", diff)
			}
			if commandCount != tt.expCmdCount {
				t.Errorf("expected commandCount=%d, got %d", tt.expCmdCount, commandCount)
			}
		})
	}
}

func TestArgsCleanArgs(t *testing.T) {
	defer types.SetupTest([]string{})()
	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "No args",
			args:     []string{"./program"},
			expected: []string{"./program", "help"},
		},
		{
			name:     "Help flag",
			args:     []string{"./program", "--help"},
			expected: []string{"./program", "help"},
		},
		{
			name:     "Version flag",
			args:     []string{"./program", "--version"},
			expected: []string{"./program", "version"},
		},
		{
			name:     "Help and command",
			args:     []string{"./program", "--help"},
			expected: []string{"./program", "help"},
		},
		{
			name:     "Commands only",
			args:     []string{"./program", "config"},
			expected: []string{"./program", "config"},
		},
		{
			name:     "Help and version",
			args:     []string{"./program", "--help", "--version"},
			expected: []string{"./program", "help"},
		},
		{
			name:     "Unrecognized flag",
			args:     []string{"./program", "-unknown"},
			expected: []string{"./program", "-unknown"},
		},
		{
			name:     "Empty arguments",
			args:     []string{"./program", ""},
			expected: []string{"./program"},
		},
		{
			name:     "Duplicate commands",
			args:     []string{"./program", "config"},
			expected: []string{"./program", "config"},
		},
		{
			name:     "Single command with flag",
			args:     []string{"./khedra", "--all"},
			expected: []string{"./khedra", "--all"},
		},
		{
			name:     "Command with subcommand",
			args:     []string{"./khedra", "config", "show"},
			expected: []string{"./khedra", "config", "show"},
		},
		{
			name:     "Command with subcommand and argument",
			args:     []string{"./khedra", "config", "show", "--key", "value"},
			expected: []string{"./khedra", "config", "show", "--key", "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			result := cleanArgs()
			if diff := cmp.Diff(tt.expected, result); diff != "" {
				t.Errorf("Test %q failed, cleanArgs() mismatch (-want +got):\n%s", tt.name, diff)
			}
		})
	}
}

func TestArgsValidateArgs(t *testing.T) {
	defer types.SetupTest([]string{})()
	tests := []struct {
		name             string
		args             []string
		expectedCmdCount int
		expectedArgCount int
		expectedError    string
	}{
		// {
		// 	name:             "Exact commands and no arguments",
		// 	args:             []string{"./program"},
		// 	expectedCmdCount: 1,
		// 	expectedArgCount: 1,
		// 	expectedError:    "",
		// },
		// {
		// 	name:             "Exact commands and arguments",
		// 	args:             []string{"./program", "config", "--flag"},
		// 	expectedCmdCount: 1,
		// 	expectedArgCount: 2,
		// 	expectedError:    "",
		// },
		// {
		// 	name:             "Too many commands",
		// 	args:             []string{"./program", "config"},
		// 	expectedCmdCount: 1,
		// 	expectedArgCount: 1,
		// 	expectedError:    "too many commands provided, expected 1",
		// },
		{
			name:             "Unknown command",
			args:             []string{"./program", "invalid"},
			expectedCmdCount: 1,
			expectedArgCount: 1,
			expectedError:    "command 'invalid' not found",
		},
		// {
		// 	name:             "Too many arguments",
		// 	args:             []string{"./program", "--flag1", "--flag2"},
		// 	expectedCmdCount: 1,
		// 	expectedArgCount: 2,
		// 	expectedError:    "too many arguments provided, expected 2",
		// },
		// {
		// 	name:             "Argument mismatch",
		// 	args:             []string{"./program", "--flag"},
		// 	expectedCmdCount: 1,
		// 	expectedArgCount: 2,
		// 	expectedError:    "unexpected number of commands (0/1) or arguments (1/2)",
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			err := validateArgs(tt.expectedCmdCount, tt.expectedArgCount)
			if (err != nil && err.Error() != tt.expectedError) || (err == nil && tt.expectedError != "") {
				t.Errorf("Test %q failed: expected error %q, got %q", tt.name, tt.expectedError, err)
			}
		})
	}
}

func TestArgsUnknownCmd(t *testing.T) {
	defer types.SetupTest([]string{})()
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "No arguments",
			args:     []string{"./program"},
			expected: "",
		},
		{
			name:     "All known commands",
			args:     []string{"./program", "config", "show"},
			expected: "",
		},
		{
			name:     "Unknown command",
			args:     []string{"./program", "invalid", "config"},
			expected: "invalid",
		},
		{
			name:     "Flags only, all invalid",
			args:     []string{"./program", "--unknown", "-key"},
			expected: "--unknown",
		},
		{
			name:     "First unknown is a flag",
			args:     []string{"./program", "--unknown", "config"},
			expected: "--unknown",
		},
		{
			name:     "First unknown is a command",
			args:     []string{"./program", "unknown", "--invalid"},
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			result := getUnknownCmd()
			if result != tt.expected {
				t.Errorf("Test %q failed: expected %q, got %q", tt.name, tt.expected, result)
			}
		})
	}
}
