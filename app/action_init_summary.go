package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getSummaryScreen(cfg *types.Config) wizard.Screen {
	_ = cfg // linter
	var summaryTitle = `Summary`
	var summarySubtitle = ``
	var summaryInstructions = `
Press enter to finish the wizard. ("b"=back, "h"=help)`
	var summaryBody = `
You've completed the wizard and your settings have been saved to the
configuation file at {cfg.General.DataFolder}.

You may re-run this wizard at any time to edit or modify the config, however
not all options are configurable. You may run khedra config edit or type
edit here to open the the actual file in your editor.
`
	var summaryReplacements = []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{summaryTitle}},
	}
	var summaryQuestions = []wizard.Question{summaryQ1}

	return wizard.Screen{
		Title:        summaryTitle,
		Subtitle:     summarySubtitle,
		Instructions: summaryInstructions,
		Body:         summaryBody,
		Replacements: summaryReplacements,
		Questions:    summaryQuestions,
		Style:        wizard.NewStyle(),
	}
}

var summaryQ1 = wizard.Question{
	Text:  "Would you like to edit the config by hand?",
	Value: "no",
}
