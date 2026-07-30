# Android /app Parity Matrix

Status reflects the native Android implementation in this repository as of
Release 1 foundation. `Manual/device` means that the route is registered but
still needs an authenticated device run against the production-compatible API.

| `/app` view | Android route | Shared API / state | Verification status |
| --- | --- | --- | --- |
| Sign in, registration | `auth` | `/api/auth/*`, session cookie | Unit: API base URL; manual/device pending |
| Onboarding, settings | `onboarding`, `settings` | `/api/settings`, navigation layout | Interface selector limited to EN/RU; manual/device pending |
| Home | `home` | `/api/session`, progress | Manual/device pending |
| AI Tutor | `tutor` | `/api/ai-tutor/*` | Manual/device pending |
| Lesson | `lesson` | `/api/lesson/*` | Manual/device pending |
| Practice | `practice` | `/api/practice` | Manual/device pending |
| Roleplay | `roleplay`, `roleplay-scenarios` | roleplay API flow | Manual/device pending |
| Shadowing | `shadowing` | `/api/shadowing/*` | Manual/device pending |
| Pronunciation | `pronunciation` | `/api/pronunciation/*` | Manual/device pending |
| Words, word game, spelling | `words`, `word-game`, `spelling` | `/api/words/*`, `/api/word-game/*`, `/api/spelling/*` | Manual/device pending |
| Vocabulary | `vocabulary` | `/api/vocabulary` | Manual/device pending |
| Phrasebook | `phrasebook` | `GET`, `POST`, `DELETE /api/phrasebook` | Unit: refresh/save/delete; debug APK built |
| Offline decks | `offline` | Room phrasebook cache | Manual/device pending |
| Level test | `level` | `/api/level-test/*` | Manual/device pending |
| Progress, dashboard, awards | `progress`, `dashboard`, `awards` | `/api/progress`, session state | Manual/device pending |
| Leaderboard, limits | `leaderboard`, `limits` | `/api/leaderboard`, session state | Manual/device pending |
| Mistakes | `mistakes`, `mistake-practice` | `/api/mistakes/*` | Manual/device pending |
| Tools | `tools` | translator, speech, image APIs | Manual/device pending |
| Premium, referral | `premium`, `referral` | plans, payment, activation APIs | Manual/device pending |
| Bug reports | `bug-report` | `/api/bug-report` | Manual/device pending |
| Completed lessons | `completed-lessons` | `/api/ai-tutor/completed` | Manual/device pending |

## Release 1 Evidence

- `NerivaApiClientTest`: normalizes the Retrofit base URL.
- `PhrasebookRepositoryTest`: server refresh, save, and delete replace the
  Room cache.
- `LanguageManagerTest`: only the supported interface locales, Russian and
  English, are exposed in settings.
- `go test . -run 'TestWebPhrasebook(GetReturnsSavedItems|PersistsInStoreAndSession)'`:
  the server serves the phrasebook GET contract used by Android.

## Remaining Release Gates

- Run authenticated route, recording, camera, upload, payment handoff, and
  relaunch smoke checks on an emulator or physical Android device.
- Add Compose UI coverage for success, loading, empty, and API-error states
  for each route before declaring full web parity.
- Complete resource-backed UI translations before enabling additional
  interface locales beyond EN/RU.
