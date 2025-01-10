package app

import (
	"fmt"
	"os"
	"time"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) initAction(c *cli.Context) error {
	_ = c // liinter
	fmt.Println("Initializing Khedra...")
	// if _, _, flags, cmdCnt := parseArgsInternal(os.Args[1:]); cmdCnt != 1 || len(flags) != 1 {
	// 		return fmt.Errorf("error in initAction: %v %v %d %d", os.Args, flags, cmdCnt, len(flags))
	// 	} else {
	// 		// _ = types.GetConfigFn()
	// 		fmt.Println("Initializing Khedra...", flags, cmdCnt, len(flags))
	// 	}
	return nil
}

func (k *KhedraApp) daemonAction(c *cli.Context) error {
	_ = c // liinter
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

	// if _, _, flags, cmdCnt := parseArgsInternal(os.Args[1:]); cmdCnt != 1 || len(flags) != 1 {
	// 		if cmdCnt > 1 {
	// 			return fmt.Errorf("only one command at a time: %s", os.Args) // error in daemonAction: %v %v %d %d", os.Args, flags, cmdCnt, len(flags))
	// 		}
	// 		return fmt.Errorf("error in daemonAction: %v %v %d %d", os.Args, flags, cmdCnt, len(flags))
	// 	} else {
	// 		fmt.Printf("[flags: %v %d %d]: Sleeping for 10 seconds", flags, cmdCnt, len(flags))
	// 		cnt := 0
	// 		for {
	// 			if cnt >= 1 {
	// 				break
	// 			}
	// 			cnt++
	// 			time.Sleep(time.Second)
	// 			fmt.Printf(".")
	// 		}
	// 		fmt.Println(".")
	// 	}
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

func (k *KhedraApp) versionAction(c *cli.Context) error {
	_ = c // liinter
	fmt.Println("khedra version " + sdk.Version())
	// if _, _, flags, cmdCnt := parseArgsInternal(os.Args[1:]); cmdCnt != 0 || len(flags) != 0 {
	// 	return fmt.Errorf("error in versionAction: %v %v %d %d", os.Args, flags, cmdCnt, len(flags))
	// } else {
	// 	// fmt.Println("In versionAction:", os.Args)
	// 	fmt.Println("khedra version "+sdk.Version(), cmdCnt, len(flags))
	// }
	return nil
}

func (k *KhedraApp) configShowAction(c *cli.Context) error {
	_ = c // liinter
	fmt.Println("In configShowAction:")
	// if _, _, flags, cmdCnt := parseArgsInternal(os.Args[1:]); cmdCnt != 2 || len(flags) != 2 {
	// 	return fmt.Errorf("error in configShowAction: %v %v %d %d", os.Args, flags, cmdCnt, len(flags))
	// } else {
	// 	fmt.Println("In configShowAction:", flags, cmdCnt, len(flags))
	// }
	// cfg, err := LoadConfig()
	// if err != nil {
	// 	return fmt.Errorf("failed to load config: %w", err)
	// }
	// bytes, err := yaml.Marshal(&cfg)
	// if err != nil {
	// 	return fmt.Errorf("failed to unmarshal config: %w", err)
	// }
	// fmt.Println(string(bytes))
	return nil
}

func (k *KhedraApp) configEditAction(c *cli.Context) error {
	_ = c // liinter
	fmt.Println("Would have edited:")
	// if _, _, flags, cmdCnt := parseArgsInternal(os.Args[1:]); cmdCnt != 2 || len(flags) != 2 {
	// 	return fmt.Errorf("error in configEditAction: %v %v %d %d", os.Args, flags, cmdCnt, len(flags))
	// } else {
	// 	fmt.Println("Would have edited", flags, cmdCnt, len(flags))
	// }
	// editor := os.Getenv("EDITOR")
	// if editor == "" {
	// 	editor = "nano"
	// }
	// configPath := types.GetConfigFn()
	// cmd := execCommand(editor, configPath)
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// if err := cmd.Run(); err != nil {
	// 	return fmt.Errorf("failed to open config for editing: %w", err)
	// }
	return nil
}
