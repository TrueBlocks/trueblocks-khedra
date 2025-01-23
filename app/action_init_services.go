package app

import (
	"encoding/json"
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

// screen|---------|---------|---------|---------|---------|---------|---|74
func getServicesScreen(cfg *types.Config) wizard.Screen {
	_ = cfg // linter
	servicesTitle := `Services Settings`
	servicesSubtitle := ``
	servicesInstructions := `Enter "yes" or "no" and press enter. ("e"=edit, "h"=help)`
	servicesBody := `
Khedra provides five services. The first, "control," exposes endpoints to 
control the other four: "scrape", "monitor", "api", and "ipfs".

You may disable/enable any combination of services, but at least one must
be enabled.

The next few screens will allow you to configure each service.
`
	servicesReplacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{"Services Settings"}},
		{Color: colors.BrightBlue, Values: []string{"\"control\"", "\"scrape\"", "\"monitor\"", "\"api\"", "\"ipfs\""}},
	}

	theMap := map[string]*wizard.Question{
		"scraper": &servicesQ1,
		"monitor": &servicesQ2,
		"api":     &servicesQ3,
		"ipfs":    &servicesQ4,
	}
	for key, question := range theMap {
		service := cfg.Services[key]
		bytes, _ := json.Marshal(service)
		question.Value = string(bytes)
	}
	servicesQuestions := []wizard.Question{servicesQ0, servicesQ1, servicesQ2, servicesQ3, servicesQ4}

	return wizard.Screen{
		Title:        servicesTitle,
		Subtitle:     servicesSubtitle,
		Body:         servicesBody,
		Instructions: servicesInstructions,
		Replacements: servicesReplacements,
		Questions:    servicesQuestions,
		Style:        wizard.NewStyle(),
	}
}

var servicesQ0 = wizard.Question{
	Question: ``,
}

func updatePrepare(key, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		service := cfg.Services[key]
		bytes, _ := json.Marshal(service)
		// q.State = colors.Yellow + string(bytes) + colors.Off
		if service.Enabled {
			return "yes", validOk(`question proceeds`, input)
		}
		return "no", validOk(`question proceeds`, input)
	}
	return q.Value, validOk(`question proceeds`, input)
}

func updateValidate(key string, input string, q *wizard.Question) (string, error) {
	if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
		service := cfg.Services[key]
		if input == "yes" {
			service.Enabled = true
			cfg.Services[key] = service
			return input, validOk("the %s sercice was enabled", key)
		} else if input == "no" {
			service.Enabled = false
			cfg.Services[key] = service
			return input, validOk("the %s sercice was disabled", key)
		} else {
			return input, fmt.Errorf(`please enter "yes" or "no"`+"%w", wizard.ErrValidate)
		}
	}
	return input, fmt.Errorf(`could not cast backing data`+"%w", wizard.ErrValidate)
}

var servicesQ1 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "scraper" service?`,
	Hint: `The "scraper" service constanly watches the blockchain and
|updates the Unchained Index with new data. If you disable it,
|your index will fall behind.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return updatePrepare("scraper", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return updateValidate("scraper", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"scraper\""}},
	},
}

var servicesQ2 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "monitor" service?`,
	Hint: `The "monitor" service watches a list of addreses for any
|appearances. Currently disabled, this feature will allow you to
|constantly keep the caches fresh for how ever many addresses you
|like. You may not enable this service.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return updatePrepare("monitor", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return updateValidate("monitor", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"monitor\""}},
	},
}

var servicesQ3 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "api" service?`,
	Hint: `The "api" service serves all of chifra's endpoints as
|described here: https://trueblocks.io/api/.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return updatePrepare("api", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return updateValidate("api", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"api\""}},
	},
}

var servicesQ4 = wizard.Question{
	//.....question-|---------|---------|---------|---------|---------|----|65
	Question: `Do you want to enable the "ipfs" service?`,
	Hint: `The "ipfs" service enables TrueBlocks' pin-by-default mechanism.
|Each time a new index chunk and bloom filter is created, if this
|service is enabled, it will automatically be pinned to IPFS.`,
	PrepareFn: func(input string, q *wizard.Question) (string, error) {
		return updatePrepare("ipfs", input, q)
	},
	Validate: func(input string, q *wizard.Question) (string, error) {
		return updateValidate("ipfs", input, q)
	},
	Replacements: []wizard.Replacement{
		{Color: colors.BrightBlue, Values: []string{"\"ipfs\""}},
	},
}
