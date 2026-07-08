# NERIVA Knowledge Base 100 Articles Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use subagent-driven-development or executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Extend the deployed NERIVA Knowledge Base from 25 to 100 Russian articles.

**Architecture:** Keep static generation. Add reusable article factories and 75 data entries to `knowledge-base-data.mjs`; update tests and hub count to use the article array length.

**Tech Stack:** Node ESM static generator, Vite public site, Playwright e2e, Caddy static serving.

---

### Task 1: Red Tests

**Files:**
- Modify: `site-react/e2e/knowledge-base.spec.ts`

- [x] Update hub expectations from 25 to 100 cards and from 8 to 10 filters.
- [x] Add category filter checks for vocabulary, CEFR, exam, and technology.
- [x] Add article checks for `/knowledge/how-language-memory-works.html` and `/knowledge/how-to-choose-language-learning-app.html`.
- [x] Add sitemap checks for representative new URLs.
- [x] Run `npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts`; expected result before implementation: fail because the site still has 25 cards and missing new pages.

### Task 2: Content Data

**Files:**
- Modify: `site-react/scripts/knowledge-base-data.mjs`

- [x] Add `exam` and `technology` categories.
- [x] Add shared extended article helpers for category-specific defaults, FAQ, plan, mistakes, and NERIVA copy.
- [x] Add articles 26-90 from the DOCX list exactly.
- [x] Add articles 91-100 for exam preparation and AI/technology topics.
- [x] Validate that article numbers are 1-100, slugs are unique, paths are unique, categories exist, and every article has FAQ content.

### Task 3: Generator And Static Output

**Files:**
- Modify: `site-react/scripts/generate-knowledge-base-pages.mjs`
- Generated: `site-react/public/knowledge/*.html`
- Generated: `site-react/public/sitemap.xml`
- Generated: `Сайт полиглота для бота/**`

- [x] Replace hardcoded `25 материалов` with `knowledgeArticles.length`.
- [x] Run `npm --prefix site-react run seo:generate`.
- [x] Run `npm --prefix site-react run build`.

### Task 4: Verification

**Commands:**
- `npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts`
- `node tools/check_encoding_artifacts.mjs`

- [x] Capture visual smoke screenshots for hub desktop, representative article desktop, and representative article mobile.
- [x] Confirm no horizontal overflow on mobile.

### Task 5: Git And Deploy

**Files:**
- Stage only task files and generated public-site artifacts.

- [ ] Commit and push to `origin/codex/ai-tutor-rebuild-fix`.
- [ ] Build `poliglot-public-site.tgz`.
- [ ] Upload through deploy-upload service.
- [ ] Run site-only deploy.
- [ ] Verify production `/knowledge/`, new article URL, sitemap, and header link.
