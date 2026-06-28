# User UI Wishes

## 2026-06-27 Social Video Outro

- Every reusable TikTok/Reels/Shorts outro should be a 2-second vertical brand insert with only the Poliglot AI logo and exact URL `poliglotAI.online` visible.
- Use the local FLUX.2 + Wan 2.1 + RIFE ComfyUI chain already configured in `comfyui_workflows`; avoid downloading new large models unless the installed set cannot produce the requested quality.
- Keep the URL exact and readable. If the video model distorts text, use the generated motion as background and overlay the exact logo/URL from the reference card in the editor.
- Prefer the more kinetic first generated background for the final outro style, but do not trust model-rendered typography. The URL should appear after the logo with an After Effects-style float/reveal, light sweep, and clean deterministic overlay.
- Outro URL typography should be clean, solid, and instantly readable on phones after social-platform compression. Avoid outline-like or overly decorative text; use the bundled Manrope-style solid brand lockup with `AI` as the mint accent.
- Do not add a second deterministic logo layer on top of an already animated generated logo. The post-render should keep the source logo motion and only replace/overlay the exact URL. Keep the lower text area clean: no decorative stripes, scanlines, or underline below the URL.

This file is the recurring UI checklist for Poliglot AI. Check it before changing the web app, public site, generated assets, or mobile layout.

## 2026-05-30 Mobile Web Menu Reorder

- Auth and registration system messages must be localized through the same 35-language V2 layer as the rest of the app. Backend auth errors should carry stable codes, and the login UI must never show generic fallback labels such as `Раздел`, `Section`, or `Mục` instead of the real error.
- Web logins may contain a dot inside the login (`friend.name`), but may not start with a dot, hyphen, or underscore. The error for an invalid start character must clearly explain the login rule in the selected interface language.
- Installed PWA/mobile browser app must always converge to the current web interface after deploy. The `/app/v2` shell, manifest, and service worker need no-store headers, and the service worker should try network first for the app shell/assets while keeping cache only as offline fallback.
- Dashboard/Progress statistic labels must stay specific, not collapse to the generic `Statistics` word; Premium/Platinum plan names must not repeat `Premium Premium`; the Phrasebook second input must say `Note or translation`. These labels need static localization and regression coverage for all 35 interface languages.
- Desktop function-ribbon and mobile bottom-menu ordering must be stored per backend user account, not only in browser `localStorage`, so the layout follows the user across reloads/devices.
- Wrong answers in Spelling must also become Mistakes dictionary entries; Lessons, Practice, and Spelling should all feed the same repair flow.
- Web bug reports sent to Telegram must include every attached screenshot as a real Telegram photo when the user uploads screenshots, not just a path or link in the text.
- Offline deck downloads must match the selected filter group exactly (`All`, `Words`, `Notes`, or `Mistakes`), TXT downloads must be readable UTF-8 on Android viewers, and every offline card must have an in-card delete button in the bottom-right corner.
- Level assessment must be available as a real exam-style test for every supported learning language.
- Level assessment must not fall back to bare vocabulary guessing for any of the 35 learning languages. Every language needs the same exam-style flow with sentence, gap-fill, grammar, meaning, naturalness, and formal-register tasks.
- The 15 newly added learning languages must use the same English level-assessment question set as the existing English test, but the visible question prompt/instruction must be localized to the user's interface language. Regression coverage should check both the A1/A2/B1/B2/C1/C2 outcomes and the interface-language prompt.
- Level-assessment task wording must be specific in every interface language: fill-gap, translation, meaning, naturalness, grammar/form, and "how to say" questions should stay recognizable, not become the same generic choose-answer label.
- The Mistakes screen owns pagination on the client and must receive the full saved mistake list from the API, newest first, so new Lesson/Practice/Spelling mistakes appear immediately on page 1.
- User visit calendars must be persisted per user in the backend database for a rolling 365 days. Browser storage can mirror the state for responsiveness, but it must not be the only source of truth.
- Mobile bottom-menu edge auto-scroll during tile dragging must be slow and controllable; dragging near the edge should nudge the rail, not fling it across the menu.
- Lesson output must not print raw `Spoken model` / `Озвученный образец` text in place of the sample. Show an actual audio control/button for the spoken sample whenever the backend returns audio text.
- Mistakes made in web Lessons and Practice must always be saved into the Mistakes dictionary. If the normal answer model forgets the `---MISTAKES---` JSON block, backend must run a strict fallback extraction instead of silently dropping the mistake.
- App startup must never show a blank mobile screen. Before React hydrates and while the session is loading, show the Poliglot AI logo/title/spinner over the animated dark glass background.
- Mobile bottom-menu reorder must visibly move the selected tile during drag, like a phone launcher. Keep the final `localStorage` write on pointer release so slow phones are protected, but do not leave the visible rail frozen while the user is dragging.
- Regression coverage for mobile menu reorder must check both phases: immediate visual movement before pointer release and saved order after reload.

## 2026-05-29 Web And Telegram 35-Language Localization

- Mobile bottom-menu reorder must stay lightweight: keep drag coordinates in refs, move only the ghost preview through `requestAnimationFrame`, avoid `localStorage` writes until pointer release, keep visible rail movement limited to small order-array updates, and keep `will-change: transform` limited to the ghost preview.
- During mobile menu edit/reorder mode, avoid expensive visual effects on menu buttons: no floating icon animation, no heavy box-shadow, no blur/backdrop filter, and no scale transforms on the actual buttons.
- Web bug reports must notify Telegram chat `185156683` (`AsaselD`) in addition to the local report log.
- Successful premium payments from Telegram Stars, YooKassa, and direct crypto must also be duplicated to Telegram chat `185156683` (`AsaselD`).
- The web app and Telegram bot must not show English UI fallback for any supported non-English interface language. Static UI copy should be owned in the codebase; do not use OpenRouter for interface labels.
- In Learn Words and Review, the source prompt word must be on the currently selected interface language. If the local dictionary lacks that prompt, call OpenRouter only for this vocabulary translation, save it into SQLite (`vocabulary_ai_translations` and `vocabulary_translations`), and reuse it; Russian or English fallback is only acceptable as a temporary model/configuration failure path.
- The source prompt must never be the same word as the correct answer. Treat dictionary rows like `creating -> creating` as missing translations and repair them through the vocabulary translation cache path.
- AI-repaired vocabulary translations must be durable: save them into the runtime SQLite cache and also back into the matching `vocabulary_words*.json` file so future SQLite rebuilds do not lose them.
- CEFR level controls must expose the full `A1/A2/B1/B2/C1/C2` range everywhere, including first-run onboarding and Settings. Learn Words should use the exact selected level first and only fallback inside the paired band.
- Mobile horizontal menu drag must feel physical: the dragged tile stays visible as a floating preview, neighboring tiles open a smooth before/after gap, edge auto-scroll works, and the tile lands in the intended spot with the order saved after reload.
- The 35-language web regression must check required V2 labels on desktop and mobile, including auth, navigation, view titles/subtitles/actions, settings, notes, and mobile/menu labels.
- Telegram localization must cover the main menu, learning actions, word learning/review prompts, settings, and tool prompts for all 35 interface languages. New interface languages must not be added as aliases to English.
- Tagalog and other Latin-script languages still need natural local wording even when English words are commonly understood; avoid leaving exact English strings such as `Settings`, `Notes`, or `Review` in required UI labels.
- After any interface-language expansion, run `go test ./...`, `npm --prefix web-react run build`, and `npm --prefix web-react run e2e` before deployment.

## 2026-05-28 V2-Only Shell / Vocabulary Fallback

- The old HTML web app is removed. Do not recreate `web/index.html`; `/app/v1` should redirect to `/app/v2/`, and public links should use `/app/v2/`.
- Learn Words, Review, Spelling, Vocabulary cards, and Telegram word study should stay responsive, but exact interface-language source prompts are now stricter than the old instant fallback rule: check SQLite first, translate/persist missing prompts through the guarded vocabulary OpenRouter path when configured, and only show the closest dictionary prompt if that path is unavailable or fails.
- Missing word translations should be saved once into local SQLite (`vocabulary_ai_translations` and `vocabulary_translations`) so later requests do not ask the model again.
- The OpenRouter word-translation prompt must behave like a professional dictionary lookup: strict JSON, common meaning, part of speech/CEFR/topic aware, up to three short native-language equivalents, no target word, no transliteration, no examples.
- The login page should reuse the same premium generative animation language as the landing hero on desktop and mobile while preserving Cloudflare, password login, registration, and Telegram OTP.
- In the public landing light theme, section label pills must remain readable; do not ship pale yellow text on a pale background.
- Missing free-license learning dictionaries for languages beyond the current dictionary-backed learning set remain a separate data task. Do not expose a language as fully dictionary-backed until source files are downloaded, license-noted, rebuilt, tested, and documented.

## 2026-05-27 Landing / Final Regression Closure

- The public landing must use the generated Poliglot logo asset in the header, not a text-only or stale V2 mark.
- The landing hero should feel premium and unified: one combined shader/product animation, polished CTA motion, concise product sections, restored reviews, FAQ, and visible pricing proof without information overload.
- The landing should follow the Busuu-style conversion structure from the provided reference HTML, but with Poliglot AI logic: clear hero, `Я хочу изучать` language chooser, course cards, proof stats, feature blocks, reviews, FAQ, and tariffs.
- Keep the old monthly prices crossed out in pricing when the visible monthly price is already the 70% discount price.
- Public landing, Privacy, Terms, web app, and Telegram interface selectors must support all 35 requested interface languages. Legal pages need the same contact/header treatment as Terms and dark-theme-readable language selectors.
- The landing must not expose a public `/app/v2` nav link; use the theme toggle where that old V2 button was.
- Mobile Awards, Mistakes, Roleplay, Tools, composer, Phrasebook/Notes, and menu editing regressions need Playwright coverage before deployment, not manual-only QA.
- Learn Words and Review must use a wide SQLite-backed CEFR-band option pool in web and Telegram. Do not accept repeated fixed wrong-answer sets as "random enough".
- Interface-language expansion must not pretend that unsupported learning dictionaries exist. Keep the learning/vocabulary list on the dictionary-backed languages until new word assets are generated.
- SQLite vocabulary randomization must behave like the old in-memory pool: pick from a broad cached array of eligible IDs, never from the first `position=0` fallback after a failed random start.
- Production vocabulary databases should stay indexed before deploys that touch words/review performance.
- Every reported mobile layout issue should be protected by a regression test when practical: scrollability, visible input, compact headers, dropdown z-index, modal viewport placement, and browser Back behavior.

## 2026-05-27 Login / Auth Notes

- The public login routes `/login`, `/app/login`, and `/app/v2/login` should use the React V2 standalone auth screen, not the old legacy login layout.
- Login must keep the current backend behavior: login/password auth, registration, Cloudflare Turnstile when enabled, and Telegram six-digit OTP through the existing popup flow.
- Standalone auth copy must be localized through `web-react/src/lib/i18n.ts` for all 35 supported interface languages. Do not leave English fallback text visible in Russian or other non-English locales.
- The login page must not show the V2 app header, ribbon, bottom navigation, or in-app menu while the user is unauthenticated.
- Login artwork should be generated through the existing local ComfyUI FLUX.2 workflow, using the already installed model chain (`flux2_dev_fp8mixed`, `mistral_3_small_flux2_bf16`, `flux2-vae`). Do not re-download models for this asset pass.
- Telegram login must be a compact row with a small logo, not a half-screen image block.
- Registration and Telegram account linking must show an explicit personal-data policy consent checkbox and link before sending auth/linking requests.
- Account recovery should open `https://t.me/AsaselD`.
- New Telegram users must see a first onboarding post that asks them to accept the personal-data policy, with Continue and Policy buttons, before the normal onboarding flow.
- The public Privacy page and `privacy-policy-i18n.js` must stay synchronized for all 35 supported interface languages on both `poliglotai.ru` and `poliglotai.online`.
- Learn Words and Review answer options must come from real dictionary distractors in the same CEFR band, not from a fixed repeated pool. This backend path is shared by web desktop, web mobile, and Telegram.
- Mobile award details must open in the currently visible viewport through a fixed/portal modal. Tapping a reward after scrolling must not jump the modal to the top of the page.
- In mobile Lesson/Practice, the "save to Notes/Phrasebook" action belongs inside the composer immediately after voice/photo controls and before/near the send controls, never as a separate low block that the user cannot see.
- Mobile Roleplay should not show a bulky `ROLEPLAY DIALOGUE`/scenario heading after a scenario is selected; keep only the useful localized action to choose another scenario.
- Mobile Mistakes must allow scrolling to the bottom of the error dictionary and drill list.

## 2026-05-26 Latest V2 Mobile / Word Regression Notes

- Mobile onboarding/start screen must fit inside the viewport and scroll internally so goal, level, language, format, and action buttons are always reachable on small phones.
- Offline, Phrases/Phrasebook, Pronunciation, and Roleplay must use the current FLUX.2 transparent menu assets on both desktop web and mobile web. Always bump the menu asset cache-buster when these PNGs change.
- The desktop function-ribbon right arrow must stay inside the app frame with visible padding from the right border.
- Mobile Lesson/Practice quick-save to Phrasebook must be compact and appear after the input/composer block, not inside the output area where it can cover task text.
- Mobile Roleplay sessions should hide the scenario brief after a scenario is chosen and show only a localized "choose another scenario" action above the dialogue.
- Pronunciation weak words, sounds, and progress history must persist per user across reloads instead of depending only on the current in-memory message list.
- Mobile More edit mode must support reordering tiles inside More as well as pin/unpin between More and the bottom navigation.
- Mobile Tools must use a selector-first flow: after a tool is selected, hide the large tool list and show the chosen tool window with a "choose tool" action to switch back.
- Learn Words answer options must not reuse the same wrong-answer set for each new word. The shared backend option generator is used by web desktop, web mobile, and Telegram, so regressions here need backend tests as well as Playwright coverage.

## 2026-05-26 V2 Regression Test Requirements

- Before every V2 deploy, run the Playwright desktop/mobile suite for the reported failure classes: 35-language menu localization, confirmed-only payments, referral invitee pagination, Telegram `Send code` visibility, Offline 10-card pagination with full TXT/JSON export, compact pronunciation map, desktop ribbon arrow bounds, mobile header controls, mobile More image tiles, editable mobile pinned menu, mobile Lesson/Practice/Roleplay/Tools input visibility, mistake dictionary scrolling, Phrasebook save persistence, and mojibake guards.
- Lesson start output must show the actual task text when the API returns `prompt`, `question`, or `task`; do not leave the user with only a generic "lesson created" message.
- Practice and Roleplay result cards must render useful `correction` and `explanation` fields, not only audio/model text.
- Mobile edit-mode animation may make icons feel floating, but the actual plus/X tap target must stay stable and clickable.
- V2 is now the only active web app; `/app/v1` is a compatibility redirect to `/app/v2/`, not a legacy shell.

## 2026-05-26 V2 Payments / Referrals / Mobile Editing Notes

- Premium payment history must show only confirmed payments. Do not add `created`, `pending`, or unpaid invoices to visible history.
- Offline, Phrases/Phrasebook, Pronunciation, and Roleplay menu artwork should be regenerated through the existing local FLUX.2 ComfyUI workflow, not a one-off checkpoint path. The current chain uses `flux2_dev_fp8mixed`, `mistral_3_small_flux2_bf16`, `flux2-vae`, chroma-key removal, and 512x512 transparent validation.
- Weekly plan cards should contain enough task detail to feel useful, but not turn into a noisy roadmap.
- Today should not start with a separate grid of quick action buttons when Daily quests already covers those actions.
- Referrals must show which users joined through the invite link, whether each reached level 3, and how much was earned from each user. Paginate invitees after 10 rows.
- Settings must show Telegram `Send code` only while Telegram is not linked.
- Mobile More must open visual tiles with the same menu images as desktop/V1-style cards.
- Mobile header quick controls must include bug report, language, theme, and icon-only logout, placed below the app frame line instead of overlapping the border.
- Mobile menus must be user-customizable: long-press to enter edit mode, X to move pinned items into More, plus/tile tap to pin items back, and Done/Cancel to leave edit mode.
- Mobile input/composer areas must remain visible above the bottom nav in Lesson, Practice, Roleplay, Listening/Pronunciation, and Tools.
- Tools on mobile must keep all tool-mode buttons selectable; no shifted/offscreen selector.
- Saved Offline cards should show 10 per page with back/next controls, while TXT/JSON export still includes the full deck.
- The daily streak label such as `1 день подряд` should be green.

## 2026-05-26 V2 Localization / Mobile Menu / Scroll Notes

- V2 menu labels must be localized for all 35 supported interface languages. This includes Roleplay, Phrases/Phrasebook, Offline, Dashboard, Invite/Referrals, V2 learning lab, Daily quests, Weekly plan, Daily route, Roleplay buttons, Offline decks, and other system section names.
- Russian labels for the requested core tabs should stay short and clean: Roleplay = `Ролевая`, Offline = `Оффлайн`, Dashboard = `Статистика`, Invite/Referrals = `Приглашения`.
- New localized strings must be normal UTF-8. Do not ship mojibake, replacement characters, or broken cp1251 text in React fallbacks, tests, docs, or generated assets.
- Pronunciation and Phrasebook menu artwork should be visually centered and scaled larger inside their buttons, both desktop and mobile.
- The Pronunciation map should be compact. Avoid repeating the same long advice text on every token; show short issue labels and one useful learner-facing hint.
- In Today, the 1-day streak panel should align cleanly to the right of the calendar on desktop and stack without overflow on mobile.
- Remove roadmap/recommendation filler such as the old "Next features / Что можно добавить дальше" block from Today when those features already have real screens.
- Offline decks should use the available width and not leave most of the screen empty because only one side column is active.
- Mobile V2 should not show the desktop top menu/ribbon. Keep the interface light: primary destinations in the bottom bar and the rest behind a bottom-right More tab.
- On mobile, logout belongs in the top-right corner as an icon-only control, with the small theme toggle immediately to its left.
- Dense V2 blocks must scroll on mobile. Do not trap content in nested fixed-height panels that prevent the user from reaching lower cards or actions.

## 2026-05-26 V2 Phrasebook / Habit / Reports Notes

- Roleplay, Offline, Pronunciation, and Phrasebook need their own generated menu assets in the same Poliglot V2 visual system, with separate light and dark variants.
- Lessons and practice should let the user save useful phrases into a personal Phrasebook. Saved phrases must stay easy to review, play back, remove, and reuse in Offline decks.
- The Today screen should include a month habit calendar: login days are visible, completed daily-goal days are filled green, missed days stay understated, and the current streak is easy to read.
- The weekly plan should be realistic, not overloaded: a few words, short practice, one listening/repeat task, and one error or phrasebook action per day.
- Claiming the completed daily goal should grant a visible bonus XP animation and must not allow duplicate claims for the same day.
- V2 needs a top-bar bug report action. The report dialog should accept a clear problem description and an optional screenshot, then save a report file for later mail-server integration.
- Offline decks must let users group cards by content type: words, favorite phrases, and mistakes. Word cards should look like learning cards, not raw JSON rows.
- Every new V2 feature needs separate desktop and mobile behavior documented before deployment, plus Playwright coverage for the main path.

## 2026-05-24 V2 Encoding / Tool Output Guard

- Never paste broken cp1251/UTF-8 mojibake into React fallbacks, prompt templates, changelog, docs, or static copy. All visible text must be normal UTF-8.
- After any text/copy/localization change, run `node tools/check_encoding_artifacts.mjs` before build/deploy.
- Each V2 tool must render only its own messages. New user/assistant messages need a stable `meta` surface id and each view must filter by that id before calling `ChatPanel`.
- Roleplay is not generic Practice. It must use the dedicated Roleplay dialog structure, wait for the learner's first answer before continuing the next question, and keep setup/correction labels user-facing.
- Settings must keep the global learning focus field so the user can change onboarding topic preferences without code changes.
- New lessons must never repeat the last 10 generated lesson tasks. If no global focus is set, rotate level-appropriate CEFR topics automatically.
- Practice must keep the last 5 learner messages as useful context, but must not dump unrelated tool history into the current answer.
- V1 and V2 awards must use one canonical set of polyglot rank names and history stories. Russian UI must show Russian award stories, not English fallback or placeholder keys.

## 2026-05-23 V2 Feature Roadmap Slice

- New feature work must still follow the point order from the 14-stage V2 plan; do not switch back to an unstructured fast pass unless the user explicitly asks again.
- Add a dedicated deployable prompt file next to the server binary and env files. It must contain all application LLM prompts, with a clear prompt header, feature/tool name, source function, and notes so the user can edit individual prompts without rebuilding the binary.
- Home/"Today" must show a practical weekly plan, daily quests, progress snapshot, and the next best action from mistakes, vocabulary, Listening, and roleplay.
- AI roleplay scenarios must be visible as real selectable blocks: restaurant, work, travel, exam, and small talk. Starting a scenario should enter the normal Practice conversation flow.
- Pronunciation needs a dashboard, not only a result block: current score, weak words/sounds heatmap, score history, retry, and Listening entry point.
- Clicking Pronunciation/Pronounce must open the pronunciation dashboard itself. It must not auto-start Listening/Audition; Listening should only start from an explicit retry/listening button inside the dashboard.
- Offline/PWA mini-decks must save useful vocabulary/mistake cards in the browser and keep a deployable service worker/manifest for `/app/v2`.
- Teacher/admin dashboard can start as a current-account dashboard, but it should expose progress, activity, Premium/payment status, referrals, leaderboard, and problem topics.

## 2026-05-23 V2 Mobile Navigation / Live QA Notes

- Mobile V2 must have its own touch navigation layer, not just a scaled desktop ribbon. The current priority row is Home, Lesson, Practice, Words, Mistakes, Premium, Settings.
- Function-ribbon neighbors may move to make space, but their icon/caption size, opacity, and scale must stay stable.
- Status banners must remain in the top header area above the function ribbon and auto-dismiss after 3 seconds.
- Do not mark live authorized QA as done while the Codex-controlled browser is still stuck on `/login` captcha.
- After every V2 UI slice, rebuild the React app, regenerate the Nx project graph, and deploy only after the Go and encoding checks pass.

## 2026-05-23 Final V2 Menu / Voice / Payment Notes

- Neighboring function-menu chips must only make room for the hovered/active chip; they must not shrink, fade, or visually "jerk" while the cursor moves.
- The recording button must keep the current Poliglot color/style, but its active state must animate like the supplied example: stop-square, expanding waveform bars, and a timer.
- Every changed desktop control needs a mobile-specific treatment with the same function and visual intent, especially ribbon controls, audio controls, upload chips, dialogs, and payment blocks.
- Settings profile/security cards should end at their content, not stretch to the bottom of the row. Telegram `Send code` opens a centered mini OTP dialog instead of expanding inline.
- Telegram two-factor `Send code` in Settings must send the code through the same Telegram delivery flow as Version 1.
- The two-factor confirmation UI must be a separate centered popup window, visually analogous to the payment dialog, not an embedded panel inside the current Settings card.
- Premium payment success and payment-not-found states must open in centered modal dialogs. Success must show subscription plan, month/year period, amount when available, and the updated current expiration date.
- Payment instructions should not say "1 в 1". They should tell the user to transfer exactly the amount/network/comment shown in requisites and warn that a mismatched payment may not pass.
- Mistake dictionary playback must use the same large waveform visual as other voiced examples.
- Header system messages stay in the header/topbar area, not down near the function ribbon.
- Dark and light backgrounds should be animated and visible through slightly transparent input/output panels.
- V1 award rank names and polyglot stories are canonical. V2 must not invent or show placeholder keys.

## 2026-05-23 V2 Follow-Up Polish Notes

- Menu hover motion must start only when hovering a real menu item, not empty ribbon space.
- Voice/audio file upload must stay a compact icon button like the microphone button; after upload, show a small closable filename chip to the right.
- Image upload must match microphone color, size, hover, and compact attachment behavior.
- Header ribbon arrows must be inset far enough from screen edges that they do not clip.
- Listening "Next phrase" must sit at the bottom-right under the phrase text block.
- Learn Words should be named as the learning block, not "choice trainer"; remove "task ready", duplicate context/example text, and result context when it is already shown elsewhere.
- Entering Learn Words, Review, and Spelling should force a fresh active task load instead of showing stale or empty "no active question" states.
- Spelling should not play/show initial source-word audio; show only the example audio immediately after the example.
- Level test result score must show the real maximum once in a polished score block; no "next word" action belongs there.
- Awards must show readable rank/title and polyglot story text; placeholder keys like `award_rank_06` and `award_story_06` must never be visible.
- Mistake/dictionary audio controls must use the waveform audio style and must not repeat comments already visible in the side panel.
- Chat/output panes must auto-scroll to the newest message when the learner adds text or new generated content arrives.
- Payment requisites must include exact-amount/network instructions, correctly formatted TON/USDT amounts, expiry time, and a manual "check payment" action.
- Settings layout should place Activation under Telegram on the right so the page reads as balanced blocks.
- Header identity should show the user login only once under "Poliglot AI"; do not duplicate it near the trophy.
- Language dropdowns must render above the header/context layers and never be hidden under adjacent panels.
- Menu and context animations should be slower, smoother, and visibly drop downward on view selection.

## 2026-05-23 V2 QA Notes

- Work strictly through the user's V2 remarks sequentially and re-check the whole V2 concept after each meaningful block.
- System status/toast messages such as lesson completion must appear at the top center under the header and disappear automatically after 3 seconds.
- Function-menu hover must smoothly expand the hovered item and push neighboring items apart.
- All language selectors must use the animated project dropdown style, including top bar, leaderboard, settings, and tool language pickers; dark theme dropdowns must not use unreadable white fills.
- Award/trophy modals must show the trophy/rank name and the polyglot history story from the old version, never generic "Level 5 / Rank 5" fallback copy.
- Mistake practice must not have a Start button, Clear must ask for confirmation, and each correct answer must have a backend TTS playback control.
- Level test answer choices should render as a polished answer-selection window, not as plain browser-looking buttons or a text-answer form.
- Dictionary cards must behave like a dictionary: do not show review/spelling counters there; show a level-appropriate example sentence under each word.
- Word examples should be generated by the configured AI model and match the word CEFR level where the backend can provide that context.
- Spelling must keep a Next word action and include playback for the example sentence.
- Learn Words should reveal word audio only after the correct answer, and the result should include a level-matched example sentence.
- Listening must include a Next phrase action.
- Voice/audio message waveform blocks must stay visually stable during navigation and typing; waveform motion should only run while audio is playing.
- All icon buttons, including record, attach image, attach audio, and file controls, need smooth hover/tap animation consistent with the rest of V2.
- File upload controls should follow the supplied shadcn/Tailwind upload examples closely while fitting the existing compact chat controls.
- Any loading process should use the supplied `Spinner` pattern instead of ad hoc loaders.
- Telegram invite/share messages in Referrals must reuse the old version's localized invite text 1:1 where available.
- Settings must include Change password with the supplied password-input styling and the existing backend account logic.
- Review and Learn Words must auto-create an active word question when entering the block; users should not see "No active word question" on first entry.
- Referrals must clearly explain program terms and show invitation actions/copy like the old version.
- The header identity must show the current user-created login, not the generic AI Tutor label or Telegram username when a web login exists.
- Context windows under selected menu items must open with a smooth, modern downward animation.
- Lesson results must include audio playback for the corrected phrase.
- Practice must include audio for the AI question after "Your turn", not only the example/model phrase.
- Translator must have a visible swap-languages arrow button between source and target selectors.
- Chat output must be independently scrollable; the message window must not trap content at the bottom when history is longer than the visible area.
- Leaderboards must expose Global top as an explicit control, not only per-language entries.
- Top function-menu buttons must visually expand from their center in both directions, not open left-to-right.
- Vocabulary should be compact and preferably two columns on desktop; every word card needs its own listen button under/near the word.
- All voice recording controls must use the animated `VoiceInput` microphone/wave/timer style. Playback/listen controls must use a speaker/play waveform style and must not show as microphone buttons.
- Voice/audio samples in V2 must use the backend GPT-4o Mini TTS/OpenRouter path, matching V1 behavior; do not use built-in browser speech synthesis for learner-facing audio.
- Keep only the requested Paper Shaders background treatment. Do not mix extra radial/orb/background overlays that hide the requested background.
- Header XP/trophy opens Awards; header brand/logo opens Settings.
- Theme toggle, payment method buttons, left/right arrows, logout, OTP, and save/approve/reject controls should follow the exact animation families supplied by the user unless a backend integration requires a small adaptation.
- Payment methods should highlight on hover/selection, keep a pending/loading state while payment is unresolved, and display full requisites in the same modal below the selected method: amount, comment/memo, wallet/address, network, payment ID, status, expiration, transaction hash, and payment link when available.
- Level test should not allow free-text answers. Remove "I do not know" from answer options and map Skip to the same backend answer value as "I do not know".
- Header must include a brightness slider for the whole interface.
- Telegram linking code flow should show a 6-digit OTP modal after Send code.
- That OTP/two-factor modal must be separate from Settings content, centered on screen, and opened as its own dialog layer like Premium payment windows.
- `Send code` must actually deliver the code to Telegram using the old V1 behavior; the web UI must not only open a local input window.
- Document/photo/file action buttons should visually match the same smooth animated button language used elsewhere.

## Web App Visual Rules

- Generated section header artwork must be visible as one masked no-repeat layer inside the header/block.
- Header artwork may visually overflow the block edge, but it must not be hidden behind gradients or clipped into a tiny fragment.
- Header text, counters, and stat cards must not cover the main artwork in a way that makes the image unreadable.
- Header artwork must never tile or duplicate horizontally.
- On desktop web, section header artwork for lesson/words/review/spelling/level/settings should sit far enough right that the visible art reaches the right edge of the header block.
- The generated app background must be visible on desktop and mobile web, in light and dark themes.
- Menu/button icons and trophy art must keep a consistent square size, transparent background, and function-specific meaning.
- Do not allow white image backplates inside icon/trophy frames.

## Mobile Web Rules

- Mobile web and desktop web must use separate forced shell modes/pages so changes for one layout do not unexpectedly break the other.
- The bottom mobile navigation must stay fixed to the viewport while scrolling.
- The bottom mobile navigation must not drift, disappear, or depend on the chat scroll container.
- The mobile top-left brand area only opens the main menu from the logo/title block. Theme, logout, language, and other header controls must not trigger menu navigation.
- Mobile award dialogs must open centered in the viewport and fit inside the screen.
- Trophy images inside mobile award dialogs must be fully visible and not clipped by the card.
- Text input/composer areas on mobile must use the available screen width and not collapse into a narrow field.
- Mobile section header artwork should fill the whole header block as a background layer, while the title text and small progress/info card keep their normal positions above it.
- Mobile main-menu section names should fit professionally: smaller type is acceptable, but labels should not break letter-by-letter and should use at most one natural wrap where possible.
- The login/auth screen must use project branding and must not show the mobile bottom navigation over the form.
- Login should stay as a separate `/login` page/state so app navigation, mobile menu, and section state do not leak into auth.
- After answers/results, mobile web must not overscroll far below the useful action. The bottom visible point should be the next primary action, for example "Next word".

## Localization Rules

- Menu labels and learner-facing controls must be localized across all 35 supported interface languages.
- Do not ship mojibake, replacement characters, broken glyphs, or placeholder question-mark strings in visible UI.
- If a translation is missing, use a clean localized/static fallback instead of corrupted text.
- Lesson/practice generated blocks must not expose English section labels like "Your task", "Situation", "Pattern", or accidental technical/style words such as "saturation" when the interface language is not English.

## Learning UX Rules

- Listening/Audition check buttons must always provide visible feedback: loading, result, or error.
- Listening/Audition pronunciation scoring must not pass the exact target phrase as an STT prompt; recognition must transcribe what the learner actually said.
- If OpenRouter STT does not return word confidence/logprobs, the app must not call the repeat perfect. Cap the score and explain that the estimate is based on transcript, speed, and gaps.
- Lesson and practice text should be structured for scanning, not dumped as one plain paragraph.
- Voice/pronunciation feedback should be clear for learners and not expose unnecessary internal model mechanics in marketing surfaces.
- Public Privacy/Terms pages should explain voice pronunciation scoring as a learner benefit and data-use note, not as an internal provider/model/JSON/confidence implementation recipe.

## 2026-05-23 Fourth V2 Audit Notes

- The `Сегодня` screen must be visible as the main learning entry point, not hidden behind the ribbon.
- Offline decks and Teacher/Admin Dashboard must preload their data without changing the selected screen after a few seconds.
- Onboarding after registration must be a real centered dialog that asks goal, level, learning language, and preferred format.
- Mistakes must be teachable: group saved errors by grammar, word order, vocabulary, politeness, and spelling, and include a visible `Потренировать похожие` action.
- Premium must show a payment history with date, plan, amount, method, and status; success dialogs must show plan, period, amount when known, and subscription expiry.
- Telegram two-factor confirmation must open above all app windows in the center of the viewport, not inside the Telegram settings card.
- New learning functions must be visible and reachable from Today: AI roleplay scenarios, pronunciation heatmap, offline mini-decks, daily quests, weekly plan, and dashboard.

## 2026-05-23 V2 Vocabulary / Roleplay / Pronounce Fix Notes

- Learn Words must show the prompt word in the selected interface language when that translation exists; Russian UI must not suddenly show the English prompt word.
- V2 tools should keep their own output surfaces. Roleplay results must stay inside Roleplay, not redirect into the generic Practice output.
- Roleplay scenarios must be localized by interface language and have their own scenario prompt structure, separate from normal practice chat.
- AI Roleplay must include more than the initial five situations; keep at least restaurant, work, travel, exam, small talk, hotel, shopping, doctor, job interview, and bank/payment.
- Pronunciation scoring must be strict for exact-repeat exercises. Bad transcripts, missing words, or substitutions must cap the score and avoid praise.
- Pronounce must be a dedicated localized block; pressing Pronounce must not auto-start or redirect into Listening/Audition.
- Vocabulary V2 loading must stay fast like V1 and must not synchronously call the AI model for every dictionary card.
- Offline decks must export a readable TXT study pack in addition to raw JSON.
- Global leaderboard rows must show which languages contributed to the user's language count.
- My Progress and Dashboard should be one Dashboard surface, not duplicated menu destinations.
- V2 should use a stable static light/dark background matching the existing glass UI instead of a moving shader background.

## 2026-05-24 V2 Output / Roleplay / Pronunciation Notes

- Every V2 menu destination should feel like its own output window. Shared infrastructure is acceptable, but the visible context must be view-specific.
- Roleplay scenario tiles should be arranged evenly in two desktop rows; after choosing a scenario, the picker should be replaced by a centered dialogue/work scene.
- Roleplay, offline decks, and every dense V2 block must be scrollable inside the app shell on desktop and mobile.
- Offline deck cards should use full-width readable rows, not clipped tiny tiles.
- Dashboard should focus on personal learning activity metrics and should not duplicate leaderboard or weak-place blocks.
- Pronunciation UI must never show raw issue keys such as `low_confidence`; weak words and sounds need learner-facing advice.
- Listening/Audition scoring must stay strict without requiring expensive models: target mismatch, missing words, substitutions, and estimated confidence must cap the score.

## 2026-05-27 Login / Mobile Regression Notes

- The login page should use the polished React/Tailwind SignIn layout, but keep current project auth: login/password, account creation, Cloudflare Turnstile, Telegram sign-in, and the existing 6-digit OTP dialog.
- Cloudflare Turnstile must render as a real challenge slot on the login/register form and the token must be submitted to the existing auth endpoints.
- Telegram sign-in from login must open Telegram with the current code workflow, then show the centered OTP dialog; do not replace it with a new unrelated auth screen.
- V2 Phrasebook/Phrases should be named Notes/Заметки in every supported interface language.
- Save-to-notes quick actions belong inside the input/composer block below the send controls and must save useful phrases, not placeholder punctuation such as `?`.
- On mobile, the compact top quick controls belong only on Today; other sections need the vertical space for output and input.
- Mobile Tools should first choose a tool, then show a focused work window with the send icon inside the input area near the lower right.
- Mobile payment requisites must never be covered by the bottom navigation after a method such as TON is selected; the modal layer should sit above nav and scroll.
- Trainer result cards should appear immediately under the active prompt/heading after an answer, not far below the option grid.
- Add Playwright desktop/mobile tests for these behaviors before deployment so the same regressions are caught automatically.

## 2026-05-27 Public Landing / Final Regression Notes

- The public landing header must not show a visible `V2` link or pill. Use a compact light/dark theme toggle there instead.
- The landing should sell the real product without overload: Daily route, Roleplay, Pronunciation, Photo/translation, Offline/PWA, Notes/mistakes, Telegram sync, and progress.
- The first landing viewport should feel unique to Poliglot AI by combining the app's dark glass visual language with a Three.js anomalous-matter shader and animated CTA buttons.
- Public Privacy and Terms language selectors must stay readable in dark theme.
- Mobile and desktop Playwright tests must cover the public landing shader, header theme toggle/no-`/app/v2` link, legal dark theme, V2 mobile Back behavior, and active Roleplay choose-another-scenario control.
- Deploy/startup must keep the site alive while large vocabulary indexes are created; dictionary optimization can run in the background, but `/healthz` and the web app should respond quickly after restart.

## 2026-05-28 Vocabulary / Mobile Scroll / Landing Notes

- Learn Words and Review must not feel like they pull wrong answers from the whole dictionary. Distractors should stay inside the learner's CEFR band: `A1/A2`, `B1/B2`, or `C1/C2`.
- Word cards should show one target-language answer per option, while the prompt/context gives several native-language meanings when the dictionary has them.
- Sparse vocabulary bands should be improved by rebuilding/enriching dictionaries, not by mixing unrelated higher-level words into easy rounds.
- Mobile Mistakes must let the user scroll through the full error dictionary; the dictionary must not be clipped halfway down the screen.
- Mobile Notes/Saved phrases must scroll below the fixed bottom menu so cards and remove/listen buttons are not hidden.
- Roleplay must provide the same kind of audio playback for the AI prompt/question that Practice has.
- The public landing light-theme hero should keep the same dark shader/glass first viewport as the dark theme because that version fits the product best.
- Every new visible label added for these areas must be covered by the 35-language localization test and the encoding-artifact check.

## 2026-05-28 Deployment Upload Notes

- Future deploys should use the permanent secured deploy-upload service instead of recreating temporary Caddy routes when PuTTY/SSH streaming breaks on large artifacts.
- The deploy-upload token must stay server-side and out of user-facing answers; Codex should fetch it through SSH only when uploading artifacts.
- Keep the Caddy deploy template in sync with the live server route so a normal deploy does not remove `poliglot-deploy-upload.service` access.

## 2026-05-28 Mobile Rail / Dictionary Notes

- Mobile web should not have a separate `More` destination anymore. All destinations belong in one horizontal bottom rail.
- The bottom rail must scroll smoothly, include the former More items, and stay visually compact with icon artwork and localized titles only.
- Long-press menu editing should feel like Android home-screen icon movement: the dragged item lifts, neighboring items open a visible gap, and the rail auto-scrolls when the finger reaches the left or right edge.
- Landing light theme should keep the first hero visually identical to the dark glass/shader hero; black text over the shader cards is not acceptable.
- Learn Words and Review should keep random words and distractors inside the correct CEFR band, while missing native-language context can be generated once by AI and then cached persistently in SQLite.

## 2026-05-28 Listening / Mistakes / Vocabulary Regression Notes

- Listening/Audition should not show the same phrase twice. Keep one active repeat panel and place the spoken-model audio inside that lower panel.
- Spelling must always show the correct target-language spelling after a correct answer or after the learner gives up.
- Mistakes must open as a list first. Selecting a mistake opens a separate practice screen, and finishing/cancelling returns to the list.
- Roleplay role cards need readable spacing between the role name and description.
- Referral, invitation, and status copy must stay localized for all 35 interface languages.
- AI-generated vocabulary clues must be in the interface/native language and must not simply repeat the hidden learning-language word.
- Mobile leaderboard and other dense lists must leave enough bottom padding so the fixed rail never covers readable rows.

## 2026-05-29 Final Regression Closure Notes

- Privacy contact cards must always show the Telegram bot as `@poliglot_ai_bot`; a blank Telegram bot card at the bottom of Privacy is a regression.
- Login should be a centered auth form over the animated background on desktop and mobile; no decorative right-side panel should return.
- Mobile menu reordering should behave like a real long-press drag: the dragged tile moves, neighbors open space, the rail auto-scrolls near edges, and the order survives reload.
- Vocabulary prompts should keep the main native-language meaning as the prompt and preserve one or two close alternate meanings as context/hint, while never showing the target answer as a clue.
- All 35 interface languages should also be real learning-language choices with local free-license dictionaries, not just UI locales.
- New dictionary imports must keep source files in `tools/sources/`, document the license, rebuild generated JSON, and run regression tests before deployment.
- Dictionary deploys must ship the generated `vocabulary_words*.json` files to `/opt/aibot` through the permanent deploy-upload package, not only the Go binary and React assets.

## 2026-05-29 Full Localization / Desktop Menu Notes

- All V2 web screens must render clean text for every supported interface language, including newly added locales; mojibake, English/Russian leftovers, and technical labels like `Section: awards` or `Mục: awards` are regressions.
- Georgian, Vietnamese, and other newly expanded locales must localize not only the left/bottom menu, but also submenu titles, roleplay scenario cards, Today blocks, dashboard labels, invite/referral text, and empty/help states.
- The desktop function ribbon should support the same real reorder behavior as mobile: the dragged tile visibly moves, neighbors open a target gap, edge drag scrolls the ribbon, and the custom order survives reload.
- Localization regressions should be broad enough to visit the major V2 screens in every interface language on desktop and mobile, not just check a small key list.

## 2026-05-29 Mobile Rail / Deep Localization Notes

- Mobile bottom-rail reorder should behave like a phone launcher dock: long-press lifts the tile, dragging immediately moves the real item through the rail, neighboring tiles make room while the finger moves, and releasing automatically saves and exits the drag without any Done/X button.
- The 35-language web localization check must include old and new languages and cover nested surfaces such as Notes, Roleplay cards, Dashboard metrics, Premium/payment requisites, Today learning-lab cards, and compact daily quests.
- Premium/payment instructions must never stay in Russian when the interface language is Thai, Vietnamese, Georgian, or any other supported non-Russian language.
- Today daily quests should stay compact: 3-4 small actionable items are preferred over a large checklist that pushes the useful content down.

## 2026-05-29 Mobile Launcher Drag / Menu Localization Notes

- Menu labels such as Spelling, Progress, and Mistakes must be localized in every one of the 35 interface languages; English fallback in non-English menus is a regression.
- Mobile web menu dragging should start as soon as the finger moves after pressing a tile, similar to HyperOS/Android launcher behavior: the dragged tile follows the finger, other tiles make room immediately, and releasing saves the tile exactly where it was dropped.
- Mobile menu drag must not self-cancel after 1-2 seconds while the finger is still held down.
- Offline Notes must paginate exactly like Vocabulary: show 10 note cards per page and use the same Back/Next page switching instead of rendering all notes in one long list.
- Mobile Roleplay scenario cards must grow with translated text instead of clipping or letting text run outside the card, and the final card must remain reachable above the bottom rail.
- Playwright regression runs should clear old service worker/cache state before opening V2 so tests always check the current build.
- Mobile bottom-rail editing should prioritize reliable reordering over animation: hold a tile for 2 seconds to enter edit mode, show a visible `Done` button above the menu on the right, save order on release, and allow multiple reorders before leaving edit mode.
- Normal mobile bottom-rail swiping must keep working outside edit mode; drag/reorder gesture blocking should apply only after edit mode is active.

## 2026-05-30 Mobile Web Menu Reorder Notes

- Mobile web menu blocks must be reorderable by touch: hold a menu tile, drag it to the desired place, release to save the new order, and keep the order after reload.
- Reorder targeting should follow the actual finger position, so it remains reliable if the menu layout changes from a horizontal rail to wrapped rows later.

## 2026-05-30 Mobile Text Fit Notes

- Mobile V2 cards must grow to fit long translated text instead of letting text run below or outside the card.
- Vocabulary cards need enough right-column space for the listen button while long translations wrap inside the left column.
- The level-test skip action should display as `Skip` even on the Russian interface; `Раздел` on that button is a regression.
- Leaderboard and roleplay outer panels must close after the final visible list item, not midway through the list while lower cards continue outside the border.
- Saving a phrase to Notes/Phrasebook from a quick-save action should fill the translation field automatically in the current interface language through OpenRouter when the user did not type a manual note.
- The Skip button text must fit as one word; `Ski` with the `p` dropped to the next line is a regression.
- Phrasebook/Notes should paginate once there are more than 10 saved phrases, using the same page pattern as Offline cards.
- Mistakes should also paginate after 10 saved errors, using the same page controls as Offline and Phrasebook.
- The first loading screen should have a visible animated background like the login screen; a mostly static dark backdrop is a regression.

## 2026-05-30 Level / Spelling / Lesson Audio Notes

- Level-assessment question instructions must always be in the interface language for all 35 interface languages.
- Level-assessment prompts must keep the original task logic: do not replace real questions such as "How do you say..." or "Fill the gap..." with a generic "Choose the correct answer" line.
- Static level-test question types such as translation, gap fill, word meaning, correct sentence, grammar form, and natural phrasing must remain visibly different after localization.
- Level-assessment examples inside questions stay in the learning language, and answer options must be generated in the selected learning language.
- Mobile Spelling should not reveal the correct spelling after a wrong attempt; reveal it only after a correct answer or when the learner chooses `Не знаю`.
- Mobile Spelling result panels must stay reachable above the fixed bottom navigation while the keyboard/input flow is active.
- Lesson example audio must read the learning-language example/model phrase, not the interface-language task instruction.

## 2026-05-30 Learn Words Result Notes

- Learn Words success results must show both sides of the learned pair: the interface-language translation and the target-language word that the learner selected.
- The Learn Words success result should offer the same compact Notes save affordance used elsewhere: note icon plus the word/translation pair below `Next`, not an unrelated large extra button.
- Saving from that result should store the target-language word as the phrase and the interface-language translation as the note/context.
- Future regression tests should cover Learn Words result content, Notes saving, all 35-language level-test task localization, wrong Spelling answer reveal rules, lesson example audio language, and clean non-mojibake UI text.
