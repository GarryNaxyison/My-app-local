# English Spark Hero Landing Redesign

## Goal

Rework the public NERIVA landing into an English-first product site that keeps the older full-screen React Sparkles hero as the first-screen visual anchor, fixes pricing and dark-theme regressions, and replaces the sticky header with a burger drawer.

## Confirmed Scope

- The first screen must use the previous `.spark-hero` visual direction with `SparklesCore`: full-screen, immersive, dark, and large.
- The current small `.bold-hero`/anomalous-matter hero is not acceptable for the landing.
- All new landing copy is English only for this stage.
- 35-language translation is explicitly deferred until after design approval.
- Public navigation must stop being a fixed/sticky header. It should become a compact top overlay with a burger button that opens a drawer.
- Loading placeholders should use skeleton states, not spinner/preloader patterns.
- Pricing must show clear plan descriptions and real struck-through old prices for paid plans.
- Dark theme must keep body, card, pricing, drawer, and legal navigation text readable.
- Mobile and desktop must have no text overlap, no horizontal overflow, and no button/card label clipping.

## Design Direction

The hero returns to the earlier cinematic Sparkles version: a dark first viewport with layered radial light, large display text, product CTA buttons, proof chips, and a glass product card. The rest of the page becomes a more deliberate product narrative in English: daily loop, modules, web/Telegram entry, pricing, reviews, and final CTA.

The visual language should lean toward a serious AI product cockpit, but without replacing the approved hero. Below the hero, sections use dark, high-contrast surfaces, restrained borders, stable card dimensions, and readable type. Pricing should feel like a clear buying decision, not a decorative dashboard.

## Required Acceptance Checks

- `.spark-hero` exists on the landing and has height at least `90vh` on desktop and mobile.
- `.bold-hero` does not render on the landing.
- The burger button is visible; opening it shows links for Features, Pricing, Reviews, Privacy, Terms, Web app, Telegram, language, and theme.
- `.public-nav` is not `position: fixed` or `position: sticky`.
- `.skeleton-line` and `.skeleton-card` exist as stable placeholders in product preview areas.
- Premium and Platinum plan cards include old prices inside `s` or `del` elements with computed `text-decoration-line` containing `line-through`.
- No visible Cyrillic text appears in landing sections while the current language is English.
- On `390x844` and `1440x1000`, checked text/control/card elements have no horizontal overflow and no incoherent overlap.
- Production build and desktop Chromium e2e pass.

## Out Of Scope

- Translating the redesigned landing into 35 languages.
- Changing the web app product UI outside the public landing.
- Deploying to the server.
