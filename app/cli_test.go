package app

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	_ "github.com/TrueBlocks/trueblocks-khedra/v5/pkg/env"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

// Testing status: not_reviewed

func TestInitializeCliCommands(t *testing.T) {
	defer types.SetupTest([]string{})()
	k := &KhedraApp{}
	cmdLine := initCli(k)
	assert.NotNil(t, cmdLine)
	assert.Equal(t, "khedra", cmdLine.Name)
	assert.Equal(t, "A tool to index, monitor, serve, and share blockchain data", cmdLine.Usage)

	commandNames := make(map[string]bool)
	for _, command := range cmdLine.Commands {
		commandNames[command.Name] = true
	}

	expectedCommands := []string{"init", "config", "version"}
	for _, cmd := range expectedCommands {
		assert.Contains(t, commandNames, cmd)
	}
}

func TestConfigShowCommand(t *testing.T) {
	defer types.SetupTest([]string{})()
	os.Args = []string{"khedra", "config", "show"}

	k := &KhedraApp{}
	cmdLine := initCli(k)

	command := getCommandByName(t, cmdLine, "config")
	assert.NotNil(t, command)

	showCommand := getSubCommandByName(t, command, "show")
	assert.NotNil(t, showCommand)

	output := captureOutput(t, func() {
		err := showCommand.Action(cli.NewContext(cmdLine, nil, nil))
		assert.NoError(t, err)
	})

	// fmt.Println("output:", string(output))
	assert.Contains(t, string(output), `general:`)
}

func TestConfigEditCommand(t *testing.T) {
	defer types.SetupTest([]string{
		"EDITOR=testing",
	})()
	os.Args = []string{"khedra", "config", "edit"}

	k := &KhedraApp{}
	cmdLine := initCli(k)

	command := getCommandByName(t, cmdLine, "config")
	assert.NotNil(t, command)

	editCommand := getSubCommandByName(t, command, "edit")
	assert.NotNil(t, editCommand)

	err := editCommand.Action(cli.NewContext(cmdLine, nil, nil))
	assert.NoError(t, err)
}

func TestCommandLineActions(t *testing.T) {
	defer types.SetupTest([]string{})()
	var testCases = []struct {
		expectError bool
		command     string
	}{
		{false, "khedra --h config daemon"},
		{false, "khedra --h config edit"},
		{false, "khedra --h config init daemon"},
		{false, "khedra --h config init"},
		{false, "khedra --h config show"},
		{false, "khedra --h config"},
		{false, "khedra --h daemon config"},
		{false, "khedra --h daemon init config"},
		{false, "khedra --h daemon init"},
		{false, "khedra --h daemon"},
		{false, "khedra --h init config"},
		{false, "khedra --h init daemon config"},
		{false, "khedra --h init daemon"},
		{false, "khedra --h init"},
		{false, "khedra --h"},

		{false, "khedra --help config daemon"},
		{false, "khedra --help config edit"},
		{false, "khedra --help config init daemon"},
		{false, "khedra --help config init"},
		{false, "khedra --help config show"},
		{false, "khedra --help config"},
		{false, "khedra --help daemon config"},
		{false, "khedra --help daemon init config"},
		{false, "khedra --help daemon init"},
		{false, "khedra --help daemon"},
		{false, "khedra --help init config"},
		{false, "khedra --help init daemon config"},
		{false, "khedra --help init daemon"},
		{false, "khedra --help init"},
		{false, "khedra --help"},

		{false, "khedra --v config edit"},
		{false, "khedra --v config init daemon"},
		{false, "khedra --v config show"},
		{false, "khedra --v daemon init config"},
		{false, "khedra --v init daemon config"},
		{false, "khedra --v"},

		{false, "khedra --version config edit"},
		{false, "khedra --version config init daemon"},
		{false, "khedra --version config init"},
		{false, "khedra --version config show"},
		{false, "khedra --version config"},
		{false, "khedra --version daemon init config"},
		{false, "khedra --version daemon init"},
		{false, "khedra --version daemon"},
		{false, "khedra --version init config"},
		{false, "khedra --version init daemon config"},
		{false, "khedra --version init daemon"},
		{false, "khedra --version init"},
		{false, "khedra --version"},

		{false, "khedra -h config daemon"},
		{false, "khedra -h config edit"},
		{false, "khedra -h config init daemon"},
		{false, "khedra -h config init"},
		{false, "khedra -h config show"},
		{false, "khedra -h config"},
		{false, "khedra -h daemon config"},
		{false, "khedra -h daemon init config"},
		{false, "khedra -h daemon init"},
		{false, "khedra -h daemon"},
		{false, "khedra -h init config"},
		{false, "khedra -h init daemon config"},
		{false, "khedra -h init daemon"},
		{false, "khedra -h init"},
		{false, "khedra -h"},

		{false, "khedra -help config daemon"},
		{false, "khedra -help config edit"},
		{false, "khedra -help config init daemon"},
		{false, "khedra -help config init"},
		{false, "khedra -help config show"},
		{false, "khedra -help config"},
		{false, "khedra -help daemon config"},
		{false, "khedra -help daemon init config"},
		{false, "khedra -help daemon init"},
		{false, "khedra -help daemon"},
		{false, "khedra -help init config"},
		{false, "khedra -help init daemon config"},
		{false, "khedra -help init daemon"},
		{false, "khedra -help init"},
		{false, "khedra -help"},

		{false, "khedra -v config daemon"},
		{false, "khedra -v config edit"},
		{false, "khedra -v config init daemon"},
		{false, "khedra -v config init"},
		{false, "khedra -v config show"},
		{false, "khedra -v config"},
		{false, "khedra -v daemon init config"},
		{false, "khedra -v daemon init"},
		{false, "khedra -v daemon"},
		{false, "khedra -v init config"},
		{false, "khedra -v init daemon config"},
		{false, "khedra -v init daemon"},
		{false, "khedra -v init"},
		{false, "khedra -v"},

		{false, "khedra -version config edit"},
		{false, "khedra -version config init daemon"},
		{false, "khedra -version config init"},
		{false, "khedra -version config show"},
		{false, "khedra -version config"},
		{false, "khedra -version daemon init config"},
		{false, "khedra -version daemon init"},
		{false, "khedra -version daemon"},
		{false, "khedra -version init config"},
		{false, "khedra -version init daemon config"},
		{false, "khedra -version init daemon"},
		{false, "khedra -version init"},
		{false, "khedra -version"},

		{false, "khedra config --h"},
		{false, "khedra config --help"},
		{false, "khedra config --v"},
		{false, "khedra config --version"},
		{false, "khedra config -h"},
		{false, "khedra config -help"},
		{false, "khedra config -v"},
		{false, "khedra config -version"},
		{false, "khedra config daemon --h"},
		{false, "khedra config daemon --help"},
		{false, "khedra config daemon --v"},
		{false, "khedra config daemon --version"},
		{false, "khedra config daemon -h"},
		{false, "khedra config daemon -help"},
		{false, "khedra config daemon -v"},
		{false, "khedra config daemon -version"},
		{false, "khedra config daemon help"},
		{false, "khedra config daemon version"},
		{false, "khedra config edit --h"},
		{false, "khedra config edit --help"},
		{false, "khedra config edit --v"},
		{false, "khedra config edit --version"},
		{false, "khedra config edit -h"},
		{false, "khedra config edit -help"},
		{false, "khedra config edit -v"},
		{false, "khedra config edit -version"},
		{false, "khedra config edit help"},
		{false, "khedra config edit version"},
		{false, "khedra config help"},
		{false, "khedra config init --h"},
		{false, "khedra config init --help"},
		{false, "khedra config init --v"},
		{false, "khedra config init --version"},
		{false, "khedra config init -h"},
		{false, "khedra config init -help"},
		{false, "khedra config init -v"},
		{false, "khedra config init -version"},
		{false, "khedra config init daemon --h"},
		{false, "khedra config init daemon --help"},
		{false, "khedra config init daemon --v"},
		{false, "khedra config init daemon --version"},
		{false, "khedra config init daemon -h"},
		{false, "khedra config init daemon -help"},
		{false, "khedra config init daemon -v"},
		{false, "khedra config init daemon -version"},
		{false, "khedra config init daemon help"},
		{false, "khedra config init daemon version"},
		{false, "khedra config init help"},
		{false, "khedra config init version"},
		{false, "khedra config show --h"},
		{false, "khedra config show --help"},
		{false, "khedra config show --v"},
		{false, "khedra config show --version"},
		{false, "khedra config show -h"},
		{false, "khedra config show -help"},
		{false, "khedra config show -v"},
		{false, "khedra config show -version"},
		{false, "khedra config show help"},
		{false, "khedra config show version"},
		{false, "khedra config version"},
		{false, "khedra config"},

		{false, "khedra daemon --h"},
		{false, "khedra daemon --help"},
		{false, "khedra daemon --v"},
		{false, "khedra daemon --version"},
		{false, "khedra daemon -h"},
		{false, "khedra daemon -help"},
		{false, "khedra daemon -v"},
		{false, "khedra daemon -version"},
		{false, "khedra daemon config --h"},
		{false, "khedra daemon config --help"},
		{false, "khedra daemon config --v"},
		{false, "khedra daemon config --version"},
		{false, "khedra daemon config -h"},
		{false, "khedra daemon config -help"},
		{false, "khedra daemon config -v"},
		{false, "khedra daemon config -version"},
		{false, "khedra daemon config help"},
		{false, "khedra daemon config version"},
		{false, "khedra daemon help"},
		{false, "khedra daemon init --h"},
		{false, "khedra daemon init --help"},
		{false, "khedra daemon init --v"},
		{false, "khedra daemon init --version"},
		{false, "khedra daemon init -h"},
		{false, "khedra daemon init -help"},
		{false, "khedra daemon init -v"},
		{false, "khedra daemon init -version"},
		{false, "khedra daemon init config --h"},
		{false, "khedra daemon init config --help"},
		{false, "khedra daemon init config --v"},
		{false, "khedra daemon init config --version"},
		{false, "khedra daemon init config -h"},
		{false, "khedra daemon init config -help"},
		{false, "khedra daemon init config -v"},
		{false, "khedra daemon init config -version"},
		{false, "khedra daemon init config help"},
		{false, "khedra daemon init config version"},
		{false, "khedra daemon init help"},
		{false, "khedra daemon init version"},
		{false, "khedra daemon version"},
		// {false, "khedra daemon"},

		{false, "khedra help config daemon"},
		{false, "khedra help config edit"},
		{false, "khedra help config init daemon"},
		{false, "khedra help config init"},
		{false, "khedra help config show"},
		{false, "khedra help config"},
		{false, "khedra help daemon config"},
		{false, "khedra help daemon init config"},
		{false, "khedra help daemon init"},
		{false, "khedra help daemon"},
		{false, "khedra help init config"},
		{false, "khedra help init daemon config"},
		{false, "khedra help init daemon"},
		{false, "khedra help init"},
		{false, "khedra help"},

		{false, "khedra init --h"},
		{false, "khedra init --help"},
		{false, "khedra init --v"},
		{false, "khedra init --version"},
		{false, "khedra init -h"},
		{false, "khedra init -help"},
		{false, "khedra init -v"},
		{false, "khedra init -version"},
		{false, "khedra init config --h"},
		{false, "khedra init config --help"},
		{false, "khedra init config --v"},
		{false, "khedra init config --version"},
		{false, "khedra init config -h"},
		{false, "khedra init config -help"},
		{false, "khedra init config -v"},
		{false, "khedra init config -version"},
		{false, "khedra init config help"},
		{false, "khedra init config version"},
		{false, "khedra init daemon --h"},
		{false, "khedra init daemon --help"},
		{false, "khedra init daemon --v"},
		{false, "khedra init daemon --version"},
		{false, "khedra init daemon -h"},
		{false, "khedra init daemon -help"},
		{false, "khedra init daemon -v"},
		{false, "khedra init daemon -version"},
		{false, "khedra init daemon config --h"},
		{false, "khedra init daemon config --help"},
		{false, "khedra init daemon config --v"},
		{false, "khedra init daemon config --version"},
		{false, "khedra init daemon config -h"},
		{false, "khedra init daemon config -help"},
		{false, "khedra init daemon config -v"},
		{false, "khedra init daemon config -version"},
		{false, "khedra init daemon config help"},
		{false, "khedra init daemon config version"},
		{false, "khedra init daemon help"},
		{false, "khedra init daemon version"},
		{false, "khedra init help"},
		{false, "khedra init version"},
		// {false, "khedra init"},

		{false, "khedra version config daemon"},
		{false, "khedra version config edit"},
		{false, "khedra version config init daemon"},
		{false, "khedra version config init"},
		{false, "khedra version config show"},
		{false, "khedra version config"},
		{false, "khedra version daemon config"},
		{false, "khedra version daemon init config"},
		{false, "khedra version daemon init"},
		{false, "khedra version daemon"},
		{false, "khedra version init config"},
		{false, "khedra version init daemon config"},
		{false, "khedra version init daemon"},
		{false, "khedra version init"},
		{false, "khedra version"},

		{false, "khedra"},

		{false, "khedra --v config init"},
		{false, "khedra --v config"},
		{false, "khedra --v daemon init"},
		{false, "khedra --v daemon"},
		{false, "khedra --v init config"},
		{false, "khedra --v init daemon"},
		{false, "khedra --v init"},

		{true, "khedra --not-a-flag config daemon"},
		{true, "khedra --not-a-flag config edit"},
		{true, "khedra --not-a-flag config init daemon"},
		{true, "khedra --not-a-flag config init"},
		{true, "khedra --not-a-flag config show"},
		{true, "khedra --not-a-flag config"},
		{true, "khedra --not-a-flag daemon config"},
		{true, "khedra --not-a-flag daemon init config"},
		{true, "khedra --not-a-flag daemon init"},
		{true, "khedra --not-a-flag daemon"},
		{true, "khedra --not-a-flag init config"},
		{true, "khedra --not-a-flag init daemon config"},
		{true, "khedra --not-a-flag init daemon"},
		{true, "khedra --not-a-flag init"},
		{true, "khedra --not-a-flag"},

		{true, "khedra config --not-a-flag"},
		{true, "khedra config daemon --not-a-flag"},
		{true, "khedra config daemon not-a-command"},
		{true, "khedra config daemon"},
		{true, "khedra config edit --not-a-flag"},
		{true, "khedra config edit not-a-command"},
		// {true, "khedra config edit"},
		{true, "khedra config init --not-a-flag"},
		{true, "khedra config init daemon --not-a-flag"},
		{true, "khedra config init daemon not-a-command"},
		{true, "khedra config init daemon"},
		{true, "khedra config init not-a-command"},
		{true, "khedra config init"},
		{true, "khedra config not-a-command"},
		{true, "khedra config show --not-a-flag"},
		{true, "khedra config show not-a-command"},
		// {true, "khedra config show"},

		{true, "khedra daemon --not-a-flag"},
		{true, "khedra daemon config --not-a-flag"},
		{true, "khedra daemon config not-a-command"},
		{true, "khedra daemon config"},
		{true, "khedra daemon init --not-a-flag"},
		{true, "khedra daemon init config --not-a-flag"},
		{true, "khedra daemon init config not-a-command"},
		{true, "khedra daemon init config"},
		{true, "khedra daemon init not-a-command"},
		{true, "khedra daemon init"},
		{true, "khedra daemon not-a-command"},

		{true, "khedra init --not-a-flag"},
		{true, "khedra init config --not-a-flag"},
		{true, "khedra init config not-a-command"},
		{true, "khedra init config"},
		{true, "khedra init daemon --not-a-flag"},
		{true, "khedra init daemon config --not-a-flag"},
		{true, "khedra init daemon config not-a-command"},
		{true, "khedra init daemon config"},
		{true, "khedra init daemon not-a-command"},
		{true, "khedra init daemon"},
		{true, "khedra init not-a-command"},

		{true, "khedra not-a-command config daemon"},
		{true, "khedra not-a-command config edit"},
		{true, "khedra not-a-command config init daemon"},
		{true, "khedra not-a-command config init"},
		{true, "khedra not-a-command config show"},
		{true, "khedra not-a-command config"},
		{true, "khedra not-a-command daemon config"},
		{true, "khedra not-a-command daemon init config"},
		{true, "khedra not-a-command daemon init"},
		{true, "khedra not-a-command daemon"},
		{true, "khedra not-a-command init config"},
		{true, "khedra not-a-command init daemon config"},
		{true, "khedra not-a-command init daemon"},
		{true, "khedra not-a-command init"},
		{true, "khedra not-a-command"},
	}
	for i, tc := range testCases {
		isGithub := os.Getenv("TB_GITHUB_TESTING") == "true"
		if !isGithub && i%30 != 0 {
			continue
		}

		t.Run(tc.command, func(t *testing.T) {
			os.Args = strings.Split(tc.command, " ")
			os.Args = cleanArgs()
			k := &KhedraApp{}
			cmdLine := initCli(k)
			if len(os.Args) > 0 {
				commandName := os.Args[1]
				command := getCommandByName(t, cmdLine, commandName)
				if command != nil {
					assert.NotNil(t, command)
					c := cli.NewContext(cmdLine, nil, nil)
					if command.Action == nil {
						fmt.Println("Command.Action not found", tc.command, commandName)
					} else {
						err := command.Action(c)
						if tc.expectError {
							assert.Error(t, err)
						} else {
							assert.NoError(t, err)
						}
					}
				} else {
					fmt.Println("Command not found", tc.command, commandName)
				}
			}
		})
	}
}

func getCommandByName(t *testing.T, cmdLine *cli.App, name string) *cli.Command {
	t.Helper()
	for _, command := range cmdLine.Commands {
		if command.Name == name {
			return command
		}
	}
	return nil
}

func getSubCommandByName(t *testing.T, command *cli.Command, name string) *cli.Command {
	t.Helper()
	for _, subcommand := range command.Subcommands {
		if subcommand.Name == name {
			return subcommand
		}
	}
	return nil
}

func captureOutput(t *testing.T, f func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}
