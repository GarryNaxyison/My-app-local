# Desktop Web App Stability Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the desktop web app use its full canvas reliably and prevent generated learning content from being contradicted or duplicated.

**Architecture:** Desktop geometry is repaired in the existing `web-react/src/styles/app.css` outside mobile media queries. Learning consistency is enforced at the Go API boundary before model calls, while phrase canonicalization is a small pure web helper reused by candidate and saved-state logic.

**Tech Stack:** React 19, TypeScript, CSS, Go, Playwright, Go testing.

---

### Task 1: Guard generated-learning consistency

**Files:**
- Modify: `web_api.go`
- Test: `web_api_feature_test.go`

- [ ] Add a failing test proving `ROLEPLAY_TOOL_V2` uses `roleplayPrompt` instead of `practicePrompt`.
- [ ] Add a failing test proving a model phrase repeated with case/spacing/terminal-punctuation changes returns no mistakes.
- [ ] Add minimal server helpers for prompt routing and normalized model-answer equivalence; use them from `/api/practice` and `/api/lesson/answer`.
- [ ] Run `go test ./... -run 'Roleplay|Lesson'` and confirm it passes.

### Task 2: Deduplicate phrasebook suggestions

**Files:**
- Create: `web-react/src/lib/phrasebook.ts`
- Modify: `web-react/src/App.tsx`
- Test: `web-react/e2e/desktop-regressions.spec.ts`

- [ ] Add a failing test for `Hello`/`Hello!` canonical equivalence.
- [ ] Implement canonical phrase keys and use them for candidates, saved state, and replacement saves.
- [ ] Run the focused web regression test and confirm it passes.

### Task 3: Stabilize desktop-only geometry

**Files:**
- Modify: `web-react/src/styles/app.css`
- Modify: `web-react/src/components/ui/logout-button.tsx`
- Test: `web-react/e2e/desktop-regressions.spec.ts`

- [ ] Add a failing desktop viewport test for a full-width active workspace and a visible logout control.
- [ ] Add desktop CSS overrides outside `@media (max-width: 760px)` for full-width/scroll ownership, desktop workspaces, auth-card horizontal containment, and Turnstile containment.
- [ ] Give the logout button an immediate static visual base and reserved width; retain confirmation behavior.
- [ ] Run the desktop test at 2048×1050 and a mobile 760px smoke assertion.

### Task 4: Verify and ship

**Files:**
- Verify: `web-react/src/styles/app.css`, `web-react/src/App.tsx`, `web_api.go`

- [ ] Run `npm --prefix web-react run build`.
- [ ] Run the targeted Go and Playwright tests.
- [ ] Inspect the diff to verify there are no Android paths or mobile media-query changes.
- [ ] Commit the isolated branch and push `codex/desktop-webapp-fixes`.
