# AI Tutor Word Report Admin Fix Design

## Status

Approved for design on 2026-06-14. The implementation must use the project runtime database, which is SQLite in this repository. There is no MySQL runtime table in the current codebase.

## Goal

Let learners report a wrong word during the AI Tutor word-learning stage, let the owner approve, reject, or correct the report from Telegram, and apply accepted corrections to both the AI Tutor lesson payload and the durable vocabulary runtime layer without exposing the server to report spam.

## Problem

AI Tutor lessons are generated and stored as JSON in `ai_tutor_lessons.payload_json`. A wrong word or translation can therefore stay in an already-approved lesson and keep appearing in future sessions. The learner currently has no direct way to flag the active word from the word-learning block.

The existing web bug report endpoint sends generic reports to Telegram, but it is not structured enough for safe automatic vocabulary replacement. A direct "user submits replacement and the database changes" flow would be unsafe because a user could spam reports, overwrite dictionary content, or force repeated Telegram/API/database work.

Pronunciation scoring is expensive. Backend routes must keep scoring unavailable on free plans even if the frontend hides controls, because a user can call API endpoints directly.

## Scope

This change covers only AI Tutor word-learning reports from the web app and the owner moderation flow in Telegram.

In scope:

- Add a "Report word error" button under `word_learn` material in the AI Tutor web view.
- Open a modal asking for the correct target-language word, the user's proposed translation, and an optional short comment.
- Submit a structured report to a new backend endpoint.
- After a successful report, automatically skip the current word by submitting `continue` to the existing AI Tutor answer endpoint.
- Store reports in SQLite with pending, accepted, rejected, applied, and failed states.
- Notify ops recipients in Telegram with inline moderation buttons.
- Let the owner choose Accept, Reject, or Fix.
- For Fix, ask the owner for text in `word - translation` format and apply that value.
- Apply accepted/fixed corrections to both AI Tutor lesson JSON and the vocabulary runtime/AI mutable layer.
- Add backend regression tests for rate limits, ownership checks, admin-only moderation, storage updates, and pronunciation premium gating.

Out of scope:

- Generic bug-report screenshots and file uploads.
- A public admin web dashboard.
- Re-generating whole AI Tutor lessons with an LLM.
- Rewriting JSON seed vocabulary files during normal runtime.
- Changing word recall, production, or story stages beyond reflecting corrected word data.

## User Experience

In the AI Tutor word-learning block, the learner sees a small secondary button below the word card:

`Report word error`

The modal contains:

- `Correct word`
- `Correct translation`
- `Comment`
- Send and cancel actions

Validation happens before submission:

- correct word: 1-80 runes;
- translation: 1-160 runes;
- comment: 0-500 runes;
- whitespace is normalized;
- control characters are removed.

After a successful response, the UI shows a short status message and immediately advances the AI Tutor session with the existing `continue` behavior. If report submission fails because of rate limits, the user stays on the same word and sees the backend error.

## Backend API

Add a new authenticated endpoint:

`POST /api/ai-tutor/word-report`

Request:

```json
{
  "session_id": "session id",
  "stage": "word_learn_1",
  "proposed_word": "correct target-language word",
  "proposed_translation": "correct interface-language translation",
  "comment": "optional learner comment"
}
```

Server-side checks:

- the current web user must be authenticated;
- `session_id` must exist and belong to the current user;
- the session's current stage must equal the requested stage;
- the stage must be `word_learn_N`;
- the lesson payload must contain a word at that index;
- proposed fields must pass length and text validation;
- duplicate pending reports for the same `telegram_id + lesson_id + word_id` are rejected or return the existing pending report id;
- rate limits must pass before Telegram notification is attempted.

Response:

```json
{
  "ok": true,
  "id": "awr_...",
  "deduplicated": false
}
```

## Data Model

Add a table in the main SQLite store:

```sql
CREATE TABLE IF NOT EXISTS ai_tutor_word_reports (
  id TEXT PRIMARY KEY,
  telegram_id INTEGER NOT NULL,
  session_id TEXT NOT NULL,
  lesson_id TEXT NOT NULL,
  stage TEXT NOT NULL,
  word_id TEXT NOT NULL,
  word_index INTEGER NOT NULL,
  learning_language TEXT NOT NULL,
  interface_language TEXT NOT NULL,
  original_word TEXT NOT NULL,
  original_translation TEXT NOT NULL,
  proposed_word TEXT NOT NULL,
  proposed_translation TEXT NOT NULL,
  user_comment TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  admin_chat_id TEXT NOT NULL DEFAULT '',
  admin_message_id INTEGER NOT NULL DEFAULT 0,
  admin_user_id INTEGER NOT NULL DEFAULT 0,
  admin_fix_prompt_message_id INTEGER NOT NULL DEFAULT 0,
  final_word TEXT NOT NULL DEFAULT '',
  final_translation TEXT NOT NULL DEFAULT '',
  decision_comment TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  decided_at TEXT NOT NULL DEFAULT '',
  applied_at TEXT NOT NULL DEFAULT '',
  error_text TEXT NOT NULL DEFAULT ''
);
```

Indexes:

- `idx_ai_tutor_word_reports_user_created` on `(telegram_id, created_at)`
- `idx_ai_tutor_word_reports_status_created` on `(status, created_at)`
- `idx_ai_tutor_word_reports_pending_word` unique partial index on `(telegram_id, lesson_id, word_id)` where `status IN ('pending','awaiting_fix')`
- `idx_ai_tutor_word_reports_fix_prompt` on `(admin_chat_id, admin_fix_prompt_message_id, status)`

## Anti-Abuse Rules

The word-report endpoint must be cheap and bounded:

- no multipart parsing and no file uploads;
- body size limited to a small JSON payload;
- per-user limit: 3 reports per 10 minutes and 10 reports per day;
- per-session limit: 2 reports per session;
- global Telegram notification limit: 30 notifications per 10 minutes; reports above this limit are stored but marked for deferred notification or sent only to logs;
- exact duplicate pending reports are deduplicated;
- all text fields are length-capped before database insert;
- Telegram messages are plain text or MarkdownV2-escaped and split within Telegram limits;
- admin actions are idempotent.

The admin moderation path must be protected:

- only configured ops/admin recipients from `TELEGRAM_OPS_RECIPIENTS` can act;
- a callback from any other Telegram user or chat is ignored and answered with a generic rejection;
- moderation callbacks include only the report id and action, never raw replacement text;
- replacement text is accepted only after an admin clicked Fix and replies to the bot's fix prompt, or includes the report id in a strict command format if reply metadata is unavailable;
- a report can transition out of pending state once.

## Telegram Moderation Flow

When a report is accepted by the endpoint, the bot sends an ops message:

```text
AI Tutor word report
ID: awr_...
User: name (telegram_id)
Lesson: lesson_id
Stage: word_learn_2
Language: en / interface ru

Original:
word - translation

User suggests:
word - translation

Comment:
...

Choose: Accept, Reject, Fix
```

Inline buttons:

- `Accept` callback: `ait_word_report|<id>|accept`
- `Reject` callback: `ait_word_report|<id>|reject`
- `Fix` callback: `ait_word_report|<id>|fix`

Accept:

- checks admin authorization;
- applies the user's proposed word and translation;
- marks report applied or failed;
- edits or sends a short status to Telegram.

Reject:

- checks admin authorization;
- marks report rejected;
- no vocabulary or lesson change is made.

Fix:

- checks admin authorization;
- marks report `awaiting_fix`;
- sends a prompt asking the admin to reply with `word - translation`;
- records `admin_chat_id`, `admin_user_id`, and the fix prompt message id.

Admin text response:

- must come from the same authorized admin chat/user;
- must be a reply to the recorded fix prompt when Telegram reply metadata is available;
- parses exactly one separator: ` - `;
- applies parsed values if valid;
- marks report applied or failed.

## Applying Corrections

Correction application must be idempotent and must use explicit transactions for each SQLite database handle it touches. Because the main app store and the vocabulary runtime can be separate SQLite databases, the implementation uses a deterministic two-step sequence: validate all data before writes, upsert the vocabulary correction in the vocabulary database transaction, then update the AI Tutor lesson/report in the main database transaction. If the second transaction fails, retrying the same admin action repeats the same vocabulary upsert and completes the main-store update, so the final state converges instead of duplicating rows.

AI Tutor layer:

- load `ai_tutor_lessons` by `lesson_id`;
- update `Payload.Words[word_index]`;
- preserve the existing `id` so `word_learning` and `word_recall` still point to the same word;
- set `target`, `interface_translation`, and `audio_text_target` to the final word;
- leave example sentence fields unchanged unless they exactly match the old word text, in which case replace that exact occurrence conservatively;
- recompute `fingerprint`;
- update `payload_json`, `fingerprint`, and `updated_at`.

Vocabulary runtime layer:

- find an existing `vocabulary_words` row by original `word_id`;
- update `word` to the final target word and `russian` to the final translation when the interface language is Russian;
- upsert `vocabulary_translations(word_id, interface_language)` to the final translation;
- if the source row exists in `vocabulary_ai_words`, mirror the same update there, including `translations_json`;
- upsert `vocabulary_ai_translations(word_id, target_language)` for the interface translation;
- if the original word id does not exist in vocabulary runtime, create an AI mutable word with the existing AI Tutor `word_id`, language, final word, final translation, level, and source `ai_tutor_report`;
- invalidate the SQLite vocabulary id cache after a successful update.

JSON seed vocabulary files remain read-only. Corrections survive runtime reimport because they are mirrored into the AI mutable layer.

## Pronunciation Premium Gate

Backend must enforce that expensive pronunciation scoring is not available to free users.

Rules:

- `/api/pronunciation/check` requires premium before transcription/scoring.
- `/api/shadowing/answer` with voice requires premium before direct audio assessment.
- lesson/practice voice answer paths keep their existing premium voice checks.
- free learned-word pronunciation audio can remain free because it is a separate text-to-speech path and does not run pronunciation scoring.
- tests must verify free users receive `402 Payment Required` before any OpenRouter scoring call is made.

## Error Handling

If Telegram notification fails, the report remains stored as pending with `error_text` populated. The user still receives success only if the report was saved.

If applying a correction fails after the admin action, the report moves to `failed` with `error_text`, and Telegram receives a failure message. The admin may retry by using Fix again only if the report is returned to pending by an explicit internal helper; normal callbacks remain idempotent.

If the AI Tutor lesson was deleted or the word index no longer matches the original word id, the report fails closed and does not update vocabulary.

## Testing

Backend tests:

- word report rejects unauthenticated requests.
- word report rejects a session owned by another user.
- word report rejects non-`word_learn` stages.
- valid report stores original/proposed fields and notifies Telegram once.
- duplicate pending report does not spam Telegram.
- per-user and per-session rate limits reject excess reports.
- non-admin Telegram callback cannot accept, reject, or fix a report.
- Accept updates `ai_tutor_lessons.payload_json` and vocabulary runtime/AI mutable rows.
- Reject closes the report without changing lesson or vocabulary rows.
- Fix waits for admin text and applies `word - translation`.
- repeated Accept/Reject/Fix actions are idempotent.
- free `/api/pronunciation/check` returns `402` without calling scoring.
- premium `/api/pronunciation/check` still works with the existing scoring/transcription fallback.

Frontend tests:

- `word_learn` renders a report button.
- modal validates required word and translation.
- successful report calls the report endpoint and then advances to the next tutor step.
- rate-limit errors keep the current word visible and show a status error.

Full verification:

- `npm test` or `npx nx run english-coach-bot:test`
- targeted Go tests for AI Tutor reports, vocabulary updates, Telegram callbacks, and pronunciation gating
- frontend build or relevant React test/e2e target

## Rollout

Ship backend protections first, then the web button. Existing generic `/api/bug-report` remains unchanged. The new report flow is narrow, authenticated, rate-limited, and admin-only for mutation, so a public user can at most create bounded pending rows and cannot directly modify vocabulary or flood Telegram.
