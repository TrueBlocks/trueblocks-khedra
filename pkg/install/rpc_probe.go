package install

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type rpcProbeResult struct {
	URL           string `json:"url"`
	OK            bool   `json:"ok"`
	ChainID       string `json:"chainId,omitempty"`
	ChainName     string `json:"chainName,omitempty"`
	ExpectedChain string `json:"expectedChainId,omitempty"`
	// ChainMismatch removed from UI error semantics; always treat reachable as OK
	ClientVersion string `json:"clientVersion,omitempty"`
	Mode          string `json:"mode"` // head or json
	Error         string `json:"error,omitempty"`
	StatusCode    int    `json:"statusCode,omitempty"`
	CheckedAt     int64  `json:"checkedAt"`
	LatencyMS     int64  `json:"latencyMs,omitempty"`
	FallbackJSON  bool   `json:"fallbackJson,omitempty"`
	Updated       bool   `json:"updated,omitempty"` // true if draft chainId was updated
}

// RpcProbeJSON provides external packages a simple way to perform a JSON probe (eth_chainId + clientVersion)
// without caching side-effects of the HTTP handler. It returns a rpcProbeResult with Mode=json.
func RpcProbeJSON(ctx context.Context, raw string) rpcProbeResult {
	return jsonProbe(ctx, raw, "")
}

var (
	probeCacheHead = map[string]rpcProbeResult{}
	probeCacheJSON = map[string]rpcProbeResult{}
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

// jsonRpcRequest / response minimal structures
type jsonRpcRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRpcRespErr `json:"error"`
}

type jsonRpcRespErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

// trimQuotes removes leading/trailing double quotes from a JSON string value.
func trimQuotes(s string) string {
	return strings.Trim(s, "\"")
}

func headProbe(ctx context.Context, raw string) rpcProbeResult {
	now := time.Now()
	pr := rpcProbeResult{URL: raw, Mode: "head", CheckedAt: now.Unix()}
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

func jsonProbe(ctx context.Context, raw string, expected string) rpcProbeResult {
	now := time.Now()
	pr := rpcProbeResult{URL: raw, Mode: "json", CheckedAt: now.Unix(), ExpectedChain: expected}
	client := &http.Client{Timeout: 6 * time.Second}
	// eth_chainId
	chainIDHex, status, rawBody, err := callSimple(ctx, client, raw, "eth_chainId")
	pr.StatusCode = status
	if err != nil {
		pr.Error = err.Error()
		// Log raw body snippet for diagnostics (truncate to 160 chars)
		snippet := string(rawBody)
		if len(snippet) > 160 {
			snippet = snippet[:160] + "…"
		}
		slog.Info("rpc json probe chainId error", "url", raw, "status", status, "err", err.Error(), "body", snippet)
		return pr
	}
	pr.ChainID = trimQuotes(chainIDHex)
	// clientVersion (best effort) – ignore error but capture status if first call succeeded.
	if cv, _, _, err2 := callSimple(ctx, client, raw, "web3_clientVersion"); err2 == nil {
		pr.ClientVersion = trimQuotes(cv)
	}
	pr.OK = true
	slog.Info("rpc json probe ok", "url", raw, "status", status, "chainId", pr.ChainID, "expected", expected)
	return pr
}

func callSimple(ctx context.Context, client *http.Client, rawURL, method string) (string, int, []byte, error) {
	payload := jsonRpcRequest{JSONRPC: "2.0", ID: 1, Method: method, Params: []interface{}{}}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, rawURL, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", resp.StatusCode, body, fmt.Errorf("status %s", resp.Status)
	}
	var jr jsonRpcResponse
	if err := json.Unmarshal(body, &jr); err != nil {
		return "", resp.StatusCode, body, err
	}
	if jr.Error != nil {
		return "", resp.StatusCode, body, fmt.Errorf("rpc error %d %s", jr.Error.Code, jr.Error.Message)
	}
	if jr.Result == nil {
		return "", resp.StatusCode, body, errors.New("empty result")
	}
	return string(jr.Result), resp.StatusCode, body, nil
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
		_ = json.NewEncoder(w).Encode(rpcProbeResult{URL: raw, OK: false, Error: "websocket RPC endpoints (ws://, wss://) are not supported; please use http:// or https://", CheckedAt: now.Unix(), Mode: mode})
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
		_ = json.NewEncoder(w).Encode(rpcProbeResult{URL: raw, OK: false, Error: "missing url", CheckedAt: now.Unix(), Mode: mode})
		return
	}
	sanitized, err := sanitizeURL(raw)
	if err != nil {
		slog.Info("rpc probe invalid url", "raw", raw, "err", err.Error(), "mode", mode)
		_ = json.NewEncoder(w).Encode(rpcProbeResult{URL: raw, OK: false, Error: err.Error(), CheckedAt: now.Unix(), Mode: mode})
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
	var pr rpcProbeResult
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
