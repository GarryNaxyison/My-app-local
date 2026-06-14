# AI Tutor Word Report Admin Fix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Keep this checklist updated as work lands.

**Goal:** Add protected AI Tutor word-correction reports from the word-learning step, route them to Telegram moderation, let the admin accept/reject/fix corrections, apply accepted/fixed words to both the AI Tutor lesson payload and the vocabulary runtime, skip the reported word in the current learning session, and lock pronunciation scoring for free users.

**Architecture:** Extend the existing AI Tutor web API, storage layer, Telegram callback handling, and React Tutor view. The browser sends a small JSON report tied to the current `session_id` and `word_learn_N` stage. The backend validates ownership and rate limits, stores a pending report, notifies Telegram ops recipients with inline actions, and immediately advances the current user past the reported word. Admin actions are idempotent and update the lesson payload plus vocabulary SQLite runtime.

**Tech Stack:** Go backend (`web_api.go`, `bot.go`, `telegram.go`, `sqlite_store.go`, `storage.go`), SQLite storage, existing Telegram Bot API client, React/TypeScript frontend (`web-react/src/App.tsx`), existing Go tests, `npm`/Nx workspace tooling where available.

---

## File Structure

Files expected to change:

```text
ai_tutor_core.go
ai_tutor_engine.go
storage.go
sqlite_store.go
vocabulary_sqlite.go
web_api.go
telegram.go
bot.go
ai_tutor_store_test.go
web_api_feature_test.go
telegram_test.go
web_app_shell_test.go
web-react/src/App.tsx
web-react/src/lib/i18n.ts
web-react/src/styles/app.css
```

Reference spec:

```text
docs/superpowers/specs/2026-06-14-ai-tutor-word-report-admin-fix-design.md
```

---

## Task 1: Add Failing Store Tests for Word Reports

- [ ] Add tests in `ai_tutor_store_test.go` before production code.
- [ ] Cover creating a pending report for a lesson/session word.
- [ ] Cover duplicate pending dedupe for the same user/session/stage.
- [ ] Cover status transition from `pending` to `accepted`/`rejected`/`fixed`.
- [ ] Cover the moderation prompt message id storage used for admin reply handling.

Target test names:

```go
func TestSQLiteStoreAITutorWordReportLifecycle(t *testing.T) {}
func TestSQLiteStoreAITutorWordReportDeduplicatesPendingStage(t *testing.T) {}
```

Expected first run:

```text
go test ./... -run 'TestSQLiteStoreAITutorWordReport'
```

Expected failure before implementation:

```text
store.createAITutorWordReport undefined
```

Implementation notes for tests:

- Use the existing temp sqlite store helpers from nearby AI Tutor store tests.
- Build a minimal lesson payload containing one word and a session on `word_learn_0`.
- Insert the report with proposed word/translation/comment and assert persisted fields.
- Re-submit the same user/session/stage while pending and assert the original report is returned instead of a new row.

---

## Task 2: Implement Store Model and SQLite Persistence

- [ ] Add report status constants and record types in `ai_tutor_core.go`.
- [ ] Add report methods to `storage` in `storage.go`.
- [ ] Implement SQLite schema migration in `sqlite_store.go`.
- [ ] Implement report create/read/status/message helpers in `sqlite_store.go`.
- [ ] Implement minimal `jsonStore` methods so tests and bot helpers that use `storage` still compile.
- [ ] Re-run the Task 1 tests and make them pass.

Data model:

```go
type aiTutorWordReportStatus string

const (
    aiTutorWordReportPending  aiTutorWordReportStatus = "pending"
    aiTutorWordReportAccepted aiTutorWordReportStatus = "accepted"
    aiTutorWordReportRejected aiTutorWordReportStatus = "rejected"
    aiTutorWordReportFixed    aiTutorWordReportStatus = "fixed"
)

type aiTutorWordReportRecord struct {
    ID                  string
    UserID              int64
    SessionID           string
    LessonID            string
    Stage               string
    WordIndex           int
    OriginalWord        string
    OriginalTranslation string
    ProposedWord        string
    ProposedTranslation string
    Comment             string
    Status              aiTutorWordReportStatus
    AdminChatID         string
    AdminMessageID      int64
    FixPromptChatID     string
    FixPromptMessageID  int64
    CreatedAt           time.Time
    UpdatedAt           time.Time
    ResolvedAt          time.Time
}
```

SQLite table:

```sql
CREATE TABLE IF NOT EXISTS ai_tutor_word_reports (
    id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    session_id TEXT NOT NULL,
    lesson_id TEXT NOT NULL,
    stage TEXT NOT NULL,
    word_index INTEGER NOT NULL,
    original_word TEXT NOT NULL,
    original_translation TEXT NOT NULL,
    proposed_word TEXT NOT NULL,
    proposed_translation TEXT NOT NULL,
    comment TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    admin_chat_id TEXT NOT NULL DEFAULT '',
    admin_message_id INTEGER NOT NULL DEFAULT 0,
    fix_prompt_chat_id TEXT NOT NULL DEFAULT '',
    fix_prompt_message_id INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    resolved_at DATETIME
);
```

Indexes:

```sql
CREATE INDEX IF NOT EXISTS idx_ai_tutor_word_reports_user_created
    ON ai_tutor_word_reports(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_ai_tutor_word_reports_session
    ON ai_tutor_word_reports(session_id, created_at);
CREATE INDEX IF NOT EXISTS idx_ai_tutor_word_reports_status
    ON ai_tutor_word_reports(status, created_at);
CREATE UNIQUE INDEX IF NOT EXISTS idx_ai_tutor_word_reports_pending_unique
    ON ai_tutor_word_reports(user_id, session_id, stage)
    WHERE status = 'pending';
```

Storage API shape:

```go
createAITutorWordReport(ctx context.Context, report aiTutorWordReportRecord) (aiTutorWordReportRecord, bool, error)
getAITutorWordReport(ctx context.Context, id string) (aiTutorWordReportRecord, bool, error)
countAITutorWordReports(ctx context.Context, userID int64, sessionID string, since time.Time) (userCount int, sessionCount int, error)
countAITutorWordReportsSince(ctx context.Context, since time.Time) (int, error)
setAITutorWordReportAdminMessage(ctx context.Context, id, chatID string, messageID int64) error
setAITutorWordReportFixPrompt(ctx context.Context, id, chatID string, messageID int64) error
resolveAITutorWordReport(ctx context.Context, id string, status aiTutorWordReportStatus, finalWord, finalTranslation string, now time.Time) (aiTutorWordReportRecord, bool, error)
findPendingAITutorWordReportByFixPrompt(ctx context.Context, chatID string, messageID int64) (aiTutorWordReportRecord, bool, error)
```

The boolean return from `createAITutorWordReport` is `true` when the returned row was already present due to duplicate pending dedupe.

---

## Task 3: Add Failing Backend API Tests

- [ ] Add web API tests in `web_api_feature_test.go` before implementation.
- [ ] Cover successful report submission from current `word_learn_N`.
- [ ] Cover validation failure when reported stage is not the current session stage.
- [ ] Cover rate limits and duplicate report handling.
- [ ] Cover that report submission advances the current stage so the next tutor fetch shows the next learning item.

Target test names:

```go
func TestAITutorWordReportCreatesTelegramNotificationAndSkipsWord(t *testing.T) {}
func TestAITutorWordReportRejectsStaleStage(t *testing.T) {}
func TestAITutorWordReportRateLimitsUser(t *testing.T) {}
```

Request shape:

```json
{
  "session_id": "session-id",
  "stage": "word_learn_0",
  "proposed_word": "correct word",
  "proposed_translation": "correct translation",
  "comment": "short note"
}
```

Assertions:

- `POST /api/ai-tutor/word-report` returns `200`.
- Response contains the refreshed AI Tutor step after the reported word.
- Stored report has original/proposed fields.
- Telegram test server receives an inline keyboard with `ait_word_report|<id>|accept`, `reject`, and `fix`.
- Free/stale/invalid requests do not call Telegram.

Expected first run:

```text
go test ./... -run 'TestAITutorWordReport'
```

Expected failure before implementation:

```text
404 Not Found
```

---

## Task 4: Implement Backend Report Endpoint

- [ ] Register `POST /api/ai-tutor/word-report` in `web_api.go`.
- [ ] Parse small JSON only; reject unknown large/multipart payloads.
- [ ] Authenticate with existing web auth/session helpers.
- [ ] Load the session and lesson, verify the requester owns the session.
- [ ] Require the supplied stage to equal the current session stage.
- [ ] Parse `word_learn_N` and load the corresponding `aiTutorWord`.
- [ ] Trim and validate:
  - proposed word: `1..80`
  - proposed translation: `1..160`
  - comment: `0..500`
- [ ] Enforce anti-abuse before Telegram:
  - user burst: max `3` reports per `10m`
  - user daily: max `10` reports per `24h`
  - session total: max `2` reports
  - global Telegram send burst: max `30` notifications per `10m`
  - duplicate pending same user/session/stage returns existing report and still advances the session once
- [ ] Store the report.
- [ ] Send Telegram ops notification only if the global notification budget allows it.
- [ ] Advance current stage with the existing `aiTutorNextStage` logic so the reported word is skipped.
- [ ] Return the same response envelope as AI Tutor answer/start endpoints for the next step.

Notification text should include:

```text
AI Tutor word report
Report: <id>
User: <user_id>
Stage: <stage>
Original: <word> - <translation>
User suggests: <word> - <translation>
Comment: <comment>
```

Inline keyboard callback data:

```text
ait_word_report|<report_id>|accept
ait_word_report|<report_id>|reject
ait_word_report|<report_id>|fix
```

---

## Task 5: Add Failing Telegram Moderation Tests

- [ ] Add Telegram callback tests in `telegram_test.go` before implementation.
- [ ] Cover admin accept applies the proposed word and translation.
- [ ] Cover admin reject marks the report rejected and does not mutate content.
- [ ] Cover admin fix sends a prompt and admin reply `word - translation` applies that override.
- [ ] Cover non-admin callback/reply is ignored or rejected.
- [ ] Cover duplicate callbacks are idempotent.

Target test names:

```go
func TestTelegramAITutorWordReportAcceptAppliesCorrection(t *testing.T) {}
func TestTelegramAITutorWordReportRejectDoesNotMutateLesson(t *testing.T) {}
func TestTelegramAITutorWordReportFixReplyAppliesCorrection(t *testing.T) {}
func TestTelegramAITutorWordReportAdminOnly(t *testing.T) {}
```

Expected first run:

```text
go test ./... -run 'TestTelegramAITutorWordReport'
```

Expected failure before implementation:

```text
unknown callback
```

---

## Task 6: Implement Telegram Admin Moderation

- [ ] Extend `telegramMessage` in `telegram.go` with `ReplyToMessage`.
- [ ] Add Telegram client helper that sends a text message and returns Telegram `message_id`.
- [ ] Add Telegram client helper that sends an inline keyboard and returns `message_id`, or update existing inline helper safely.
- [ ] In `bot.handleCallbackQuery`, route `ait_word_report|...` callback data.
- [ ] Allow only Telegram ops recipients from existing `telegramOpsRecipients`.
- [ ] For `accept`:
  - resolve report to `accepted`
  - apply `ProposedWord`/`ProposedTranslation`
  - answer callback and send short admin confirmation
- [ ] For `reject`:
  - resolve report to `rejected`
  - do not mutate lesson/vocabulary
  - answer callback and send short admin confirmation
- [ ] For `fix`:
  - send prompt: `Ответьте на это сообщение в формате: word - translation`
  - store the prompt chat id and message id on the report
- [ ] In `bot.handleMessage`, before normal user text handling, detect admin replies to a stored fix prompt.
- [ ] Parse the first supported separator:
  - `" - "`
  - `" — "`
  - `"-"` as a fallback after trim
- [ ] Validate final word/translation with the same caps as user report fields.
- [ ] Resolve report to `fixed` and apply the admin override.

Idempotency:

- Resolving a non-pending report should return the existing resolved row without reapplying mutations.
- Callback answers should still succeed with a short “already resolved” message.

---

## Task 7: Add Failing Correction Application Tests

- [ ] Add focused tests for applying accepted/fixed corrections.
- [ ] Verify both the AI Tutor lesson payload and vocabulary runtime are updated.
- [ ] Verify lesson fingerprint changes after the payload mutation.
- [ ] Verify audio target text follows the corrected learning-language word.

Target test names:

```go
func TestApplyAITutorWordReportCorrectionUpdatesLessonAndVocabulary(t *testing.T) {}
func TestApplyAITutorWordReportCorrectionPreservesWordIDAndExamples(t *testing.T) {}
```

Expected first run:

```text
go test ./... -run 'TestApplyAITutorWordReportCorrection'
```

Expected failure before implementation:

```text
undefined: applyAITutorWordReportCorrection
```

Vocabulary assertions:

- `vocabulary_words.word` equals the final learning-language word.
- `vocabulary_translations.text` for the interface language equals the final translation.
- `vocabulary_ai_words` and `vocabulary_ai_translations` contain the corrected values for AI runtime reads.
- The SQLite vocabulary ID cache is invalidated.

---

## Task 8: Implement Correction Application Across Both Layers

- [ ] Add `applyAITutorWordReportCorrection(ctx, store, report, finalWord, finalTranslation, now)` helper.
- [ ] Load the lesson record by `report.LessonID`.
- [ ] Mutate `lesson.Payload.Words[report.WordIndex]`:
  - preserve `ID`
  - set `Target`
  - set `InterfaceTranslation`
  - set `AudioTextTarget`
  - keep existing examples unless an exact occurrence replacement is safe
- [ ] Recompute lesson fingerprint with the existing lesson fingerprint helper.
- [ ] Save the updated lesson via `saveAITutorLesson`.
- [ ] Update the SQLite vocabulary runtime through a helper in `vocabulary_sqlite.go`.
- [ ] Keep vocabulary helper tolerant of missing rows by creating the AI mutable row where possible.

Suggested helper shape:

```go
func applyAITutorWordReportCorrection(ctx context.Context, st storage, report aiTutorWordReportRecord, finalWord, finalTranslation string, now time.Time) error
```

Vocabulary helper shape:

```go
func sqliteVocabularyApplyAITutorWordCorrection(ctx context.Context, original aiTutorWordReportRecord, finalWord, finalTranslation string) error
```

Runtime update rules:

- Locate `vocabulary_words` by the original word id if the AI Tutor word has a vocabulary id, else by `(language, word)`.
- Update `vocabulary_words.word` and `vocabulary_words.russian` where applicable.
- Upsert `vocabulary_translations` for the interface language.
- Upsert `vocabulary_ai_words` and `vocabulary_ai_translations` so AI-generated mutable content matches.
- Call `invalidateSQLiteVocabularyIDCache()`.

If no vocabulary database is configured, do not fail the Telegram moderation path; the lesson payload update is still authoritative for AI Tutor. Log the missing vocabulary runtime.

---

## Task 9: Add Frontend Source Tests

- [ ] Add source-level frontend tests in `web_app_shell_test.go` before implementation.
- [ ] Assert the Tutor word-learning UI includes the report button, dialog fields, endpoint call, and skip-after-submit flow.
- [ ] Assert copy keys exist for Russian and English.

Target test name:

```go
func TestWebAppAITutorWordReportUIWiring(t *testing.T) {}
```

Expected first run:

```text
go test ./... -run 'TestWebAppAITutorWordReportUIWiring'
```

Expected failure before implementation:

```text
missing ai tutor word report UI wiring
```

---

## Task 10: Implement React Word Report UI

- [ ] In `TutorView`, add dialog state for word report fields.
- [ ] Render a bottom action button on `word_learn` steps:

```tsx
<Button type="button" variant="ghost" onClick={openWordReportDialog}>
  <Bug className="..." />
  {t("ai_tutor_word_report_button")}
</Button>
```

- [ ] Add dialog fields:
  - correct learning-language word
  - correct translation
  - optional comment
- [ ] Submit to `POST /api/ai-tutor/word-report` with current `session_id` and current `stage`.
- [ ] On success, close the dialog, clear fields, and update the Tutor step from the API response so the next word is shown.
- [ ] Disable submit while pending and surface validation/API errors inline.
- [ ] Add copy keys in `web-react/src/lib/i18n.ts` for English and Russian.
- [ ] Add restrained styles in `web-react/src/styles/app.css` without creating nested cards.

UX constraints:

- The modal should not describe internal workflow or Telegram moderation.
- The current word should remain visible behind the modal contextually, but the modal inputs should be the active focus.
- The button appears only for `word_learn` steps, not story/quiz/speaking steps.

---

## Task 11: Add Pronunciation Free-Tier Regression Test

- [ ] Add or verify a test in `web_api_feature_test.go` proving free users cannot trigger pronunciation scoring.
- [ ] The test must fail if `/api/pronunciation/check` calls transcription/scoring providers before checking premium.

Target test name:

```go
func TestPronunciationCheckRequiresPremiumBeforeProviderCall(t *testing.T) {}
```

Expected behavior:

- Free authenticated user receives `402 Payment Required`.
- Mock provider transport is not called.
- No pronunciation attempt row is created for free users.

Note: if this test passes immediately, the requirement is already satisfied in production code; keep the regression test and avoid unnecessary production changes.

---

## Task 12: Verification

- [ ] Run targeted Go tests after each task group:

```powershell
go test ./... -run 'TestSQLiteStoreAITutorWordReport|TestAITutorWordReport|TestTelegramAITutorWordReport|TestApplyAITutorWordReportCorrection|TestWebAppAITutorWordReportUIWiring|TestPronunciationCheckRequiresPremiumBeforeProviderCall'
```

- [ ] Run broad backend tests:

```powershell
go test ./...
```

- [ ] Run frontend checks available in the workspace:

```powershell
npm --prefix web-react run build
```

- [ ] If an Nx target is discoverable, run the relevant build/test target through the repo's normal script or Nx command.
- [ ] Try Playwright/browser verification if the local browser policy allows it; otherwise record the policy blocker and rely on build/source tests.

Expected final status:

- All targeted tests pass.
- Broad Go tests pass.
- React build passes.
- Worktree contains only intended implementation/test/doc changes.
- Changes are committed and pushed to `origin/codex/ai-tutor-rebuild-fix`.

---

## Rollback Notes

- The new table is additive and should be safe to leave in place.
- If Telegram moderation must be disabled quickly, keep storing reports but skip Telegram notification when no ops recipients are configured.
- If vocabulary runtime update fails, the lesson payload update remains the canonical AI Tutor correction and should be logged for manual cleanup.
