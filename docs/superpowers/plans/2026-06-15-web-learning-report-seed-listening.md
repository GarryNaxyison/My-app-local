# Web Learning Report, Seed, and Listening Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix the web learning flows so ordinary Words can report bad words, completed Lessons expose a visible next-lesson action, generated AI texts vary every run, and Listening differs clearly from Pronunciation.

**Architecture:** Keep the existing Go API and React v2 shell. Reuse the AI Tutor word-report moderation storage and Telegram admin buttons where possible, but route ordinary Words reports through a vocabulary-only correction path that updates SQLite vocabulary tables without requiring an AI Tutor lesson payload.

**Tech Stack:** Go `net/http`, SQLite via `modernc.org/sqlite`, React/Vite, Playwright, GitHub remote `origin`.

---

### Task 1: Spec Commit

**Files:**
- Create: `docs/superpowers/plans/2026-06-15-web-learning-report-seed-listening.md`

- [ ] Save this plan.
- [ ] Commit with `docs: add web learning report seed plan`.

### Task 2: Ordinary Words Report

**Files:**
- Modify: `web_api.go`
- Modify: `ai_tutor_word_reports.go`
- Modify: `vocabulary_sqlite.go`
- Modify: `web-react/src/lib/types.ts`
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Test: `web_api_feature_test.go`
- Test: `telegram_test.go`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [x] Add a failing backend test for `POST /api/words/report`.
- [x] Add a failing Telegram moderation test proving Accept/Fix on a vocabulary report updates SQLite vocabulary tables.
- [x] Add a failing Playwright test proving `.context-display--words` shows `Сообщить об ошибке` and posts the current `word_id`.
- [x] Return `word_id`, `word`, `translation`, and `reportable` from `/api/words/next`.
- [x] Add `/api/words/report` with validation for active word, proposed word, proposed translation, and optional comment.
- [x] Extend report correction logic so stage `vocabulary_word` applies only to SQLite vocabulary tables.
- [x] Reuse the AI Tutor admin keyboard and adjust Telegram report text for vocabulary reports.
- [x] Add the Words report button and dialog in React.
- [x] Verify targeted Go tests and Playwright test.
- [x] Commit with `feat: add word trainer report flow`.

### Task 3: Visible Next Lesson CTA

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [x] Add a failing Playwright assertion that `Следующий урок` is visible inside `.context-display--lesson` after a lesson answer.
- [x] Render an in-panel CTA when `lessonHasCompletedAnswer` is true.
- [x] Keep the existing composer button as a fallback.
- [x] Verify desktop and mobile Playwright coverage for the lesson flow.
- [x] Commit with `fix: show next lesson cta in lesson panel`.

### Task 4: AI Variation Seed

**Files:**
- Modify: `prompts.go`
- Modify: `prompts_test.go`
- Modify if needed: `course_tutor.go`
- Test: `prompts_test.go`
- Test if changed: `course_tutor_test.go`

- [x] Add failing tests that two lesson prompt builds with different counters include different `Variation seed` values.
- [x] Add failing tests that two practice prompt builds with different counters include different `Variation seed` values.
- [x] Add a prompt helper that derives stable-visible variation text from scope, level, count, focus, and recent context.
- [x] Include the seed in lesson and practice prompts with an instruction not to reveal it to the learner.
- [x] Inspect roleplay/dialogue prompt builders and add the same seed rule where AI generates dialogue text.
- [x] Verify prompt tests and affected tutor tests.
- [ ] Commit with `fix: add variation seeds to ai prompts`.

### Task 5: Listening vs Pronunciation UX

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [x] Add a failing Playwright test that Listening hides the target text before the learner answers.
- [x] Keep the Listening audio button visible before answer.
- [x] Show the target text only after the answer is checked.
- [x] Keep Pronunciation target text visible before recording.
- [x] Verify desktop/mobile smoke around Listening and Pronunciation.
- [ ] Commit with `fix: separate listening and pronunciation flows`.

### Task 6: Unified Notes Strip

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] Add a failing Playwright test for AI Tutor feedback that includes both recommendations and mistakes.
- [ ] Assert the UI renders exactly one notes block in the tutor context after feedback.
- [ ] Assert the UI does not show duplicate `Сохранить в заметки` blocks.
- [ ] Assert the UI never renders `Ошибки: Сохранить в заметки`.
- [ ] Merge recommendation and mistake phrase candidates before rendering the note strip.
- [ ] Keep real mistakes visible only in the mistakes block, without turning the note-strip label into an error heading.
- [ ] Verify the AI Tutor feedback smoke test.
- [ ] Commit with `fix: unify ai tutor note strip`.

### Task 7: Final Verification and Deploy

**Files:**
- No source edits unless verification finds a defect.

- [ ] Run targeted Go tests for word report, prompt seed, and lesson API.
- [ ] Run `npm --prefix web-react run build`.
- [ ] Run targeted Playwright tests for Words, Lesson, Listening, and Pronunciation.
- [ ] Run encoding artifact check.
- [ ] Push all commits to `origin/codex/ai-tutor-rebuild-fix`.
- [ ] Deploy backend and React assets to production.
- [ ] Verify production `/healthz`, `/app/`, service status, and new frontend asset markers.
