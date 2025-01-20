package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func getSummaryScreen(cfg *types.Config) wizard.Screen {
	_ = cfg // linter
	var summaryScreen = wizard.Screen{
		Title:    `Summary`,
		Subtitle: ``,
		Body: `
This wizard helps you configure Khedra. It will walk you through four
sections General, Services, Chains, and Logging.

You may quit the wizard at any time by typing "q" or "quit". The next
time you run it, it will continue where you left off. Type "help" at any
point to get more information.
`,
		Instructions: `Select an option or hit enter. ("h"=help, "e"=edit, "q"=quit)`,
		Replacements: []wizard.Replacement{
			{Color: colors.Yellow, Values: []string{"Summary"}},
		},
		Questions: []wizard.Question{
			{
				Text:  "Would you like to edit the config by hand?",
				Value: "no",
			},
		},
		Style: wizard.NewStyle(),
	}

	return summaryScreen
}
