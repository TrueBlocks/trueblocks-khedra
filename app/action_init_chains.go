package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func getChainsScreen(cfg *types.Config) wizard.Screen {
	_ = cfg // linter
	var chainsScreen = wizard.Screen{
		Title:    `Chains Settings`,
		Subtitle: ``,
		Body: `
You may use Khedra against any EVM chain, however in order for it to
work, you must provide an RPC endpoint for the chain. Generally, the
tools work much better with a locally running node, but it does work
with remove RPCs (just not as fast).
`,
		Instructions: `Type your answers and press enter. ("b"=back, "q"=quit)`,
		Replacements: []wizard.Replacement{
			{Color: colors.Yellow, Values: []string{"Chains Settings", "remove <chain>"}},
		},
		Questions: []wizard.Question{
			{
				Text: `Which chains do you want to index? (Enter a chain's name directly
            to add chains or "remove <chain>" to remove them.)`,
				Value: "mainnet, gnosis, sepolia, optimism",
			},
		},
		Style: wizard.NewStyle(),
	}

	return chainsScreen
}
