package app

import (
	"encoding/json"
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
func cPrepare(key, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		if key == "mainnet" {
			if ch, ok := cfg.Chains[key]; !ok {
				ch.Name = key
				ch.Enabled = true
				cfg.Chains[key] = ch
			}

			if ch, ok := cfg.Chains[key]; !ok {
				log.Fatal("chain not found")
			} else {
				if !ch.HasValidRpc() {
					bytes, _ := json.Marshal(&ch)
					q.State = string(bytes)
					msg := fmt.Sprintf("no rpcs for chain %s ", key)
					return strings.Join(ch.RPCs, ","), fmt.Errorf(msg+"%w", wizard.ErrValidate)
				}
			}
		}
	}
	return input, validOk("skip - have all rpcs", input)
}

// --------------------------------------------------------
func cValidate(key string, input string, q *wizard.Question) (string, error) {
	if _, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
			if key == "mainnet" {
				if ch, ok := cfg.Chains[key]; !ok {
					log.Fatal("chain not found")
				} else {
					ch.RPCs = strings.Split(input, ",")
					cfg.Chains[key] = ch
					if !ch.HasValidRpc() {
						bytes, _ := json.Marshal(&ch)
						q.State = string(bytes)
						msg := fmt.Sprintf("no rpcs for chain %s ", key)
						return strings.Join(ch.RPCs, ","), fmt.Errorf(msg+"%w", wizard.ErrValidate)
					}
				}
			}
		}
	}
	return input, nil
}

// --------------------------------------------------------
var c0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}

// --------------------------------------------------------
var c1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Please provide an RPC for Ethereum mainnet?`,
	Hint: `Khedra requires an Ethereum mainnet RPC. It needs to read
|state from the Unchained Index smart contract. Type "help" for
|more information.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		// qq := ChainQuestion{Question: *q}
		return cPrepare("mainnet", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		// qq := ChainQuestion{Question: *q}
		return cValidate("mainnet", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.Green, Values: []string{"\"help\"", "Unchained Index"}},
	},
}

// --------------------------------------------------------
var c2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Which chains do you want to index?`,
	Hint: `Enter a comma separated list of chains to index. The wizard will
|ask you next for RPCs. Enter "chains" to open a large list of
|EVM chains. Use the shortNames from that list to name your
|chains. When you publish your index, others' indexes will
|match (e.g., mainnet, gnosis, optimism, sepolia, etc.)`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		// qq := ChainQuestion{Question: *q}
		return cPrepare("chain", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		// qq := ChainQuestion{Question: *q}
		return cValidate("chain", input, q)
	},
}

// type ChainQuestion struct {
// 	wizard.Question
// }
