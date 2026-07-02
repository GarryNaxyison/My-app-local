# NERIVA Full Rebrand, Edge Redirect, And Bot Transition Design

## Goal

Make NERIVA the primary public product identity across the website, SEO output, web app links, Telegram-facing copy, and operational defaults. `neriva.ru` becomes the canonical public domain. `poliglotai.ru` and `poliglotai.online` become redirect-only legacy domains. The old Telegram bot becomes a transition bot that sends users to the new `@NERIVAapp_bot`.

## Recommended Architecture

Use `neriva.ru` as the public front door and keep the VPS as the dynamic origin. Cloudflare Pages serves the public static site, while the Cloudflare Worker proxies dynamic paths such as `/app`, `/login`, `/api`, `/healthz`, and payment/webhook callbacks to `api.neriva.ru`. If the VPS origin fails, the Worker returns the static maintenance page instead of exposing a raw origin failure.

Domain migration should happen at the edge first. Requests for `poliglotai.ru`, `www.poliglotai.ru`, `poliglotai.online`, and `www.poliglotai.online` should permanently redirect to `https://neriva.ru{path}{query}` before reaching the VPS. Caddy should also contain the same public-domain redirect as a backup, but Caddy must not be the only redirect layer because it cannot work when the VPS is down.

Telegram migration uses two runtime modes from the same codebase:

- Primary mode runs the full NERIVA bot with the new bot token and `WEB_TELEGRAM_LOGIN_BOT=NERIVAapp_bot`.
- Transition mode runs the old bot token with a minimal handler that does not require AI/payment/database dependencies and only sends a button to `https://t.me/NERIVAapp_bot`.

The new Telegram token is a secret and must not be written to source, docs, tests, generated assets, or commit history.

## Scope

In scope:

- Replace user-facing `Poliglot AI`, `POLIGLOT`, and Russian `Полиглот AI` brand references with `NERIVA` in active product surfaces.
- Update landing pages, SEO metadata, Open Graph/Twitter metadata, generated static SEO pages, sitemap, robots, public legal pages, maintenance pages, and visible site copy.
- Update web app shell/legal fallbacks and Telegram-facing bot copy so new links prefer NERIVA.
- Update default runtime URLs to `neriva.ru` / `api.neriva.ru` where callbacks can be safely moved.
- Keep old API callback hosts accepted temporarily where payment providers may still call them, but do not present old public domains as canonical.
- Add edge redirect behavior for legacy public domains.
- Add Caddy backup redirects for legacy public domains.
- Add old-bot transition mode and NERIVA bot username configuration.
- Update tests and generated public outputs.
- Update operator docs for the new deployment shape without committing secrets.

Out of scope:

- Editing historical changelog entries only to rewrite history.
- Committing the Telegram bot token or any production secret.
- Changing legal company/operator identity unless separate business/legal text is provided.
- Deploying Cloudflare dashboard settings or DNS records without available credentials.
- Removing temporary support for old API callback hosts before provider settings are confirmed migrated.

## Domain And Availability Behavior

`neriva.ru` is the canonical public site. Cloudflare should serve it even when the VPS is down:

- Static public pages are served by Cloudflare Pages.
- Dynamic paths are proxied by the Worker to `https://api.neriva.ru`.
- Origin failures return `/maintenance.html` from the static site.

Legacy public domains redirect at the edge:

- `https://poliglotai.ru/path?x=1` -> `https://neriva.ru/path?x=1`
- `https://www.poliglotai.ru/path?x=1` -> `https://neriva.ru/path?x=1`
- `https://poliglotai.online/path?x=1` -> `https://neriva.ru/path?x=1`
- `https://www.poliglotai.online/path?x=1` -> `https://neriva.ru/path?x=1`

API aliases can remain temporarily accepted for callbacks and rollout safety:

- `api.poliglotai.ru`
- `api.poliglotai.online`

They should not be used in newly generated public links.

## Telegram Behavior

The new full bot uses the new token from production secrets and presents itself as NERIVA. Public links, web login links, referral links, premium/payment deep links, and site CTA links use `@NERIVAapp_bot`.

The old bot uses its old token in transition mode. It should answer `/start`, `/menu`, and ordinary messages with a short migration notice and an inline URL button:

- Button text: `Open NERIVA`
- Button URL: `https://t.me/NERIVAapp_bot`

Transition mode should avoid OpenRouter, payment provider, and storage requirements. This allows the old transition service to keep responding even if only its Telegram token and target username are configured.

## Testing Strategy

Use test-first changes for behavioral code:

- Go tests for config defaults, transition-mode validation, NERIVA bot username fallback, and transition keyboard/link output.
- Worker tests for legacy host redirects, path/query preservation, canonical host handling, and VPS-origin maintenance fallback.
- Playwright tests for public SEO metadata, static SEO pages, sitemap, robots, and visible NERIVA branding.
- Focused checks that active public outputs no longer contain user-facing `Poliglot AI` brand strings except explicitly allowed legacy redirect/API compatibility docs.

Generated assets should be regenerated from source scripts rather than hand-edited when practical.

## Acceptance Criteria

- `neriva.ru` is canonical in active SEO/meta/sitemap/robots/static pages.
- `poliglotai.ru` and `poliglotai.online` public hosts redirect to `neriva.ru` with path and query preserved.
- The Worker can serve static NERIVA pages and return maintenance content when the VPS origin fails.
- Caddy contains backup redirects for old public hosts and serves/proxies NERIVA hosts.
- New public Telegram links point to `https://t.me/NERIVAapp_bot`.
- The old bot can run in transition-only mode and show a button to `@NERIVAapp_bot`.
- The new Telegram token is not committed.
- Focused Go, public-site, Worker, and build verifications pass.
