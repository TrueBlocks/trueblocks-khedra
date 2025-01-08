package app

import (
	"log/slog"

	"github.com/urfave/cli/v2"
)

func initializeCli() *cli.App {
	return &cli.App{
		Name:  "khedra",
		Usage: "A tool to index, monitor, serve, and share blockchain data",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initializes Khedra",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "mode",
						Usage: "Initialization mode (all, blooms, none)",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: init", "mode", c.String("mode"))
					return nil
				},
			},
			{
				Name:  "scrape",
				Usage: "Controls the blockchain scraper",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "state",
						Usage: "Scraper state (on, off)",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: scrape", "state", c.String("state"))
					return nil
				},
			},
			{
				Name:  "api",
				Usage: "Starts or stops the API server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "state",
						Usage: "API state (on, off)",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: api", "state", c.String("state"))
					return nil
				},
			},
			{
				Name:  "sleep",
				Usage: "Sets the duration between updates",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "duration",
						Usage: "Sleep duration in seconds",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: sleep", "duration", c.Int("duration"))
					return nil
				},
			},
		},
	}
}

/*
var (
	ErrMissingArgument = errors.New("missing argument")
	ErrInvalidValue    = errors.New("invalid value")
)

func (a *App) ParseArgs() (bool, []services.Servicer, error) {
	var activeServices []services.Servicer

	hasValue := func(i int) bool {
		return i+1 < len(os.Args) && os.Args[i+1][0] != '-'
	}

	handleInit := func(i int) (int, error) {
		if hasValue(i) {
			if mode, err := validateMode(os.Args[i+1]); err == nil {
				a.InitMode = mode
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --init: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --init", ErrMissingArgument)
	}

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

	handleSleep := func(i int) (int, error) {
		if hasValue(i) {
			if sleep, err := validateSleep(os.Args[i+1]); err == nil {
				a.Sleep = sleep
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --sleep: %w", err)
			}
		}
		return i, fmt.Errorf("%w for --sleep", ErrMissingArgument)
	}

	a.Logger.Debug("Parsing command line arguments", "args", os.Args)

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		var err error
		switch arg {
		case "--scrape":
			if i, err = handleService(i, Scrape); err != nil {
				return true, nil, err
			}
		case "--api":
			if i, err = handleService(i, Api); err != nil {
				return true, nil, err
			}
		case "--ipfs":
			if i, err = handleService(i, Ipfs); err != nil {
				return true, nil, err
			}
		case "--monitor":
			if i, err = handleService(i, Monitor); err != nil {
				return true, nil, err
			}
		case "--init":
			i, err = handleInit(i)
		case "--sleep":
			i, err = handleSleep(i)
		case "--version":
			fmt.Println("trueblocks-node " + sdk.Version())
			return false, nil, nil
		default:
			if arg != "--help" {
				return true, nil, fmt.Errorf("unknown option:%s\n%s", os.Args[i], helpText)
			}
			fmt.Printf("%s\n", helpText)
			return false, nil, nil
		}
		if err != nil {
			return true, nil, err
		}
	}

	if len(activeServices) == 0 && os.Getenv("TEST_MODE") != "true" {
		return true, nil, fmt.Errorf("you must enable at least one of the services\n%s", helpText)
	}

	controlService := services.NewControlService(a.Logger)
	activeServices = append([]services.Servicer{controlService}, activeServices...)

	a.Logger.Debug("Command line parsing complete", "services", len(activeServices))
	return true, activeServices, nil
}

func validateEnum[T ~string](value T, validOptions []T, name string) (T, error) {
	for _, option := range validOptions {
		if value == option {
			return value, nil
		}
	}
	return value, fmt.Errorf("invalid value for %s: %s", name, value)
}

func validateMode(value string) (InitMode, error) {
	return validateEnum(InitMode(value), []InitMode{All, Blooms, None}, "mode")
}

func validateOnOff(value string) (OnOff, error) {
	return validateEnum(OnOff(value), []OnOff{On, Off}, "onOff")
}

func validateSleep(value string) (int, error) {
	var sleep int
	if _, err := fmt.Sscanf(value, "%d", &sleep); err != nil || sleep < 1 {
		return 1, fmt.Errorf("invalid value for sleep: %s", value)
	}
	return sleep, nil
}

const helpText = `Usage: trueblocks-node <options>

Options:
---------
 --init     [all|blooms|none*]   download from the unchained index smart contract (default: none)
 --scrape   [on|off*]            enable/disable the Unchained Index scraper (default: off)
 --api      [on|off*]            enable/disable API server (default: off)
 --ipfs     [on|off*]            enable/disable IPFS daemon (default: off)
 --monitor  [on|off*]            enable/disable address monitoring (currently disabled, default: off)
 --sleep    int                  the number of seconds to sleep between updates (default: 30)
 --version                       display the version string
 --help                          display this help text

Notes:
-------
If --scrape is on, --init must be either blooms or all. If you choose --all, you must always choose --all.

Environment:
-------------
You MUST export the following values to the environment:

  TB_NODE_DATAFOLDER:    A directory to store the indexer's data (required, created if necessary)
  TB_NODE_MAINNETRPC: A valid RPC endpoint for Ethereum mainnet (required)

You MAY also export these environment variables:

  TB_NODE_CHAINS:     A comma-separated list of chains to index (default: "mainnet")
  TB_NODE_<CHAIN>RPC: For each CHAIN in the TB_NODE_CHAINS list, a valid RPC endpoint
                      (example: TB_NODE_SEPOLIARPC=http://localhost:8548)

You may put these values in a .env file in the current folder. See env.example.`
*/
