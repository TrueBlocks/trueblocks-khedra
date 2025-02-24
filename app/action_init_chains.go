package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/wizard"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getChainsScreen() wizard.Screen {
	cTitle := `Chain Settings`
	cSubtitle := ``
	cInstructions := `Type your answers and press enter. ("b"=back, "q"=quit)`
	cBody := `
Khedra will index any EVM chain. The only requirement is a working RPC
endpoint. You may index more than one chain.

The wizard allows you to enter the name of a chain and then asks you for
an RPC endpoint for that chain. It won't proceed until you provide one.
An Ethereum "mainnet" RPC is required.

The code prefers fast local endpoints, although remote endpoints do work.
If you are rate limited (likely), use the sleep option. See "help".
`
	cReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{cTitle}},
		{Color: colors.Green, Values: []string{
			"sleep", "mainnet", "\"edit\"",
		}},
	}
	cQuestions := []wizard.Questioner{&c0, &c1, &c2}
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
var c0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}

// --------------------------------------------------------
var c1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Please provide an RPC for Ethereum mainnet.`,
	Hint: `Khedra requires a valid, reachable RPC for Mainnet
|Ethereum. It must read state from the Unchained Index smart
|contract. When you press enter, the RPC will be validated.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		_ = input
		return prepare[types.Chain](q, func(cfg *types.Config) (string, types.Chain, error) {
			if _, ok := cfg.Chains["mainnet"]; !ok {
				cfg.Chains["mainnet"] = types.NewChain("mainnet", 1)
			}
			copy := cfg.Chains["mainnet"]
			copy.Name = ""
			return strings.Join(copy.RPCs, ","), copy, validContinue()
		})
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		_ = input
		return confirm[types.Chain](q, func(cfg *types.Config) (string, types.Chain, error) {
			copy, ok := cfg.Chains["mainnet"]
			if !ok {
				log.Fatal("chain mainnet not found")
			}
			copy.RPCs = strings.Split(input, ",")
			if !copy.HasValidRpc(2) {
				copy.Name = ""
				return strings.Join(copy.RPCs, ","), copy, fmt.Errorf(`no rpcs for chain mainnet %w`, wizard.ErrValidate)
			}
			cfg.Chains["mainnet"] = copy
			copy.Name = ""
			return input, copy, validOk("mainnet rpc set to %s", input)
		})
	},
	Replacements: []wizard.Replacement{
		{Color: colors.Green, Values: []string{"\"help\"", "Unchained Index"}},
	},
}

// --------------------------------------------------------
var c2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to index other chains?`,
	Hint: `You may index as many chains as you wish. All you need
|is a separate, fast RPC endpoint for each chain. If
|you do want to index another chain, type "edit" to open
|the file in your editor. Adding your own chains should be
|obvious. Save your work to return to this screen.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		q.Screen.Instructions = `Type "edit" to add another chain or press enter to continue.`
		return input, validContinue()
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		_ = q
		if input != "edit" && len(input) > 0 {
			return "", fmt.Errorf(`"edit" is the only valid response %w`, wizard.ErrValidate)
		}
		return input, validContinue()
	},
}

type ChainQuestion struct {
	wizard.Question
}
