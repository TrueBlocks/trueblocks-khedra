package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/control"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/install"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) daemonAction(c *cli.Context) error {
	defer func() {
		// cleanup the control file on exit
		_ = os.Remove(control.Path())
	}()

	_ = c // linter
	// Ensure logger is initialized for first-run case
	if k.logger == nil {
		k.logger = types.NewLogger(types.Logging{Level: "info"})
	}

	// On first run, start control service immediately for wizard access
	if !install.Configured() {
		if err := k.initializeControlSvc(); err != nil {
			return err
		}
		if os.Getenv("KHEDRA_EMBED") != "1" {
			fmt.Printf("Khedra is not configured. Please complete setup in your browser: http://localhost:%d\n", k.controlSvc.Port())
			utils.System("open http://localhost:" + fmt.Sprintf("%d", k.controlSvc.Port()))
		}
		// Let daemon continue to blocking loop - live-reload will handle config changes
	}

	// Only load config if we're now configured (either was already, or just completed setup)
	if install.Configured() {
		if err := k.loadConfigIfInitialized(); err != nil {
			return err
		}
		k.logger.Info("Starting khedra daemon...config loaded...")
	}

	if err := k.handleWaitForNode(); err != nil {
		return err
	}

	// Only validate RPCs if we have a complete, finalized configuration
	// Skip validation if wizard is still in progress (draft exists) or config is incomplete
	draft, _ := install.LoadDraft()
	if install.Configured() && draft == nil {
		for _, ch := range k.config.Chains {
			if ch.Enabled {
				if !HasValidRpc(&ch, 4) {
					return fmt.Errorf("chain %s has no valid RPC", ch.Name)
				}
				k.logger.Progress("Connected to", "chain", ch.Name)
			}
		}
	}
	k.logger.Info("Processing chains...", "chainList", k.config.EnabledChains())
	k.logger.Info("Paths:", "indexPath", k.config.IndexPath())
	k.logger.Info("", "cachePath", k.config.CachePath())

	rootFolder := config.PathToRootConfig()

	os.Setenv("XDG_CONFIG_HOME", rootFolder)
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_SETTINGS_INDEXPATH", k.config.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", k.config.CachePath())
	for key, ch := range k.config.Chains {
		if ch.Enabled {
			envKey := "TB_CHAINS_" + strings.ToUpper(key) + "_RPCPROVIDER"
			os.Setenv(envKey, ch.RPCs[0])
		}
	}

	for _, env := range os.Environ() {
		if (strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_")) && strings.Contains(env, "=") {
			parts := strings.Split(env, "=")
			if len(parts) > 1 {
				k.logger.Progress("environment", parts[0], parts[1])
			} else {
				k.logger.Progress("environment", parts[0], "<empty>")
			}
		}
	}

	k.logger.Progress("Starting services", "services", k.config.ServiceList(true /* enabledOnly */))

	configFn := filepath.Join(rootFolder, "trueBlocks.toml")
	if file.FileExists(configFn) {
		k.logger.Info("Config file found", "fn", configFn)
		if !k.chainsConfigured(configFn) {
			k.logger.Error("Config file not configured", "fn", configFn)
			return fmt.Errorf("config file not configured")
		}
	} else {
		k.logger.Warn("Config file not found", "fn", configFn)
		if err := k.createChifraConfig(rootFolder); err != nil {
			k.logger.Error("Error creating config file", "error", err)
			return err
		}
	}

	// Initialize the control service -- we need it for daemon
	_ = k.initializeControlSvc()
	if err := k.serviceManager.StartAllServices(); err != nil {
		k.logger.Panic("%s", err.Error())
	}

	// Delegate signal handling & graceful cleanup to the ServiceManager implementation.
	k.serviceManager.HandleSignals()
	k.logger.Info("daemon running; press Ctrl+C to initiate graceful shutdown (managed by ServiceManager)")
	select {} // block until ServiceManager handles signal and exits process
}

func (k *KhedraApp) chainsConfigured(configFn string) bool {
	chainStr := k.config.EnabledChains()
	chains := strings.Split(chainStr, ",")

	k.logger.Info("chifra config loaded")
	k.logger.Info("checking", "configFile", configFn, "nChains", len(chains))

	contents := file.AsciiFileToString(configFn)
	for _, chain := range chains {
		search := "[chains." + chain + "]"
		if !strings.Contains(contents, search) {
			msg := fmt.Sprintf("config file {%s} does not contain {%s}", configFn, search)
			k.logger.Error(msg)
			return false
		}
	}
	return true
}

func (k *KhedraApp) createChifraConfig(rootFolder string) error {
	if err := file.EstablishFolder(rootFolder); err != nil {
		return err
	}

	chainStr := k.config.EnabledChains()
	chains := strings.Split(chainStr, ",")
	for _, chain := range chains {
		if err := k.createChainConfigFolder(rootFolder, chain); err != nil {
			return err
		}
	}

	tmpl, err := template.New("tmpl").Parse(configTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, &k.config); err != nil {
		return err
	}
	if len(buf.String()) == 0 {
		return fmt.Errorf("empty config file")
	}

	configFn := filepath.Join(rootFolder, "trueBlocks.toml")
	err = file.StringToAsciiFile(configFn, buf.String())
	if err != nil {
		return err
	}
	k.logger.Info("Created config file", "configFile", configFn, "nChains", len(chains))
	return nil
}

// For monitor --watch
// 14080,apps,Accounts,monitors,acctExport,watch,w,,visible|docs|notApi,4,switch,<boolean>,,,,,continually scan for new blocks and extract data as per the command file
// 14090,apps,Accounts,monitors,acctExport,watchlist,a,,visible|docs|notApi,,flag,<string>,,,,,available with --watch option only&#44; a file containing the addresses to watch
// 14100,apps,Accounts,monitors,acctExport,commands,d,,visible|docs|notApi,,flag,<string>,,,,,available with --watch option only&#44; the file containing the list of commands to apply to each watched address
// 14110,apps,Accounts,monitors,acctExport,batch_size,b,8,visible|docs|notApi,,flag,<uint64>,,,,,available with --watch option only&#44; the number of monitors to process in each batch
// 14120,apps,Accounts,monitors,acctExport,run_count,u,,visible|docs|notApi,,flag,<uint64>,,,,,available with --watch option only&#44; run the monitor this many times&#44; then quit
// 14130,apps,Accounts,monitors,acctExport,sleep,s,14,visible|docs|notApi,,flag,<float64>,,,,,available with --watch option only&#44; the number of seconds to sleep between runs
// 14160,apps,Accounts,monitors,acctExport,n3,,,,,note,,,,,,The --watch option requires two additional parameters to be specified: `--watchlist` and `--commands`.
// 14170,apps,Accounts,monitors,acctExport,n4,,,,,note,,,,,,Addresses provided on the command line are ignored in `--watch` mode.
// 14180,apps,Accounts,monitors,acctExport,n5,,,,,note,,,,,,Providing the value `existing` to the `--watchlist` monitors all existing monitor files (see --list).

func (k *KhedraApp) createChainConfigFolder(rootFolder string, chain string) error {
	chainConfig := filepath.Join(rootFolder, "config", chain)
	if err := file.EstablishFolder(chainConfig); err != nil {
		return fmt.Errorf("failed to create folder %s: %w", chainConfig, err)
	}

	baseURL := "https://raw.githubusercontent.com/TrueBlocks/trueblocks-core/refs/heads/master/src/other/install/per-chain"
	url, err := url.JoinPath(baseURL, chain, "allocs.csv")
	if err != nil {
		return err
	}
	allocFn := filepath.Join(chainConfig, "allocs.csv")
	dur := 100 * 365 * 24 * time.Hour // 100 years
	if _, err := downloadAndStore(url, allocFn, dur); err != nil {
		k.logger.Warn(fmt.Errorf("failed to download and store allocs.csv for chain %s: %w", chain, err).Error())
		// It's not an error to not have an allocation file. IsArchiveNode assumes archive if not present.
		return nil
	}
	k.logger.Progress("Creating chain config", "chainConfig", allocFn)
	k.logger.Progress("Creating chain config", "source", url)

	return nil
}

var configTmpl string = `[version]
  current = "{{.Version}}"

[settings]
  cachePath = "{{.CachePath}}"
  defaultChain = "mainnet"
  indexPath = "{{.IndexPath}}"

[keys]
  [keys.etherscan]
    apiKey = ""

[chains]
{{- range .Chains}}
  [chains.{{.Name}}]
    chain = "{{.Name}}"
    chainId = "{{.ChainID}}"
    remoteExplorer = "{{.RemoteExplorer}}"
    rpcProvider = "{{ index .RPCs 0 }}"
    symbol = "{{.Symbol}}"
{{end -}}
`

/*

// HandleWatch starts the monitor watcher
func (opts *MonitorsOptions) HandleWatch(rCtx *output.RenderCtx) error {
	opts.Globals.Cache = true
	scraper := NewScraper(colors.Magenta, "MonitorScraper", opts.Sleep, 0)

	var wg sync.WaitGroup
	wg.Add(1)
	// Note that this never returns in normal operation
	go opts.RunMonitorScraper(&wg, &scraper)
	wg.Wait()

	return nil
}

// RunMonitorScraper runs continually, never stopping and freshens any existing monitors
func (opts *MonitorsOptions) RunMonitorScraper(wg *sync.WaitGroup, s *Scraper) {
	defer wg.Done()

	chain := opts.Globals.Chain
	tmpPath := filepath.Join(config.PathToCache(chain), "tmp")

	s.ChangeState(true, tmpPath)

	runCount := uint64(0)
	for {
		if !s.Running {
			s.Pause()

		} else {
			monitorList := opts.getMonitorList()
			if len(monitorList) == 0 {
				logger.Error(types.Usage("No monitors found. Use 'chifra list' to initialize a monitor.").Error())
				return
			}

			if canceled, err := opts.Refresh(monitorList); err != nil {
				logger.Error(err)
				return
			} else {
				if canceled {
					return
				}
			}

			runCount++
			if opts.RunCount != 0 && runCount >= opts.RunCount {
				return
			}

			sleep := opts.Sleep
			if sleep > 0 {
				ms := time.Duration(sleep*1000) * time.Millisecond
				if !opts.Globals.TestMode {
					logger.Info(fmt.Sprintf("Sleeping for %g seconds", sleep))
				}
				time.Sleep(ms)
			}
		}
	}
}

type Command struct {
	Fmt    string `json:"fmt"`
	Folder string `json:"folder"`
	Cmd    string `json:"cmd"`
	Cache  bool   `json:"cache"`
}

func (c *Command) fileName(addr base.Address) string {
	return filepath.Join(c.Folder, addr.Hex()+"."+c.Fmt)
}

func (c *Command) resolve(addr base.Address, before, after int64) string {
	fn := c.fileName(addr)
	if file.FileExists(fn) {
		if strings.Contains(c.Cmd, "export") {
			c.Cmd += fmt.Sprintf(" --first_record %d", uint64(before+1))
			c.Cmd += fmt.Sprintf(" --max_records %d", uint64(after-before+1)) // extra space won't hurt
		} else {
			c.Cmd += fmt.Sprintf(" %d-%d", before+1, after)
		}
		c.Cmd += " --append --no_header"
	}
	c.Cmd = strings.ReplaceAll(c.Cmd, "  ", " ")
	ret := c.Cmd + " --fmt " + c.Fmt + " --output " + c.fileName(addr) + " " + addr.Hex()
	if c.Cache {
		ret += " --cache"
	}
	return ret
}

func (c *Command) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

func (opts *MonitorsOptions) Refresh(monitors []monitor.Monitor) (bool, error) {
	theCmds, err := opts.getCommands()
	if err != nil {
		return false, err
	}

	batches := batchSlice[monitor.Monitor](monitors, opts.BatchSize)
	for i := 0; i < len(batches); i++ {
		addrs := []base.Address{}
		countsBefore := []int64{}
		for _, mon := range batches[i] {
			addrs = append(addrs, mon.Address)
			countsBefore = append(countsBefore, mon.Count())
		}

		batchSize := int(opts.BatchSize)
		fmt.Printf("%s%d-%d of %d:%s chifra export --freshen",
			colors.BrightBlue,
			i*batchSize,
			min(((i+1)*batchSize)-1, len(monitors)),
			len(monitors),
			colors.Green)
		for _, addr := range addrs {
			fmt.Printf(" %s", addr.Hex())
		}
		fmt.Println(colors.Off)

		canceled, err := opts.FreshenMonitorsForWatch(addrs)
		if canceled || err != nil {
			return canceled, err
		}

		for j := 0; j < len(batches[i]); j++ {
			mon := batches[i][j]
			countAfter := mon.Count()

			if countAfter > 1000000 {
				// TODO: Make this value configurable
				fmt.Println(colors.Red, "Too many transactions for address", mon.Address, colors.Off)
				continue
			}

			if countAfter == 0 {
				continue
			}

			logger.Info(fmt.Sprintf("Processing item %d in batch %d: %d %d\n", j, i, countsBefore[j], countAfter))

			for _, cmd := range theCmds {
				countBefore := countsBefore[j]
				if countBefore == 0 || countAfter > countBefore {
					utils.System(cmd.resolve(mon.Address, countBefore, countAfter))
					// o := opts
					// o.Globals.File = ""
					// _ = o.Globals.PassItOn("acctExport", chain, cmd, []string{})
				} else if opts.Globals.Verbose {
					fmt.Println("No new transactions for", mon.Address.Hex(), "since last run.")
				}
			}
		}
	}
	return false, nil
}

func batchSlice[T any](slice []T, batchSize uint64) [][]T {
	var batches [][]T
	for i := 0; i < len(slice); i += int(batchSize) {
		end := i + int(batchSize)
		if end > len(slice) {
			end = len(slice)
		}
		batches = append(batches, slice[i:end])
	}
	return batches
}

func GetExportFormat(cmd, def string) string {
	if strings.Contains(cmd, "json") {
		return "json"
	} else if strings.Contains(cmd, "txt") {
		return "txt"
	} else if strings.Contains(cmd, "csv") {
		return "csv"
	}
	if len(def) > 0 {
		return def
	}
	return "csv"
}

func (opts *MonitorsOptions) cleanLine(lineIn string) (cmd Command, err error) {
	line := strings.ReplaceAll(lineIn, "[{ADDRESS}]", "")
	if strings.Contains(line, "--fmt") {
		line = strings.ReplaceAll(line, "--fmt", "")
		line = strings.ReplaceAll(line, "json", "")
		line = strings.ReplaceAll(line, "csv", "")
		line = strings.ReplaceAll(line, "txt", "")
	}
	line = utils.StripComments(line)
	if len(line) == 0 {
		return Command{}, nil
	}

	folder, err := opts.getOutputFolder(line)
	if err != nil {
		return Command{}, err
	}

	_ = file.EstablishFolder(folder)
	return Command{Cmd: line, Folder: folder, Fmt: GetExportFormat(lineIn, "csv"), Cache: opts.Globals.Cache}, nil
}

func (opts *MonitorsOptions) getCommands() (ret []Command, err error) {
	lines := file.AsciiFileToLines(opts.Commands)
	for _, line := range lines {
		// orig := line
		if cmd, err := opts.cleanLine(line); err != nil {
			return nil, err
		} else if len(cmd.Cmd) == 0 {
			continue
		} else {
			ret = append(ret, cmd)
		}
	}
	return ret, nil
}

func (opts *MonitorsOptions) getOutputFolder(orig string) (string, error) {
	okMap := map[string]bool{
		"export": true,
		"list":   true,
		"state":  true,
		"tokens": true,
	}

	cmdLine := orig
	parts := strings.Split(strings.ReplaceAll(cmdLine, "  ", " "), " ")
	if len(parts) < 1 || parts[0] != "chifra" {
		s := fmt.Sprintf("Invalid command: %s. Must start with 'chifra'.", strings.Trim(orig, " \t\n\r"))
		logger.Panic(s)
	}
	if len(parts) < 2 || !okMap[parts[1]] {
		s := fmt.Sprintf("Invalid command: %s. Must start with 'chifra export', 'chifra list', 'chifra state', or 'chifra tokens'.", orig)
		logger.Panic(s)
	}

	cwd, _ := os.Getwd()
	cmdLine += " "
	folder := "unknown"
	if parts[1] == "export" {
		if strings.Contains(cmdLine, "-p ") || strings.Contains(cmdLine, "--appearances ") {
			folder = filepath.Join(cwd, parts[1], "appearances")
		} else if strings.Contains(cmdLine, "-r ") || strings.Contains(cmdLine, "--receipts ") {
			folder = filepath.Join(cwd, parts[1], "receipts")
		} else if strings.Contains(cmdLine, "-l ") || strings.Contains(cmdLine, "--logs ") {
			folder = filepath.Join(cwd, parts[1], "logs")
		} else if strings.Contains(cmdLine, "-t ") || strings.Contains(cmdLine, "--traces ") {
			folder = filepath.Join(cwd, parts[1], "traces")
		} else if strings.Contains(cmdLine, "-n ") || strings.Contains(cmdLine, "--neighbors ") {
			folder = filepath.Join(cwd, parts[1], "neighbors")
		} else if strings.Contains(cmdLine, "-C ") || strings.Contains(cmdLine, "--accounting ") {
			folder = filepath.Join(cwd, parts[1], "accounting")
		} else if strings.Contains(cmdLine, "-A ") || strings.Contains(cmdLine, "--statements ") {
			folder = filepath.Join(cwd, parts[1], "statements")
		} else if strings.Contains(cmdLine, "-b ") || strings.Contains(cmdLine, "--balances ") {
			folder = filepath.Join(cwd, parts[1], "balances")
		} else {
			folder = filepath.Join(cwd, parts[1], "transactions")
		}

	} else if parts[1] == "list" {
		folder = filepath.Join(cwd, parts[1], "appearances")

	} else if parts[1] == "state" {
		if strings.Contains(cmdLine, "-l ") || strings.Contains(cmdLine, "--call ") {
			folder = filepath.Join(cwd, parts[1], "calls")
		} else {
			folder = filepath.Join(cwd, parts[1], "blocks")
		}

	} else if parts[1] == "tokens" {
		if strings.Contains(cmdLine, "-b ") || strings.Contains(cmdLine, "--by_acct ") {
			folder = filepath.Join(cwd, parts[1], "by_acct")
		} else {
			folder = filepath.Join(cwd, parts[1], "blocks")
		}
	}

	if strings.Contains(folder, "unknown") {
		return "", fmt.Errorf("unable to determine output folder for command: %s", cmdLine)
	}

	if abs, err := filepath.Abs(filepath.Join(opts.Globals.Chain, folder)); err != nil {
		return "", err
	} else {
		return abs, nil
	}
}

func (opts *MonitorsOptions) getMonitorList() []monitor.Monitor {
	var monitors []monitor.Monitor

	monitorChan := make(chan monitor.Monitor)
	go monitor.ListWatchedMonitors(opts.Globals.Chain, opts.Watchlist, monitorChan)

	for result := range monitorChan {
		switch result.Address {
		case base.NotAMonitor:
			logger.Info(fmt.Sprintf("Loaded %d monitors", len(monitors)))
			close(monitorChan)
		default:
			if result.Count() > 500000 {
				logger.Warn("Ignoring too-large address", result.Address)
				continue
			}
			monitors = append(monitors, result)
		}
	}

	return monitors
}

			if opts.Watch {
				if opts.Globals.IsApiMode() {
					return types.Usage("The {0} options is not available from the API", "--watch")
				}

				if len(opts.Globals.File) > 0 {
					return types.Usage("The {0} option is not allowed with the {1} option. Use {2} instead.", "--file", "--watch", "--commands")
				}

				if len(opts.Commands) == 0 {
					return types.Usage("The {0} option requires {1}.", "--watch", "a --commands file")
				} else {
					cmdFile, err := filepath.Abs(opts.Commands)
					if err != nil || !file.FileExists(cmdFile) {
						return types.Usage("The {0} option requires {1} to exist.", "--watch", opts.Commands)
					}
					if file.FileSize(cmdFile) == 0 {
						logger.Panic(types.Usage("The file you specified ({0}) was found but contained no commands.", cmdFile).Error())
					}
				}

				if len(opts.Watchlist) == 0 {
					return types.Usage("The {0} option requires {1}.", "--watch", "a --watchlist file")
				} else {
					if opts.Watchlist != "existing" {
						watchList, err := filepath.Abs(opts.Watchlist)
						if err != nil || !file.FileExists(watchList) {
							return types.Usage("The {0} option requires {1} to exist.", "--watch", opts.Watchlist)
						}
						if file.FileSize(watchList) == 0 {
							logger.Panic(types.Usage("The file you specified ({0}) was found but contained no addresses.", watchList).Error())
						}
					}
				}

				if err := index.IsInitialized(chain, config.ExpectedVersion()); err != nil {
					if (errors.Is(err, index.ErrNotInitialized) || errors.Is(err, index.ErrIncorrectHash)) && !opts.Globals.IsApiMode() {
						logger.Panic(err)
					}
					return err
				}

				if opts.BatchSize < 1 {
					return types.Usage("The {0} option must be greater than zero.", "--batch_size")
				}
			} else {

			if opts.BatchSize != 8 {
				return types.Usage("The {0} option is not available{1}.", "--batch_size", " without --watch")
			} else {
				opts.BatchSize = 0
			}

			if opts.RunCount > 0 {
				return types.Usage("The {0} option is not available{1}.", "--run_count", " without --watch")
			}

			if opts.Sleep != 14 {
				return types.Usage("The {0} option is not available{1}.", "--sleep", " without --watch")
			}

*/

// TODO: Search for this function in trueblocks-core/src/apps/pkg/utils. It's identical to here.
// TODO: Make that function public and remove this one.
func downloadAndStore(url, filename string, dur time.Duration) ([]byte, error) {
	if file.FileExists(filename) {
		lastModDate, err := file.GetModTime(filename)
		if err != nil {
			return nil, err
		}
		if time.Since(lastModDate) < dur {
			data, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// If the file doesn't exist remotely, store an empty file
		if err := os.WriteFile(filename, []byte{}, 0644); err != nil {
			return nil, err
		}
		// Optionally update its mod time
		_ = file.Touch(filename)
		return []byte{}, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status %d %s for URL %s",
			resp.StatusCode, resp.Status, url)
	}

	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var prettyData []byte
	if json.Valid(rawData) {
		var jsonData interface{}
		if err := json.Unmarshal(rawData, &jsonData); err != nil {
			return nil, err
		}
		prettyData, err = json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return nil, err
		}
	} else {
		prettyData = rawData
	}

	if err := os.WriteFile(filename, prettyData, 0644); err != nil {
		return nil, err
	}

	_ = file.Touch(filename)

	return prettyData, nil
}
