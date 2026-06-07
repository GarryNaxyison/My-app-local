# Poliglot AI Redesign Plan

## Product Direction

Poliglot AI should feel like a premium AI language coach, not a generic bot wrapper and not a playful school app. The product promise is simple: one AI tutor across Telegram and web that helps people speak, write, listen, translate, review mistakes, learn words, and track progress in 20 interface languages.

The visual language is `Language Intelligence Cockpit`: a calm, high-end workspace with clear controls, language routes, voice waves, progress rings, memory tiles, and premium learning signals. Generated images must support the interface instead of competing with it.

## Brand Attributes

- Confident: the product feels reliable enough for daily learning.
- International: the design works for Latin, Cyrillic, CJK, Georgian, Armenian, and Central Asian language contexts.
- Focused: the web app is a working tool, not a marketing hero inside the app.
- Premium but warm: rich materials and lighting, without dark luxury clutter.
- AI-native: subtle neural routes, translation flows, voice waveforms, and learning signals.

## Visual System

### Palette

- `Ink`: `#101827` for primary text.
- `Graphite`: `#516173` for secondary copy.
- `Cloud`: `#F6F8FC` for app background.
- `Surface`: `#FFFFFF` for panels.
- `Indigo`: `#2D5BFF` as the primary action color.
- `Teal`: `#00A88F` for AI, practice, voice, and tools.
- `Mint`: `#DDF8EF` for soft success states.
- `Gold`: `#F3B84B` for Premium, Platinum, awards, and active achievement.
- `Plum`: `#7C5CFF` only as a small secondary accent.
- `Rose`: `#E45A6F` only for errors and mistake review.

Dark mode keeps the same identity but shifts surfaces to `#0D1320`, `#141E2D`, and uses softened Indigo/Teal/Gold accents.

### Materials

- App surfaces: restrained glass only where useful for depth.
- Navigation: solid, quiet, high-contrast, designed for repeated use.
- Generated imagery: polished 3D enamel/glass/soft metal, consistent camera and lighting.
- Legal pages: editorial cards, clear table styles, no decorative clutter.

### Shape And Spacing

- Primary radius: `10px`.
- Compact controls: `8px`.
- Large marketing bands: `18px` only for hero media and legal shells.
- Desktop app grid: sidebar, workspace, inspector.
- Mobile app grid: top identity/status bar, content, bottom navigation.

### Typography

- Default: Inter/system fallback.
- No viewport-based font scaling.
- Landing H1: strong but readable, with strict max width.
- App headings: compact and functional.
- Legal text: comfortable line height, narrow measure, sticky table of contents on desktop.

## App Screen Map

### Shared App Shell

- Sidebar: product identity, account state, navigation groups.
- Workspace: topbar, view stage, content/chat/tool area, composer where relevant.
- Inspector: progress, daily limits, plan/account actions.
- Mobile: top identity/status, language select, primary action nav, overflow sheet.

### App Views

- `auth`: premium login/register with Telegram-link path, no broken language fallback.
- `home`: command dashboard with daily action cards and current learning state.
- `lesson`: generated lesson card with clear answer flow and next actions.
- `practice`: conversation workspace, prompt starters, correction states.
- `words`: new-word flow, pronunciation, save/review states.
- `word-game`: quiz card, answer choices, result feedback.
- `spelling`: input-focused spelling drill, give-up path, correct answer reveal.
- `vocabulary`: searchable/reviewable learned words.
- `level`: CEFR-style assessment cards.
- `progress`: XP, streak, lessons, practice, voice/photo usage.
- `awards`: 20 rank trophies with consistent trophy art.
- `leaderboard`: compact comparison by language.
- `premium`: Free/Premium/Platinum cards, monthly/yearly, Stars/TON/USDT/card.
- `limits`: daily usage meters and upgrade paths.
- `mistakes`: mistake dictionary and practice loop.
- `tools`: translator, voice-to-text, image translation, GPT agent link.
- `referral`: invite link, balance, bonus explanation.
- `settings`: interface language, learning language, password, logout.

## Public Site Map

### Landing

Sections to redesign:

- Hero: Poliglot AI, Telegram + web, 20 languages, one progress system.
- Country/language proof: 20 interface languages and separate learning language.
- Product capabilities: lessons, practice, voice, photo, dictionary, mistakes.
- Launch discount: Free/Premium/Platinum price story.
- Bot workflow: why Telegram matters.
- Feature grid: concise, icon-led, no repetitive marketing blocks.
- Plan comparison: clear free vs paid limits.
- Real tasks: travel, work, study, pronunciation, reading, translation.
- How it works: choose language, link Telegram, practice daily.
- Pricing: Free/Premium/Platinum, RUB/Stars/USDT.
- CTA: open web app and Telegram.
- FAQ.

### Terms

Keep all legal content, but redesign the reading experience:

- Shared site header and language picker.
- Legal hero consistent with landing.
- Sticky contents on desktop.
- Responsive pricing table.
- Better spacing for long translated paragraphs.

### Privacy

Same legal system as Terms:

- Data protection hero.
- Sticky contents.
- Clear cards for data, purpose, payments, voice/photo, retention, rights.
- Translation-safe paragraphs.

## Asset Strategy

Generated assets must use one master style:

`Premium AI language cockpit, polished 3D enamel and glass objects, deep indigo studio background, teal neural routes, mint highlights, selective gold for achievement and premium, soft cinematic product lighting, clean negative space, no text, no letters, no UI screenshots, no people, no hands.`

### Web App Assets

- Plan cards: `plan-free-light.png`, `plan-free-dark.png`, `plan-premium-light.png`, `plan-premium-dark.png`, `plan-platinum-light.png`, `plan-platinum-dark.png`.
- View headers: home, lesson, practice, progress, awards, words, word-game, spelling, vocabulary, level, leaderboard, premium, limits, mistakes, tools, referral, settings.
- Icons: one light/dark icon pair per app view, all same camera/lighting/material.
- Theme selection: `/app` switches image pairs from `html[data-theme]`; legacy non-themed filenames are light-theme fallbacks only.

### Public Site Assets

- `poliglot-ai-avatar.jpg`: main landing hero and online bot hero.
- `times.jpg`: languages/features.
- `primi.jpg`: use cases/pricing.
- `write.jpg`: start/how-it-works.
- `yspex.jpg`: CTA/success.
- `site-legal-shield.png`: Terms/Privacy shared legal hero.

### Asset QA

Reject any image with:

- fake text, pseudo letters, numbers, flags, logos, watermarks;
- phones, dashboards, browser screenshots, documents, certificate plaques;
- people, faces, hands;
- isolated generic orb or sphere;
- mismatched palette or camera style.

## Figma Deliverable

The Figma file should include:

- cover page with design direction;
- color and type tokens;
- components: button, input, select, nav item, card, plan card, legal shell;
- desktop app wireframe;
- mobile app wireframe;
- landing wireframe;
- Terms/Privacy wireframe;
- asset board with required filenames and prompts.

Figma is a design source for implementation and QA, while code remains the production source.

## I18n Requirements

Supported site languages:

`ru`, `en`, `es`, `de`, `fr`, `it`, `zh`, `ja`, `ko`, `tg`, `uz`, `tt`, `hy`, `kk`, `ky`, `ka`, `uk`, `pl`, `ro`, `pt`.

Checks:

- all `data-i18n` keys exist for all languages or have phrase fallback;
- landing, Terms, and Privacy do not show Russian in non-Russian language mode;
- no mojibake markers such as `Ð`, `Ñ`, `Рџ`, `РЎ`, `вЂ`;
- CJK lines fit in buttons/cards;
- German/French/Uzbek long strings wrap without overflow;
- pricing and legal tariff names stay consistent.

## Implementation Order

1. Create this plan and Figma board.
2. Update ComfyUI generator to produce the unified asset pack.
3. Generate contact sheets and inspect visual consistency.
4. Replace app visual assets and CSS tokens.
5. Redesign app shell and views with minimal behavioral risk.
6. Redesign public landing and legal pages.
7. Rebuild public site translations.
8. Run Go tests and JavaScript syntax checks.
9. Verify desktop and mobile in browser.
10. Deploy `/opt/aibot/aibot` plus `/opt/aibot/web`.
11. Deploy public site to `/var/www/poliglot...`.
12. Verify production URLs, assets, headers, language switching, and mobile.

## Definition Of Done

- The app and public site look like one product.
- Generated images share one visible art direction.
- HTML is served from disk next to the bot, not embedded in the binary.
- `/app/assets/...` works behind Caddy.
- 20 languages switch correctly on landing, Terms, and Privacy.
- Desktop 1280x720 and mobile 390x844 are visually checked.
- Production bot and site are deployed and healthy.
