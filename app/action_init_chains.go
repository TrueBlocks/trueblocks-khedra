package app

import (
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getChainsScreen(cfg *types.Config) wizard.Screen {
	_ = cfg // linter	// linter
	var chainsTitle = `Chain Settings`
	var chainsSubtitle = ``
	var chainsInstructions = `Type your answers and press enter. ("b"=back, "q"=quit)`
	var chainsBody = `
Khedra will index any number of EVM chains, however it requires an
RPC endpoint for each to do so. Fast, dedicated local endpoints are
preferred. Likely, you will get rate limited if you point to a remote
endpoing, but if you do, you may use the Sleep option to slow down
operation. See "help".

You may add chains to the list by typing the chain's name. Remove chains
with "remove <chain>". Or, an easier way is to edit the configuration
file directly by typing "edit". The mainnet chain is required.
`
	var chainsReplacements = []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{chainsTitle}},
		{Color: colors.Green, Values: []string{"remove <chain>", "Sleep", "mainnet", "\"edit\""}},
	}

	var chainsQuestions = []wizard.Question{chainsQ0, chainsQ1}
	for i := 0; i < len(chainsQuestions); i++ {
		chainsQuestions[i].Value = strings.ReplaceAll(chainsQuestions[i].Value, "{cfg.ChainList}", cfg.ChainList())
	}

	return wizard.Screen{
		Title:        chainsTitle,
		Subtitle:     chainsSubtitle,
		Body:         chainsBody,
		Instructions: chainsInstructions,
		Replacements: chainsReplacements,
		Questions:    chainsQuestions,
		Style:        wizard.NewStyle(),
	}
}

// --------------------------------------------------------
var chainsQ0 = wizard.Question{}

var chainsQ1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Text: `Which chains do you want to index? (Enter a chain's name
	directly to add chains or "remove <chain>" to remove them.)`,
	Value: "{cfg.ChainList}",
}
