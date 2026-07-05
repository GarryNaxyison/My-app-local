# Mobile Web App Repair Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Repair mobile web app usability without changing desktop behavior, and safely reset corrupted learner state while preserving account and commercial data.

**Architecture:** Mobile UI changes stay behind existing `@media (max-width: 760px)` rules and mobile-specific selectors in `web-react/src/styles/app.css`, with minimal React markup changes in `web-react/src/App.tsx`. Backend changes add explicit tests for free word audio, monotonic rating points, and a learner-state reset helper that avoids deleting users/accounts/payments/referrals.

**Tech Stack:** React 19, Vite, Playwright, Go, SQLite.

---

### Task 1: Mobile UI Regression Tests

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [x] Extend the existing mobile lesson test to assert the lesson textarea is large enough, has no fake internal scroll for short text, stays above `.mobile-bottom-nav-v2`, and the send SVG is visible.
- [x] Extend the existing mobile practice/tools test with the same textarea checks for practice and tools, plus `.tools-submit-v2 > svg` visibility.
- [x] Extend the mobile leaderboard bottom-menu guard to assert `window.scrollY === 0` after container scrolling and `document.documentElement.scrollWidth <= window.innerWidth + 1`.
- [x] Run focused Playwright tests and confirm they fail before CSS/markup fixes.

### Task 2: Backend Regression Tests

**Files:**
- Modify: `tools_test.go`
- Modify: `sqlite_store_test.go`
- Modify: `web_api_feature_test.go`

- [x] Add a test proving visible leaderboard score never drops below XP after XP/word progress changes.
- [x] Keep/extend word pronunciation tests so `/api/words/pronunciation` remains free for learned words.
- [x] Add a SQLite reset test proving learner state is wiped while premium/account/referral fields survive.
- [x] Run focused Go tests and confirm new tests fail before implementation where behavior is missing.

### Task 3: Mobile UI Implementation

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`

- [x] Replace premium payment modal close glyph with a back icon and accessible back label.
- [x] Remove redundant mobile composer header copy for lesson/practice/tool input where placeholder already describes the task.
- [x] Set mobile textarea geometry to readable defaults: larger min-height, hidden/auto-resizing internal overflow for short text, no fake scroll when content fits.
- [x] Rework mobile layout constraints so `.context-display` scrolls, `body/window` does not, functional panels cannot be scrolled out of view, and frame borders span the viewport correctly.
- [x] Ensure send buttons render visible lucide icons on mobile.

### Task 4: Backend Implementation

**Files:**
- Modify: `storage.go`
- Modify: `sqlite_store.go`
- Optionally create: `tools/reset_learning_state.go` or a repo-local script only if needed for operational reset.

- [x] Make global visible rating points monotonic with XP by using a score floor of user XP.
- [x] Add `resetLearningState` on `sqliteStore` that clears learner tables and resets learning counters/mode/history without touching account/commercial fields.
- [x] Provide an operational command/script path to run the reset safely.

### Task 5: Verification and Sync

**Files:**
- Update seekstone note with results.

- [x] Run focused Go tests.
- [x] Run `npm --prefix web-react run build`.
- [x] Run focused Playwright mobile tests.
- [x] Review `git diff` for desktop selector bleed.
- [x] Commit and push to `https://github.com/GarryNaxyison/My-app-local.git` if verification passes.
