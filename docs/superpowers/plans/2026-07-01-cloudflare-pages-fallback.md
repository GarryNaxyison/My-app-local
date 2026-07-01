# Cloudflare Pages Fallback Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a Cloudflare Free-ready static hosting path for the React landing and branded 404/maintenance fallback pages.

**Architecture:** Keep the existing VPS deploy output intact. Add a Cloudflare-specific Vite build output at `site-react/dist`, static error pages, Cloudflare Pages metadata files, and a Worker template that can proxy dynamic paths to the VPS origin and serve maintenance HTML on origin failure.

**Tech Stack:** React, Vite, TypeScript, Playwright, Cloudflare Pages, Cloudflare Workers.

---

## File Structure

- Modify `site-react/package.json` to add `build:cloudflare`.
- Modify `site-react/vite.config.ts` to allow an environment-selected output directory and include `404.html` and `maintenance.html` as Vite inputs.
- Create `site-react/scripts/build-cloudflare.mjs` to run the Cloudflare build cross-platform.
- Create `site-react/404.html` and `site-react/maintenance.html` as Vite entry HTML files.
- Create `site-react/src/ErrorPageApp.tsx` and `site-react/src/errorPages.css` for shared branded error UI.
- Create `site-react/public/_headers` and `site-react/public/_redirects` for Pages static behavior.
- Create `site-react/cloudflare/worker.js` and `site-react/cloudflare/README.md` as the Worker fallback template.
- Modify `site-react/e2e/deploy-output.spec.ts` and create `site-react/e2e/cloudflare-output.spec.ts`.

## Tasks

### Task 1: Cloudflare Build Output Tests

**Files:**
- Modify: `site-react/e2e/deploy-output.spec.ts`
- Create: `site-react/e2e/cloudflare-output.spec.ts`

- [ ] Add failing assertions that the normal build output contains `404.html` and `maintenance.html`.
- [ ] Add a failing Cloudflare output test that expects `site-react/dist` to contain `poliglot-ai.html`, `404.html`, `maintenance.html`, `_headers`, `_redirects`, `assets/site-react`, and SEO OG images.
- [ ] Run `npm --prefix site-react run e2e -- e2e/deploy-output.spec.ts e2e/cloudflare-output.spec.ts` and confirm the new checks fail because files/scripts do not exist yet.

### Task 2: Static Error Pages

**Files:**
- Create: `site-react/404.html`
- Create: `site-react/maintenance.html`
- Create: `site-react/src/ErrorPageApp.tsx`
- Create: `site-react/src/errorPages.css`
- Modify: `site-react/src/main.tsx`
- Modify: `site-react/vite.config.ts`

- [ ] Add React rendering for `404` and `maintenance` page modes based on `data-page-kind`.
- [ ] Add compact, production-oriented Russian-first copy and English alternates.
- [ ] Add Vite inputs for `404.html` and `maintenance.html`.
- [ ] Run the focused deploy output test and confirm it passes.

### Task 3: Cloudflare Pages Build Mode

**Files:**
- Modify: `site-react/package.json`
- Create: `site-react/scripts/build-cloudflare.mjs`
- Modify: `site-react/vite.config.ts`
- Create: `site-react/public/_headers`
- Create: `site-react/public/_redirects`

- [ ] Add `build:cloudflare` script.
- [ ] Make Vite use `site-react/dist` only when `PUBLIC_SITE_OUT_DIR=dist`.
- [ ] Add Pages `_headers` for no-cache HTML and immutable hashed assets.
- [ ] Add Pages `_redirects` for root landing fallback without touching dynamic app/API paths.
- [ ] Run `npm --prefix site-react run build:cloudflare` and confirm `site-react/dist` contains required files.

### Task 4: Worker Fallback Template

**Files:**
- Create: `site-react/cloudflare/worker.js`
- Create: `site-react/cloudflare/README.md`
- Create: `site-react/e2e/cloudflare-worker.spec.ts`

- [ ] Add a Worker module that proxies `/app`, `/login`, `/api`, `/healthz`, and callback paths to `ORIGIN_BASE_URL`.
- [ ] Serve `maintenance.html` when origin fetch throws or returns a 502/503/504 status.
- [ ] Add unit-style Playwright/Node test coverage for path matching and fallback response.
- [ ] Document Cloudflare dashboard variables and route setup.

### Task 5: Visual And Full Verification

**Files:**
- No production file changes expected.

- [ ] Run `npm --prefix site-react run build`.
- [ ] Run `npm --prefix site-react run build:cloudflare`.
- [ ] Run `npm --prefix site-react run e2e -- e2e/public-seo.spec.ts e2e/deploy-output.spec.ts e2e/cloudflare-output.spec.ts`.
- [ ] Start preview and capture desktop/mobile screenshots of `/404.html` and `/maintenance.html`.
- [ ] Review final diff, then request an independent code review.
- [ ] Commit and push verified changes to `origin/codex/ai-tutor-rebuild-fix`.
