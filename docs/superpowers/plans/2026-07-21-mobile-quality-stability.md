# Mobile Quality Stability Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use `executing-plans` to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Remove mobile workspace overlaps, retain pronunciation advice, make server audio retryable, and cache static assets after the first visit.

**Architecture:** Retain the existing React, Go, and service-worker boundaries. Change only the report identity and focused mobile render/CSS paths; wrap current TTS calls with one bounded retry; cache only static `/app/` resources.

**Tech Stack:** React 19, TypeScript, Vite, Playwright, Go `net/http`, OpenRouter TTS, Cache Storage.

---

### Task 1: Retain distinct pronunciation advice

**Files:**
- Modify: `web-react/src/App.tsx:1244-1275, 6774-6824, 7492-7545`
- Modify: `web-react/e2e/web-smoke.spec.ts:2419-2445, 4060-4077`

- [ ] Write this failing regression:

```ts
localStorage.setItem("poliglot-pronunciation-v2:demor22", JSON.stringify([
  { expected: "Can we meet earlier?", spoken: "Can we meet earlier?", score: 70, tips: ["Keep the final consonant clear."], createdAt: "2026-07-21T10:00:00.000Z" },
  { expected: "Can we meet earlier?", spoken: "Can we meet earlier?", score: 70, tips: ["Link meet and earlier smoothly."], createdAt: "2026-07-21T10:01:00.000Z" },
]));
await page.goto("/app/?view=pronunciation");
await expect(page.locator(".pronunciation-history-v2__rows > div")).toHaveCount(2);
await expect(page.locator(".pronunciation-history-v2")).toContainText("Link meet and earlier smoothly.");
```

- [ ] Run: `npx playwright test e2e/web-smoke.spec.ts -g "keeps distinct tips" --reporter=list`; expect failure because the old fingerprint omits tips.
- [ ] Change `pronunciationFingerprint` to serialize expected/spoken/score/feedback/problem words/phoneme issues/stress/rhythm/intonation/tips. Keep the existing sixteen-report cap and compact-only tip limit.
- [ ] Re-run: `npx playwright test e2e/web-smoke.spec.ts -g "keeps distinct tips|pronunciation map is compact" --reporter=list`; expect PASS.
- [ ] Commit: `git add web-react/src/App.tsx web-react/e2e/web-smoke.spec.ts && git commit -m "fix: retain distinct pronunciation advice"`.

### Task 2: Stabilize mobile listening, tools, and bottom navigation

**Files:**
- Modify: `web-react/src/App.tsx:7180-7289, 7472-7486, 8727-8849`
- Modify: `web-react/src/styles/app.css:9033-9059, 9978-10546`
- Modify: `web-react/e2e/web-smoke.spec.ts:2578-2606, 2735-2777, 4449-4467`

- [ ] Write a mobile regression that navigates to shadowing, asserts exactly one `.audio-wave-button-v2`, then navigates to tools, sends `Bonjour`, and checks that `.tools-submit-v2.bottom <= source-select.top` and `source-select.bottom <= target-select.top`.
- [ ] Run: `npx playwright test e2e/web-smoke.spec.ts -g "mobile listening and populated translator controls" --reporter=list`; expect failure before the final mobile rules.
- [ ] Add `compact?: boolean` to `AudioActionRow`; pass it to `AudioWaveButton` for the active shadowing model. Use `text + wordId` as the player key. In the final mobile CSS block use a normal-flow single-column grid for `.tool-input-shell-v2` and `.language-row-v2`; remove old textarea padding that served absolute submit buttons; keep `.mobile-bottom-nav-v2` in grid row 2 with no fixed/absolute positioning.
- [ ] Re-run: `npx playwright test e2e/web-smoke.spec.ts -g "mobile listening and populated translator controls|mobile pronunciation blocks keep vertical order|mobile leaderboard stays readable" --reporter=list`; expect PASS.
- [ ] Commit: `git add web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts && git commit -m "fix: stabilize mobile listening and tools layout"`.

### Task 3: Retry server audio; never use browser TTS

**Files:**
- Modify: `web_api.go:2526-2575, 3573-3631`
- Modify: `web_api_feature_test.go:815-930`
- Modify: `web-react/src/components/ui/audio-wave-button.tsx:64-121`
- Modify: `web-react/e2e/web-smoke.spec.ts:4034-4058`

- [ ] Write `TestWebTranslatorSpeechRetriesProviderOnceThenReturnsAudio`: arrange a premium user and an HTTP transport that returns a temporary error once then WAV bytes; call `/api/tools/translator-speech`; assert HTTP 200, two provider calls, non-empty body, and `audio/wav`.
- [ ] Run: `go test ./... -run TestWebTranslatorSpeechRetriesProviderOnceThenReturnsAudio -count=1`; expect failure because the first provider error is returned directly.
- [ ] Add and use this shared helper after existing cache lookups:

```go
func (api *webAPI) synthesizeSpeechWithRetry(ctx context.Context, text string) ([]byte, error) {
  var lastErr error
  for attempt := 0; attempt < 2; attempt++ {
    audio, err := api.bot.openrouter.synthesizeSpeech(ctx, api.cfg.OpenRouterTTSModel, api.cfg.OpenRouterTTSVoice, text)
    if err == nil && len(audio) > 0 { return audio, nil }
    lastErr = err
  }
  return nil, lastErr
}
```

- [ ] In `AudioWaveButton`, change the error label to a localized retry action while retaining the error in `title`; the existing click repeats the server request. Do not introduce `speechSynthesis`.
- [ ] Re-run: `go test ./... -run "TestWebTranslatorSpeechRetriesProviderOnceThenReturnsAudio|TestWriteAudioResponsesUseDetectedWAVMetadata" -count=1`; then run a Playwright route test with 502 then success and assert the same player retries.
- [ ] Commit: `git add web_api.go web_api_feature_test.go web-react/src/components/ui/audio-wave-button.tsx web-react/e2e/web-smoke.spec.ts && git commit -m "fix: retry server audio generation"`.

### Task 4: Cache shell and static assets, never API data

**Files:**
- Modify: `web-react/public/offline-deck-sw.js`
- Modify: `web-react/e2e/web-smoke.spec.ts:777-780`

- [ ] Write a source-policy regression asserting the worker: excludes `url.pathname.startsWith("/api/")`, cache-first serves static assets, cache version changes, and navigation falls back to cached `/app/`.
- [ ] Run: `npx playwright test e2e/web-smoke.spec.ts -g "PWA caches static app assets" --reporter=list`; expect failure against the network-first worker.
- [ ] Bump `CACHE_NAME`. On install cache `/app/`, boot logo, and same-origin JS/CSS/image URLs parsed from the fetched app document. Cache-first `GET /app/assets/` and static JS/CSS/image/font requests, writing successful fetches; use network-first only for navigation and never call `event.respondWith` for `/api/`.
- [ ] Re-run: `npx playwright test e2e/web-smoke.spec.ts -g "PWA caches static app assets|PWA service worker cache" --reporter=list && npm run build`; expect PASS and build exit 0.
- [ ] Commit: `git add web-react/public/offline-deck-sw.js web-react/e2e/web-smoke.spec.ts web && git commit -m "perf: cache static app shell assets"`.

### Task 5: Verify and release

**Files:**
- Verify: `web_api_feature_test.go`, `web-react/e2e/web-smoke.spec.ts`

- [ ] Run: `go test ./...`; expect PASS.
- [ ] Run: `npm run build`; expect TypeScript and Vite exit 0.
- [ ] Run: `npx playwright test e2e/web-smoke.spec.ts -g "pronunciation|mobile listening|mobile leaderboard|audio buttons|PWA caches static app assets" --reporter=list`; expect PASS.
- [ ] Run: `git diff --check && git status --short`; expect no whitespace errors and only task files plus generated `web` assets.
- [ ] Push: `git push origin HEAD`.
