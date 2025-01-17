package app

import (
	"fmt"
	"os"
	"time"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) daemonAction(c *cli.Context) error {
	_ = c // linter
	fn := types.GetConfigFnNoCreate()
	if !coreFile.FileExists(fn) {
		return fmt.Errorf("not initialized you must run `khedra init` first")
	}

	_, err := k.ConfigMaker()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	for _, chain := range k.config.Chains {
		for _, rpc := range chain.RPCs {
			if err := validate.TryConnect(chain.Name, rpc, 5); err != nil {
				return err
			}
		}
	}
	fmt.Printf("Sleeping for 10 seconds")
	cnt := 0
	for {
		if cnt >= 1 {
			break
		}
		cnt++
		if os.Getenv("TEST_MODE") != "true" {
			time.Sleep(time.Second)
		}
		fmt.Printf(".")
	}
	fmt.Println(".")

	// if _, proceed, err := app.Load Config(); !proceed {
	// 	return
	// } else if err != nil {
	// 	k.Fatal(err.Error())
	// } else {
	// k.Info("Starting Khedra with", "services", len(k.ActiveServices))
	// // TODO: The following should happen in Load Config
	// for _, svc := range k.ActiveServices {
	// 	if controlSvc, ok := svc.(*services.ControlService); ok {
	// 		controlSvc.AttachServiceManager(k)
	// 	}
	// }
	// // TODO: The previous should happen in Load Config
	// if err := k.StartAllServices(); err != nil {
	// 	a.Fatal(err)
	// }
	// HandleSignals()

	// 	select {}
	// }
	return nil
}
