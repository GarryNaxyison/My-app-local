# Web Roleplay Errors Ribbon Localization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Localize roleplay and mistake category UI across all 35 app locales and fix desktop ribbon arrow placement.

**Architecture:** Keep localization in `web-react/src/lib/i18n.ts` and consume it through `appCopy` from `web-react/src/App.tsx`. Keep ribbon behavior in the existing morphing button, adding a compact mode used only by the function ribbon.

**Tech Stack:** React, TypeScript, Vite, Playwright, Nx workspace scripts.

---

### Task 1: Localization Regression Tests

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Write failing tests**

Add assertions that every `appLocaleCodes` entry has non-empty localized values for `mistake_group_*`, `roleplay_title`, `roleplay_subtitle`, `roleplay_scenario_*`, and `roleplay_scenario_*_description`, and that Russian roleplay values do not contain English fallback fragments.

- [ ] **Step 2: Run test to verify it fails**

Run: `npm --prefix web-react run e2e -- --project=desktop-chromium --grep "localized for all 35|roleplay scenario"`

Expected: FAIL before localization dictionaries are added.

### Task 2: All-Locale Copy

**Files:**
- Modify: `web-react/src/lib/i18n.ts`
- Modify: `web-react/src/App.tsx`

- [ ] **Step 1: Implement copy dictionaries**

Add all-locale records for mistake category labels, roleplay UI labels, and scenario title/description keys.

- [ ] **Step 2: Use localized scenario keys**

Update `roleplayScenarioTitle` and `roleplayScenarioDescription` to prefer `appCopy(lang, key)` before older local fallbacks.

- [ ] **Step 3: Run localization tests**

Run: `npm --prefix web-react run e2e -- --project=desktop-chromium --grep "localized for all 35|roleplay scenario"`

Expected: PASS.

### Task 3: Ribbon Geometry

**Files:**
- Modify: `web-react/src/components/ui/morphing-arrow-button.tsx`
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add compact arrow mode**

Add `size?: "default" | "compact"` to `MorphingArrowButton`, use smaller dimensions in compact mode, and pass `size="compact"` from `FunctionRibbon`.

- [ ] **Step 2: Normalize ribbon CSS**

Use fixed grid columns for arrows, remove negative right-arrow margins, and keep arrows centered inside their columns.

- [ ] **Step 3: Run desktop geometry test**

Run: `npm --prefix web-react run e2e -- --project=desktop-chromium --grep "desktop ribbon right arrow"`

Expected: PASS with both arrows inside the app frame and the right arrow placed to the right of the ribbon content.

### Task 4: Verification And GitHub Sync

**Files:**
- Verify changed files only and commit.

- [ ] **Step 1: Build**

Run: `npm run web:build`

Expected: exit 0.

- [ ] **Step 2: Focused Playwright checks**

Run focused Playwright checks for localization and ribbon geometry.

- [ ] **Step 3: Commit and push**

Run `git add`, `git commit`, and `git push` to the configured GitHub remote.
