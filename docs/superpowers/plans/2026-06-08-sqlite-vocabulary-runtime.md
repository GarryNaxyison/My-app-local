# SQLite Vocabulary Runtime Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make SQLite the only runtime vocabulary store while preserving AI-filled vocabulary and keeping JSON dictionaries read-only seed files.

**Architecture:** Seed JSON files are imported into SQLite and tracked with source metadata. AI-created words and translations live in dedicated SQLite tables and are replayed into hot lookup tables after seed imports. Bot traffic never rewrites vocabulary JSON files.

**Tech Stack:** Go, `database/sql`, SQLite, existing vocabulary tests, existing Nx/go test scripts.

---

## File Structure

- Modify `vocabulary_test.go`: replace the JSON persistence expectation with SQLite-only persistence tests and update source freshness expectations.
- Modify `vocabulary_sqlite.go`: remove runtime JSON rewrite from AI translation writes and compare JSON modification time in source freshness checks.
- Modify `README.md`: clarify that production must preserve `vocabulary.sqlite` or AI export backups.
- Modify `docs/reference/VOCABULARY_SOURCES.md`: document JSON as read-only seed and SQLite as the mutable AI layer.

## Task 1: Red Tests For SQLite-Only AI Translations

- [ ] Change `TestSQLiteVocabularyAITranslationPersistsToJSONDictionary` into `TestSQLiteVocabularyAITranslationDoesNotMutateJSONSeed`.

```go
func TestSQLiteVocabularyAITranslationDoesNotMutateJSONSeed(t *testing.T) {
	word := vocabWord{
		ID:       "en:creating",
		Language: "en",
		English:  "creating",
		Russian:  "creating",
		Translations: map[string]string{
			"ru": "creating",
		},
		Level: "A2",
	}
	configureSQLiteVocabularyForTest(t, []vocabWord{word})
	path, ok := vocabularyPath("en")
	if !ok {
		t.Fatal("missing English vocabulary path")
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := sqliteVocabularyAITranslationSet(word, "ru", "test-model", "prompt", "создание; сотворение"); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatal("AI translation mutated read-only JSON seed")
	}
	value, ok, err := sqliteVocabularyAITranslationGet(word.ID, "ru")
	if err != nil || !ok || value != "создание; сотворение" {
		t.Fatalf("sqliteVocabularyAITranslationGet = %q, %v, %v; want SQLite persistence", value, ok, err)
	}
}
```

- [ ] Run `go test ./... -run TestSQLiteVocabularyAITranslationDoesNotMutateJSONSeed -count=1`.
- [ ] Expected result before implementation: FAIL with `AI translation mutated read-only JSON seed`.

## Task 2: Red Test For JSON Source Freshness

- [ ] Rename `TestSQLiteVocabularySourceFreshIgnoresMTimeOnlyChanges` to `TestSQLiteVocabularySourceFreshDetectsMTimeOnlyChanges`.
- [ ] Change its final assertion to:

```go
if fresh {
	t.Fatal("expected mtime-only vocabulary changes to require reimport")
}
```

- [ ] Run `go test ./... -run TestSQLiteVocabularySourceFreshDetectsMTimeOnlyChanges -count=1`.
- [ ] Expected result before implementation: FAIL because the current code ignores mtime-only changes.

## Task 3: Minimal Production Fix

- [ ] In `sqliteVocabularySourceFresh`, require stored `mod_time_unix` to equal `info.ModTime().Unix()`.
- [ ] In `sqliteVocabularyAITranslationSet`, delete the call to `persistVocabularyTranslationToJSON`.
- [ ] Keep `persistVocabularyTranslationToJSON` only if tests still reference it; otherwise remove it and related unused imports.
- [ ] Run:

```bash
go test ./... -run "TestSQLiteVocabularyAITranslationDoesNotMutateJSONSeed|TestSQLiteVocabularySourceFreshDetectsMTimeOnlyChanges|TestSQLiteVocabularyAITranslationPersistsAndJoins|TestSQLiteVocabularyImportSkipsStaleAITranslations|TestSQLiteVocabularyAIWordReplayAfterJSONImport" -count=1
```

- [ ] Expected result after implementation: PASS.

## Task 4: Documentation

- [ ] Update `README.md` low-memory section to say runtime AI vocabulary is in `vocabulary.sqlite` and must be backed up/preserved.
- [ ] Update `docs/reference/VOCABULARY_SOURCES.md` to say JSON seeds are not mutated by bot traffic and AI writes are SQLite-only.
- [ ] Run `node tools/check_encoding_artifacts.mjs`.

## Final Verification

- [ ] Run `go test ./... -count=1`.
- [ ] Run `node tools/check_encoding_artifacts.mjs`.
- [ ] Run `git diff --check`.
- [ ] Commit and push to `origin/main`.
