package main

import (
	"github.com/TrueBlocks/trueblocks-khedra/v2/app"
)

func main() {
	k := app.NewKhedraApp()
	k.Debug("Starting Khedra")
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
	k.Run() // HandleSignals()

	// 	select {}
	// }
}
