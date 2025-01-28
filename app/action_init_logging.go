package app

import (
	"fmt"

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
		return confirm[types.Logging](q, func(cfg *types.Config) (string, types.Logging, error) {
			copy := types.Logging{Filename: cfg.Logging.Filename, ToFile: cfg.Logging.ToFile}
			switch input {
			case "no":
				cfg.Logging.ToFile = false
				copy.ToFile = cfg.Logging.ToFile
				return input, copy, validOk(`logs will be reported to screen only`, "")
			case "yes":
				cfg.Logging.ToFile = true
				copy.ToFile = cfg.Logging.ToFile
				return input, copy, validOk(`logs will be stored at %s`, cfg.Logging.Filename)
			}
			return input, copy, fmt.Errorf(`value must be either "yes" or "no" %w`, wizard.ErrValidate)
		})
	},
}

// --------------------------------------------------------
var l2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `What log level do you want to enable (debug, info, warn, error)?`,
	Hint:     `Select a log level from the list.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return prepare[types.Logging](q, func(cfg *types.Config) (string, types.Logging, error) {
			copy := types.Logging{Level: cfg.Logging.Level}
			return cfg.Logging.Level, copy, validContinue()
		})
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return confirm[types.Logging](q, func(cfg *types.Config) (string, types.Logging, error) {
			copy := types.Logging{Level: input}
			if input != "debug" && input != "info" && input != "warn" && input != "error" {
				err := fmt.Errorf(`value must be either "debug", "info", "warn", or "error"%w`, wizard.ErrValidate)
				return input, copy, err
			}
			cfg.Logging.Level = input
			return input, copy, validOk(`logging level will be "%s"`, input)
		})
	},
}
