package install

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/utils"
)

var (
	probeCacheHead = map[string]rpc.PingResult{}
	probeCacheJSON = map[string]rpc.PingResult{}
	probeCacheMu   sync.Mutex
	probeTTL       = 30 * time.Second

	// rate limiting state
	ratelimitMu      sync.Mutex
	sessionWindow    = 10 * time.Second
	perSessionLimit  = 5
	globalLimitTotal = 60
	sessionHits      = map[string][]time.Time{}
	globalHits       []time.Time
)

func sanitizeURL(raw string) (string, error) {
	if raw == "" {
		return "", errors.New("empty url")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("scheme must be http or https")
	}
	if u.Host == "" {
		return "", errors.New("missing host")
	}
	// scrub credentials
	u.User = nil
	return u.String(), nil
}

func headProbe(ctx context.Context, raw string) rpc.PingResult {
	now := time.Now()
	pr := rpc.PingResult{URL: raw, Mode: "head", CheckedAt: now.Unix()}
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, raw, nil)
	client := &http.Client{Timeout: 4 * time.Second}
	start := time.Now()
	if resp, err := client.Do(req); err == nil {
		resp.Body.Close()
		pr.StatusCode = resp.StatusCode
		// Consider 2xx/3xx success; also treat common HEAD rejections from JSON-RPC gateways as reachable
		switch resp.StatusCode {
		case 200, 201, 202, 204, 301, 302, 307, 308, 400, 401, 403, 405, 415, 429, 503:
			pr.OK = true
			if resp.StatusCode >= 400 { // informative but not fatal
				pr.Error = resp.Status
			}
		default:
			// Some providers return other 4xx codes but are still reachable; treat any non-404 < 500 as reachable but flagged.
			if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 404 {
				pr.OK = true
				pr.Error = resp.Status
			} else {
				pr.Error = resp.Status
			}
		}
		elapsed := time.Since(start).Milliseconds()
		if pr.OK {
			slog.Info("rpc head probe ok", "url", raw, "status", resp.StatusCode, "ms", elapsed)
		} else {
			slog.Info("rpc head probe not ok", "url", raw, "status", resp.StatusCode, "err", pr.Error, "ms", elapsed)
		}
	} else {
		pr.Error = err.Error()
		slog.Info("rpc head probe error", "url", raw, "err", err.Error())
	}
	return pr
}

func jsonProbe(ctx context.Context, raw string, expected string) rpc.PingResult {
	_ = ctx
	result, err := rpc.PingRpc(raw)
	if err != nil {
		slog.Info("rpc json probe error", "url", raw, "err", err.Error())
	}
	if result == nil {
		return rpc.PingResult{URL: raw, Mode: "json", CheckedAt: time.Now().Unix(), ExpectedChain: expected}
	}
	result.Mode = "json"
	result.ExpectedChain = expected
	if result.OK {
		slog.Info("rpc json probe ok", "url", raw, "chainId", result.ChainID, "expected", expected)
	}
	return *result
}

// RpcProbeHandler returns minimal information about chain RPC reachability for quick UI validation.
// Query param: url=... (single). It lightly caches recent results to avoid spamming endpoints.
func RpcProbeHandler(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("url")
	mode := r.URL.Query().Get("mode") // head|json (default head)
	expected := r.URL.Query().Get("expected")
	w.Header().Set("Content-Type", "application/json")
	now := time.Now()
	slog.Debug("rpc probe request", "url", raw, "mode", mode, "expected", expected)
	// Explicitly reject websocket schemes with clearer message (UI previously showed generic 'unreachable')
	if strings.HasPrefix(strings.ToLower(raw), "ws://") || strings.HasPrefix(strings.ToLower(raw), "wss://") {
		_ = json.NewEncoder(w).Encode(rpc.PingResult{URL: raw, OK: false, Error: "websocket RPC endpoints (ws://, wss://) are not supported; please use http:// or https://", CheckedAt: now.Unix(), Mode: mode})
		return
	}
	// rudimentary session-based rate limiting (session id passed via header or cookie later; fallback remote addr)
	sessionID := r.Header.Get("X-Khedra-Session")
	if sessionID == "" {
		sessionID = r.RemoteAddr
	}
	globalExceeded := false
	ratelimitMu.Lock()
	// prune old global hits
	cutoff := now.Add(-sessionWindow)
	var kept []time.Time
	for _, t := range globalHits {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	globalHits = kept
	if len(globalHits) >= globalLimitTotal {
		globalExceeded = true
	}
	// session hits
	hits := sessionHits[sessionID]
	kept = kept[:0]
	for _, t := range hits {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	hits = kept
	if !globalExceeded && len(hits) < perSessionLimit {
		hits = append(hits, now)
		globalHits = append(globalHits, now)
		sessionHits[sessionID] = hits
		ratelimitMu.Unlock()
	} else {
		sessionHits[sessionID] = hits
		ratelimitMu.Unlock()
		retry := int(sessionWindow.Seconds())
		w.WriteHeader(http.StatusTooManyRequests)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": "rate_limited", "retryAfterSec": retry})
		return
	}
	if raw == "" {
		_ = json.NewEncoder(w).Encode(rpc.PingResult{URL: raw, OK: false, Error: "missing url", CheckedAt: now.Unix(), Mode: mode})
		return
	}
	sanitized, err := sanitizeURL(raw)
	if err != nil {
		slog.Info("rpc probe invalid url", "raw", raw, "err", err.Error(), "mode", mode)
		_ = json.NewEncoder(w).Encode(rpc.PingResult{URL: raw, OK: false, Error: err.Error(), CheckedAt: now.Unix(), Mode: mode})
		return
	}
	if mode == "" {
		mode = "head"
	}
	slog.Info("rpc probe start", "url", sanitized, "mode", mode, "expected", expected)
	// cache lookup
	probeCacheMu.Lock()
	switch mode {
	case "head":
		if res, ok := probeCacheHead[sanitized]; ok && now.Sub(time.Unix(res.CheckedAt, 0)) < probeTTL {
			slog.Debug("rpc probe cache hit", "url", sanitized, "mode", mode)
			probeCacheMu.Unlock()
			_ = json.NewEncoder(w).Encode(res)
			return
		}
	case "json":
		cacheKey := sanitized // ignore expected chain id for caching; we accept whatever comes back
		if res, ok := probeCacheJSON[cacheKey]; ok && now.Sub(time.Unix(res.CheckedAt, 0)) < probeTTL {
			slog.Debug("rpc probe cache hit", "url", sanitized, "mode", mode)
			probeCacheMu.Unlock()
			_ = json.NewEncoder(w).Encode(res)
			return
		}
	}
	probeCacheMu.Unlock()

	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()
	var pr rpc.PingResult
	startProbe := time.Now()
	if mode == "json" {
		pr = jsonProbe(ctx, sanitized, expected)
	} else {
		pr = headProbe(ctx, sanitized)
		// Fallback: some gateways reject HEAD outright; if not OK try JSON quickly
		if !pr.OK {
			jp := jsonProbe(ctx, sanitized, expected)
			// If JSON succeeded, prefer that result and mark fallback
			if jp.OK {
				jp.FallbackJSON = true
				pr = jp
			}
		}
	}
	pr.LatencyMS = time.Since(startProbe).Milliseconds()
	slog.Info("rpc probe result", "url", sanitized, "mode", pr.Mode, "ok", pr.OK, "status", pr.StatusCode, "err", pr.Error, "chainId", pr.ChainID, "updated", pr.Updated)
	slog.Debug("rpc probe result", "url", sanitized, "mode", pr.Mode, "ok", pr.OK, "status", pr.StatusCode, "err", pr.Error, "chainId", pr.ChainID, "updated", pr.Updated)
	probeCacheMu.Lock()
	if pr.Mode == "head" {
		probeCacheHead[sanitized] = pr
	} else {
		probeCacheJSON[sanitized] = pr
	}
	probeCacheMu.Unlock()

	// Always adopt returned chainId (ignore expected) and derive chain name via ChainList
	if pr.Mode == "json" && pr.OK && pr.ChainID != "" {
		cidStr := strings.TrimPrefix(strings.ToLower(pr.ChainID), "0x")
		if n, err := strconv.ParseUint(cidStr, 16, 64); err == nil {
			if item := utils.GetChainListItem("~/.khedra", int(n)); item != nil {
				pr.ChainName = item.Name
			}
			chainKey := r.URL.Query().Get("chain")
			if chainKey != "" { // update existing draft chain slot if provided
				if d, _ := LoadDraft(); d != nil {
					if ch, ok := d.Config.Chains[chainKey]; ok {
						if ch.ChainID != int(n) {
							oldID := ch.ChainID
							ch.ChainID = int(n)
							d.Config.Chains[chainKey] = ch
							_ = SaveDraftAtomic(d)
							pr.Updated = true
							slog.Info("rpc probe chainId updated", "chain", chainKey, "old", oldID, "new", ch.ChainID)
						}
					}
				}
			}
		}
	}
	_ = json.NewEncoder(w).Encode(pr)
}
