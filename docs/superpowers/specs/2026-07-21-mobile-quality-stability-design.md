# Mobile quality and instant-repeat loading

**Status:** approved design, awaiting written-spec review
**Date:** 2026-07-21

## Goal

Repair the mobile workspaces shown in the supplied Android screenshots, make pronunciation feedback durable across attempts, require server-backed audio, and make repeat visits load the app shell and navigation assets from cache.

## Scope

### 1. Mobile workspace geometry

The listening/shadowing screen will render one primary, full-width task card. Its spoken-model control will use the compact audio treatment and will not be duplicated in both the transcript and composer for the same active phrase. The response field will use the same available width as other workspace cards.

The translation tools screen will have one deterministic mobile flow:

1. Input text
2. Send button
3. Source language selector
4. Swap control
5. Target language selector
6. Saved-phrase strip

Controls are normal-flow grid items on mobile. No mobile button or label may be positioned over an input, selector, or saved-phrase strip. The content panel uses natural content height; it does not reserve a blank viewport-sized tail.

The page uses one authoritative mobile layout rule set: app grid with a scrolling content row and a separate bottom-navigation row. The content row accounts for safe areas, while the navigation never overlays it.

### 2. Pronunciation history and map

Pronunciation-report identity includes the complete instructional payload: expected/spoken text, score, problem words, phoneme issues, diagnostics, overall feedback, and tips. Only exact duplicate reports collapse.

Persisted reports and live reports merge through this identity, retaining up to the existing history limit. The pronunciation map combines unique word/phoneme findings while the report view preserves every distinct system tip associated with retained attempts.

### 3. Server-backed audio

Audio remains server-generated. Browser speech synthesis is explicitly out of scope.

The client continues to request word pronunciation first and text-to-speech second when a word-specific route cannot serve audio. It must not treat server failure as a passive, permanent `Audio unavailable` state: it exposes a retry action and preserves the server diagnostic for troubleshooting.

The API routes validate non-empty text and target language, synthesize valid WAV or MP3 bytes, return the matching content type, and surface provider failure as a clear HTTP error. A cache of successful server audio responses may be used to make repeated playback reliable and fast; it must be keyed by normalized text, language, configured voice, and configured model, and must not serve a different language or voice.

### 4. Static-asset caching

The service worker precaches the application document, boot logo, currently referenced application JS/CSS chunks, and the mobile-navigation icons after a successful visit. Hashed static assets use cache-first with background refresh; navigation requests use network-first with the cached app shell as offline fallback. API, user content, and authentication responses are never cached by the service worker.

Installing a new worker uses a new cache version and removes only older application caches after activation.

## Non-goals

- No browser-TTS fallback.
- No redesign of desktop screens or navigation information architecture.
- No caching of API data, recordings, authentication, or learner content.
- No wholesale responsive-layout rewrite beyond the affected mobile workspaces.

## Acceptance criteria

- At a 360–430 px viewport, listening, pronunciation, and tools controls do not overlap; content fits within the scroll panel and the final element remains visible above bottom navigation.
- The listening prompt and response field match the content width of the other mobile sections; there is one active spoken-model player for the active phrase.
- Submitted translation results do not overlap source/target labels, selectors, swap control, send button, or saved phrases, and the tools panel has no empty forced-height area.
- Two assessment attempts with otherwise similar data but different tips both remain after reload and their advice is visible in the pronunciation feedback/map.
- Word and translator speech APIs return playable audio with correct content type for valid input. A provider failure is actionable and retryable, never silently represented as permanently unavailable audio.
- A second app visit can render the cached shell and already-visited hashed assets/navigation icons without waiting for each icon request; API responses remain uncached.
- New Playwright regression tests cover the populated translation state, compact listening player, mobile clearance above navigation, pronunciation advice persistence across reload, audio request/failure-retry behavior, and service-worker caching strategy.

## Verification

Run focused Playwright regressions, Go API feature tests covering TTS routes, the full web build, and the relevant mobile E2E suite. Inspect service-worker cache names and response policies in a browser test; verify no `/api/` response is written to Cache Storage.
