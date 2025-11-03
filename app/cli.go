package app

import (
	"fmt"
	"os"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/colors"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v6"
	"github.com/urfave/cli/v2"
)

func initCli(k *KhedraApp) *cli.App {
	os.Args = cleanArgs()

	showError := func(c *cli.Context, showHelp bool, err error) {
		_, _ = c.App.Writer.Write([]byte("\n" + colors.Red + "Error: " + err.Error() + colors.Off + "\n\n"))
		if showHelp {
			_ = cli.ShowAppHelp(c)
		}
	}

	var onUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		_ = isSubcommand
		showError(c, true, err)
		return nil
	}

	return &cli.App{
		Name:    "khedra",
		Usage:   "A tool to index, monitor, serve, and share blockchain data",
		Version: sdk.Version(),
		Commands: []*cli.Command{
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
			{
				Name:         "pause",
				Usage:        "Pause the given service (one of scraper, monitor, all)",
				OnUsageError: onUsageError,
				Action: func(c *cli.Context) error {
					return k.pauseAction(c)
				},
			},
			{
				Name:         "unpause",
				Usage:        "Unpause the given service (one of scraper, monitor, all)",
				OnUsageError: onUsageError,
				Action: func(c *cli.Context) error {
					return k.unpauseAction(c)
				},
			},
		},
		OnUsageError: onUsageError,
		CommandNotFound: func(c *cli.Context, command string) {
			_ = command
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
