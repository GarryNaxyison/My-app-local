# Landing SEO/AEO Cross-Domain Design

Date: 2026-06-29

## Summary

Poliglot AI needs the existing public landing page to become easier to discover in search and AI answer experiences without a visual redesign. The Russian domain `poliglotai.ru` is the primary surface for Russian-speaking users. The `poliglotai.online` mirror is the primary surface for international users.

This design adds a technical SEO layer and a small answer-oriented content layer to the current `site-react` public landing. It keeps the current layout, product positioning, language selector, relative `/app/` links, legal pages, and React/Vite build flow.

## Goals

- Make the current landing clearer to crawlers and answer engines for Russian queries first.
- Preserve `poliglotai.ru` as the Russian canonical host and `poliglotai.online` as the international canonical host.
- Add metadata and structured data that describe Poliglot AI as an AI language tutor, web app, and Telegram bot.
- Add a compact visible FAQ/AEO section with direct answers for common product-discovery queries.
- Add `robots.txt` and `sitemap.xml` entries for the public landing and legal pages.
- Cover the change with focused Playwright checks and existing build checks.

## Non-Goals

- No visual redesign of the landing.
- No new marketing funnel, blog, or separate SEO landing cluster in this iteration.
- No forced redirect from `.online` to `.ru` or the reverse.
- No paid SEO promises, keyword stuffing, hidden text, or mass-generated pages.
- No reliance on `llms.txt` or custom "AI-only" markup as a ranking mechanism.

## Current Context

- Public site source is in `site-react`.
- Main landing HTML entry is `site-react/poliglot-ai.html`.
- Public React entry is `site-react/src/PublicSiteApp.tsx`, currently rendering `EnglishSparkLanding`.
- The site supports `?lang=<code>` and persists the selected interface language with `poliglot_site_language`.
- The language script updates `document.documentElement.lang`, document title, translated text, translated attributes, and internal HTML links.
- Public app links are intentionally relative (`/app/`) so users stay on the current host: `.ru/app/` or `.online/app/`.
- The project already supports both domains in CORS, docs, Caddy examples, and legal text.

## Domain And Canonical Strategy

The two public domains must be treated as localized regional surfaces, not duplicate mirrors where one domain canonicalizes to the other.

- `https://poliglotai.ru/poliglot-ai.html` is canonical for Russian (`ru`).
- `https://poliglotai.online/poliglot-ai.html?lang=en` is canonical for English/international (`en`) when the page is requested in English.
- Other supported non-Russian interface languages use the `.online` host with their `?lang=<code>` URL.
- `x-default` points to `https://poliglotai.online/poliglot-ai.html?lang=en`.
- The Russian page links alternates for `ru`, `en`, and other supported locales.
- The international page links alternates for `ru` on `.ru` and supported non-Russian locales on `.online`.

Implementation should use static tags where safe and a small runtime metadata helper where the current host and `?lang=` must influence canonical and alternate URLs.

## Metadata

Add or update page metadata for the landing:

- `title`
  - Russian: `Poliglot AI - AI-репетитор английского и языков в Telegram`
  - English: `Poliglot AI - AI language tutor in Telegram and web app`
- `description`
  - Russian should mention AI lessons, speaking practice, Telegram, web app, pronunciation, photo translation, mistakes, and Premium.
  - English should mention AI language tutor, speaking practice, Telegram bot, web app, voice, photo translation, mistakes, and premium plans.
- `robots`: `index, follow`.
- `canonical`: host-aware as described above.
- `alternate` `hreflang`: all 35 supported interface languages plus `x-default`.
- Open Graph and Twitter tags:
  - title and description localized by current language.
  - type `website`.
  - URL set to canonical URL.
  - image set to an absolute public URL on the current canonical host.
- Keep `theme-color`, favicon, viewport, and existing script order intact.

The metadata helper must not translate or mutate the viewport meta tag.

## Structured Data

Add JSON-LD for the landing. The data should be truthful and compact.

Recommended graph:

- `Organization`
  - name `Poliglot AI`
  - url based on canonical host
  - logo URL
  - sameAs social URLs already used by the public site
  - contact email `supportpoliglotai@gmail.com`
- `WebSite`
  - name `Poliglot AI`
  - url based on canonical host
  - inLanguage from current language
- `SoftwareApplication`
  - name `Poliglot AI`
  - applicationCategory `EducationalApplication`
  - operatingSystem `Web, Telegram`
  - offers for Free, Premium, and Platinum with current public prices
  - description localized for `ru` and `en`
- `FAQPage`
  - question/answer pairs must match visible FAQ/AEO content on the page.

Use JSON-LD because Google recommends it for structured data. Do not mark up content that is not visible or is materially different from the page.

## Visible FAQ/AEO Content

Add a compact answer-oriented section to the existing landing, near the current FAQ/pricing/final CTA area where it fits without changing the hero. The section should contain direct answers that can stand alone in snippets and AI answers.

Russian priority questions:

- `Что такое Poliglot AI?`
- `Можно ли учить английский с ИИ в Telegram?`
- `Чем AI-репетитор отличается от обычного приложения со словами?`
- `Можно ли тренировать произношение и speaking?`
- `Можно ли переводить текст с фото?`
- `Есть ли бесплатный тариф?`

English/international equivalents:

- `What is Poliglot AI?`
- `Can I practice English with an AI tutor in Telegram?`
- `How is an AI tutor different from a vocabulary app?`
- `Can I practice pronunciation and speaking?`
- `Can I translate text from photos?`
- `Is there a free plan?`

The answers should be short, factual, and product-specific. They should not promise guaranteed language results, exam outcomes, or employment outcomes.

## Robots And Sitemap

Add public static search files under `site-react/public` so Vite copies them to the built public site.

`robots.txt`:

- Allow normal crawling.
- Reference both controlled sitemap URLs:
  - `https://poliglotai.ru/sitemap.xml`
  - `https://poliglotai.online/sitemap.xml`

`sitemap.xml`:

- Use one static sitemap file that lists both controlled domain families.
- Include the landing and legal pages:
  - `/poliglot-ai.html`
  - `/privacy.html`
  - `/terms.html`
  - `/agreement.html`
  - `/consent.html`
- Include alternate links for Russian and English landing variants at minimum.
- Include `.ru` Russian URLs and `.online` international URLs.
- Include `xhtml:link` alternates for all 35 supported interface languages on the landing URL.
- Submit both domain properties in Search Console so cross-domain sitemap entries are understood as controlled surfaces.

## Runtime Metadata Helper

Add a focused helper rather than scattering DOM metadata updates through the app.

Responsibilities:

- Determine current host family:
  - Russian host: `poliglotai.ru` and `www.poliglotai.ru`.
  - International host: `poliglotai.online` and `www.poliglotai.online`.
  - Local/dev host falls back to Russian defaults for `ru`, `.online` defaults for non-Russian `lang`.
- Determine current language from `window.poliglotSiteI18n.currentLanguage()` or the same URL/localStorage/html fallback.
- Compute canonical URL.
- Upsert `link[rel="canonical"]`.
- Upsert `link[rel="alternate"][hreflang]`.
- Upsert Open Graph/Twitter URL/title/description tags.
- Upsert one JSON-LD script with stable `id`.
- Re-run after the existing `poliglot-language-change` event.

The helper should live in `site-react/src/seoMetadata.ts`, following the current project style. It must avoid broad DOM rewrites and must not conflict with `site-i18n.js`.

## Testing

Add focused Playwright coverage in `site-react/e2e/public-site.spec.ts` or a new public SEO spec:

- Landing has a canonical link for `ru` on a local/default request.
- `?lang=en` produces an English/international canonical or expected local fallback.
- The page has `hreflang` entries for `ru`, `en`, and `x-default`.
- The page has a JSON-LD script containing `Organization`, `WebSite`, `SoftwareApplication`, and `FAQPage`.
- Visible FAQ/AEO questions exist in Russian and are localized when switching to English.
- `robots.txt` and `sitemap.xml` are served by the Vite dev server.
- Existing app links remain relative `/app/` and are not forced to one domain.

Verification commands:

- `npm --prefix site-react run build`
- `npm --prefix site-react run e2e`

If browser MCP is blocked by the system Chrome remote debugging policy, use the project Playwright test runner as the verification fallback and report the MCP limitation.

## Risks And Mitigations

- Cross-domain duplicate content risk: mitigate with host-aware canonical and `hreflang`, not cross-canonicalizing all pages to one host.
- Runtime metadata mismatch: keep metadata computation in one helper and test `ru` and `en`.
- Overclaiming in AEO answers: keep answers factual and avoid guaranteed outcomes.
- Existing user changes in `site-i18n.js` and Playwright tests: do not overwrite them; make narrow patches and review diffs before committing.
- Static sitemap host mismatch: one sitemap intentionally includes both controlled domains, and both domain properties must be submitted in Search Console.

## References

- Google Search Central: SEO Starter Guide - https://developers.google.com/search/docs/fundamentals/seo-starter-guide
- Google Search Central: Optimizing for generative AI features - https://developers.google.com/search/docs/fundamentals/ai-optimization-guide
- Google Search Central: Structured data intro - https://developers.google.com/search/docs/appearance/structured-data/intro-structured-data
- Google Search Central: Structured data policies - https://developers.google.com/search/docs/appearance/structured-data/sd-policies
- Schema.org FAQPage - https://schema.org/FAQPage
