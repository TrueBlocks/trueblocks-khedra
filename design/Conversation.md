# Conversation State Summary (for Workspace Reopen)

Date: 2025-08-14

## 1. Goal
Implement a local-first HTML installation wizard for Khedra (replacing terminal wizard) plus post‑install dashboard and desktop embedding support, keeping everything loopback-only and minimal.

## 2. Key Design Docs
- `InstallFlow.md`: Authoritative wizard specification (7 steps: welcome, paths, index, chains, services, logging, summary; atomic persistence; session mgmt; `/install/state`; embed mode).
- `Dashboard.md`: Post-install control panel (services, chains, paths, actions, log tail) with `/dashboard/state` JSON contract.
- `Embedding.md`: Desktop (Wails/WebView) embedding playbook and re‑entrancy behavior.
- `DocumentationOverhaulPlan.md`: Deferred plan for full documentation rewrite once new flow is shipped.

## 3. Core Decisions
- Endpoints standardized: step pages under `/install/<step>`, apply via POST `/install/summary`, state via GET `/install/state`, RPC test via GET `/install/rpc-test`.
- Single session enforced (409 on conflicting POST; takeover after 5m inactivity) with hidden session token field (to implement).
- Draft persistence: `config.draft.json` (atomic temp + fsync + rename). One backup of previous final config on apply; draft deleted only on success.
- Index modes enum: `bloom | full | scratch` (naming may map to bloom/index/from-scratch internally).
- Rate limiting (minimal) for `/install/rpc-test`: 5 requests / 10s window per client; 60 cap total.
- Accessibility: minimal baseline (labels, aria-current on progress, skip link).
- Re-run setup: non-destructive incremental edit; optional future reset deferred.
- Local-only: all endpoints bind loopback (future unix socket acceptable); no external assets/telemetry.

## 4. Open Implementation Work (See ToDoList)
High-level phases:
1. Install mode detection + `/install/state` + root redirect.
2. Draft persistence + atomic write helpers.
3. Templates & routing for each step (GET/POST pattern, PRG).
4. Validation logic + RPC test + rate limiting.
5. Embed mode chrome suppression & data attributes.
6. Progress component.
7. Session enforcement + takeover.
8. Apply/finalization (atomic config write + backup).
9. Dashboard gating / link integration.
10. Accessibility baseline tasks.
11. Tests (unit + light integration, minimal scope per design).
12. CLI `khedra init` adjustment (open URL, exit) & deprecation of terminal wizard.
13. Documentation cross-links (not the full rewrite).
14. Future scaffolds (SSE stub, reset stub) without full implementation.

## 5. Current Implementation Status (Partial)
- New package `pkg/install/` added:
  - `state.go`: `State` struct, `SessionStore` (not yet used), handler stub returning placeholder currentStep (`welcome`).
  - `validity.go`: `Configured()` shallow check (file existence only; does not inspect mainnet RPC yet).
- Terminal wizard still active (`init` command unchanged) – HTTP wizard not wired.
- No root redirect or `/install/state` route mounted yet (control service source in sibling module not visible in current workspace segment).
- Session ID generation not yet implemented or persisted.
- No HTML templates or pages created yet.
- `ai/ToDoList.md`: Tasks & steps with checkpoint rows; all currently unstarted ("-"), including Task 1 steps (scaffolding partially present but not complete/integrated).

## 6. Next Needed Actions (When Full Repo Open)
1. Locate Control Service implementation (created via `services.NewControlService`) and extend its router to:
   - Mount `/install/state` using install.Handler(sessionStore, version, install.Configured()).
   - Intercept `/` to redirect to `/install/welcome` when !Configured().
2. Implement session ID generation (short random hex) set on first state call or welcome POST.
3. Deepen `Configured()` to ensure mainnet RPC presence (parse config). 
4. Create minimal placeholder `/install/welcome` HTML (static) before full templates (Task 3) to validate redirect flow.
5. Mark Task 1 steps complete in `ai/ToDoList.md` once wired.

## 7. Risks / Watch Items
- Need clear extension point for HTTP routes; if none exists, may introduce shim or modify SDK package (coordinate with maintainers).
- Avoid premature complex validation; adhere to “minimal first pass” principle.
- Keep apply atomic write semantics consistent across draft and final config.

## 8. Definition of Done (Initial Milestone)
Browser hitting control service root on fresh install serves (or redirects to) `/install/welcome`; `/install/state` returns `configured:false` with currentStep; creating `config.yaml` through (future) apply process flips state to `configured:true` and root now serves dashboard.

## 9. Reference Files
- `design/InstallFlow.md`
- `design/Dashboard.md`
- `design/Embedding.md`
- `ai/ToDoList.md`

---
This summary should be pasted or retained when reopening the larger mono-repo so implementation can continue seamlessly.
