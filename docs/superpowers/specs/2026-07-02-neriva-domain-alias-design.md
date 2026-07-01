# Neriva Domain Alias Design

## Goal

Make `neriva.ru` the new production domain for the existing landing, web app, API, health check, payment callbacks, and deploy upload route while keeping `poliglotai.ru` and `poliglotai.online` working without redirects.

## Recommended Architecture

Treat `neriva.ru`, `www.neriva.ru`, and `api.neriva.ru` as additional hostnames for the current VPS deployment. Caddy should route the new hostnames through the same handlers as the existing Poliglot domains: static public site from `/var/www/poliglotai`, web app and API traffic to `localhost:8080`, and deploy upload to `localhost:19081`.

Runtime configuration should make Neriva the primary user-facing URL for new links and payment returns. Old domains remain accepted in CORS and Caddy so existing sessions, bookmarks, and provider callbacks continue to work during the rebrand.

## DNS Requirements

Create these records wherever `neriva.ru` DNS is managed:

- `neriva.ru` / apex: `A` record to `186.246.45.123`.
- `www.neriva.ru`: `CNAME` to `neriva.ru`, or an `A` record to `186.246.45.123`.
- `api.neriva.ru`: `A` record to `186.246.45.123`.
- Add `AAAA` records only after the VPS has confirmed working IPv6.

If `neriva.ru` is managed through Cloudflare like the existing domains, REG.RU should delegate to Cloudflare nameservers and the DNS records should be created in Cloudflare.

## Scope

In scope:

- Add Neriva hostnames to the deploy Caddyfile.
- Update sample/runtime environment documentation so Neriva is the primary app/payment URL and accepted CORS origin.
- Update Go production defaults, Telegram legal/app fallbacks, and browser legal-origin fallbacks so newly generated links prefer Neriva.
- Update public-site SEO, sitemap, robots, static page generation, and Cloudflare Worker defaults so Neriva is the canonical public host.
- Update deploy instructions and README examples to use Neriva endpoints.
- Preserve all existing Poliglot domains without redirects.

Out of scope:

- Redirecting `poliglotai.ru` or `poliglotai.online` to Neriva.
- Rebranding product copy, Telegram bot username, assets, or legal company text.
- Deploying to the VPS from this local change without explicit production credentials.
- Editing ignored local `.env` secrets into git.

## Acceptance Criteria

- `deploy/caddy/Caddyfile.updated` includes `neriva.ru`, `www.neriva.ru`, and `api.neriva.ru`.
- Existing `poliglotai.ru` and `poliglotai.online` hostnames remain present.
- `.env.example` documents Neriva as the primary `WEB_APP_URL`, `WEB_PAYMENT_RETURN_URL`, and `YOOKASSA_RETURN_URL`.
- Runtime defaults and legal links prefer Neriva, while old hosts remain accepted where aliases are needed.
- Generated SEO pages, sitemap, robots, and Cloudflare Worker defaults use Neriva as the canonical public host.
- Documentation uses Neriva for deploy/upload, webhook, and verification examples.
- Verification confirms no `redir` rule was added for old domains.
