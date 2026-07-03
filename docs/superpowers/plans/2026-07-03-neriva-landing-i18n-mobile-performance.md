# Neriva Landing I18n Mobile Performance Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the NERIVA landing fully localized, faster, theme-aware, and separately composed for mobile while preserving the existing hero animation.

**Architecture:** Move landing copy into a typed local content module and render from current language instead of relying on post-load DOM translation. Keep SEO/AEO pages RU/EN, lazy-load heavy language/legal scripts only where needed, and split CSS by real theme/mobile states.

**Tech Stack:** React, TypeScript, Vite, Playwright, static public output in `Сайт полиглота для бота`.

---

### Task 1: Landing Content Contract

**Files:**
- Create: `site-react/src/landingContent.ts`
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Test: `site-react/e2e/public-site.spec.ts`

- [ ] Add tests that switching to `ru`, `en`, `es`, and `de` changes visible landing copy and that every configured locale has a complete content object.
- [ ] Create typed landing content for 35 locales. RU and EN are authored directly; the remaining 33 use concise curated marketing copy, not browser machine translation at runtime.
- [ ] Render all landing strings, alt text, FAQ, plan labels, proof labels, and CTA copy from `landingContent.ts`.
- [ ] Keep SEO/AEO static pages RU/EN only.

### Task 2: Theme Contract

**Files:**
- Modify: `site-react/src/englishSparkLanding.css`
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Test: `site-react/e2e/public-site.spec.ts`

- [ ] Add tests that the theme toggle changes landing background/text/button colors and hero preset.
- [ ] Replace hard-coded white/black `!important` blocks with CSS variables for `html[data-site-theme="light"]` and `html[data-site-theme="dark"]`.
- [ ] Keep the hero animation in both themes, with light/dark colors tuned separately.

### Task 3: Mobile Landing Composition

**Files:**
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Modify: `site-react/src/englishSparkLanding.css`
- Test: `site-react/e2e/public-site.spec.ts`

- [ ] Add mobile viewport tests at `390x844` for no horizontal overflow, compact CTA buttons, and bounded first-screen height.
- [ ] Add a mobile-only concise product rail and compact CTA sections.
- [ ] Reduce mobile section padding, hide secondary desktop-only screenshot density, and prevent full-width stretched buttons from elongating the page.

### Task 4: Performance Cleanup

**Files:**
- Modify: `site-react/src/PublicSiteApp.tsx`
- Modify: `site-react/vite.config.ts`
- Modify: `site-react/public/*.html`
- Test: `site-react/e2e/public-site.spec.ts`

- [ ] Add tests that `/poliglot-ai.html` does not include `site-phrases.js` or `legal-documents-i18n.js` in the initial HTML.
- [ ] Load legal documents only on legal pages.
- [ ] Remove `site-phrases.js` from landing startup and use React content instead.
- [ ] Keep `three` in its own async/vendor chunk and do not remove the hero animation.

### Task 5: Final Build, Screenshots, Deploy

**Files:**
- Generated: `Сайт полиглота для бота/*`
- Generated screenshots: `tmp/landing-redesign-20260703/*`

- [ ] Run `npm --prefix site-react run build`.
- [ ] Run focused Playwright tests for public landing.
- [ ] Capture desktop and mobile screenshots in light and dark themes.
- [ ] Commit, push, and deploy through the existing `neriva.ru/__codex_deploy_upload` flow.
