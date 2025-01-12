package app

import (
	"fmt"
	"os"
	"time"

	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) daemonAction(c *cli.Context) error {
	_ = c // liinter
	k.ConfigMaker()
	fmt.Printf("Sleeping for 10 seconds")
	cnt := 0
	for {
		if cnt >= 10 {
			break
		}
		cnt++
		if os.Getenv("TEST_MODE") != "true" {
			time.Sleep(time.Second)
		}
		fmt.Printf(".")
	}
	fmt.Println(".")

	// if _, proceed, err := app.LoadConfig(); !proceed {
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
