# Neriva Landing SEO Split Design

Date: 2026-07-03

## Context

The current public landing is too long for the main conversion job. It includes the hero, route cards, cockpit explanation, scenario cards, device-flow explanation, daily loop, feature modules, mistake memory, entry-point explanation, pricing, reviews, comparison blocks, long search answers, SEO guides, and final CTA.

The approved direction is a shorter main landing plus separate crawlable SEO pages. Removed long blocks must not disappear; their content should move into static SEO pages.

## Main Landing

Keep the landing focused on conversion:

1. Hero with NERIVA, one concrete promise, and two CTAs: Web app and Telegram.
2. Live product scenario showing one short lesson flow.
3. Four key functions: AI Tutor, Voice Coach, Photo Practice, Mistake Review.
4. Pricing with Telegram Stars and YooKassa/SBP only.
5. Giveaway/community block with the Telegram channel and current social links.
6. Short FAQ.
7. Compact guides links in the footer/late page area only.

Remove from the main body:

- AI tutor vs vocabulary app comparison cards.
- Long SEO/AEO answer cards.
- Big Web app vs Telegram explanation.
- Long system mechanics blocks such as route/cockpit/device-flow/memory-loop explanations.
- Repeated scenario/user-story sections that make the landing feel like a generic SaaS page.

## SEO Pages

Keep the existing generated static SEO architecture, but align the page set to the requested routes:

- `/ai-tutor.html`
- `/speaking-practice.html`
- `/pronunciation.html`
- `/photo-translation.html`
- `/telegram-language-bot.html`
- English alternates under `/en/`.

Each page should carry the long content removed from the main page:

- `/ai-tutor.html`: guided AI Tutor loop and “not just a chat / not just vocabulary” comparison.
- `/speaking-practice.html`: roleplay, conversation scenarios, and daily speaking practice.
- `/pronunciation.html`: voice score, weak words, shadowing, and listening loop.
- `/photo-translation.html`: OCR/photo translation, menus/signs/tasks, and context practice.
- `/telegram-language-bot.html`: Web app vs Telegram explanation and sync model.

## Translation Model

The main landing and the static SEO/AEO pages use different translation rules:

- Main landing: all 35 interface languages must keep visible marketing/product copy in a maintained creative translation layer, not in the generated Google Translate block.
- Static SEO/AEO pages: Russian and English only. Root pages are Russian, `/en/` pages are English, and no static SEO slugs are generated for the other interface languages.
- The landing translation generator must not collect `EnglishSparkLanding.tsx` copy into `site-phrases.js`.
- Tests should fail if curated landing strings appear in the generated machine block, or if SEO/AEO pages expose non-RU/EN static paths.

## Visual Direction

Move the main landing toward a Swiss/editorial product feel:

- white/black surfaces;
- strict grid and hairline rules;
- one blue accent;
- fewer card clusters;
- real product UI fragments and scenario images;
- no neon/glass SaaS overload.

Existing images under `site-react/public/assets/scenarios` and testimonial/social assets remain available, but the main page should use fewer of them.

## Testing

Update Playwright expectations before implementation:

- landing no longer renders comparison cards or long answer cards;
- landing no longer renders course/device/memory/review overload sections;
- landing still renders hero, demo, four module cards, pricing, short FAQ, community, and compact SEO guide links;
- payment copy excludes TON, USDT, crypto, blockchain, and RollyPay;
- static SEO page test expects the five requested slugs and sitemap entries.
- curated landing copy is available in all 35 interface languages through `englishSparkCreativeCopy`;
- generated translation output does not contain the new curated landing marketing strings;
- static SEO/AEO pages stay limited to Russian and English.

## Acceptance Criteria

- The landing is materially shorter and focused on the first-screen product decision.
- Removed long content remains reachable on generated SEO pages.
- The main landing has visible Web app and Telegram CTAs.
- Payment copy mentions only Telegram Stars and YooKassa/SBP.
- SEO pages are generated from one source of truth and included in the sitemap.
- Russian and English static SEO pages have canonical/hreflang metadata.
- Landing copy across 35 interface languages uses the maintained creative override layer rather than machine-generated `site-phrases.js` entries.
- Static SEO/AEO pages remain RU/EN-only and do not create crawlable 35-language SEO duplicates.
- Build and focused Playwright checks pass.
