package app

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) initAction(c *cli.Context) error {
	_ = c // liinter
	k.ConfigMaker()
	fmt.Println("Initializing Khedra...")
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
