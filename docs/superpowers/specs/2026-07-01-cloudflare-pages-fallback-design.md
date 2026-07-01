# Cloudflare Pages Fallback Design

## Goal

Host the public Poliglot AI landing as a static React/Vite site on Cloudflare Free while keeping the VPS backend as the origin for the web app, API, and payment callbacks.

## Recommended Architecture

Cloudflare Pages serves the public website for `poliglotai.ru` and `poliglotai.online`. The current VPS remains the origin for dynamic paths such as `/app`, `/login`, `/api`, `/healthz`, and payment/webhook callbacks.

The repository keeps the existing VPS build path unchanged and adds a Cloudflare-specific build path that writes the same static site into `site-react/dist`.

## Error And Fallback Behavior

The static site includes:

- `404.html` for missing public pages.
- `maintenance.html` for planned work or backend outage messaging.

Cloudflare Free cannot replace Cloudflare's native 522/5xx error pages through Custom Errors. To show branded maintenance content during origin failures, the project will add a Worker template that can proxy dynamic paths to the VPS origin and return the static maintenance page if the origin fetch fails.

## Scope

In scope:

- Add a Cloudflare Pages build command and output directory.
- Add static `404.html` and `maintenance.html` pages to the React/Vite build.
- Add Pages routing headers/redirects where useful for static hosting.
- Add a Worker template for dynamic-path fallback.
- Add Playwright/build-output tests for the new static files and Cloudflare build output.

Out of scope:

- Cloudflare dashboard setup, DNS changes, or production deployment without Cloudflare credentials.
- Moving the Go backend or Telegram bot to Cloudflare.
- Paid Cloudflare Custom Errors.

## Acceptance Criteria

- `npm --prefix site-react run build` still produces the existing VPS deploy output.
- `npm --prefix site-react run build:cloudflare` produces `site-react/dist`.
- `site-react/dist` contains `poliglot-ai.html`, SEO assets, `404.html`, `maintenance.html`, `_headers`, and `_redirects`.
- Public landing SEO tests still pass.
- Error pages are visually inspectable on desktop and mobile.
