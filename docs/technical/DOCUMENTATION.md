# Project Documentation

## Agent MCP Startup

General MCP startup rules are documented in `MCP_STARTUP.md`. At project start, verify Figma, 21st.dev Magic, Ref, shadcn, Nx MCP, and Linux/Ubuntu access when server work may be needed.

Nx MCP can be kept available for local tools as a standalone Streamable HTTP server:

```powershell
npx -y nx-mcp@latest "C:\Users\Admin\Documents\New project" --transport http --port 9921 --no-minimal --disableTelemetry
```

The server should log `Nx MCP server (Streamable HTTP) listening on port 9921`. The stdio form from Codex config, `npx nx mcp --workspacePath "C:\Users\Admin\Documents\New project" --no-minimal --disableTelemetry`, exits immediately when launched manually without an MCP client attached.

## Current Regression Contract

Before deploying changes to the public site or V2 app, run the focused suite that protects the latest reported regressions:

```powershell
go test ./...
npm --prefix web-react run build
npm --prefix web-react run e2e
npm --prefix site-react run build
npm --prefix site-react run e2e
node tools/check_encoding_artifacts.mjs
```

The web app Playwright suite covers desktop and mobile separately for login/auth, mobile onboarding scroll, mobile composer placement, Roleplay, Tools, Mistakes, Awards, payment modal placement, browser Back behavior, localized labels, confirmed-only payment history, referral pagination, Offline pagination/export, Phrasebook persistence, login canvas animation, and Learn Words option variation. The public-site suite covers the landing logo, combined hero shader, `/app/` app links, theme toggle, 35-language copy checks, discounted pricing, light-theme section label readability, reviews/FAQ, and dark-theme legal page readability.

## Web And Telegram Interface Localization

The React web app and Telegram bot support the same 35 interface languages: Russian, English, Spanish, German, French, Italian, Chinese, Japanese, Korean, Tajik, Uzbek, Tatar, Armenian, Kazakh, Kyrgyz, Georgian, Ukrainian, Polish, Romanian, Portuguese, Arabic, Bengali, Czech, Greek, Hindi, Hungarian, Indonesian, Dutch, Swedish, Tamil, Telugu, Thai, Tagalog, Turkish, and Vietnamese.

Interface copy is static product UI text and must stay in the repository. Do not use OpenRouter or another runtime model for web or Telegram labels. OpenRouter remains only for guarded missing vocabulary translation cache warming described in the vocabulary section.

React V2 localization lives in `web-react/src/lib/i18n.ts`. Rich hand-written copy exists for the original language set, while the 35-language pass now includes derived localized UI copy and a guard that replaces exact English fallback strings for required non-English labels. Playwright test `v2 required labels are localized for all 35 interface languages` checks these required labels on desktop and mobile.

Web auth/register system errors use stable backend error codes such as `invalid_login_start`, `invalid_login_chars`, `login_taken`, `invalid_credentials`, and `captcha_required`. The standalone React login page maps those codes to `auth_error_*` copy via `appCopy`, so user-facing errors stay localized and must never display generic fallback text like `Р Р°Р·РґРµР»`, `Section`, or `Mб»Ґc`. Login normalization allows dots inside a login (`friend.name`) but rejects a leading dot, hyphen, or underscore.

Telegram localization lives in `i18n.go`. New interface languages must use local `uiCopy`/generated compact copy and `toolUICopy` values, never `uiAliases` to English. Go test `TestEveryNonEnglishInterfaceLanguageAvoidsEnglishTelegramFallback` checks that every non-English interface language has localized core menu and tool labels.

Level assessment prompts must keep the question/task chrome in the interface language for all 35 interface languages, while answer options and target examples remain in the learning language. Every supported learning language must have a static 36-question exam-style set in `levelAssessmentQuestionsByLanguage`; do not route level determination through `vocabularyLevelAssessmentQuestions`, because that becomes word guessing rather than a level test. The backend maps source question types such as fill-gap, translation, meaning, naturalness, and grammar-form checks to localized task labels; do not collapse them back to the generic `Choose answer` copy. Go tests `TestLevelAssessmentQuestionsExistForEveryLearningLanguage` and `TestLevelAssessmentTaskTypesLocalizedForEveryInterfaceLanguage` protect this contract.

Learn Words/Review result cards must keep the learner-facing explanation in the interface language while still showing the target-language word that was just learned. On a correct Learn Words answer, the success card shows both the interface-language translation and the learned target word, then exposes the same Notes quick-save chip pattern under `Next`; the saved note stores the target word as the phrase and the interface-language translation as the note/context. Playwright test `learn words result shows target word and saves word pair to notes` protects this behavior.

## Public Landing and Legal Localization

The public landing source is `site-react/src/PublicSiteApp.tsx`; shared public styles live in `site-react/src/styles.css`. The landing header logo is served from `site-react/public/assets/brand-logo-mini.png` and must stay in sync with the generated app brand assets. The first viewport uses one combined shader/product hero, not two unrelated hero animations.

The current landing follows a Busuu-style product flow without copying Busuu branding: a simple promise, `РЇ С…РѕС‡Сѓ РёР·СѓС‡Р°С‚СЊ` language chooser, course cards, proof stats, concise feature blocks, reviews, FAQ, and pricing. The content must stay about Poliglot AI features: web app, PWA, Telegram, daily route, Roleplay, Pronunciation, Photo/translation, Notes, mistake dictionary, and A1-C2 learning levels.

Poliglot AI now has 35 interface and learning languages for the public site, React web app, Telegram selector, and local vocabulary dictionaries: Russian, English, Spanish, German, French, Italian, Chinese, Japanese, Korean, Tajik, Uzbek, Tatar, Armenian, Kazakh, Kyrgyz, Georgian, Ukrainian, Polish, Romanian, Portuguese, Arabic, Bengali, Czech, Greek, Hindi, Hungarian, Indonesian, Dutch, Swedish, Tamil, Telugu, Thai, Tagalog, Turkish, and Vietnamese. Do not expose a new language as a learning option until its `vocabulary_words*.json` file is generated and covered by the all-language vocabulary regression.

Translations for landing, Privacy, and Terms are generated by:

```powershell
node tools\rebuild_public_site_translations.mjs
```

The generator scans React source plus `site-react/src/legacyLegalContent.ts`, writes `site-phrases.js`, and mirrors it into both the React public assets and the legacy static site assets. The browser-side decoder in `site-react/public/assets/site-i18n.js` must only repair actual mojibake; it must not decode normal accented Latin text such as French, German, Spanish, Polish, or Romanian.

Legal pages must keep the same contact/header treatment as Terms, including `@AsaselD`, email/contact fields, and dark-theme-safe language selectors. Any new visible public copy must be checked for all 35 supported interface languages and then validated with `node tools/check_encoding_artifacts.mjs`.

## Vocabulary SQLite Performance

Learn Words and Review answer choices are generated by the backend and shared by web desktop, web mobile, and Telegram. The authoritative fix is in Go/SQLite, not in React. Distractors should be sampled from a wide random range in the same CEFR band as the correct word: `A1/A2`, `B1/B2`, or `C1/C2`.

The current SQLite path restores the old wide in-memory random behavior by caching arrays of matching word IDs. `Learn Words` picks a random eligible ID from the cached array after filtering learned IDs. `Review` keeps the correct answer inside the user's personal review pool, but samples wrong answers from cached SQLite ID arrays for the same CEFR band while excluding the current word and the user's learned/review IDs. Do not reintroduce the old `random position -> fallback to position 0` behavior; sparse early pools will repeat the first few words again.

The production SQLite vocabulary database must keep indexes for language/level/rank/position word lookup and language/word translation lookup. The current production check confirmed `vocabulary_words` and `vocabulary_translations` are indexed for the randomized option path, so do not replace the cached ID-array path with broad `ORDER BY RANDOM()` scans.

Missing interface-language prompts must prefer the learner's selected interface language over Russian or English fallback. The backend first checks the local dictionary and `vocabulary_ai_translations`; if the exact word/interface pair is still missing and OpenRouter is configured, it performs a guarded synchronous lexicographic translation with the `vocabulary.translation.*` prompt, persists the result in `vocabulary_ai_translations`, and mirrors it into `vocabulary_translations`. Future requests for that word/language pair read from SQLite and do not call AI again. The nearest dictionary fallback, usually English or Russian, is only a temporary error/configuration fallback when the model path is unavailable or fails.

Vocabulary prompts must never reveal the hidden answer. A dictionary value that equals or contains the target headword, including bad rows such as English `creating` with Russian gloss `creating`, is treated as missing and must go through the SQLite/OpenRouter translation path instead of being shown to the learner.

Successful AI vocabulary translations are also written back to the matching source JSON dictionary (`vocabulary_words*.json`) under `translations[interface_language]`; Russian translations additionally update the legacy `russian` field. Writes are serialized and atomic through a temporary file plus rename. During SQLite JSON import, existing `vocabulary_ai_translations` are replayed into `vocabulary_translations`, so cached AI work survives either a normal runtime lookup or a SQLite rebuild.

Learner CEFR level is stored in the backend profile and applies to web and Telegram word selection. New Learn Words entries try the exact selected level first; if that exact pool is exhausted, fallback remains inside the paired band only: `A1/A2`, `B1/B2`, or `C1/C2`. Wrong-answer distractors use the same paired-band rule and must not backfill from another CEFR band.

Mobile web navigation uses the horizontal rail as the primary menu. Long-press starts a one-shot drag, not a persistent edit mode: the tile lifts under the finger, movement immediately reorders the real rail items, neighboring tiles visibly make room, the rail auto-scrolls near edges, and releasing saves the order in `poliglot-mobile-nav-rail-v2:{accountKey}` and exits drag mode automatically without a Done/X control.

The AI translation prompt must stay lexicographic rather than conversational: strict JSON only, translate the hidden target-language lemma into the learner's native/interface language, match part of speech/CEFR/topic/existing glosses, return up to three short natural equivalents separated by `; `, and never include the target word, transliteration, examples, grammar notes, or Markdown.

## Runtime

Current routing note: older documentation below may still mention the original private V2 rollout. As of 2026-06-07, the public web shell is `/app`. `/app/mobile` and `/app/desktop` force the two supported web layouts, while `/app?shell=mobile` and `/app?shell=desktop` remain QA overrides. `/app/v1` and `/app/v2` are compatibility redirects to `/app` and must not be used as public links. `/login` and `/app/login` serve the standalone React auth shell. The built app lives in `web/index.html` plus `web/assets`; do not recreate `web/v2`.

Poliglot AI is a Go Telegram bot with one active React web-app shell served by the same Go HTTP server. The current React web app is served on `/` and `/app`; versioned `/app/v1` and `/app/v2` URLs redirect to `/app` for old links. The React source lives in `web-react`, and production web files are emitted into `web/index.html` and `web/assets` with:

```powershell
npm run web:build
```

The Go backend remains the source of truth for auth, lessons, practice, vocabulary, Premium, referrals, Phrasebook/Notes, and settings. The web client calls the existing `/api/...` endpoints; do not duplicate authoritative backend state in the browser.

The binary serves the web shell by reading `web/index.html` for `/` and `/app`, from the current working directory or from the directory beside the executable. Production keeps `/opt/aibot/web` next to `/opt/aibot/aibot` instead of baking HTML into the binary. The HTTP handler injects a forced mobile or desktop shell into the `<html>` tag before sending the shell: real mobile user agents receive `data-app-shell="mobile"`, desktop receives `data-app-shell="desktop"`, and QA can force either with `/app/mobile`, `/app/desktop`, `/app?shell=mobile`, or `/app?shell=desktop`. The React shell reads those attributes on startup. `/tools` remains available for the public website or a separate landing/tool route.

The React visual system is a Vite/Tailwind/shadcn-compatible web app. It uses a horizontally scrollable, scroll-snap image-tile ribbon at the top instead of the old sidebar/workspace/inspector prototype. Desktop shows a narrow profile/status bar with interface language, theme, logout, XP progress, and the current trophy thumbnail placed to the left of `LVL`; below it the tile ribbon can scroll left/right and the selected mode opens one large terminal-style context display that fills the remaining viewport height. Mobile keeps the same model with a compact top row, large visible tile artwork, fixed bottom navigation, and task/result panels sized for touch use. Both light and dark themes now use a stable static gradient background tuned to the current glass cockpit design; moving shader layers are disabled for this pass.

The current web QA direction is stricter than the first React pass. The background should stay static and theme-aware: no drifting shader, no extra decorative orbs, and no motion that competes with the task surface. The top function ribbon uses centered visual expansion for chips, not left-to-right width growth. The header brand button opens Settings, and the XP/trophy progress button opens Awards. Header controls include the animated solar theme toggle, a brightness slider, interface language selection, and the animated logout confirmation button.

The 2026-05-23 V2 correction pass adds a stricter interaction contract. System status messages are top-center under the header and auto-dismiss after 3 seconds. Language selectors use the shared animated dropdown everywhere they appear, including Settings, Leaderboard, and Tools, with dark-theme-safe colors. The function ribbon expands hovered items while neighboring items make room. Each selected context view drops down with a short smooth animation, and primary/icon buttons share the same hover/tap motion language. Loading states use `web-react/src/components/ui/spinner.tsx`.

The final 2026-05-23 V2 refinement tightens that contract further. Neighboring function-ribbon chips must make room without shrinking, fading, or changing their own scale. The active recording control keeps the Poliglot visual skin but follows the supplied animated pattern: stop-square, growing waveform, and elapsed timer. Payment checks and webhook/session refresh success states open centered dialogs using `web-react/src/components/ui/dialog.tsx`; success dialogs show plan, period, amount when available, and the renewed subscription expiration, while not-found checks show a centered warning. Payment instructions should tell users to transfer exactly the shown amount, network, and comment/memo and warn that mismatched payments may not pass. Awards prefer the old V1 polyglot rank/story text as the canonical fallback, and V2 panels stay slightly transparent so the static theme background remains visible in both themes.

The current mobile web implementation has its own bottom navigation layer (`MobileBottomNav` in `web-react/src/App.tsx`) instead of relying on the desktop function ribbon. Mobile hides the desktop ribbon, exposes primary destinations with larger touch targets, and puts secondary destinations behind the More sheet. Any new web view must be checked in both desktop and mobile compositions: desktop can use the wide ribbon/context display, while mobile needs a direct next action, no horizontal overflow, and dialogs that fit inside `390x844` and `430x932` viewports.

The 2026-05-23 feature slice adds real surfaces for the roadmap items that were previously only recommendations. Home now acts as the "Today" screen with a weekly learning plan, daily quests, progress meters, and next actions. `roleplay` starts one of ten AI practice scenarios and renders the result inside Roleplay instead of redirecting into the generic Practice surface. `pronunciation` is now a standalone pronunciation workout: model text/audio, user upload/record, scoring, weak sounds, and corrected sample audio. `offline` builds browser mini-decks from vocabulary and mistakes, stores them in `localStorage`, supports JSON and readable TXT export, and registers `/app/offline-deck-sw.js` plus `manifest.webmanifest` for PWA/offline shell caching. Installed PWA shells use a network-first service-worker strategy for `/app` and `/app/assets`, cache only as offline fallback, and the Go server sends no-store headers for the HTML shell, manifest, and service worker so mobile installed apps converge to the current web UI after deploy. `dashboard` starts the teacher/admin surface with current learner progress, activity, Premium/referral state, leaderboard snapshot, and problem topics; the old Progress route normalizes into this Dashboard. These views reuse existing Go API state first; deeper multi-user admin and reminder automation should be added as backend endpoints instead of storing authoritative data in the browser.

The fourth 2026-05-23 V2 audit slice makes those surfaces explicit and removes the navigation regressions that hid them. `normalizeView` defaults to `home`, and entering `offline` or `dashboard` preloads Vocabulary, Mistakes, and Leaderboard data with `navigate: false`, so background loading cannot switch the user into another mode. The Today screen includes a `V2 learning lab` strip with direct entries for AI roleplay, Pronunciation, Offline decks, and Dashboard. The first-run onboarding dialog asks for learning goal, CEFR level, target language, and preferred format, then saves the profile settings and closes. Mistakes are now grouped by grammar, word order, vocabulary, politeness, and spelling, and the `РџРѕС‚СЂРµРЅРёСЂРѕРІР°С‚СЊ РїРѕС…РѕР¶РёРµ` action starts a drill from the same group. Premium records browser-local payment history rows for created/checked/webhook-confirmed payments while the backend remains the source of truth for actual payment status. Telegram two-factor input is rendered through `otpdialog.tsx` into a global portal on `document.body`, keeping it centered above all app panels.

The V2 vocabulary/roleplay/pronunciation correction pass fixes the latest user-facing regressions. Learn Words chooses the prompt word from the selected interface language when a translation exists, so Russian UI can show the Russian prompt while the answer remains in the learning language. `/api/vocabulary` now keeps V1-style fast list loading and does not call the AI example generator for every dictionary card during page load. `ROLEPLAY_TOOL_V2` requests are routed through a dedicated backend `roleplayPrompt` and editable `roleplay.scenario.user` runtime prompt instead of the generic practice prompt. Pronunciation exact-repeat scoring caps low-similarity or missing-word attempts and blocks praise for clearly wrong repeats. Global leaderboard rows now show the language names behind the language count, and Offline decks can export a readable TXT study pack for learners.

The 2026-05-24 V2 repair pass makes message output scoped per tool. React chat messages carry a `meta` surface id, so Lesson, Practice, Roleplay, Listening, Learn Words, Review, Spelling, Level, Mistakes, and each Tools mode render only their own relevant messages instead of leaking the whole shared history into every output window. When adding a new tool or menu surface, tag both user-side and assistant-side messages and filter by that tag before passing data to `ChatPanel`.

The 2026-05-26 V2 learning slice adds a personal Phrasebook, habit calendar, daily bonus, bug report flow, and richer Offline grouping. Phrasebook data is now persisted through `GET/POST/DELETE /api/phrasebook`, returned in `/api/session`, and stored in `phrasebook_entries` for SQLite or `phrasebook` in the JSON profile; browser `localStorage` is only an optimistic cache/fallback per account. Habit data is browser-local in `poliglot-habit-v2:{accountKey}` and records login/completion status for the monthly calendar; server-side duplicate protection for the XP bonus lives behind `POST /api/daily/claim`, which stores claimed dates in `daily_bonus_claims` for SQLite and `daily_bonus_claims` in the JSON profile. Bug reports are authenticated multipart posts to `POST /api/bug-report`; the server writes a JSONL index to `bug_reports/bug_reports.jsonl` and screenshots to `bug_reports/screenshots`. The Offline view can group all cards, vocabulary, saved phrases, or mistakes, and TXT export includes learner-facing fields instead of raw JSON.

The web localization layer in `web-react/src/lib/i18n.ts` now exposes 35 interface locale codes. The original 20-language override layer still contains richer hand-localized copy, and the 15 newly added UI languages have targeted core labels plus English fallback for secondary text until full hand translation is added. Keep all overrides in normal UTF-8 and expand them whenever a new menu destination, card title, or system button is introduced.

The 2026-05-27 auth/privacy pass keeps the standalone login screen in `web-react/src/App.tsx` backed by the shadcn-compatible `web-react/src/components/ui/sign-in.tsx` component. `/login` and `/app/login` must preserve `/api/auth/login`, `/api/auth/register`, Cloudflare Turnstile token submission, Telegram six-digit OTP login, and the shared centered OTP popup. Registration and Telegram linking require the personal-data policy checkbox and submit `privacy_consent` plus the current policy URL. Account recovery opens `https://t.me/AsaselD`. Login hero artwork is generated by the existing ComfyUI FLUX.2 workflow and served from `web/assets/auth-login-hero-light.png` and `web/assets/auth-login-hero-dark.png`.

The public Privacy policy is generated into both the React site and the legacy static asset path by `tools/update_privacy_policy_assets.mjs`. It updates `site-react/src/legacyLegalContent.ts`, `site-react/public/assets/privacy-policy-i18n.js`, and `РЎР°Р№С‚ РїРѕР»РёРіР»РѕС‚Р° РґР»СЏ Р±РѕС‚Р°/assets/privacy-policy-i18n.js`. Keep the generated policy text translated for all 35 interface languages and keep the Telegram contact `@AsaselD` plus bot link `@poliglot_ai_bot` visible. First-time Telegram `/start` users receive the privacy onboarding prompt before the usual profile onboarding.

Vocabulary answer options are generated server-side and are shared by web desktop, web mobile, and Telegram. `wordOptions` and the SQLite option query should sample distractors from the same CEFR band as the correct answer: `A1/A2`, `B1/B2`, or `C1/C2`, excluding duplicate visible words and legacy duplicate IDs. Do not fix Learn Words or Review answer repetition only in React; the backend path is authoritative.

Mobile award details are portal-rendered to `document.body` and use a fixed viewport modal. This is intentional: after a user scrolls the Awards page and taps a trophy, the details card must appear in the current visible area, not at the top of the scrolled content.

Mobile V2 now hides the desktop top bar and function ribbon. The phone layout uses fixed top-right quick controls for report bug, interface language, theme, and icon-only logout, a fixed bottom navigation row for the five primary sections, and a bottom-right More tab that opens the remaining sections as a compact sheet with the same menu artwork as desktop. The bottom nav is user-editable: long-press enters edit mode, X removes a pinned item into More, tapping a More tile in edit mode pins it back, and the Done control exits editing. The main mobile content scrolls through `.context-display`, and nested view panels should stay `height: auto` / `overflow: visible` unless a specific component requires its own internal scroll.

The 2026-05-26 second V2 cleanup pass adds stricter state rules for payments, referrals, Offline decks, and mobile input surfaces. Premium payment history is browser-local but must only store confirmed/paid payments; created and pending invoices are filtered out on read and are not written when a payment window opens. The Referrals view consumes `user.referral_invitees` from `/api/session` and paginates invited learners 10 per page with joined date, level/XP, level-3 status, and earned amount. Offline decks paginate visible cards 10 per page, but TXT/JSON export must always use the full filtered deck. Mobile Lesson, Practice, Roleplay, Listening, and Tools keep their composer/input panels visible above the bottom nav. Tools uses a selector-first mobile flow: choose translator, voice-to-text, or photo translation first, then open the dedicated work window.

Pronunciation and Offline have additional presentation rules after the cleanup pass. Pronunciation menu artwork is intentionally larger and centered, while the pronunciation map should show compact issue labels and avoid repeating the same long feedback on each token. Offline decks should use the available page width on desktop with readable multi-column cards and stack cleanly on phones. Today places the streak panel beside the calendar on desktop and removes obsolete roadmap filler now that the listed features have real screens.

The latest 2026-05-26 mobile regression pass tightens several V2 rules. Mobile onboarding is height-constrained and scrolls inside the dialog so format/level/language selection and submit buttons stay reachable. Lesson/Practice Phrasebook quick-save chips sit after the composer on phones and are compact enough not to cover output. Mobile Roleplay hides the scenario brief after entry and keeps only a localized choose-another-scenario button. Pronunciation weak words/sounds and score history are cached per account under `poliglot-pronunciation-v2:{accountKey}` so a reload does not erase the dashboard. Mobile More edit mode supports reordering tiles inside More, not only moving pinned bottom-nav items. Menu PNGs for Offline, Phrases/Phrasebook, Pronunciation, and Roleplay are referenced through a versioned asset URL to force both desktop and mobile to fetch the current FLUX.2 transparent assets. Learn Words answer options are generated through the shared Go `wordOptions` path used by web and Telegram; wrong answers should vary across rounds and have backend regression coverage.

The current menu artwork for Offline, Phrases/Phrasebook, Pronunciation, and Roleplay is generated by the local ComfyUI FLUX.2 split workflow in `tools/generate_comfy_brand_assets.mjs`: `UNETLoader` uses `flux2_dev_fp8mixed.safetensors`, `CLIPLoader(type=flux2)` uses `mistral_3_small_flux2_bf16.safetensors`, and `VAELoader` uses `flux2-vae.safetensors`. Transparent icon generation must keep the chroma-key/postprocess validation path and write 512x512 PNGs plus `brand-assets-manifest.json` entries with `qa_status: generated_transparent_validated`.

Desktop/mobile behavior for the current and newly added V2 functionality is tracked in `docs/product/WEB_V2_FUNCTIONALITY.md`. Update that file whenever a V2 feature changes its layout, navigation, or main user flow.

Settings includes a user-editable global learning focus. The field is saved as `learning_focus` through `/api/settings`, stored in both JSON and SQLite stores, and returned in `/api/session`. `lessonPrompt`, `practicePrompt`, and `roleplayPrompt` receive that focus as a soft topic preference so the user can stop the app from repeating one onboarding theme such as airports or hotels without changing code.

Lesson generation also keeps anti-repeat state. The backend stores `lesson_history` with the last 10 generated lesson prompts in both JSON and SQLite storage, passes that history into `lessonPrompt`, and asks the model to avoid the same topics, situations, example sentences, and communicative tasks. If the global learning focus is empty, the prompt rotates through a CEFR-specific topic pool for the selected level instead of falling back to one fixed airport/hotel/travel theme. Practice keeps only the last 5 learner messages as short context, so it can remember the current exchange without dragging the whole conversation into every new answer.

Editable AI prompts are indexed in `docs/reference/APP_PROMPTS.md` and loaded at runtime from `app_prompts.json` by default. The path can be changed with `APP_PROMPTS_FILE`; on the server this file should sit next to the Go binary and `.env` files so prompt text can be adjusted without rebuilding the binary. Each prompt entry keeps `feature`, `tool`, `function`, `notes`, and `template` fields. The Go code falls back to the compiled prompt when a key is absent, but invalid JSON is treated as a startup error to avoid silently running with broken prompt text. Current runtime keys cover the shared coach system prompt, lesson generation, lesson feedback, practice chat, translator prompts, image translation, vocabulary hints/examples, Listening phrase/feedback, and pronunciation coach JSON.

Encoding rule for future changes: all new source, prompt, and documentation text must be normal UTF-8. Do not paste cp1251/UTF-8 mojibake into fallbacks, copy maps, prompts, changelog entries, or docs. Before build/deploy, run `node tools/check_encoding_artifacts.mjs`; if it fails, fix the source text instead of adding broad skip rules. The only allowed mojibake handling code is the escaped detector/repair implementation in `web-react/src/lib/i18n.ts`.

Regression test rule for web UI changes: update and run `web-react/e2e/web-smoke.spec.ts` before deployment. The suite must keep covering mobile and desktop separately for localization, confirmed-only payment history, referral invitee pagination, Offline pagination with full export, Telegram `Send code` visibility, pronunciation compactness and persisted history, mobile onboarding scrollability, mobile header controls, More image tiles/descriptions/reordering, editable mobile nav, Lesson/Practice/Roleplay/Tools input visibility, mistake dictionary scrolling, Phrasebook API persistence, desktop ribbon arrow bounds, current FLUX.2 asset URLs, Learn Words answer-option variation, and mojibake guards. Backend route/state changes should also have Go coverage in `web_api_feature_test.go`, `vocabulary_test.go`, or related tests.

Live production QA can be blocked by the `/login` anti-bot challenge. When that happens, local build/test/Nx checks are still valid, but do not mark authenticated live QA complete until the same browser session used by Codex reaches `/app` with a real account session.

Learning surfaces also keep old-version behavior while using V2 styling. Mistake practice has no Start button, asks before Clear, and exposes correct-answer audio. Awards use rank names and polyglot history stories instead of generic level fallback text. Vocabulary is a dictionary view, so it hides review/spelling counters and shows an example sentence when the backend provides one. Word learning and spelling expose word/example audio only at the right answer/result stage and keep the next-word action visible. Listening includes a next-phrase action. Lesson/practice audio controls should include corrected/model phrases and the AI follow-up question after "Your turn" via backend TTS text.

Web Settings must include password change for normal web accounts using the local password-input styling in `web-react/src/components/ui/password-input.tsx`, while keeping server-side password validation in Go. Telegram two-factor setup keeps the old V1 delivery behavior: pressing `Send code` calls `/api/auth/telegram/start`, opens the returned Telegram bot deep link with the `web_...` start payload, then shows a separate centered OTP/two-factor dialog layer, visually analogous to the Premium payment dialog, rather than embedding the input inside the Settings card. The dialog verifies the 6-digit code through `/api/auth/telegram/status` and resend creates a fresh Telegram auth request. Referrals must explain the 7-day invitee reward, delayed inviter reward at XP level 3, 20% direct and 5% indirect balance logic, and use the old localized Telegram invitation template where it exists. Header identity should prefer the user-created web login before Telegram display names.

Reusable V2 UI components now live under `web-react/src/components/ui`, matching the shadcn-compatible component path. Important components include `button.tsx` for joly-style hover/tap buttons and approve/reject variants, `voice-input.tsx` for all recording controls, `audio-wave-button.tsx` for backend TTS/playback rows, `animated-theme-toggle.tsx`, `interfaces-slider.tsx`, `morphing-arrow-button.tsx`, `morph-button.tsx`, `logout-button.tsx`, `input.tsx`, `label.tsx`, and `otpdialog.tsx`. Do not reintroduce browser `speechSynthesis` for voiced samples; V2 listen buttons should call the existing backend audio endpoints such as `/api/tools/translator-speech` and `/api/words/pronunciation`.

V2 payment UX shows details inside the same payment modal after a method is selected. The selected method must be visually highlighted and can stay in the morph/loading state while payment is pending. Requisites should show all fields supplied by the backend, especially payment ID, status, amount/currency, network, wallet address, memo/comment, expiration, transaction hash, and invoice/payment link. The level test must not provide a free-text answer field; the Skip action sends the backend `-1` answer that previously represented "I do not know".

React v2 must preserve the old app's functional meaning, not just its route names. The current v2 routes use the existing Go endpoints directly: lessons, practice, Listening, word learning/review, spelling, level test, vocabulary, mistakes, leaderboard, Premium, referrals, tools, progress, awards, and settings. Chat-like workflows (`lesson`, `practice`, `shadowing`, and `tools`) render through the shared chat display and show returned AI fields such as `feedback`, `reply`, `mistakes`, examples/context, spoken model text, pronunciation score, weak words, and tips; learner text/voice/image submissions are also appended as user-side chat messages. Voice-ready content is exposed as inline waveform listen controls backed by the Go/OpenRouter TTS endpoints, not browser speech synthesis. Trainer workflows (`words`, `word-game`, `spelling`, `level`, and `mistakes`) auto-start or auto-load where possible, use question/answer panels with selected/correct/wrong answer state, and keep results in the main context area. `shadowing` sends the target phrase back to `/api/shadowing/answer` and supports microphone recording through `MediaRecorder` plus audio file upload; vocabulary calls `/api/vocabulary?page=N` and renders compact paginated word cards when the backend provides pagination. Data workflows auto-load where possible: awards open centered trophy/story modals, leaderboard loads top-10 rows with a global-top button plus language dropdown, Premium plan buttons open a payment-method modal, and Referrals include copy/share/Telegram invite actions plus the current `usdt_rub_rate`. The Tools view calls `/api/tools/translator`, `/api/tools/voice-text`, and `/api/tools/image-translate`; the Mistakes view calls `/api/mistakes`, `/api/mistakes/practice/start`, `/api/mistakes/practice/answer`, and `/api/mistakes/clear`. Desktop QA should check the top ribbon and central context panel at 1440x900 and 1366x768 with 100% browser zoom. Mobile QA should check 390x844 and 430x932 with no horizontal overflow and usable touch targets.

The public landing, Privacy, and Terms are also rebuilt as a React/Vite multi-page site in `site-react`. It uses a shadcn-style `site-react/src/components/ui` folder, Tailwind CSS, and TypeScript. The current production landing has no public `/app/v2` shortcut in the header; the former V2 pill is replaced by a light/dark theme toggle. The first viewport combines the local `anomalous-matter-hero.tsx` Three.js shader with the existing Sparkles layer, animated CTA buttons, and concise copy for the real app features: Daily route, Roleplay, Pronunciation, Photo/translation, Offline/PWA, Notes/mistakes, Telegram sync, and progress. Build it with:

```powershell
npm run site:build
```

or through Nx:

```powershell
npx nx run english-coach-bot:site-build
```

The build writes `poliglot-ai.html`, `privacy.html`, `terms.html`, and hashed assets under `РЎР°Р№С‚ РїРѕР»РёРіР»РѕС‚Р° РґР»СЏ Р±РѕС‚Р°/assets/site-react`, which is the folder deployed to `/var/www/poliglotai`. The public React site and the React `/app` interface both expose the 35 supported languages listed above; the Go learning vocabulary set is dictionary-backed for the same 35 languages.

Authentication uses the current React shell at `/login` and `/app/login`. Caddy must route `/login` and `/login/*` to the Go web app together with `/app/*`; otherwise the public site fallback will serve the landing page instead of auth. Unauthenticated `/app` visitors are redirected to `/login`, where the user can switch the interface language and theme before registering or signing in. Authenticated visits to `/login` return to `/app`; Telegram profile completion can stay on the login shell while the user chooses whether to finish a new profile or link an existing web account.

Public website links to the web app should use relative `/app/` URLs. This keeps `poliglotai.ru` users on the Russian domain and `poliglotai.online` users on the European mirror instead of forcing either group to a single host. The landing page also rewrites visible app URL labels from the current host, so Russian-domain visitors see `poliglotai.ru/app` and online-mirror visitors see `poliglotai.online/app`.

The public Privacy and Terms pages keep the legal body text inside the React legal shell. Their first block uses the same restrained Sparkles treatment as the landing, and `legacyLegalContent.ts` is generated from the prior static pages so legal information is not lost during design work. The public Privacy page also keeps `РЎР°Р№С‚ РїРѕР»РёРіР»РѕС‚Р° РґР»СЏ Р±РѕС‚Р°/assets/privacy-policy-i18n.js` available for existing translation behavior. That file and the shared `legal_card_*`, `legal_updated_date`, and `legal_site_scope` keys in `assets/site-i18n.js` must stay translated for every supported site language. Public-site regression coverage lives in `site-react/e2e/public-site.spec.ts` and must keep checking the theme toggle, `/app/` app links, shader canvas on desktop and mobile, key feature headings, light-theme section-label readability, dark-theme language select readability, and the `@AsaselD` contact card.

When the React web app changes, upload `web/index.html`, `web/assets`, `web/manifest.webmanifest`, and `web/offline-deck-sw.js`. Do not upload or recreate `web/v2`; versioned paths are redirect-only. Caddy should continue reverse-proxying `/app`, `/app/*`, `/login`, and `/login/*` to Go; `deploy/caddy/Caddyfile.current` and `deploy/caddy/Caddyfile.updated` include those matchers.

## Learning Modes

Practice accepts text, voice, and image context. In the web app, the Practice composer exposes microphone/file input plus photo upload and camera capture; in Telegram, photos sent while Practice mode is active are routed into the same coach flow. The image is described through the OpenRouter vision path first, then the resulting compact context is added to the practice prompt as coach-only information so the model can discuss what is in the photo without grading the image text as the learner's writing.

The core study loop is shared between Telegram and `/app`: lessons create one active-recall micro-task, practice keeps a short adaptive conversation history, mistakes are saved for repair drills, and words move into vocabulary after enough correct review/spelling answers.

The CEFR level assessment stores the active question index, score, and per-attempt question order in `user.Mode`. Web `/api/level-test/start` resumes that active mode after a refresh instead of restarting from question one. New attempts create a shuffled order, and both Telegram callbacks and web answers keep using the stored order until the test is completed or the user chooses a manual level. Level-test prompts must keep their actual task logic localized into the interface language; do not flatten them into a generic `Choose the correct answer` label. Target-language examples/options stay in the learning language.

Daily Telegram reminders use the saved `interface_language` from the user profile. `system_i18n_clean_reminders.go` is the final clean 20-language override layer for reminder settings, reminder buttons, reminder footer text, and the rotating daily reminder messages. `reminderTextForDate` also uses a clean localized menu label instead of older UI button text, so CJK and other multilingual reminders do not inherit legacy mojibake from old copy maps. Any future reminder copy change must update that clean map and keep `TestReminderTextIsCleanAndLocalizedForEveryInterfaceLanguage` passing.

Listening / `РђСѓРґРёСЂРѕРІР°РЅРёРµ` is the dedicated listen-and-repeat pronunciation loop. Telegram keeps `/shadowing` and `/repeat` as compatible commands, while menus show `РђСѓРґРёСЂРѕРІР°РЅРёРµ`; the web app keeps the `shadowing` route with `POST /api/shadowing/start` and `POST /api/shadowing/answer`. It generates one target phrase in the learning language, plays it with OpenRouter TTS, accepts a voice repeat through `OPENROUTER_STT_MODEL=openai/gpt-4o-transcribe`, scores the transcript against the target phrase, and returns pronunciation score, accent strength, fluency, weak words, and tips. The STT prompt must not include the exact target phrase, so recognition is not biased toward the answer. If OpenRouter does not return word confidence/logprobs, the score is capped and marked as estimated instead of being allowed to look perfect from transcript-only evidence. It requires Premium because it uses voice recognition.

Voice answers in lesson and practice use the same OpenRouter-only pronunciation stack. The server normalizes audio to 16 kHz mono when `ffmpeg` is available, transcribes through `openai/gpt-4o-transcribe` via OpenRouter, then sends only compact transcript/confidence/fluency data to `OPENROUTER_PRONUNCIATION_MODEL=google/gemini-3.1-flash-lite`. The coach model never receives audio; it turns technical scoring into learner-friendly feedback. Lesson/practice prompts now receive localized section labels for all supported interface languages, and the web client normalizes older English labels before displaying structured blocks.

Lesson and practice audio is split into purposeful target-language clips when available. Lesson generation must include a localized `Model phrase` / example block that contains only one clean learning-language sentence, and web lesson `question_audio_text` reads that model phrase instead of the interface-language task instruction. Practice answer audio reads the full localized `Model phrase` from the coach reply, and the question clip reads the final target-language question. The `Model phrase` audio/title must contain only the target-language sentence: interface-language translations in trailing parentheses are stripped before TTS. The fallback is still the saved correction text, but the web and Telegram paths now prefer the complete phrase so learners hear the full target-language context.

Web app copy is layered as API session copy first, then clean static browser copy, then fallback strings. The browser includes a segment-based mojibake repair pass for old cp1251/UTF-8 mixed strings, but new UI text should still be added as normal UTF-8 in React source and Go source. The Tools translator surface has explicit localized labels for `Text В· voice В· photo`, file/camera/audio controls, and translator prompts; Telegram translator menu labels are localized through `toolUICopyFor`. Telegram translator language buttons are one-per-row and use `languageButtonLabel`, which shows the native language name plus the interface-language name in parentheses when they differ. Translator/pronunciation audio messages are tracked as recent pronunciation messages so the next generated audio can delete the previous clip before Telegram queues old audio.

The lesson and practice prompts intentionally stay on the configured OpenRouter model. Lesson generation now asks for a real-life micro-scenario, one teacher notice, useful chunks, and one production task. Practice feedback corrects the highest-impact issue first and recycles recent weak phrases when useful, instead of turning every exchange into a long grammar lecture.

## Account Linking

Website and Telegram accounts can share one learning profile. If a logged-in web user links Telegram and only one side has meaningful progress, the profiles merge automatically. If both sides already contain progress, the web app returns a signed short-lived merge challenge and asks which side should be primary: website or Telegram. The selected primary profile keeps its level, learning language, interface settings, and account identity, while XP, lesson/practice counts, vocabulary, mistakes, referrals, and other compatible progress are merged from the other side.

If a Telegram-authenticated user is asked to complete a web profile and enters the login/password of an existing web account, the server verifies that password and then uses the same merge-choice flow. This avoids a dead end where the UI only says "login is taken" even though the person owns both accounts.

Web Settings show the linked Telegram identity as a read-only status on both desktop and mobile: the session DTO includes `telegram_account` with the saved display name and numeric Telegram ID, and the Settings UI displays it without any unlink or replacement controls.

## Referrals

Referral registration immediately gives the invited user 7 days of Premium and records `invited_by` plus the inviter's `referral_count`. The inviter's 7-day Premium reward is delayed until the invited user reaches XP level 3 (`referralRewardXPLevel`). The `referral_level_rewarded` flag on the invited profile makes that level reward one-time only, while referral purchase balances still credit direct and indirect inviters separately.

## 2026 Redesign System

The current visual direction is documented in `docs/product/REDESIGN_PLAN.md` and in the Figma file `Poliglot AI Redesign System 2026`. The product should read as a premium AI language coach with a `Language Intelligence Cockpit` style: deep ink, indigo, teal, mint, selective gold, compact controls, calm panels, language-route lines, speech bubbles, waveform ribbons, progress rings, memory tiles, and trophy/premium signals.

The updated public-site direction also has a dedicated Figma landing file: `Poliglot AI Landing 2026` (`https://www.figma.com/design/G93sWSxT6kkcHqszrzqVQf`). It contains editable desktop and mobile frames for the refreshed site: hero, product modes, pronunciation scoring, 35-language proof, reviews, tariffs, Privacy, and Terms. Public HTML should follow the current Busuu-inspired product flow when landing content changes.

The public landing, Privacy, and Terms should describe voice learning from a learner-benefit perspective: lesson, practice, and Listening can accept voice input, compare the answer with the study phrase, and return an educational pronunciation hint. These public pages must not expose provider/model names, JSON examples, confidence thresholds, or the internal implementation recipe; keep those details in engineering docs and code comments only.

The public landing hero is intentionally compact on both desktop and mobile. It should expose the main promise, two CTA buttons, and a visible hint of the following feature section in the first viewport. Mobile `/app` auth uses a full-width top visual and an inner-padded form to avoid horizontal overflow in narrow browsers and Telegram WebViews.

Production web app images are stored in `web/assets`:

- app plan cards: `plan-free-light.png`, `plan-free-dark.png`, `plan-premium-light.png`, `plan-premium-dark.png`, `plan-platinum-light.png`, and `plan-platinum-dark.png`;
- app-wide backgrounds: `app-background-light.png` and `app-background-dark.png`;
- app section headers: `header-*-light.png` and `header-*-dark.png`;
- app navigation/action icons: `icon-*-light.png` and `icon-*-dark.png`;
- desktop right-inspector backgrounds: `panel-progress-light.png`, `panel-progress-dark.png`, `panel-limits-light.png`, `panel-limits-dark.png`, `panel-account-light.png`, and `panel-account-dark.png`;
- trophy artwork: `award-01.png` through `award-20.png`;
- compatibility aliases: `plan-*.png`, `header-*.png`, and `icon-*.png` are regenerated as light-theme fallbacks;
- public legal visual: `site-legal-shield.png` stays as the Terms/Privacy hero.

The public landing intentionally uses the user-uploaded static-site images in `РЎР°Р№С‚ РїРѕР»РёРіР»РѕС‚Р° РґР»СЏ Р±РѕС‚Р°/assets`: `poliglot-ai-avatar.jpg` for the hero/first block, `times.jpg` for languages/features, `primi.jpg` for use cases/pricing, `write.jpg` for start/how-it-works, and `yspex.jpg` for CTA/success. These are not overwritten by the web app asset generator.

The final PNG pack is deterministic and generated by:

```powershell
C:\Users\Admin\Documents\ComfyUI\.venv\Scripts\python.exe tools\create_poliglot_premium_assets.py
```

ComfyUI remains available for exploratory high-end bitmap generation through `tools/generate_comfy_brand_assets.mjs`. The default checkpoint is `sd3.5_large_fp8_scaled.safetensors` with the SD3.5 node path (`ModelSamplingSD3`, `CLIPTextEncodeSD3`, `EmptySD3LatentImage`, and optional `SkipLayerGuidanceSD3`). The generator also supports FLUX.2 Dev test runs with the ComfyUI split workflow: `UNETLoader` loads `flux2_dev_fp8mixed.safetensors`, `CLIPLoader(type=flux2)` loads `mistral_3_small_flux2_bf16.safetensors`, `VAELoader` loads `flux2-vae.safetensors`, then `Flux2Scheduler`, `EmptyFlux2LatentImage`, and `SamplerCustomAdvanced` produce the image. Juggernaut XL, SDXL, and FLUX.1 Dev remain comparison/fallback checkpoints. FLUX.2 Dev output is allowed on the current test server for visual review, but commercial production use still needs license and access review.

The social ad video workflow lives in `comfyui_workflows` and is copied into the local ComfyUI user workflow folder:

- source workflow: `C:\Users\Admin\Documents\New project\comfyui_workflows\ad_shorts_reels_tiktok_native.json`;
- installed workflow copy: `C:\Users\Admin\Documents\ComfyUI\user\default\workflows\ad_shorts_reels_tiktok_native.json`;
- model installer: `C:\Users\Admin\Documents\New project\comfyui_workflows\install_ad_short_models.ps1`;
- model manifest: `C:\Users\Admin\Documents\New project\comfyui_workflows\model_manifest.json`;
- placeholder inputs: `C:\Users\Admin\Documents\ComfyUI\input\ad_product_reference.png`, `ad_structure_reference.png`, and `ad_voiceover.wav`.

The default chain is intentionally disk-conscious: product reference -> FLUX.2 Dev keyframe -> CLIP Vision image lock -> Wan 2.1 I2V 480p 14B FP8 scaled -> VAE Decode -> native RIFE 4.26 frame interpolation -> Create Video with `ad_voiceover.wav` -> Save MP4. The optional branch uses the existing `RealESRGAN_x4plus.pth` model to upscale approved clips to 1080x1920. Generate motion at 480x832 first, then enable the upscale save node only for final candidates.

The first ad prompt is `C:\Users\Admin\Documents\New project\comfyui_workflows\poliglot_elite_20s_ad_prompt.md`. It defines a 20-second English XTTS voiceover, a premium young-adult business presenter direction, web-app/site visuals based on `web/assets`, and four 5-second visual beats. Prefer generating four short Wan clips and stitching them over one 20-second Wan pass; this keeps VRAM lower and usually gives cleaner product shots.

The first visual QA pass uses only two 5-second tests before expanding to the full 12-clip pack: `01_web_hero_presenter` shows the desktop web interface on a laptop with the presenter, and `02_mobile_voice_pronunciation` shows the mobile web voice/pronunciation interface on a smartphone. Their prompts are documented in `C:\Users\Admin\Documents\New project\comfyui_workflows\poliglot_2x5s_test_ad_variants.md`; generated variant workflows are under `C:\Users\Admin\Documents\New project\comfyui_workflows\test_variants`.

Install only the missing compact video files:

```powershell
PowerShell -ExecutionPolicy Bypass -File .\comfyui_workflows\install_ad_short_models.ps1
```

The workflow keeps XTTS and Whisper as input/output lanes instead of mandatory custom nodes so it loads in the current native ComfyUI install. Put XTTS output into `C:\Users\Admin\Documents\ComfyUI\input\ad_voiceover.wav` before running. After the MP4 is saved, run Whisper/subtitle burn-in as a post step or install dedicated ComfyUI subtitle nodes later. The parked ControlNet node is a sidecar for SD3.5/SDXL structure work; the main FLUX.2 + Wan path uses native reference latent and CLIP Vision conditioning for UI/product lock.

Run a Flux2 asset pass after the model pack is installed:

```powershell
$env:COMFY_BASE='http://127.0.0.1:8000'
$env:COMFY_CHECKPOINT='flux2_dev_fp8mixed.safetensors'
node tools\generate_comfy_brand_assets.mjs plans
node tools\generate_comfy_brand_assets.mjs headers
node tools\generate_comfy_brand_assets.mjs icons
node tools\generate_comfy_brand_assets.mjs backgrounds
node tools\generate_comfy_brand_assets.mjs panels
node tools\generate_comfy_brand_assets.mjs awards
```

To regenerate only the Listening / `РђСѓРґРёСЂРѕРІР°РЅРёРµ` assets in the same style:

```powershell
$env:COMFY_BASE='http://127.0.0.1:8000'
$env:COMFY_CHECKPOINT='flux2_dev_fp8mixed.safetensors'
$env:ONLY_ASSETS='header-shadowing-light,header-shadowing-dark,icon-shadowing-light,icon-shadowing-dark'
node tools\generate_comfy_brand_assets.mjs all
```

Transparent menu icons and award trophies use a stricter generation path than wide backgrounds. `tools/generate_comfy_brand_assets.mjs` prompts these files as object-only product renders on a flat `#00ff00` chroma-key background, then `tools/postprocess_transparent_asset.py` removes the key and normalizes the result:

- `icon-*-light.png` and `icon-*-dark.png` are fixed 512x512 transparent PNGs, with the visible object scaled to the same fill target and rejected if it becomes a tiny center mark, long horizontal strip, square tile, backplate, or opaque card.
- `award-01.png` through `award-20.png` are fixed 768x768 transparent PNGs, with broad trophy silhouettes, matching camera distance, and a tier prompt that progresses from steel/copper/bronze/silver/gold through rare crystal legendary forms.
- Each transparent asset is generated, post-processed, and validated before the generator advances to the next file. Failed attempts are retried with a deterministic seed offset and recorded in `tmp/asset-validation/*.json`.
- Section headers are prompted as one cohesive wide banner with no tiled pattern or repeated duplicate objects. CSS uses `background-repeat: no-repeat` and mask-friendly placement so the image fades inside the block instead of repeating across it.

If Hugging Face access to `black-forest-labs/FLUX.2-dev` has been accepted and `HF_TOKEN` or `HUGGINGFACE_TOKEN` is set, `tmp/download_flux2_dev_fp16_gated.ps1` can install the gated full-precision `flux2-dev.safetensors` into `ComfyUI\models\diffusion_models`.

Final production assets must be checked for fake text, pseudo letters, flags, people, screens, documents, and watermarks before deployment. The procedural generator is intentionally kept as the production-safe fallback because it guarantees no broken AI lettering and writes `web/assets/brand-assets-manifest.json` with theme, model, seed, prompt, dimensions, and QA status metadata.

## Public Website I18n

The public landing uses `РЎР°Р№С‚ РїРѕР»РёРіР»РѕС‚Р° РґР»СЏ Р±РѕС‚Р°/assets/site-i18n.js` for both `data-i18n` keys and phrase-level fallback translation. New landing text that is written directly in HTML or generated from JavaScript arrays must be added to `landingEnglishPhrases` or converted to a `data-i18n` key. This prevents non-Russian language views from showing Russian in marketing sections, pricing cards, comparison rows, FAQ answers, footer links, and translated attributes.

The phrase translator also decodes mojibake-looking source text before lookup, so older misdecoded landing strings can still match the clean Russian source phrase and render a non-Russian fallback.

`tools/rebuild_public_site_translations.mjs` rebuilds the generated phrase block for `poliglot-ai.html`, `terms.html`, and `privacy.html`. It writes the `public-site-translations-generated` block in `РЎР°Р№С‚ РїРѕР»РёРіР»РѕС‚Р° РґР»СЏ Р±РѕС‚Р°/assets/site-phrases.js`, deep-merging generated values with hand-written phrase translations instead of replacing the whole entry.

After changing public landing, Terms, or Privacy copy, run `node tools\rebuild_public_site_translations.mjs`, then syntax-check `site-i18n.js`, `site-phrases.js`, and the rebuild script. Use `publicSiteFinalOverrides` in `site-i18n.js` only for intentional manual corrections, and keep those overrides partial when only one language needs to change.

Common local command:

```powershell
go run .
```

Recommended low-memory settings:

```env
VOCABULARY_DIR=data/vocabulary
VOCABULARY_DATABASE_PATH=vocabulary.sqlite
ACTIVATION_KEYS_DATABASE_PATH=activation_keys.sqlite
ACTIVATION_KEYS_FILE=activation_keys.txt
MEMORY_LIMIT_MB=512
```

Keep `vocabulary_words*.json` next to the binary or point `VOCABULARY_DIR` to their directory.

## Serial Premium / Platinum Keys

Serial keys use their own SQLite database (`ACTIVATION_KEYS_DATABASE_PATH`, default `activation_keys.sqlite`) and their own import file (`ACTIVATION_KEYS_FILE`, default `activation_keys.txt`). This keeps manually issued Premium and Platinum keys separate from the main user database.

Import file format:

```text
XXXX-XXXX-XXXX-XXXX 30 premium optional-note
XXXX-XXXX-XXXX-XXXX 365 premium optional-note
XXXX-XXXX-XXXX-XXXX 30 platinum optional-note
XXXX-XXXX-XXXX-XXXX 365 platinum optional-note
```

`30` grants one month; `365` grants one year. `month`, `year`, `premium_30d`, `premium_365d`, `platinum_30d`, and `platinum_365d` are accepted aliases. The third field is the tier: `premium` or `platinum`. The key format is always four groups of four symbols; generated keys use the alphabet `ABCDEFGHJKLMNPQRSTUVWXYZ23456789` to avoid visually ambiguous characters.

Generate keys:

```bash
node tools/generate_activation_keys.mjs --count 100 --period month --out activation_keys.txt
node tools/generate_activation_keys.mjs --count 25 --period year --out activation_keys.txt
node tools/generate_activation_keys.mjs --count 10 --period month --tier platinum --out activation_keys.txt
```

Production Ubuntu command after deployment. It is installed as root-only:

```bash
poliglot-keys 100 month
poliglot-keys 25 year
poliglot-keys 10 month platinum
poliglot-keys 5 year platinum
```

The command creates `/opt/aibot/activation_keys_YYYYMMDD_HHMMSS.txt` and appends the same `XXXX-XXXX-XXXX-XXXX` keys to `/opt/aibot/activation_keys.txt`, which the bot imports automatically. If a learner enters a key that was already redeemed, the web app shows a separate "already activated" message instead of mixing it with "not found."

When a user redeems a key in web Settings or in the desktop account inspector activation form, the server imports the text file, validates the key, deletes it from active keys, writes a redemption record, and extends the selected paid tier. Both forms support Enter-to-activate and share the same `/api/premium/activation-key` path. The Settings form shows the localized result inline; the desktop right-inspector form also shows that same activated, already-used, or not-found result in a centered confirmation card so the response is visible away from the lower-right panel. A used key cannot become active again just because it remains in the text file.

## Encoding Safety

Visible UI copy must stay UTF-8. Do not bulk-decode or paste text from a mojibake-rendered terminal into React source, Go source, or public-site JavaScript. Add new labels through clean override maps or normal UTF-8 string literals, then run:

```bash
node tools/check_encoding_artifacts.mjs
```

The check fails on common broken Cyrillic, smart-quote, replacement-character, and broken emoji artifacts in the web app, API error strings, Telegram copy, and public-site shell files. The browser also rejects old mixed-encoding values via `cleanCopyValue`, but the source check must pass before deploy.

## Web Auth Protection

Cloudflare Turnstile can be enabled for web login, registration, and Telegram profile completion:

```env
WEB_TURNSTILE_SITE_KEY=your_site_key
WEB_TURNSTILE_SECRET_KEY=your_secret_key
```

When both values are set, the web app renders the Turnstile challenge and sends its token with auth requests. The server verifies the token before creating sessions or accounts. The API also has a global per-IP/per-session limiter for `/api/*`, plus stricter auth-specific limits, so logged-in users cannot spam buttons into unlimited server work.

`X-Forwarded-For` and `X-Real-IP` are trusted only when the TCP peer is loopback or matches `TRUSTED_PROXY_CIDRS`:

```env
TRUSTED_PROXY_CIDRS=127.0.0.1/32,::1/128,10.0.0.0/8
```

Payment creation endpoints accept `Idempotency-Key`. Repeating the same payment-create request with the same user, product, channel, and key reuses the provider-side idempotence/order token; direct crypto returns the existing pending invoice instead of creating another row.

Application rate limits are not volumetric DDoS protection. Production must keep an upstream edge layer in front of Caddy/backend, such as provider firewall rules, Cloudflare/WAF rules, or a verified Caddy build with rate-limit support.

## Direct Crypto Payments

The project supports direct crypto payments without Cryptomus, NOWPayments, or another payment processor. The same Go HTTP server that serves `/app` also exposes authenticated web endpoints:

- `POST /api/premium/crypto/payment` - creates a unique direct TON/USDT invoice for a Premium or Platinum product.
- `POST /api/premium/crypto/check` - checks the blockchain and activates Premium or Platinum after the matching transfer is found.

The web Premium screen and the Telegram bot expose direct crypto payment buttons. Supported methods are native TON, USDT on TON, and USDT TRC20. TON invoices use the merchant TON wallet, an exact amount in nanotons, and a unique required comment like `POLIGLOT:ton_...`. USDT on TON uses the same TON owner address by default, the official USDT Jetton master, 6-decimal microUSDT units, and the same required comment. USDT TRC20 uses the configured TRON address and a unique exact microUSDT amount because standard TRC20 transfers do not carry a reliable payment comment. After payment, the server checks the relevant blockchain, stores the transaction hash, extends the selected paid plan, and credits referral rewards. Telegram invoice expiry is displayed in the user's selected UTC offset; the web app displays the invoice expiry in the browser locale.

If `CRYPTO_TON_YEAR_AMOUNT` or `CRYPTO_TON_PLATINUM_YEAR_AMOUNT` is empty, the yearly TON invoice amount is derived from the configured monthly TON amount using the matching RUB price ratio. This keeps yearly TON checkout available on the web and in the bot when the monthly TON price is configured but the yearly TON env value was missed.

In Telegram, the Premium menu is a two-step flow: first choose a tariff, then choose a payment method for that tariff. The selected tariff submenu shows Telegram Stars, YooKassa/SBP when enabled, and only the crypto methods configured for that exact product.

Telegram Stars invoices encode the selected product, Telegram user ID, and creation timestamp in the invoice payload. Pre-checkout approval and successful-payment activation both validate the payload user ID and exact Stars amount against the current Telegram payer, so a malformed or mismatched payload is rejected instead of activating a paid plan.

On the web Premium screen, a newly created TON/USDT invoice is auto-checked only every 15 minutes while its status is still pending. The "Check payment" button remains immediate for manual checks. Once the server-side API/webhook or a manual/auto check confirms successful payment, or the invoice expires, the browser stops further checks and shows the paid/expired state without repeatedly repainting the invoice block.

## Telegram Stars From Web

The web Premium screen opens a Telegram deep link such as `https://t.me/Poliglot_AI_bot?start=buy_premium_30d` for Stars payment in a new tab only; the current `/app` page must stay open. When the user opens the bot, the `/start buy_...` payload sends the matching `currency=XTR` invoice. `POST /api/premium/stars` remains available for linked Telegram accounts and also returns the same bot URL after sending an invoice to the bot chat.

Stars payment completion still follows the normal Telegram pre-checkout and successful-payment flow, so Premium/Platinum activation, duplicate-payment protection, and referral rewards stay shared with bot-originated Stars invoices.

Required server configuration:

```env
CRYPTO_TON_WALLET=EQ...
CRYPTO_TON_MONTH_AMOUNT=1.25
CRYPTO_TON_YEAR_AMOUNT=10
CRYPTO_TON_PLATINUM_MONTH_AMOUNT=3
CRYPTO_TON_PLATINUM_YEAR_AMOUNT=30
CRYPTO_TONCENTER_API_KEY=
CRYPTO_TONCENTER_BASE_URL=https://toncenter.com/api/v2
CRYPTO_TONAPI_KEY=
CRYPTO_TONAPI_BASE_URL=https://tonapi.io/v2
CRYPTO_TONAPI_WEBHOOK_KEY=change_me_to_a_long_random_secret_for_tonapi_webhook
CRYPTO_USDT_MONTH_AMOUNT=
CRYPTO_USDT_YEAR_AMOUNT=
CRYPTO_USDT_PLATINUM_MONTH_AMOUNT=
CRYPTO_USDT_PLATINUM_YEAR_AMOUNT=
CRYPTO_USDT_RUB_RATE=72
CRYPTO_USDT_TON_WALLET=
CRYPTO_USDT_TON_JETTON_MASTER=EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs
CRYPTO_USDT_TRC20_WALLET=TFoRoWRxgYWzmnJAwsom51DmgnNyYD368M
CRYPTO_USDT_TRC20_CONTRACT=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
CRYPTO_TRONGRID_API_KEY=
CRYPTO_TRONGRID_BASE_URL=https://api.trongrid.io
CRYPTO_PAYMENT_TTL_MINUTES=60
```

`CRYPTO_TONAPI_KEY` is the API key from TonConsole and is the recommended production setting for native TON checks and current USDT/RUB rates. If it is present, native TON payment checking asks TonAPI first; TON Center remains a fallback and is retried without a rejected TON Center key. USDT-TON checks use TON Center API v3 Jetton transfers. USDT-TRC20 checks use TronGrid's TRC20 transaction history endpoint; `CRYPTO_TRONGRID_API_KEY` is optional but recommended for production rate limits.

If explicit `CRYPTO_USDT_*_AMOUNT` values are empty, USDT invoices are calculated from the RUB plan price using TonAPI `/v2/rates` for the USDT TON jetton against RUB. The rate is cached briefly and `CRYPTO_USDT_RUB_RATE` stays as the fallback if TonAPI is unavailable. TRC20 invoices may add a small uniqueness offset.

If you create the API key in TonConsole, put it into `CRYPTO_TONAPI_KEY` (not `CRYPTO_TONCENTER_API_KEY`). `CRYPTO_TONAPI_WEBHOOK_KEY` must be only the secret token, not the full URL. You can then add a TonAPI webhook in TonConsole with this endpoint:

```text
https://poliglotai.ru/tonapi/webhook/<CRYPTO_TONAPI_WEBHOOK_KEY>
```

Then subscribe the webhook to the merchant TON wallet account. The webhook is used as a wake-up signal: the server receives the transaction hash/LT from TonAPI, rechecks pending TON and USDT-TON invoices against the blockchain, and activates the matching Premium or Platinum invoice automatically.

Example Caddy shape:

```caddyfile
poliglotai.ru, www.poliglotai.ru, poliglotai.online, www.poliglotai.online {
	reverse_proxy 127.0.0.1:8080
}
```

The saved production Caddy snapshot is in `deploy/caddy/Caddyfile.current`; the updated multi-domain candidate is in `deploy/caddy/Caddyfile.updated`.

The data model stores `payment_id`, `tx_hash`, `network`, `telegram_id`, `product`, `status`, and `expires_at`, so new direct crypto methods can be added without replacing the subscription model.

## Pricing and Plan Matrix

| Plan | Monthly | Yearly | Limits |
| --- | --- | --- | --- |
| Free | 0 RUB | 0 RUB | 5 lessons, 15 practice messages, no voice |
| Premium | 300 RUB / 150 Stars / about 4.15 USDT | 3000 RUB / 1500 Stars / about 41.50 USDT | 50 lessons, 200 practice messages, 20 voice messages up to 30 seconds |
| Platinum | 590 RUB / 300 Stars / about 8.15 USDT | 5900 RUB / 3000 Stars / about 81.50 USDT | 100 lessons, 500 practice messages, 60 voice messages up to 30 seconds |

The visible paid prices are launch prices with the 70% discount already applied. The public landing shows rounded crossed-out base prices for the first launch month: Premium about 1000 RUB/month and 10000 RUB/year, Platinum about 2000 RUB/month and 20000 RUB/year, with matching rounded Stars and USDT equivalents.

The web app Premium screen groups purchases into adaptive Free, Premium, and Platinum cards. The live cards use theme-aware PNG pairs from `web/assets/plan-free-light.png`, `web/assets/plan-free-dark.png`, `web/assets/plan-premium-light.png`, `web/assets/plan-premium-dark.png`, `web/assets/plan-platinum-light.png`, and `web/assets/plan-platinum-dark.png`; each card also keeps inline SVG artwork as a fallback layer if an image route fails. The production-safe PNG pack is generated by `tools/create_poliglot_premium_assets.py`, while `tools/generate_comfy_brand_assets.mjs` remains available for ComfyUI experiments.

## Web Auth

The login entry points `/login` and `/app/login` now serve the React web bundle in standalone auth mode. The page is implemented by `web-react/src/components/ui/sign-in.tsx`, under the current shadcn-compatible component path `web-react/src/components/ui` because the Vite alias `@` resolves to `web-react/src`.

The React login screen keeps the existing backend contracts:

- `/api/auth/login` for login/password sign-in.
- `/api/auth/register` for account creation with optional referral code and selected interface language.
- `/api/auth/telegram/start` and `/api/auth/telegram/status` for Telegram sign-in through the existing six-digit OTP dialog.
- `session.captcha.provider === "turnstile"` with `session.captcha.site_key` for Cloudflare Turnstile. The token is sent as `captcha_token` to login/register requests.

The auth page must not render the V2 app header, ribbon, bottom navigation, or any in-app menu controls. If the session is already authenticated, the standalone login route redirects to `/app`.

Standalone auth copy is part of the V2 35-language localization layer in `web-react/src/lib/i18n.ts`. Keep `auth_login_hint`, `auth_register_hint`, `auth_support`, `auth_subtitle`, `auth_referral_placeholder`, `auth_finish_account`, `auth_words`, Turnstile/Telegram labels, and password visibility labels localized before deployment; otherwise non-English login screens will fall back to English.

Main web app sections receive a view-specific branded header treatment through `stage-*` classes. `setChat` adds the current view class to `home-stage` sections automatically, and fallback `stage-panel` styling covers first-message panels such as Practice. Desktop and phone layouts use the same CSS variables, now backed by generated `web/assets/header-*-light.png` and `web/assets/header-*-dark.png` images for Home, Lesson, Practice, Progress, Awards, Words, Word Game, Spelling, Vocabulary, Level, Leaderboard, Premium, Limits, Mistakes, Tools, Referrals, and Settings. Action cards, desktop sidebar buttons, mobile bottom navigation, and mobile overflow actions use matching `web/assets/icon-*-light.png` and `web/assets/icon-*-dark.png` packs. The main menu uses larger generated artwork that fills the action-card height, while the global shell uses `app-background-light.png` and `app-background-dark.png` behind translucent UI. The desktop right inspector uses generated `panel-progress-*`, `panel-limits-*`, and `panel-account-*` backgrounds under real HTML text.

Production Caddy proxies `/app` and `/app/*` to the Go server while root `/assets/*` belongs to the static public site. For that reason the web app references generated visuals through `/app/assets/...`. The Go router also keeps `/assets/...` routes for local/dev compatibility. HTML is intentionally disk-only, while asset handlers first try `web/assets` from the working directory or beside the executable before falling back to embedded image files.

The Telegram main keyboard includes the web app button as a primary entry point. The translator/tools inline keyboard no longer duplicates the web app button.

## Vocabulary

Large vocabulary JSON files are source files, not the preferred production runtime format. Startup creates these read-only dictionary tables in `VOCABULARY_DATABASE_PATH` (`vocabulary.sqlite` by default):

- `vocabulary_words` - one normalized row per learning word.
- `vocabulary_translations` - interface-language prompts by word id.
- `vocabulary_sources` - JSON file size, modification time, word count, and import timestamp.

The import is automatic and idempotent. If a `vocabulary_words*.json` file changes size or modification time, that language is reimported. If SQLite vocabulary rows are not available, the old streaming JSON reader is still used as a fallback.

Hot lookups now use indexed SQLite reads for `id -> word`, learned-word page metadata, and review-game option selection. This avoids repeated JSON parsing under concurrent users while keeping memory stable.

The backend opens the vocabulary database and starts HTTP before long dictionary sync/index work. JSON reimport and optimization-index creation run in the background, so `/healthz`, login, payments, and static web routes stay available during deploys or first-start dictionary maintenance.

The vocabulary database is separate from the user `DATABASE_PATH` so daily user-data backups do not copy the static dictionary tables. It can be regenerated from the JSON files if removed.

The learn-word flow also uses SQLite for the next-word lookup and answer options, so tapping `Learn words` does not scan the full dictionary. Learned vocabulary pages show the newest mastered words first.

Learn-word and review distractors are sampled only inside the user's CEFR band: `A1/A2`, `B1/B2`, or `C1/C2`. Exact-level lookup runs first; if the exact level is exhausted, fallback stays inside the same band and does not jump forward into the full language pool. If a sparse band cannot provide three distractors, the UI may show fewer options rather than mixing unrelated levels. The fix is to enrich the dictionary band, not to pull random words from all levels.

Dictionary prompts keep one primary native-language value as the prompt and one target-language word per answer option. Additional native-language meanings from the dictionary are exposed through the context/hint line when available, so the learner sees semantic context without revealing the target-language answer.

## Web Mobile

On phone-sized screens, the web app shows a compact top header with the Poliglot AI brand, the current learning language, the CEFR level, a small XP-level progress chip, account chip, Premium remaining chip, and interface-language selector. The header shows the remaining Premium time in a short form instead of a full expiration date. `Tools` has its own bottom navigation tab on mobile and opens the compact in-app Tools view instead of a separate app page. The mobile Tools picker is stacked top-to-bottom as full-width cards so it matches the app header width and does not require horizontal scrolling. The mobile bottom navigation does not include a separate `Main menu` button.

Tools and referral cards use final CSS overrides in `web-react/src/styles/app.css` so late redesign rules do not reintroduce opaque panels or narrow text columns. The Tools voice-limit card is left-aligned on desktop; on mobile its plan/status pill sits on a second line below the counter. The mobile Referrals hero keeps the copy in a wider text column and lets the invite/balance card breathe below it.

Clean web-only localization overrides are kept in the React i18n layer and related V2 copy helpers, not in a legacy HTML file. These overrides intentionally sit after the broad interface copy cleanup layer and before tool-surface copy, so Referrals and empty Vocabulary states do not fall back to English or older mixed-encoding strings when Chinese, CJK, Armenian, Georgian, or other supported interface languages are selected.

The mobile leaderboard uses compact cards, tighter rows, right-aligned scores, and horizontally scrollable language tabs so the list fits comfortably on phones. The leaderboard stage is sized to the chat container, not raw viewport width, so it aligns with the mobile header and menu. Leaderboard language labels are shown in English names for every interface language.

Mobile spelling tracks wrong attempts for the active word. After three wrong answers, the web app offers `РќРµ Р·РЅР°СЋ`; choosing it reveals the correct answer without awarding XP or increasing spelling progress, then lets the learner move to another word while keeping the skipped word in the spelling pool. Mobile Learn Words scrolls to the `Next word` action after a correct answer.

Mistake pagination renders previous/next controls only when another page actually exists. Disabled edge controls are omitted instead of shown as inert buttons, and labels must come from localized copy so mojibake arrows do not appear in mobile or desktop shells.

## Awards

The web app has an `Awards` view tied to the existing `xp_level` scale from 1 to 20. Each unlocked level has generated trophy PNG artwork in `web/assets/award-01.png` through `web/assets/award-20.png`, a polyglot rank label, and a short world-language-history story. Clicking an unlocked trophy opens the story. Locked trophies do not open a modal; their artwork is gray and blurred until the level is reached.

Award rank names and stories are localized for every interface language currently supported by the bot: `ru`, `en`, `es`, `de`, `fr`, `it`, `zh`, `ja`, `ko`, `tg`, `uz`, `tt`, `hy`, `kk`, `ky`, `ka`, `uk`, `pl`, `ro`, and `pt`. The historical story is the same in meaning for every user; only the language changes according to `interface_language`.

The browser safe-display filter must allow clean Cyrillic, CJK, Armenian, Georgian, and other supported-script award copy. It should reject only actual mixed-encoding artifacts; otherwise Russian award stories can fall back to English even when localized data exists.

When a web action returns a user object with a higher `xp_level` than the previous session user, the app shows a short trophy celebration with animated fireworks and stores the seen level in `localStorage` to avoid replaying old awards on login.

## Gamification Direction

The gamification system should stay shared across Telegram bot, desktop web, and mobile web so every surface reflects the same learner identity:

- Core identity: XP level 1-20, CEFR level, learning language, polyglot rank, and trophy collection.
- Short-term motivation: daily streak, daily quests, weekly quests, and small completion rewards for lessons, practice, spelling, review, and tools.
- Mastery loops: word mastery badges, mistake-cleanup badges, pronunciation/voice badges, and translation-tool badges.
- Social loops: leaderboards by all users and by learning language, seasonal leaderboards, referral boosts, and public rank cards.
- Celebration loops: level-up animation in web/mobile, compact level-up message in Telegram, and trophy reveal wherever the user earned the level.
- Retention loops: reminder text should mention streak/quest progress when useful, but never block learning behind streak pressure.

The next implementation step is to move award metadata and localized stories into a shared Go/JSON source so Telegram messages, web, and future mobile app screens use the same ranks and stories without duplicating text.

## Daily Reminders

Daily Telegram reminders are selected by local due time but rendered in the user's selected interface language. `dueReminderUsers` reads `interface_language` with each reminder target, `sendDueReminders` passes it into `reminderTextForDate`, and the send call uses the same localized `ui(user)` copy for Telegram buttons. Languages with full `systemUICopy` entries use their native reminder copy; compact fallback languages build reminder text from their localized UI labels instead of falling back to English or Russian.

## Referrals

Referral data is exposed in three places:

- Web sidebar: `Referrals`.
- Mobile web menu: `More -> Referrals`.
- Telegram bot: main menu `Referrals` and `/referral`.

The referral block shows the personal code, share link, invitation count, referral balance, share/copy actions, and withdrawal status.

The web referral view must use localized labels from the clean static overrides for dashboard title, invitations, referral balance, share/copy/withdraw buttons, the 20%/5% balance hint, and withdrawal unavailable text. The visible `USDT/RUB` token can stay as a currency-pair label, but surrounding text should remain in the selected interface language.

The referral block is intentionally not duplicated in Settings. Referral sharing uses the selected interface language and addresses the invited person directly: it describes Poliglot AI, the key bot/web features, and the 7-day Premium bonus received from the link.

## Mobile Web

The mobile layout uses one horizontal bottom navigation rail. There is no separate `More` destination: former secondary destinations are part of the same scrollable rail. Long-press enters edit mode, the dragged item lifts, neighboring items open a visible gap, and the rail auto-scrolls near the left or right edge while dragging.

Lesson, practice, spelling, and mistake-practice screens share a compact bottom composer with a fixed-height text field and send button. Text inputs across auth, settings, Telegram-code entry, the composer, and translator textarea use one shared visual style. Text-entry screens do not focus the typing field automatically, so the user opens the keyboard by tapping the field.

## V2 Learning Interaction Update

Each V2 menu destination now has its own output context through `data-view` and view-specific surfaces. Roleplay opens a dedicated scenario dialogue window after choosing a tile; the scenario picker is no longer the active output while the roleplay session is running. Offline decks render full-width scrollable cards and can export readable TXT study packs.

Dashboard is focused on the learner's own activity, not leaderboard or weak-topic repetition. It shows XP, level, Premium/referrals, words learned, practiced phrases, completed lessons, review rounds, voice attempts, and current mistake workload.

Pronunciation exact-repeat scoring is intentionally strict. The backend caps scores for missing words, substitutions, low text similarity, and estimated confidence. The frontend translates raw technical issue codes into learner-facing advice before showing weak words/sounds.

Roleplay reuses the same chat audio path as practice. `question_audio_text` is generated from labeled `Your turn`/question lines first; if a roleplay reply has no final question mark, the backend falls back to the last short learner-facing dialogue line. The frontend keeps question audio clips even when the formatted text already includes a spoken-model line.

Listening/Audition keeps a single active work surface: the duplicate upper phrase card is filtered out, and the spoken-model audio clip is rendered inside the lower repeat panel next to the target phrase and next-phrase action.

Mistakes use a two-step flow on desktop and mobile. The first screen is the dictionary/list with category filters. Choosing one saved mistake opens a dedicated practice screen; a correct answer clears the selected practice and returns to the list, while the back action cancels the practice without changing the dictionary.

Spelling API responses include `correct_answer` for both correct answers and give-up answers. The web trainer renders that target-language spelling only after a correct answer or an explicit give-up action, so a normal wrong attempt can say "try again" without revealing the answer early.

## Vocabulary Context Cache

Learn Words and Review keep selection inside the active CEFR band: `A1/A2`, `B1/B2`, or `C1/C2`. When a dictionary row lacks a useful native-language context clue, the bot/web API can call the configured vocabulary model once, sanitize the answer so it does not reveal the target word, and persist the hint/example in SQLite table `vocabulary_ai_cache`. Later rounds reuse that cached value instead of calling the model again.

Vocabulary clue generation explicitly asks the model for native/interface-language meanings that do not equal or repeat the hidden target-language word. The sanitizer drops any generated clue that repeats the target word, including same-surface borrowed words, and falls back to dictionary data or a neutral localized hint instead.

## 35 Learning Dictionaries

All 35 interface languages are now available as local dictionary-backed learning languages. The 15 added learning dictionaries are Arabic, Bengali, Czech, Greek, Hindi, Hungarian, Indonesian, Dutch, Swedish, Tamil, Telugu, Thai, Tagalog, Turkish, and Vietnamese.

The new dictionaries are generated from Kaikki/Wiktextract JSONL dumps stored in `tools/sources/kaikki-<code>.jsonl`. Kaikki data is derived from Wiktionary and follows the same free-license terms as Wiktionary. `tools/build_multilingual_vocab.py` includes script-aware filtering for Arabic, Bengali, Greek, Hindi, Tamil, Telugu, and Thai so imported headwords stay in the expected writing system instead of accepting Latin fallback entries from the dump.

After adding or refreshing a dictionary source, run `python tools/build_multilingual_vocab.py`, then `go test ./...`. The Go regression suite verifies that every interface language has a learning dictionary, every generated dictionary loads, and new script dictionaries do not start with Latin headwords.

## Permanent Deploy Upload Service

The Ubuntu server has a persistent `poliglot-deploy-upload.service` for large deploy artifacts when PuTTY `pscp` or streamed SSH uploads fail.

- Public HTTPS route: `https://poliglotai.ru/__codex_deploy_upload/<filename>`.
- The service itself binds only to `127.0.0.1:19081`; Caddy proxies the route and enforces a `450MB` request body limit.
- Uploads require the `X-Deploy-Token` header. The token is stored on the server in `/opt/aibot/deploy-upload/deploy-upload.env` and should not be printed in user-facing answers.
- Allowed filenames are intentionally narrow: `aibot-linux-amd64`, `poliglot-app-web.tgz`, `poliglot-public-site.tgz`, `Caddyfile.deploy`, `deploy_react_server.sh`, and smoke-test file `deploy-upload-smoke.txt`.
- `poliglot-app-web.tgz` must include both `web/` and `data/vocabulary/vocabulary_words*.json`; `deploy_react_server.sh` copies those dictionaries into `/opt/aibot` before restarting the service.
- The permanent upload service also accepts optional chunk headers `X-Deploy-Chunk-Start` and `X-Deploy-Total-Bytes` for the same allowlisted filenames. This is a resilience feature for unstable large PUT transfers, not a separate temporary route.
- Reusable operator prompt and commands are in `deploy/SERVER_UPLOAD_PROMPT.md`.

## 2026-05-29 Regression Closure

Privacy pages are rendered by the React public site plus `assets/privacy-policy-i18n.js`. The generated privacy policy data for Russian and English must include three contacts: email, admin Telegram, and `Bot: @poliglot_ai_bot`. The runtime renderer also falls back to that bot handle if an older locale block misses the third item, so `privacy.html?lang=ru` cannot show an empty Telegram bot card.

The web login screen is intentionally a single centered form over the generative background. The old right-side desktop panel is not part of the DOM. The auth regression verifies the Cloudflare slot, login API payload, Telegram OTP path, absence of side panels, centered form geometry, and full-viewport background canvas on desktop and mobile.

The mobile bottom menu remains one horizontal rail. Long-press enters edit mode, pointer-capture failures are ignored for WebView/test compatibility, drag movement reorders the real rail items, edge movement auto-scrolls the rail, and the final order is stored in `poliglot-mobile-nav-rail-v2:<account>`. The Playwright regression verifies reorder, auto-scroll, and persistence after reload.

Vocabulary Learn Words and Review use the SQLite dictionary tables and indexes documented above. Prompts choose the primary native/interface-language dictionary value; alternate native meanings are exposed as context/hints and capped so the learner sees useful semantic context without a long synonym dump. If a precise prompt is missing, the vocabulary model must answer the strict dictionary JSON prompt; sanitized results are stored in `vocabulary_ai_translations` and mirrored into `vocabulary_translations`. Distractors stay inside the active CEFR band (`A1/A2`, `B1/B2`, or `C1/C2`).

## 35-Language V2 Localization Guard

The React V2 app uses `appCopy(interface_language, key, fallback)` as the clean client-side source for visible labels. Server-provided copy from session payloads is treated as lower priority for V2 UI text because older server dictionaries may still contain mojibake or generic technical fallbacks. Values that look like `Section: <key>`, `Mб»Ґc: <key>`, exact English fallback text, Russian fallback text in non-Russian locales, or known mojibake sequences are rejected before display and replaced with localized client fallbacks. This guard covers nested V2 surfaces including Notes/Phrasebook, Roleplay scenario cards, Dashboard metrics, Premium/payment requisites, learning-lab cards, and compact daily quests.

Roleplay scenario titles and descriptions use the same guard. If a localized scenario string is missing or corrupt, the app falls back through the scenario-specific `appCopy` key and then to a neutral localized roleplay label instead of showing mojibake or English helper text.

The Playwright regression `main V2 screens do not expose mojibake or fallback labels across interface languages` opens Today, Roleplay, Pronunciation, Offline, Mistakes, Leaderboard, Spelling, Tools, and Dashboard for every non-English interface language on both desktop and mobile. It fails on known mojibake artifacts and exposed English/Russian technical fallback labels. Because this test covers hundreds of page loads, it has an explicit extended timeout.

## Function Ribbon Reordering

The desktop function ribbon stores the user's reordered menu IDs in localStorage under `poliglot-function-ribbon-v2:<account>`. Dragging a chip sets visible before/after drop targets, auto-scrolls the ribbon near horizontal edges, and writes the final order after drop. Invalid or stale stored IDs are ignored and the current navigation list fills any missing destinations.
