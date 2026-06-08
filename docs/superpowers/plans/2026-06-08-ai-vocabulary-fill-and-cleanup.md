# AI Vocabulary Fill And Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add durable AI-generated vocabulary cards on top of the existing SQLite dictionary and remove ignored raw dictionary downloads that are not needed at runtime.

**Architecture:** Keep JSON dictionaries as seed/fallback. Add a `vocabulary_ai_words` table and mirror accepted AI words into `vocabulary_words` plus `vocabulary_translations` so all existing lookup paths keep working. Use the existing deterministic vocabulary ID as the duplicate guard.

**Tech Stack:** Go, SQLite via `database/sql`, existing OpenRouter client, existing vocabulary tests.

---

## File Structure

- Modify `vocabulary_sqlite.go`: add AI word table, indexes, insert/upsert helpers, duplicate checks.
- Modify `vocabulary.go`: add generated-card validation and conversion helpers.
- Modify `prompts.go`: add strict JSON prompt for one AI-generated vocabulary card.
- Modify `bot.go`: when no local unlearned word exists, request one AI card and persist it.
- Modify `vocabulary_test.go`: add red/green tests for insert, duplicate rejection, and sanitizer behavior.
- Modify `docs/reference/VOCABULARY_SOURCES.md`: document that `tools/sources` is rebuild-only and can be absent in runtime deployments.
- Local cleanup: remove ignored `tools/sources/*` after verifying the resolved path.

## Task 1: SQLite AI Word Persistence

- [ ] Add failing tests:
  - `TestSQLiteVocabularyAIWordSetPersistsAndFindsWord`
  - `TestSQLiteVocabularyAIWordSetRejectsDuplicateWord`
- [ ] Add `vocabulary_ai_words` table in `initSQLiteVocabularyTables`.
- [ ] Add index/primary-key-backed helper:
  - `sqliteVocabularyWordExists(language, text string) (bool, error)`
  - `sqliteVocabularyAIWordSet(word vocabWord, model string, prompt string) (bool, error)`
- [ ] Verify with:
  - `go test ./... -run "TestSQLiteVocabularyAIWordSet"`
- [ ] Commit:
  - `feat: persist ai vocabulary words`

## Task 2: AI Card Validation And Prompt

- [ ] Add failing tests:
  - `TestSanitizeGeneratedVocabularyWordRejectsTargetAsTranslation`
  - `TestSanitizeGeneratedVocabularyWordNormalizesLevelAndID`
- [ ] Add strict JSON card type and sanitizer:
  - `generatedVocabularyWord`
  - `sanitizeGeneratedVocabularyWord`
- [ ] Add prompt builder:
  - `vocabularyGenerationPrompt(user userState, forbidden []string) []chatMessage`
- [ ] Verify with:
  - `go test ./... -run "GeneratedVocabulary|VocabularyGeneration"`
- [ ] Commit:
  - `feat: validate ai vocabulary cards`

## Task 3: Word Lesson AI Fallback

- [ ] Update `bot.startWordLesson` so it tries `nextUnlearnedWord` first.
- [ ] If no word is available and OpenRouter vocabulary model is configured, request one AI card, reject duplicates, persist it, and show it through the existing word lesson UI.
- [ ] Keep current local fallback messages when AI is unavailable or validation fails.
- [ ] Verify with:
  - `go test ./... -run "Vocabulary|WordLesson"`
- [ ] Commit:
  - `feat: fill vocabulary gaps with ai`

## Task 4: Cleanup Raw Sources

- [ ] Verify `tools/sources` resolves inside the repository.
- [ ] Remove ignored files under `tools/sources`.
- [ ] Update vocabulary documentation.
- [ ] Verify `git status --short --ignored tools/sources`.
- [ ] Commit docs only if documentation changed:
  - `docs: document vocabulary source cleanup`

## Final Verification

- [ ] `go test ./... -run "Vocabulary|WordLesson"`
- [ ] `go test ./...`
- [ ] `node tools/check_encoding_artifacts.mjs`
- [ ] `git status --short`
- [ ] `git push origin main`
