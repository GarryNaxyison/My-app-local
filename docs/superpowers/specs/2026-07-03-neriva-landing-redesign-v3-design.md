# NERIVA Landing Redesign V3 - Spec It

Date: 2026-07-03

## Decision

Rebuild the landing design before another implementation pass. The current page reads as AI-generated because it combines oversized screenshots, weak section concepts, empty proof boxes, abrupt rectangular backgrounds, and generic product copy.

This spec replaces the previous V2 assumptions where they conflict:

- Main landing supports Russian and English only for now.
- Other 35 interface languages are paused until the final landing structure is approved.
- Light/dark theme switching is removed from the public landing for now.
- The landing ships as a polished light-only product page.
- Dark mode is out of scope for this redesign and must not be implemented in the next pass.

## Current Problems To Fix

- In light theme, some headings and labels are unreadable, including "AI-репетитор в вебе и Telegram".
- Under the hero CTA buttons, the three proof cells render as empty outlined rectangles.
- Product screenshots in several sections are too large; the viewer sees a blurred crop instead of useful UI evidence.
- The "Четыре функции в одном профиле" block is too primitive and must become a controlled product carousel.
- Section transitions are too hard; large background rectangles expose block edges and make the page look assembled from templates.
- Pricing must not use a screenshot as the main visual. It needs a designed tariff card system.
- Telegram, mobile, community, and final CTA blocks currently look weak and generic.
- The final CTA is effectively a black screen in dark mode; removing theme switching avoids this until design is stable.
- Copy is too wordy in places and too generic in others.

## Reference Research

Use these as design references, not as copy-paste dependencies:

- 21st.dev: source for crafted React/Tailwind blocks and non-generic marketing layout density.
- 21st.dev carousel components: reference for controlled feature/product carousel with nav and indicators.
- 21st.dev testimonials/community components: reference for social proof and community blocks.
- Magic UI: reference for restrained motion accents and animated attention, not full-page effect overload.
- Aceternity UI: reference for sticky reveal, product carousel, pricing sections, and tactile motion.
- shadcn/ui and shadcn blocks: reference for accessible primitives, pricing cards, tabs, buttons, and clean layout discipline.

GitHub sources checked:

- `serafimcloud/21st`
- `21st-dev/magic-mcp`
- `magicuidesign/magicui`
- `birobirobiro/awesome-shadcn-ui`
- Aceternity React community ports were found, but they are not strong enough to adopt directly.

Skill decision:

- Use `ui-ux-pro-max` for UX/design-system checklist.
- Use `ckm:ui-styling` for shadcn/Radix/Tailwind-style component thinking.
- Use `frontend-design` as a visual anchor guardrail against generic AI SaaS output.
- Do not rely on the missing `ui-ux-pro-max` CLI scripts in this local install; the scripts directory is not available as a runnable search tool.

## Visual Direction

Direction: editorial product page with a crafted component system.

Base:

- One light theme only.
- White and near-white surfaces, black text, neutral gray secondary text.
- One confident accent blue, used sparingly for primary CTAs and active states.
- No dark theme, no theme toggle, no black final CTA.

Layout:

- Product-led editorial grid, not glass SaaS cards.
- Wide desktop uses asymmetric composition: strong copy column + intentionally framed product UI.
- Mobile uses a separate composition, not a squeezed desktop layout.
- Section transitions use whitespace, subtle hairlines, overlapping visual rhythm, or soft neutral bands; avoid visible hard-edged blocks.

Motion:

- Keep the existing light hero animation, but it must stay transparent and behind the product/copy.
- Motion must not create a visible rectangular canvas background.
- Use 1-2 additional motion accents only if they support carousel/product state changes.

## Main Landing Structure

### 1. Hero

Goal: in the first viewport, understand what NERIVA is and choose Web app or Telegram.

Content:

- H1: short RU/EN product statement.
- One short subline: lesson, answer, correction, repeat.
- CTAs: "Веб-приложение" and "Telegram".
- Product proof: no empty boxes. Replace current three proof cells with either:
  - compact inline stats with filled text, or
  - a real mini status strip from the dashboard.

Visual:

- Current light hero animation stays.
- Dashboard screenshot must be framed so key UI is legible.
- Avoid huge screenshot scale where only chrome/top bars are visible.

### 2. Live Lesson Scenario

Goal: show the actual learning loop before feature marketing.

Content:

- One scenario: travel or work phrase.
- Steps: phrase -> user answer -> correction -> repeat tomorrow.
- Keep copy human and short.

Visual:

- Use a composed lesson mockup, not a raw full screenshot.
- Crop the AI tutor screenshot into 2-3 meaningful fragments: task, user answer, correction.
- Add small labels only where they clarify the loop.

### 3. Four Functions Carousel

Replace the primitive four-card grid with a controlled product carousel inspired by 21st.dev / Aceternity / Magic UI patterns.

Functions:

- AI tutor: task, answer, correction, explanation.
- Voice: score, weak words, shadowing/audio review.
- Photo: OCR, menu/sign/task translation, extracted phrase.
- Mistakes: saved errors, corrected phrase, next review.

Behavior:

- Desktop: one active feature at a time, with side navigation or segmented tabs.
- Active panel shows one clear product visual and one concise explanation.
- Inactive features stay visible as compact selectors, not full equal cards.
- Mobile: horizontal snap carousel or compact tabs; no long stacked screenshot tape.

Visual:

- Screenshots are cropped and masked into product fragments.
- Each slide must show the exact useful area of the UI.
- Do not show a giant full desktop screenshot if the text inside becomes unreadable.

### 4. Pricing

Replace screenshot-led pricing with designed tariff cards.

Plans:

- Free
- Premium
- Platinum

Rules:

- Cards must show price, daily/monthly limits, and best-for line.
- Premium can be visually emphasized if that is the recommended default.
- Payment line: "Telegram Stars и YooKassa/SBP".
- No pricing screenshot as the main visual.
- Optional small "limits meter" or "access chips" can make the cards feel product-specific.

### 5. Telegram Companion

Rebuild the weak Telegram block as a product companion module.

Message:

- Telegram is a quick practice channel, not a second product.
- Same profile continues in the Web app.

Visual:

- Use the light Telegram screenshot as a phone/chat frame.
- Add 2-3 action chips: "Голос", "Фото", "Напоминание".
- Avoid generic paragraph-heavy copy.

CTA:

- "Открыть Telegram".

### 6. Mobile Web App

Rewrite and redesign the "На телефоне всё короче" block.

Message:

- Mobile is for quick check-ins, not a compressed desktop.
- Suggested RU copy direction: "Короткая практика с телефона" or "Один экран для прогресса и повторения".

Visual:

- 2 or 3 phone frames maximum.
- Show progress, lesson/correction, and mistakes or voice.
- No long vertical tape of screenshots.

### 7. Community / Giveaways

Replace the weak "Розыгрыши и новости продукта" block.

Goal:

- Make it feel like a living product channel, not a footer ad.

Possible design:

- Editorial band with channel preview, three recent content chips, and Premium giveaway note.
- Social icons/links stay, but they must not be the whole design.

Copy:

- Short and concrete.
- Example RU direction: "Канал NERIVA: обновления, мини-уроки и розыгрыши Premium".

### 8. Final CTA

Replace the black-screen CTA with a light final action band.

Content:

- H2: "Начните первый урок".
- Subline: one sentence.
- Buttons: Web app, Telegram.
- Optional micro-flow: "Фраза -> ответ -> правка -> повтор".

Visual:

- Light surface only.
- Can use an abstracted mini lesson card, not a screenshot.
- No black background.

### 9. FAQ

Keep short.

Topics:

- Web app vs Telegram.
- Free/Premium/Platinum.
- Voice/photo practice.
- Data/cookies in plain language.

## SEO / AEO Pages

Keep separate SEO/AEO pages, but only Russian and English:

- `/ai-tutor`
- `/speaking-practice`
- `/pronunciation`
- `/photo-translation`
- `/telegram-language-bot`

Long explanations, comparisons, and "how it works" content stay there, not on the main landing.

## Language Rules

For this redesign phase:

- Main landing: RU and EN only.
- Language selector must expose only Russian and English on the landing.
- Do not delete the product's broader app language support.
- Do not generate 35 landing translations until the final visual/copy structure is approved.

## Copy Rules

Russian is primary; English mirrors it.

Use:

- short lesson;
- answer;
- correction;
- weak words;
- voice;
- photo;
- saved mistakes;
- repeat;
- Web app;
- Telegram.

Avoid:

- "seamless";
- "unlock potential";
- "AI-powered platform";
- long generic marketing paragraphs;
- duplicate explanations between main landing and SEO pages.

## Mobile Rules

Mobile is a separate layout.

- Hero must fit cleanly within the first screen without giant CTAs.
- Feature carousel must use snap/tabs and one visible feature panel.
- Hide or compress secondary proof.
- No horizontal overflow.
- No screenshot stack that makes the page feel endless.
- Keep touch targets at least 44px.

## Component / File Direction For Future Plan

The implementation plan must split the landing instead of keeping one large component:

- `LandingHero`
- `LessonScenario`
- `FeatureCarousel`
- `PricingCards`
- `TelegramCompanion`
- `MobilePreview`
- `CommunityBand`
- `FinalCta`
- `LandingFaq`
- `landingCopy` for RU/EN only

Do not implement this spec in the spec-writing step. Create an implementation plan after approval.

## Acceptance Criteria

- Landing has one light theme and no visible theme toggle.
- Main landing exposes only RU/EN language options.
- Hero label and H1 are readable in light mode.
- No empty proof rectangles under hero buttons.
- Hero animation remains transparent and does not create a canvas background block.
- Screenshots are legible because they are cropped/framed as product evidence.
- Four functions are presented through a polished carousel/showcase, not a generic four-card grid.
- Pricing is a designed tariff-card section, not a screenshot section.
- Telegram block clearly presents the bot as a companion channel and looks product-specific.
- Mobile block is compact and has stronger copy.
- Community/giveaway block feels like a product channel module, not a weak footer ad.
- Final CTA is light, readable, and not a black screen.
- Section transitions are visually smooth; no abrupt rectangle edges.
- Removed long explanations stay on RU/EN SEO pages.
- E2E checks cover RU/EN, no dark theme toggle, no 35-language selector on landing, no horizontal overflow, and hero proof not empty.
- Visual review screenshots are required before deploy.
