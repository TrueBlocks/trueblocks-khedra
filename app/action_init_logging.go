package app

import (
	"fmt"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/wizard"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getLoggingScreen() wizard.Screen {
	lTitle := `Logging Settings`
	lSubtitle := ``
	lInstructions := `Type your answer and press enter. ("q"=quit, "b"=back, "h"=help)`
	lBody := `
The Logging group of options helps you set up Khedra's logging options.

Choose your folders carefully. The logs can get quite large depending
on the configuration. As always, type "help" to get more information.

You may use $HOME or ~/ in your paths to refer to your home directory.
`
	lReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{lTitle}},
		{Color: colors.Green, Values: []string{"$HOME", "~/"}},
	}
	lQuestions := []wizard.Questioner{&l0, &l1, &l2}
	lStyle := wizard.NewStyle()

	return wizard.Screen{
		Title:        lTitle,
		Subtitle:     lSubtitle,
		Body:         lBody,
		Instructions: lInstructions,
		Replacements: lReplacements,
		Questions:    lQuestions,
		Style:        lStyle,
	}
}

// --------------------------------------------------------
func lValidate(key string, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		switch key {
		case "enable":
			switch input {
			case "yes":
				cfg.Logging.ToFile = true
				err := cfg.WriteToFile(types.GetConfigFnNoCreate())
				if err != nil {
					fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
				}
				path := filepath.Join(cfg.General.DataFolder, cfg.Logging.Filename)
				return input, validOk(`logs will be stored at %s`, path)
			case "no":
				cfg.Logging.ToFile = false
				err := cfg.WriteToFile(types.GetConfigFnNoCreate())
				if err != nil {
					fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
				}
				return input, validOk(`logs will be reported to screen only`, cfg.Logging.Filename)
			default:
				return input, fmt.Errorf(`value must be either "yes" or "no" %w`, wizard.ErrValidate)
			}
		case "level":
			if input != "debug" && input != "info" && input != "warn" && input != "error" {
				err := fmt.Errorf(`value must be either "debug", "info", "warn", or "error"%w`, wizard.ErrValidate)
				return input, err
			}
			cfg.Logging.Level = input
			err := cfg.WriteToFile(types.GetConfigFnNoCreate())
			if err != nil {
				fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
			}
			msg := fmt.Errorf(`logging level will be "%s" %w`, input, wizard.ErrValidateMsg)
			return input, msg
		}
	}
	return input, validOk(`don't skip`, input)
}

// --------------------------------------------------------
var l0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}

// --------------------------------------------------------
var l1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable file-based logging?`,
	Hint: `Logging to the screen is always enabled. If you enable file-based
|logging, Khedra will also write log files to disk.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return prepare[types.Logging](q, func(cfg *types.Config) (string, types.Logging, error) {
			copy := types.Logging{ToFile: cfg.Logging.ToFile, Filename: cfg.Logging.Filename}
			if cfg.Logging.ToFile {
				return "yes", copy, validContinue()
			}
			return "no", copy, validContinue()
		})
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return lValidate("enable", input, q)
	},
}

// --------------------------------------------------------
var l2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `What log level do you want to enable (debug, info, warn, error)?`,
	Hint:     `We need a hint`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return prepare[types.Logging](q, func(cfg *types.Config) (string, types.Logging, error) {
			copy := types.Logging{Level: cfg.Logging.Level}
			return cfg.Logging.Level, copy, validContinue()
		})
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return lValidate("level", input, q)
	},
}
