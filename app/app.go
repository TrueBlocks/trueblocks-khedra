package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	Cli        *cli.App
	config     *types.Config
	fileLogger *slog.Logger
	progLogger *slog.Logger
}

func NewKhedraApp() *KhedraApp {
	os.Args = cleanArgs(os.Args)
	k := &KhedraApp{}
	k.Cli = initializeCli(k)

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	k.config = &cfg
	k.fileLogger, k.progLogger = types.NewLoggers(cfg.Logging)

	return k
}

func (k *KhedraApp) Run(args []string) error {
	return k.Cli.Run(args)
}

/*
// Feature is a type that represents the features of the app
type Feature string

func (f Feature) String() string {
	return string(f)
}

const (
	// Scrape represents the scraper feature. The scraper may not be disabled.
	Scrape Feature = "scrape"
	// Api represents the API feature. The api is On by default. Disable it
	// with the `--api off` option.
	Api Feature = "api"
	// Ipfs represents the IPFS feature. The ipfs is Off by default. Turn it
	// on with the `--ipfs on` option.
	Ipfs Feature = "ipfs"
	// Monitor represents the monitor feature. The monitor is Off by default. Enable
	// it with the `--monitor on` option.
	Monitor Feature = "monitor"
)

// InitMode is a type that represents the initialization for the Unchained Index. It
// applies to the `--init` option.
type InitMode string

const (
	// All cause the initialization to download both the bloom filters and the index
	// portions of the Unchained Index.
	All InitMode = "all"
	// Blooms cause the initialization to download only the bloom filters portion of
	// the Unchained Index.
	Blooms InitMode = "blooms"
	// None cause the app to not download any part of the Unchained Index. It will be
	// built from scratch with the scraper.
	None InitMode = "none"
)

// OnOff is a type that represents a boolean value that can be either "on" or "off".
type OnOff string

const (
	// On is the "on" value for a feature. It applies to the `--monitor` and `--api` options.
	On OnOff = "on"
	// Off is the "off" value for a feature. It applies to the `--monitor` and `--api` options.
	Off OnOff = "off"
)

// App is the main struct for the app. It contains the logger, the configuration, and the
// state of the app.
type App struct {
	Logger   *slog.Logger
	Config   config.Config
	InitMode InitMode
	Scrape   OnOff
	Api      OnOff
	Ipfs     OnOff
	Monitor  OnOff
	Sleep    int
	BlockCnt int
	Level    slog.Level
}

// NewApp creates a new App instance with the default values.
func NewApp() *App {
	blockCnt := 2000
	if bc, ok := os.LookupEnv("TB_NODE_BLOCKCNT"); ok {
		blockCnt = int(base.MustParseUint64(bc))
	}

	custom Logger, level := NewCustomLogger()
	app := &App{
		Logger:   custom Logger,
		Level:    level,
		Sleep:    6,
		Scrape:   Off,
		Api:      Off,
		Ipfs:     Off,
		Monitor:  Off,
		InitMode: Blooms,
		BlockCnt: blockCnt,
		Config: config.Config{
			ProviderMap: make(map[string]string),
		},
	}

	return app
}

// IsOn returns true if the given feauture is enabled. It returns false otherwise.
func (a *App) IsOn(feature Feature) bool {
	switch feature {
	case Scrape:
		return a.Scrape == On
	case Api:
		return a.Api == On
	case Ipfs:
		return a.Ipfs == On
	case Monitor:
		return a.Monitor == On
	}
	return false
}

// State returns "on" or "off" depending if the feature is on or off.
func (a *App) State(feature Feature) string {
	if a.IsOn(feature) {
		return "on"
	}
	return "off"
}

func (a *App) Fatal(err error) {
	fmt.Printf("Error: %s%s%s\n", colors.Red, err.Error(), colors.Off)
	os.Exit(1)
}
*/

func parseArgsInternal(args []string) (hasHelp bool, hasVersion bool, commands []string, nonFlagCount int) {
	commands = []string{}
	if len(args) == 0 {
		hasHelp = true
		return
	}

	helpForms := map[string]bool{
		"--help": true, "-help": true, "help": true,
		"--h": true, "-h": true,
	}

	versionForms := map[string]bool{
		"--version": true, "-version": true, "version": true,
		"--v": true, "-v": true,
	}

	for i, arg := range args {
		if helpForms[arg] {
			hasHelp = true
			continue
		}
		if versionForms[arg] {
			hasVersion = true
			continue
		}
		commands = append(commands, arg)
		if i != 0 && len(arg) == 0 || arg[0] != '-' {
			nonFlagCount++
		}
	}

	return
}

func cleanArgs(args []string) []string {
	programName := args[:1] // program name

	hasHelp, hasVersion, commands, _ := parseArgsInternal(args[1:])
	if hasHelp {
		result := append(programName, "help")
		if len(commands) > 0 {
			return append(result, commands[0])
		}
		return result
	}

	if hasVersion {
		return append(programName, "version")
	}

	return append(programName, commands...)
}
