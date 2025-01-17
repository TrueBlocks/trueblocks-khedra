package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
)

var servicesTitle = `Services`
var servicesSubtitle = `pg 1/1`
var servicesBody = `
This wizard helps you configure Khedra. It will walk you through four
sections General, Services, Chains, and Logging.

You may quit the wizard at any time by typing "q" or "quit". The next
time you run it, it will continue where you left off. Type "help" at any
point to get more information.

Type your answer below a press enter to continue.
`
var servicesOptions = wizard.Option{
	Replacements: []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{"KHEDRA WIZARD"}},
		{Color: colors.Green, Values: []string{"General", "Services", "Chains", "Logging"}},
	},
}
