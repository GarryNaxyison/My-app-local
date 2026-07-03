# Neriva Landing Redesign V2 - Spec It

Date: 2026-07-03

## Decision

Do not make the landing "prettier" in the same structure. Change the structure.

The landing becomes a short product-first page. Long comparison, explanation, and SEO/AEO content moves into separate crawlable pages.

No implementation, public deploy, domain/DNS, Cloudflare, or email-routing work is part of this spec. The RU domain flow does not need Cloudflare for this task.

## Guardrails

- Review product screenshots before design implementation.
- Use light-theme product screenshots for the main visual system.
- Do not deploy until the landing screenshots are reviewed and explicitly approved.
- Do not add Cloudflare setup, proxying, DNS changes, Workers, or Pages changes.
- Keep removed long blocks reachable through SEO pages; do not delete their value.
- Keep test data realistic and anonymous: no tokens, email addresses, Telegram IDs, or private account data.

## Main Landing Structure

1. Hero
   - Brand: Neriva.
   - One direct explanation of the product: AI language practice in Web app and Telegram.
   - Two primary actions: Web app and Telegram.
   - First visual anchor: real product UI, not abstract AI decoration.

2. Live Product Scenario
   - Show one realistic session, for example travel or work English.
   - Flow: user phrase -> tutor correction -> voice/text answer -> mistake saved for review.
   - This should feel like a real lesson, not a generic feature grid.

3. Four Key Functions
   - AI tutor: task, answer, correction, explanation.
   - Voice/pronunciation: score, weak words, shadowing/audio practice.
   - Photo translation: OCR, menu/sign/task translation, useful phrase extraction.
   - Mistakes/repetition: saved errors, corrected phrases, weak words, review loop.

4. Pricing
   - Free, Premium, Platinum.
   - Show limits clearly.
   - Payment copy: Telegram Stars and YooKassa/SBP only.

5. Giveaway / Channel / Community
   - Telegram channel/community block.
   - Premium key giveaway can stay as a compact trust/community module.

6. Short FAQ
   - 4-6 questions max.
   - Focus: Web app vs Telegram, pricing, languages, privacy/test data, how review works.

## Move From Main Landing To SEO Pages

These blocks must leave the main landing body:

- "AI tutor vs vocabulary app" comparisons.
- Long SEO/AEO answers.
- Large Web app / Telegram explanations.
- Long "how the system works" sections.
- Repeated generic AI SaaS feature-card clusters.

## SEO Pages

Create or keep focused pages:

- `/ai-tutor`
- `/speaking-practice`
- `/pronunciation`
- `/photo-translation`
- `/telegram-language-bot`

Static SEO/AEO pages are Russian and English only.

Suggested content split:

- `/ai-tutor`: tutor loop, correction, explanation, comparison with simple vocabulary apps.
- `/speaking-practice`: roleplays, text/voice answers, daily speaking scenarios.
- `/pronunciation`: scoring, weak words, shadowing, audio examples.
- `/photo-translation`: OCR, menus, signs, worksheets, phrase extraction.
- `/telegram-language-bot`: why Telegram exists, quick start, reminders, Web app link/sync.

Each SEO page needs:

- page-specific H1;
- compact product intro;
- screenshots or UI fragments;
- FAQ;
- links back to Web app and Telegram;
- RU/EN canonical and hreflang only.

## Translation Rule

Main landing copy and SEO copy use different translation rules.

- Main landing: visible marketing/product copy needs human-edited creative translations for 35 interface languages.
- SEO/AEO pages: Russian and English only.
- Do not push long SEO copy through the 35-language static SEO generator.
- Do not let machine translation overwrite curated landing voice.

## Visual Direction

Direction: Swiss/editorial product landing.

- Base: white, black, neutral grays.
- Accent: one strong blue, for example `#002FA7`.
- Layout: strict grid, large but controlled headings, left-aligned rhythm.
- Details: 1px rules, high contrast, fewer containers.
- Product UI: real screenshots and cropped interface fragments.
- Cards: only where they represent real repeated items or pricing; no nested card stacks.

Avoid:

- generic dark "neural SaaS" look;
- glassmorphism-heavy sections;
- purple/blue gradient overload;
- decorative blobs/orbs;
- abstract AI illustrations replacing product evidence;
- long identical feature cards.

## Product Screenshot Inputs

Current captured PNG set:

- `tmp/product-screenshots/landing-assets-20260703/01-dashboard-progress-desktop.png`
- `tmp/product-screenshots/landing-assets-20260703/02-ai-tutor-lesson-correction-desktop.png`
- `tmp/product-screenshots/landing-assets-20260703/03-mistakes-review-desktop.png`
- `tmp/product-screenshots/landing-assets-20260703/04-voice-pronunciation-score-desktop.png`
- `tmp/product-screenshots/landing-assets-20260703/05-photo-translation-desktop.png`
- `tmp/product-screenshots/landing-assets-20260703/06-notes-phrasebook-desktop.png`
- `tmp/product-screenshots/landing-assets-20260703/07-premium-plans-desktop.png`
- `tmp/product-screenshots/landing-assets-20260703/08-telegram-bot-demo-mobile.png`
- `tmp/product-screenshots/landing-assets-20260703/08b-telegram-real-dark-reference.png`
- `tmp/product-screenshots/landing-assets-20260703/09-mobile-home-progress.png`
- `tmp/product-screenshots/landing-assets-20260703/10-mobile-lesson-correction.png`
- `tmp/product-screenshots/landing-assets-20260703/11-mobile-mistakes.png`
- `tmp/product-screenshots/landing-assets-20260703/12-mobile-voice-pronunciation.png`

Screenshot usage:

- Hero / early product proof: dashboard + lesson correction.
- Live product scenario: AI tutor lesson correction.
- Feature section: mistakes, voice, photo, notes.
- Pricing section: Premium plans.
- Telegram section: use the real Telegram dark screenshot only as a reference or temporary proof. For final visual consistency, request or capture a light-theme Telegram screenshot if available; otherwise keep a clearly styled product frame that does not pretend dark UI is light.
- Mobile section: home, lesson, mistakes or voice.

## Copy Tone

Russian primary, English alternate.

Use concrete product language:

- short lesson;
- answer by voice;
- correction;
- weak words;
- saved mistakes;
- photo translation;
- Web app;
- Telegram bot.

Avoid:

- "revolutionary AI-powered platform";
- "unlock your potential";
- "seamless journey";
- generic AI claims that do not show the product.

## Acceptance Criteria

- The main landing is materially shorter than the previous long SEO-heavy page.
- First screen clearly says Neriva and offers Web app / Telegram.
- A real product scenario appears before generic feature summaries.
- The four key functions are visible and tied to product screenshots.
- Pricing shows Free / Premium / Platinum and limits.
- Main page does not include long comparison or SEO answer blocks.
- Removed content remains available on the five SEO pages.
- SEO pages are RU/EN only.
- Main landing translation layer supports 35 human-edited interface languages.
- Visual design is light, editorial, product-led, and avoids generic AI SaaS/glass styling.
- Screenshot review happens before implementation.
- No deployment happens until visual screenshots are reviewed and approved.
- No Cloudflare work is introduced by this redesign task.
