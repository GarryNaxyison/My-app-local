# Shared AI-generated vocabulary

## Goal

Replace the imported local word banks with one durable, shared SQLite dictionary populated by validated AI-generated vocabulary. Learners select unused words from that shared pool; when their level-appropriate pool is exhausted, the app asks the vocabulary model for a new card and stores it for future learners.

## Scope and constraints

- Stop importing and deploying `data/vocabulary/vocabulary_words*.json` as runtime dictionaries.
- On startup, retain cards in `vocabulary_ai_words` and remove legacy imported rows from the hot `vocabulary_words`/`vocabulary_translations` tables and source manifest.
- Apply one one-time SQLite migration that clears all users' learning progress (learned words, mistakes, phrasebook entries, habit/streak and XP counters, lesson history and assignments, active/completed AI Tutor sessions, answers, reviews, and reports). Preserve accounts, credentials, language preferences, reminders, referrals, and navigation settings.
- Only AI-generated shared cards are eligible for new-word selection and answer-choice distractors.
- Keep the shared dictionary keyed by language plus normalized word ID, level, translations, and metadata in SQLite.
- Reject duplicate suggestions against both the shared AI dictionary and any legacy rows still present during rollout. Ask the model again within the bounded generation attempt limit; never insert or teach the duplicate.
- Generate cards for the learner's CEFR level and learning language. Validate word, translation, CEFR level, and serialization before inserting.
- Do not create backups during migration or deployment.

## Observable behavior

1. A new learner is offered cards already in the shared SQLite pool for their learning language and CEFR band.
2. The same shared word can be learned by different users; each learner's learned state remains independent.
3. Once a learner has used all suitable pool cards, the system requests a new level-matched card, saves it to SQLite, and makes it available to other learners.
4. If generation returns an existing word, that candidate is rejected and another request is made; no duplicate row or repeated card is created.
5. Legacy imported dictionaries are not re-imported at startup or included in the server web/app deployment package.
6. The one-time global progress reset runs exactly once; accounts and profile preferences survive, and future progress remains intact across restarts.
7. A new learner can start an AI Tutor lesson with no learned words and an empty dictionary. After lesson quality approval, valid single words with Russian translations seed the shared SQLite pool; multiword chunks are skipped and duplicates remain rejected.

## Acceptance criteria

- Runtime initialization migrates away legacy static rows without deleting AI-generated shared rows.
- Selection and distractor queries use only AI-generated shared rows.
- Empty-pool selection falls through to generation rather than a local JSON fallback.
- Duplicate generation is retried and duplicate insertion remains protected transactionally.
- Go tests cover migration, AI-only selection, and duplicate rejection/retry.
- Deployment packaging and runtime documentation no longer require the old language JSON files.
- SQLite tests prove the reset clears progress, preserves accounts/preferences, and does not repeat on later startup.
- Go regression test proves AI Tutor opens for a new learner with an empty SQLite pool and seeds only eligible unique single-word entries without extra model calls.
