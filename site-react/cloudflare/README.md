# Cloudflare Pages And Worker Setup

## Pages

Use Cloudflare Pages for the public landing.

- Project root: `site-react`
- Build command: `npm run build:cloudflare`
- Build output directory: `dist`
- Production domains: `poliglotai.ru`, `www.poliglotai.ru`, `poliglotai.online`, `www.poliglotai.online`

The build keeps the current VPS deploy path untouched. It only changes output when `PUBLIC_SITE_OUT_DIR=dist` is set by `scripts/build-cloudflare.mjs`.

## Worker Fallback

Use the Worker only on dynamic routes, not on the whole public site.

Recommended routes:

- `poliglotai.ru/app*`
- `poliglotai.ru/login*`
- `poliglotai.ru/api/*`
- `poliglotai.ru/healthz`
- `poliglotai.ru/payment/success*`
- `poliglotai.ru/tonapi/webhook*`
- `poliglotai.ru/yookassa/webhook*`
- `poliglotai.ru/rollypay/webhook/*`
- repeat the same routes for `poliglotai.online`

Worker variables:

- `ORIGIN_BASE_URL=https://origin.poliglotai.ru`
- `MAINTENANCE_PAGE_URL=https://poliglotai.ru/maintenance.html`

Point `origin.poliglotai.ru` to the VPS. Keep it protected and do not use it as the public landing domain.
