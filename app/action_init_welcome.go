package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/boxes"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func getWelcomeScreen(cfg *types.Config) wizard.Screen {
	_ = cfg // linter
	var welcomeScreen = wizard.Screen{
		Title:    `KHEDRA WIZARD`,
		Subtitle: `Index, monitor, serve, and share blockchain data.`,
		Body: `
Welcome to Khedra, the world's only local-first indexer/monitor for
EVM blockchains. This wizard will help you configure Khedra. There are
three groups of settings: General, Services, and Chains.

Type "q" or "quit" to quit, "b" or "back" to return to a previous screen,
or "help" to get more information.
`,
		Instructions: `Press enter to continue.`,
		Replacements: []wizard.Replacement{
			{Color: colors.Yellow, Values: []string{"KHEDRA WIZARD"}},
			{Color: colors.Green, Values: []string{"\"q\"", "\"quit\"", "\"b\"", "\"back\"", "\"help\""}},
		},
		Style: wizard.Style{
			Outer:   boxes.Single | boxes.All,
			Inner:   boxes.Double | boxes.All,
			Justify: boxes.Center,
		},
	}

	return welcomeScreen
}
