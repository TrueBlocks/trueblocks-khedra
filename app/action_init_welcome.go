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
Welcome to Khedra, the world's only local-first indexer/monitor for
EVM blockchains. This wizard will help you configure Khedra. There are
three groups of settings: General, Services, and Chains.

Type "q" or "quit" to quit, "b" or "back" to return to a previous screen,
or "help" to get more information.
`
	wReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{wTitle}},
		{Color: colors.Green, Values: []string{
			"\"q\"", "\"quit\"", "\"b\"", "\"back\"", "\"help\"",
		}},
	}
	wStyle := wizard.NewStyle()
	wStyle.Justify = boxes.Center

	return wizard.Screen{
		Title:        wTitle,
		Subtitle:     wSubtitle,
		Instructions: wInstructions,
		Body:         wBody,
		Replacements: wReplacements,
		Style:        wStyle,
	}
}
