package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/boxes"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/wizard"
)

func getWelcomeScreen() wizard.Screen {
	wTitle := `KHEDRA WIZARD`
	wSubtitle := `Index, monitor, serve, and share blockchain data.`
	wInstructions := `Press enter to continue.`
	wBody := `
Welcome to Khedra, a local-first indexer/monitor for EVM blockchains. This
wizard will walk you through step by step to config the app.

Type "help" at any time, "q" for "quit" to quit, "b" or "back" to return
to a previous screen, or "edit" to open the configuration file.
`
	wReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{wTitle}},
		{Color: colors.Green, Values: []string{
			"\"q\"", "\"quit\"", "\"b\"", "\"back\"", "\"help\"", "\"edit\"", "Khedra",
		}},
	}
	wQuestions := []wizard.Questioner{&w0}
	wStyle := wizard.NewStyle()
	wStyle.Justify = boxes.Center

	return wizard.Screen{
		Title:        wTitle,
		Subtitle:     wSubtitle,
		Instructions: wInstructions,
		Body:         wBody,
		Questions:    wQuestions,
		Replacements: wReplacements,
		Style:        wStyle,
	}
}

// --------------------------------------------------------
var w0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}
