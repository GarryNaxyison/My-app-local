# Premium AI Tutor Cockpit Landing Design

## Context

The public React landing lives in `site-react/src/PublicSiteApp.tsx` and builds into `Сайт полиглота для бота/`. The private learning web app lives separately in `web-react` and is served from `/app/`. This design updates the public landing only.

The current landing already has a full-screen Three.js hero animation through `GenerativeArtScene` inside `.landing-hero__matter`. That animation is a hard constraint: the refresh must keep the same full-screen animated hero system and build a more premium product story around it.

The user approved the `Premium AI Tutor Cockpit` direction after requesting competitor and marketing analysis. The goal is not a generic language-learning landing. The page must look expensive, technological, and credible enough to sell a Premium AI language product.

## Market Read

The language-learning market is crowded. The strongest competitors sell a clear emotional outcome:

- Duolingo Max: AI roleplay and answer explanation on top of a habit-forming mass-market product.
- Busuu: structured courses, guided progress, and a practical learning path.
- Babbel: practical adult learning, real-life dialogs, speech recognition, and live/class value.
- ELSA Speak: a focused AI pronunciation coach.
- Speak: AI speaking tutor and conversation confidence.
- Praktika: AI avatar tutors and speaking practice.
- Memrise and Mondly: broad language catalog, video/chatbot/AR-style practice, and app-based repetition.

The common winning message is not "many lessons". It is confidence: speak more often, make mistakes safely, receive feedback, and see measurable progress.

NERIVA should compete by combining the premium AI tutor story with a broader practical ecosystem: web app, PWA, Telegram, daily route, roleplay, pronunciation, photo practice, vocabulary, mistakes, notes, offline decks, and one learning profile.

## Approved Strategy

Use the `Premium AI Tutor Cockpit` strategy.

The landing should feel like entering a high-end AI learning cockpit:

- The first viewport is cinematic and product-led, with the existing full-screen hero animation.
- The UI demo should feel alive and useful, not decorative.
- Product sections should show what the learner does today, what the AI corrects, and what Premium unlocks.
- The page should make Premium feel like the natural way to study seriously, while keeping Free as a credible start.

## Non-Negotiable Requirements

- Keep the existing full-screen hero animation and its placement. Do not replace `GenerativeArtScene`, shrink the canvas into a card, or change the animation engine.
- Keep public app CTAs on `/app/`.
- Keep Terms and Privacy links in the header and footer. Do not remove or hide legal navigation.
- Keep the 35 interface language promise and language selector.
- Design desktop web and mobile web separately. The final layout must be intentionally responsive, not merely stacked.
- Keep the current legal pages working: `privacy.html` and `terms.html`.
- Do not promise unsupported dictionary coverage for every interface language. The landing may say 35 interface languages, but learning dictionaries remain separate product capability.
- Do not invent ratings, institution logos, user counts, certificates, or outcomes that are not supported by the product.

## Product Truths To Surface

Use only capabilities already represented in the project:

- Daily route / Today screen: weekly plan, daily quests, habit calendar, XP, streak, and daily bonus.
- AI Tutor lessons: guided lesson flow, learner answer checks, feedback, XP, and structured progression.
- Practice: text, voice, image/file controls, model phrases, and saved phrases.
- Roleplay: realistic scenarios and corrected conversation practice.
- Pronunciation and listening: waveform, score, weak words/sounds, tips, next phrase, and history.
- Photo and translation: image/file context can become translation and practice.
- Vocabulary, Review, Spelling, Mistakes, Notes/Phrasebook: weak material becomes repeatable training.
- Offline/PWA: saved decks can be used outside a stable connection.
- One profile across web app, PWA, and Telegram.
- Premium and Platinum plans with higher limits and voice/photo-heavy use.
- 35 interface languages in the public language selector.

## Positioning

Primary positioning:

> NERIVA is a premium AI tutor that turns each day into a guided language session: lesson, conversation, pronunciation, photo practice, mistakes, and progress in one web + Telegram profile.

Core promise:

> Speak, make mistakes, get precise feedback, and return to weak spots every day.

Avoid:

- "Just another app to memorize words."
- "Learn any language instantly."
- "Human-teacher replacement" claims.
- Overloaded lists of internal modules without explaining user value.

## Message Map

### Hero Message

Hero should make the premium product obvious in the first 5 seconds.

Recommended Russian direction:

- Eyebrow: `Premium AI-репетитор для ежедневной практики`
- H1: `Говорите на новом языке с AI-репетитором, который помнит ваши ошибки`
- Body: `Урок, диалог, произношение, фото-перевод, словарь ошибок и прогресс собираются в один профиль web, PWA и Telegram.`
- Primary CTA: `Начать бесплатно`
- Secondary CTA: `Открыть Telegram`

The hero should keep:

- brand signal in the header;
- language picker;
- proof chips for `35 языков интерфейса`, `A1-C2`, and `web + PWA + Telegram`;
- a live product demo card;
- enough first-viewport height for the animation to feel immersive.

### Product Proof

The page should show six high-value jobs:

1. Daily route: the learner knows what to do today.
2. Roleplay: the learner practices travel, work, interview, and everyday situations.
3. Pronunciation: the learner sees score, weak words, and next phrase.
4. Photo practice: menus, signs, exercises, and files turn into study material.
5. Mistakes and vocabulary: errors become review, spelling, saved notes, and offline decks.
6. One profile: web app, PWA, and Telegram keep the same progress.

### Premium Value

Premium must be framed as seriousness and depth, not only "more limits".

Premium story:

- more daily guided learning;
- voice and photo-heavy practice;
- more roleplay and practice capacity;
- deeper correction loop through mistakes, notes, vocabulary, and offline review.

Platinum story:

- intensive preparation for travel, work, exams, or daily speaking;
- highest limits and the least friction for frequent sessions.

### Objection Handling

FAQ should answer:

- Can I use only Telegram?
- What is shared between web app and Telegram?
- What changes in Free, Premium, and Platinum?
- Can I study on mobile / PWA / offline?
- Are Terms and Privacy available?

## Page Architecture

Keep the existing block list, but rewrite the story around the cockpit strategy:

1. Header: brand, features, pricing, reviews, Privacy, Terms, language selector, theme toggle, `/app/` CTA.
2. Hero cockpit: full-screen animated background, strong premium offer, app/Telegram CTAs, language chooser, proof chips, interactive AI tutor demo.
3. Premium product strip: outcome cards for speaking, pronunciation, photo practice, and mistake memory.
4. Feature grid: six real capability cards with concise value-first copy.
5. Workflow: choose goal/language/level, complete a short session, repeat weak points.
6. Pricing: Free, Premium, Platinum with the current 300 ₽ and 590 ₽ visible offer and crossed-out old prices if retained.
7. Reviews: three realistic user stories, not fake metric claims.
8. Progress panel: XP, streak, mistakes, notes, weak words, pronunciation history, and awards as one retention loop.
9. FAQ: practical objections and legal access.
10. Final CTA: `Начать бесплатно` and optional Telegram secondary action.
11. Footer: keep Privacy, Terms, contacts, Telegram bot, and app link.

## Desktop Web Design

Desktop should feel like a premium cockpit, not a content blog.

Hero:

- Keep the animation visible across the full hero.
- Use a dark veil for readability but preserve animated depth.
- Use a two-column composition: left offer and proof, right product cockpit/demo.
- Keep the next section hinted below the first viewport where practical.
- Demo card should show realistic product states: today's route, roleplay, pronunciation, photo, weak words, score.

Below hero:

- Use dense but polished SaaS-like bands.
- Avoid nested cards and decorative-only sections.
- Use compact, high-contrast cards with clear hierarchy.
- Pricing should be easy to compare at desktop width without excessive vertical scrolling.

## Mobile Web Design

Mobile is not a collapsed desktop page. It needs its own rhythm.

Hero mobile:

- Keep the animated background but reduce overlay density so text stays readable.
- H1 must not overflow or overlap the demo.
- Primary CTA should be immediately visible.
- Language picker should become a compact, touch-safe grid or horizontal group.
- Product demo should appear after the primary offer and not push the CTA below an unreachable fold.

Mobile sections:

- Cards stack in one column with stable heights where needed.
- Pricing cards must be readable without horizontal scroll.
- Terms and Privacy links must remain reachable in nav/footer.
- Language selector must remain usable and readable in light/dark themes.
- No text should overlap in Russian or translated languages.

## Visual Direction

Use a premium technology palette based on the existing identity:

- dark glass first viewport;
- cyan and gold highlights;
- restrained white/ink typography;
- crisp 8px-radius panels consistent with the current design;
- subtle motion on CTAs and demo states;
- product UI panels instead of decorative marketing illustrations.

Avoid:

- purple-blue generic gradients as the dominant theme;
- beige/cream luxury templates;
- dark static hero without product proof;
- huge marketing cards that do not show real product functionality;
- new decorative orbs/blobs.

## Technical Design

Primary implementation files:

- `site-react/src/PublicSiteApp.tsx`: copy, arrays, demo tabs, feature sections, pricing text, reviews, FAQ, CTA targets.
- `site-react/src/styles.css`: responsive layout, cockpit visual hierarchy, mobile refinements, pricing/product proof styling.
- `site-react/e2e/public-site.spec.ts`: update `/app/` expectation, preserve legal links, 35 languages, hero canvas, block counts, pricing, mobile/desktop checks.
- Built output under `Сайт полиглота для бота/` after `npm run site:build`.

Do not edit:

- `web-react` app behavior;
- Go payment or premium backend behavior;
- legal document content except existing generated public-site output if build touches it;
- `GenerativeArtScene` internals unless a browser verification shows a rendering bug.

## 21st.dev / Figma / MCP Use

Use available MCP context where practical:

- Nx: use local Nx CLI fallback if MCP tool is unavailable.
- Playwright/browser: required for desktop and mobile visual verification.
- GitHub: commit and push verified changes.
- 21st.dev Magic: use for component pattern inspiration only if available; if API rendering takes 5-6 minutes, do not block implementation on it unless a specific component pattern is needed.
- Figma: use if a concrete Figma file/link is provided or available in MCP; otherwise do not invent Figma-derived details.

## Verification

Run before claiming completion:

```bash
npm --prefix site-react run build
npm --prefix site-react run e2e -- --reporter=line
```

Browser verification:

- desktop viewport around 1440x900;
- mobile viewport around 390x844;
- hero canvas exists, is visible, and still covers the hero;
- CTA links resolve to `/app/`;
- Privacy and Terms links are visible/reachable;
- language selector exposes all 35 interface locales;
- no `/app/v2/` public CTA remains;
- no text overlap in hero, language picker, product proof, pricing, reviews, FAQ, footer;
- pricing old/new prices remain readable;
- mobile first viewport feels intentionally designed, not merely stacked.

## Out Of Scope

- Rebuilding the private web app.
- Changing backend limits, payments, auth, Telegram bot behavior, or vocabulary data.
- Adding unsupported claims, user counts, external ratings, or partner logos.
- Removing legal pages or language support.
- Replacing the hero animation.

## Open Assumptions

- The current React landing source is the source of truth, and static files in `Сайт полиглота для бота/` are build output.
- The approved public app target is `/app/`.
- Russian source copy is acceptable as the primary copy surface, with existing public-site localization machinery preserving interface-language support.
- The first implementation pass can be done locally without waiting for a 21st.dev render unless a specific component generation need appears.
