# Telegram Listening and Pronunciation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Expose current Listening and Pronunciation practice in Telegram and keep their logic distinct from each other.

**Architecture:** Keep Listening on the existing `shadowing` flow: the learner listens first, sends voice, and the bot compares what was heard to the target. Add a separate Telegram Pronunciation exact-repeat flow that reuses the existing pronunciation phrase builder, TTS audio, transcription, and `pronunciationModeExact` assessor.

**Tech Stack:** Go Telegram bot handlers, existing `shadowing.go` pronunciation helpers, localized `uiCopy`, existing HTTP-backed Telegram test harness.

---

### Task 1: Telegram Menu Coverage

**Files:**
- Modify: `telegram_test.go`
- Modify: `telegram.go`
- Modify: `i18n.go`

- [ ] Add a failing test that the main menu exposes both `menu_shadowing` and `menu_pronunciation`.
- [ ] Add a failing localization test that every interface language has a non-empty Pronunciation label and the menu button uses that label.
- [ ] Add `Pronunciation` to `uiCopy` and runtime fallbacks.
- [ ] Add Pronunciation buttons to the main and learning Telegram menus.
- [ ] Run the targeted menu tests until they pass.
- [ ] Commit with `feat: expose pronunciation in telegram menus`.

### Task 2: Telegram Pronunciation Flow

**Files:**
- Modify: `telegram_test.go`
- Modify: `shadowing.go`
- Modify: `bot.go`

- [ ] Add a failing test that `menu_pronunciation` starts a standalone exact-repeat pronunciation mode, sends the target phrase, and requests a voice reply.
- [ ] Add `pronunciationMode(phrase)` / `parsePronunciationMode(mode)` helpers using the same safe base64 format as shadowing.
- [ ] Add `startPronunciation` that checks premium voice access, builds a phrase with `buildPronunciationPhrase`, stores the pronunciation mode, sends localized instructions, and sends model audio.
- [ ] Route `menu_pronunciation` callbacks to `startPronunciation`.
- [ ] Route voice messages in pronunciation mode through `buildPronunciationAssessment(..., pronunciationModeExact)` and show the localized pronunciation result with a next-pronunciation button.
- [ ] Run targeted Telegram, shadowing, and pronunciation tests until they pass.
- [ ] Commit with `feat: add telegram pronunciation practice`.

### Task 3: Verification and Deploy

**Files:**
- No source edits unless verification finds a defect.

- [ ] Run targeted Go tests for Telegram menus, pronunciation, shadowing, and existing AI Tutor/web regressions.
- [ ] Run `npm --prefix web-react run build`.
- [ ] Run targeted Playwright web smoke tests for lesson, listening, pronunciation, words, and AI Tutor.
- [ ] Run the encoding artifact check.
- [ ] Push all commits to `origin/codex/ai-tutor-rebuild-fix`.
- [ ] Deploy backend and web assets to production.
- [ ] Verify `/healthz`, `/app/`, Telegram service status, and production asset/runtime markers.
