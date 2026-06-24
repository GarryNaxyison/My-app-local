# Mobile UI Regressions Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Repair the reported mobile UI and behavior regressions in the web app.

**Architecture:** Keep changes close to the existing React monolith and Go web API. Add one API endpoint for single mistake deletion, then wire UI behavior and mobile layout fixes through existing components and CSS.

**Tech Stack:** Go web API, React 19, TypeScript, Vite, Playwright, Nx.

---

### Task 1: Regression Tests

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`
- Modify: `web_api_feature_test.go`

- [ ] Add Playwright tests for lesson top scroll, level reentry, mistake hint/delete/similar, tools layout, and audio containment.
- [ ] Add Go test for `POST /api/mistakes/delete`.
- [ ] Run focused tests and confirm they fail before production changes.

### Task 2: API And React Behavior

**Files:**
- Modify: `web_api.go`
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/lib/types.ts`
- Modify: `web-react/src/components/ui/chat-interface.tsx`

- [ ] Register and implement `/api/mistakes/delete`.
- [ ] Add `deleteMistake` client action and per-card delete buttons.
- [ ] Track active mistake index and make "Practice similar" choose a different same-category mistake where possible.
- [ ] Stop chat panels from auto-scrolling to the bottom on mount.
- [ ] Preserve level tab state when returning to it.

### Task 3: Mobile Layout

**Files:**
- Modify: `web-react/src/styles/app.css`

- [ ] Make Premium/free copy wrap and stay inside plan cards.
- [ ] Make audio buttons and upload cards shrink/wrap inside mobile cards.
- [ ] Rework mobile tools layout so send/file/language controls do not overlap.
- [ ] Add compact mistake hint and card action styling.

### Task 4: Verification And Sync

**Files:**
- No new source files expected.

- [ ] Run focused Go and Playwright tests.
- [ ] Run `npm --prefix web-react run build`.
- [ ] Run `npm run test:go` or the closest passing project baseline.
- [ ] Commit and push to `origin`.
