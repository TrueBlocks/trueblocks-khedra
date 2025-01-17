package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-khedra/v2/app/wizard"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) initAction(c *cli.Context) error {
	_ = c // linter
	if _, err := k.ConfigMaker(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	steps := []wizard.Step{
		wizard.NewScreen(welcomeTitle, welcomeSubtitle, welcomeBody, welcomeOptions),
		wizard.NewScreen(generalTitle, generalSubtitle, generalBody, generalOptions),
		wizard.NewScreen(servicesTitle, servicesSubtitle, servicesBody, servicesOptions),
		wizard.NewScreen(chainsTitle, chainsSubtitle, chainsBody, chainsOptions),
		wizard.NewScreen(loggingTitle, loggingSubtitle, loggingBody, loggingOptions),
	}

	w := wizard.NewWizard(steps, "")

	if err := w.Run(); err != nil {
		return err
	}

	return nil
}

var generalScreen = `Where should Khedra store its data?`

var servicesScreen = `There are four services. Which do you want to enable?`

var chainsScreen = `Which chains do you want to support?`

var loggingScreen = `Where should we store log files?`

// func (k *KhedraApp) Initialize() error {
// 	k.Info("Test log: Scraper initialization started")
// 	k.Info("Initializing unchained index")

// 	if s.initMode != "none" {
// 		reports := make([]*scraperReport, 0, len(s.configTargets))
// 		for _, chain := range s.configTargets {
// 			if rep, err := s.initOneChain(chain); err != nil {
// 				if !strings.HasPrefix(err.Error(), "no record found in the Unchained Index") {
// 					k.Warn("Warning", "msg", err)
// 				} else {
// 					k.Warn("No record found in the Unchained Index for chain", "chain", chain)
// 				}
// 			} else {
// 				reports = append(reports, rep)
// 			}
// 		}

// 		for _, report := range reports {
// 			reportScrape(k, report)
// 		}
// 	}
// 	k.Info("Scraper initialization complete")
// 	return nil
// }
