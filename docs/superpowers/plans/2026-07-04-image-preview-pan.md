# Image Preview Pan Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add drag-to-pan behavior to zoomed public-site image previews and remove native preview scrollbars.

**Architecture:** Keep the behavior local to `EnglishSparkLanding.tsx`. Use React state for scale and pan, pointer events for dragging, and CSS custom properties for the image transform.

**Tech Stack:** React 19, TypeScript, CSS, Playwright.

---

### Task 1: Regression Test

**Files:**
- Modify: `site-react/e2e/landing-v3.spec.ts`

- [ ] **Step 1: Write the failing test**

Extend `landing v3 product screenshots open a focused preview` to zoom the preview, verify the frame has no scroll overflow, drag the image, and expect `data-preview-pan` to change from `0,0`.

- [ ] **Step 2: Run test to verify it fails**

Run: `npx playwright test e2e/landing-v3.spec.ts -g "landing v3 product screenshots open a focused preview" --project=desktop-chromium`

Expected: FAIL because the frame still has scroll overflow or the pan data remains `0,0`.

### Task 2: Preview Pan Implementation

**Files:**
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Modify: `site-react/src/englishSparkLanding.css`

- [ ] **Step 1: Implement minimal behavior**

Add `previewPan`, drag state refs, pointer handlers, and CSS variables `--preview-pan-x` / `--preview-pan-y`. Clamp the interaction to zoomed previews and reset pan when scale returns to `1`.

- [ ] **Step 2: Hide scrollbars by removing native overflow**

Change `.image-preview__frame` from `overflow: auto` to `overflow: hidden`, add drag cursors, and combine `translate(...) scale(...)` on the image.

- [ ] **Step 3: Run focused verification**

Run: `npm run build`

Run: `npx playwright test e2e/landing-v3.spec.ts -g "landing v3 product screenshots open a focused preview" --project=desktop-chromium`

Expected: both commands exit `0`.
