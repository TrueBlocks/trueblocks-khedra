package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getChainsScreen() wizard.Screen {
	cTitle := `Chain Settings`
	cSubtitle := ``
	cInstructions := `Type your answers and press enter. ("b"=back, "q"=quit)`
	cBody := `
Khedra will index any number of EVM chains, however it requires an
RPC endpoint for each to do so. Fast, dedicated local endpoints are
preferred. Likely, you will get rate limited if you point to a remote
endpoing, but if you do, you may use the Sleep option to slow down
operation. See "help".

You may add chains to the list by typing the chain's name. Remove chains
with "remove <chain>". Or, an easier way is to edit the configuration
file directly by typing "edit". The mainnet chain is required.
`
	cReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{cTitle}},
		{Color: colors.Green, Values: []string{
			"remove <chain>", "Sleep", "mainnet", "\"edit\"",
		}},
	}
	cQuestions := []wizard.Question{c0, c1}
	cStyle := wizard.NewStyle()

	return wizard.Screen{
		Title:        cTitle,
		Subtitle:     cSubtitle,
		Body:         cBody,
		Instructions: cInstructions,
		Replacements: cReplacements,
		Questions:    cQuestions,
		Style:        cStyle,
	}
}

// --------------------------------------------------------
func cPrepare(key, input string, q *wizard.Question) (string, error) {
	return input, nil
}

// --------------------------------------------------------
func cValidate(key string, input string, q *wizard.Question) (string, error) {
	return input, nil
}

// --------------------------------------------------------
var c0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}

// --------------------------------------------------------
var c1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Which chains do you want to index? (Enter a chain's name
|directly to add chains or "remove <chain>" to remove them.)`,
}
