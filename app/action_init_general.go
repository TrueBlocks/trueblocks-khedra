package app

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getGeneralScreen(cfg *types.Config) wizard.Screen {
	var generalTitle = `General Settings`
	var generalSubtitle = ``
	var generalInstructions = `
Type your answer and press enter. ("q"=quit, "b"=back, "h"=help)`
	var generalBody = `
The General group of options controls where Khedra stores the Unchained
Index and its caches. It also helps you choose a download strategy for
the index and helps you set up Khedra's logging options.

Choose your folders carefully. The index and logs can get quite large
depending on the configuration. As always, type "help" to get more
information.

You may use $HOME or ~/ in your paths to refer to your home directory.`
	var generalReplacements = []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{generalTitle}},
		{Color: colors.Green, Values: []string{"Unchained\nIndex", "Unchained Index", "$HOME", "~/"}},
	}
	var generalQuestions = []wizard.Question{generalQ0, generalQ1, generalQ2, generalQ3, generalQ4, generalQ5}

	for i := 0; i < len(generalQuestions); i++ {
		q := &generalQuestions[i]
		q.Value = strings.ReplaceAll(q.Value, "{cfg.General.DataFolder}", cfg.General.DataFolder)
		for j := 0; j < len(q.Messages); j++ {
			q.Messages[j] = strings.ReplaceAll(q.Messages[j], "{cfg.Logging.Filename}", cfg.Logging.Filename)
		}
	}

	return wizard.Screen{
		Title:        generalTitle,
		Subtitle:     generalSubtitle,
		Instructions: generalInstructions,
		Body:         generalBody,
		Replacements: generalReplacements,
		Questions:    generalQuestions,
		Style:        wizard.NewStyle(),
	}
}

// --------------------------------------------------------
var generalQ0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}

// --------------------------------------------------------
var generalQ1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Would you like to create the Unchained Index from scratch
|(starting at block zero) or download from IPFS?`,
	Hint: `Downloading is faster (a few hours). Building from scratch is
|more secure but much slower (depending on the chain, perhaps as
|long as a few days).`,
	Value: "download",
	Validate: func(input string, q *wizard.Question) (string, error) {
		switch input {
		case "download":
			return input, validOk(q.Messages[0], input)
		case "scratch":
			return input, validOk(q.Messages[1], input)
		default:
			return input, fmt.Errorf(q.Messages[2]+"%w", wizard.ErrValidate)
		}
	},
	Messages: []string{
		`the index will be downloaded`,
		`the index will be built from scratch`,
		`value must be either "download" or "scratch"`,
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"download", "scratch"}},
	},
}

// --------------------------------------------------------
var generalQ2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to download only bloomFilters or the entireIndex?`,
	Hint: `Downloading blooms takes less time and is smaller (4gb), but is
|slower when searching. Downloading the entire index takes longer
|and is larger (180gb), but is much faster during search.`,
	Value: "entireIndex",
	Validate: func(input string, q *wizard.Question) (string, error) {
		switch input {
		case "bloomFilters":
			return input, validOk(q.Messages[0], input)
		case "entireIndex":
			return input, validOk(q.Messages[1], input)
		default:
			return input, fmt.Errorf(q.Messages[2]+"%w", wizard.ErrValidate)
		}
	},
	Messages: []string{
		`only bloom filters will be downloaded`,
		`both bloom filters and index chunks will be downloaded`,
		`value must be either "bloomFilters" or "entireIndex"`,
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"bloomFilters", "entireIndex"}},
	},
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		if q.Screen.Questions[0].Value == "scratch" {
			return input, validSkipNext(`question skipped`, input)
		}
		return input, validOk(`question proceeds`, input)
	},
}

// --------------------------------------------------------
var generalQ3 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Where do you want to store the Unchained Index and the
|binary caches?`,
	Value: `{cfg.General.DataFolder}`,
	Validate: func(input string, q *wizard.Question) (string, error) {
		path, err := utils.ResolveValidPath(input)
		if err != nil {
			return input, fmt.Errorf(q.Messages[2]+"%w", path, wizard.ErrValidate)
		}

		if !file.FolderExists(path) {
			file.EstablishFolder(path)
			return input, validWarn(q.Messages[0], path)
		}

		return input, validOk(q.Messages[1], path)
	},
	Messages: []string{
		`"%s" was created`,
		`the index will be stored at %s`,
		"unable to create folder: %s",
	},
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		if q.Screen.Questions[2].Value == "bloomFilters" {
			q.Hint = `The bloom filters take up about 5-10gb and the caches may get
|quite large depending on your usage, so choose a folder where you
|can store up to 100gb.`
			q.Hint = strings.ReplaceAll(q.Hint, "\n|", "\n          ")
			return input, validOk(`question proceeds`, input)
		}
		q.Hint = `The index takes up about 120-150gb and the caches may get quite
|large depending on your usage, so choose a folder where you can
|store up to 300gb.`
		q.Hint = strings.ReplaceAll(q.Hint, "\n|", "\n          ")
		return input, validOk(`question proceeds`, input)
	},
}

// --------------------------------------------------------
var generalQ4 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: "Do you want to enable file-based logging?",
	Value:    "no",
	Hint: `Logging to the screen is always enabled. If you enable file-based
|logging, Khedra will also write log files to disk.`,
	Validate: func(input string, q *wizard.Question) (string, error) {
		prevScreen := 2
		path := filepath.Join(q.Screen.Questions[prevScreen].Value, q.Messages[3])
		switch input {
		case "yes":
			return input, validOk(q.Messages[0], path)
		case "no":
			return input, validOk(q.Messages[1], path)
		default:
			return input, fmt.Errorf(q.Messages[2]+"%w", wizard.ErrValidate)
		}
	},
	Messages: []string{
		`logs will be stored at %s`,
		`logs will be reported to screen only`,
		`value must be either "yes" or "no"`,
		`{cfg.Logging.Filename}`,
	},
}

// --------------------------------------------------------
var generalQ5 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: "What log level do you want to enable (debug, info, warn, error)?",
	Value:    "info",
	Validate: func(input string, q *wizard.Question) (string, error) {
		if input != "debug" && input != "info" && input != "warn" && input != "error" {
			err := fmt.Errorf(`value must be either "debug", "info", "warn", or "error"%w`, wizard.ErrValidate)
			return input, err
		}
		msg := fmt.Errorf(q.Messages[0]+"%w", input, wizard.ErrValidateMsg)
		return input, msg
	},
	Messages: []string{
		`logging level will be "%s"`,
	},
}
