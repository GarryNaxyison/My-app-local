# NERIVA Android App — Design Spec

**Date:** 2026-07-11
**Target:** Native Android app (`.aab` for Google Play) reproducing the full NERIVA web app functionality.
**Stack:** Kotlin + Jetpack Compose + Material 3 + Retrofit + Room.

## 1. Product Summary

NERIVA is an AI language tutor operating across web, Telegram, and (now) Android — sharing one user profile, progress, and premium status. The Android app connects to the existing production API at `api.neriva.ru` and reproduces all 28 web views.

## 2. Architecture

**Pattern:** MVVM + Repository + Clean Architecture layers

```
ui/          → Compose screens + ViewModels + Navigation
domain/      → Use cases, business models
data/        → Retrofit API, Room DB, Repositories
di/          → Dependency injection (manual service locator)
```

## 3. Design Tokens (from REDESIGN_PLAN.md + app.css)

### Colors — Light
| Token | Hex | Usage |
|-------|-----|-------|
| bg | `#EDF7FF` | App background |
| surface | `#FFFFFF` (86% opacity) | Panels/cards |
| text | `#07111F` | Primary text |
| muted | `#5C6C82` | Secondary copy |
| faint | `#8391A5` | Tertiary |
| primary (Indigo) | `#2D5BFF` | Primary actions |
| teal | `#00A88F` | AI/practice/voice |
| gold | `#F3B84B` | Premium/awards |
| rose | `#E45A6F` | Errors/mistakes |
| plum | `#7C5CFF` | Secondary accent |

### Colors — Dark
| Token | Hex | Usage |
|-------|-----|-------|
| bg | `#020511` | App background |
| surface | `rgba(13,22,39,0.74)` | Panels |
| text | `#F6F9FF` | Primary text |
| muted | `#A9B7CD` | Secondary |
| primary | `#7AA2FF` | Primary actions |
| teal | `#45D8BB` | AI/practice |
| gold | `#F4C96D` | Premium |
| rose | `#FF8395` | Errors |

### Typography
- Font: Inter (system fallback)
- No viewport-based scaling; compact functional headings

### Shape
- Primary radius: 10dp
- Compact controls: 8dp
- Input border-radius: 16dp

### Shadows
- Card shadow: soft elevation
- Primary action: prominent elevation

## 4. Screen Map (28 views → Compose screens)

| # | Screen | Route | Key API |
|---|--------|-------|---------|
| 1 | Auth (SignIn/Register) | `auth` | `/api/auth/login`, `/register` |
| 2 | Home (Today) | `home` | `/api/session`, `/api/daily/claim` |
| 3 | AI Tutor Lesson | `tutor` | `/api/ai-tutor/*` |
| 4 | Practice (Chat) | `practice` | `/api/practice` |
| 5 | Roleplay | `roleplay` | `/api/practice` (scenario) |
| 6 | Shadowing | `shadowing` | `/api/shadowing/*` |
| 7 | Pronunciation | `pronunciation` | `/api/pronunciation/*` |
| 8 | Learn Words | `words` | `/api/words/next`, `/answer` |
| 9 | Word Game | `word-game` | `/api/word-game/*` |
| 10 | Spelling | `spelling` | `/api/spelling/*` |
| 11 | Vocabulary | `vocabulary` | `/api/vocabulary` |
| 12 | Phrasebook (Notes) | `phrasebook` | `/api/phrasebook` |
| 13 | Offline Decks | `offline` | Local Room |
| 14 | Level Test | `level` | `/api/level-test/*` |
| 15 | Progress/Dashboard | `progress` | `/api/progress` |
| 16 | Awards | `awards` | `/api/progress` |
| 17 | Leaderboard | `leaderboard` | `/api/leaderboard` |
| 18 | Limits | `limits` | `/api/session` |
| 19 | Mistakes | `mistakes` | `/api/mistakes`, `/practice/*` |
| 20 | Tools | `tools` | `/api/tools/*` |
| 21 | Premium | `premium` | `/api/premium/*` |
| 22 | Settings | `settings` | `/api/settings`, `/auth/password` |
| 23 | Referral | `referral` | `/api/session` |
| 24 | Bug Report | (dialog) | `/api/bug-report` |
| 25 | Onboarding | `onboarding` | `/api/settings` |
| 26 | Navigation Edit | (sheet) | `/api/navigation-layout` |
| 27 | Lesson History | `lesson-history` | `/api/ai-tutor/completed` |
| 28 | Profile | `profile` | `/api/auth/profile` |

## 5. Navigation

- **Bottom nav** (mobile-first, matching web mobile pattern):
  - Pinned tabs: Home, Lesson, Practice, Words, More
  - "More" opens grid of remaining views
  - Editable: long-press to reorder/pin
- **Top bar**: language selector, theme toggle, bug report icon
- Navigation state persisted via `/api/navigation-layout`

## 6. API Integration

### Auth
- Session cookie: `poliglot_web_session` = `userID.HMAC-SHA256-signature`
- Cookie managed by OkHttp `CookieJar` (persistent via EncryptedSharedPreferences)
- Auto-login: on launch, call `GET /api/session` → if authenticated, go to Home; else Auth screen
- CORS: backend allows requests with no `Origin` header (Android sends none) ✓
- Captcha (Cloudflare Turnstile): embedded WebView for login/register captcha verification

### Rate Limits
- 240 req/min per IP, 180 req/min per user
- Client-side throttle not needed (within limits for single user)

### Endpoints (57 total)
All documented in `web_api.go:register()` (lines 146-224). Key groups:
- Auth: register, login, telegram/start|status|merge, profile, password, logout
- Session: session, settings, navigation-layout, daily/claim
- Learning: ai-tutor/*, lesson/*, practice, shadowing/*, pronunciation/*
- Vocabulary: words/*, word-game/*, spelling/*, vocabulary
- Mistakes: mistakes, mistakes/practice/*
- Phrasebook: phrasebook
- Tools: tools/voice-text, image-translate, translator, translator-speech
- Progress: progress, premium/*, leaderboard, bug-report

## 7. Data Models

Kotlin data classes mirroring `web-react/src/lib/types.ts`:
- `SessionData`, `UserProfile`, `AiTutorResponse`, `AiTutorStep`
- `WordChallenge`, `SpellingChallenge`, `VocabularyItem`
- `MistakeItem`, `PhrasebookItem`, `PremiumPlan`
- `LeaderboardEntry`, `TranslatorResult`, `LevelQuestion`

Serialization: `kotlinx.serialization` with `@SerialName` for snake_case JSON keys.

## 8. Offline Support

- **Room DB** caches: phrasebook, offline decks, pronunciation history, vocabulary pages
- Downloadable decks (JSON) stored locally for practice without network
- Session cached for instant launch, refreshed in background

## 9. Audio

- **TTS playback**: backend returns WAV blobs → `MediaPlayer`
- **Voice input**: `MediaRecorder` (AAC) → upload as multipart to `/api/practice` or `/api/pronunciation/check`
- **Pronunciation**: record → send → receive score + weak words feedback

## 10. i18n

- 35 interface languages
- Copy strings delivered in `/api/session` response (`copy`, `tool_copy`, `system_copy`, `premium_copy`)
- Android `strings.xml` fallback for `en`/`ru`
- Language switch via `/api/settings` → refreshes session copy

## 11. Security

- `EncryptedSharedPreferences`: cookie persistence, any cached user data
- `network_security_config.xml`: HTTPS-only, optional certificate pinning for `api.neriva.ru`
- No hardcoded secrets; API uses cookie session

## 12. Build & Signing

- **Min SDK:** 26 (Android 8.0) — covers ~95% of devices
- **Target SDK:** 35
- **Build variants:**
  - `debug`: connects to `api.neriva.ru`, uses debug keystore
  - `release`: connects to `api.neriva.ru`, R8 minified, signed AAB
- **R8/ProGuard rules:** keep Retrofit interfaces, Room entities, kotlinx.serialization, Compose
- **Output:** `app/build/outputs/bundle/release/app-release.aab`
- **App ID:** `ru.neriva.app`

## 13. Definition of Done

- All 28 screens functional against production API
- Light/dark theme matching web design tokens
- At least ru/en interface languages working
- Audio playback + voice input functional
- Offline phrasebook/decks accessible
- `.aab` builds successfully
- App installs and runs on Android 8.0+ emulator
