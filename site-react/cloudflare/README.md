# Cloudflare Pages And Worker Setup

## Pages

Use Cloudflare Pages for the public landing.

- Project root: `site-react`
- Build command: `npm run build:cloudflare`
- Build output directory: `dist`
- Production domains: `neriva.ru`, `www.neriva.ru`
- Reserve Cloudflare routes: `cf.neriva.ru`, `fallback.neriva.ru`, `reserve.neriva.ru`
- Legacy public domains redirect to `https://neriva.ru{uri}`: `poliglotai.ru`, `www.poliglotai.ru`, `poliglotai.online`, `www.poliglotai.online`

The build keeps the current VPS deploy path untouched. It only changes output when `PUBLIC_SITE_OUT_DIR=dist` is set by `scripts/build-cloudflare.mjs`.

## Worker Fallback

Keep `neriva.ru` and `www.neriva.ru` DNS-only when Russian access is the priority. DNS-only traffic bypasses Cloudflare, so the Worker cannot automatically serve the same hostname when the VPS is down.

Use proxied Worker custom domains for reserve access:

- `cf.neriva.ru`
- `fallback.neriva.ru`
- `reserve.neriva.ru`

Those reserve hostnames are expected to serve static landing, SEO, legal, robots, sitemap, and maintenance pages from Cloudflare Pages even if the VPS is down. Dynamic `/app`, `/api`, payment, webhook, and deploy-upload paths still require the origin API; when the origin is unavailable, the Worker returns the maintenance page.

Use the Worker only on dynamic routes for DNS-only primary domains, not on the whole public site.

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
