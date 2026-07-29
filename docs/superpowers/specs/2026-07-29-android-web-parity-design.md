# Android Web Parity Design

## Scope

Bring the native Android application to functional parity with `/app`. The
Android app remains native Jetpack Compose; it does not embed the web app.

## Bootstrap Requirement

Android must start before any parity work can be verified. The Retrofit base
URL is normalized to a trailing slash and startup has a focused test that
constructs the API client without throwing.

## Shared Data Contract

The Go API remains the source of truth for authenticated user data. Android
uses the same server endpoints and session cookie as `/app`.

Phrasebook items must be fetched from, created on, and deleted through the
server API. Room may cache successful server results for offline display, but
must not become a separate unsynchronized phrasebook.

## Parity Groups

Implementation and verification are organized into independent groups:

1. Authentication, session restoration, onboarding, settings, theme, and
   interface/learning language selection.
2. Learning flows: AI Tutor, lessons, practice, roleplay, shadowing,
   pronunciation, words, word game, spelling, vocabulary, mistakes, and
   phrasebook.
3. Service flows: tools, image/voice upload, progress, awards, leaderboard,
   referral, premium plans, payment handoff, activation keys, bug reports,
   offline decks, level test, and limits.
4. Localization and quality: all supported Android resource locales contain
   product text rather than accidental English fallback, and release flows have
   automated API, UI, and device-level coverage.

## Native Interaction Rules

Each Android screen may use mobile-native layout and controls, but must expose
the same completed user outcome, server state changes, permissions, errors,
and access restrictions as its `/app` counterpart. Android-only enhancements
are allowed only when they do not replace or weaken the shared flow.

Audio and image features upload through the same API contract as web. Payment
screens hand off to the server-provided payment URL and return to a refreshed
session state; card details are never handled in the app.

## Testing

- Unit tests cover API URL construction, cookie persistence, repository sync,
  and serialization contracts.
- Compose UI tests cover each screen's success, loading, empty, and server
  error state.
- Instrumented smoke tests run the authenticated core flow on an emulator:
  sign in, open a learning flow, submit an answer/recording fixture, save a
  phrasebook item, and confirm session state after relaunch.
- A parity checklist maps every `/app` view and API action to an Android
  screen and test before release.

## Release Criteria

- Debug and release APKs build successfully.
- Android starts with no Retrofit initialization exception.
- Shared phrasebook is consistent between web and Android after refresh.
- Each parity group has passing automated coverage and documented manual
  device verification.
- No unsupported interface language silently falls back to unrelated product
  copy.
