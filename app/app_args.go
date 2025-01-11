package app

import (
	"fmt"
)

func parseArgsInternal(args []string) (hasHelp bool, hasVersion bool, commands []string, nonFlagCount int) {
	commands = []string{}
	if len(args) == 0 {
		hasHelp = true
		return
	}

	helpForms := map[string]bool{
		"--help": true, "-help": true, "help": true,
		"--h": true, "-h": true,
	}

	versionForms := map[string]bool{
		"--version": true, "-version": true, "version": true,
		"--v": true, "-v": true,
	}

	for i, arg := range args {
		if helpForms[arg] {
			hasHelp = true
			continue
		}
		if versionForms[arg] {
			hasVersion = true
			continue
		}
		commands = append(commands, arg)
		if i != 0 && len(arg) == 0 || arg[0] != '-' {
			nonFlagCount++
		}
	}

	return
}

func cleanArgs(args []string) []string {
	programName := args[:1] // program name

	hasHelp, hasVersion, commands, _ := parseArgsInternal(args[1:])
	if hasHelp {
		result := append(programName, "help")
		if len(commands) > 0 {
			return append(result, commands[0])
		}
		return result
	}

	if hasVersion {
		return append(programName, "version")
	}

	return append(programName, commands...)
}

func validateArgs(args []string, expectedCmdCount, expectedFlagCount int) error {
	_, _, flags, cmdCnt := parseArgsInternal(args)
	if cmdCnt != expectedCmdCount || len(flags) != expectedFlagCount {
		if cmdCnt > expectedCmdCount {
			var err error
			if unknown := getUnknownCmd(args); len(unknown) > 0 {
				err = fmt.Errorf("command '%s' not found", unknown)
			} else {
				err = fmt.Errorf("use only one command at a time")
			}
			return err
		}
		return fmt.Errorf("argument mismatch: %v %v %d %d", args, flags, cmdCnt, len(flags))
	}
	return nil
}
