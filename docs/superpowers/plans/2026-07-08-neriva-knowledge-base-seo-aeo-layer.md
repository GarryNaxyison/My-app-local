# NERIVA Knowledge Base SEO/AEO Layer Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add systematic SEO/AEO and conversion upgrades to the NERIVA Knowledge Base through category pillar pages, topic-specific CTAs, and in-body internal links.

**Architecture:** Keep the static generator as the source of truth. Add category route metadata to `knowledge-base-data.mjs`, render pillar pages and article link blocks from `generate-knowledge-base-pages.mjs`, and style the new blocks in the existing Knowledge Base stylesheet.

**Tech Stack:** Node ESM static generator, Vite public-site build, Playwright e2e, Caddy static hosting.

---

### Task 1: Failing E2E Coverage

**Files:**
- Modify: `site-react/e2e/knowledge-base.spec.ts`

- [x] Add a test that opens `/knowledge/vocabulary.html`, expects a successful response, title text, category article cards, and no non-vocabulary cards.
- [x] Add a test that opens an article and expects `.article-route-links`, a category pillar link, three article links, and `.article-cta--vocabulary`.
- [x] Extend sitemap coverage to require `/knowledge/start.html`, `/knowledge/vocabulary.html`, `/knowledge/english.html`, `/knowledge/speaking.html`, `/knowledge/exam.html`, and `/knowledge/technology.html`.
- [x] Run `npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts`.
- [x] Expected result before implementation: fail because pillar pages and the new article route/CTA classes do not exist.

### Task 2: Category Route Data

**Files:**
- Modify: `site-react/scripts/knowledge-base-data.mjs`

- [x] Add `knowledgeCategoryRoutes` with one route object for each existing category.
- [x] Add `categoryRouteFor(categoryId)` for article rendering.
- [x] Add `articlesForCategory(categoryId)` for pillar rendering.
- [x] Update `knowledgeSitemapUrls()` to include hub, category route paths, and article paths.
- [x] Keep existing article validation and add validation that every category has one route and every route path is unique.

### Task 3: Static Rendering

**Files:**
- Modify: `site-react/scripts/generate-knowledge-base-pages.mjs`

- [x] Import `knowledgeCategoryRoutes`, `categoryRouteFor`, and `articlesForCategory`.
- [x] Render category pillar pages with `CollectionPage`, `ItemList`, and `BreadcrumbList` JSON-LD.
- [x] Add `renderArticleRouteLinks(article)` after the direct-answer block.
- [x] Replace the generic CTA markup with `renderArticleCta(article)` using category-specific route CTA copy.
- [x] Generate `/knowledge/<category>.html` files before article pages.

### Task 4: Styling

**Files:**
- Modify: `site-react/public/assets/knowledge-base.css`

- [x] Add styles for `.knowledge-cluster`, `.cluster-route`, `.cluster-cta`, and `.article-route-links`.
- [x] Extend `.article-cta` with category-safe modifier classes without changing the established NERIVA palette.
- [x] Add mobile rules so route links and pillar pages do not cause horizontal overflow.

### Task 5: Build, Verify, Commit, Deploy

**Files:**
- Generated: `site-react/public/knowledge/*.html`
- Generated: `site-react/public/sitemap.xml`
- Generated: `Сайт полиглота для бота/**`

- [x] Run `npm --prefix site-react run build`.
- [x] Run `npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts`.
- [x] Run `node tools/check_encoding_artifacts.mjs`.
- [x] Check generated counts: 111 Knowledge Base HTML files in public output and 111 Knowledge Base URLs in sitemap.
- [ ] Stage only task files and generated artifacts; do not touch unrelated local dirty files.
- [ ] Commit and push to `origin/codex/ai-tutor-rebuild-fix`.
- [ ] Deploy the public-site bundle using the existing deploy-upload site-only flow.
- [ ] Verify production `/knowledge/vocabulary.html`, article CTA/internal links, sitemap, and read-only health URLs.
