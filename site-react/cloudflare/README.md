# Cloudflare Pages And Worker Setup

## Pages

Use Cloudflare Pages for the public landing.

- Project root: `site-react`
- Build command: `npm run build:cloudflare`
- Build output directory: `dist`
- Production domains: `neriva.ru`, `www.neriva.ru`
- Legacy public domains redirect to `https://neriva.ru{uri}`: `poliglotai.ru`, `www.poliglotai.ru`, `poliglotai.online`, `www.poliglotai.online`

The build keeps the current VPS deploy path untouched. It only changes output when `PUBLIC_SITE_OUT_DIR=dist` is set by `scripts/build-cloudflare.mjs`.

## Worker Fallback

Use the Worker only on dynamic routes, not on the whole public site.

Recommended routes:

- `neriva.ru/app*`
- `neriva.ru/login*`
- `neriva.ru/api/*`
- `neriva.ru/healthz`
- repeat the same app/API/health routes for `www.neriva.ru`
- route legacy public hosts to the Worker redirect before the VPS origin
- keep payment success and webhook routes on legacy API callback hosts until the payment providers are explicitly reconfigured

Worker variables:

- `ORIGIN_BASE_URL=https://api.neriva.ru`
- `MAINTENANCE_PAGE_URL=https://neriva.ru/maintenance.html`

Point `api.neriva.ru` to the VPS. Keep direct origin hostnames protected and do not use them as the public landing domain.
