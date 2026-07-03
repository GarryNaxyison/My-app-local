# Neriva Landing SEO Split Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Shorten the Neriva public landing, move removed long comparison/explanation content into generated SEO pages, and keep landing translations creative across 35 interface languages while static SEO/AEO stays RU/EN-only.

**Architecture:** Keep `site-react/src/EnglishSparkLanding.tsx` as the interactive main landing and keep static SEO pages generated from `site-react/scripts/static-seo-pages-data.mjs`. Update Playwright tests first so the landing split is verified by behavior, not screenshots.

**Tech Stack:** React 19, Vite, TypeScript, static Node HTML generator, Playwright.

---

### Task 1: Red Test For Short Landing

**Files:**
- Modify: `site-react/e2e/public-site.spec.ts`

- [ ] Replace old expectations for the overloaded landing with expectations that `.comparison-card`, `.answer-card`, `.course-strip`, `.device-flow`, `.memory-loop-section`, and `.review-card` are absent from the main landing.
- [ ] Keep expectations for `.landing-hero`, `.hero-demo`, four `.module-card` items, three `.plan-card` items, `.community-section`, `.landing-faq`, and compact `.seo-guides` links.
- [ ] Assert `.payment-methods` contains `Telegram Stars` and `YooKassa/SBP`, and does not contain `TON`, `USDT`, `crypto`, `blockchain`, or `RollyPay`.
- [ ] Run: `npm --prefix site-react run e2e -- public-site.spec.ts -g "landing presents the approved English spark hero product site"`
- [ ] Expected before implementation: FAIL because old long sections still exist and the new community/FAQ sections do not.

### Task 2: Red Test For Requested SEO Slugs

**Files:**
- Modify: `site-react/e2e/static-seo-pages.spec.ts`

- [ ] Replace old static SEO page slugs with `/ai-tutor.html`, `/speaking-practice.html`, `/pronunciation.html`, `/photo-translation.html`, and `/telegram-language-bot.html`, plus `/en/` alternates.
- [ ] Keep metadata, FAQ, related links, sitemap, and landing guide checks.
- [ ] Run: `npm --prefix site-react run e2e -- static-seo-pages.spec.ts -g "serves Russian static SEO pages"`
- [ ] Expected before implementation: FAIL because the requested pages are not generated yet.

### Task 3: Implement Short Landing

**Files:**
- Modify: `site-react/src/EnglishSparkLanding.tsx`
- Modify: `site-react/src/englishSparkLanding.css`

- [ ] Remove route/cockpit/scenario/device/memory/review/comparison/answer blocks from the main render.
- [ ] Keep hero, live demo, one product-scenario block, four key modules, pricing, community, FAQ, guide links, and final CTA.
- [ ] Add `community-section` and `landing-faq` using concise product copy and existing social/channel URLs.
- [ ] Use Swiss/editorial CSS: light surfaces, black text, hairline rules, one blue accent, fewer glass cards.

### Task 4: Move Removed Content Into SEO Pages

**Files:**
- Modify: `site-react/scripts/static-seo-pages-data.mjs`
- Modify: `site-react/src/landingSeoContent.ts`

- [ ] Rename static SEO slugs and guide links to the requested route set.
- [ ] Expand each SEO page sections so the removed long explanations live there.
- [ ] Keep all pages useful and visible; do not hide SEO text.
- [ ] Keep canonical/hreflang relationships and sitemap generation from the shared page list.

### Task 5: Add Translation Guardrails

**Files:**
- Modify: `site-react/public/assets/site-i18n.js`
- Modify: `tools/rebuild_public_site_translations.mjs`
- Modify: `site-react/e2e/public-site.spec.ts`
- Modify: `site-react/e2e/static-seo-pages.spec.ts`

- [ ] Add `englishSparkCreativeCopy` with maintained landing copy overrides for all 35 interface languages.
- [ ] Exclude `site-react/src/EnglishSparkLanding.tsx` from generated machine translation collection.
- [ ] Assert production `site-i18n.js` contains the curated landing source strings and `site-phrases.js` generated blocks do not.
- [ ] Assert sitemap/static SEO pages expose only Russian root pages and English `/en/` alternates.
- [ ] Assert non-RU/EN static SEO paths return 404.

### Task 6: Verify

**Files:**
- Generated: `site-react/public/*.html`
- Generated: `site-react/public/en/*.html`
- Generated: `site-react/public/sitemap.xml`

- [ ] Run: `npm --prefix site-react run build`
- [ ] Run focused tests: `npm --prefix site-react run e2e -- public-site.spec.ts static-seo-pages.spec.ts public-seo.spec.ts`
- [ ] Inspect the landing in browser at `http://127.0.0.1:5175/poliglot-ai.html?lang=en` or a preview build if a server is needed.
- [ ] Commit and push verified changes.
