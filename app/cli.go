package app

import (
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
	"github.com/urfave/cli/v2"
)

func initCli(k *KhedraApp) *cli.App {
	os.Args = cleanArgs()

	showError := func(c *cli.Context, showHelp bool, err error) {
		_, _ = c.App.Writer.Write([]byte("\n" + colors.Red + "Error: " + err.Error() + colors.Off + "\n\n"))
		if showHelp {
			cli.ShowAppHelp(c)
		}
	}

	var onUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		showError(c, true, err)
		return nil
	}

	return &cli.App{
		Name:    "khedra",
		Usage:   "A tool to index, monitor, serve, and share blockchain data",
		Version: sdk.Version(),
		Commands: []*cli.Command{
			{
				Name:         "init",
				Usage:        "Initializes Khedra",
				OnUsageError: onUsageError,
				Action: func(c *cli.Context) error {
					if err := validateArgs(1, 1); err != nil {
						return err
					}
					return k.initAction(c)
				},
			},
			{
				Name:         "daemon",
				Usage:        "Runs Khedra's services",
				OnUsageError: onUsageError,
				Action: func(c *cli.Context) error {
					if err := validateArgs(1, 1); err != nil {
						return err
					}
					return k.daemonAction(c)
				},
			},
			{
				Name:         "version",
				Usage:        "Displays Khedra's version",
				Hidden:       true,
				OnUsageError: onUsageError,
				Action: func(c *cli.Context) error {
					if err := validateArgs(0, 0); err != nil {
						return err
					}
					return k.versionAction(c)
				},
			},
			{
				Name:  "config",
				Usage: "Manages Khedra configuration",
				Subcommands: []*cli.Command{
					{
						Name:         "edit",
						Usage:        "Opens the configuration file for editing",
						OnUsageError: onUsageError,
						Action: func(c *cli.Context) error {
							if err := validateArgs(2, 2); err != nil {
								return err
							}
							return k.configEditAction(c)
						},
					},
					{
						Name:         "show",
						Usage:        "Displays the current configuration",
						OnUsageError: onUsageError,
						Action: func(c *cli.Context) error {
							if err := validateArgs(2, 2); err != nil {
								return err
							}
							return k.configShowAction(c)
						},
					},
				},
				OnUsageError: onUsageError,
			},
		},
		OnUsageError: onUsageError,
		CommandNotFound: func(c *cli.Context, command string) {
			var err error
			if unknown := getUnknownCmd(); len(unknown) > 0 {
				err = fmt.Errorf("command '%s' not found", unknown)
			} else {
				err = fmt.Errorf("use only one command at a time")
			}
			showError(c, true, err)
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				showError(c, true, err)
			}
		},
	}
}

/*
AT LEAST ONE SERVICE (OUT OF MONITOR, SCRAPER, API) MUST BE ENABLED
AT LEAST ONE VALID CHAIN WITH ACTIVE RPC MUST BE PROVIDED
A MAINNET RPC MUST BE PROVIDED

	handleService := func(i int, feature Feature) (int, error) {
		if hasValue(i) {
			if mode, err := validateOnOff(os.Args[i+1]); err == nil {
				switch feature {
				case Scrape:
					a.Scrape = mode
					if a.IsOn(Scrape) {
						scrapeSvc := services.NewScrapeService(
							a.Logger,
							string(a.InitMode),
							a.Config.Targets,
							a.Sleep,
							a.BlockCnt,
						)
						activeServices = append(activeServices, scrapeSvc)
					}
				case Api:
					a.Api = mode
					if a.IsOn(Api) {
						apiSvc := services.NewApiService(a.Logger)
						activeServices = append(activeServices, apiSvc)
					}
				case Ipfs:
					a.Ipfs = mode
					if a.IsOn(Ipfs) {
						ipfsSvc := services.NewIpfsService(a.Logger)
						activeServices = append(activeServices, ipfsSvc)
					}
				case Monitor:
					a.Monitor = mode
					if a.IsOn(Monitor) {
						monSvc := services.NewMonitorService(a.Logger)
						activeServices = append(activeServices, monSvc)
					}
				}
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --%s: %w", feature.String(), err)
			}
		}
		return i, fmt.Errorf("%w for --%s", ErrMissingArgument, feature.String())
	}
	controlService := services.NewControlService(a.Logger)
	activeServices = append([]services.Servicer{controlService}, activeServices...)
}
*/
