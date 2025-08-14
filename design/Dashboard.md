# Dashboard Specification (Post-Install)

Scope: Local control panel after configuration. Install flow: `InstallFlow.md`. Embedding: `Embedding.md`.

## 1. Objectives
- At-a-glance service & chain status.
- Pause/unpause actions (only pausable services).
- Paths & basic storage metrics (lightweight; avoid heavy scans inline).
- Re-run setup entry point / config download.
- Embeddable via `?embed=1` (chrome suppression, retains semantics).

## 2. State Endpoint
`GET /dashboard/state` JSON (example):
```json
{
  "version": "vX.Y.Z",
  "services": [
    {"name":"control","state":"running","pausable":false,"port":8338},
    {"name":"scraper","state":"running","pausable":true,"port":8081},
    {"name":"monitor","state":"paused","pausable":true,"port":8082},
    {"name":"api","state":"running","pausable":false,"port":8080},
    {"name":"ipfs","state":"running","pausable":false,"port":5001}
  ],
  "chains": [
    {"name":"mainnet","enabled":true,"rpc":"https://...","height":19123456,"reachable":true}
  ],
  "paths": {"data":"/Users/alex/khedra-data","logs":"/Users/alex/khedra-logs"},
  "logTail": ["2025-08-14T12:05:01Z service scraper cycle ..."],
  "pausedSummary": {"paused": ["monitor"], "totalPausable":2},
  "schema": 1
}
```
Caching: optional short (≈1s) cache to smooth bursts.

## 3. Functional Areas
1. Services table (Name / State / Port / Action)
2. Chains list (reachable + optional height)
3. Paths & storage metrics (size computed async / cached)
4. Actions (Re-run Setup, Download Config, Refresh)
5. Optional Log Tail (collapsible)
6. About (version/build/local-first tagline)

## 4. Wireframe (ASCII)
```
┌──────────────────────────────────────────────────────────────────────────────┐
│ Header: Khedra • vX.Y.Z • Local-First                                       │
├────────────────────────────┬────────────────────────────┬────────────────────┤
│ Services                   │ Chains                     │ Actions            │
│ control  running   -       │ mainnet height 19123456 ✓  │ Re-run Setup       │
│ scraper  running  8081 ⏸   │                            │ Download cfg       │
│ monitor  paused   8082 ▶   │                            │ Refresh            │
│ api      running  8080 -   │                            │                    │
│ ipfs     running  5001 -   │                            │                    │
├────────────────────────────┴─────────────┬──────────────┴────────────────────┤
│ Paths & Storage                          │ About / Build / Policy            │
└──────────────────────────────────────────┴────────────────────────────────────┘
```

## 5. HTML Skeleton (Rendered)
Below is a minimal, directly rendered HTML excerpt (no enclosing `<body>` so it can safely appear inside this Markdown). Attributes mirror the spec; trim or extend as needed.

<div id="dashboard-skeleton" data-view="dashboard" data-embed="0" style="border:1px solid #222;padding:0.75rem;font:14px/1.3 system-ui, sans-serif;max-width:980px;background:#111;color:#eee;">
  <header data-role="header" style="display:flex;gap:.5rem;align-items:baseline;margin-bottom:.75rem;">
    <h1 style="font-size:1.25rem;margin:0;">Khedra</h1>
    <span data-role="version" style="opacity:.7;">vX.Y.Z</span>
  </header>
  <main data-layout="grid" style="display:grid;grid-template-columns:2fr 1fr 1fr;gap:.75rem;">
  <section data-panel="services" style="border:1px solid #333;padding:.5rem;">
      <h2 style="margin-top:0;font-size:1rem;">Services</h2>
      <table data-table="services" style="width:100%;font-size:.75rem;border-collapse:collapse;">
        <thead><tr><th align="left">Name</th><th align="left">State</th><th align="left">Port</th><th></th></tr></thead>
        <tbody>
          <tr><td>control</td><td>running</td><td>-</td><td></td></tr>
          <tr><td>scraper</td><td>running</td><td>8081</td><td><button data-action="pause" style="font-size:.6rem;">⏸</button></td></tr>
          <tr><td>monitor</td><td>paused</td><td>8082</td><td><button data-action="unpause" style="font-size:.6rem;">▶</button></td></tr>
        </tbody>
      </table>
    </section>
  <section data-panel="chains" style="border:1px solid #333;padding:.5rem;">
      <h2 style="margin-top:0;font-size:1rem;">Chains</h2>
      <ul data-list="chains" style="padding-left:1rem;margin:0;font-size:.75rem;">
        <li>mainnet <span style="opacity:.6;">19123456 ✓</span></li>
      </ul>
    </section>
  <aside data-panel="actions" style="border:1px solid #333;padding:.5rem;">
      <h2 style="margin-top:0;font-size:1rem;">Actions</h2>
      <button data-action="rerun-setup" style="display:block;margin:0 0 .25rem 0;">Re-run Setup</button>
      <button data-action="download-config" style="display:block;margin:0 0 .25rem 0;">Download Config</button>
      <button data-action="refresh" style="display:block;">Refresh</button>
    </aside>
  <section data-panel="paths" style="border:1px solid #333;padding:.5rem;grid-column: span 2;">
      <h2 style="margin-top:0;font-size:1rem;">Paths & Storage</h2>
      <dl style="display:grid;grid-template-columns:max-content 1fr;gap:.25rem .5rem;font-size:.65rem;margin:0;">
        <dt>Data</dt><dd>/Users/alex/khedra-data</dd>
        <dt>Logs</dt><dd>/Users/alex/khedra-logs</dd>
        <dt>Size</dt><dd data-metric="data-size" style="opacity:.6;">(calculating...)</dd>
      </dl>
    </section>
  <section data-panel="about" style="border:1px solid #333;padding:.5rem;">
      <h2 style="margin-top:0;font-size:1rem;">About</h2>
      <p style="font-size:.65rem;margin:0;">All processing is local. <em>No telemetry.</em></p>
    </section>
  <section data-panel="logs" data-collapsible="true" style="border:1px solid #333;padding:.5rem;grid-column: span 3;">
      <h2 style="margin-top:0;font-size:1rem;">Log Tail</h2>
      <pre data-log="tail" style="background:#111;color:#9f9;padding:.5rem;font-size:.55rem;max-height:6rem;overflow:auto;margin:0;">2025-08-14T12:05:01Z service scraper cycle ...</pre>
    </section>
  </main>
</div>

<p style="font-size:.6rem;opacity:.55;">Rendered HTML snippet – view source to inspect <code>data-*</code> attributes.</p>

## 6. Interaction Contract
 
| Action | Endpoint | Result |
|--------|----------|--------|
| Pause | POST /service/<name>/pause | State → paused |
| Unpause | POST /service/<name>/unpause | State → running |
| Refresh | GET /dashboard/state | Update UI |
| Re-run Setup | GET /install/summary?embed=1 | Review wizard |
| Download Config | GET /config.yaml | Download |

## 7. Accessibility & Performance Notes
- Semantic table; manual refresh initially (avoid aggressive polling).
- Async size calculation for storage metrics; cache results.
- Optional SSE future.

## 8. Future Enhancements
- SSE push stream.
- Disk usage bar.
- Bulk pause/unpause.
- Chain latency badges.

---
Dashboard spec end.
