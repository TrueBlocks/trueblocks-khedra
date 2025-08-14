# Control Service Review (Local‑First Rewrite)

Date: 2025-08-14

Primary Directive: Khedra (and all of TrueBlocks) is 100% local‑first. Control of a running process is strictly a local machine concern. Remote orchestration is out of scope by design – not deferred, simply not part of the philosophy.

This rewrite reframes the earlier review through that lens and focuses the remediation plan on eliminating ambiguity in local discovery rather than adding remote‑oriented safeguards.

## 1. Current Behavior (Through Local‑First Lens)
The *control service* is started automatically in `app/action_daemon.go` via `services.NewControlService(logger)`. It binds to the first free port in a hard‑coded descending list (8338–8335). The CLI discovers it by blindly probing those ports’ root paths (`/`). The first HTTP 200 wins. This heuristic is the core problem: discovery is probabilistic instead of authoritative.

### 1.1 Discovery Flow (Today)
1. `pause` / `unpause` commands call `findControlServiceURL()`.
2. That scans ports [8338,8337,8336,8335], hitting `http://localhost:<port>` with a generic ping.
3. First port responding to TCP connect + HTTP 200 is assumed to be the control service.
4. Commands then call `/<endpoint>?name=<service>` with HTTP GET.
5. Response JSON is parsed into `[]map[string]string` and printed with colored free‑text statuses.

### 1.2 isRunning() Guard
`app.isRunning()` repeats the same scan; any unrelated HTTP listener returning 200 on those ports is interpreted as “Khedra already running.” This is a false positive risk that directly impedes startup.

### 1.3 Configuration Exposure
Port choice is implicit (range only). There is no declared control service configuration and no authoritative run metadata (PID + chosen port). Without a canonical artifact, every consumer must guess.

### 1.4 API Semantics
- Mutations (`/pause`, `/unpause`) use HTTP GET (semantically wrong, though purely local).
- Lack of a distinctive identity endpoint (e.g. `/control/info`); root probing cannot differentiate control vs any other tiny server.
- Response shape is untyped free‑text.

### 1.5 Service Enumeration
Hard‑coded list (`"scraper, monitor, all"`) forces scattered edits; the authoritative list lives in runtime state but is not queried.

### 1.6 Local‑Only Risk Summary

| Category | Risk | Local Impact |
|----------|------|--------------|
| False Identification | Random process misread as control | Pause/unpause fails or acts on wrong target |
| Startup Block | Port occupant blocks launch | User incorrectly told daemon is running |
| Weak Semantics | GET for mutations | Ambiguity; tool churn if later fixed |
| Fragile Parsing | Free‑text status strings | CLI coloring / messaging drifts with upstream text |
| Static Range | Fixed port window | Collisions on systems with those ports in use |
| Hidden State | No PID/port artifact | CLI must guess; scripts unreliable |

## 2. Design Tenets & Scope
### 2.1 Local‑First Tenets
1. No remote control surface will be added; control stays loopback or tighter (domain socket).
2. Deterministic discovery over heuristic probing.
3. Human‑inspectable, minimal on‑disk artifact is the source of truth.
4. Explicit identity handshake (no accidental cross‑process confusion).
5. Strong preference for OS primitives (PID existence, filesystem) over network scans.

### 2.2 Concrete Goals
- Deterministic discovery (single authoritative metadata file + identity endpoint).
- Eliminate false positives in `isRunning()`.
- Structured, versioned JSON responses.
- Dynamic enumeration of pausable services (no hard‑coded list in CLI).
- Correct HTTP semantics (POST for mutations) while keeping purely local binding.
- Optionally offer a Unix domain socket alternative (future extension) – noted but not required immediately.

### 2.3 Explicit Non‑Goals
- Remote management (any hostname other than 127.0.0.1 / ::1 or local socket).
- Network authentication layers (shared secrets, tokens) – unnecessary if never exposed remotely.
- TLS (no remote exposure implies no need). 
- Expanded health/metrics (orthogonal to discovery fix).

## 3. Proposed Adjustments
### 3.1 Config Additions (`config.yaml`)
```yaml
services:
  control:
    enabled: true      # default true
    port: 8338         # single explicit port (no fallback scan)
    socket: ""         # (optional future) path to unix domain socket; if set, overrides port
```
Rules:
- If `socket` is set, only the socket is used (no TCP listener).
- If `enabled: false`, CLI pause/unpause emits clear instruction (“control service disabled in config”).

### 3.2 Run Metadata Artifact
Write (and treat as authoritative) JSON at `~/.khedra/run/control.json`:
```json
{
  "pid": 12345,
  "port": 8338,
  "socket": "",
  "started": "2025-08-14T12:34:56Z",
  "version": "vX.Y.Z",
  "schema": 1
}
```
Lifecycle:
- Created immediately after successful bind.
- On graceful shutdown, attempt removal (best effort; stale file is tolerated via validation logic).
- Stale detection: if PID absent or endpoint handshake fails, overwrite with new instance.

### 3.3 Discovery (Authoritative)
1. Read metadata file.
2. Validate PID (signal 0) – if dead, treat as not running.
3. If `socket` populated: connect & GET `/control/info` over that socket. Else GET `http://127.0.0.1:<port>/control/info`.
4. Verify response contains expected signature fields (`service:"control"`, `schema:1`).
5. Only then proceed; otherwise treat daemon as not running.

### 3.4 Endpoint Set (Local Loopback Only)

| Endpoint | Method | Purpose | Response |
|----------|--------|---------|----------|
| `/control/info` | GET | Identity + versions + pausable list + states | ControlInfo |
| `/control/pause` | POST | Pause one or more services | ControlResults |
| `/control/unpause` | POST | Unpause one or more services | ControlResults |
| `/control/status` | GET | Alias of `/control/info` | ControlInfo |

Request body for pause/unpause:
```json
{ "services": ["scraper", "monitor"], "all": false }
```
Alternatively allow query param `all=true` if body omitted.

### 3.5 Authentication
None. Local‑first directive: control channel is never exposed off loopback / socket. Security boundary is the host user context.

### 3.6 Structured Types (Go)
Internal types (versioned with `schema` in metadata + `ControlInfo`):
```go
type ControlInfo struct {
  Version  string           `json:"version"`
  PID      int              `json:"pid"`
  Port     int              `json:"port"`
  Services []ServiceStatus  `json:"services"`
}

type ServiceStatus struct {
  Name   string `json:"name"`
  State  string `json:"state"` // running|paused|not-pausable
}

type ControlResult struct {
  Name     string `json:"name"`
  Previous string `json:"previous"`
  Current  string `json:"current"`
  Changed  bool   `json:"changed"`
  Error    string `json:"error,omitempty"`
}

type ControlResults struct {
  Results []ControlResult `json:"results"`
}
```
CLI formatting logic transitions to these structs (color based on `Changed`+`Current`).

### 3.7 Migration Strategy
Phase A: Add metadata file + new endpoints; CLI prefers them. Keep legacy `/pause`, `/unpause`, `/isPaused` for one release with a single deprecation notice.

Phase B: Remove legacy endpoints and heuristic port scan; refuse to operate if metadata missing (explicit guidance to restart daemon).

### 3.8 isRunning() Logic
1. Read metadata.
2. Check PID alive.
3. Perform identity request to `/control/info`.
4. Success → running; failure → not running (no port scans). Stale file silently ignored.

### 3.9 Error Handling & Timeouts
Timeouts tuned to local expectations (e.g. 1s connect, 2s total). Retries unnecessary on localhost – fail fast and report.

### 3.10 Minimal Test Coverage
1. Discovery success (valid metadata + live endpoint).
2. Stale metadata (dead PID) → treated as not running.
3. Pause request (POST) returns structured result and CLI formats correctly.
4. Legacy fallback path still works during Phase A (environment flag to simulate absence of new endpoint).

## 4. Implementation Steps (Ordered)
1. Add `control` config section (enabled, port, socket placeholder) in `pkg/types/config.go` with defaults.
2. Enforce single explicit port bind (fail fast if unavailable rather than scanning). Optionally skip if socket mode later.
3. After successful bind, write metadata file (atomic write: write temp then rename).
4. Implement `/control/info`, `/control/pause`, `/control/unpause`, `/control/status` handlers (wrapping SDK service manager).
5. Add structured types and JSON marshaling.
6. Update CLI discovery to use metadata + identity, remove port scan; implement legacy fallback for Phase A.
7. Switch pause/unpause to POST with JSON body (legacy GET fallback).
8. Update `isRunning()` to metadata approach.
9. Add tests (discovery, stale metadata, pause formatting, legacy fallback).
10. Update documentation (`command_line_interface.md`) marking legacy endpoints deprecated (Phase A). Remove their mention in Phase B.

## 5. Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| SDK lacks explicit port override | If unavoidable, document fixed port; still write metadata to disambiguate |
| User scripts rely on legacy GET endpoints | Maintain for one release with clear warning |
| Stale metadata after crash | PID + identity handshake validation before acceptance |
| Port already bound | Fail fast with actionable error; user updates config |

## 6. Acceptance Criteria
- Daemon start produces `~/.khedra/run/control.json` with pid, port, schema, version.
- `khedra pause scraper` performs POST `/control/pause` (verified via debug log) and succeeds.
- `isRunning()` returns false when unrelated service listens on 8338 (no valid metadata + identity).
- Legacy GET endpoints still operate with single deprecation warning during Phase A.
- Removing metadata file while daemon runs causes next CLI call to reconstruct it (self‑healing) or instruct restart.

## 7. Future (Still Local) Enhancements (Not in Current Work)
- Unix domain socket activation (eliminate TCP even on loopback).
- Event stream (SSE/WebSocket) purely local for UI tooling.
- Fine‑grained per‑service status detail (uptime counters) if ever needed for UX.

## 8. Next Action
Proceed to implement Steps 1–5 (core service + metadata + endpoints) then update CLI (Steps 6–8) and add tests (Step 9).

---
Rewritten with explicit local‑first directive and focused on deterministic discovery.
