# Web UI Regression Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix the roleplay, pronunciation, notes, and tools UI regressions in the React web app.

**Architecture:** Keep the existing single-file React view structure and CSS system. Add targeted state plumbing for phrasebook saved state and roleplay voice submission, then solve layout regressions with scoped class changes.

**Tech Stack:** React 19, TypeScript, Vite, Playwright, Nx.

---

## File Map

- Modify `web-react/e2e/web-smoke.spec.ts`: add focused regression tests before production edits.
- Modify `web-react/src/App.tsx`: pass phrasebook state to views, add roleplay voice controls, reuse quick-save state, and expose phrase saves where missing.
- Modify `web-react/src/styles/app.css`: adjust roleplay, pronunciation, ribbon, tools, and quick-save styles.

### Task 1: Regression Tests

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [x] **Step 1: Add tests for quick-save active state, roleplay voice answer, mobile pronunciation, tools compactness, and desktop ribbon containment.**

Use existing app fixtures and selectors. Assertions must inspect real DOM dimensions and computed color/state, not snapshots.

- [x] **Step 2: Run the focused Playwright tests to verify they fail before implementation.**

Run: `npm --prefix web-react run e2e -- --grep "regression: .*notes|regression: roleplay accepts voice|regression: mobile pronunciation|regression: tools selector|regression: desktop function ribbon"`

Expected: at least the new tests fail for missing active bookmark state and missing roleplay voice control.

### Task 2: Notes State and Coverage

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`

- [x] **Step 1: Add phrasebook saved lookup helpers.**

Derive a normalized phrase set from `phrasebook` and pass it through `ViewRendererProps`.

- [x] **Step 2: Update `PhraseQuickSave` and word-result save button.**

Buttons receive active state from the saved lookup and set `aria-pressed`, an active class, and active bookmark color immediately after `savePhrase`.

- [x] **Step 3: Add quick-save entries in roleplay and tools.**

Roleplay exposes the latest roleplay phrases. Tools exposes text from the latest tool answer and current draft when meaningful.

- [x] **Step 4: Run focused notes test and keep it green.**

Run: `npm --prefix web-react run e2e -- --grep "notes"`

Expected: quick-save bookmark changes active state after click and phrasebook persistence still passes.

### Task 3: Roleplay Text and Voice Layout

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`

- [x] **Step 1: Extend roleplay submission to accept an optional voice file.**

Use existing `apiForm` + `voice` upload pattern and keep text submission unchanged.

- [x] **Step 2: Add `FileControls` to the roleplay composer.**

Show voice upload/recording only, clear the voice file after successful submit, and keep the composer pinned below the dialogue content.

- [x] **Step 3: Tighten desktop roleplay session sizing.**

Make the dialogue card span the content width with the scroll area above the composer.

- [x] **Step 4: Run focused roleplay tests.**

Run: `npm --prefix web-react run e2e -- --grep "roleplay"`

Expected: desktop and mobile roleplay layout tests pass, including voice submission.

### Task 4: Pronunciation and Tools Layout

**Files:**
- Modify: `web-react/src/styles/app.css`

- [x] **Step 1: Make mobile pronunciation a strict vertical document flow.**

Remove grid heights that force overlap; ensure each block has `min-height: 0`, visible overflow where needed, and sufficient bottom padding above mobile nav.

- [x] **Step 2: Compact tool selector buttons.**

Reduce selector min-height/padding, keep icon and label aligned, and let input/output use more width and height.

- [x] **Step 3: Keep desktop ribbon arrows inside their shell.**

Use fixed shell columns for arrows and avoid negative margins/transforms that move buttons outside the menu frame.

- [x] **Step 4: Run focused layout tests.**

Run: `npm --prefix web-react run e2e -- --grep "pronunciation|tools selector|function ribbon"`

Expected: no overlapping blocks, compact tool buttons, and contained ribbon arrows.

### Task 5: Verification and Sync

**Files:**
- No additional production files.

- [x] **Step 1: Build the web app.**

Run: `npm --prefix web-react run build`

Expected: TypeScript and Vite build exit 0.

- [x] **Step 2: Run targeted Playwright regressions.**

Run: `npm --prefix web-react run e2e -- --grep "roleplay|pronunciation|tools selector|function ribbon|notes|phrasebook"`

Expected: all targeted tests pass.

- [x] **Step 3: Open the app in the in-app browser and verify critical screens visually.**

Check desktop roleplay, desktop tools, mobile pronunciation, and mobile roleplay.

- [x] **Step 4: Commit and push.**

Run: `git add docs/superpowers/specs/2026-06-13-web-ui-regression-fixes-design.md docs/superpowers/plans/2026-06-13-web-ui-regression-fixes.md web-react/e2e/web-smoke.spec.ts web-react/src/App.tsx web-react/src/styles/app.css`

Run: `git commit -m "fix: stabilize web app layouts and notes"`

Run: `git push`
