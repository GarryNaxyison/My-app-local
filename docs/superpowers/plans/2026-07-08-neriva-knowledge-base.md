# NERIVA Knowledge Base Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build and deploy a Stitch-style `/knowledge/` section with 25 Russian articles, search, filters, article pages, sitemap entries, and a Knowledge Base link in the public header.

**Architecture:** Keep the Knowledge Base as generated static HTML under `site-react/public/knowledge/` so every article is indexable. Add a dedicated data file for article content and a generator that reuses the existing static SEO build flow. Keep the React landing change limited to the header link before the language selector.

**Tech Stack:** React 19 public header, Node ESM static generator, CSS, Playwright, Vite, existing VPS deploy flow.

---

### Task 1: Lock The Public Contract

**Files:**
- Create: `site-react/e2e/knowledge-base.spec.ts`

- [x] **Step 1: Write the failing Playwright test**

The test must assert `/knowledge/`, search/filter behavior, one article page, mobile article layout, and sitemap entries.

- [x] **Step 2: Run test to verify it fails**

Run:

```bash
npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts
```

Expected: FAIL because `/knowledge/` and article pages do not exist yet.

### Task 2: Add Knowledge Content And Generator

**Files:**
- Create: `site-react/scripts/knowledge-base-data.mjs`
- Create: `site-react/scripts/generate-knowledge-base-pages.mjs`
- Modify: `site-react/scripts/generate-static-seo-pages.mjs`

- [ ] **Step 1: Create article data**

Add 25 article records with `slug`, `title`, `category`, `cluster`, `readingTime`, `excerpt`, `keywords`, `tags`, `directAnswer`, `why`, `plan`, `mistakes`, `neriva`, and `faq`.

- [ ] **Step 2: Create generator**

Generate `/knowledge/index.html` and `/knowledge/<slug>.html` from the article data. The generator must write deterministic HTML, inline JSON data for search, JSON-LD for articles, and link `/assets/knowledge-base.css`.

- [ ] **Step 3: Wire sitemap**

Export a `knowledgeSitemapUrls()` helper and import it in `generate-static-seo-pages.mjs` so the existing `sitemap.xml` includes `/knowledge/` and all article URLs.

- [ ] **Step 4: Run the failing test**

Run:

```bash
npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts
```

Expected: Knowledge page assertions may still fail until CSS/header work lands, but routes should no longer be 404.

### Task 3: Add Stitch-Style CSS And Interactions

**Files:**
- Create: `site-react/public/assets/knowledge-base.css`
- Create generated inline script in `site-react/scripts/generate-knowledge-base-pages.mjs`

- [ ] **Step 1: Implement tokens and desktop layout**

Port the Stitch token values into local CSS: dark canvas, glass header, 80px nav, max-width 1280px, Sora/Inter/JetBrains Mono font families, large hub heading, search bar, filters, result cards, article answer callout, sticky TOC, article body, FAQ, related cards, and footer.

- [ ] **Step 2: Implement mobile layout**

Below 768px, use the mobile article composition: sticky compact header, search below brand row, collapsible TOC, compact related rows, no horizontal overflow, 44px touch targets.

- [ ] **Step 3: Implement search/filter JS**

The hub script must filter cards by text query and category, update result count, set active filter state, and restore all cards on reset.

- [ ] **Step 4: Run Knowledge Base tests**

Run:

```bash
npm --prefix site-react run e2e -- e2e/knowledge-base.spec.ts
```

Expected: PASS.

### Task 4: Add Header Link To The React Public Site

**Files:**
- Modify: `site-react/src/PublicSiteApp.tsx`
- Modify: `site-react/src/englishSparkLanding.css`

- [ ] **Step 1: Insert link before language selector**

In `SiteNavDrawer`, render a `nav-knowledge-link` anchor before `<select data-site-language-select>`. Use `База знаний` for `ru` and `Knowledge Base` otherwise.

- [ ] **Step 2: Style the link**

Match the current drawer header: pill border, muted text, primary hover, responsive short label on narrow screens.

- [ ] **Step 3: Run landing/public tests**

Run:

```bash
npm --prefix site-react run e2e -- e2e/public-seo.spec.ts e2e/static-seo-pages.spec.ts e2e/knowledge-base.spec.ts
```

Expected: PASS.

### Task 5: Build, Visual Verify, Commit, Push, Deploy

**Files:**
- Generated: `site-react/public/knowledge/index.html`
- Generated: `site-react/public/knowledge/*.html`
- Generated: `site-react/public/sitemap.xml`

- [ ] **Step 1: Build**

Run:

```bash
npm --prefix site-react run build
```

Expected: PASS and generated knowledge pages included in build output.

- [ ] **Step 2: Visual screenshots**

Use local Playwright screenshots for `/knowledge/`, the first article desktop, and the first article mobile. Compare against `site-react/tmp/knowledge-reference/*.png`.

- [ ] **Step 3: Commit and push**

Stage only files touched by this work, commit, and push to `https://github.com/GarryNaxyison/My-app-local.git`.

- [ ] **Step 4: Deploy**

Use the existing VPS deploy-upload/SSH flow from `deploy/SERVER_UPLOAD_PROMPT.md`. Verify:

```text
https://neriva.ru/knowledge/
https://neriva.ru/knowledge/can-you-learn-a-language-yourself.html
https://neriva.ru/sitemap.xml
https://neriva.ru/neriva.html
```

Expected: all return 200, production landing header contains `База знаний` before the language selector, and sitemap contains Knowledge Base URLs.
