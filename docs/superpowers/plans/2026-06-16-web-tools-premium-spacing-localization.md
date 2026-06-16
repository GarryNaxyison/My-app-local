# Web Tools Premium Spacing Localization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the Tools workspace more compact vertically and finish Russian Free/Premium plan copy before deployment.

**Architecture:** This is a narrow React app slice. Tests live in `web-react/e2e/web-smoke.spec.ts`, product UI copy is owned by `web-react/src/lib/i18n.ts` plus fallback plan definitions in `web-react/src/App.tsx`, and layout is owned by `web-react/src/styles/app.css`.

**Tech Stack:** React, Vite, TypeScript, CSS, Playwright.

---

### Task 1: Regression Test

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [x] **Step 1: Write the failing test**

Add a test that opens `/app/?view=premium`, verifies Russian Free card copy includes "Для начала" and no English fallback fragments, then opens `/app/?view=tools` on desktop and asserts the composer panel starts close to the tool switcher.

- [x] **Step 2: Run the focused test to verify it fails**

Run: `npm --prefix web-react run e2e -- --grep "tools spacing and free plan Russian copy"`

Expected: FAIL before implementation because the Free plan still exposes old/English fragments or the Tools composer sits too low.

### Task 2: Copy And Layout Fix

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/lib/i18n.ts`
- Modify: `web-react/src/styles/app.css`

- [x] **Step 1: Update Russian plan copy**

Make the Free plan label "Для начала" and replace English fallback feature fragments with Russian copy. Keep "AI Tutor" and "Premium" as branded product names.

- [x] **Step 2: Tighten Tools top spacing**

Adjust `.tools-layout-v2` desktop layout so the chat panel and composer align closer to the tool switcher without changing button dimensions.

- [x] **Step 3: Run the focused test**

Run: `npm --prefix web-react run e2e -- --grep "tools spacing and free plan Russian copy"`

Expected: PASS.

### Task 3: Verification And Deploy

**Files:**
- Verify repository output and deployment artifacts only.

- [x] **Step 1: Run full checks**

Run:

```powershell
node tools/check_encoding_artifacts.mjs
npm --prefix web-react run build
npm --prefix web-react run e2e
```

- [x] **Step 2: Browser-check the UI**

Open the local web app in a browser, inspect `/app/?view=tools` and `/app/?view=premium`, and capture/confirm the corrected layout and copy.

- [ ] **Step 3: Commit and push**

Commit only files changed for this task and push to the configured GitHub remote.

- [ ] **Step 4: Deploy**

Use the repository deployment instructions for the React web app, then verify `https://poliglotai.ru/app/` and service health.
