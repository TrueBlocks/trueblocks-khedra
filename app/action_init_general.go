package app

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/wizard"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getGeneralScreen() wizard.Screen {
	gTitle := `General Settings`
	gSubtitle := ``
	gInstructions := `Type your answer and press enter. ("q"=quit, "b"=back, "h"=help)`
	gBody := `
The General group of options controls where Khedra stores the Unchained
Index and its caches. It also helps you choose a download strategy for
the index.

Choose your folders carefully. The index and logs can get quite large
depending on the configuration. As always, type "help" to get more
information.

You may use $HOME or ~/ in your paths to refer to your home directory.`
	gReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{gTitle}},
		{Color: colors.Green, Values: []string{
			"Unchained\nIndex", "Unchained Index", "$HOME", "~/",
		}},
	}
	gQuestions := []wizard.Question{g0, g1, g2, g3}
	gStyle := wizard.NewStyle()

	return wizard.Screen{
		Title:        gTitle,
		Subtitle:     gSubtitle,
		Instructions: gInstructions,
		Body:         gBody,
		Replacements: gReplacements,
		Questions:    gQuestions,
		Style:        gStyle,
	}
}

// --------------------------------------------------------
func gPrepare(key, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		show := types.General{}
		switch key {
		case "strategy":
			input = cfg.General.Strategy
			show.Strategy = input
		case "detail":
			input = cfg.General.Detail
			show.Detail = input
			if cfg.General.Strategy == "scratch" {
				return input, validSkipNext(`question skipped`, input)
			}
		case "folder":
			input = cfg.General.DataFolder
			show.DataFolder = input
			if cfg.General.Detail == "bloomFilters" {
				q.Hint = `The bloom filters take up about 5-10gb and the caches may get
|quite large depending on your usage, so choose a folder where you
|can store up to 100gb.`
			} else {
				q.Hint = `The index takes up about 120-150gb and the caches may get quite
|large depending on your usage, so choose a folder where you can
|store up to 300gb.`
			}
			q.Hint = strings.ReplaceAll(q.Hint, "\n|", "\n          ")
		}
		bytes, _ := json.Marshal(show)
		q.State = string(bytes)
	}
	return input, validOk(`don't skip`, input)
}

// --------------------------------------------------------
func gValidate(key string, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		switch key {
		case "strategy":
			msgs := []string{
				`the index will be downloaded`,
				`the index will be built from scratch`,
				`value must be either "download" or "scratch"`,
			}
			switch input {
			case "download":
				cfg.General.Strategy = input
				return input, validOk(msgs[0], input)
			case "scratch":
				cfg.General.Strategy = input
				return input, validOk(msgs[1], input)
			default:
				return input, fmt.Errorf(msgs[2]+"%w", wizard.ErrValidate)
			}
		case "detail":
			msgs := []string{
				`only bloom filters will be downloaded`,
				`both bloom filters and index chunks will be downloaded`,
				`value must be either "bloomFilters" or "entireIndex"`,
			}
			switch input {
			case "bloomFilters":
				cfg.General.Detail = input
				return input, validOk(msgs[0], input)
			case "entireIndex":
				cfg.General.Detail = input
				return input, validOk(msgs[1], input)
			default:
				return input, fmt.Errorf(msgs[2]+"%w", wizard.ErrValidate)
			}
		case "folder":
			msgs := []string{
				`"%s" was created`,
				`the index will be stored at %s`,
				"unable to create folder: %s",
			}
			path, err := utils.ResolveValidPath(input)
			if err != nil {
				return input, fmt.Errorf(msgs[2]+"%w", path, wizard.ErrValidate)
			}

			cfg.General.DataFolder = input
			if !file.FolderExists(path) {
				file.EstablishFolder(path)
				return input, validWarn(msgs[0], path)
			}
			return input, validOk(msgs[1], path)
		}
	}
	return input, validOk(`don't skip`, input)
}

// --------------------------------------------------------
var g0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}

// --------------------------------------------------------
var g1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Would you like to create the Unchained Index from scratch
|(starting at block zero) or download from IPFS?`,
	Hint: `Downloading is faster (a few hours). Building from scratch is
|more secure but much slower (depending on the chain, perhaps as
|long as a few days).`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return gPrepare("strategy", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return gValidate("strategy", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"download", "scratch"}},
	},
}

// --------------------------------------------------------
var g2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to download only bloomFilters or the entireIndex?`,
	Hint: `Downloading blooms takes less time and is smaller (4gb), but is
|slower when searching. Downloading the entire index takes longer
|and is larger (180gb), but is much faster during search.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return gPrepare("detail", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return gValidate("detail", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"bloomFilters", "entireIndex"}},
	},
}

// --------------------------------------------------------
var g3 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Where do you want to store the Unchained Index and the
|binary caches?`,
	Hint: `<set on load>`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return gPrepare("folder", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return gValidate("folder", input, q)
	},
}
