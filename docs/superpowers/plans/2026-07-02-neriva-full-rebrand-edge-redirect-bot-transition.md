# NERIVA Full Rebrand, Edge Redirects, Bot Transition Plan

## Scope

- Make NERIVA the public brand across active landing, SEO, app, bot text, examples, and generated public assets.
- Keep `neriva.ru` as the canonical public site and app origin.
- Redirect `poliglotai.ru`, `www.poliglotai.ru`, `poliglotai.online`, and `www.poliglotai.online` to `https://neriva.ru{uri}`.
- Keep API legacy host aliases only for migration callbacks while production webhooks are moved to `api.neriva.ru`.
- Add a transition-only Telegram mode for the old bot token that replies with a button to `@NERIVAapp_bot`.
- Keep the Cloudflare Worker fallback so `neriva.ru` can serve static/maintenance content when the VPS is down.

## Verification

- Add RED tests for Cloudflare legacy host redirects.
- Add RED tests for NERIVA config defaults and transition-only bot configuration.
- Add RED tests for the transition bot keyboard URL and copy.
- Regenerate SEO pages and run public-site build.
- Run web app build after rebrand string changes.
- Run Go tests and focused Playwright Worker tests.
