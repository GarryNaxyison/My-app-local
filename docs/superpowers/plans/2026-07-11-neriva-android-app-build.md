# NERIVA Android App — Build Plan

**Date:** 2026-07-11
**Spec:** `docs/superpowers/specs/2026-07-11-neriva-android-app-design.md`

## Phases

### Phase 1: Environment + Scaffold
1. Set `ANDROID_HOME`, `JAVA_HOME`
2. Create `android-app/` Gradle project (package `ru.neriva.app`)
3. Configure dependencies: Compose BOM, Material 3, Retrofit, OkHttp, Room, Coil, kotlinx.serialization, EncryptedSharedPreferences
4. Create AVD for testing

### Phase 2: Data Layer
1. NERIVA theme (Color, Type, Shape tokens)
2. API interface (Retrofit `NerivaApi` — 57 endpoints)
3. Data models (kotlinx.serialization)
4. Repositories (Auth, Session, Tutor, Vocabulary, Pronunciation, Mistakes, Premium, Tools, Progress)
5. Room DB (Phrasebook, OfflineDeck, PronunciationHistory DAOs)
6. CookieJar (persistent)

### Phase 3: Auth + Navigation
1. SignIn / Register screens with Turnstile WebView
2. Telegram OTP login
3. Compose Navigation (NavHost + bottom nav scaffold)
4. Auto-login flow

### Phase 4: Core Screens
1. HomeScreen (Today, daily bonus, habit calendar)
2. TutorScreen (AI lesson flow)
3. PracticeScreen (chat with AI, voice/image)
4. VocabularyScreen (dictionary + Learn Words + Spelling)
5. ProgressScreen (XP, streak, awards)
6. SettingsScreen

### Phase 5: Extended Screens
Roleplay, Shadowing, Pronunciation, WordGame, Phrasebook, Offline, LevelTest, Mistakes, Tools, Premium, Referral, Leaderboard, Limits, Awards, LessonHistory, Profile

### Phase 6: Polish + AAB
1. Animations, splash, app icons
2. R8/ProGuard rules
3. Signing config
4. `./gradlew bundleRelease` → `app-release.aab`

### Phase 7: Test + Deploy
1. Run on emulator
2. Test core flows against production API
3. Commit + push to GitHub
