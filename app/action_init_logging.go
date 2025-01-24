package app

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
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
	lQuestions := []wizard.Question{l0, l1, l2}
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
func lPrepare(key, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		lCopy := types.Logging{
			Filename: cfg.Logging.Filename,
			Level:    cfg.Logging.Level,
		}
		switch key {
		case "enable":
			if len(cfg.Logging.Filename) > 0 {
				input = "yes"
			} else {
				input = "no"
			}
		case "level":
			input = cfg.Logging.Level
		}
		bytes, _ := json.Marshal(lCopy)
		q.State = string(bytes)
	}
	return input, validOk(`don't skip`, input)
}

// --------------------------------------------------------
func lValidate(key string, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		switch key {
		case "enable":
			msgs := []string{
				`logs will be stored at %s`,
				`logs will be reported to screen only`,
				`value must be either "yes" or "no"`,
			}
			switch input {
			case "yes":
				cfg.Logging.Filename = "khedra.log"
				path := filepath.Join(cfg.General.DataFolder, cfg.Logging.Filename)
				return input, validOk(msgs[0], path)
			case "no":
				cfg.Logging.Filename = ""
				return input, validOk(msgs[1], cfg.Logging.Filename)
			default:
				return input, fmt.Errorf(msgs[2]+"%w", wizard.ErrValidate)
			}
		case "level":
			msgs := []string{
				`logging level will be "%s"`,
			}
			if input != "debug" && input != "info" && input != "warn" && input != "error" {
				err := fmt.Errorf(`value must be either "debug", "info", "warn", or "error"%w`, wizard.ErrValidate)
				return input, err
			}
			cfg.Logging.Level = input
			msg := fmt.Errorf(msgs[0]+"%w", input, wizard.ErrValidateMsg)
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
		return lPrepare("enable", input, q)
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
		return lPrepare("level", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return lValidate("level", input, q)
	},
}
