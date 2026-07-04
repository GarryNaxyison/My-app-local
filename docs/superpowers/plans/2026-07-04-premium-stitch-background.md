# Premium Stitch Background Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the primitive grid landing background with the Stitch-style premium tech atmosphere and add a CSS-only animated border accent.

**Architecture:** Keep the existing React structure intact. Implement the background and animated contour as CSS-only layers in `site-react/src/englishSparkLanding.css`, and cover the behavior with Playwright checks in `site-react/e2e/landing-v3.spec.ts`.

**Tech Stack:** React, Vite, CSS, Playwright.

---

### Task 1: Premium Background And Animated Border

**Files:**
- Modify: `site-react/src/englishSparkLanding.css`
- Test: `site-react/e2e/landing-v3.spec.ts`

- [ ] **Step 1: Update visual CSS**

Replace the grid background on `.english-spark-landing` with layered midnight gradients, add ambient drift, and add `nerivaBorderOrbit` to key glass surfaces using `background-origin: border-box`.

- [ ] **Step 2: Add regression coverage**

In `landing-v3.spec.ts`, assert that the landing background no longer contains linear grid lines and that a representative card has a running border animation.

- [ ] **Step 3: Verify**

Run:

```powershell
npm run build
npx playwright test e2e/landing-v3.spec.ts --project=desktop-chromium
node ../tools/check_encoding_artifacts.mjs
```

- [ ] **Step 4: Deploy and verify production**

Deploy the generated static bundle to `neriva.ru`, then inspect production styles and hero stability with Playwright.
