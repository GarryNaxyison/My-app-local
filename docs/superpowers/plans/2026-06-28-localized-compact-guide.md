# Localized Compact Guide Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Localize the five-page web guide across all 35 interface languages and remove the visible modal header that consumes vertical space.

**Architecture:** Keep the guide component structure intact, but render the guide title and description only for assistive technology. Move all guide copy into an explicit locale map in the i18n layer so non-Russian languages no longer depend on generated placeholder strings.

**Tech Stack:** React 19, TypeScript, Vite, Playwright, Nx.

---

### Task 1: Failing Guide Tests

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add failing localization assertions**

Add assertions to the existing all-locale copy test that `app_guide_step_1_detail_1`, `app_guide_step_5_result`, and non-Russian `app_guide_body` are populated with meaningful guide text and do not collapse to the old placeholder shape.

- [ ] **Step 2: Add failing visible-header assertion**

Update `guide button opens localized quick start guide and mobile home shows level XP panel` so it expects no visible `h2` text `Как начать` and no visible dialog description text `Пять страниц про главные сценарии...`.

- [ ] **Step 3: Run the targeted tests and confirm RED**

Run: `cd web-react && npx playwright test e2e/web-smoke.spec.ts --grep "copy stays localized|guide button opens"`

Expected: the new assertions fail before production code changes.

### Task 2: Compact Guide Header

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`

- [ ] **Step 1: Keep accessible title/description only**

In `AppGuideDialog`, remove the visible `DialogTitle` and `DialogDescription` from the body flow and replace them with visually hidden equivalents.

- [ ] **Step 2: Tighten spacing**

Reduce the guide body top padding and gap so the modal starts with the compact icon/page row and then the active guide card.

- [ ] **Step 3: Run the guide UI test and confirm GREEN for header behavior**

Run: `cd web-react && npx playwright test e2e/web-smoke.spec.ts --grep "guide button opens"`

Expected: the guide still opens, has one active card, keeps `1/5`, and does not show the removed header copy.

### Task 3: Guide Localization

**Files:**
- Modify: `web-react/src/lib/i18n.ts`
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add explicit guide copy**

Add a `localizedGuideCopy` record keyed by `AppLocaleCode`, covering every guide key for all locales.

- [ ] **Step 2: Apply the map after generated placeholders**

Merge `localizedGuideCopy[code]` into `localeOverrides[code]` so explicit translations win over the old generated strings.

- [ ] **Step 3: Remove guide keys from fallback-leak allow list where explicit translations now exist**

Keep `pay_stars` and other intentional branded values allowed, but stop allowing guide keys to hide fallback leaks.

- [ ] **Step 4: Run localization tests and confirm GREEN**

Run: `cd web-react && npx playwright test e2e/web-smoke.spec.ts --grep "copy stays localized"`

Expected: all 35 locales have non-empty guide copy without mojibake, i18n keys, or old placeholder collapse.

### Task 4: Build, Review, Deploy

**Files:**
- Built output: `web/v2/index.html`
- Built output: `web/v2/assets`

- [ ] **Step 1: Run full web build**

Run: `npm run web:build`

Expected: TypeScript and Vite build finish with exit code 0.

- [ ] **Step 2: Inspect git diff**

Run: `git diff -- web-react/src/App.tsx web-react/src/lib/i18n.ts web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts docs/superpowers/specs/2026-06-28-localized-compact-guide-design.md docs/superpowers/plans/2026-06-28-localized-compact-guide.md`

Expected: only guide UI, guide copy, tests, and docs changed.

- [ ] **Step 3: Commit and push**

Run: `git add <changed files> && git commit -m "fix: localize compact app guide" && git push`

Expected: commit and push succeed on the current branch.

- [ ] **Step 4: Deploy**

Use the project deployment script/checklist discovered from `deploy/` and verify the running service or published static app after deployment.
