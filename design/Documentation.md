# Documentation Overhaul Plan (Deferred)

Purpose: After the new install wizard, control endpoints, dashboard, and embedding mechanics are implemented and validated, we will execute a full documentation rewrite ("the book", all project README files, and other Markdown), explicitly excluding anything under `./ai/` (design notes stay historical/log-style).

## 1. Trigger (Do NOT Start Yet)
Rewrite begins only after explicit maintainer confirmation that:
1. New HTML install flow is merged & shipped (feature flag or default behavior).
2. Control service redesign (metadata file + structured endpoints) is live.
3. Dashboard + `/dashboard/state` endpoint stable.
4. Embed mode (`?embed=1`) is exercised successfully inside a TrueBlocks/Wails prototype.

## 2. Scope of Rewrite
- `book/src/` entire content: reorganize chapters around: Introduction, Local‑First Principles, Installation (HTML wizard), Configuration Reference, Services, Dashboard Usage, Embedding Guide, Operations & Tuning, Upgrade & Migration, FAQ.
- Top-level `README.md`: concise positioning, quick start (install + run + open local URL), link to Book.
- Any nested package / module READMEs (synchronize terminology: “control service”, “install wizard”, “dashboard”).
- `docs/` illustrative assets: regenerate diagrams (Mermaid where suitable) + optional fresh screenshots from the implemented UI.

## 3. Out of Scope
- `./ai/` directory (kept as internal design history, not polished prose).
- Marketing / website materials (handled separately if needed).
- Non-English translations (defer until English baseline stabilized).

## 4. Content Strategy
- Single source of truth for configuration keys: auto-generate a table from Go `types.Config` (code→doc) to eliminate drift.
- Replace narrative “wizard step” screenshots with: (a) one composite annotated screenshot, (b) textual step list referencing stable endpoint paths.
- Emphasize local-first security considerations (loopback binding, absence of remote control, explicit non-goals) in an early chapter.
- Provide an Embedding Quick Start (see Embedding Playbook) with minimal Wails integration example.

## 5. Process
1. Inventory existing pages; map to new structure (matrix: old file → new destination / retired).
2. Generate config reference stub via a small Go utility (iterate struct fields, emit Markdown).
3. Draft new chapter outlines (bulleted) → quick review → fill prose.
4. Insert updated diagrams (Mermaid + selective PNGs if clarity needed) with version tags.
5. Run link integrity & lint (markdownlint) pass, confirm zero warnings.
6. Final approval & merge; archive prior book commit hash in CHANGELOG for traceability.

## 6. Quality Gates for Rewrite
- No TODO/FIXME markers remain in published docs.
- Each endpoint referenced has either an example curl or JSON schema snippet.
- “Local-only” phrase appears in Introduction, Security, and Embedding chapters.
- Config key table auto-generated (checked in) + generation script committed.

## 7. Migration Notes Section (Future Chapter)
Will summarize differences vs. terminal wizard era: discovery, pause/unpause endpoints, absence of heuristic port scan, new metadata file, embedding behavior. (This historical delta appears only once; elsewhere we write present tense.)

## 8. Definition of Done (Docs Rewrite)
All scoped files updated, cross-links valid, config reference script added, maintainers sign-off, release notes point to new book section anchors.

---
Prepared; do not begin execution until triggers met.
