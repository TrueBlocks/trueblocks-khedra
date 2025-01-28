package app

import (
	"encoding/json"
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/wizard"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getServicesScreen() wizard.Screen {
	sTitle := `Services Settings`
	sSubtitle := ``
	sInstructions := `Enter "yes" or "no" and press enter. ("e"=edit, "h"=help)`
	sBody := `
Khedra provides five services. The first, "control," exposes endpoints to 
control the other four: "scrape", "monitor", "api", and "ipfs".

You may disable/enable any combination of services, but at least one must
be enabled.

The next few screens will allow you to configure each service.
`
	sReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{sTitle}},
		{Color: colors.BrightBlue, Values: []string{
			"\"control\"", "\"scrape\"", "\"monitor\"", "\"api\"", "\"ipfs\"",
		}},
	}
	sQuestions := []wizard.Questioner{&s0, &s1, &s2, &s3, &s4}
	sStyle := wizard.NewStyle()

	return wizard.Screen{
		Title:        sTitle,
		Subtitle:     sSubtitle,
		Body:         sBody,
		Instructions: sInstructions,
		Replacements: sReplacements,
		Questions:    sQuestions,
		Style:        sStyle,
	}
}

// --------------------------------------------------------
func sPrepare(key, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		service := cfg.Services[key]
		bytes, _ := json.Marshal(service)
		q.State = string(bytes)
		if service.Enabled {
			q.Value = "yes"
			return "yes", validOk(`don't skip`, input)
		}
		q.Value = "no"
		return "no", validOk(`don't skip`, input)
	}
	return input, validOk(`don't skip`, input)
}

// --------------------------------------------------------
func sValidate(key string, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		service := cfg.Services[key]
		if input == "yes" {
			service.Enabled = true
			cfg.Services[key] = service
			err := cfg.WriteToFile(types.GetConfigFnNoCreate())
			if err != nil {
				fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
			}
			return input, validOk(`the %s service was enabled`, key)
		} else if input == "no" {
			service.Enabled = false
			cfg.Services[key] = service
			err := cfg.WriteToFile(types.GetConfigFnNoCreate())
			if err != nil {
				fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
			}
			return input, validOk(`the %s service was disabled`, key)
		} else {
			return input, fmt.Errorf(`value must be either "yes" or "no" %w`, wizard.ErrValidate)
		}
	}
	return input, fmt.Errorf(`could not cast backing data`+"%w", wizard.ErrValidate)
}

// --------------------------------------------------------
var s0 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
}

// --------------------------------------------------------
var s1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "scraper" service?`,
	Hint: `The "scraper" service constanly watches the blockchain and
|updates the Unchained Index with new data. If you disable it,
|your index will fall behind.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return sPrepare("scraper", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return sValidate("scraper", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"scraper\""}},
	},
}

// --------------------------------------------------------
var s2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "monitor" service?`,
	Hint: `The "monitor" service watches a list of addreses for any
|appearances. Currently disabled, this feature will allow you to
|constantly keep the caches fresh for how ever many addresses you
|like. You may not enable this service.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return sPrepare("monitor", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return sValidate("monitor", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"monitor\""}},
	},
}

// --------------------------------------------------------
var s3 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "api" service?`,
	Hint: `The "api" service serves all of chifra's endpoints as
|described here: https://trueblocks.io/api/.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return sPrepare("api", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return sValidate("api", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"api\""}},
	},
}

// --------------------------------------------------------
var s4 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "ipfs" service?`,
	Hint: `The "ipfs" service enables TrueBlocks' pin-by-default mechanism.
|Each time a new index chunk and bloom filter is created, if this
|service is enabled, it will automatically be pinned to IPFS.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return sPrepare("ipfs", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return sValidate("ipfs", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"ipfs\""}},
	},
}
