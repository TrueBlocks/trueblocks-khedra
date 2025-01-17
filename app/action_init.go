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
		{Name: "Welcome", Prompt: "Welcome to Khedra. We're going to walk you through.", Metadata: map[string]interface{}{"default": "yes"}},
		{Name: "General", Prompt: "Where do you want to store your data?", Metadata: map[string]interface{}{"default": "25"}},
		{Name: "Services", Prompt: "There are four services. Which do you want to enable?", Metadata: map[string]interface{}{"default": "25"}},
		{Name: "Chains", Prompt: "Which chains do you want to support?", Metadata: map[string]interface{}{"default": "New York"}},
		{Name: "Logging", Prompt: "What logging options do you want to enable?", Metadata: map[string]interface{}{"default": "New York"}},
	}

	w := wizard.NewWizard(steps, ">")

	if err := w.Run(); err != nil {
		return err
	}

	return nil
}

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
