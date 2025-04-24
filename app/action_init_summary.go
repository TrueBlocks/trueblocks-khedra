package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/wizard"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getSummaryScreen() wizard.Screen {
	sumTitle := `Summary`
	sumSubtitle := ``
	sumInstructions := `
Press enter to finish the wizard. ("b"=back, "h"=help)`
	sumBody := `
You've completed the wizard and your settings have been saved to the
configuation file at {cfg.General.DataFolder}.

You may re-run this wizard at any time to edit or modify the config, however
not all options are configurable. You may run khedra config edit or type
edit here to open the the actual file in your editor.
`
	sumReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{sumTitle}},
	}
	sumQuestions := []wizard.Questioner{&sum0}
	sumStyle := wizard.NewStyle()

	return wizard.Screen{
		Title:        sumTitle,
		Subtitle:     sumSubtitle,
		Instructions: sumInstructions,
		Body:         sumBody,
		Replacements: sumReplacements,
		Questions:    sumQuestions,
		Style:        sumStyle,
	}
}

// --------------------------------------------------------
var sum0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}
