# Changelog

## 2026-06-29

- Refined the dynamic outro into v12: masked the unwanted model-rendered text stripe in the center and replaced the stiff URL reveal with a comet-style letter/particle assembly. Final MP4: `C:\Users\Admin\Documents\ComfyUI\output\outro\poliglot_outro_tech_warp_v12_comet_text.mp4`.
- Added a new dynamic Poliglot AI 2-second tech-warp outro v11 after rejecting the flatter/static space direction: holographic tunnel depth, logo shard assembly, light-object URL arrival, and deterministic readable `poliglotAI.online` lockup.
- Generated the checked 1080x1920 render pass at `C:\Users\Admin\Documents\ComfyUI\output\outro\poliglot_outro_tech_warp_v11.mp4`, with preview frames/contact sheet under `tmp/comfy-outro`.
- Queued the matching ComfyUI v11 Wan/RIFE API workflow and used its output as the moving source for the final composite at `C:\Users\Admin\Documents\ComfyUI\output\outro\poliglot_outro_tech_warp_v11_comfy_composite.mp4`.
- Added the v11 reference image, API prompt, prompt notes, and compositor script in `comfyui_workflows` so the same direction can be regenerated through the running ComfyUI API without distorting the final URL.

## 2026-06-28

- Added a polished After Effects-style Poliglot AI outro render pass that reuses the more kinetic first ComfyUI motion background, masks distorted model text, keeps the generated logo animation, and overlays only the exact `poliglotAI.online` URL with a floating reveal.
- Reworked the URL typography to use bundled Manrope ExtraBold with solid white text and a mint `AI` accent for better Shorts/Reels/TikTok readability after compression.
- Cleaned the final composition so the post-render no longer overlays a second logo on top of the generated logo animation and no longer draws decorative stripes or underline below the URL.
- Rendered the 2-second vertical MP4 at `C:\Users\Admin\Documents\ComfyUI\output\outro\poliglot_outro_ae_reveal_v9_clean.mp4` and saved preview frames under `tmp/comfy-outro`.

## 2026-06-27

- Added a dedicated Poliglot AI 2-second TikTok/Reels/Shorts outro pack for ComfyUI: FLUX.2 keyframe prompt, Wan 2.1 I2V prompt, native ComfyUI workflow, API prompt, exact `poliglotAI.online` reference card, silent audio placeholder, and an API runner script.
- Copied the outro workflow and input reference files into the local ComfyUI workspace so the insert can be opened in ComfyUI Desktop or queued through the API without downloading new models.

## 2026-05-30

- Fixed auth/register system errors so backend responses include stable auth error codes, the React login page translates those errors through the V2 localization layer, and invalid login messages no longer fall back to generic `Р Р°Р·РґРµР»`/`Section` labels. Added regression coverage proving `friend.name` is valid while `.friend` returns a clear coded error.
- Deployed the auth-error localization fix through the permanent deploy-upload service on `.online` without creating backups, restarted `aibot.service`, and verified `.online`/`.ru` `/healthz`, `/app/v2/`, upload health, installed binary size, and fresh V2 asset `index-DAu6soW5.js`.
- Updated the PWA refresh path so installed mobile app shells fetch `/app/v2` from the network first, bump the service-worker cache to `poliglot-v2-offline-decks-20260530`, skip waiting on new workers, and serve HTML/manifest/service-worker entrypoints with no-store headers.
- Deployed the PWA refresh fix through the permanent deploy-upload service without creating server backups, restarted `aibot.service`, and verified public `.online`/`.ru` `/app/v2/` no-store headers plus the new service-worker cache marker.
- Updated the Learn Words success result so it shows the learned target-language word alongside the interface-language translation, then offers the same Notes quick-save chip under `Next`.
- Added Playwright regression coverage for the Learn Words success result, saving that word pair to Notes, wrong Spelling attempts not revealing the answer, full 35-language level-test localization, and clean non-mojibake V2 surfaces.
- Replaced the vocabulary-guessing level-assessment fallback for the 15 newly added learning languages with static 36-question exam sets, so all 35 learning languages now use a real level test instead of bare word guessing.
- Localized level-assessment task labels for all 35 interface languages, so static English-style questions keep their real intent (`fill gap`, `translation`, `meaning`, `natural option`, `correct form`, etc.) instead of falling back to a generic choose-answer prompt outside Russian.
- Added regression coverage that verifies CEFR level determination thresholds (A1 through C2) for the 15 newly added learning languages, with prompts localized through the user's interface language and answer options generated in the selected learning language.
- Fixed level-assessment prompt localization so the visible task/question is always in the interface language, while the example/options stay in the learning language; the new-language fallback now keeps a real question form instead of flattening prompts to a generic `Choose answer`, and uses target-language vocabulary options without slow full-dictionary scans.
- Restored the original level-test task wording logic for static questions: translation, gap-fill, meaning, grammar, and naturalness prompts now stay distinct instead of collapsing into the generic `Choose the correct answer` label.
- Redeployed the restored level-test question wording to production after verifying `go test ./...`, `npm --prefix web-react run build`, public `/healthz`, public `/app/v2/`, and server-side `aibot.service`.
- Fixed lesson example audio so the lesson generator emits a dedicated target-language model phrase and web lesson TTS reads that phrase instead of the interface-language task instruction.
- Fixed mobile Spelling results so a wrong attempt no longer reveals the correct spelling before the learner gives up or answers correctly, added bottom padding for the mobile result panel, and replaced the generic `Р Р°Р·РґРµР»` audio label with the localized example label.
- Deployed the level-assessment, lesson-audio, and mobile Spelling fix to `/opt/aibot` through the permanent deploy-upload service, restarted `aibot.service`, verified public `/healthz` plus `/app/v2/`, and kept `go test ./...`, `npm --prefix web-react run build`, `npm --prefix web-react run e2e`, and `node tools/check_encoding_artifacts.mjs` green.
- Fixed the V2 Dashboard/Progress statistic labels so metric cards no longer collapse to the generic `Statistics` label, localized the affected metric names for all 35 interface languages, corrected Premium/Platinum plan titles and tiers, and changed the Phrasebook second field placeholder to `Note or translation` with 35-language regression coverage.
- Persisted desktop function-ribbon and mobile bottom-menu ordering per backend user through `/api/navigation-layout`, with browser storage kept only as a local cache.
- Saved wrong Spelling answers into the Mistakes dictionary, alongside the existing Lesson/Practice mistake capture and fallback extractor.
- Changed web bug-report Telegram notifications to send all attached screenshots as actual Telegram photos; the first photo carries the full report caption and later photos carry the same report ID with screenshot numbering.
- Fixed Offline deck exports so TXT/JSON downloads follow the currently selected group (`All`, `Words`, `Notes`, or `Mistakes`), added a UTF-8 BOM for readable TXT files on Android viewers, and added per-card delete buttons inside offline cards.
- Added regression coverage proving every supported learning language can produce level-assessment questions, including vocabulary-backed fallback tests for languages without a hand-written static grammar set.
- Fixed the web Mistakes API for client-side pagination: it now returns the full mistake list with newest items first, while keeping a first-page slice for API consumers.
- Persisted the web habit/visit calendar in the backend store for each user with a rolling 365-day window, and returned it through `/api/session` so it no longer depends only on browser storage.
- Slowed mobile bottom-menu edge auto-scroll while dragging reordered tiles so the rail moves in controlled small steps instead of jumping sideways.
- Removed raw lesson feedback text for `Spoken model`/`РћР·РІСѓС‡РµРЅРЅС‹Р№ РѕР±СЂР°Р·РµС†`; spoken samples now stay available through the audio button/control.
- Added a backend fallback mistake extractor for web lessons and practice so learner errors still reach the Mistakes dictionary when the main model response omits or breaks the `---MISTAKES---` JSON block.
- Added Mistakes pagination at 10 cards per page, matching the Offline and Phrasebook paging pattern.
- Made the pre-hydration loading screen background visibly animated with drifting shader-like gradients instead of a static backdrop.
- Added mobile Phrasebook pagination at 10 saved notes per page, matching the Offline deck paging controls.
- Fixed the level-test Skip button's inner text span so `Skip` stays on one line instead of wrapping the final letter below the button label.
- Made mobile leaderboard and roleplay panels grow to the end of their lists so card content no longer continues below a prematurely closed block border.
- Added OpenRouter-backed phrasebook auto-translation for saved phrases when no manual note/translation is supplied, and covered it with a backend regression test.
- Fixed the mobile V2 level-test Skip button so it no longer falls back to `Р Р°Р·РґРµР»`, and tightened mobile vocabulary/level card wrapping so long text stays inside the block above the bottom rail.
- Restored the branded Poliglot AI app boot screen before React hydration: the HTML shell now shows the logo, Poliglot AI title, spinner, and animated dark glass background immediately instead of a blank screen.
- Matched the in-app `sessionLoading` screen to the same branded animated loader and added Playwright coverage that blocks JavaScript to prove the preloader is visible before hydration.
- Fixed mobile web bottom-menu reordering so dragging a tile changes the visible rail order immediately while still saving the final order only on pointer release.
- Updated the mobile Playwright regression to assert live rail movement before release plus persistence after reload.

## 2026-05-29

- Optimized mobile bottom-menu reordering for low-end web mobile: drag pointer coordinates now live in refs, the floating preview is moved through `requestAnimationFrame`, the rail order and `localStorage` are changed only on release, and heavy menu shadows/blur/scale animations are disabled during edit mode.
- Routed web bug-report submissions to the Telegram recipient `185156683` (`AsaselD`) while keeping the local JSONL report log, and added backend coverage for the notification payload.
- Duplicated successful payment notifications to Telegram recipient `185156683` (`AsaselD`) for Telegram Stars, YooKassa, and direct crypto premium activations, with regression coverage for the Stars success path.
- Tightened the latest V2 localization pass for reported Thai/Vietnamese/Georgian-style regressions: Notes/Phrasebook, Roleplay scenarios, Dashboard progress metrics, Premium plan/payment requisites, learning-lab cards, and daily labels now route through guarded localized copy instead of leaking English, Russian, or mojibake text.
- Made the Today daily-quest block more compact by showing the first four focused items instead of the larger full checklist.
- Reworked mobile bottom-rail reordering to match phone-launcher behavior: long-press lifts the tile under the finger, movement immediately changes the actual rail order, edge drag still auto-scrolls, and releasing saves the order and exits drag mode automatically without a separate Done/X control.
- Fixed Learn Words/Review prompt selection for missing interface-language word prompts: the backend now checks the SQLite AI translation cache first, then performs a guarded OpenRouter lexicographic translation into the selected interface language and persists it before falling back to RU/EN only on configuration or model errors.
- Added a regression test that proves a cached interface-language prompt is used before a Russian/English dictionary fallback, preventing the source prompt from drifting between Russian and English when the learner switches UI languages.
- Hardened vocabulary prompt filtering so broken dictionary rows such as `creating -> creating` cannot leak the hidden answer into the question; same-language headwords and same-answer glosses are rejected before cache/OpenRouter fallback.
- Aligned the first-run web onboarding level selector with Settings/backend CEFR support by exposing the full `A1/A2/B1/B2/C1/C2` range.
- Made AI vocabulary translations durable across SQLite rebuilds by writing successful translations back into the matching `vocabulary_words*.json` source dictionary and replaying existing AI translations into `vocabulary_translations` during JSON import.
- Fixed mobile web menu reordering: long-press drag now shows a floating tile preview, opens a visible before/after gap in the horizontal rail, auto-scrolls at the rail edges, commits the order on release, and keeps the saved order after reload.
- Finished the 35-language localization closure for the React V2 web app: required menu, auth, view title/subtitle/action, navigation description, notes, settings, and regression copy now avoids English fallback for every supported non-English interface language.
- Added Playwright coverage that checks the required V2 labels across all 35 interface languages on both desktop and mobile, including the newly added Arabic, Bengali, Czech, Greek, Hindi, Hungarian, Indonesian, Dutch, Swedish, Tamil, Telugu, Thai, Tagalog, Turkish, and Vietnamese UI languages.
- Finished Telegram UI localization for the 15 newly added interface languages by replacing the old English alias fallback with local static copies for main menu, learning actions, word/review prompts, settings, and tool prompts.
- Added Go coverage that fails if any non-English Telegram interface language falls back to English for core menu and tool strings.
- Verified the pass with `go test ./...`, `npm --prefix web-react run build`, and `npm --prefix web-react run e2e`.

## 2026-05-28

- Promoted the React V2 web app to the only active web shell: deleted the old `web/index.html` legacy app, changed `/app/v1` to redirect to `/app/v2/`, and updated public-site app links to point directly at `/app/v2/`.
- Added a fast missing-translation path for Learn Words, Review, Spelling, Vocabulary, and Telegram word cards. If a word has no prompt in the learner's interface language, the app returns the nearest existing dictionary prompt immediately, warms an OpenRouter lexicographic translation in the background, and persists the result into SQLite for future requests.
- Added the new `vocabulary_ai_translations` SQLite table plus indexes and joins into `vocabulary_translations`, so one AI-resolved word-language pair becomes a normal local dictionary entry after the first successful generation.
- Tightened Learn Words loading performance by removing the old prompt-availability SQL join from next-word selection and reusing cached SQLite ID arrays by learning language and CEFR level. Word selection remains limited to the user's learning language and CEFR band.
- Added the landing hero-style generative canvas animation to the React login page on desktop and mobile while keeping the existing Cloudflare, login/password, and Telegram OTP auth flow.
- Fixed light-theme landing section eyebrow pills with an explicit dark text fill and updated regression coverage so the public landing stays readable in light mode.
- Updated regression tests for the V2-only routing, login canvas animation, `/app/v2/` public links, light-theme section labels, and persisted AI vocabulary translations. Verified with `go test ./...`, `npm --prefix web-react run build`, `npm --prefix web-react run e2e`, `npm --prefix site-react run build`, and `npm --prefix site-react run e2e`.

## 2026-05-27

- Reworked the public landing into a Busuu-inspired product flow for Poliglot AI: clear hero, language chooser, course cards, proof stats, concise feature blocks, reviews, FAQ, and pricing while keeping the unique Three.js shader/product demo.
- Expanded the interface language catalog from 20 to 35 languages across the public site, React web app, and Telegram language selector; learning/vocabulary languages remain on the existing dictionary-backed set until new vocabularies are added.
- Regenerated public-site phrase translations for the expanded 35-language catalog and added Playwright coverage for the Busuu-style landing flow, 35-language selector, legal contacts, dark legal language dropdowns, crossed-out sale prices, and desktop/mobile shader rendering.
- Fixed the SQLite vocabulary randomization regression for Learn Words and Review: SQLite now samples from cached matching ID arrays instead of falling back to the first `position >= 0` match, so sparse early CEFR pools no longer repeat the same words/options.
- Added regression tests for sparse SQLite pools where eligible A1 words sit at the beginning of the database and a large non-matching tail follows; the tests require wide variation for both next-word selection and answer distractors.
- Finished the landing completion pass: the public header now uses the generated Poliglot brand logo asset instead of the old text-only mark, the hero keeps one combined premium shader/product animation, and the landing again includes concise reviews, FAQ, feature proof, and discounted pricing with the old one-month prices crossed out before the 70% sale price.
- Regenerated and hardened public-site localization for the landing, Privacy, and Terms across all 20 interface languages; fixed the browser mojibake decoder so normal accented Latin text is not corrupted, and added contact/header parity for the legal pages.
- Added the missing V2 regression closures for the reported mobile issues: compact Awards header/modal behavior, 2/3 Mistakes practice plus scrollable dictionary, leaderboard language dropdown above the panel, stable mobile text/chat bubbles, and refreshed tests for desktop and mobile.
- Verified the production vocabulary path: Learn Words and Review now sample wider CEFR-band distractor pools from SQLite instead of a repeated tiny set, and the production vocabulary database has indexes for language/level/position and translation lookup.
- Updated the regression suite for this pass and verified locally with `go test ./...`, `npm --prefix web-react run build`, `npm --prefix web-react run e2e`, `npm --prefix site-react run build`, `npm --prefix site-react run e2e`, and `node tools/check_encoding_artifacts.mjs`.
- Finished the public landing refresh: removed the visible `/app/v2` nav link, replaced it with a light/dark theme toggle, added the Three.js anomalous-matter shader hero in `site-react/src/components/ui/anomalous-matter-hero.tsx`, and rewrote the landing around the core app value points without overloading the page.
- Added animated CTA/button polish, concise feature cards for Daily route, Roleplay, Pronunciation, Photo/translation, Offline/PWA, and Progress, plus public-site Playwright coverage for desktop/mobile shader rendering, theme toggle, dark legal select readability, and contact cards.
- Changed vocabulary startup so the backend opens `vocabulary.sqlite`, serves `/healthz`, and then runs heavy JSON sync/index creation in the background; deploys no longer keep the public app down while SQLite builds optimization indexes.
- Completed the latest V2 regression pass: browser Back now navigates inside the app on mobile as well as desktop, active mobile Roleplay keeps the localized choose-another-scenario action visible without the bulky scenario header, daily bonus claims are server-locked for 24 hours, and vocabulary prompt selection uses wider randomized pools.
- Added backend coverage for 24-hour daily bonus locking and randomized vocabulary selection, and verified the full pass with `go test ./...`, `npm --prefix web-react run build`, `npm --prefix web-react run e2e`, `npm --prefix site-react run build`, `npm --prefix site-react run e2e`, and `node tools/check_encoding_artifacts.mjs`.
- Added the next V2 auth/mobile/privacy regression pass: award details now render through a fixed portal so mobile taps open the trophy card in the currently visible viewport instead of near the top of the scrolled page; mobile Phrasebook save actions stay inside the composer after voice/photo controls; mobile Roleplay no longer spends space on the selected scenario title block; and the mobile mistake dictionary regained bottom scrolling.
- Reworked standalone login around the shadcn-compatible `SignInPage` component while keeping the existing backend auth contract: login/password, registration, Cloudflare Turnstile, Telegram six-digit OTP, account recovery through `https://t.me/AsaselD`, compact Telegram login row, 20-language auth cards, and required personal-data policy consent for registration and Telegram linking.
- Updated the privacy-policy assets for both public domains and the Telegram bot onboarding: `privacy-policy-i18n.js` now covers all 20 interface languages, the public Privacy page contains the new PoliglotAI personal-data policy with `@AsaselD`, and first-time Telegram `/start` shows Continue/Policy before onboarding.
- Fixed the Privacy asset delivery path after smoke testing: the generated policy asset is clean UTF-8, the Privacy page references it with `?v=20260527-data-policy`, and Caddy no longer treats `/assets/privacy-policy-i18n.js` as an immutable static asset.
- Fixed shared vocabulary distractor generation for web desktop, web mobile, and Telegram so Learn Words and Review no longer recycle the same wrong answers; options are sampled from dictionary words in the same CEFR band (`A1/A2`, `B1/B2`, `C1/C2`) with uniqueness guards.
- Extended regression coverage for the new pass: Go tests cover CEFR-band distractors and bot privacy onboarding, web shell tests cover 20-language privacy-policy assets, and Playwright covers auth privacy consent, mobile award viewport placement, mobile mistake scrolling, compact save-to-notes, payment modal placement, tools send-button layout, and trainer result placement.
- Integrated the shadcn-compatible React `SignInPage` at `web-react/src/components/ui/sign-in.tsx` and switched `/login`, `/app/login`, and `/app/v2/login` to the standalone V2 auth shell while preserving existing login/password auth, Cloudflare Turnstile token submission, and Telegram 6-digit OTP login flow.
- Localized the standalone auth page copy for all 20 interface languages, including login/register hints, account recovery, divider text, referral-code placeholder, Telegram completion status, and Words/auth testimonial labels, so Russian no longer falls back to English on the public login page.
- Added and updated Playwright coverage for the React auth shell, Turnstile slot, Telegram OTP flow, mobile payment modal placement, compact Phrasebook quick-save, mobile Tools input, Today-only mobile quick controls, and the new auth localization keys.
- Verified this pass with `npm --prefix web-react run build`, `npm --prefix web-react run e2e`, `node tools/check_encoding_artifacts.mjs`, `npm run check`, and a public Playwright desktop/mobile `/login` smoke. Deployed the refreshed V2 static bundle to `/opt/aibot/web/v2`; backup saved at `/opt/aibot/deploy-backups/20260527-v2-auth-localization-static`, with public `/healthz`, `/login`, `/app/v2`, and `index-B2GYJ_Vi.js` verified.

## 2026-05-26

- Fixed the latest V2 mobile regressions: onboarding/start dialog now scrolls inside short phone viewports, mobile Phrasebook quick-save sits compactly after the composer, selected Roleplay sessions hide the scenario brief on mobile, mobile Tools hides the large selector after choosing a tool, and More edit mode supports reordering tiles inside More.
- Forced the current FLUX.2 transparent menu assets for Offline, Phrases/Phrasebook, Pronunciation, and Roleplay across desktop and mobile with a versioned `/app/assets/icon-...png?v=...` URL, and moved the desktop ribbon right arrow further inside the app frame.
- Persisted Pronunciation weak-word/history data per account in `localStorage`, so the dashboard survives reloads instead of resetting with the message list.
- Fixed shared Learn Words answer-option sampling for web desktop, web mobile, and Telegram by using a process-level RNG instead of reseeding option sampling from the current timestamp on every call.
- Added regression coverage for mobile onboarding scrollability, current FLUX.2 menu asset URLs, mobile Phrasebook quick-save placement, mobile Roleplay brief hiding, Pronunciation history persistence after reload, mobile More reordering, mobile Tools selector hiding, and Learn Words wrong-answer variation.
- Verified this pass with `npm --prefix web-react run build`, `go test ./... -run TestWordOptionsVaryDistractorsAcrossRounds -count=1`, `npm --prefix web-react run e2e`, and `npm run check`.
- Deployed the latest V2 mobile/Learn Words fix to `/opt/aibot`; backup saved at `/opt/aibot/deploy-backups/20260527-0020-v2-mobile-word-options`, `aibot.service` restarted, and public `/healthz`, `/app`, `/app/v2`, `/app/v1`, `index-CVWyFmYo.js`, `index-DfVI7AHk.css`, current FLUX.2 menu icons, and `brand-assets-manifest.json` verified.
- Added the requested regression test layer for the latest V2 mobile/UI failures: Playwright now covers mobile lesson prompt readability and 2/3 output vs 1/3 input layout, Roleplay session ordering without overlap, mobile Tools selector-first flow, mistake dictionary scrolling, localized More-menu descriptions, editable floating mobile nav, desktop ribbon arrow bounds, Phrasebook API persistence, confirmed-only payments, Offline 10-card pagination with full export, referral invitee pagination, Telegram Send-code visibility, compact pronunciation output, and mojibake guards.
- Fixed the issues exposed while adding those tests: lesson start messages now show `prompt`/task text when the API returns it, practice/roleplay result cards include `correction` and `explanation`, mobile edit-mode floating animation no longer moves the actual tap target, and `/app` now serves the React V2 app while legacy V1 remains available at `/app/v1`.
- Verified the regression pass with `npm --prefix web-react run build`, `npm --prefix web-react run e2e`, `go test ./...`, `node tools/check_encoding_artifacts.mjs`, and `npm run check`.
- Deployed the V2 regression-test pass to `/opt/aibot`; backup saved at `/opt/aibot/deploy-backups/20260526-2340-v2-regression-tests`, `aibot.service` restarted, and public `/healthz`, `/app`, `/app/v2`, `/app/v1`, `index-0dUw8-rz.js`, and `index-CHSeDokX.css` verified.
- Completed the second V2 mobile/referral/offline/payment pass: Premium history now keeps only confirmed payments, Today daily quests replace the old quick-action block and the weekly plan has denser tasks, all output surfaces regained mobile scrolling, Tools stay selectable on phones, mobile quick controls include bug report/language/theme/logout, and mobile More supports image tiles plus long-press pin/unpin editing.
- Added referral invitee visibility in V2: `/api/session` now returns invited learners with joined date, level/XP, level-3 status, and per-user earnings, while the Referrals screen paginates invitees 10 per page and shows the human-readable earned amount.
- Updated Offline decks so saved/offline copy is localized, cards are paginated 10 per page, and TXT/JSON export still downloads the full deck instead of only the visible page.
- Regenerated the Offline, Phrases, Pronunciation, and Roleplay light/dark menu PNGs through the local FLUX.2 ComfyUI workflow (`flux2_dev_fp8mixed`, `mistral_3_small_flux2_bf16`, `flux2-vae`) with transparent 512x512 validation and refreshed `brand-assets-manifest.json`.
- Expanded Playwright desktop/mobile coverage for confirmed-payment filtering, referral invitee pagination, Telegram `Send code` visibility, Offline pagination/export, mobile More image tiles, mobile pin editing, input visibility, Tools selection, scrollability, and mojibake guards.
- Verified this pass with `npm --prefix web-react run build`, `go test ./...`, `npm --prefix web-react run e2e`, and `npm run check`.
- Deployed the second V2 payment/referral/mobile/offline pass to `/opt/aibot`; backup saved at `/opt/aibot/deploy-backups/20260526-223406-v2-referrals-mobile`, `aibot.service` restarted, and public `/healthz`, `/app/v2`, `index-C4nUHB4M.js`, `index-DqC1jeH0.css`, FLUX.2 icons, and `brand-assets-manifest.json` verified.
- Localized the requested V2 menu/system labels across all 20 interface languages, including Roleplay, Phrases/Phrasebook, Offline, Dashboard, Invite/Referrals, V2 learning lab, Daily quests, Weekly plan, Daily route, and Offline decks copy, while keeping the strings clean UTF-8.
- Reworked React V2 mobile navigation: the desktop top menu/ribbon is hidden on phones, primary tabs stay in the bottom bar, the bottom-right More tab opens the rest of the sections, logout is an icon-only top-right control, and the theme toggle is a smaller icon placed to its left.
- Improved V2 layout polish for the requested blocks: pronunciation and phrase icons are larger and centered, the pronunciation map now removes repeated long advice text, Today places the 1-day streak panel neatly to the right of the calendar, the old "Next features" recommendation block was removed, Offline decks use the available width, and mobile dense blocks now scroll normally.
- Added Playwright coverage for the 20-language V2 menu labels, desktop Today/Offline/Pronunciation layout, compact pronunciation map, mobile quick controls, More menu, and mobile scroll behavior.
- Changed Nx Go test execution to use a local vocabulary directory during tests, avoiding production `.env` paths such as `/opt/aibot` during Windows `npm run check`.
- Deployed the V2 localization/mobile cleanup build to `/opt/aibot/web/v2`; backup saved at `/opt/aibot/deploy-backups/20260526-205908-v2-mobile-localization`, `aibot.service` restarted, and public `/healthz`, `/app/v2`, `index-DG4RaOUL.js`, and `index-BFuJzIXN.css` verified.
- Added the next V2 web learning slice: Roleplay, Pronunciation, Offline, and Phrasebook now use their own generated menu artwork targets; V2 has a personal Phrasebook/Favorite phrases surface sourced from lesson/practice messages, Offline decks can group vocabulary, phrasebook, and mistake cards, the Today screen includes a monthly habit calendar plus modest daily plan, and completed daily goals can grant a once-per-day bonus XP animation.
- Added authenticated bug reports from the V2 top bar. Reports accept a text description plus optional screenshot, save uploaded files under `bug_reports/screenshots`, and append metadata to `bug_reports/bug_reports.jsonl` until mail delivery is added later.
- Added React Aria calendar UI, V2 Playwright desktop/mobile smoke coverage for Today, Phrasebook, Offline grouping, bug reports, Roleplay, Pronunciation, and mobile navigation, plus documentation for desktop/mobile V2 functionality.
- Installed and configured the read-only MySQL MCP server as `mysql-readonly` in the local Codex config, and documented startup/env-var/SELECT-only usage in `MCP_STARTUP.md`.

## 2026-05-24

- Added lesson anti-repeat memory: lesson generation now stores the last 10 lesson prompts, sends them back into the lesson prompt, rotates CEFR-appropriate topics when no global focus is set, and keeps practice context to the last 5 learner messages. V1 award rank/story rendering now uses the same canonical polyglot stories as V2 and no longer rejects clean Russian award text as mojibake.
- Deployed the lesson anti-repeat / award-story sync pass to `/opt/aibot`, including the rebuilt backend, updated legacy `/app` HTML, current React `/app/v2` build, and `app_prompts.json`; verified production `/healthz`, `/app`, `/app/v2`, and V2 assets.
- Fixed the V2 encoding regression again: React fallback text and docs no longer contain mojibake artifacts, `cleanAppText` repairs mixed cp1251/UTF-8 fragments using escaped detector patterns, and `node tools/check_encoding_artifacts.mjs` now covers React V2 files. Tool, lesson, practice, roleplay, listening, word, spelling, level, mistake, and translator messages are now tagged by surface so each V2 tool shows only its own message history. Roleplay now uses the dedicated `ROLEPLAY_TOOL_V2` flow without showing a premature "next question" block before the learner replies. Settings now exposes a global learning focus field saved to the backend and injected into lesson/practice/roleplay prompts.
- Deployed the V2 encoding/roleplay/focus fix to `/opt/aibot`, uploaded the updated `app_prompts.json`, replaced `/opt/aibot/web/v2`, restarted `aibot.service`, and verified server-side `/healthz`, public `/app`, public `/app/v2`, and fresh V2 asset `index-BaaRj7Zt.js`.
- Completed the V2 output/roleplay/pronunciation/offline/dashboard pass: each menu destination now gets a view-specific context surface, Roleplay enters a dedicated dialogue scene after scenario selection, scenario tiles form an even desktop two-row grid, dense V2 blocks are scrollable, Offline deck cards render as full-width rows, Dashboard no longer duplicates leaderboard/weak-topic blocks, and Pronunciation no longer exposes raw technical issue keys.
- Tightened exact-repeat Listening/Pronunciation scoring in Go: target mismatch, missing words, substitutions, and estimated confidence now cap the score more aggressively while keeping feedback learner-facing.
- Deployed the pass to `/opt/aibot`, restarted `aibot.service`, and verified production `/healthz`, `/app`, `/app/v2`, plus fresh V2 assets `index-BljsgA2m.js` / `index-CIAalwsj.css`.

## 2026-05-23

- Completed the V2 vocabulary/roleplay/pronunciation correction pass: Learn Words now honors the selected interface-language prompt word when a translation exists, Vocabulary no longer performs per-card synchronous AI example generation during list load, AI Roleplay has ten localized scenarios and stays in its own output surface with a dedicated backend roleplay prompt, Pronunciation exact-repeat scoring is stricter and no longer praises clearly wrong repeats, Offline decks export readable TXT study packs, Global leaderboard rows show language names, Progress is merged into Dashboard, and V2 now uses static light/dark backgrounds. React build emits fresh V2 assets `index-BYNw5Trl.css` / `index-DwXYwmtU.js`.
- Deployed the V2 vocabulary/roleplay/pronunciation correction pass to `/opt/aibot`, uploaded `app_prompts.json` next to the binary, restarted `aibot.service`, and verified production `/healthz`, `/app`, `/app/v2`, plus the fresh V2 JS/CSS assets.
- Added runtime-editable AI prompts: `APP_PROMPTS_FILE` now defaults to `app_prompts.json`, the server loads prompt templates at startup, core lesson/practice/translator/vocabulary/listening/pronunciation prompts can be edited without rebuilding the binary, and `docs/reference/APP_PROMPTS.md` documents feature/tool/function ownership for navigation. Also fixed the V2 Pronunciation/Pronounce tile so it opens the pronunciation dashboard instead of auto-starting Listening/Audition. React build emits fresh asset `index-BUoc4oQk.js`.

- Completed the fourth V2 learning-app audit slice: the default screen is now the visible `РЎРµРіРѕРґРЅСЏ` dashboard, Home is exposed in the ribbon/mobile nav, Offline and Dashboard preload data without redirecting users into Vocabulary/Mistakes, Settings Telegram 2FA uses a global centered portal popup, the post-registration onboarding dialog collects goal/level/language/format, Mistakes now groups entries by grammar, word order, vocabulary, politeness, and spelling with a `РџРѕС‚СЂРµРЅРёСЂРѕРІР°С‚СЊ РїРѕС…РѕР¶РёРµ` drill action, Premium now records local payment history and shows webhook/auto-check success entries, and the Today screen highlights the new V2 learning lab modes. React build emits fresh V2 assets `index-DIhZOhmh.css` / `index-BJCVGuii.js`.
- Deployed the fourth V2 audit slice to `/opt/aibot/web/v2`, restarted `aibot.service`, and verified internal/public `/healthz`, `/app/v2`, and the fresh V2 JS/CSS assets from the server.
- Added the next point-by-point V2 roadmap slice: Home now has a real weekly learning plan and daily quests, the function ribbon exposes AI Roleplay, Pronunciation, Offline Decks, and Teacher Dashboard views, Pronunciation shows a heatmap/history surface from voice assessment data, Offline Decks save browser mini-decks and register a PWA service worker/manifest, Dashboard summarizes progress, activity, Premium, referrals, leaderboard, and weak topics. React build emits fresh V2 assets `index-Cwac7skq.css` / `index-DTL-RdVj.js`.
- Deployed the point-by-point V2 feature slice to `/opt/aibot`, rebuilt the Go backend so `/app/v2/manifest.webmanifest` and `/app/v2/offline-deck-sw.js` are served as explicit no-cache PWA files, restarted `aibot.service`, and verified internal plus public `/healthz`, `/app/v2`, V2 JS/CSS, manifest, and service worker endpoints.
- Added the next V2 implementation slice after the captcha handoff: mobile gets its own bottom navigation row, function-ribbon neighbors now make room without shrinking/fading, header status banners stay above the ribbon layer, Settings cards align to content height, translucent panels show the animated background more clearly, and the latest React build emits fresh V2 assets `index-6B3MD3Yv.css` / `index-DsYghMDs.js`. Local verification passed; live authorized QA is still blocked because the current Codex in-app browser session remains on the anti-bot login challenge.
- Deployed that V2 slice to `/opt/aibot`, restarted `aibot.service`, refreshed Caddy so public `/healthz` proxies to Go, and verified production `/healthz`, `/app`, `/app/v2`, and the new V2 JS/CSS assets.
- Started implementing the full V2 improvement roadmap with the Settings two-factor slice: `Send code` now creates the same Telegram bot deep-link flow as V1, opens the centered OTP dialog as a separate popup layer, verifies the site code through `/api/auth/telegram/status`, refreshes the web session after success, and keeps resend tied to a fresh Telegram auth request. Cleaned OTP/payment visible copy and added regression coverage for the Telegram code flow.
- Completed the final V2 menu/voice/payment refinement pass: menu neighbors no longer shrink or fade on hover, the recording button uses the animated mic/stop/wave/timer pattern while keeping the current visual style, payment success and not-found states open as centered dialogs with plan/period/expiry details, payment instructions were cleaned, mistake-dictionary audio uses the shared waveform, V1 award ranks/stories are restored as canonical fallbacks, Settings OTP opens in a centered mini dialog, animated backgrounds show through more transparent panels, and mobile ribbon/audio/control spacing was tightened.
- Deployed the final V2 refinement to `/opt/aibot`, cleaned stale `/app/v2` hashed assets on the server, restarted `aibot.service`, and verified production `/healthz`, `/app`, `/app/v2`, and the new V2 JS/CSS assets.
- Completed the follow-up React `/app/v2` polish pass from the user's screenshot review: menu hover now shifts only after hovering a concrete item, compact audio/image attachment buttons match the microphone control, ribbon arrows are inset, Listening next action sits bottom-right, word/spelling result blocks remove duplicate "task ready"/context/example text, level scores show the real max once, award modals ignore placeholder rank/story keys, chat output auto-scrolls, payment requisites include exact-amount instructions/expiry/check action, Settings stacks Activation under Telegram, and language menus render above the header context layer with slower drop-down animation.
- Completed the strict 21-point React `/app/v2` QA pass: top-center 3-second system status, animated/dark-safe language dropdowns, expanding function ribbon hover, smooth context panel drop-in, mistake clear confirmation/no start button/correct-answer audio, trophy names and polyglot stories, dictionary examples without review counters, restored next/audio controls in spelling/listening/word trainers, stable audio-wave rendering, animated upload controls, spinner loading states, Settings password change, V1 Telegram referral invite copy/terms, login-based header identity, and backend audio text for lesson corrections plus practice follow-up questions.
- Updated private React `/app/v2` around the latest QA checklist: translator language swap, scrollable chat output, explicit Global top leaderboard button, top ribbon chips that visually expand from their center, compact two-column vocabulary cards, and per-word GPT-4o Mini TTS playback through the existing backend pronunciation endpoint.
- Replaced voice recording surfaces with the provided animated `VoiceInput` pattern and replaced listen/playback surfaces with waveform audio controls that request backend TTS instead of browser `speechSynthesis`; voice/text/image submissions continue to append learner-side messages into the chat stream.
- Simplified V2 backgrounds to the requested Paper Shaders background only, added header brightness control, wired the brand click to Settings and the XP/trophy click to Awards, refreshed the polyglot award story fallback, and integrated the requested animated theme toggle, morphing arrows, joly-style buttons, payment morph buttons, and logout confirmation animation.
- Expanded Premium payment UX so selected payment methods remain visually highlighted/loading, and returned payment details render in the same modal with payment ID, status, amount, network, wallet, comment/memo, expiration, transaction hash, and payment link when present.
- Updated Settings and trainer flows: Telegram code send opens a 6-digit OTP modal, save/activation controls use the requested animated button styles, and the level test removes typed/free-form answers while Skip sends the same `-1` backend answer used for "I do not know".
- Verified this pass with `npm run build` in `web-react` and `go test ./...`; local browser opening reached the V2 static app but unauthenticated preview redirects to `/login`, so authenticated visual QA still needs a real session or mock API server that remains attached.

## 2026-05-22

- Fixed the legacy `/app` v1 referral display so the live TonAPI-backed `usdt_rub_rate` from `/api/session` is shown in referral balance cards and withdrawal hints instead of a stale hidden/fallback rate.
- Updated private React `/app/v2` with the latest QA fixes: Premium plan buttons now open a centered payment modal, the current trophy sits left of LVL in the header, Home/Limits are hidden from the top function ribbon, lesson/word/spelling/vocabulary/mistake/leaderboard flows auto-start or auto-load, Enter submits in chat/tool/trainer inputs, and user replies now appear in the conversation stream.
- Restyled v2 around the requested animated AI chat feel and gradient function menu: fixed-size context panels, scrollable output, compact input, voice/listen buttons inside chat bubbles, centered award story modals, referral Telegram/copy/share actions with USDT/RUB, and an app-wide Paper Shaders animated background for both light and dark themes.
- Reverified the latest v2 pass with `npm --prefix web-react run build`, `go test ./...`, `node tools/check_encoding_artifacts.mjs`, in-app browser console checks, and Playwright/Chrome smoke passes at 1440x900, 1366x768, 390x844, and 430x932.
- Reworked private React `/app/v2` again around the requested narrow header, XP/current-award progress, fixed language selector/theme/logout controls, large image tile ribbon with working horizontal scroll, and one full-height terminal-style context panel with no extra internal header.
- Restored v2 functional parity details: lesson/practice/tools/Listening replies now surface structured AI fields, word/spelling/level/mistake trainers show answer state plus result blocks, Listening submits the target phrase with mic/file audio, vocabulary uses 10-word pagination, and awards are clickable with unlocked/locked states.
- Fixed v2 mobile usability: large tile artwork remains visible, status banners no longer cover controls or intercept send buttons, and the mobile task/result flow keeps the answer input and returned result in the same context area.
- Reverified v2 with clean `npm --prefix web-react run build`, `go test ./...`, `node tools/check_encoding_artifacts.mjs`, and repeated Playwright/Chrome smoke passes across 1440x900, 1366x768, 390x844, and 430x932.
- Rebuilt private React `/app/v2` as a Tailwind/shadcn-compatible standalone Vite app with an iOS-style horizontal function ribbon, desktop/mobile-specific layouts, a single central context display, dark `FallingPattern` background, light `ShaderBackground`, and real Go API wiring for lessons, practice, Listening, words, spelling, level test, vocabulary, mistakes, tools, Premium, referrals, limits, progress, awards, leaderboard, and settings.
- Added the React v2 shadcn foundation under `web-react/src`: `@` alias, `lib/utils.ts`, `components/ui/button.tsx`, `falling-pattern.tsx`, `shaders-hero-section.tsx`, and `chat-interface.tsx`; v2 logout and `/app/v2/login` continue to leave users on the old `/login`, and the legacy `/app` shell still exposes no public v2 link.
- Added browser-side mojibake repair to React v2 localization so old 20-language strings from the API or local fallback maps render as readable UTF-8 instead of broken mojibake text.
- Deployed the rebuilt Go backend and private React v2 app to production under `/app/v2`, preserving public landing, Terms, and Privacy checksums and keeping `/app/v2/login` redirected to `/login`.
- Kept the public landing, Privacy, and Terms deployment path on the previous static site while v2 is still private, and removed the visible v2 launch link from the legacy `/app` shell.
- Restored the legacy `/login` as the default unauthenticated/login/logout destination: React `/app/v2/login` now redirects to `/login`, and React logout leaves v2 for the old login page.
- Restored the desktop inspector serial-key form in the legacy web app and wired it to the same activation-key API/status path as Settings.
- Expanded the legacy login fallback copy and language list so the standalone auth page can render the supported 20 interface languages even before a session is loaded.
- Reworked React `/app/v2` tool and mistake flows: Tools now call the real translator, voice-to-text, and image-translate endpoints; Mistakes now supports loading, starting a specific repair drill, submitting the corrected answer, clearing the dictionary, and visible success status.
- Tightened React `/app/v2` desktop/mobile separation: Vite no longer proxies its own `/app/v2/assets`, build output is cleaned on each build, desktop keeps sidebar/workspace/inspector, and mobile uses topbar plus fixed bottom navigation without horizontal overflow.
- Verified the v2 pass with `npm --prefix web-react run build`, `go test ./...`, `node tools/check_encoding_artifacts.mjs`, and a Playwright/Chrome smoke test covering desktop Tools, Mistakes drill, and mobile Home layout.
- Rebalanced the React public landing, Privacy, and Terms pages after visual QA: the first block now uses the lighter `SparklesCore` particle treatment, the landing keeps `BackgroundPaths`, pricing, and testimonial components only where they help the structure, and oversized legal/landing headings are constrained for desktop and mobile.
- Preserved the previous Privacy and Terms legal body copy 1:1 inside the new React legal shell, while reducing decorative motion so the pages stay readable and all existing information remains present.
- Fixed the React `/app/v2` shell at 100% zoom: desktop content now fits the sidebar/workspace/inspector grid, mobile login no longer clips copy or auth buttons, and the Practice/Listening action buttons route through working handlers instead of inert placeholders.
- Completed browser-side auth/login copy coverage for all 20 interface languages used by the React `/app/v2` shell.
- Cleaned `/api/session` web language DTOs so login/settings selectors receive readable names for all 20 languages instead of legacy mojibake display strings.
- Rebuilt the public landing, Privacy, and Terms as a React/Vite multi-page site under `site-react`, with Tailwind CSS, TypeScript, shadcn-style `components/ui`, integrated Sparkles/BackgroundPaths/pricing/testimonial components, and output assets under `assets/site-react`.
- Added the `site:build` script and Nx `site-build` target so the public React site can be built and deployed separately from the Go backend and `/app/v2` React app.
- Added local 20-language browser copy for the React `/app/v2` interface and auth language selection, while keeping the Go API/backend as the source of truth for account data and learning state.
- Added regression checks that require the React app localization layer and verify the public React landing/legal build outputs for all three public pages.
- Split web app deployment into legacy `/app` v1 and React `/app/v2`: Vite now builds into `web/v2`, Go serves a separate v2 SPA route and v2 asset route, Caddy matchers explicitly include `/app/v2`, and v1 shows a fixed `New React app` launch button.
- Replaced the legacy inline `/app` frontend with a Vite/React 19 application in `web-react`, using the 21st.dev-generated dashboard concepts as source material: animated dot-matrix background, glass cockpit shell, premium desktop sidebar, mobile top/bottom navigation, auth shell, learning modules, chat/output stream, and Premium/Settings surfaces.
- Connected the React build to the existing Go backend without replacing API logic: `npm run web:build` emits `web/v2/index.html` plus hashed assets under `web/v2/assets`, and the Go static asset handlers now serve Vite JS/CSS/font/image files with correct MIME types.
- Updated web shell regression tests to validate the built React entrypoint and React-to-Go API contracts instead of legacy inline-script markers.
- Added `MCP_STARTUP.md` with the default startup checklist for Figma, 21st.dev Magic, Ref, shadcn, Nx MCP, and Linux/Ubuntu access across projects.
- Documented the persistent Nx MCP startup command for this workspace, including the Streamable HTTP mode on port 9921 and the difference from the stdio Codex config command.
- Added a disk-conscious ComfyUI workflow package for Shorts/Reels/TikTok ads: FLUX.2 keyframes, Wan 2.1 I2V FP8 motion, CLIP Vision image lock, native RIFE interpolation, voiceover combine, optional RealESRGAN upscale, model manifest, and guarded installer.
- Added the first elite Poliglot AI 20-second English ad script/prompt for ComfyUI/XTTS, with a business-suited young adult presenter, web-app visuals, and a four-shot 5-second production plan.
- Added the first two 5-second ComfyUI ad test variants and a queue script: one desktop web-app presenter clip and one mobile voice/pronunciation presenter clip.
- Softened the web app glass cards/panels on desktop and mobile, left-aligned the desktop Tools voice-limit card, and moved the mobile Tools active-status pill onto its own line so it no longer overlaps the voice counter.
- Improved the mobile Referrals hero text flow and balance card width so referral copy reads in normal lines instead of collapsing into a narrow column.
- Added clean web referral and vocabulary-empty localization overrides for all 20 interface languages, preventing Chinese and other languages from falling back to English in Referrals and Vocabulary empty states.
- Changed Telegram tool-mode back buttons for voice/text and image translation to return to the Tools menu instead of jumping to the main menu.
- Changed referral Premium rewards so the invited user still receives 7 days immediately, while the inviter receives 7 days only after that invited user reaches XP level 3; added storage migration and regression coverage for one-time reward issuance.
- Refreshed the public Privacy Policy with a concise Russian-law-focused structure, preserved the existing legal-page design/assets, and added rendered policy copy for all 20 site languages.
- Tuned web app header mini-cards and the Tools voice-limit card to use softer semi-transparent glass on mobile and desktop, added the Site block with `poliglotai.ru` and `poliglotai.com` to Terms/Privacy, and fixed dark-theme legal language selects/update notes across all site languages.
- Fixed web Premium payment actions: Telegram Stars now opens only in a new Telegram tab without replacing the current app page, and direct TON/USDT invoices auto-check quietly every 15 minutes only while still pending, while the manual "Check payment" button remains immediate.

## 2026-05-21

- Added current USDT/RUB lookup through TonAPI `/v2/rates`, with cached rates in referral balances and dynamic USDT invoice calculation when explicit `CRYPTO_USDT_*_AMOUNT` values are empty.
- Changed web Telegram Stars checkout to open the bot through a product deep link (`/start buy_...`), so mobile browsers and desktop Telegram can continue the Stars invoice inside the bot.
- Localized the web paginator next button across all 20 interface languages and added right/left arrow labels for vocabulary and mistakes pagination.
- Added regression coverage for weighted CEFR level-test scoring and Telegram Stars deep-link payloads.
- Added a read-only Telegram account identity line to web Settings on desktop and mobile, showing the linked Telegram name and ID without offering unlink controls.
- Fixed Russian award modal copy in the web app by allowing clean Cyrillic rank/story text through the safe-display filter instead of falling back to English.
- Made the desktop right-inspector serial-key form show the same activated/used/not-found result in a centered confirmation card while keeping the inline Settings-style status text.
- Reworked the Telegram translator language picker into larger single-column buttons with native names plus an interface-language hint in parentheses, so scripts like `ж—Ґжњ¬иЄћ` are understandable from Russian and other UI languages.
- Localized Russian Telegram translator prompts, removed stale tool prompts when switching between image/text translation modes, and tracked/deleted previous pronunciation/translator audio messages before sending a new one to avoid Telegram replaying old audio in sequence.
- Fixed the web mistakes paginator so inactive previous/next controls are not rendered and the next button no longer shows broken mojibake text.
- Added a compact activation-key form directly under the desktop web account inspector panel; Enter submits the key through the same Premium/Platinum activation path as Settings.
- Made the CEFR level assessment resume active progress after page refresh and shuffle question order on new attempts while keeping each attempt's order stable across web and Telegram answers.
- Localized award rank names and polyglot story text for all 20 web interface languages, adding clean CJK/Armenian/Georgian overrides and a 20-language validation pass so broken legacy mojibake cannot leak into the trophy modal.
- Added the current trophy thumbnail plus XP level/progress chip to the desktop web header, matching the mobile shortcut behavior and opening Awards on click/keyboard activation.
- Added a clean 20-language Telegram reminder layer for daily notifications and reminder settings, plus regression tests that verify every supported interface language sends localized reminder text without mojibake.
- Cleaned mojibake from the web app, recorder statuses, static composer buttons, and Go API error messages; added `tools/check_encoding_artifacts.mjs` so common broken Cyrillic, smart-quote, and emoji artifacts are caught before deploy.
- Added one-time serial Premium/Platinum keys: a separate `activation_keys.sqlite` database, `activation_keys.txt` import file, web Settings activation UI with 20-language copy, clear already-activated errors, `tools/generate_activation_keys.mjs`, and the root-only Ubuntu `poliglot-keys` helper for `XXXX-XXXX-XXXX-XXXX` month/year key generation.
- Added optional Cloudflare Turnstile protection for web login/registration/profile setup, a global `/api/*` per-IP/per-session rate limiter, and a browser-side short click guard for action buttons.
- Localized the Telegram `Listening` / `РђСѓРґРёСЂРѕРІР°РЅРёРµ` menu button through the shared UI copy instead of hardcoding Russian text, and made web Settings language changes redraw the active screen immediately.
- Finalized the desktop/mobile app header-art balance: lesson, words, review, spelling, level, and settings headers now keep separate desktop/mobile placement rules, with desktop art reaching the right edge and mobile art filling the full header block behind unchanged text/progress content.
- Cleaned the public Privacy/Terms pronunciation copy and dark-theme legal styling: removed provider/model/JSON/confidence implementation details from user-facing legal text, kept the first hero as one branded background image, and made pronunciation/contact cards readable on dark backgrounds.
- Prevented practice pronunciation audio from reading interface-language translations in parentheses; `Model phrase` is now prompted and sanitized as the target-language-only sentence for both Telegram and web.
- Routed public `/login` through Caddy to the Go web app, so the standalone auth page is available on `poliglotai.ru/login` instead of falling back to the landing page.
- Tightened the final mobile bottom-nav rule with fixed viewport bounds and verified it at a 390px mobile viewport; the "More" button no longer clips off the right edge.
- Added a standalone `/login` / `/app/login` auth page with theme and interface-language controls, redirecting unauthenticated `/app` visitors there and authenticated login visitors back into `/app`.
- Fixed Telegram profile completion so entering the login/password of an existing web account opens the merge-choice flow instead of only reporting that the login is taken.
- Hardened web app view state: each section keeps its own rendered content and composer draft when users move around, and stale async responses can no longer redraw a different active section.
- Reworked the final desktop/mobile section-header CSS so generated art uses a reserved masked visual lane, while menu card icons, mobile bottom navigation, and the mobile composer stay stable.
- Restyled the `/app` login/profile-completion screen with current Poliglot AI branding, logo artwork, clean Russian copy, Telegram-login card, and mobile auth rules that hide the bottom navigation instead of letting it cover the form.
- Reworked web app section headers again so desktop header art has its own visual zone, mobile header art remains visible as a masked layer, and focus/stat cards do not sit directly on top of the main copy.
- Added photo/camera input to the Practice composer in the web app and Telegram practice photo handling: image context is analyzed by the OpenRouter vision path and fed to the practice coach without treating the photo text as learner writing.
- Made the desktop account chip and mobile login pill open Settings, matching the requested header shortcut behavior.
- Added Telegram/web account merge choice: when both sides already have progress, the web app asks whether website or Telegram progress should be primary, then merges the other side into it.
- Softened section header artwork so generated images behave as one translucent masked layer instead of a rectangular photo block; mobile headers now show a larger visible masked image, while menu card icons are smaller and less intrusive on narrow screens.
- Restored full desktop Listening/Audition result text by using compact pronunciation cards only on mobile; desktop no longer truncates the pronunciation feedback with an ellipsis.
- Made the mobile top-bar Premium pill open the Premium section and the level/trophy progress chip open Awards.
- Removed stale `bot.html` online-chat prototype links from Terms and Privacy navigation/footer links and pointed them to the real `/app` web application.
- Restored Russian web-app localization by repairing mojibake copy in the browser instead of dropping it to English fallbacks, and added a visible-copy repair pass for dynamic auth/settings/tool messages.
- Split the web app shell into forced mobile/desktop delivery modes (`/app?shell=mobile`, `/app?shell=desktop`, and automatic user-agent selection) so mobile CSS fixes no longer force desktop layout changes.
- Reworked generated header artwork placement again as a larger masked no-repeat layer with separate desktop/mobile scale rules, keeping the art visible while allowing controlled overflow outside the header block.
- Compactified the mobile Listening/Audition result card so transcript, feedback, pronunciation score, accent, fluency, and the first weak word stay scannable instead of pushing the whole page downward.
- Fixed practice correction audio to prefer the full localized `Model phrase` block, so the first generated voice speaks the whole corrected phrase instead of only the final corrected word.
- Added the Yandex.Metrika `109326597` counter to the web app shell.
- Fixed Listening/Audition pronunciation scoring so OpenRouter STT no longer receives the exact target phrase as a recognition prompt; the recognizer now transcribes what the learner actually said instead of being biased toward the answer.
- Capped pronunciation scores when OpenRouter does not return word-level confidence/logprobs, so Telegram/web voice repeats cannot be marked perfect from transcript-only evidence; feedback now explains the stricter estimated-confidence mode.
- Localized lesson/practice prompt section labels across the 20 interface languages and added web-side normalization for old English headings like `Situation`, `Your task`, and accidental `saturation` leakage.
- Hardened the mobile web app shell: header artwork is a visible masked no-repeat layer, the global background is visible on mobile, bottom navigation stays fixed, award modals stay centered, the composer uses full width, and the logo/title tap target alone opens the menu.
- Added `docs/tracking/USER_UI_WISHES.md` as the recurring checklist for generated art, mobile layout, localization, pronunciation scoring, and result-scroll behavior before future web/app changes.
- Reworked the public landing and Figma landing page away from internal speech-scoring mechanics toward sales-focused benefits: hear, repeat, understand weak spots, and continue learning in Telegram or web app.
- Created a new editable Figma design file `Poliglot AI Landing 2026` with desktop and mobile landing frames covering refreshed branding, product sections, pronunciation scoring, 20-language positioning, reviews, pricing, Privacy, and Terms.
- Refreshed the public landing with the same structure: updated hero copy, added the pronunciation-engine explanation, added a reviews section, added a Privacy/Terms product-trust band, and expanded navigation/localization for the new sections.
- Updated Privacy and Terms to describe the OpenRouter GPT-4o Transcribe pronunciation flow, text-only Gemini coach report, and the educational/non-exam nature of pronunciation scoring.
- Renamed the learner-facing Shadowing surface to Listening / `РђСѓРґРёСЂРѕРІР°РЅРёРµ` while keeping `/shadowing`, `/repeat`, and internal route names compatible.
- Switched OpenRouter speech recognition defaults and env examples to `OPENROUTER_STT_MODEL=openai/gpt-4o-transcribe`; added `OPENROUTER_PRONUNCIATION_MODEL=google/gemini-3.1-flash-lite` for the cheap text-only pronunciation coach.
- Added pronunciation assessment to every voice-learning path: Listening exact repeat, Telegram lesson/practice voice replies, and web lesson/practice voice replies now normalize audio, transcribe through OpenRouter STT, compare transcript/confidence/fluency, and return score, accent strength, weak words, and tips.
- Added voice recording/file input to the web lesson and practice composer, with structured pronunciation result cards after AI feedback.
- Added a new Shadowing learning mode for Telegram and the web app: the learner listens to one generated target phrase, repeats it aloud, gets STT-based similarity scoring, concise pronunciation/rhythm feedback, XP, and a next-phrase action.
- Added `/shadowing` and `/repeat` Telegram commands, a main-menu Shadowing button, web `/api/shadowing/start` and `/api/shadowing/answer` endpoints, a dedicated web screen with recording/file upload, and shared local scoring tests.
- Strengthened lesson and practice prompts around active recall, reusable chunks, highest-impact correction, and adaptive recycling of recent weak phrases while staying on the same OpenRouter model path.
- Added Flux2 prompt definitions for unified Shadowing header/icon artwork and connected the app-wide background plus right account panel background as no-repeat cover/masked layers in the web shell.
- Rebuilt the ComfyUI prompt rules for menu icons and awards so generated button artwork is function-specific, transparent-background, same-scale, and checked per file before moving to the next asset.
- Added a transparent asset post-processing validator for icon and trophy PNGs: chroma-key removal, fixed square canvas sizes, normalized object fill, transparent corner checks, narrow-object rejection, and tile/backplate rejection.
- Fixed generated header/background placement so section artwork renders as one masked no-repeat image, added a generated background layer to the Tools hero, and gave Settings the same top header treatment as other app sections.
- Regenerated and validated 34 transparent menu icon PNGs and 20 transparent trophy PNGs, refreshed the manifest validation metadata, and checked the result with contact sheets plus local `/app` screenshots.
- Bumped the `/app/assets/...` cache-buster from `flux2-20260520` to `flux2-20260521` so deployment of the regenerated icon/header/trophy pack forces browsers to fetch the new PNG files.

## 2026-05-20

- Expanded the Flux2 web refresh beyond the original plan/header/icon pack: added app-wide light/dark background images, right-inspector panel backgrounds for Progress/Limits/Account, regenerated Free plan artwork, generated 20 trophy PNGs, and restored larger main-menu card art with glass overlays.
- Added Tools and Settings cards back into the main web menu with generated icon artwork, and changed the desktop home/sidebar entry to use the generated main-menu icon path instead of the stale logo image.
- Strengthened Flux2 app-background and panel prompts to avoid fake UI text, dashboard rows, labels, chat-message lettering, avatars, and numbers, then regenerated the global backgrounds and right-side account panel backgrounds with clean negative space for real HTML text.
- Added reminder regression coverage proving daily Telegram reminder text follows the user's selected `interface_language`, including a changed-language SQLite target and localized compact fallback languages.
- Rebuilt and redeployed the Linux bot binary after the reminder localization checks so production Telegram reminders run with the current language-aware scheduler code.
- Completed the FLUX.2 Dev test-art refresh for the web app with 74 theme-aware premium PNG assets: Free/Premium/Platinum plan cards, section headers, and action/navigation icons now use separate light and dark files generated through the ComfyUI Flux2 split workflow.
- Deployed `tmp/poliglot-web-flux2-refresh.tar.gz` to `/opt/aibot/web` on `root@186.246.45.123`; production `/app` now serves the Flux2 manifest and matching Flux2 PNG hashes.
- Added a `v=flux2-20260520` cache-buster to `/app/assets/...` references in `web/index.html` so browsers that cached the previous week-long PNG responses fetch the refreshed web app visuals immediately.
- Reworked weak generated concepts for Vocabulary, Learn Words, and Level Assessment to remove folder-like trays, blank placeholder icons, gauge scales, tick marks, and pseudo-text-like micro details.
- Rebuilt the lesson/words/word-game/spelling/level content surfaces as unified `task-stage` panels so their backgrounds and header artwork match the rest of the premium web app design.
- Removed stale unthemed generated assets (`header-*.png`, `icon-*.png`, old plan PNG/SVG files, and replaced public `site-*` section images) while keeping backend compatibility for legacy header/icon/plan URLs by serving the light-theme asset as a fallback.
- Verified the refreshed pack locally with contact sheets, app/public asset audits, `node --check`, inline web script syntax checks, `go test ./...`, and browser screenshots for landing plus `/app` desktop/mobile light/dark.
- Began the Flux2 test-art refresh for the web app: fixed the duplicated dark-theme header image layer, removed secondary decorative image overlays from main action cards, strengthened the no-folders/no-screens/no-pseudo-text prompt system, and added a FLUX.2 Dev ComfyUI workflow path to `tools/generate_comfy_brand_assets.mjs`.
- Added local installer scripts for the FLUX.2 Dev ComfyUI model pack (`flux2_dev_fp8mixed.safetensors`, `mistral_3_small_flux2_bf16.safetensors`, and `flux2-vae.safetensors`) while documenting that the full BFL `flux2-dev.safetensors` fp16 file is gated behind Hugging Face access.
- Restored the user-uploaded public landing imagery with the fixed mapping: `poliglot-ai-avatar.jpg`, `times.jpg`, `primi.jpg`, `write.jpg`, and `yspex.jpg`; the online bot page now also uses `poliglot-ai-avatar.jpg` for its hero/Open Graph image.
- Added `docs/tracking/PROJECT_TRACKING.md` as an Obsidian-style project board with sprint status, asset matrix, QA checklist, deployment checklist, and decision log.
- Changed web app plan/header/icon artwork to theme-aware production files (`*-light.png` and `*-dark.png`) and updated `/app` CSS/HTML to switch them from `html[data-theme]`.
- Removed decorative pseudo-letter CSS marks from the web app stage, action-card, and panel visual treatments.
- Updated the local asset manifest with theme, model, seed, prompt, dimensions, source, and QA status metadata, plus a generated contact sheet at `tmp/brand-assets-theme-contact.png`.
- Reworked `tools/generate_comfy_brand_assets.mjs` so SD3.5 Large FP8 is the default ComfyUI workflow with SD3.5-specific nodes, while SDXL/Juggernaut remains available as a fallback/comparison path.
- Added `docs/product/REDESIGN_PLAN.md` as the source-of-truth design plan for the Poliglot AI 2026 `Language Intelligence Cockpit` redesign.
- Created a Figma redesign board for the product direction, tokens, app shell, mobile shell, landing, legal pages, and required asset matrix.
- Rebuilt the final visual asset pack as one consistent no-text/no-pseudo-letter system for plan cards, app headers, app icons, and public-site hero/legal/section visuals.
- Added `tools/create_poliglot_premium_assets.py` to generate deterministic premium PNG assets without fake text, screens, flags, people, or broken AI lettering.
- Updated the web app desktop and mobile visual shell with the final cockpit palette, calmer grid background, stronger glass panels, refined auth screen, and improved generated-asset placement.
- Updated the public landing, online bot page, Terms, and Privacy visual system to use the new site assets, darker premium hero treatment, refreshed navigation, legal shells, and unified token colors.
- Tightened the public landing hero for desktop and mobile so the first viewport shows the core offer, CTA buttons, and a visible hint of the next section instead of trapping users in an oversized hero.
- Fixed the mobile web app auth shell so `/app` no longer overflows horizontally or clips login copy on narrow screens.
- Verified public-site language switching locally for English, Uzbek, and Portuguese on landing/legal pages after the redesign.
- Fixed yearly Premium and Platinum web crypto checkout so TON remains available even when only the monthly TON amount is configured; yearly TON amounts now fall back to the monthly amount scaled by the RUB plan price ratio.
- Changed the Telegram Premium flow to two steps: the user first chooses the tariff, then sees Stars, YooKassa, and configured crypto methods for that selected tariff.
- Reworked the web Premium block into aligned Free, Premium, and Platinum cards with dedicated plan artwork, monthly/yearly rows, and direct TON buttons for yearly tariffs.
- Added unique, understated header backgrounds for web app sections such as Practice, Progress, Awards, Tools, Premium, Limits, Mistakes, and related mobile views.
- Generated real branded PNG tariff images in `web/assets/plan-free.png`, `web/assets/plan-premium.png`, and `web/assets/plan-platinum.png`, with SVG fallbacks kept in the Premium cards.
- Regenerated the web app visual system through local ComfyUI/Juggernaut XL: plan artwork, section headers, and action/navigation icons now live in `web/assets/header-*.png` and `web/assets/icon-*.png` and are embedded by the Go web server.
- Changed web app image URLs to `/app/assets/...` and added matching Go routes so generated app visuals load correctly behind the production Caddy `/app/*` reverse proxy instead of falling through to the static landing site.
- Restored `web/index.html` to clean UTF-8 text without mojibake or a BOM so Russian and other localized web app copy renders correctly after deployment.
- Stopped embedding `web/index.html` into the Go binary; the web app shell is now served from the `web/index.html` file next to the running bot, while assets still keep embedded fallbacks.
- Compacted and branded the desktop right inspector so the Progress, Limits, and Account panels fit cleanly at a 1280x720 desktop viewport, with the Premium action visible at the top of the account panel.
- Added distinct branded treatments for main menu action cards, including per-feature accent colors, module marks, and stronger icon panels for desktop and mobile.
- Completed the local `.env` crypto payment block with the missing USDT-TON wallet placeholder and TronGrid settings.
- Hardened Telegram Stars payment handling by validating that invoice payloads match the payer Telegram ID and the expected Stars amount before pre-checkout approval or Premium activation.
- Refreshed the web app visual system for desktop and mobile with cleaner branded glass panels, adaptive tariff grids, embedded PNG plan artwork with inline SVG fallback, and more polished mobile navigation.

## 2026-05-27

- Replaced the legacy `/login` shell with the React V2 sign-in screen at `/login`, `/app/login`, and `/app/v2/login`, using the new `web-react/src/components/ui/sign-in.tsx` component.
- Preserved existing auth behavior on the new login screen: web login/register still use `/api/auth/login` and `/api/auth/register`, Cloudflare Turnstile renders in its own slot, and Telegram login opens the existing 6-digit OTP dialog through `/api/auth/telegram/start` and `/api/auth/telegram/status`.
- Renamed the V2 Phrasebook/Phrases surface to localized Notes/Р—Р°РјРµС‚РєРё copy across all 20 supported interface languages and forced client-side V2 labels to avoid older server copy overriding the new naming.
- Moved Lesson/Practice quick-save actions into the input/composer section below the send controls, and filtered out placeholder/question-mark candidates so Practice saves real phrases instead of `?`.
- Adjusted mobile V2 layout: quick header controls only appear on Today, Tools uses an in-composer send icon near the text input, result cards appear directly below trainer headings, and payment requisites modals sit above the bottom nav with scrollable mobile padding.
- Expanded Playwright desktop/mobile coverage for auth login, Telegram OTP, Turnstile-safe layout, mobile payment requisites, quick-save phrase filtering, trainer result placement, mobile Tools send controls, and Notes localization.

## 2026-05-19

- Rebuilt generated public-site phrase translations for the landing, Terms, and Privacy pages across all supported site languages, with safer deep-merge behavior and protected-token spacing.
- Fixed public landing fallback translations so non-Russian language views no longer show Russian in the marketing sections, plan cards, comparison table, FAQ, footer, and metadata.
- Fixed newly added Terms and Privacy pricing/payment text so Free, Premium, Platinum, Telegram Stars, voice/photo limits, and Privacy section 6 localize instead of staying mixed Russian/Korean/English.
- Added a public landing launch-offer story with crossed-out rounded base prices, visible 70% discounted prices, and first-month launch discount messaging.
- Added web Telegram Stars checkout: the web Premium screen now has Stars buttons that send a Telegram invoice to the linked bot chat.
- Changed locked web awards so they no longer open story modals; locked trophy artwork is now gray and blurred until the level is reached.
- Rebuilt the public landing around the current bot value proposition: AI tutor, Telegram/web sync, lessons, practice, voice, photo translation, mistake review, vocabulary, CEFR assessment, and multilingual UI.
- Added the Free/Premium/Platinum plan ladder across bot, web API, web app, landing, Terms, Privacy, README, and documentation.
- Updated paid prices: Premium 300 RUB/month or 150 Stars/month, 3000 RUB/year or 1500 Stars/year; Platinum 590 RUB/month or 300 Stars/month, 5900 RUB/year or 3000 Stars/year.
- Added USDT pricing equivalents: Premium about 4.15 USDT/month and 41.50 USDT/year; Platinum about 8.15 USDT/month and 81.50 USDT/year, using a 72 RUB/USDT configuration baseline.
- Added Platinum payment products, plan persistence, limits, invoices, crypto payment mapping, web DTO fields, and localized Premium menu copy.
- Added the web app button to the main Telegram keyboard and removed the duplicate web button from the translator/tools inline keyboard.
- Updated landing i18n hero copy for every supported site language and refreshed legal tariff badges.

- Added SQLite-backed runtime vocabulary tables in a separate `VOCABULARY_DATABASE_PATH` database with automatic JSON import and metadata-based resync.
- Added indexed vocabulary lookups for word id resolution and learned-word review options.
- Added indexed SQLite selection for the next learn-word prompt and answer options to avoid full dictionary scans.
- Changed learned vocabulary pages to show newest mastered words first.
- Added a direct mobile Tools tab and expanded the mobile web header with interface language, login, theme, and logout controls.
- Returned the mobile Tools tab to the main web app flow at `/app?view=tools`, with legacy `/app/tools` links redirecting to that view.
- Made the mobile leaderboard more compact with tighter cards, rows, scores, and language tabs.
- Made the shared mobile composer compact for practice, spelling, lesson answers, and mistake practice.
- Unified the web text-field styling across chat composer, auth/settings inputs, Telegram code entry, and translator text areas.
- Refined the mobile web header with compact account, Premium remaining, and language chips instead of a wide nickname row.
- Added a web Awards screen with 20 unique level trophies, polyglot ranks, trophy story modals, and a celebration animation when a new XP level is unlocked.
- Added direct TON crypto payments without a payment processor: SQLite invoices, Tonkeeper payment links, TON Center transaction checks, web API endpoints, web Premium buttons, Telegram Premium buttons, and automatic Premium activation.
- Added TonConsole/TonAPI support for direct TON checks, including `CRYPTO_TONAPI_KEY`, TonAPI fallback when TON Center rejects a key, and `/tonapi/webhook/<secret>` webhook activation.
- Added direct USDT payments for USDT-TON and USDT-TRC20, including USDT amount configuration, the official TON USDT Jetton master, TronGrid TRC20 checks, and exact unique TRC20 invoice amounts.
- Added a full web Premium direct TON payment flow with invoice details, copy actions, wallet opening, status checking, and browser-local expiry display.
- Added web payment confirmation celebrations for successful Premium payments, including direct TON auto-checking, a centered success card, fireworks, and a thank-you message.
- Changed Telegram TON invoice expiry text to use the user's selected timezone offset instead of raw UTC.
- Improved the web Tools page layout for mobile with a dedicated hero, icon tool picker, compact tool headers, and cleaner form spacing.
- Tightened the mobile Tools view scale with a shorter hero, compact tool tabs, and denser tool forms.
- Changed Awards trophies to render as larger designer trophy artwork in the grid instead of circle-like badges.
- Added crypto payment environment settings for the Caddy-hosted Go server: `CRYPTO_TON_WALLET`, `CRYPTO_TON_MONTH_AMOUNT`, `CRYPTO_TON_YEAR_AMOUNT`, `CRYPTO_TONCENTER_API_KEY`, `CRYPTO_TONCENTER_BASE_URL`, `CRYPTO_TONAPI_KEY`, `CRYPTO_TONAPI_BASE_URL`, `CRYPTO_TONAPI_WEBHOOK_KEY`, USDT TON/TRC20 settings, TronGrid settings, and `CRYPTO_PAYMENT_TTL_MINUTES`.
- Saved the current production Caddyfile as `deploy/caddy/Caddyfile.current` and added `deploy/caddy/Caddyfile.updated` for `poliglotai.ru`, `poliglotai.online`, their `www` hosts, API hosts, and crypto/web payment webhook paths.
- Improved mobile award trophy rendering by disabling SVG drop-shadow filters on phone layouts and stabilizing trophy dimensions to avoid intermittent mobile repaint gaps.
- Changed leaderboard language labels to English names.
- Kept streaming JSON vocabulary as a fallback when SQLite vocabulary rows are not available.
- Added SQLite vocabulary import coverage in tests.
- Reduced vocabulary memory pressure by removing the retained full-dictionary cache path from regular lookups.
- Optimized learned vocabulary pages so each page resolves word metadata with one streamed dictionary pass instead of repeated scans per row.
- Limited vocabulary-based level assessment generation to bounded streamed samples.
- Added a standalone referrals view to the web app desktop and mobile menus.
- Added a separate referrals menu item and `/referral` command to the Telegram bot.
- Improved mobile layout density for leaderboard rows, tool picker cards, and referral controls.
- Kept the Tools workspace inside the web app at `/app?view=tools` instead of routing the mobile tab through `/tools`.
- Changed the mobile Tools picker from a horizontal rail to full-width stacked cards so it fits the app header width on phones.
- Fixed mobile Tools icons so closed SVG shapes in the image and GPT agent tools cannot render as gray filled squares in mobile WebViews.
- Added a mobile web spelling give-up flow after three wrong attempts: `РќРµ Р·РЅР°СЋ` reveals the correct answer without XP or spelling progress, and keeps the word in the spelling pool.
- Improved mobile Learn Words scrolling so a correct answer brings the `Next word` action into view on phones.
- Tightened the mobile leaderboard and bottom navigation to match the Tools view width and keep overflow-menu screens visually consistent.
- Changed public-site app links to relative `/app` links so `poliglotai.online` opens the app on the online mirror instead of jumping to `poliglotai.ru`.
- Fixed mobile leaderboard container alignment so the rating block lines up with the app header and menu instead of drifting right.
- Fixed mobile leaderboard rendering on browsers with a wide layout viewport by forcing the mobile app shell from touch/screen detection instead of leaving the desktop sidebar visible.
- Changed public-site visible app URL labels to follow the current host: `poliglotai.ru/app` on the Russian domain and `poliglotai.online/app` on the online mirror.
- Fixed mojibake text in the public landing hero app mockup and replaced the broken Premium star glyph with a FontAwesome icon.
- Removed duplicate Tools from the main web action grid.
- Removed referral controls from Settings in the web app and Telegram settings menu.
- Changed referral share text to address the invited person in the selected interface language and explain the 7-day Premium bonus.
- Stopped auto-focusing the web composer when opening lesson, practice, spelling, or mistake-practice screens.
- Updated V2 roleplay so scenario selection enters a dedicated dialogue surface instead of appending output below the scenario grid.
- Made roleplay scenario tiles render as an even two-row desktop grid and added scrollable V2 context surfaces for roleplay and offline decks.
- Replaced the Dashboard leaderboard/weak-topic blocks with learning activity metrics: learned words, practiced phrases, lessons, review rounds, voice attempts, and active mistakes.
- Cleaned Pronunciation weak-word rendering so learner-facing cards no longer expose raw `low_confidence`/system issue keys.
- Tightened exact-repeat pronunciation scoring caps around target similarity, missing words, substitutions, estimated confidence, and non-praising feedback.
- Restricted learn-word and review answer distractors to the active CEFR band (`A1/A2`, `B1/B2`, `C1/C2`) instead of backfilling from the whole language when a band is sparse.
- Rebuilt vocabulary JSON files and added deterministic native-language context values for dictionary prompts while keeping target-language answer options as single words.
- Added roleplay question-audio fallback so mobile roleplay dialogue can play the AI prompt even when the model answers without a final question mark.
- Fixed mobile Mistakes and Notes scrolling so dictionaries/cards are not clipped or hidden behind the bottom navigation.
- Added missing mobile regression localization keys for mistake practice, saved notes, and roleplay placeholders across all 35 interface languages.
- Matched the public landing hero light theme to the dark shader/glass hero so the first viewport keeps the same contrast and visual identity.
- Added/updated Go and Playwright regressions for CEFR-banded vocabulary selection, native dictionary context, roleplay audio, mobile mistakes scrolling, mobile Notes bottom-nav clearance, and all-language labels.
- Added a permanent secured deploy-upload service on the Ubuntu server for large artifacts when SSH/SCP streaming is unstable.
- Documented the reusable deploy-upload prompt in `deploy/SERVER_UPLOAD_PROMPT.md` and kept the deploy Caddy template route in sync.
- Replaced the mobile `More` bottom-menu destination with one horizontal scrollable navigation rail containing all menu items.
- Added mobile nav edit-mode regressions for long-press drag reordering, visible drop gaps, no block descriptions, and absence of the old More sheet.
- Expanded the public landing light-theme contrast overrides so the first shader/glass viewport keeps readable white hero/card text while other light sections keep dark readable headings.
- Added persistent SQLite cache for AI-generated vocabulary hints/examples so missing native-language context is generated once and reused instead of repeatedly calling the model.
- Simplified Listening/Audition so the duplicate upper phrase card is hidden and the spoken-model audio lives inside the active repeat/input panel.
- Changed Mistakes into a two-step flow: the dictionary list is the first screen, and selecting an error opens a dedicated practice screen with a return action.
- Added spelling-result output of the correct target-language word for both correct answers and give-up answers.
- Hardened AI vocabulary hint prompts and sanitization so native-language clues do not repeat the hidden target word.
- Improved mobile bottom-rail drag preview, mobile leaderboard bottom clearance, roleplay role-card spacing, and referral localization fallback keys.
- Added/updated regressions for listening layout, spelling correct-answer display, mistake list/practice navigation, mobile nav reflow, leaderboard clearance, and 35-language referral labels.
- Fixed Privacy runtime rendering so the bottom Telegram bot contact card keeps `@poliglot_ai_bot` on `privacy.html?lang=ru` for both production domains.
- Centered the login/auth form over the animated background, kept the Cloudflare/login/register/Telegram OTP flows, and added layout regressions that reject the old right desktop panel.
- Hardened the mobile bottom rail reorder flow: long-press edit survives pointer-capture failures, drag auto-scrolls near rail edges, neighboring tiles reflow, and the saved order persists after reload.
- Tightened vocabulary prompt context so the primary native-language value is the prompt and up to two alternate meanings are kept as context without leaking the target answer.
- Re-ran the Go vocabulary regressions covering Russian prompts, OpenRouter JSON/cache sanitization, CEFR-band distractors, randomization, and same-answer leak prevention.
- Downloaded Kaikki/Wiktextract dumps for the 15 former UI-only languages and rebuilt the vocabulary set so all 35 interface languages are now selectable learning dictionaries.
- Added Arabic, Bengali, Czech, Greek, Hindi, Hungarian, Indonesian, Dutch, Swedish, Tamil, Telugu, Thai, Tagalog, Turkish, and Vietnamese vocabulary JSON files plus script-safety regression coverage.
- Updated the permanent deploy package so `poliglot-app-web.tgz` carries all `vocabulary_words*.json` files and the server deploy script backs up and installs those dictionaries into `/opt/aibot`.
- Hardened the permanent `poliglot-deploy-upload.service` with same-file chunk headers for unstable large PUT uploads while keeping the same token, allowlist, and Caddy route.
- Finished the 35-language V2 localization hardening pass: web copy now prefers clean client-side locale data over stale server fallback strings, repairs common mojibake before rendering, and rejects generic `Section: key`/`Mб»Ґc: key` labels.
- Added clean Georgian and Vietnamese overrides for V2 navigation, Today, Roleplay scenarios, and scenario descriptions, including a fallback path that avoids old mojibake scenario labels.
- Added a Playwright desktop/mobile regression that opens the main V2 screens across every non-English interface language and fails on mojibake or exposed English/Russian technical fallback labels.
- Added persistent desktop function-ribbon reordering with visible drag/drop targets, edge auto-scroll, saved order per account, and Playwright persistence coverage.
- Added critical 35-language menu label coverage for Spelling, Progress, and Mistakes, including their nav descriptions and view titles, so English fallback leaks fail Playwright.
- Filled missing non-English V2 menu overrides for Spelling/Progress/Mistakes in Arabic, Bengali, Czech, Greek, Hindi, Hungarian, Indonesian, Dutch, Swedish, Tamil, Telugu, Thai, Tagalog, Turkish, Uzbek, and Vietnamese.
- Improved mobile bottom-rail dragging so movement starts reordering immediately after touch movement, disables mobile callout/context-menu interruption, and keeps HyperOS-style drop persistence after release.
- Added Offline Notes pagination coverage so the Notes group shows 10 cards per page and moves remaining notes to the next page with the same controls as Vocabulary.
- Fixed mobile Roleplay scenario cards so long localized text expands the card instead of spilling past the border, and the last card scrolls above the fixed bottom menu.
- Hardened V2 Playwright setup to clear service worker/cache state before tests so future visual and localization regressions cannot pass or fail because of stale built assets.
- Simplified mobile bottom-rail editing: long-press now enters a stable edit mode after 2 seconds, reordering saves only on release, the floating drag animation was removed, and a visible Done button exits edit mode.
- Restored horizontal swiping for the mobile bottom rail by allowing pan gestures outside edit mode and limiting gesture blocking to active reorder mode.
# 2026-05-30

- Improved mobile web menu reordering: touch hold enters edit mode faster, drag targeting follows the actual finger position, and the Playwright regression now checks editable mobile rail reorder behavior.
