# Project Map

Generated/updated from Nx metadata on 2026-05-26 after the V2 phrasebook, habit calendar, bug-report, offline grouping, and Playwright coverage pass. The current Nx project graph is saved at `tmp/nx/project-graph.html`; supporting metadata is in `tmp/nx/projects.json`. Nx graph generation completed, while Nx also printed a non-blocking local git-base warning because this workspace has no `main` ref.

## Nx Index

Nx reports two workspace projects:

- `english-coach-bot` - root Go application and orchestrator for backend, Telegram bot, legacy web app, React v2 app, and public site builds.
- `web-static` - static web dependency declared by the root app.

Primary `english-coach-bot` targets:

- `build` - `go build -ldflags "-s -w" -o aibot .`
- `test` - `go test ./...`
- `check` - `go test ./...` then `node tools/check_encoding_artifacts.mjs`
- `web-build` - `npm --prefix web-react run build`, output to `web/v2`
- `web-dev` - `npm --prefix web-react run dev`
- `site-build` - `npm --prefix site-react run build`
- `generate-activation-keys` - serial key generator
- `rebuild-public-translations` - public site translation rebuild

Nx environment snapshot:

- Node `24.15.0`
- npm `11.14.1`
- Nx `21.6.11`
- Nx MCP HTTP endpoint was verified on port `9921` during this pass.

## Runtime Surfaces

- Go backend and Telegram bot: root `*.go` files, with entry points in `main.go`, Telegram flow in `bot.go`, web API in `web_api.go`, auth in `web_auth.go`, persistence in `sqlite_store.go` and `storage.go`.
- Legacy private web app: `web/index.html` plus `web/assets`.
- React private web app V2: `web-react/src`, built into `web/v2`.
- Editable AI prompts: `docs/reference/APP_PROMPTS.md` is the human navigation index; `app_prompts.json` is the runtime prompt file deployed next to the binary/env files and selected with `APP_PROMPTS_FILE`.
- Public landing/legal site: `site-react/src`, built into the public static-site folder.
- Shared docs and tracking: `docs/technical/DOCUMENTATION.md`, `docs/tracking/CHANGELOG.md`, `docs/tracking/USER_UI_WISHES.md`, `docs/tracking/PROJECT_TRACKING.md`, `docs/reference/PROJECT_MAP.md`, `MCP_STARTUP.md`, `docs/product/WEB_V2_FUNCTIONALITY.md`.

## V2 React Structure

- Main app and API wiring: `web-react/src/App.tsx`
- App styling and animations: `web-react/src/styles/app.css`
- Typed API data: `web-react/src/lib/types.ts`
- Habit calendar component: `web-react/src/components/ui/calendar-rac.tsx`
- Playwright V2 smoke suite: `web-react/e2e/v2-smoke.spec.ts`
- Playwright config: `web-react/playwright.config.ts`
- Vite local asset middleware for dev/preview `/app/assets`: `web-react/vite.config.ts`
- Shared UI components: `web-react/src/components/ui`
- Current added/updated controls: `file-upload.tsx`, compact `audio-upload-card.tsx`, `spinner.tsx`, `password-input.tsx`, `audio-wave-button.tsx`, `voice-input.tsx`, `dialog.tsx`, `otpdialog.tsx`, `chat-interface.tsx`, `calendar-rac.tsx`, and `morphing-arrow-button.tsx`
- Mobile V2 navigation: `MobileBottomNav` in `web-react/src/App.tsx`, styled by `.mobile-bottom-nav-v2` in `web-react/src/styles/app.css`
- Pronunciation navigation rule: the `pronunciation` view opens `PronunciationDashboardView` only; Listening/Audition starts from explicit dashboard buttons.
- Current learning surfaces: `HomeView` (`Today`, daily quests, weekly plan, V2 learning lab, habit calendar, daily bonus), `RoleplayView` with ten localized scenarios and a dedicated dialogue session surface, `PronunciationDashboardView` with cleaned learner-facing weak-word labels, `PhrasebookView` for saved/favorite lesson phrases, `OfflineDecksView` with grouped JSON/TXT deck cards, `TeacherDashboardView` merged with the old Progress destination and focused on learning activity metrics, grouped `MistakesView`, `PremiumView` with payment history, and `SettingsView` with global Telegram OTP.

## Backend Areas Used By V2

- Lesson/practice audio text: `bot.go`, `web_api.go`
- Vocabulary examples: `web_api.go`, vocabulary helpers in `bot.go`; `/api/vocabulary` keeps V1-style fast list loading and does not call AI per card.
- Password change: `web_api.go` authenticated `/api/auth/password`
- Referrals: `referral.go`, `premium_i18n.go`, web DTOs in `web_api.go`
- Awards/ranks: localized copy surfaced through session/progress data and browser fallbacks
- Premium crypto payment requisites and payment result dialogs: `crypto_payments.go`, `web_api.go`, and the V2 payment modal/session polling in `web-react/src/App.tsx`
- Runtime prompt overrides: `runtime_prompts.go`, `app_prompts.json`, and prompt call sites in `prompts.go`, `translator.go`, `shadowing.go`, `pronunciation.go`, and `openrouter.go`. `ROLEPLAY_TOOL_V2` requests are routed through `roleplayPrompt` instead of the generic practice prompt.
- Lesson anti-repeat state: `lesson_history` is stored in `storage.go` and `sqlite_store.go`, trimmed by `trimLessonHistory` in `prompts.go`, and passed into `lessonPrompt` from Telegram and Web API lesson start paths. Practice context is trimmed to `practiceMemoryLimit = 5`.
- Legacy/V2 awards: V2 canonical award names/stories live in `web-react/src/App.tsx`; V1 uses matching browser fallback arrays in `web/index.html`, with clean Cyrillic award text allowed through the safe display filter.
- Daily bonus claims: `web_api.go` exposes `/api/daily/claim`; `storage.go` and `sqlite_store.go` persist one claim per date and add XP once.
- Bug report storage: `web_api.go` exposes `/api/bug-report`; runtime reports are appended to `bug_reports/bug_reports.jsonl` with optional screenshots in `bug_reports/screenshots`.
- V2 generated assets added in this pass: `web/assets/icon-roleplay-light.png`, `web/assets/icon-roleplay-dark.png`, `web/assets/icon-pronunciation-light.png`, `web/assets/icon-pronunciation-dark.png`, `web/assets/icon-offline-light.png`, `web/assets/icon-offline-dark.png`, `web/assets/icon-phrasebook-light.png`, and `web/assets/icon-phrasebook-dark.png`.

## Verification Index

Use these commands for this project:

```powershell
npm --prefix web-react run build
go test ./...
node tools/check_encoding_artifacts.mjs
npm --prefix web-react run e2e
npx nx show projects --json
npx nx graph --file=tmp\nx\project-graph.html
```

Current verification from the 2026-05-26 V2 phrasebook / habit / bug-report pass:

- `npm --prefix web-react run build` passed.
- `go test ./...` passed.
- `node tools/check_encoding_artifacts.mjs` passed.
- `npm --prefix web-react run e2e` passed: 8 Playwright tests across desktop and mobile Chromium.
- Nx graph and metadata were generated under `tmp/nx`.
- Latest local V2 build assets: `web/v2/assets/index-WWostreF.css` and `web/v2/assets/index-Bf2gn7t3.js`.
- Authenticated live QA remains pending unless the current Codex browser session has a real `/app/v2` account session.
