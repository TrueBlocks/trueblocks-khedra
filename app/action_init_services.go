package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func getServicesScreen(cfg *types.Config) wizard.Screen {
	_ = cfg // linter
	var servicesScreen = wizard.Screen{
		Title:    `Services Settings`,
		Subtitle: ``,
		Body: `
This wizard helps you configure Khedra. It will walk you through four
sections General, Services, Chains, and Logging.

You may quit the wizard at any time by typing "q" or "quit". The next
time you run it, it will continue where you left off. Type "help" at any
point to get more information.
`,
		Instructions: `Type your answers and press enter. ("b"=back, "q"=quit)`,
		Replacements: []wizard.Replacement{
			{Color: colors.Yellow, Values: []string{"Services Settings"}},
		},
		Style: wizard.NewStyle(),
	}

	return servicesScreen
}
