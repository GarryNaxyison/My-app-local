# Bold Product Cockpit Landing Design

## Context

The public landing is the React app in `site-react`, not the private `web-react` app shell. The source files for this redesign are:

- `site-react/src/PublicSiteApp.tsx`
- `site-react/src/styles.css`
- `site-react/e2e/public-site.spec.ts`

The generated static output under `Сайт полиглота для бота/` must be updated only through the `site-react` build, not by manual HTML edits.

The current landing already has user-facing product content for NERIVA: AI Tutor, roleplay, pronunciation, photo practice, mistakes, notes, offline decks, web app, PWA, Telegram, Free, Premium, Platinum, and 35 interface languages. The redesign should use these real product capabilities, not invented metrics, institutions, certificates, or fake usage numbers.

## Approved Direction

Use the **Bold Product Cockpit** direction.

The landing should feel like a strong product launch page for a premium AI language tutor. It should not read as a generic language-learning feature list. The first impression should be bold, designed, and product-specific: large typography, high contrast, clear product mockups, and a confident message about daily guided practice.

Approved user decisions:

- Keep the existing hero animation only in the first block.
- Do not repeat the hero animation in later sections.
- Use two equal primary entry points: Web app and Telegram.
- Build designer product mockups directly in code, based on real product functions.
- Do not use Inter or the misspelled Inrer. Use the design skill font direction: Unbounded for display headlines and Manrope for body/interface copy.
- Use the `ui-ux-pro-max` and `frontend-design` guidance for accessibility, responsive layout, typography, visual hierarchy, and product-led page structure.

## Visual System

The chosen visual character is Bold Product, implemented as a disciplined product cockpit rather than a chaotic page.

Typography:

- Display: `Unbounded`
- Body and UI: `Manrope`
- Fallbacks: `ui-sans-serif`, `system-ui`, `-apple-system`, `BlinkMacSystemFont`, `"Segoe UI"`, `sans-serif`
- No `Inter` in the root font stack or landing-specific overrides.

Color and material:

- Primary surface: deep ink / black-blue for the hero and selected contrast sections.
- Product surfaces: white and near-white panels with high-contrast ink text.
- Accents: electric blue, acid cyan, and warm gold.
- Avoid a one-note purple/blue gradient page.
- Avoid beige/cream luxury styling.
- Avoid decorative gradient orbs, bokeh blobs, or visual noise.
- Use crisp panels, strong section boundaries, and controlled 8px radii unless a specific product mockup needs a slightly larger device frame.

Motion:

- Keep the existing `GenerativeArtScene` hero animation in the first block.
- Use small UI transitions for tabs, hover states, and mockup state changes.
- Respect reduced motion.
- Do not add heavy scroll-linked animation or repeat animated hero backgrounds below the first viewport.

Iconography:

- Use lucide-react icons where icons are needed.
- No emoji icons.
- Icon-only controls require accessible labels.

## Page Architecture

### 1. Header

The header should remain practical and conversion-focused:

- Brand: NERIVA with existing logo asset.
- Navigation: product sections, pricing, reviews, Privacy, Terms.
- Language selector remains available.
- Theme toggle remains available if it still works reliably after redesign.
- Two equal entry links can appear in compact form: Web app and Telegram.

The header must not hide Privacy or Terms.

### 2. Hero: Bold Product Cockpit

The hero is the only section with the existing full-screen animated background.

Hero content:

- Strong short headline about speaking with an AI tutor that remembers mistakes.
- One concise supporting paragraph.
- Two equal CTA buttons:
  - `Открыть Web app` -> `/app/`
  - `Открыть Telegram` -> `https://t.me/NERIVAapp_bot`
- One compact proof row, not a crowded grid.
- One large code-built product mockup showing an AI Tutor session.

Hero mockup:

- Tabs or segmented states: `Урок`, `Диалог`, `Голос`, `Фото`.
- Content should be based on real product flows:
  - lesson explanation and answer check;
  - roleplay with a practical phrase;
  - voice score / weak words / shadowing;
  - photo text turned into translation and practice.

Remove first-viewport overload:

- Move the 8 language cards below the hero or remove them from the first viewport.
- Move path chips and repeated proof claims below the hero.
- Do not make the hero a dense feature matrix.

### 3. What You Do Today

This is the first section after the hero and should establish the daily learning loop.

Content:

1. Get a focused task.
2. Answer with text, voice, or photo.
3. See the correction and repeat the weak spot.

Design:

- Bright, high-contrast band.
- One large code-built daily route mockup.
- Three clear steps, each with a real action and outcome.
- No decorative-only cards.

### 4. AI Tutor Modules

Show the product breadth through designer product mockups, not plain feature cards.

Modules:

- AI Tutor: guided lesson, answer check, XP.
- Roleplay: travel/work/exam/speaking scenario with corrected phrase.
- Voice Coach: pronunciation score, weak words, next phrase.
- Photo Practice: image/file text becomes translation, note, and practice prompt.

Design:

- Four substantial mockup tiles.
- Each tile has one clear job, one visual product state, and concise copy.
- Use real feature names where they help comprehension, but avoid internal jargon overload.

### 5. Mistake Memory Loop

This section explains the product's strongest differentiator: mistakes become training material.

Content:

- Mistakes
- Weak words
- Notes / Phrasebook
- Review
- Spelling
- Offline decks
- XP, streak, awards, daily bonus

Design:

- Dark or high-contrast system section.
- Visual metaphor: a correction loop or memory board built in code.
- Make the loop feel like a product system, not a list of modules.

### 6. Web App + Telegram

Web app and Telegram are equal product entry points.

Content:

- Web app: better for focused sessions, dashboard, progress, longer lessons, pricing.
- Telegram: better for fast practice, reminders, voice, quick continuation.
- One profile keeps progress, Premium, mistakes, notes, and practice connected.

Design:

- Two equal panels or side-by-side product mockups.
- Neither CTA should appear secondary.
- Mobile layout should keep both choices visible and touch-safe.

### 7. Pricing

Pricing must be clear, not decorative.

Plans:

- Free: credible start without payment.
- Premium: regular daily learning with AI Tutor, voice, photo, and expanded limits.
- Platinum: intensive use with the highest daily limits.

Requirements:

- Keep current prices and old prices only if they are already supported by the product copy.
- Show payment methods: Stars, YooKassa, TON, USDT.
- Compare limits clearly.
- Use product-value language, not only "more limits".

### 8. Reviews And Final CTA

Reviews:

- Use realistic short user stories.
- Do not invent ratings, user counts, institution logos, certificates, or outcome guarantees.

Final CTA:

- Repeat the two equal entry points:
  - Web app
  - Telegram
- Keep footer legal links, contacts, bot link, and app link.

## Responsive Design

Desktop:

- Hero uses a bold two-column product cockpit layout.
- Product mockups should feel large and intentional, not small decorative previews.
- Pricing should compare plans without excessive vertical scrolling.
- Later sections should use clear bands with strong section jobs.

Mobile:

- Hero animation remains in the first block but content density is reduced.
- Headline must not overflow in Russian or translated languages.
- Equal CTAs must remain visible, touch-safe, and readable.
- Product mockups stack with stable dimensions.
- No horizontal scroll.
- Language selector, Privacy, and Terms must remain reachable.

## Accessibility And UX Requirements

- Maintain color contrast of at least 4.5:1 for normal text where practical.
- Keep touch targets at least 44px high.
- Preserve visible focus states.
- Use semantic buttons/links.
- Provide meaningful labels for icon-only controls.
- Respect `prefers-reduced-motion`.
- Reserve space for images/mockups to avoid layout shift.
- Avoid text overlap across desktop and mobile.

## Technical Scope

Primary files:

- `site-react/src/PublicSiteApp.tsx`
- `site-react/src/styles.css`
- `site-react/e2e/public-site.spec.ts`

Build output:

- Generated through `npm --prefix site-react run build`.
- Output goes to `Сайт полиглота для бота/`.

Do not modify:

- `web-react` behavior.
- Go backend or payment logic.
- Legal document source content, except generated output changed by the public-site build.
- `GenerativeArtScene` internals unless browser verification shows a rendering bug.

## MCP And Tooling Notes

- 21st.dev Magic was used for component inspiration. The useful pattern is a hero with a mockup, but implementation should follow this repository's existing React and CSS structure.
- Figma MCP is available, but no Figma file or node URL was provided. Do not claim Figma-derived design details.
- Nx MCP endpoint was present, but one Nx MCP query returned an internal error. Use local npm/Nx commands as fallback.
- Playwright/browser verification is required after implementation for desktop and mobile screenshots/layout checks.
- Subagents were used only for read-only routine context checks. Their findings confirmed that the public landing is `site-react` and that Inter is currently declared in `site-react/src/styles.css`.

## Test Strategy

Update Playwright tests to protect durable contracts rather than overfitting to every marketing sentence.

Tests should verify:

- Public landing loads from `poliglot-ai.html`.
- Hero animation canvas exists only in the first hero block.
- No `Inter` appears in the active root font stack.
- Web app and Telegram CTAs are both present and visually equal enough to be primary choices.
- Privacy and Terms links remain visible.
- Language selector remains usable.
- Pricing includes Free, Premium, Platinum and supported payment methods.
- Product sections for AI Tutor, Voice Coach, Photo Practice, mistake loop, Web app, and Telegram are present.
- Mobile viewport has no horizontal overflow or text clipping.

## Acceptance Criteria

- The landing uses the Bold Product Cockpit direction.
- The existing hero animation is preserved only in the first block.
- Later sections are redesigned without hero animation reuse.
- The page uses code-built designer product mockups based on real NERIVA functions.
- Web app and Telegram are equal primary CTAs.
- Inter/Inrer is not used.
- The page remains responsive and readable on mobile.
- Legal links and 35 interface-language promise remain available without overpromising learning dictionary coverage.
- The generated static public site is rebuilt from React sources.
