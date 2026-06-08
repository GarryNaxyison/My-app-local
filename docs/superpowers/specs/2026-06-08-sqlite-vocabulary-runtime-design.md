# SQLite Vocabulary Runtime Design

## Status

Approved for implementation on 2026-06-08 after the production server ran out of RAM during vocabulary work.

## Goal

Make SQLite the only runtime vocabulary store, keep JSON dictionaries as read-only seed files, and preserve all AI-filled vocabulary without spending AI tokens again after restart, reimport, or deploy.

## Problem

The current AI translation write path stores the translation in SQLite and then calls `persistVocabularyTranslationToJSON`. That helper reads the full language JSON file, unmarshals it into `[]vocabWord`, marshals it again, and replaces the source file. English and Russian seed dictionaries are tens of megabytes on disk, so this can create large transient heap pressure on a small VPS. It also makes user traffic mutate seed files, which should be rebuild artifacts rather than hot runtime state.

The SQLite import metadata stores JSON modification time for diagnostics, but low-memory runtime should not reimport on mtime-only changes. Routine deploys can rewrite timestamps without changing content, and reimporting large dictionaries during deploy is unnecessary pressure on a 1 GB VPS.

## Design

Use a two-layer vocabulary model:

1. Seed layer: `data/vocabulary/vocabulary_words*.json`.
   - Read-only at runtime.
   - Imported into `vocabulary_words` and `vocabulary_translations`.
   - Tracked in `vocabulary_sources` by path, size, modification time, word count, and import timestamp.
2. AI mutable layer: `vocabulary_ai_words` and `vocabulary_ai_translations`.
   - This is the durable value created with paid AI calls.
   - JSON reimport must never delete these tables.
   - After seed import, AI words and AI translations are replayed into the hot `vocabulary_words` and `vocabulary_translations` tables.

## Runtime Rules

- `sqliteVocabularyAITranslationSet` writes only to SQLite.
- AI translations are still mirrored into `vocabulary_translations` for normal lookup.
- AI words are still mirrored into `vocabulary_words` and `vocabulary_translations`.
- Seed JSON files are never rewritten by normal bot traffic.
- JSON mtime-only changes do not trigger reimport when size and SQLite word count still match.
- Production deploy must preserve `vocabulary.sqlite` unless an explicit migration/restore step is being performed.

## Indexing

SQLite remains the indexed runtime store. Required indexes:

- `vocabulary_words(language, position)`
- `vocabulary_words(language, level, position)`
- `vocabulary_words(language, frequency_rank)`
- `vocabulary_words(language, level, frequency_rank, position)`
- `vocabulary_words(language, position, id)`
- `vocabulary_translations(language)`
- `vocabulary_translations(language, word_id)`
- `vocabulary_ai_cache(word_id, kind, interface_language)`
- `vocabulary_ai_translations(target_language, learning_language)`
- `vocabulary_ai_words(language, level)`
- `vocabulary_ai_words(language, word)`

JSON is not indexed separately. The JSON "index" is the `vocabulary_sources` metadata row plus SQLite import tables.

## Backup Requirement

Backups must include `vocabulary.sqlite` or an export of `vocabulary_ai_words` and `vocabulary_ai_translations`. Losing only seed JSON is recoverable. Losing the AI mutable layer means spending AI tokens again.

## Tests

- AI translation insert must not change the JSON seed file.
- AI translation must remain available from SQLite after insert.
- AI translation must survive JSON reimport when the seed word still exists.
- Stale AI translation must not resurrect a deleted seed word.
- JSON mtime-only change must be treated as fresh to avoid deploy-time reimports.
- Existing SQLite index coverage test must keep passing.
