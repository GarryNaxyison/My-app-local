# Neriva Domain Alias Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `neriva.ru` as a full production hostname set for the existing VPS landing, web app, API, callbacks, and deploy upload route without redirecting old Poliglot domains.

**Architecture:** Keep one shared Caddy routing setup and add Neriva hostnames to the existing site/API blocks. Make Neriva the documented and runtime primary URL in environment examples, Go defaults, legal links, SEO/static-site output, Cloudflare Worker defaults, and deploy docs while preserving old domains as accepted aliases.

**Tech Stack:** Caddy, Go runtime configuration, repository deployment docs, PowerShell/Git Bash verification.

---

## File Structure

- Modify `deploy/caddy/Caddyfile.updated` to add Neriva hosts to the production Caddy config.
- Modify `.env.example` to make Neriva the primary documented runtime URL and keep old domains in `WEB_CORS_ORIGINS`.
- Modify `config.go`, `telegram.go`, `web-react/src/App.tsx`, and focused Go tests for runtime/legal defaults.
- Modify `site-react` SEO sources, Cloudflare Worker defaults, e2e expectations, and generated public HTML/sitemap/robots output.
- Modify `README.md` payment/Caddy examples to reference Neriva.
- Modify `deploy/SERVER_UPLOAD_PROMPT.md`, `docs/technical/WEB_API.md`, and `docs/technical/DOCUMENTATION.md` deploy/webhook/API examples to reference Neriva while keeping the same VPS IP.

## Tasks

### Task 1: Caddy Hostname Alias

**Files:**
- Modify: `deploy/caddy/Caddyfile.updated`

- [ ] Add `neriva.ru` and `www.neriva.ru` to the main site block.
- [ ] Add `api.neriva.ru` to the API reverse proxy block.
- [ ] Keep `poliglotai.ru`, `www.poliglotai.ru`, `poliglotai.online`, `www.poliglotai.online`, `api.poliglotai.ru`, and `api.poliglotai.online`.
- [ ] Verify with text checks that no `redir` directive was introduced.

### Task 2: Runtime Environment Examples

**Files:**
- Modify: `.env.example`

- [ ] Keep `YOOKASSA_RETURN_URL=https://api.poliglotai.ru/payment/success`.
- [ ] Include Neriva and old Poliglot domains in `WEB_CORS_ORIGINS`.
- [ ] Set `WEB_APP_URL=https://neriva.ru/app`.
- [ ] Keep `WEB_PAYMENT_RETURN_URL=https://poliglotai.ru/app?payment=success`.
- [ ] Do not edit the ignored local `.env` file in git.

### Task 3: Operator Documentation

**Files:**
- Modify: `README.md`
- Modify: `deploy/SERVER_UPLOAD_PROMPT.md`
- Modify: `docs/technical/WEB_API.md`
- Modify: `docs/technical/DOCUMENTATION.md`

- [ ] Update payment return examples to Neriva.
- [ ] Update deploy upload endpoint examples to `https://neriva.ru/__codex_deploy_upload/...`.
- [ ] Update webhook examples to `api.neriva.ru`.
- [ ] Update verification URLs to `https://neriva.ru/...`.
- [ ] Keep the server IP `186.246.45.123` unchanged.

### Task 4: Runtime Defaults and Legal Links

**Files:**
- Modify: `config.go`
- Modify: `telegram.go`
- Modify: `web-react/src/App.tsx`
- Modify: `config_test.go`
- Modify: `premium_test.go`
- Modify: `web_app_shell_test.go`

- [ ] Add focused tests for Neriva production defaults and legal-origin handling.
- [ ] Make Go defaults prefer `https://neriva.ru/app`, `https://api.neriva.ru`, and old aliases in CORS.
- [ ] Make Telegram legal/app fallbacks and web React legal fallbacks prefer `https://neriva.ru`.
- [ ] Keep old host handling where current-host aliases remain supported.

### Task 5: Public Site SEO and Cloudflare Defaults

**Files:**
- Modify: `site-react/src/landingSeoContent.ts`
- Modify: `site-react/scripts/static-seo-pages-data.mjs`
- Modify: `site-react/scripts/generate-static-seo-pages.mjs`
- Modify: `site-react/cloudflare/worker.js`
- Modify: `site-react/cloudflare/README.md`
- Modify: `site-react/e2e/public-seo.spec.ts`
- Modify: `site-react/e2e/static-seo-pages.spec.ts`
- Modify: `site-react/e2e/cloudflare-worker.spec.ts`
- Modify generated public site HTML, sitemap, and robots outputs.

- [ ] Make Neriva the canonical public-site origin.
- [ ] Regenerate static SEO output.
- [ ] Update Cloudflare Worker defaults to `https://api.neriva.ru` and the Neriva Pages host.
- [ ] Keep old Poliglot domains documented as accepted aliases, not redirects.

### Task 6: Verification

**Files:**
- No production file changes expected.

- [ ] Run focused text checks for Neriva hostnames and old Poliglot hostnames.
- [ ] Run focused Go tests for Neriva config, RollyPay return URL, and legal links.
- [ ] Run focused Playwright e2e tests for public SEO, static SEO pages, and Cloudflare Worker defaults.
- [ ] Run site build to regenerate public output.
- [ ] Run Cloudflare build check.
- [ ] Inspect `git diff --check`.
- [ ] Review final diff before commit.
