# English Spark Hero Landing Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the approved English-first landing redesign with the previous full-screen Sparkles hero, burger drawer navigation, skeleton placeholders, fixed pricing, and readable dark/mobile layouts.

**Architecture:** Replace the temporary `BoldProductLanding` rendering path with a new focused landing component that reuses `SparklesCore` for the first viewport and keeps pricing/navigation concerns testable via stable class names. Keep translation generation out of this iteration; English landing copy is authoritative until design approval.

**Tech Stack:** React, TypeScript, Framer Motion, lucide-react, Vite, Playwright.

---

## File Structure

- Modify `site-react/e2e/public-site.spec.ts`: update RED tests from the rejected `.bold-hero` contract to the approved `.spark-hero`, burger drawer, skeleton, pricing, English-copy, contrast, and no-overlap contract.
- Create `site-react/src/EnglishSparkLanding.tsx`: English landing sections, Sparkles hero, skeleton product preview, pricing cards, and stable CTA/data attributes.
- Modify `site-react/src/PublicSiteApp.tsx`: render `EnglishSparkLanding` on the landing page and adjust `SiteNav` to burger/drawer behavior while preserving legal pages.
- Create `site-react/src/englishSparkLanding.css`: scoped landing styles, full-screen hero, drawer, skeletons, pricing, dark readability, responsive layout.
- Modify `site-react/src/main.tsx`: import the new CSS and remove the rejected bold landing CSS import if unused.
- Modify generated static files under `Сайт полиглота для бота/` via `npm --prefix site-react run build`.

## Task 1: RED Contract Tests

- [ ] Update `site-react/e2e/public-site.spec.ts` to assert:
  - landing renders `.english-spark-landing`;
  - `.spark-hero` exists, `.bold-hero` does not;
  - hero height is at least 90% of viewport on desktop/mobile;
  - burger button opens `.nav-drawer`;
  - nav is not sticky/fixed;
  - skeleton elements exist;
  - old prices use `s`/`del` and computed line-through;
  - landing text has no Cyrillic in English mode;
  - mobile has no overflow/overlap.
- [ ] Run `npm --prefix site-react run e2e -- --project=desktop-chromium`.
- [ ] Confirm expected failure references missing `.english-spark-landing` or current `.bold-hero` contract.

## Task 2: Landing Component

- [ ] Create `site-react/src/EnglishSparkLanding.tsx`.
- [ ] Implement arrays for modules, pricing, reviews, and proof chips in English.
- [ ] Build a full-screen `.spark-hero` using `SparklesCore`, large H1, CTA buttons with `data-entry`, proof chips, skeleton preview, and product card.
- [ ] Add sections for daily loop, modules, entry routes, pricing, reviews, and final CTA.

## Task 3: Navigation Drawer

- [ ] Modify `SiteNav` in `site-react/src/PublicSiteApp.tsx` to render a burger button and drawer on public pages.
- [ ] Drawer includes Features, Pricing, Reviews, Privacy, Terms, language select, theme toggle, Web app, Telegram.
- [ ] Keep legal page theme and language behavior working.

## Task 4: Styles

- [ ] Create `site-react/src/englishSparkLanding.css`.
- [ ] Set `.spark-hero` to full viewport with old visual direction and stable responsive constraints.
- [ ] Style `.nav-drawer`, `.skeleton-line`, `.skeleton-card`, pricing old prices, dark theme text, and mobile layouts.
- [ ] Ensure no sticky/fixed header and no text overlap at `390x844`.

## Task 5: Build, Visual Verification, Commit

- [ ] Run `npm --prefix site-react run build`.
- [ ] Run `npm --prefix site-react run e2e -- --project=desktop-chromium`.
- [ ] Run Playwright visual metrics for desktop/mobile hero height, overflow, overlap, drawer visibility, old price line-through, and canvas/sparkles presence.
- [ ] Stage only related files, leaving unrelated `web-react`/`web` changes untouched.
- [ ] Commit with `feat: redesign landing around spark hero`.
- [ ] Push `codex/ai-tutor-rebuild-fix`.

## Self-Review

- Spec coverage: every confirmed scope item maps to Task 1 assertions and Tasks 2-4 implementation.
- Placeholder scan: no TBD/TODO placeholders.
- Type consistency: class names are stable across tests and implementation plan.
