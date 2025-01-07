package main

import (
	"github.com/TrueBlocks/trueblocks-khedra/v2/app"
)

func main() {
	k := app.NewKhedraApp()
	if _, proceed, _ := app.LoadConfig(); !proceed {
		return
		// } else if err != nil {
		// 	k.Fatal(err)
	} else {
		// k.Info("Starting Khedra with", "services", len(k.ActiveServices))
		// // TODO: The following should happen in LoadConfig
		// for _, svc := range k.ActiveServices {
		// 	if controlSvc, ok := svc.(*services.ControlService); ok {
		// 		controlSvc.AttachServiceManager(k)
		// 	}
		// }
		// // TODO: The previous should happen in LoadConfig

		// if err := k.StartAllServices(); err != nil {
		// 	a.Fatal(err)
		// }
		k.Run() // HandleSignals()

		select {}
	}
}
