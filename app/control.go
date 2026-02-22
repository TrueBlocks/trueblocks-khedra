package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/file"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/control"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/coordinator"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/install"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	"github.com/TrueBlocks/trueblocks-sdk/v6/services"
)

func (k *KhedraApp) initializeControlSvc() error {
	if k.controlSvc != nil {
		return nil
	}

	if k.logger == nil {
		k.logger = types.NewLogger(types.Logging{Level: "info"})
	}
	if k.config == nil {
		cfg := types.NewConfig()
		k.config = &cfg
	}

	// ----------------------------------------------------------------------------------
	// Write initial control metadata (port from config); ignore errors (best-effort)
	k.controlSvc = services.NewControlService(k.logger.GetLogger())
	meta := control.NewMetadata(k.controlSvc.Port(), k.config.Version())
	_ = control.Write(meta)

	// Create coordinator for scraper-monitor coordination
	k.coordinator = coordinator.NewScraperMonitorCoordinator(k.logger.GetLogger())

	// Create all services using factory
	factory := NewServiceFactory(k.config, k.logger, k.coordinator)
	activeServices := factory.CreateAllServices(k.controlSvc)

	k.serviceManager = services.NewServiceManager(activeServices, k.logger.GetLogger())
	k.controlSvc.AttachServiceManager(k.serviceManager)

	// Add handlers AFTER serviceManager is created so dashboard state handler can access it
	_ = k.addHandlers()

	k.logger.Info("Control service initialized", "services", len(activeServices))
	return nil
}

func (k *KhedraApp) addHandlers() error {
	k.logger.Info("Adding control handlers")

	// ----------------------------------------------------------------------------------
	// Session store shared across state handler (placeholder; expanded later with inactivity logic)
	installSession := install.NewSessionStore()

	// ----------------------------------------------------------------------------------
	// Install state handler
	k.controlSvc.AddHandler("/install/state", func(w http.ResponseWriter, r *http.Request) {
		install.Handler(installSession, k.config.Version(), install.Configured())(w, r)
	})

	// ----------------------------------------------------------------------------------
	// Download current config (final if present, else synthesize from draft)
	k.controlSvc.AddHandler("/config.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		cfgFn := types.GetConfigFnNoCreate()
		if file.FileExists(cfgFn) {
			w.Header().Set("Content-Disposition", "attachment; filename=\"config.yaml\"")
			http.ServeFile(w, r, cfgFn)
			return
		}
		// Fallback: build YAML from draft in-memory
		d, err := install.LoadDraft()
		if err != nil || d == nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("config not found"))
			return
		}
		cfg := d.Config
		var buf bytes.Buffer
		if tpl, err := template.New("cfg").Parse(strings.TrimSpace(types.ConfigTemplate()) + "\n"); err == nil {
			_ = tpl.Execute(&buf, &cfg)
		}
		processed := types.RemoveZeroLines(buf.String())
		w.Header().Set("Content-Disposition", "attachment; filename=\"config.draft.yaml\"")
		if processed != "" {
			_, _ = w.Write([]byte(processed))
		} else {
			_, _ = w.Write([]byte("# draft config could not be rendered"))
		}
	})

	// ----------------------------------------------------------------------------------
	// RPC probe alias normalization (new simplified JSON shape) + rate limiting
	var rpcProbeDeprecLogged bool
	var rpcProbeMu sync.Mutex
	type probeWindow struct {
		count int
		reset time.Time
	}
	perSession := map[string]*probeWindow{}
	globalWindow := &probeWindow{count: 0, reset: time.Now().Add(10 * time.Second)}
	const perSessionLimit = 5
	const globalLimit = 60
	const windowDur = 10 * time.Second

	// ----------------------------------------------------------------------------------
	// /install/rpc-test
	k.controlSvc.AddHandler("/install/rpc-test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = w.Write([]byte(`{"error":"GET only"}`))
			return
		}
		// Rate limiting (per-session + global)
		postedSession := strings.TrimSpace(r.Header.Get("X-Khedra-Session"))
		var retryAfter time.Duration
		limited := false
		rpcProbeMu.Lock()
		now := time.Now()
		// reset global window if elapsed
		if now.After(globalWindow.reset) {
			globalWindow.count = 0
			globalWindow.reset = now.Add(windowDur)
		}
		// locate session window
		if postedSession != "" {
			pw, ok := perSession[postedSession]
			if !ok || now.After(pw.reset) {
				pw = &probeWindow{count: 0, reset: now.Add(windowDur)}
				perSession[postedSession] = pw
			}
			if pw.count >= perSessionLimit {
				limited = true
				retryAfter = pw.reset.Sub(now)
			} else if globalWindow.count >= globalLimit {
				limited = true
				retryAfter = globalWindow.reset.Sub(now)
			} else {
				pw.count++
				globalWindow.count++
			}
		} else { // no session -> apply only global limit
			if now.After(globalWindow.reset) {
				globalWindow.count = 0
				globalWindow.reset = now.Add(windowDur)
			}
			if globalWindow.count >= globalLimit {
				limited = true
				retryAfter = globalWindow.reset.Sub(now)
			} else {
				globalWindow.count++
			}
		}
		rpcProbeMu.Unlock()
		if limited {
			ra := int(math.Ceil(retryAfter.Seconds()))
			if ra < 1 {
				ra = 1
			}
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]any{"error": "rate_limited", "retryAfterSec": ra})
			return
		}
		url := strings.TrimSpace(r.URL.Query().Get("rpc"))
		if url == "" { // fallback to body (legacy) or default mainnet chain rpc
			_ = r.ParseForm()
			if v := r.FormValue("rpc"); v != "" {
				url = v
			}
		}
		if url == "" { // attempt to use mainnet from config if enabled
			if ch, ok := k.config.Chains["mainnet"]; ok && len(ch.RPCs) > 0 {
				url = ch.RPCs[0]
			}
		}
		res, err := rpc.PingRpc(url)
		if err != nil {
			k.logger.Debug("RPC ping failed", "url", url, "error", err)
		}
		payload := map[string]any{
			"ok":            res.OK,
			"chainId":       res.ChainID,
			"chainName":     res.ChainName,
			"clientVersion": res.ClientVersion,
			"error":         res.Error,
			"latencyMs":     res.LatencyMS,
		}
		if !res.OK {
			w.WriteHeader(http.StatusBadGateway)
		}
		_ = json.NewEncoder(w).Encode(payload)
	})

	// ----------------------------------------------------------------------------------
	// /install/rpc_probe
	k.controlSvc.AddHandler("/install/rpc_probe", func(w http.ResponseWriter, r *http.Request) {
		if !rpcProbeDeprecLogged {
			k.logger.Warn("/install/rpc_probe is deprecated; use /install/rpc-test")
			rpcProbeDeprecLogged = true
		}
		install.RpcProbeHandler(w, r) // legacy full response for backward compatibility
	})

	// ----------------------------------------------------------------------------------
	// Dashboard state endpoint (initial minimal implementation per spec)
	k.controlSvc.AddHandler("/dashboard/state", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.Header().Set("Content-Type", "application/json")
		// Build services slice, sorted alphabetically for stable UI
		var servicesJSON []map[string]any
		var names []string
		for name := range k.config.Services {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			svc := k.config.Services[name]
			state := "running"
			// Query the actual ServiceManager for real-time pause state
			if k.serviceManager != nil {
				if results, err := k.serviceManager.IsPaused(name); err == nil && len(results) > 0 {
					// Find the result for this specific service
					for _, result := range results {
						if result["name"] == name {
							if result["status"] == "paused" {
								state = "paused"
							}
							break
						}
					}
				}
			}
			servicesJSON = append(servicesJSON, map[string]any{
				"name":     name,
				"state":    state,
				"pausable": true, // all services now show pause/unpause button
				"port":     svc.Port,
			})
		}
		// Chains slice
		var chainsJSON []map[string]any
		for name, ch := range k.config.Chains {
			if !ch.Enabled {
				continue
			}
			rpc := ""
			if len(ch.RPCs) > 0 {
				rpc = ch.RPCs[0]
			}
			// Default blank chain info (no height fetch)
			entry := map[string]any{
				"name":    name,
				"enabled": ch.Enabled,
				"rpc":     rpc,
			}
			chainsJSON = append(chainsJSON, entry)
		}
		paths := map[string]string{
			"data":  k.config.IndexPath(),
			"cache": k.config.CachePath(),
			"logs":  k.config.Logging.Folder,
		}
		// Optimized log tail: only attempt when logging to file enabled
		var logTail []string
		logToFile := k.config.Logging.ToFile
		if logToFile {
			logFile := filepath.Join(k.config.Logging.Folder, k.config.Logging.Filename)
			if file.FileExists(logFile) {
				// Efficient tail read of last 15 non-empty lines without loading entire file
				if f, err := os.Open(logFile); err == nil {
					defer f.Close()
					// Seek from end in chunks
					const maxLines = 15
					const chunkSize = 8 * 1024
					var (
						pos  int64
						buf  []byte
						stat os.FileInfo
					)
					if stat, err = f.Stat(); err == nil {
						pos = stat.Size()
						var collected []string
						rem := ""
						for pos > 0 && len(collected) < maxLines {
							readSize := chunkSize
							if pos-int64(readSize) < 0 {
								readSize = int(pos)
							}
							pos -= int64(readSize)
							_, _ = f.Seek(pos, 0)
							buf = make([]byte, readSize)
							if _, err := f.Read(buf); err != nil {
								break
							}
							chunk := string(buf) + rem
							parts := strings.Split(chunk, "\n")
							// Last element may be partial; keep as rem for next iteration
							rem = parts[0]
							for i := len(parts) - 1; i >= 1 && len(collected) < maxLines; i-- { // skip parts[0] (carried)
								ln := strings.TrimSpace(parts[i])
								if ln != "" {
									collected = append(collected, ln)
								}
							}
						}
						// collected currently newest-first; reverse
						for i := len(collected) - 1; i >= 0; i-- {
							logTail = append(logTail, collected[i])
						}
					}
				}
			}
		}
		pausedSummary := map[string]any{"paused": []string{}, "totalPausable": 0}
		var paused []string
		// Query ServiceManager for actual paused services
		if k.serviceManager != nil {
			for name := range k.config.Services {
				if results, err := k.serviceManager.IsPaused(name); err == nil && len(results) > 0 {
					if results[0]["status"] == "paused" {
						paused = append(paused, name)
					}
				}
			}
		}
		sort.Strings(paused)
		pausedSummary["paused"] = paused
		pausedSummary["totalPausable"] = len(k.config.Services)
		resp := map[string]any{
			"version":         k.config.Version(),
			"services":        servicesJSON,
			"chains":          chainsJSON,
			"paths":           paths,
			"logTail":         logTail,
			"logToFile":       logToFile,
			"loggingFilename": k.config.Logging.Filename,
			"pausedSummary":   pausedSummary,
			"schema":          1,
		}
		enc := json.NewEncoder(w)
		_ = enc.Encode(resp)
	})

	// ----------------------------------------------------------------------------------
	// Dynamic chain add/remove endpoints for new UI
	k.controlSvc.AddHandler("/install/chain_add", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		rpcURL := strings.TrimSpace(r.URL.Query().Get("rpc"))
		if rpcURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"missing rpc"}`))
			return
		}
		// Probe JSON directly (reachability assumed if returns)
		res, err := rpc.PingRpc(rpcURL)
		if err != nil {
			k.logger.Debug("RPC ping failed during chain add", "url", rpcURL, "error", err)
		}
		if !res.OK || res.ChainID == "" {
			w.WriteHeader(http.StatusBadGateway)
			b, _ := json.Marshal(res)
			_, _ = w.Write(b)
			return
		}
		// Validate that chainId exists in local ChainList; reject if unknown
		cidStr := strings.TrimPrefix(strings.ToLower(res.ChainID), "0x")
		cidNum, errParse := strconv.ParseUint(cidStr, 16, 64)
		if errParse != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"ok":false,"error":"invalid chainId %s"}`, res.ChainID)
			return
		}
		if utils.GetChainListItem("~/.khedra", int(cidNum)) == nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"ok":false,"error":"unknown chainId %s not found in chain list"}`, res.ChainID)
			return
		}
		// Load draft and append / replace chain keyed by chain name or fallback numeric id
		draft, _ := install.LoadDraft()
		if draft == nil {
			draft = install.NewDraftFromConfig("")
		}
		if draft.Config.Chains == nil {
			draft.Config.Chains = map[string]types.Chain{}
		}
		name := res.ChainName
		// Prefer canonical chain list name if available (ensures consistency)
		if item := utils.GetChainListItem("~/.khedra", int(cidNum)); item != nil && item.Name != "" {

			toKey := func(s string) string {
				reg, _ := regexp.Compile("[^a-zA-Z0-9_.-]+")
				return reg.ReplaceAllString(strings.TrimSpace(s), "_")
			}
			name = toKey(item.Name)
		}
		if name == "" { // attempt to map via KnownChains by id first
			for _, kc := range install.KnownChains() {
				if uint64(kc.ChainID) == cidNum {
					name = kc.Name
					break
				}
			}
			if name == "" { // fallback synthetic
				name = fmt.Sprintf("chain-%d", cidNum)
			}
		}
		if cidNum == 1 && name != "mainnet" {
			name = "mainnet"
		}
		name = strings.ToLower(name)
		ch, exists := draft.Config.Chains[name]
		if !exists {
			ch = types.NewChain(name, int(cidNum))
		}
		ch.ChainID = int(cidNum)
		if len(ch.RPCs) == 0 {
			ch.RPCs = []string{rpcURL}
		} else if ch.RPCs[0] != rpcURL {
			ch.RPCs[0] = rpcURL
		}
		ch.Enabled = true
		draft.Config.Chains[name] = ch
		if err := install.SaveDraftAtomic(draft); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"ok":false,"error":"save failed"}`))
			return
		}
		// Create response with RPC validity
		chainWithStatus := map[string]any{
			"name":     ch.Name,
			"chainId":  ch.ChainID,
			"rpcs":     ch.RPCs,
			"enabled":  ch.Enabled,
			"rpcValid": true, // If we got here, probe was successful so RPC is valid
		}
		b, _ := json.Marshal(map[string]any{"ok": true, "chain": chainWithStatus, "probe": res})
		_, _ = w.Write(b)
	})

	// ----------------------------------------------------------------------------------
	// /isntall/chain_remove
	k.controlSvc.AddHandler("/install/chain_remove", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		name := strings.TrimSpace(r.URL.Query().Get("name"))
		if name == "" || name == "mainnet" {
			_, _ = w.Write([]byte(`{"ok":false}`))
			return
		}
		draft, _ := install.LoadDraft()
		if draft != nil && draft.Config.Chains != nil {
			delete(draft.Config.Chains, name)
			_ = install.SaveDraftAtomic(draft)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	// ----------------------------------------------------------------------------------
	// Simple ping endpoint to verify new binary deployed
	k.controlSvc.AddHandler("/install/ping", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		_, _ = w.Write([]byte("pong"))
	})

	// ----------------------------------------------------------------------------------

	// Unified live-update endpoint for config feedback (draft or real)
	k.controlSvc.AddHandler("/live-update/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Handle POST requests to update draft config
		if r.Method == http.MethodPost {
			// Load current draft
			draft, err := install.LoadDraft()
			if err != nil || draft == nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"failed to load draft config"}`))
				return
			}

			// Parse form data and update draft
			if err := r.ParseForm(); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error":"invalid form data"}`))
				return
			}

			// Apply form updates to draft config
			install.ApplyFormToDraft(draft, r.Form)

			// Save draft to disk immediately
			if err := install.SaveDraft(draft); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"failed to save draft config"}`))
				return
			}
		}

		// Return current config (GET or after POST update)
		var cfg any
		var source string

		// Check if we're in wizard mode by looking at the request path or headers
		// During wizard steps, ALWAYS show draft; only on dashboard show real config
		isWizardStep := false
		if referer := r.Header.Get("Referer"); referer != "" {
			isWizardStep = strings.Contains(referer, "/install/") && !strings.Contains(referer, "/dashboard")
		}
		// Also check current path from request context if available
		if !isWizardStep && r.URL.Path != "" {
			isWizardStep = strings.HasPrefix(r.URL.Path, "/install/") && !strings.Contains(r.URL.Path, "/dashboard")
		}

		draft, _ := install.LoadDraft()
		if isWizardStep && draft != nil {
			// During wizard steps, always show draft
			cfg = draft.Config
			source = "draft"
		} else if !isWizardStep && install.Configured() {
			// On dashboard, show real config if it exists
			loaded, err := LoadConfig()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"failed to load config"}`))
				return
			}
			cfg = loaded
			source = "real"
		} else if draft != nil {
			// Fallback: if draft exists, show it
			cfg = draft.Config
			source = "draft"
		} else {
			// No config available
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"no config available"}`))
			return
		}

		// Calculate estimates for the config response
		var meta map[string]any
		if configStruct, ok := cfg.(types.Config); ok {
			diskGB, hours := install.EstimateIndex(configStruct.General.Strategy, configStruct.General.Detail)
			meta = map[string]any{
				"EstDiskGB": diskGB,
				"EstHours":  hours,
			}
		}

		payload := map[string]any{"source": source, "config": cfg}
		if meta != nil {
			payload["config"] = map[string]any{
				"General":  cfg.(types.Config).General,
				"Chains":   cfg.(types.Config).Chains,
				"Services": cfg.(types.Config).Services,
				"Logging":  cfg.(types.Config).Logging,
				"Meta":     meta,
			}
		}
		b, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"failed to marshal config"}`))
			return
		}
		_, _ = w.Write(b)
	})

	// ----------------------------------------------------------------------------------
	// SetRootHandler
	k.controlSvc.SetRootHandler(func(w http.ResponseWriter, r *http.Request) {
		configured := install.Configured()

		// Determine persistent embed preference: query param overrides and sets cookie; cookie persists.
		const embedCookieName = "KHEDRA_EMBED"
		embedPref := false
		if qp := r.URL.Query().Get("embed"); qp != "" { // explicit user choice this request
			embedPref = qp == "1"
			http.SetCookie(w, &http.Cookie{Name: embedCookieName, Value: map[bool]string{true: "1", false: "0"}[embedPref], Path: "/", Expires: time.Now().Add(365 * 24 * time.Hour)})
		} else if c, err := r.Cookie(embedCookieName); err == nil {
			embedPref = c.Value == "1"
		}

		// Determine persistent debug preference (mirrors embed). If debug=1 show side panel with config JSON.
		const debugCookieName = "KHEDRA_DEBUG"
		debugPref := false
		if qp := r.URL.Query().Get("debug"); qp != "" { // explicit user choice this request
			debugPref = qp == "1"
			http.SetCookie(w, &http.Cookie{Name: debugCookieName, Value: map[bool]string{true: "1", false: "0"}[debugPref], Path: "/", Expires: time.Now().Add(365 * 24 * time.Hour)})
		} else if c, err := r.Cookie(debugCookieName); err == nil {
			debugPref = c.Value == "1"
		}

		// Helper function to build URLs with embed parameters when in embed mode
		buildURL := func(path string, extraParams ...string) string {
			if embedPref || len(extraParams) > 0 {
				u, _ := url.Parse(path)
				q := u.Query()
				if embedPref {
					q.Set("embed", "1")
					q.Set("debug", map[bool]string{true: "1", false: "0"}[debugPref])
				}
				// Add any extra parameters
				for i := 0; i < len(extraParams); i += 2 {
					if i+1 < len(extraParams) {
						q.Set(extraParams[i], extraParams[i+1])
					}
				}
				u.RawQuery = q.Encode()
				return u.String()
			}
			return path
		}

		// If not yet configured and user directly hits /dashboard, redirect them
		// into the install flow instead of falling through to the control service
		// default root (which exposes raw endpoint JSON).
		if !configured && r.URL.Path == "/dashboard" {
			http.Redirect(w, r, buildURL("/install/welcome"), http.StatusSeeOther)
			return
		}

		// Allowed wizard steps (in order) and quick lookup map.
		allowedWizardSteps := []string{"welcome", "paths", "chains", "index", "services", "logging", "summary"}
		allowedWizardStepSet := map[string]struct{}{}
		for _, s := range allowedWizardSteps {
			allowedWizardStepSet[s] = struct{}{}
		}

		// Helper to set the wizard step cookie consistently.
		setWizardStepCookie := func(step string, ttl time.Duration) {
			http.SetCookie(w, &http.Cookie{Name: "KHEDRA_WIZARD_STEP", Value: step, Path: "/", Expires: time.Now().Add(ttl)})
		}

		// Resume wizard based on cookie if user lands on root/dashboard.
		if r.URL.Path == "/" || r.URL.Path == "/dashboard" {
			if c, err := r.Cookie("KHEDRA_WIZARD_STEP"); err == nil {
				step := c.Value
				if _, ok := allowedWizardStepSet[step]; ok {
					http.Redirect(w, r, buildURL("/install/"+step), http.StatusFound)
					return
				}
			}
		}

		// Intercept install flow OR redirect root if not configured; otherwise serve dashboard
		if strings.HasPrefix(r.URL.Path, "/install") || (!configured && r.URL.Path == "/") {
			// Prepare debug JSON: during install always show draft; on dashboard show applied config (fallback to draft). Only if debug enabled.
			var debugConfigJSON string
			if debugPref {
				if strings.HasPrefix(r.URL.Path, "/install") { // draft-first view
					if d, err := install.LoadDraft(); err == nil && d != nil {
						if b, err2 := json.MarshalIndent(map[string]any{"file": install.DraftFilePath(), "config": d.Config}, "", "  "); err2 == nil {
							debugConfigJSON = string(b)
						}
					}
				}
				if debugConfigJSON == "" { // outside install or draft missing
					if cfg, err := LoadConfig(); err == nil {
						if b, err2 := json.MarshalIndent(map[string]any{"file": types.GetConfigFnNoCreate(), "config": cfg}, "", "  "); err2 == nil {
							debugConfigJSON = string(b)
						}
					} else if d, err3 := install.LoadDraft(); err3 == nil && d != nil {
						if b, err4 := json.MarshalIndent(map[string]any{"file": install.DraftFilePath(), "config": d.Config}, "", "  "); err4 == nil {
							debugConfigJSON = string(b)
						}
					}
				}
				if debugConfigJSON == "" {
					debugConfigJSON = "{}"
				}
			}

			serveStep := func(stepIdx int, tmplName string, data map[string]any) {
				files := []string{"templates/base.html", "templates/progress.html", "templates/" + tmplName}
				tmpl, err := loadTemplates(files...)
				if err != nil {
					k.logger.Error("template parse failed", "err", err, "tmpl", tmplName)
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("template error: " + err.Error()))
					return
				}
				if data == nil {
					data = map[string]any{}
				}
				steps := install.StepOrder
				data["Steps"] = steps
				data["Embed"] = embedPref
				data["Debug"] = debugPref
				if debugPref {
					data["DebugConfig"] = debugConfigJSON
				}
				// Always include session id so forms can post it back; Create one early if needed.
				data["SessionID"] = installSession.EnsureID()
				// Corruption flag: show one-time banner if a recent corruption replacement occurred
				if !configured {
					if install.ConsumeCorruptionFlag(24 * time.Hour) { // one-time consumption
						data["DraftCorrupt"] = true
					} else {
						data["DraftCorrupt"] = false
					}
				}
				// Simplest mapping: index comes directly from the URL-selected step.
				if stepIdx >= 0 {
					data["CurrentStepIndex"] = stepIdx
				} else {
					data["CurrentStepIndex"] = -1
				}
				if stepIdx >= 0 && stepIdx < len(steps) {
					data["StepName"] = steps[stepIdx]
					// Set/update cookie for resume; 7-day ttl
					setWizardStepCookie(steps[stepIdx], 7*24*time.Hour)
				} else {
					data["StepName"] = "dashboard"
				}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				_ = tmpl.ExecuteTemplate(w, "base.html", data)
			}

			// Extract posted session (if any) prior to handling specific steps for POST requests.
			const inactivityWindow = 5 * time.Minute
			var postedSession string
			if r.Method == http.MethodPost {
				_ = r.ParseForm()
				postedSession = strings.TrimSpace(r.FormValue("session"))
				if postedSession == "" { // allow header override (e.g. JS clients)
					postedSession = strings.TrimSpace(r.Header.Get("X-Khedra-Session"))
				}
				if postedSession != "" { // enforce if we have one; else treat as conflict
					// capture previous before possible takeover for logging
					prevID, prevLast := installSession.Get()
					status, takeover, current, last := installSession.Enforce(postedSession, inactivityWindow)
					if status == "conflict" {
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusConflict)
						_ = json.NewEncoder(w).Encode(map[string]any{"error": "session_conflict", "activeSession": current, "lastActivity": last.UTC().Format(time.RFC3339), "inactivitySecRequired": int(inactivityWindow.Seconds())})
						return
					}
					if takeover {
						w.Header().Set("X-Khedra-Session-Takeover", "1")
						if prevID != "" && prevID != current { // log only when a real replacement happened
							k.logger.Warn("session takeover", "prevSession", prevID, "prevLast", prevLast.UTC().Format(time.RFC3339), "newSession", current, "remote", r.RemoteAddr)
						}
					}
				} else if r.Method == http.MethodPost { // missing token on mutating request
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusConflict)
					_ = json.NewEncoder(w).Encode(map[string]any{"error": "session_required"})
					return
				}
			}

			switch r.Method {
			case http.MethodPost:
				k.logger.Info("install POST", "path", r.URL.Path, "ts", time.Now().Format(time.RFC3339))
			case http.MethodGet:
				k.logger.Info("install GET", "path", r.URL.Path, "ts", time.Now().Format(time.RFC3339))
			}

			// normalize trailing slash (except root and /install/ which maps to welcome)
			if r.URL.Path != "/" && len(r.URL.Path) > 1 && strings.HasSuffix(r.URL.Path, "/") && r.URL.Path != "/install/" {
				r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			}

			// -----// -----------------------------------------------------------------------------
			// Reset endpoint: removes draft & sends user to welcome to fully restart progress
			if r.URL.Path == "/install/reset" {
				if r.Method == http.MethodPost || r.Method == http.MethodGet { // accept both for simplicity
					if dpath := install.DraftFilePath(); dpath != "" {
						if err := os.Remove(dpath); err != nil && !os.IsNotExist(err) {
							k.logger.Warn("failed removing draft on reset", "err", err, "file", dpath)
						}
					}
					// Also remove any previous yaml backup to avoid confusion (best effort)
					cfgFn := types.GetConfigFnNoCreate()
					prev := filepath.Join(filepath.Dir(cfgFn), "config.draft.json") // already deleted; just in case
					_ = os.Remove(prev)
					// Reset wizard cookie to welcome
					setWizardStepCookie("welcome", 7*24*time.Hour)
					http.Redirect(w, r, buildURL("/install/welcome", "reset", "1"), http.StatusSeeOther)
					return
				}
			}

			if r.URL.Path == "/" && !configured {
				http.Redirect(w, r, buildURL("/install/welcome"), http.StatusFound)
				return
			}

			// Ensure draft exists (seed on first entry)
			if d, _ := install.LoadDraft(); d == nil {
				_ = install.SaveDraftAtomic(install.NewDraftFromConfig(""))
			}

			// -----// -----------------------------------------------------------------------------
			// Welcome
			if r.URL.Path == "/install" || r.URL.Path == "/install/" || r.URL.Path == "/install/welcome" {
				if r.Method == http.MethodPost {
					k.logger.Info("welcome submit", "remote", r.RemoteAddr, "ua", r.UserAgent())
					http.Redirect(w, r, buildURL("/install/paths"), http.StatusSeeOther)
					return
				}
				k.logger.Info("welcome view", "remote", r.RemoteAddr, "ua", r.UserAgent())
				serveStep(0, "welcome.html", map[string]any{"Reset": r.URL.Query().Get("reset")})
				return
			}

			// -----// -----------------------------------------------------------------------------
			// Paths step (phased validation)
			if r.URL.Path == "/install/paths" {
				if r.Method == http.MethodPost {
					if err := r.ParseForm(); err != nil {
						serveStep(1, "paths.html", map[string]any{"Errors": []install.FieldError{{Field: "general.dataFolder", Code: "bad_form", Message: "Invalid form submission"}}})
						return
					}
					df := strings.TrimSpace(r.FormValue("dataFolder"))
					draft, _ := install.LoadDraft()
					if draft == nil {
						draft = install.NewDraftFromConfig("")
					}
					draft.Config.General.DataFolder = df
					if err := install.SaveDraftAtomic(draft); err == nil {
						install.ClearCorruptionFlag()
					}
					ferrs := install.ValidateDraftPhase(draft, "step:paths")
					if len(ferrs) > 0 {
						serveStep(1, "paths.html", map[string]any{"DataFolder": draft.Config.General.DataFolder, "Errors": ferrs})
						return
					}
					http.Redirect(w, r, buildURL("/install/chains"), http.StatusSeeOther)
					return
				}
				draft, _ := install.LoadDraft()
				df := ""
				if draft != nil {
					df = draft.Config.General.DataFolder
				}
				serveStep(1, "paths.html", map[string]any{"DataFolder": df})
				return
			}

			// -----// -----------------------------------------------------------------------------
			// Index step (strategy + detail) with phased validation
			if r.URL.Path == "/install/index" {
				if r.Method == http.MethodPost {
					if err := r.ParseForm(); err != nil {
						serveStep(3, "index.html", map[string]any{"Errors": []install.FieldError{{Field: "general.strategy", Code: "bad_form", Message: "Invalid form submission"}}})
						return
					}
					strategy := r.FormValue("strategy")
					if strategy == "" {
						strategy = "download"
					}
					detail := r.FormValue("detail")
					if detail == "" {
						detail = "index"
					}
					if detail == "blooms" {
						detail = "bloom"
					}
					draft, _ := install.LoadDraft()
					if draft == nil {
						draft = install.NewDraftFromConfig("")
					}
					install.UpdateIndexStrategy(draft, strategy, detail)
					if err := install.SaveDraftAtomic(draft); err == nil {
						install.ClearCorruptionFlag()
					}
					ferrs := install.ValidateDraftPhase(draft, "step:index")
					if len(ferrs) > 0 {
						serveStep(3, "index.html", map[string]any{"Strategy": strategy, "Detail": detail, "Disk": draft.Meta.EstDiskGB, "Hours": draft.Meta.EstHours, "Errors": ferrs})
						return
					}
					http.Redirect(w, r, buildURL("/install/services"), http.StatusSeeOther)
					return
				}
				draft, _ := install.LoadDraft()
				strategy := "download"
				detail := "index"
				if draft != nil {
					if draft.Config.General.Strategy != "" {
						strategy = draft.Config.General.Strategy
					}
					if draft.Config.General.Detail != "" {
						detail = draft.Config.General.Detail
					}
				}
				disk, hours := 0, 0
				if draft != nil && draft.Meta.EstDiskGB > 0 && draft.Meta.EstHours > 0 {
					disk, hours = draft.Meta.EstDiskGB, draft.Meta.EstHours
				} else {
					disk, hours = install.EstimateIndex(strategy, detail)
				}
				serveStep(3, "index.html", map[string]any{"Strategy": strategy, "Detail": detail, "Disk": disk, "Hours": hours})
				return
			}

			// -----// -----------------------------------------------------------------------------
			// Chains step (phased validation)
			if r.URL.Path == "/install/chains" {
				if r.Method == http.MethodPost {
					if err := r.ParseForm(); err != nil {
						serveStep(2, "chains.html", map[string]any{"Errors": []install.FieldError{{Field: "chains.mainnet.rpc", Code: "bad_form", Message: "Invalid form submission"}}})
						return
					}
					action := r.FormValue("action")
					draft, _ := install.LoadDraft()
					if draft == nil {
						draft = install.NewDraftFromConfig("")
					}
					if draft.Config.Chains == nil {
						draft.Config.Chains = map[string]types.Chain{}
					}
					if action == "add" {
						name := strings.TrimSpace(r.FormValue("chain_to_add"))
						if install.IsKnownChain(name) {
							if _, exists := draft.Config.Chains[name]; !exists {
								var cid int
								for _, kc := range install.KnownChains() {
									if kc.Name == name {
										cid = kc.ChainID
										break
									}
								}
								ch := types.NewChain(name, cid)
								if name != "mainnet" {
									ch.Enabled = false
								}
								ch.RPCs = []string{"https://"}
								draft.Config.Chains[name] = ch
							}
						}
						if err := install.SaveDraftAtomic(draft); err == nil {
							install.ClearCorruptionFlag()
						}
						// fall through to GET render with validation errors if any
					} else if rem := r.FormValue("action_remove"); rem != "" && rem != "mainnet" {
						delete(draft.Config.Chains, rem)
						if err := install.SaveDraftAtomic(draft); err == nil {
							install.ClearCorruptionFlag()
						}
					} else { // update existing (Next button pressed)
						for name, ch := range draft.Config.Chains {
							enField := name + "_enabled" // Fix: match the actual form field names
							rpcField := "chain_rpc_" + name
							ch.Enabled = r.FormValue(enField) == "on"
							if rv := strings.TrimSpace(r.FormValue(rpcField)); rv != "" {
								if len(ch.RPCs) == 0 {
									ch.RPCs = []string{rv}
								} else {
									ch.RPCs[0] = rv
								}
							}
							draft.Config.Chains[name] = ch
						}
						if err := install.SaveDraftAtomic(draft); err == nil {
							install.ClearCorruptionFlag()
						}
						// ONLY validate when trying to proceed to next step (Next button)
						ferrs := install.ValidateDraftPhase(draft, "step:chains")
						if len(ferrs) == 0 {
							http.Redirect(w, r, buildURL("/install/index"), http.StatusSeeOther)
							return
						}
						// If validation fails, re-render with errors to block progression
					}
					// Render with potential errors (draft already loaded/modified above)
				}
				draft, _ := install.LoadDraft()
				existing := map[string]bool{}
				var chainRows []types.Chain
				for name, ch := range draft.Config.Chains {
					existing[name] = true
					chainRows = append(chainRows, ch)
				}
				var addable []install.KnownChain
				for _, kc := range install.KnownChains() {
					if !existing[kc.Name] {
						addable = append(addable, kc)
					}
				}
				sort.Slice(chainRows, func(i, j int) bool { return chainRows[i].ChainID < chainRows[j].ChainID })
				var ordered []types.Chain
				for _, ch := range chainRows {
					if ch.Name == "mainnet" {
						ordered = append([]types.Chain{ch}, ordered...)
					}
				}
				if len(ordered) == 0 || ordered[0].Name != "mainnet" {
					ordered = chainRows
				} else {
					for _, ch := range chainRows {
						if ch.Name != "mainnet" {
							ordered = append(ordered, ch)
						}
					}
				}
				// Add RPC validity information to each chain
				type ChainWithStatus struct {
					types.Chain
					RpcValid bool `json:"rpcValid"`
				}
				var orderedWithStatus []ChainWithStatus
				for _, ch := range ordered {
					chainWithStatus := ChainWithStatus{
						Chain:    ch,
						RpcValid: HasValidRpc(&ch, 2), // Quick check with 2 tries
					}
					orderedWithStatus = append(orderedWithStatus, chainWithStatus)
				}
				// No validation on entry - allow users to fix invalid RPCs
				serveStep(2, "chains.html", map[string]any{"Chains": orderedWithStatus, "Addable": addable, "Errors": []install.FieldError{}})
				return
			}

			// -----// -----------------------------------------------------------------------------
			// Services step (phased validation)
			if r.URL.Path == "/install/services" {
				if r.Method == http.MethodPost {
					if err := r.ParseForm(); err != nil {
						serveStep(4, "services.html", map[string]any{"Errors": []install.FieldError{{Field: "services", Code: "bad_form", Message: "Invalid form submission"}}})
						return
					}
					draft, _ := install.LoadDraft()
					if draft == nil {
						draft = install.NewDraftFromConfig("")
					}
					for name := range draft.Config.Services {
						_, present := r.Form[name+"_enabled"]
						svc := draft.Config.Services[name]
						svc.Enabled = present
						draft.Config.Services[name] = svc
					}
					if err := install.SaveDraftAtomic(draft); err == nil {
						install.ClearCorruptionFlag()
					}
					ferrs := install.ValidateDraftPhase(draft, "step:services")
					if len(ferrs) > 0 {
						servicesMap := map[string]bool{}
						for name, svc := range draft.Config.Services {
							servicesMap[name] = svc.Enabled
						}
						serveStep(4, "services.html", map[string]any{"Services": servicesMap, "Errors": ferrs})
						return
					}
					http.Redirect(w, r, buildURL("/install/logging"), http.StatusSeeOther)
					return
				}
				draft, _ := install.LoadDraft()
				servicesMap := map[string]bool{}
				for name, svc := range draft.Config.Services {
					servicesMap[name] = svc.Enabled
				}
				ferrs := install.ValidateDraftPhase(draft, "step:services")
				serveStep(4, "services.html", map[string]any{"Services": servicesMap, "Errors": ferrs})
				return
			}

			// -----// -----------------------------------------------------------------------------
			// Logging step (phased validation)
			if r.URL.Path == "/install/logging" {
				if r.Method == http.MethodPost {
					if err := r.ParseForm(); err != nil {
						serveStep(5, "logging.html", map[string]any{"Errors": []install.FieldError{{Field: "logging.level", Code: "bad_form", Message: "Invalid form submission"}}})
						return
					}
					lvl := r.FormValue("level")
					if lvl == "" {
						lvl = "info"
					}
					toFile := r.FormValue("toFile") == "1"
					folder := r.FormValue("folder")
					filename := r.FormValue("filename")
					maxSizeStr := r.FormValue("maxSize")
					maxBackupsStr := r.FormValue("maxBackups")
					maxAgeStr := r.FormValue("maxAge")
					compress := r.FormValue("compress") == "1"
					maxSize, _ := strconv.Atoi(maxSizeStr)
					if maxSize <= 0 {
						maxSize = 10
					}
					maxBackups, _ := strconv.Atoi(maxBackupsStr)
					if maxBackups <= 0 {
						maxBackups = 3
					}
					maxAge, _ := strconv.Atoi(maxAgeStr)
					if maxAge <= 0 {
						maxAge = 10
					}
					draft, _ := install.LoadDraft()
					if draft == nil {
						draft = install.NewDraftFromConfig("")
					}
					lg := draft.Config.Logging
					lg.Level = lvl
					lg.ToFile = toFile
					if folder != "" {
						lg.Folder = folder
					}
					if filename != "" {
						lg.Filename = filename
					}
					lg.MaxSize = maxSize
					lg.MaxBackups = maxBackups
					lg.MaxAge = maxAge
					lg.Compress = compress
					draft.Config.Logging = lg
					if err := install.SaveDraftAtomic(draft); err == nil {
						install.ClearCorruptionFlag()
					}
					ferrs := install.ValidateDraftPhase(draft, "step:logging")
					if len(ferrs) > 0 {
						serveStep(5, "logging.html", map[string]any{"Level": lg.Level, "ToFile": lg.ToFile, "Folder": lg.Folder, "Filename": lg.Filename, "MaxSize": lg.MaxSize, "MaxBackups": lg.MaxBackups, "MaxAge": lg.MaxAge, "Compress": lg.Compress, "Errors": ferrs})
						return
					}
					http.Redirect(w, r, buildURL("/install/summary"), http.StatusSeeOther)
					return
				}
				draft, _ := install.LoadDraft()
				data := map[string]any{}
				if draft != nil {
					lg := draft.Config.Logging
					data["Level"] = lg.Level
					data["ToFile"] = lg.ToFile
					data["Folder"] = lg.Folder
					data["Filename"] = lg.Filename
					data["MaxSize"] = lg.MaxSize
					data["MaxBackups"] = lg.MaxBackups
					data["MaxAge"] = lg.MaxAge
					data["Compress"] = lg.Compress
				} else {
					data["Level"] = "info"
				}
				serveStep(5, "logging.html", data)
				return
			}

			// -----// -----------------------------------------------------------------------------
			// Summary step (validate & apply)
			if r.URL.Path == "/install/summary" {
				if r.Method == http.MethodPost {
					if err := install.ApplyDraft(); err != nil {
						draft, _ := install.LoadDraft()
						ferrs := install.ValidateDraftPhase(draft, "final")
						serveStep(6, "summary.html", map[string]any{"Draft": draft, "Errors": ferrs, "Error": err.Error()})
						return
					}
					if cfg, err := LoadConfig(); err != nil {
						k.logger.Error("Failed to reload config after applying draft", "error", err)
					} else {
						k.config = &cfg
						k.logger = types.NewLogger(cfg.Logging)
						k.logger.Info("Config reloaded after install wizard completion. Services will pick up changes naturally.")
					}

					// NOTE: Skip service restart to avoid disrupting the control service.
					// Services will pick up config changes on their next iteration.
					// if err := k.RestartAllServices(); err != nil {
					//	k.logger.Error("Service restart failed after wizard finish", "error", err)
					// }
					setWizardStepCookie("dashboard", 30*24*time.Hour)
					http.Redirect(w, r, buildURL("/dashboard"), http.StatusSeeOther)
					return
				}
				draft, _ := install.LoadDraft()
				ferrs := install.ValidateDraftPhase(draft, "final")
				serveStep(6, "summary.html", map[string]any{"Draft": draft, "Errors": ferrs})
				return
			}

			// Unimplemented future steps placeholder
			if strings.HasPrefix(r.URL.Path, "/install/") {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				_, _ = w.Write([]byte("Install step not yet implemented."))
				return
			}
		}

		k.logger.Debug("root handler", "configured", configured, "path", r.URL.Path)
		// If configured and hitting root or /dashboard, render dashboard
		if configured && (r.URL.Path == "/" || r.URL.Path == "/dashboard") {
			// Prepare debug JSON for dashboard if debug enabled
			var debugConfigJSON string
			if debugPref {
				if cfg, err := LoadConfig(); err == nil {
					if b, err2 := json.MarshalIndent(map[string]any{"file": types.GetConfigFnNoCreate(), "config": cfg}, "", "  "); err2 == nil {
						debugConfigJSON = string(b)
					}
				}
				if debugConfigJSON == "" {
					debugConfigJSON = "{}"
				}
			}

			// minimal data for now
			files := []string{"templates/base.html", "templates/progress.html", "templates/dashboard.html"}
			tmpl, err := loadTemplates(files...)
			if err != nil {
				k.logger.Error("template parse failed", "err", err, "tmpl", "dashboard.html")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("template error: " + err.Error()))
				return
			}
			data := map[string]any{
				"Chains":           k.config.Chains,
				"Services":         k.config.Services,
				"Steps":            nil,
				"CurrentStepIndex": -1,
				"Embed":            embedPref,
				"Debug":            debugPref,
				"StepName":         "dashboard",
			}
			if debugPref {
				data["DebugConfig"] = debugConfigJSON
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_ = tmpl.ExecuteTemplate(w, "base.html", data)
			return
		}

		// Fallback
		k.controlSvc.DefaultRootHandler()(w, r)
	})

	// ----------------------------------------------------------------------------------
	// Control info endpoint returning metadata
	k.controlSvc.AddHandler("/control/info", func(w http.ResponseWriter, r *http.Request) {
		_ = r
		w.Header().Set("Content-Type", "application/json")
		var regenerated bool
		var meta control.Metadata
		if svcCfg, ok := k.config.Services["control"]; ok {
			m, regen, err := control.EnsureMetadata(svcCfg.Port, k.config.Version())
			if err == nil {
				meta = m
				regenerated = regen
			} else { // fall back to best-effort existing read; error surfaced in payload
				meta, _ = control.Read()
			}
		}
		payload := map[string]any{
			"ok":          true,
			"metadata":    meta,
			"regenerated": regenerated,
		}
		_ = json.NewEncoder(w).Encode(payload)
	})

	return nil
}
