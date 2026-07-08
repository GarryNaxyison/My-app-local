# Knowledge Search Navigation SEO Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add live topic suggestions and faster topic navigation to the NERIVA knowledge hub, and tighten root landing SEO for "изучение языков".

**Architecture:** Keep the knowledge base as generated static HTML with one small inline script. Reuse `knowledgeArticles`, `knowledgeCategoryRoutes`, and existing `data-search` card indexes. Reuse `landingSeoCopy` so static and runtime SEO metadata stay aligned.

**Tech Stack:** Node ESM generator, static HTML/CSS, Playwright e2e, Vite/React landing metadata.

---

### Task 1: RED Tests

**Files:**
- Modify: `site-react/e2e/knowledge-base.spec.ts`
- Modify: `site-react/e2e/public-seo.spec.ts`

- [ ] Add Playwright assertions that `/knowledge/` renders `.knowledge-topic-nav`, `.knowledge-search__suggestions`, and no suggestions before typing.
- [ ] Add Playwright assertions that typing `пам` shows suggestion links including `Как работает память при изучении языков?`, and clicking it opens `/knowledge/how-language-memory-works.html`.
- [ ] Add Playwright assertions that the quick topic nav exposes all 10 route links and `/knowledge/vocabulary.html`.
- [ ] Update Russian landing SEO constants so tests expect title/description with `изучение языков`.
- [ ] Run targeted tests and confirm they fail because the new UI/metadata does not exist yet:

```bash
cd site-react
npx playwright test e2e/knowledge-base.spec.ts e2e/public-seo.spec.ts --grep "knowledge hub|search|Russian social metadata|Russian canonical"
```

### Task 2: Knowledge Hub Implementation

**Files:**
- Modify: `site-react/scripts/generate-knowledge-base-pages.mjs`
- Modify: `site-react/public/assets/knowledge-base.css`

- [ ] Add a generated suggestions list under the search input using article data.
- [ ] Add a generated quick topic navigation block above the route/card sections.
- [ ] Extend `renderHubScript()` so suggestions update from the same normalized search logic as card filtering.
- [ ] Add CSS for desktop and mobile dropdown/index states without changing existing card markup.
- [ ] Regenerate static pages:

```bash
cd site-react
npm run seo:generate
```

### Task 3: Landing SEO Implementation

**Files:**
- Modify: `site-react/src/landingSeoContent.ts`
- Modify generated/static output through build or SEO generation as needed.

- [ ] Update Russian title and description in `landingSeoCopy.ru`.
- [ ] Keep English metadata semantically equivalent but do not expand static SEO pages beyond RU/EN.
- [ ] Run the targeted SEO tests and fix any constant mismatches.

### Task 4: Verification

**Files:**
- Verify generated outputs under `site-react/public`.

- [ ] Run:

```bash
cd site-react
npx playwright test e2e/knowledge-base.spec.ts e2e/public-seo.spec.ts
npm run build
```

- [ ] Inspect `git diff --stat` and confirm changes are scoped to knowledge UX, landing SEO metadata, generated public files, and docs.
