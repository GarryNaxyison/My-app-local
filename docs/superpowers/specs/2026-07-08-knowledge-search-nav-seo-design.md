---
project: NERIVA
date: 2026-07-08
status: approved
tags:
  - spec-it
  - neriva
  - knowledge-base
  - seo
---

# NERIVA Knowledge Search, Navigation, and SEO Design

## Problem

The knowledge hub has 100 articles, so a flat list forces users to scroll too much. The existing search filters cards, but it does not show matching topics while the user types. Separately, the root landing SEO snippet can miss the exact intent phrase "изучение языков", even though supporting SEO pages already target it.

## Approved Approach

Implement option 1:

- Keep the current `/knowledge/?q=...` filtering behavior.
- Add a live suggestion dropdown under the knowledge search input with the best matching article topics.
- Add a compact topic navigation layer on the hub so users can jump by category instead of scrolling through 100 cards.
- Tighten Russian landing SEO title/description so the root page clearly says NERIVA is for "изучение английского и языков".

## Knowledge Hub Behavior

The search input should still filter `.knowledge-card` elements using the current partial-token search. When the user types, a dropdown appears under the input with up to six matching articles. Each suggestion shows the article title, category, and reading time. Clicking a suggestion navigates to the article. Pressing Enter keeps the existing filtered-list behavior unless a suggestion is actively focused by keyboard.

The quick navigation should be visible near the top of the hub and should link to the existing category route pages. It should show all 10 topic routes, each with the current article count and a few representative article titles. This gives a useful overview without creating new URLs.

## SEO Behavior

The Russian landing title and description should include:

- `AI-репетитор английского`
- `изучение языков`
- `Telegram` and `web app`

The change should update static pre-JavaScript metadata and runtime metadata/JSON-LD through the existing `landingSeoCopy` path.

## Constraints

- Do not add new routable URLs.
- Do not rewrite the knowledge hub as a React page.
- Preserve current sitemap behavior.
- Preserve 100 generated article cards and all category route pages.
- Keep static SEO/AEO pages RU/EN only.

## Verification

- Playwright e2e must cover the suggestions dropdown, suggestion click navigation, quick topic navigation, and landing SEO metadata.
- Existing knowledge search tests for partial query and `/knowledge/?q=...` must keep passing.
- Build must regenerate public knowledge pages and static metadata.
