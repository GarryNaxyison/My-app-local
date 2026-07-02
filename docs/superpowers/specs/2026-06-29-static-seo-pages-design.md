# Static SEO Pages Design

Date: 2026-06-29

## Summary

NERIVA needs a small bilingual SEO page cluster that promotes the product as a full language learning web app. Telegram remains a companion channel for quick mobile practice, notifications, voice messages, and synced study state, not the flagship product positioning.

The first wave adds five focused Russian pages for `poliglotai.ru` and five matching English pages for `poliglotai.online`. These pages are static HTML, generated from a shared content source, and included in the public sitemap. They should be indexable and reachable through normal links, but they must not take visual priority over the main landing page.

## Goals

- Add search-focused pages for people looking for an AI English tutor, speaking practice, pronunciation training, practical English for work/travel, and a language learning web app.
- Position NERIVA as a web app first.
- Provide Russian and English versions with correct canonical and `hreflang` relationships.
- Keep the main landing focused on the product, pricing, and app entry.
- Make the SEO pages lightweight, crawlable, maintainable, and covered by focused checks.
- Avoid hidden links, keyword stuffing, doorway-page behavior, or low-value generated text.

## Non-Goals

- No visual redesign of the main landing.
- No blog engine, CMS, comments, author profiles, or publishing workflow in this iteration.
- No mass page generation beyond the approved five topics and two languages.
- No claim that SEO pages guarantee ranking, traffic, language fluency, exam results, or job outcomes.
- No paid-link or guest-post strategy in this step.

## Current Context

- Public site source lives in `site-react`.
- The build outputs into the static public site directory used for deployment.
- `site-react/public` already contains static search assets such as `robots.txt`, `sitemap.xml`, social preview images, and site language assets.
- Existing landing SEO metadata and JSON-LD are handled in React/browser code.
- The SEO page cluster does not need React interactivity.

## Recommended Architecture

Use generated static HTML pages instead of React pages.

The implementation should add:

- a source content module for SEO page data;
- a generator script that writes static HTML pages;
- a shared static CSS file for the SEO pages;
- sitemap updates generated from the same page list;
- focused tests that verify metadata, content, canonical links, language alternates, and sitemap inclusion.

This keeps every public SEO page as normal HTML while avoiding duplicated hand-written markup across ten files.

## URL Strategy

The first wave contains five page topics.

Russian canonical URLs:

- `https://poliglotai.ru/ai-english-tutor.html`
- `https://poliglotai.ru/english-speaking-practice.html`
- `https://poliglotai.ru/english-pronunciation-trainer.html`
- `https://poliglotai.ru/english-for-work-and-travel.html`
- `https://poliglotai.ru/language-learning-web-app.html`

English canonical URLs:

- `https://poliglotai.online/en/ai-english-tutor.html`
- `https://poliglotai.online/en/english-speaking-practice.html`
- `https://poliglotai.online/en/english-pronunciation-trainer.html`
- `https://poliglotai.online/en/english-for-work-and-travel.html`
- `https://poliglotai.online/en/language-learning-web-app.html`

Each page pair must include:

- `rel="canonical"` pointing to the canonical URL for the current language;
- `hreflang="ru"` pointing to the Russian canonical page;
- `hreflang="en"` pointing to the English canonical page;
- `hreflang="x-default"` pointing to the English canonical page.

If the same static directory is served by both domains, duplicate host variants are acceptable only when their canonical links point back to the intended canonical domain.

## Page Topics

### AI English Tutor

Search intent: users looking for an online AI English tutor.

Primary message: NERIVA guides the learner through lessons, answers, corrections, repetition, progress, and practice inside a web app.

### English Speaking Practice

Search intent: users looking for speaking practice without scheduling a human tutor.

Primary message: NERIVA provides roleplay, free practice, feedback, saved mistakes, and repeatable sessions.

### English Pronunciation Trainer

Search intent: users looking for pronunciation, voice, and shadowing practice.

Primary message: NERIVA supports voice practice, weak-word feedback, shadowing, listening loops, and progress history.

### English For Work And Travel

Search intent: users preparing for practical English situations.

Primary message: NERIVA helps rehearse travel, work calls, messages, cafes, hotels, airports, and everyday scenarios.

### Language Learning Web App

Search intent: users comparing language learning apps, AI tutors, and web apps.

Primary message: NERIVA is a full web app with lessons, practice, vocabulary, mistakes, notes, progress, subscriptions, and optional Telegram sync.

## Page Structure

Each static SEO page should use the same structure:

1. Header with NERIVA brand, language switch, compact navigation back to the landing, and primary CTA to `/app/`.
2. Hero-like article intro with H1, short value proposition, and app CTA.
3. Three to five product-specific sections answering the page's search intent.
4. A concise comparison or "when this helps" section.
5. FAQ with four to six direct questions and answers.
6. Internal links to the other SEO pages.
7. Footer links to the main landing, app, legal pages, and social profiles.

The main landing should link to the SEO pages only through a small footer or "Guides" area. These links must be visible and crawlable, but they should not appear in the hero, main CTA area, or pricing flow.

## Metadata

Each page needs language-specific metadata:

- unique `title`;
- unique `meta description`;
- `robots` set to `index, follow`;
- canonical link;
- Russian, English, and `x-default` alternates;
- Open Graph and Twitter metadata;
- JSON-LD graph containing `Organization`, `WebSite`, `SoftwareApplication`, `Article` or `WebPage`, `BreadcrumbList`, and `FAQPage` where the visible FAQ matches the markup.

The application schema should describe NERIVA as an `EducationalApplication` with operating systems `Web, PWA, Telegram`.

## Content Rules

- The copy must be useful to a learner, not a block of repeated keywords.
- Use web app first wording.
- Mention Telegram only as optional sync, quick access, or companion practice.
- Keep claims factual and product-specific.
- Avoid fabricated reviews, fake author bios, fake statistics, and unsupported rankings.
- Do not hide SEO text from users.

## Sitemap And Robots

`robots.txt` keeps both sitemap references:

- `https://poliglotai.ru/sitemap.xml`
- `https://poliglotai.online/sitemap.xml`

`sitemap.xml` should include:

- existing landing and legal pages;
- all five Russian canonical SEO pages;
- all five English canonical SEO pages;
- `xhtml:link` alternates for each Russian/English page pair.

Both domain properties should be submitted in Google Search Console after deployment.

## Testing

Add focused checks for the static SEO pages:

- each generated page returns HTML with a real H1 and no empty React root dependency;
- each page has canonical, `hreflang`, description, Open Graph, and JSON-LD;
- each page has visible FAQ content matching its JSON-LD FAQ;
- all ten canonical URLs appear in the sitemap;
- the main landing contains only a compact footer or guides link area for the SEO pages;
- the primary app CTA remains `/app/`.

Recommended verification commands:

- `npm --prefix site-react run build`
- `npm --prefix site-react run e2e`

If Playwright browser automation is not needed for the spec-only step, defer it to implementation verification.

## Implementation Boundaries

The implementation should change only the public site and SEO assets:

- `site-react` source, scripts, public assets, and e2e tests;
- generated output under the public deployment directory if the project keeps built assets committed;
- no changes to Go backend behavior, payments, authentication, or learning logic.

Existing unrelated working tree changes must not be reverted.

## Risks And Mitigations

- Duplicate cross-domain pages: mitigate with canonical and `hreflang` tags.
- Thin content: keep each page focused on a distinct search intent and actual product capability.
- Main landing clutter: expose SEO pages only through footer/guides links.
- Maintenance drift: generate pages, sitemap entries, and language alternates from one source of truth.
- Overclaiming: avoid guaranteed outcomes and unsupported comparisons.

## Approval Gate

After this spec is reviewed, the next step is an implementation plan. Implementation should not start until the written spec is approved.
