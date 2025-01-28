package app

import (
	"fmt"

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
	gQuestions := []wizard.Questioner{&g0, &g1, &g2, &g3}
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
		return prepare[types.General](q, func(cfg *types.Config) (string, types.General, error) {
			copy := types.General{Strategy: cfg.General.Strategy}
			return copy.Strategy, copy, nil
		})
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return confirm[types.General](q, func(cfg *types.Config) (string, types.General, error) {
			cfg.General.Strategy = input
			copy := types.General{Strategy: cfg.General.Strategy}
			switch input {
			case "download":
				return input, copy, validOk(`the index will be downloaded`, input)
			case "scratch":
				cfg.General.Detail = "index"
				copy.Detail = "index"
				return input, copy, validOk(`the index will be built from scratch`, input)
			}
			return input, copy, fmt.Errorf(`value must be either "download" or "scratch" %w`, wizard.ErrValidate)
		})
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"download", "scratch"}},
	},
}

// --------------------------------------------------------
var g2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to download only bloom filters or the entire index?`,
	Hint: `Downloading bloom fiters takes less time and is smaller (4gb),
|but is slower when searching. Downloading the entire index takes
|longer and is larger (180gb), but is much faster during search.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return prepare[types.General](q, func(cfg *types.Config) (string, types.General, error) {
			copy := types.General{Detail: cfg.General.Detail}
			if cfg.General.Strategy == "scratch" {
				return copy.Detail, copy, validSkipNext()
			}
			return copy.Detail, copy, nil
		})
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return confirm[types.General](q, func(cfg *types.Config) (string, types.General, error) {
			cfg.General.Detail = input
			copy := types.General{Detail: cfg.General.Detail}
			switch input {
			case "bloom":
				cfg.General.Detail = input
				return input, copy, validOk(`only bloom filters will be downloaded`, input)
			case "index":
				cfg.General.Detail = input
				return input, copy, validOk(`both bloom filters and index chunks will be downloaded`, input)
			default:
				return input, copy, fmt.Errorf(`value must be either "bloom" or "index" %w`, wizard.ErrValidate)
			}
		})
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"bloom", "index"}},
	},
}

// --------------------------------------------------------
var g3 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Where do you want to store the Unchained Index and the
|binary caches?`,
	Hint: `<set on load>`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return prepare[types.General](q, func(cfg *types.Config) (string, types.General, error) {
			copy := types.General{DataFolder: cfg.General.DataFolder}
			if cfg.General.Detail == "bloom" {
				q.Hint = bloomHint
			} else {
				q.Hint = indexHint
			}
			return cfg.General.DataFolder, copy, nil
		})
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return confirm[types.General](q, func(cfg *types.Config) (string, types.General, error) {
			cfg.General.DataFolder = input
			copy := types.General{DataFolder: cfg.General.DataFolder}
			path, err := utils.ResolveValidPath(input)
			if err != nil {
				return input, copy, fmt.Errorf("unable to create folder: %s %w", path, wizard.ErrValidate)
			}
			if !file.FolderExists(path) {
				file.EstablishFolder(path)
				return input, copy, validWarn(`"%s" was created`, path)
			}
			return input, copy, validOk(`the index will be stored at %s`, path)
		})
	},
}

// --------------------------------------------------------
var bloomHint = `The bloom filters take up about 5-10gb and the caches may get
|quite large depending on your usage, so choose a folder where you
|can store up to 100gb.`

// --------------------------------------------------------
var indexHint = `The index takes up about 120-150gb and the caches may get quite
|large depending on your usage, so choose a folder where you can
|store up to 300gb.`
