# Landing Goal Courses Redesign Design

## Context

The public landing page is implemented in `site-react/src/PublicSiteApp.tsx` and built into `Сайт полиглота для бота/` with `npm run site:build`. The `/app/` web application is a separate Vite app under `web-react`; this redesign must not rebuild or change `/app/` unless explicitly requested.

The current landing already has the required block list:

- navigation
- hero with animated visual, CTA, language chooser, proof chips, and interactive demo
- course strip
- features
- workflow
- pricing
- reviews
- progress
- FAQ
- final start CTA
- footer and legal pages

The redesign must keep that block list and the main product concept: Poliglot AI is an AI language tutor available in web app, PWA, and Telegram with one learning profile.

## Approved Direction

Use the `Goal Courses` direction selected by the user. The landing should feel closer to Babbel in its first decision: the learner chooses a goal, and the product assembles the route. It should still retain Duolingo-like momentum through short sessions, progress, streak/XP, awards, and repeatable daily practice.

The user also requested a specific hero visual constraint: use the old landing's full-screen animated hero background. The new goal-based content should sit on top of that full-screen animation instead of replacing it with a static course-card layout.

## Product Claims To Surface

Use only functionality that exists in the codebase and docs:

- AI Tutor lessons with scenario, teaching point, model answer, slot-aware practice, final word check, and review.
- Lesson and practice modes in Telegram and web.
- Pronunciation and shadowing with STT/TTS, scoring, weak words or sounds, tips, and progress history.
- Photo/text translation and practice from image context.
- Vocabulary training, review, spelling, phrasebook, mistakes, and offline decks.
- Web app v2 routes for home, lessons, practice, tutor, shadowing, pronunciation, words, vocabulary, phrasebook, mistakes, tools, progress, premium, referrals, settings, and leaderboards.
- One profile across web app, PWA, and Telegram.
- Free, Premium, and Platinum plans.
- Plan limits: Free has 5 lessons and 15 practice messages per day with no voice; Premium has 50 lessons, 200 practice messages, and 20 voice messages up to 30 seconds; Platinum has 100 lessons, 500 practice messages, and 60 voice messages up to 30 seconds.
- Payments through Telegram Stars, YooKassa/SBP, TON, and USDT.
- Referral rewards and activation keys can be mentioned only as secondary trust/monetization details, not as the main hero promise.

## Landing Content Design

### Navigation

Keep the current nav structure. Rename the primary app CTA from a generic "Open web app" style to a stronger action, such as "Начать бесплатно" or "Пройти бесплатный урок". Keep the Telegram path available as a secondary CTA in the hero and footer.

### Hero

Hero must be full-screen or first-viewport dominant, with the old animated hero background visible across the viewport. The copy should lead with the selected goal-based promise:

- H1 direction: "Выберите цель, Poliglot AI соберет маршрут".
- Supporting copy: state that travel, work, exam preparation, and conversation practice can combine AI Tutor, roleplay, voice, photo translation, vocabulary, mistakes, and progress in one profile.
- Primary CTA: "Выбрать цель и начать бесплатно".
- Secondary CTA: "Открыть Telegram-бота".
- Goal selector: show compact choices such as "Путешествия", "Работа", "Экзамен", and "Разговорная речь".
- Proof chips: keep `35 языков интерфейса`, `A1-C2`, and `web + PWA + Telegram`; add or rotate in `Free старт` if space allows.
- Demo card: keep the interactive tabs, but change them to demonstrate a goal route: "Урок", "Диалог", "Голос", "Фото". The demo should read like a mini course flow, not a generic app preview.

### Course Strip

Keep the course strip block. Reframe the cards around outcomes instead of internal feature names:

- "Для поездки": survival phrases, photo translation, hotel/airport/cafe roleplay.
- "Для работы": small talk, meetings, interview answers, correction and model phrases.
- "Для экзамена и речи": structured practice, pronunciation, shadowing, weak-word review.

Each card should include a clear "start this route" style CTA while still linking to `/app/`.

### Features

Keep six feature cards. Update the content to cover the actual current functionality:

- AI Tutor route: scenario, teaching point, model answer, final check.
- Roleplay and practice: real-life dialogues and answer correction.
- Pronunciation and shadowing: STT/TTS, score, weak words/sounds, tips.
- Photo and translation: text from image into translation, note, and practice.
- Vocabulary, mistakes, and offline: review, spelling, phrasebook, offline decks.
- Progress and motivation: XP, streak, awards, level, leaders, daily bonus.

### Workflow

Keep three steps. The revised flow:

1. Choose goal, language, and level.
2. Complete a short guided session in web app, PWA, or Telegram.
3. Repeat weak words, mistakes, and voice issues from the same profile.

### Pricing

Keep Free, Premium, and Platinum. Make plan value specific:

- Free: "попробовать маршрут", 5 lessons, 15 practice messages, no voice.
- Premium: 50 lessons, 200 practice messages, 20 voice messages, voice/photo features, mistakes and vocabulary.
- Platinum: 100 lessons, 500 practice messages, 60 voice messages, highest limits for intensive preparation.

Show monthly prices currently used in the landing and README: Premium 300 ₽ and Platinum 590 ₽. Mention Stars, YooKassa/SBP, TON, and USDT as payment methods in compact supporting copy.

### Reviews

Keep three review cards. Rewrite them as goal-specific stories:

- travel user who handled menu/hotel phrases with photo and roleplay.
- work user who practices interviews or meetings in Telegram and reviews mistakes in web.
- pronunciation user who improves weak words through shadowing.

Do not invent external ratings, logos, or unverifiable user counts.

### Progress

Keep the progress block. Tie it to retention and confidence:

- mistakes, notes, weak words, pronunciation history, XP, streak, awards, and leaders stay with the same profile.
- copy should say what the learner does next, not just that progress exists.

### FAQ

Keep four FAQ items. Update answers to include:

- Telegram-only use is possible.
- web app and Telegram share one profile.
- what Free/Premium/Platinum change.
- PWA/offline decks and saved mistakes/phrasebook.

### Final CTA

Keep the start panel. Use the same goal-based CTA as hero:

- heading: "Выберите цель и пройдите первый маршрут".
- CTA: "Начать бесплатно".
- optional secondary link to Telegram if layout allows without adding a new block.

## Visual Design

The hero should use the old full-screen animated background. It must stay readable with a dark veil or contrast layer. The layout should hint at the next section on desktop and mobile, but the first viewport should clearly communicate the brand, goal selector, CTA, and live product demo.

The rest of the page should avoid a one-hue palette. Keep the current restrained blue/cyan/gold identity, but use goal cards and pricing highlights to create clearer visual rhythm. Cards remain at 8px radius or less where the existing system already uses that shape.

No new marketing-only landing section should be added. All improvements must happen inside the existing block list.

## Technical Design

Primary files to modify during implementation:

- `site-react/src/PublicSiteApp.tsx` for copy, arrays, CTA labels, plan text, and block content.
- `site-react/src/styles.css` for hero full-screen animation treatment, goal selector layout, responsive refinements, and visual hierarchy.
- `site-react/e2e/public-site.spec.ts` for updated assertions that match the revised copy and preserved block counts.
- Built output under `Сайт полиглота для бота/` after running `npm run site:build`.

Do not edit `web-react` or `web/` for this landing task.

## Verification

Run:

```bash
npm --prefix site-react run build
npm --prefix site-react run e2e
```

Then use Playwright/browser verification for:

- desktop viewport around 1440x900.
- mobile viewport around 390x844.
- hero animation canvas or animated visual is present, visible, and not blank.
- hero text and CTA are readable over the animation.
- no overlapping text in hero, pricing, course cards, feature cards, FAQ, and final CTA.
- preserved block counts: 8 language cards, 3 course cards, 6 feature cards, 3 review cards, 4 FAQ cards, 3 pricing cards.
- all app CTAs still resolve to `/app/` or `/app/v2/` according to current site i18n/link normalization.

## Out Of Scope

- Changing the `/app/` web application.
- Adding new public landing sections.
- Changing payment backend behavior.
- Changing legal documents beyond text automatically affected by existing public site i18n.
- Replacing the brand identity or logo.

## Open Assumptions

- The old hero animation means the existing animated hero background already present in the public site source, not an external archived design.
- The Russian landing copy is the source text; existing `site-i18n.js` phrase translation machinery will translate visible text where matching phrase translations are available.
