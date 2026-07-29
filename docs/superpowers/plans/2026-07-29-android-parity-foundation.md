# Android Parity Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the Android app boot reliably and use server-synchronized phrasebook data as the base for native web parity.

**Architecture:** Retrofit client creation is tested as an isolated factory contract. Phrasebook operations are routed through the existing `/api/phrasebook` contract, with Room as a cache populated from server data rather than an independent source of truth.

**Tech Stack:** Kotlin, Jetpack Compose, Retrofit, OkHttp, Room, JUnit, MockWebServer.

---

### Task 1: Add Android Test Support And API Client Regression Test

**Files:**
- Modify: `android-app/gradle/libs.versions.toml`
- Modify: `android-app/app/build.gradle.kts`
- Create: `android-app/app/src/test/java/ru/neriva/app/data/api/NerivaApiClientTest.kt`
- Modify: `android-app/app/src/main/java/ru/neriva/app/data/api/NerivaApiClient.kt`
- Modify: `android-app/app/src/main/java/ru/neriva/app/NERIVAApp.kt`

- [ ] Add JUnit and MockWebServer test dependencies.
- [ ] Write a failing test that `NerivaApiClient.create("https://api.neriva.ru", jar)` does not throw and resolves `api/session` under the host root.
- [ ] Run `gradlew.bat :app:testDebugUnitTest` and confirm the test fails because Retrofit rejects the non-terminated base URL.
- [ ] Normalize the base URL in one client factory helper, including a trailing slash, and use it from `NERIVAApp`.
- [ ] Re-run the unit test and `:app:assembleDebug`.

### Task 2: Synchronize Phrasebook Through The Server Contract

**Files:**
- Modify: `android-app/app/src/main/java/ru/neriva/app/data/repo/Repositories.kt`
- Modify: `android-app/app/src/main/java/ru/neriva/app/data/db/Daos.kt`
- Modify: `android-app/app/src/main/java/ru/neriva/app/ui/screens/Screens.kt`
- Create: `android-app/app/src/test/java/ru/neriva/app/data/repo/PhrasebookRepositoryTest.kt`

- [ ] Write failing repository tests for refresh, save, and delete using a fake `NerivaApi`.
- [ ] Add repository methods that fetch server items, replace the Room cache, save or remove server items, then refresh cache.
- [ ] Change the phrasebook screen to load from the repository and render cache plus loading/error state.
- [ ] Run repository tests and `:app:assembleDebug`.

### Task 3: Establish The Parity Test Matrix

**Files:**
- Create: `android-app/PARITY_MATRIX.md`
- Create: `android-app/app/src/androidTest/java/ru/neriva/app/NavigationSmokeTest.kt`

- [ ] Map every `/app` view and action to an Android route, API call, and verification state.
- [ ] Add Compose navigation smoke coverage for all registered routes and authenticated-route session errors.
- [ ] Run `gradlew.bat :app:connectedDebugAndroidTest` on an available emulator; report emulator absence as a verification gap rather than a passed test.
