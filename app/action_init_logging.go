package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
)

var loggingTitle = `Logging`
var loggingSubtitle = `pg 1/1`
var loggingBody = `
This wizard helps you configure Khedra. It will walk you through four
sections General, Services, Chains, and Logging.

You may quit the wizard at any time by typing "q" or "quit". The next
time you run it, it will continue where you left off. Type "help" at any
point to get more information.

Type your answer below a press enter to continue.
`
var loggingOptions = wizard.Option{
	Replacements: []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{"KHEDRA WIZARD"}},
		{Color: colors.Green, Values: []string{"General", "Services", "Chains", "Logging"}},
	},
}

// var loggingScreen = `
// 80 ┌──────────────────────────────────────────────────────────────────────────────┐
//    │                                                                              │
//    │                   34 ╔════════════════════════════════╗                      │
//    │                      ║             KHEDRA             ║                      │
//    │                      ║                                ║                      │
//    │                      ║     Index, monitor, serve,     ║                      │
//    │                      ║   and share blockchain data.   ║                      │
//    │                      ╚════════════════════════════════╝                      │
//    │                                                                              │
//    │                                                                              │
//    │   This wizard will help you configure Khedra. The wizard walks you through   │
//    │   four sections General, Services, Chains, and Logging.                      │
//    │                                                                              │
//    │   Type "help" at any point to get more information. You may quit at any      │
//    │   time. The wizard will start over where you left off next time. Type "q"    │
//    │   or "quit" at any time to exit. Press enter to continue.                    │
//    │                                                                              │
//    └──────────────────────────────────────────────────────────────────────────────┘
// `
