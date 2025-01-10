package app

import (
	"reflect"
	"testing"
)

func TestParseArgsInternal(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		expHelp    bool
		expVersion bool
		expCmds    []string
	}{
		{
			name:       "No args",
			args:       []string{},
			expHelp:    true,
			expVersion: false,
			expCmds:    []string{},
		},
		{
			name:       "Only help",
			args:       []string{"--help"},
			expHelp:    true,
			expVersion: false,
			expCmds:    []string{},
		},
		{
			name:       "Only version",
			args:       []string{"--version"},
			expHelp:    false,
			expVersion: true,
			expCmds:    []string{},
		},
		{
			name:       "Help and version",
			args:       []string{"--help", "--version"},
			expHelp:    true,
			expVersion: true,
			expCmds:    []string{},
		},
		{
			name:       "Commands only",
			args:       []string{"init", "config", "edit"},
			expHelp:    false,
			expVersion: false,
			expCmds:    []string{"init", "config", "edit"},
		},
		{
			name:       "Commands with help",
			args:       []string{"init", "config", "--help"},
			expHelp:    true,
			expVersion: false,
			expCmds:    []string{"init", "config"},
		},
		{
			name:       "Commands with version",
			args:       []string{"init", "config", "--version"},
			expHelp:    false,
			expVersion: true,
			expCmds:    []string{"init", "config"},
		},
		{
			name:       "Commands with help and version",
			args:       []string{"init", "--help", "config", "--version"},
			expHelp:    true,
			expVersion: true,
			expCmds:    []string{"init", "config"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasHelp, hasVersion, commands, _ := parseArgsInternal(tt.args)
			if hasHelp != tt.expHelp {
				t.Errorf("expected hasHelp=%v, got %v", tt.expHelp, hasHelp)
			}
			if hasVersion != tt.expVersion {
				t.Errorf("expected hasVersion=%v, got %v", tt.expVersion, hasVersion)
			}
			if !reflect.DeepEqual(commands, tt.expCmds) {
				t.Errorf("expected commands=%v, got %v", tt.expCmds, commands)
			}
		})
	}
}

func TestCleanArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expectErr bool
		expected  []string
	}{
		{
			name:      "No args",
			args:      []string{"./program"},
			expectErr: false,
			expected:  []string{"./program", "help"},
		},
		{
			name:      "Help flag",
			args:      []string{"./program", "--help"},
			expectErr: false,
			expected:  []string{"./program", "help"},
		},
		{
			name:      "Version flag",
			args:      []string{"./program", "--version"},
			expectErr: false,
			expected:  []string{"./program", "version"},
		},
		{
			name:      "Help and command",
			args:      []string{"./program", "--help", "init"},
			expectErr: false,
			expected:  []string{"./program", "help", "init"},
		},
		{
			name:      "Commands only",
			args:      []string{"./program", "init", "config"},
			expectErr: true,
			expected:  []string{"./program", "init", "config"},
		},
		{
			name:      "Help and version",
			args:      []string{"./program", "--help", "--version"},
			expectErr: false,
			expected:  []string{"./program", "help"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanArgs(tt.args)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
