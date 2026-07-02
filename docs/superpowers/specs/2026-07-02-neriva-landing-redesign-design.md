# Neriva Landing Redesign

## Context

Current landing tries to do too much on one page: product pitch, SEO answers, feature explanations, platform explanations, comparisons, pricing, and legal/commercial trust copy. This makes the page long and gives it a generic AI-generated SaaS feel.

The new direction is a shorter main landing for conversion, with long SEO content moved to dedicated pages.

## Goals

- Make the main page feel like Neriva, not a template landing.
- Reduce the main page length and keep first-screen decisions simple.
- Use real product UI as the visual anchor.
- Keep SEO coverage, but move long answers to focused pages.
- Keep payment copy compliant: Telegram Stars and YooKassa/SBP only.
- Keep social/community links current: Telegram channel, Instagram, YouTube, TikTok.

## Main Page Structure

1. Hero
   - Brand signal: Neriva.
   - One concrete promise: language practice with AI tutor in Web app and Telegram.
   - Primary actions: Web app and Telegram.
   - Visual anchor: real UI fragment or screenshot, not abstract glass cards.

2. Live Product Scenario
   - Show how one short session works:
     - learner starts from a weak phrase or topic;
     - AI tutor explains;
     - learner answers by voice/text;
     - system saves mistakes and words for review.
   - This block should feel like a product walkthrough, not a marketing feature grid.

3. Four Key Functions
   - AI tutor.
   - Speaking/voice practice.
   - Photo translation.
   - Mistakes and spaced repetition.
   - Each item needs one specific user outcome, not generic AI claims.

4. Pricing
   - Free entry plus paid limits.
   - Payment methods: Telegram Stars and YooKassa/SBP.
   - No RollyPay, TON, USDT, crypto wallets, or blockchain copy.

5. Giveaway / Community
   - Mention monthly premium key giveaways as channel activity.
   - Link Telegram channel: https://t.me/neriva_app.
   - Keep Instagram: https://www.instagram.com/neriva.ru.
   - Keep YouTube: https://www.youtube.com/@neriva_app.
   - Keep TikTok if already present.

6. Short FAQ
   - 4-6 questions max.
   - Focus on Web app vs Telegram, pricing, data/privacy, and languages.

## Move Off The Main Page

These blocks should not be on the main landing by default:

- "AI tutor vs vocabulary app" comparisons.
- Long SEO answers.
- Big Web app / Telegram explanations.
- Long "how the system works" blocks.
- Repeated feature cards with similar copy.

## SEO Pages

Create dedicated pages:

- `/ai-tutor`
- `/speaking-practice`
- `/pronunciation`
- `/photo-translation`
- `/telegram-language-bot`

Each page should have:

- one strong page-specific H1;
- short product intro;
- 3-5 real product sections;
- relevant screenshots/UI fragments;
- FAQ;
- links back to Web app and Telegram.

## Visual Direction

Use a Swiss/editorial product direction:

- white/black base;
- strict grid;
- fewer cards;
- one accent color;
- real screenshots and UI fragments;
- large but controlled typography;
- no glassmorphism-heavy layout;
- no decorative gradient blobs;
- no generic "AI SaaS" neon/purple styling.

The page should look designed around the actual app, not around generic feature cards.

## Copy Direction

Russian and English text should be shorter and more concrete.

Avoid:

- "unlock your potential";
- "revolutionary AI-powered platform";
- "seamless learning journey";
- generic AI tutor claims without product detail.

Prefer:

- "Разберите фразу, ответьте голосом, получите исправления и повторите ошибки позже.";
- "Start a short lesson, answer by voice, and save mistakes for review.";
- product-specific nouns: AI tutor, voice answer, saved mistake, weak word, photo translation.

## Screenshot Inputs

Use these existing/user-provided assets:

- Telegram app menu screenshot: `C:\Users\Admin\Pictures\Скрины приложения\Telegram APP.PNG`
- Current generated web screenshots: `tmp\landing-screenshots\webapp`

Needed later from product if screenshots are missing:

- AI tutor lesson screen.
- Voice/speaking practice screen.
- Photo translation result.
- Mistakes/repetition screen.
- Premium/pricing screen.

## Implementation Notes

- Keep current legal routes and generated legal documents.
- Keep Yandex Metrika disclosure in legal/cookie/privacy copy.
- Main landing should link to SEO pages, but SEO pages should carry long explanations.
- Existing public site output directory remains `Сайт полиглота для бота`.

## Acceptance Criteria

- Main page is materially shorter than current version.
- Main page does not contain long SEO comparison blocks.
- Payment copy mentions only Telegram Stars and YooKassa/SBP.
- No public landing/legal copy mentions RollyPay, TON, USDT, crypto wallets, or blockchain payments.
- The first screen clearly says Neriva and gives Web app / Telegram actions.
- Separate SEO pages exist for the listed topics.
- Visual design uses real product UI and avoids generic AI SaaS styling.
