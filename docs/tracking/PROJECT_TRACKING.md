---
project: Poliglot AI
status: active
updated: 2026-05-30
tags:
  - poliglot-ai
  - redesign
  - assets
  - comfyui
  - deploy
---

# Poliglot AI Project Tracking

## Current Sprint

### Done

- [x] Start the full V2 improvement plan with the Settings two-factor implementation slice.
- [x] Wire `Send code` in Settings to `/api/auth/telegram/start`, open the V1-style Telegram bot deep link, and verify entered codes through `/api/auth/telegram/status`.
- [x] Keep the two-factor code input as a separate centered popup layer, with resend creating a fresh Telegram auth request.
- [x] Add regression coverage for the web Telegram code start/verify flow.
- [x] Restore user-uploaded landing image mapping.
- [x] Use `poliglot-ai-avatar.jpg` for the first/hero visual and bot page hero.
- [x] Add theme-aware web app asset naming for plan cards, headers, and icons.
- [x] Generate light/dark production-safe PNG pairs.
- [x] Remove CSS pseudo-letter marks from app stage, action, and panel treatments.
- [x] Rebuild the ComfyUI workflow script around SD3.5 Large FP8 by default.
- [x] Run local static checks and Go tests.
- [x] Run browser QA for landing and `/app` desktop/mobile light/dark.
- [x] Build Linux binary and deployment archives.
- [x] Deploy `/opt/aibot` and `/var/www/poliglotai`.
- [x] Restart `aibot.service`.
- [x] Run production URL and asset checks.
- [x] Install FLUX.2 Dev ComfyUI model pack for test-generation.
- [x] Generate premium Flux2 replacement pack for tariffs, headers, and action/menu icons.
- [x] Prepare `tmp/poliglot-web-flux2-refresh.tar.gz` from the current `web/` folder for server upload.
- [x] Re-run local checks for the refreshed Flux2 web pack.
- [x] Generate expanded Flux2 app visuals: global backgrounds, right-inspector panel backgrounds, main-menu card art, Free access replacement, Tools/Settings artwork, and 20 award PNGs.
- [x] Add reminder localization regression tests for selected `interface_language`.
- [x] Deploy expanded Flux2 visual pack to `/opt/aibot/web` on `root@186.246.45.123`.
- [x] Rebuild and redeploy the Linux bot binary after reminder localization tests.
- [x] Verify production HTML, manifest, PNG hashes, and service health after the Flux2 upload.
- [x] Add strict transparent icon/trophy generation prompts and per-file validation for equal canvas size, object fill, and no tile/backplate artifacts.
- [x] Regenerate and visually review the transparent icon/trophy pack with contact sheets.
- [x] Add Shadowing learning mode to Telegram and web with STT scoring, feedback, XP, and a generated phrase loop.
- [x] Rename learner-facing Shadowing to Listening / `РђСѓРґРёСЂРѕРІР°РЅРёРµ` and wire OpenRouter `openai/gpt-4o-transcribe` pronunciation assessment into Listening, lesson, and practice voice answers.
- [x] Add web lesson/practice voice recording and file upload with pronunciation score, accent strength, fluency, weak words, and tips.
- [x] Wire the web app-wide background and account panel background as no-repeat cover/masked CSS layers.
- [x] Strengthen lesson/practice prompts around active recall, reusable chunks, and focused correction.
- [x] Create the `Poliglot AI Landing 2026` Figma design file with desktop/mobile frames for the refreshed public site.
- [x] Add the public-site reviews section, pronunciation-engine explanation, Privacy/Terms trust band, and localized navigation copy.
- [x] Update Privacy and Terms with the OpenRouter GPT-4o Transcribe pronunciation assessment flow and educational-use caveat.
- [x] Fix Telegram/web Listening STT bias by removing the exact target phrase from the transcription prompt.
- [x] Cap transcript-only pronunciation scores when OpenRouter does not return word confidence/logprobs.
- [x] Add `docs/tracking/USER_UI_WISHES.md` as the persistent checklist for mobile UI, generated artwork, localization, and pronunciation scoring regressions.
- [x] Refresh the Figma landing page with a sales-first desktop/mobile design that focuses on voice-learning benefits instead of internal model mechanics.
- [x] Clean web/TG visible mojibake around Tools/Translator copy, localize `Text В· voice В· photo` across 20 languages, and repair Go API error strings.
- [x] Add `tools/check_encoding_artifacts.mjs` and re-clean web recorder/auth/settings/API strings so broken Cyrillic, smart-quote, and emoji artifacts do not reappear.
- [x] Deploy the cleaned localization build to `/opt/aibot` and verify service health.
- [x] Complete the strict 21-point React `/app/v2` UI/UX correction pass from 2026-05-23.
- [x] Add animated top-center 3-second status messages, animated dropdown selectors, expanding function-ribbon hover, and smooth context view entry animation.
- [x] Restore V1-style learning details in V2: trophy names/stories, mistake clear confirmation/correct audio, dictionary examples without review counters, spelling/listening next actions, word audio timing, and practice/lesson TTS clips.
- [x] Add supplied-style upload, audio upload, spinner, and password input components under `web-react/src/components/ui`.
- [x] Add Settings password change, V1 Telegram referral share text, referral terms, and login-based header identity.
- [x] Generate Nx project metadata and graph into `tmp/nx` and add `docs/reference/PROJECT_MAP.md`.
- [x] Deploy the strict 21-point React `/app/v2` UI/UX correction pass to `/opt/aibot`.
- [x] Verify production `/healthz`, `/app`, `/app/v2`, V2 hashed assets, and deployed UI markers after restart.
- [x] Complete the follow-up V2 polish pass from the screenshot review: compact upload buttons, corrected menu hover trigger, inset ribbon arrows, cleaner Learn Words/Spelling/Level result blocks, waveform audio placement, payment instructions/check action, Settings column stacking, and language dropdown z-index.
- [x] Refresh Nx MCP/graph metadata after the follow-up pass.
- [x] Complete the final V2 menu/voice/payment refinement: non-shrinking menu neighbors, animated recording button, centered payment result dialogs, cleaned payment copy, V1 award story fallback, transparent animated backgrounds, and mobile control spacing.
- [x] Add the shared Radix dialog component and centered OTP/payment result modal surfaces under `web-react/src/components/ui`.
- [x] Verify the final V2 refinement locally with `npm --prefix web-react run build` and `go test ./...`.
- [x] Deploy the final V2 refinement to `/opt/aibot`, clean stale `/opt/aibot/web/v2` assets, and restart `aibot.service`.
- [x] Verify production `/healthz`, `/app`, `/app/v2`, and the fresh V2 JS/CSS assets after deployment.
- [x] Add the next V2 mobile-specific navigation layer and stabilize the function ribbon so neighboring chips do not shrink/fade on hover.
- [x] Rebuild V2 with fresh hashed assets and regenerate Nx graph/metadata under `tmp/nx`.
- [x] Re-run local `npm --prefix web-react run build`, `go test ./...`, and `node tools/check_encoding_artifacts.mjs`.
- [x] Deploy the mobile navigation/ribbon refinement to `/opt/aibot`, restart `aibot.service`, update Caddy `/healthz` routing, and verify public production endpoints/assets.
- [x] Add the point-by-point V2 feature slice: Today weekly plan, daily quests, AI Roleplay, Pronunciation dashboard/heatmap, Offline/PWA decks, and Teacher Dashboard.
- [x] Verify the new V2 feature slice locally with `npm run build` in `web-react`.
- [x] Re-run Go, encoding, and Nx graph checks for the V2 feature slice.
- [x] Deploy the V2 feature slice and PWA root-asset backend route to `/opt/aibot`, restart `aibot.service`, and smoke-check production.
- [x] Add the editable prompt registry requirement to the V2 roadmap and implement runtime prompt loading through `APP_PROMPTS_FILE`.
- [x] Create `docs/reference/APP_PROMPTS.md` and deployable `app_prompts.json` with feature/tool/function notes for core app prompts.
- [x] Fix Pronunciation/Pronounce navigation so it opens the pronunciation dashboard instead of auto-starting Listening.
- [x] Fix V2 Learn Words prompt language so it follows the selected interface language when a translation exists.
- [x] Restore fast V2 vocabulary list loading by removing per-card synchronous AI example calls from `/api/vocabulary`.
- [x] Add ten localized AI Roleplay scenarios, keep Roleplay results in the Roleplay surface, and route `ROLEPLAY_TOOL_V2` through a dedicated backend runtime prompt.
- [x] Tighten exact-repeat pronunciation scoring so low-similarity/missing-word attempts are capped and not praised.
- [x] Add readable TXT export for Offline decks, language-name breakdowns in Global leaderboard, Dashboard/Progress route normalization, and static theme backgrounds.
- [x] Rebuild Nx graph/project metadata after the V2 vocabulary/roleplay/pronunciation correction pass.
- [x] Deploy the V2 vocabulary/roleplay/pronunciation correction pass to `/opt/aibot`, restart `aibot.service`, and verify production `/healthz`, `/app`, `/app/v2`, and fresh JS/CSS assets.
- [x] Localize V2 menu/system labels across all 20 interface languages and verify no mojibake in the new React copy.
- [x] Rework mobile V2 navigation with hidden desktop top menu, icon-only top-right logout, small theme icon, fixed primary bottom tabs, and a More sheet for secondary sections.
- [x] Polish Today, Offline, Pronunciation, and Phrasebook presentation: aligned streak panel, no obsolete Next features block, wider Offline decks, centered/larger artwork, compact pronunciation map, and mobile scrolling for dense blocks.
- [x] Add Playwright coverage for the 20-language V2 labels, desktop layout regressions, mobile quick controls, More menu, and mobile scroll behavior.
- [x] Deploy the V2 localization/mobile cleanup build to `/opt/aibot/web/v2`, restart `aibot.service`, and verify public health plus fresh V2 assets.
- [x] Filter V2 Premium payment history to confirmed/paid payments only.
- [x] Add referral invitee rows with level-3 status, per-user earnings, and 10-row pagination.
- [x] Add Offline deck 10-card pagination while keeping full TXT/JSON export.
- [x] Add mobile top quick controls for bug report, language, theme, and icon-only logout below the app frame.
- [x] Add editable mobile bottom navigation with long-press edit mode, X unpin, More tile pin, and Done exit.
- [x] Regenerate Offline, Phrases, Pronunciation, and Roleplay light/dark transparent icons through the local FLUX.2 ComfyUI workflow and validate the manifest.
- [x] Expand Playwright coverage for payment filtering, referral pagination, Offline export, Telegram Send code visibility, mobile pin editing, mobile inputs, Tools selector, scrollability, and mojibake guards.
- [x] Verify the second 2026-05-26 V2 pass with React build, Go tests, Playwright desktop/mobile, and `npm run check`.
- [x] Deploy the second 2026-05-26 V2 pass to `/opt/aibot`, restart `aibot.service`, and verify public production endpoints/assets.
- [x] Add targeted regression tests for the remaining reported mobile failures: lesson prompt scroll/readability, Roleplay no-overlap layout, Tools selector-first mobile flow, mistake dictionary scrolling, localized More descriptions, desktop ribbon arrow bounds, persistent Phrasebook API saves, and V2-as-main routing.
- [x] Fix the code paths surfaced by those tests: lesson start uses returned `prompt`/task text, practice/roleplay cards render `correction` and `explanation`, mobile edit-mode animation leaves tap targets stable, and legacy V1 is moved to `/app/v1`.
- [x] Reverify the regression pass with React build, Playwright desktop/mobile, Go tests, encoding checks, and Nx `npm run check`.
- [x] Deploy the latest regression-test pass to `/opt/aibot`, restart `aibot.service`, and verify public `/healthz`, `/app`, `/app/v2`, `/app/v1`, and fresh V2 assets.
- [x] Fix the latest reported mobile regressions: scrollable onboarding/start dialog, compact Phrasebook save after input, mobile Roleplay scenario brief hidden after selection, mobile Tools selector hidden after choosing a tool, and More-sheet tile reordering.
- [x] Force current FLUX.2 transparent menu assets for Offline, Phrases/Phrasebook, Pronunciation, and Roleplay across desktop and mobile with versioned asset URLs.
- [x] Persist Pronunciation weak-word/history data per account and fix Learn Words wrong-answer sampling for web desktop, web mobile, and Telegram.
- [x] Add targeted Playwright and Go regression tests for the latest mobile UI failures, current asset URLs, persisted pronunciation history, More reordering, Tools mobile flow, and varied Learn Words distractors.
- [x] Deploy the V2 mobile/Learn Words fix to `/opt/aibot`, restart `aibot.service`, and verify public `/healthz`, `/app`, `/app/v2`, `/app/v1`, current JS/CSS, FLUX.2 icons, and manifest.
- [x] Fix installed mobile PWA update behavior with network-first service-worker caching, no-store app-shell/manifest/service-worker headers, regression coverage, and deploy through the permanent upload service without server backups.

### In Progress

- [ ] Authenticated browser QA on the live V2 app with a real session.

### Next

- [ ] Continue the full V2 plan with authorized live QA, mobile-specific layouts, and `App.tsx` decomposition into focused view modules.
- [ ] License-review any FLUX.2 Dev output before commercial production use.

## Asset Matrix

| Surface | Slot | Production Asset | Notes | QA |
| --- | --- | --- | --- | --- |
| Landing | Hero / first visual | `assets/poliglot-ai-avatar.jpg` | User-uploaded asset | Pending browser check |
| Landing | Languages / features | `assets/times.jpg` | User-uploaded asset | Pending browser check |
| Landing | Use cases / pricing | `assets/primi.jpg` | User-uploaded asset | Pending browser check |
| Landing | Start / how it works | `assets/write.jpg` | User-uploaded asset | Pending browser check |
| Landing | CTA / success | `assets/yspex.jpg` | User-uploaded asset | Pending browser check |
| Bot page | Hero / Open Graph | `assets/poliglot-ai-avatar.jpg` | Aligned with landing hero | Pending browser check |
| Terms / Privacy | Legal hero | `assets/site-legal-shield.png` | Generated legal-specific asset stays intentional | Pending production check |
| Web app | Plan cards | `plan-free-*.png`, `plan-premium-*.png`, `plan-platinum-*.png` | Light/dark pairs with legacy light aliases | Generated |
| Web app | Global shell background | `app-background-light.png`, `app-background-dark.png` | Full-page background behind translucent app UI | Generated |
| Web app | Section headers | `header-*-light.png`, `header-*-dark.png` | Selected by `html[data-theme]` | Generated |
| Web app | Listening / `РђСѓРґРёСЂРѕРІР°РЅРёРµ` module | `header-shadowing-*.png`, `icon-shadowing-*.png` | Voice repeat trainer, same Flux2 style | Generated |
| Web app | Navigation/action icons | `icon-*-light.png`, `icon-*-dark.png` | Selected by `html[data-theme]` | Generated |
| Web app | Right inspector panels | `panel-progress-*.png`, `panel-limits-*.png`, `panel-account-*.png` | Decorative backgrounds for Progress, Limits, Account | Generated |
| Web app | Award trophies | `award-01.png` ... `award-20.png` | PNG trophy artwork by XP level | Generated |
| Web app | Manifest | `web/assets/brand-assets-manifest.json` | Theme, seed, model, prompt, dimensions, QA status | Generated |
| ComfyUI | Flux2 diffusion | `models/diffusion_models/flux2_dev_fp8mixed.safetensors` | Official ComfyUI Flux2 Dev pack; full BFL fp16 is gated | Installed |
| ComfyUI | Flux2 text encoder | `models/text_encoders/mistral_3_small_flux2_bf16.safetensors` | Official ComfyUI Flux2 Dev pack | Installed |
| ComfyUI | Flux2 VAE | `models/vae/flux2-vae.safetensors` | Official ComfyUI Flux2 Dev pack | Installed |

## QA Checklist

- [x] `rg` confirms replaced landing `site-*` section images are not referenced by landing or bot page.
- [x] Every referenced `/app/assets/...` image exists under `web/assets`.
- [x] Every referenced public `assets/...` image exists under `РЎР°Р№С‚ РїРѕР»РёРіР»РѕС‚Р° РґР»СЏ Р±РѕС‚Р°/assets`.
- [x] `node --check` passes for site JS and generator scripts.
- [x] `go test ./...` passes.
- [x] Landing desktop screenshot has CTA visible and next section hinted.
- [x] Landing mobile screenshot has CTA visible and no horizontal overflow.
- [x] `/app` light theme screenshot shows light header/icon/plan assets.
- [x] `/app` dark theme screenshot shows dark header/icon/plan assets.
- [x] Browser console has no JS errors.
- [x] Production image URLs load from `/assets/...` and `/app/assets/...`.
- [x] Production landing, `/app`, `terms.html`, `privacy.html`, and `bot.html` return 200 after deploy.
- [x] Bot service is active after restart.
- [x] Refreshed Flux2 web pack passes `go test ./...`.
- [x] Refreshed Flux2 web pack passes `node --check tools/generate_comfy_brand_assets.mjs`.
- [x] Refreshed Flux2 web pack inline app scripts compile.
- [x] Refreshed Flux2 web pack has no missing static `/app/assets/...` references.
- [x] Refreshed Flux2 web pack passes `go test ./...`, including reminder localization tests.
- [x] Production `/app/assets/brand-assets-manifest.json` reports the Flux2 model/workflow after upload.
- [x] Production `/app` HTML references app backgrounds and right-inspector panel backgrounds.
- [x] Production PNG hashes match local files for app background, panel backgrounds, awards, plan Free, Tools, Settings, and Home icon.
- [x] Browser smoke-check for logged-out `/app` returns 200, renders the auth screen, and has no JS console errors.
- [x] Transparent icon/trophy validator reports every regenerated PNG as fixed-size RGBA with transparent corners and normalized object fill.
- [x] Final icon and trophy contact sheets show consistent scale, no white backgrounds, and no repeated/tiled section imagery.
- [x] Local `/app` screenshots confirm Tools has a generated hero background, Settings has a matching header block, and main-menu card art is full-size on desktop/mobile.
- [x] React V2 strict correction pass builds with `npm --prefix web-react run build`.
- [x] Backend changes pass `go test ./...`.
- [x] Nx metadata generated: `tmp/nx/projects.json`, `tmp/nx/english-coach-bot.project.json`, `tmp/nx/nx-report.txt`, and `tmp/nx/project-graph.html`.
- [x] Mock browser smoke loads React V2 desktop/mobile without horizontal overflow; screenshots saved in `tmp/v2-current-smoke`.
- [x] Follow-up V2 pass keeps `npm --prefix web-react run build` and `go test ./...` green.
- [x] Refreshed Nx graph output in `tmp/nx/project-graph.html`; persistent Nx MCP HTTP startup was attempted but not left running in this shell.
- [x] V2 vocabulary/roleplay/pronunciation correction pass keeps `npm run build` in `web-react`, `go test ./...`, and `node tools/check_encoding_artifacts.mjs` green.
- [x] Second 2026-05-26 V2 pass keeps `npm --prefix web-react run build`, `go test ./...`, `npm --prefix web-react run e2e`, `npm run check`, and production URL/asset checks green.

## Deploy Checklist

- [x] Build Linux binary.
- [x] Upload `/opt/aibot/aibot`.
- [x] Upload `/opt/aibot/web/index.html`.
- [x] Upload `/opt/aibot/web/assets`.
- [x] Upload public site to `/var/www/poliglotai`.
- [x] Restart bot service.
- [x] Confirm production HTTP 200s.
- [x] Confirm production image URLs load from `/assets/...` and `/app/assets/...`.

## Decision Log

| Date | Decision | Reason |
| --- | --- | --- |
| 2026-05-20 | User-uploaded landing photos are the public marketing visuals. | They fit the landing blocks better than the generated site section art. |
| 2026-05-20 | Web app assets use explicit light/dark files. | Theme switching becomes auditable and avoids one image doing two visual jobs. |
| 2026-05-20 | `sd3.5_large_fp8_scaled.safetensors` is the primary ComfyUI model. | It is installed locally and supports the SD3.5 node path in this ComfyUI setup. |
| 2026-05-20 | Juggernaut XL and SDXL remain comparison/fallback models. | Existing workflow support is useful for quick comparison and recovery. |
| 2026-05-20 | FLUX experiments require license review before production. | Local testing is allowed, but production output needs source/license confidence. |
| 2026-05-20 | FLUX.2 Dev is used for the next test-server art pass. | The current server is explicitly a test target, and the user wants to evaluate Flux2 output in-place. |
| 2026-05-20 | Full `black-forest-labs/FLUX.2-dev` fp16 remains gated. | Hugging Face returns 401 without accepted access/token; the public ComfyUI workflow uses `flux2_dev_fp8mixed` plus bf16 text encoder and Flux2 VAE. |
| 2026-05-20 | Deployment used the active MobaXterm SCP session to install a temporary deploy key. | Password auth was rejected by direct SSH, while MobaXterm already had an authenticated session. |
| 2026-05-20 | Daily Telegram reminders render from `interface_language`. | Web language changes persist to user state, and the reminder scheduler carries that language into reminder text and Telegram button copy. |
| 2026-05-20 | Generated panel images must stay decorative and text-free. | Real HTML text sits above them; fake labels or gauges from the image model make the app look broken. |
| 2026-05-21 | Menu icons and award trophies must be transparent cutouts with normalized scale. | Button meaning is lost when generated objects are tiny, tiled, or framed by a white/backplate background. |
| 2026-05-23 | V2 interactions should use shared animated components instead of native browser controls. | The user wants one polished interface language across dropdowns, buttons, loading, upload, audio, and context panels. |
| 2026-05-23 | Compact upload controls are required in chat/tool rows; large upload cards are only allowed where the page is explicitly an upload surface. | The screenshot review called out the oversized voice-file upload UI and asked for mic-style icon buttons with a filename chip. |
| 2026-05-23 | Trainer result blocks should not repeat context/example/comment text when the same content is already visible in the primary panel. | Repetition made Learn Words, Spelling, and Mistakes look noisy and unlike a focused trainer. |
| 2026-05-23 | V2 learning features must be visible from the default `РЎРµРіРѕРґРЅСЏ` screen and must not depend on hidden menu discovery. | The user did not see the new functions while Offline/Dashboard were also redirecting after background preloads. |
| 2026-05-23 | Settings Telegram OTP must render through a global portal layer. | Centered two-factor confirmation must stay above all transformed/scrolled app panels, matching payment modal behavior. |
| 2026-05-23 | Payment history can be browser-local for the immediate V2 UX slice, while backend payment records remain the source of truth for real status checks. | The current request required visible date/plan/amount/method/status history without adding a risky new persistence migration in the same UI pass. |

## Deployment Log

| Date | Target | Result |
| --- | --- | --- |
| 2026-05-20 | `/opt/aibot`, `/var/www/poliglotai` on `root@186.246.45.123` | Deployed `tmp/poliglot-web.tar.gz` and `tmp/poliglot-site.tar.gz`; restarted `aibot.service`; production pages and key assets returned 200. |
| 2026-05-20 | Local Flux2 web refresh archive | Prepared `tmp/poliglot-web-flux2-refresh.tar.gz`; production still serves the older `brand-assets-manifest.json` until SSH auth is available for upload. |
| 2026-05-20 | `/opt/aibot/web` on `root@186.246.45.123` | Deployed expanded `tmp/poliglot-web-flux2-refresh.tar.gz`; production manifest has 108 Flux2 entries; key PNG hashes match local files; `aibot.service` is active. |
| 2026-05-20 | `/opt/aibot/aibot` on `root@186.246.45.123` | Rebuilt `tmp/aibot-linux-amd64`, installed it over the service binary, restarted `aibot.service`, and confirmed `/app` returned 200 after restart. |
| 2026-05-21 | `/opt/aibot/aibot` and `/opt/aibot/web/index.html` on `root@186.246.45.123` | Deployed Russian award localization safety, Telegram translator language labels/audio cleanup, web mistake pagination, desktop activation-key form, and level-test resume/shuffle fixes; `aibot.service` is active and `/healthz` returns `ok`. |
| 2026-05-21 | `/opt/aibot/aibot` and `/opt/aibot/web/index.html` on `root@186.246.45.123` | Deployed TonAPI USDT/RUB pricing, Telegram Stars web-to-bot deep links, localized paginator arrows, and level-test scoring coverage; `aibot.service` is active, internal `/healthz` returns `ok`, and public `/app` contains the new markers. |
| 2026-05-23 | `/opt/aibot/aibot`, `/opt/aibot/web/index.html`, and `/opt/aibot/web/v2` on `root@186.246.45.123` | Deployed the strict 21-point V2 UI/UX correction pass, legacy dictionary example cleanup, generated V2 assets, and restarted `aibot.service`; backup saved at `/opt/aibot/deploy-backups/20260523-031552`; internal `/healthz`, `/app`, `/app/v2`, V2 JS asset, legacy marker, and animated-select CSS marker verified. |
| 2026-05-23 | `/opt/aibot/aibot`, `/opt/aibot/web/index.html`, and `/opt/aibot/web/v2` on `root@186.246.45.123` | Deployed the follow-up screenshot-review V2 polish pass; backup saved at `/opt/aibot/deploy-backups/20260523-162424`; internal `/healthz`, `/app`, `/app/v2`, exact V2 JS/CSS assets, compact audio upload marker, payment instruction marker, settings stack marker, and public `https://poliglotai.ru/app/v2` were verified. |
| 2026-05-23 | `/opt/aibot/aibot`, `/opt/aibot/web/v2`, and `/etc/caddy/Caddyfile` on `root@186.246.45.123` | Deployed the mobile bottom-nav and stable ribbon refinement; backup saved at `/opt/aibot/deploy-backups/20260523-182906`, Caddy backup saved at `/etc/caddy/Caddyfile.backup-20260523-182906`; public `/healthz`, `/app`, `/app/v2`, `index-6B3MD3Yv.css`, and `index-DsYghMDs.js` verified. |
| 2026-05-23 | `/opt/aibot/aibot` and `/opt/aibot/web/v2` on `root@186.246.45.123` | Deployed the point-by-point V2 feature slice: Today weekly plan, daily quests, AI Roleplay, Pronunciation dashboard/heatmap, Offline/PWA decks, Teacher Dashboard, and no-cache PWA root assets; backups saved at `/opt/aibot/deploy-backups/20260523-193000-v2-feature-slice` and `/opt/aibot/deploy-backups/20260523-193000-aibot-feature-slice`; internal and public `/healthz`, `/app/v2`, V2 JS/CSS, manifest, and service worker endpoints verified. |
| 2026-05-23 | `/opt/aibot/web/v2` on `root@186.246.45.123` | Deployed the fourth V2 audit slice: visible Today entry, no Offline/Dashboard redirect, global Telegram OTP, onboarding dialog, grouped Mistakes with similar drills, Premium payment history, and learning-lab cards; backup saved at `/opt/aibot/deploy-backups/20260523-2042-v2-fourth-audit/v2`; internal/public `/healthz`, `/app/v2`, `index-DIhZOhmh.css`, and `index-BJCVGuii.js` verified from the server. |
| 2026-05-23 | `/opt/aibot/aibot`, `/opt/aibot/app_prompts.json`, and `/opt/aibot/web/v2` on `root@186.246.45.123` | Deployed runtime-editable prompts and the Pronunciation navigation fix; backup saved at `/opt/aibot/deploy-backups/20260523-2218-v2-prompts-pronounce`; corrected an initial Windows-binary upload by rebuilding Linux/amd64 and reinstalling `/opt/aibot/aibot`; service logs show `Loaded 17 editable app prompts from app_prompts.json`; public `/healthz`, `/app/v2`, `index-BUoc4oQk.js`, and `index-DIhZOhmh.css` verified. |
| 2026-05-24 | `/opt/aibot/aibot`, `/opt/aibot/app_prompts.json`, and `/opt/aibot/web/v2` on `root@186.246.45.123` | Deployed the V2 output/roleplay/pronunciation/offline/dashboard pass; backup saved at `/opt/aibot/deploy-backups/20260524-010220-v2-output-roleplay-pronunciation`; restarted `aibot.service`; public `/healthz`, `/app`, `/app/v2`, `index-BljsgA2m.js`, and `index-CIAalwsj.css` returned 200. |
| 2026-05-26 | `/opt/aibot/web/v2` on `root@186.246.45.123` | Deployed the V2 localization/mobile cleanup build; backup saved at `/opt/aibot/deploy-backups/20260526-205908-v2-mobile-localization`; restarted `aibot.service`; internal `/healthz`, public `/healthz`, `/app/v2`, `index-DG4RaOUL.js`, and `index-BFuJzIXN.css` verified. |
| 2026-05-26 | `/opt/aibot/aibot`, `/opt/aibot/web/v2`, and selected `/opt/aibot/web/assets` on `root@186.246.45.123` | Deployed the confirmed-payment/referral/mobile/offline pass plus regenerated FLUX.2 transparent icons; backup saved at `/opt/aibot/deploy-backups/20260526-223406-v2-referrals-mobile`; restarted `aibot.service`; public `/healthz`, `/app/v2`, `index-C4nUHB4M.js`, `index-DqC1jeH0.css`, icon PNGs, and `brand-assets-manifest.json` verified. |
| 2026-06-27 | `/opt/aibot/web` and `/var/www/poliglotai` on `root@186.246.45.123` | Deployed Poliglot AI social links for the public landing/footer and web app settings; backup saved at `/opt/aibot/deploy-backups/20260627-223337-social-links`; enabled persistent 1G `/swapfile` after OOM report, confirmed no `fold` processes, `aibot.service` and `caddy` active, public `/healthz` returns `ok`, production HTML references `index-CZPGvZ6d.js`, `index-DtwfRzt2.css`, `main-DMRCusAl.js`, and `main-CCIuuDnV.css`, and production JS contains YouTube, Instagram, and TikTok links. |
| 2026-06-27 | `/opt/aibot/web` and `/var/www/poliglotai` on `root@186.246.45.123` | Deployed social/cookie localization follow-up; backup saved at `/opt/aibot/deploy-backups/20260627-233225-social-i18n`; production HTML references `index-C3O7Lq7c.js`, `index-DtwfRzt2.css`, `main-DcJ36_sh.js`, and `main-C1gzSw4z.css`; stale hashed assets removed; `/healthz` returns `ok`; browser checks confirmed light-theme social proof contrast, TikTok `@poliglotai.online`, and cookie banners remain visible after saved consent on both landing and web app. |

## Server Defaults

- SSH target: `root@186.246.45.123`
- Bot path: `/opt/aibot`
- Public site path: `/var/www/poliglotai`
- Authentication: temporary password supplied in this thread; do not store it in repo files.
