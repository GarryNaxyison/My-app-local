# AI Vocabulary Fill And Cleanup Design

## Status

Approved for implementation on 2026-06-08 by user direction: keep SQLite as the hot vocabulary store, fill gaps through AI, and remove excess downloaded raw dictionary sources.

## Goal

Keep the existing open-source JSON dictionaries as the seed/fallback layer, add durable AI-generated word cards to SQLite when local vocabulary cannot provide a suitable new word, and remove local raw source downloads that are not needed at runtime.

## Current State

- `data/vocabulary/vocabulary_words*.json` contains the deployable learning dictionaries and is still required for first SQLite import and fallback.
- `vocabulary_sqlite.go` imports JSON into `vocabulary_words` and `vocabulary_translations`, then uses indexed SQLite lookups for hot paths.
- `vocabulary_ai_translations` already persists AI repairs for missing translations.
- `tools/sources/*` contains very large raw TEI/JSONL downloads used only for rebuilding dictionaries and is ignored by git.

## Design

Use a hybrid vocabulary model:

1. Static seed dictionaries remain in `data/vocabulary`.
2. AI-generated words are stored in a separate durable SQLite table, then mirrored into `vocabulary_words` and `vocabulary_translations`.
3. Duplicate checks use the existing deterministic vocabulary ID (`language:normalized-word`) so lookups remain primary-key fast.
4. The word lesson flow asks AI for one new card only when the local/SQLite pool cannot produce an unlearned word.
5. AI output is accepted only after validation: target word exists, CEFR level is valid, translation does not reveal the answer, and duplicate IDs are rejected.
6. Ignored raw source downloads under `tools/sources` can be removed locally to reduce disk usage. They can be downloaded again if a full dictionary rebuild is needed.

## Non-Goals

- Do not delete `data/vocabulary/*.json`.
- Do not replace the whole dictionary pipeline with live AI calls.
- Do not generate bulk 30k-word dictionaries in one request.

## Performance

Checking whether an AI-suggested word already exists is O(log n) through `vocabulary_words.id` primary key. A 30k-per-language vocabulary is small for SQLite when keyed by language and normalized word ID.

## Testing

- SQLite tests cover AI word insert, duplicate rejection, and normal `findVocabWord` retrieval.
- Sanitizer tests cover rejecting translations that equal/reveal the target word.
- Existing vocabulary tests continue to cover JSON import, SQLite lookup, and word selection.
