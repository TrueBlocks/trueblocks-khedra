# Embedding & Re-Entrancy Guide

Scope: Desktop (Wails/other) embedding of wizard & dashboard. Install: `InstallFlow.md`; dashboard: `Dashboard.md`.

## 1. Principles
- Single Source (served HTML authoritative)
- Stateless pages (server reconstructs state per GET)
- Chrome suppression via `?embed=1`
- Theming via data attributes (e.g. `data-panel`, `data-role`)

## 2. State Endpoint
`GET /install/state` → `{configured,false|true,currentStep,sessionId,version,schema}`

## 3. Startup Algorithm (Pseudo)
```ts
async function boot(){
  try {
    const r = await fetch('/install/state',{cache:'no-store'});
    const st = await r.json();
    if(!st.configured){ openWebView(`/install/${st.currentStep}?embed=1`); return; }
    openWebView('/dashboard?embed=1');
  } catch { /* prompt to start daemon */ }
}
```

## 4. Re-Entrancy Sequence
Desktop polls until configured then switches to dashboard (or future SSE).

## 5. Failure / Recovery
 
| Scenario | Behavior | Desktop Action |
|----------|----------|----------------|
| Draft corrupt | Wizard restarts earliest step | Reload webview |
| Validation error | 422 with inline errors | User edits |
| Daemon stopped | Poll fails | Prompt restart |

## 6. CSS Contract
`data-embed="1"` hides header/footer; host may inject overrides after base stylesheet.

## 7. Local Exposure Guarantee
Endpoints bound to loopback only; desktop must not proxy externally.

## 8. Future SSE
Potential `/install/events` & `/dashboard/events` (deferred initial release).

## 9. Removal of Terminal Wizard
CLI `khedra init` prints & opens local URL; legacy TUI removed post adoption.

---
Embedding guide end.
