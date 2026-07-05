# NERIVA Marketing Completion Design

## Goal

Improve the NERIVA landing copy and section flow without changing the existing premium dark visual system. The page should sound like normal Russian from a real product owner, not generic AI marketing text.

## Scope

- Keep the current `EnglishSparkLanding` structure, dark premium style, product screenshots, buttons, pricing cards, and responsive behavior.
- Add a stronger pain/outcome block after the hero.
- Add scenario lanes for travel, work, exam, and pronunciation.
- Make proof more concrete through product screenshots and plain-language outcomes.
- Make Premium feel like the natural daily plan, not only a higher limit.
- Rewrite Russian copy to be more direct and human.

## Implementation Notes

- Main content stays centralized in `site-react/src/landingContent.ts`.
- Rendering changes stay in `site-react/src/EnglishSparkLanding.tsx`.
- Styling additions stay in `site-react/src/englishSparkLanding.css`, reusing the existing cards, gradients, borders, type scale, and screenshot presentation.
- Existing product screenshots in `site-react/public/assets/product` are enough for this pass; new screenshots are optional only if verification shows a content gap.

## Verification

- Build the site with `npm run build` in `site-react`.
- Use local Playwright screenshot/text extraction on `http://127.0.0.1:5175/neriva.html`.
- Check desktop and mobile viewports for text overflow, blocked CTAs, broken layout, and unnatural Russian.
