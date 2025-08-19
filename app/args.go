package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// parseArgsInternal processes command-line arguments using os.Args, identifying help and version
// variants, separating the remaining arguments into a non-duplicative array (excluding the
// program name), and counting non-flag (i.e., command) arguments.
func parseArgsInternal() (hasHelp bool, hasVersion bool, argsOut []string, commandCount int) {
	argsOut = []string{}
	seen := map[string]bool{}

	args := os.Args
	if len(args) < 2 {
		hasHelp = true
		return
	}
	args = args[1:]

	helpForms := map[string]bool{
		"--help": true, "-help": true, "help": true,
		"--h": true, "-h": true,
	}

	versionForms := map[string]bool{
		"--version": true, "-version": true, "version": true,
		"--v": true, "-v": true,
	}

	for _, arg := range args {
		arg = strings.TrimSpace(arg)
		if len(arg) == 0 {
			continue
		}
		if helpForms[arg] {
			hasHelp = true
			continue
		}
		if versionForms[arg] {
			hasVersion = true
			continue
		}
		if !seen[arg] {
			argsOut = append(argsOut, arg)
			seen[arg] = true
			if arg[0] != '-' {
				commandCount++
			}
		}
	}

	return
}

// cleanArgs processes command-line arguments to produce a cleaned set by appending
// "help" if a help flag is present, "version" if a version flag is present, or preserving
// commands and flags otherwise. Defaults to "help" when no arguments are provided.
func cleanArgs() []string {
	programName := os.Args[:1] // program name
	hasHelp, hasVersion, cleanedArgs, _ := parseArgsInternal()

	if hasHelp {
		if len(cleanedArgs) > 0 {
			return append(programName, "help", cleanedArgs[0])
		}
		return append(programName, "help")
	}

	if hasVersion {
		return append(programName, "version")
	}

	return append(programName, cleanedArgs...)
}

// validateArgs checks if the number of commands and arguments matches the expected counts.
// Returns an error if the counts are incorrect or if any unknown commands are found.
func validateArgs(expectedCmdCount, expectedArgCount int) error {
	_, _, cleanedArgs, cmdCnt := parseArgsInternal()
	unknown := getUnknownCmd()
	if len(unknown) > 0 {
		return fmt.Errorf("command '%s' not found", unknown)
	}

	if cmdCnt != expectedCmdCount || len(cleanedArgs) != expectedArgCount {
		if cmdCnt > expectedCmdCount {
			return fmt.Errorf("too many commands provided, expected %d", expectedCmdCount)
		}

		if len(cleanedArgs) > expectedArgCount {
			return fmt.Errorf("too many arguments provided, expected %d", expectedArgCount)
		}

		return fmt.Errorf(
			"unexpected number of commands (%d/%d) or arguments (%d/%d)",
			cmdCnt, expectedCmdCount, len(cleanedArgs), expectedArgCount,
		)
	}

	return nil
}

// getUnknownCmd identifies the first invalid argument (if any) from the command-line
// input. It checks each argument against a dynamically extracted list of valid commands
// and flags. Returns the first unrecognized command or flag as a string, or an empty
// string if all arguments are valid.
func getUnknownCmd() string {
	validCommands, validFlags := extractCmdsAndFlags()
	for i, arg := range os.Args {
		if i > 0 {
			isFlag := strings.HasPrefix(arg, "-")
			if isFlag {
				if !validFlags[arg] {
					return arg
				}
			} else {
				if !validCommands[arg] {
					return arg
				}
			}
		}
	}
	return ""
}

// extractCmdsAndFlags extracts valid commands and flags from the CLI application structure.
// It traverses the `cli.App` configuration to build maps of recognized commands and flags,
// including subcommands and their associated flags.
// Returns two maps: one for valid commands and one for valid flags.
func extractCmdsAndFlags() (map[string]bool, map[string]bool) {
	var tmpK = &KhedraApp{}
	tmpCli := initCli(tmpK)
	validCommands := map[string]bool{}
	validFlags := map[string]bool{}

	var processCommands func(commands []*cli.Command)
	processCommands = func(commands []*cli.Command) {
		for _, cmd := range commands {
			validCommands[cmd.Name] = true
			for _, flag := range cmd.Flags {
				if name := flag.Names(); len(name) > 0 {
					validFlags["--"+name[0]] = true
				}
			}
			if len(cmd.Subcommands) > 0 {
				processCommands(cmd.Subcommands)
			}
		}
	}

	processCommands(tmpCli.Commands)
	validCommands["help"] = true
	return validCommands, validFlags
}
