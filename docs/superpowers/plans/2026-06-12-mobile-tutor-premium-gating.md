# Mobile Tutor And Premium Gating Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship mobile AI Tutor ordering, mobile overflow protection, and plan-tier messaging that gates AI Tutor to paid plans.

**Architecture:** Keep behavior client-side in the existing React app. CSS handles mobile ordering and overflow. `PremiumView` derives feature bullets from the existing `PremiumPlan` product/title data, so no backend migration is required.

**Tech Stack:** React 19, TypeScript, Vite, Tailwind CSS entry stylesheet, Playwright e2e.

---

### Task 1: Regression Tests

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [x] Add a mobile AI Tutor layout test inside the existing AI Tutor server-driven test after mobile navigation opens the tutor.
- [x] Add assertions that `.tutor-context-v2` top is above `.tutor-plan-v2` top.
- [x] Add `document.documentElement.scrollWidth <= window.innerWidth + 1`.
- [x] Add a premium tier-copy test that opens `/app/?view=premium` and checks Free excludes AI Tutor while Premium/Platinum include it.
- [x] Run `npm --prefix web-react run e2e -- --project mobile-chromium --grep "AI Tutor server-driven"` and confirm the new mobile ordering assertion fails before implementation.

### Task 2: AI Tutor Mobile Order And Overflow

**Files:**
- Modify: `web-react/src/styles/app.css`

- [x] At the mobile tutor breakpoint, set `.tutor-session-v2` to a single-column grid with bounded width.
- [x] Set `.tutor-context-v2 { order: 1; min-width: 0; max-width: 100%; }`.
- [x] Set `.tutor-plan-v2 { order: 2; min-width: 0; max-width: 100%; }`.
- [x] Add targeted wrap guards for tutor headers, cards, buttons, progress steps, and audio rows.
- [x] Rerun the mobile AI Tutor test and confirm it passes.

### Task 3: Premium Feature Ladder

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Modify: `web-react/src/lib/i18n.ts`

- [x] Add `premiumPlanFeatures(plan, copy)` near the existing premium plan helpers.
- [x] Render the feature list in each `.plan-card-v2`.
- [x] Update default and RU/EN plan copy to make AI Tutor paid-only and position Platinum as the higher-limit tier.
- [x] Add CSS for compact plan feature lists with locked/available states and mobile wrapping.
- [x] Run the premium tier-copy test and confirm it passes.

### Task 4: Verification And Sync

**Files:**
- Built output may update under `web/`.

- [x] Run `npm --prefix web-react run build`.
- [x] Run targeted Playwright tests for AI Tutor mobile and premium copy.
- [x] Run the project check target if feasible: `npm run check`.
- [x] Review `git diff`.
- [ ] Commit and push to `origin` from branch `codex/ai-tutor-rebuild-fix`.
