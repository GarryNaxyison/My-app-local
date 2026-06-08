package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode"
)

func seedVocabularyRandomForTest(t *testing.T, seed int64) {
	t.Helper()
	vocabularyRandom.mu.Lock()
	previous := vocabularyRandom.rng
	vocabularyRandom.rng = rand.New(rand.NewSource(seed))
	vocabularyRandom.mu.Unlock()
	t.Cleanup(func() {
		vocabularyRandom.mu.Lock()
		vocabularyRandom.rng = previous
		vocabularyRandom.mu.Unlock()
	})
}

func configureSQLiteVocabularyForTest(t *testing.T, rows []vocabWord) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("VOCABULARY_DIR", dir)
	payload, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "vocabulary_words.json"), payload, 0o644); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "vocabulary.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		configureSQLiteVocabulary(nil)
		_ = db.Close()
	})
	if err := syncSQLiteVocabularyFromJSON(db); err != nil {
		t.Fatal(err)
	}
	configureSQLiteVocabulary(db)
}

func testVocabWord(english string, level string, rank int) vocabWord {
	return vocabWord{
		English:       english,
		Russian:       "ru " + english,
		Translations:  map[string]string{"ru": "ru " + english},
		Level:         level,
		FrequencyRank: rank,
	}
}

func TestVocabularyLoadsLargeDictionary(t *testing.T) {
	if got := vocabularySize(); got < minimumVocabularySize {
		t.Fatalf("vocabularySize() = %d, want at least %d", got, minimumVocabularySize)
	}
	if _, ok := findVocabWord("make"); !ok {
		t.Fatal("expected common word make in vocabulary")
	}
	if _, ok := findVocabWord("apple"); !ok {
		t.Fatal("expected curated basic word apple in vocabulary")
	}
	for language, minimum := range map[string]int{
		"ru": 40000,
		"es": 10000,
		"de": 19000,
		"fr": 19000,
		"it": 12000,
		"zh": 6000,
		"ja": 3000,
		"ko": 5000,
		"tg": 2000,
		"uz": 2500,
		"tt": 1400,
		"hy": 13000,
		"kk": 9000,
		"ky": 2100,
		"ka": 10000,
		"uk": 23000,
		"pl": 20000,
		"ro": 29000,
		"pt": 8000,
	} {
		if got := len(vocabularyForLanguage(language)); got < minimum {
			t.Fatalf("vocabularyForLanguage(%q) = %d, want at least %d", language, got, minimum)
		}
	}
}

func TestSQLiteVocabularyRuntimeImportsAndFindsWords(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("VOCABULARY_DIR", dir)
	vocabularyJSON := `[
		{"english": "alpha", "russian": "\u0430\u043b\u044c\u0444\u0430", "translations": {"es": "alfa", "ru": "\u0430\u043b\u044c\u0444\u0430"}, "level": "A1", "frequency_rank": 1},
		{"english": "bravo", "russian": "\u0431\u0440\u0430\u0432\u043e", "translations": {"es": "bravo", "ru": "\u0431\u0440\u0430\u0432\u043e"}, "level": "A1", "frequency_rank": 2},
		{"english": "charlie", "russian": "\u0447\u0430\u0440\u043b\u0438", "translations": {"es": "charlie", "ru": "\u0447\u0430\u0440\u043b\u0438"}, "level": "A1", "frequency_rank": 3},
		{"english": "delta", "russian": "\u0434\u0435\u043b\u044c\u0442\u0430", "translations": {"es": "delta", "ru": "\u0434\u0435\u043b\u044c\u0442\u0430"}, "level": "A1", "frequency_rank": 4},
		{"english": "echo", "russian": "\u044d\u0445\u043e", "translations": {"es": "eco", "ru": "\u044d\u0445\u043e"}, "level": "A2", "frequency_rank": 5}
	]`
	if err := os.WriteFile(filepath.Join(dir, "vocabulary_words.json"), []byte(vocabularyJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "vocabulary.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		configureSQLiteVocabulary(nil)
		_ = db.Close()
	})
	if err := syncSQLiteVocabularyFromJSON(db); err != nil {
		t.Fatal(err)
	}
	configureSQLiteVocabulary(db)

	word, ok := findVocabWord("alpha")
	if !ok {
		t.Fatal("expected SQLite-backed vocabulary lookup to find alpha")
	}
	if word.ID != "en:alpha" || word.Language != "en" {
		t.Fatalf("unexpected word identity: %#v", word)
	}
	if got := wordTranslation(word, "es"); got != "alfa" {
		t.Fatalf("wordTranslation(alpha, es) = %q, want alfa", got)
	}
	next, ok := nextUnlearnedWord(userState{LearningLanguage: "en", InterfaceLanguage: "ru", Level: "A1"})
	if !ok || normalizeLearningLanguage(next.Language) != "en" || normalizeCEFRLevel(next.Level) != "A1" {
		t.Fatalf("nextUnlearnedWord() = %#v, %v; want random en A1 word", next, ok)
	}
	options := wordOptions(word)
	if len(options) != 4 {
		t.Fatalf("wordOptions(alpha) returned %d options, want 4: %#v", len(options), options)
	}

	count := 0
	if err := forEachVocabularyWord("en", func(word vocabWord) bool {
		count++
		return true
	}); err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Fatalf("SQLite vocabulary iteration count = %d, want 5", count)
	}
}

func TestSQLiteVocabularySourceFreshIgnoresMTimeOnlyChanges(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("VOCABULARY_DIR", dir)
	path := filepath.Join(dir, "vocabulary_words.json")
	vocabularyJSON := `[
		{"english": "alpha", "russian": "\u0430\u043b\u044c\u0444\u0430", "level": "A1", "frequency_rank": 1},
		{"english": "bravo", "russian": "\u0431\u0440\u0430\u0432\u043e", "level": "A1", "frequency_rank": 2}
	]`
	if err := os.WriteFile(path, []byte(vocabularyJSON), 0o644); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "vocabulary.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	if err := syncSQLiteVocabularyFromJSON(db); err != nil {
		t.Fatal(err)
	}
	original, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	changedModTime := original.ModTime().Add(2 * time.Hour)
	if err := os.Chtimes(path, changedModTime, changedModTime); err != nil {
		t.Fatal(err)
	}
	changed, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if changed.Size() != original.Size() || changed.ModTime().Unix() == original.ModTime().Unix() {
		t.Fatalf("test setup did not change only mtime: original=%v changed=%v", original, changed)
	}

	fresh, err := sqliteVocabularySourceFresh(db, "en", path, changed)
	if err != nil {
		t.Fatal(err)
	}
	if !fresh {
		t.Fatal("expected mtime-only vocabulary changes to stay fresh")
	}
}

func TestWordOptionsVaryDistractorsAcrossRounds(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("VOCABULARY_DIR", dir)
	words := []string{
		"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf", "hotel",
		"india", "juliet", "kilo", "lima", "mango", "november", "oscar", "papa",
	}
	var builder strings.Builder
	builder.WriteString("[")
	for index, word := range words {
		if index > 0 {
			builder.WriteString(",")
		}
		fmt.Fprintf(
			&builder,
			`{"english":%q,"russian":%q,"translations":{"ru":%q},"level":"A1","frequency_rank":%d}`,
			word,
			"ru "+word,
			"ru "+word,
			index+1,
		)
	}
	builder.WriteString("]")
	if err := os.WriteFile(filepath.Join(dir, "vocabulary_words.json"), []byte(builder.String()), 0o644); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "vocabulary.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		configureSQLiteVocabulary(nil)
		_ = db.Close()
	})
	if err := syncSQLiteVocabularyFromJSON(db); err != nil {
		t.Fatal(err)
	}
	configureSQLiteVocabulary(db)

	sets := map[string]bool{}
	for _, correct := range words[:6] {
		word, ok := findVocabWord("en:" + correct)
		if !ok {
			t.Fatalf("test vocabulary word not found: %s", correct)
		}
		options := wordOptions(word)
		if len(options) != 4 {
			t.Fatalf("wordOptions(%s) returned %d options, want 4: %#v", correct, len(options), options)
		}
		seenText := map[string]bool{}
		wrong := []string{}
		for _, option := range options {
			if seenText[option.English] {
				t.Fatalf("wordOptions(%s) returned duplicate visible option %q: %#v", correct, option.English, options)
			}
			seenText[option.English] = true
			if option.ID != word.ID {
				wrong = append(wrong, option.English)
			}
		}
		sort.Strings(wrong)
		sets[strings.Join(wrong, ",")] = true
	}
	if len(sets) < 2 {
		t.Fatalf("wrong answer variants did not change across words: %#v", sets)
	}
}

func TestWordOptionsUseCEFRBandDistractors(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		testVocabWord("alpha", "A1", 1),
		testVocabWord("bravo", "A2", 2),
		testVocabWord("charlie", "A2", 3),
		testVocabWord("delta", "A2", 4),
		testVocabWord("echo", "A2", 5),
		testVocabWord("foxtrot", "B1", 6),
		testVocabWord("golf", "B2", 7),
	})

	correct, ok := findVocabWord("en:alpha")
	if !ok {
		t.Fatal("expected test word")
	}
	options := wordOptions(correct)
	if len(options) != 4 {
		t.Fatalf("wordOptions returned %d options, want 4: %#v", len(options), options)
	}
	for _, option := range options {
		if option.ID == correct.ID {
			continue
		}
		if !vocabWordInOptionBand(option, correct) {
			t.Fatalf("expected A-band distractor for %s, got %s (%s): %#v", correct.English, option.English, option.Level, options)
		}
	}
}

func TestLearnedWordOptionsUseDictionaryBandNotFixedLearnedPool(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		testVocabWord("alpha", "A1", 1),
		testVocabWord("bravo", "A1", 2),
		testVocabWord("charlie", "A1", 3),
		testVocabWord("delta", "A1", 4),
		testVocabWord("echo", "A2", 5),
		testVocabWord("foxtrot", "A2", 6),
		testVocabWord("golf", "A2", 7),
		testVocabWord("hotel", "A2", 8),
		testVocabWord("india", "B1", 9),
		testVocabWord("juliet", "B2", 10),
	})

	correct, ok := findVocabWord("en:alpha")
	if !ok {
		t.Fatal("expected test word")
	}
	learnedIDs := []string{"en:alpha", "en:bravo", "en:charlie", "en:delta"}
	options := learnedWordOptions(correct, learnedIDs)
	if len(options) != 4 {
		t.Fatalf("learnedWordOptions returned %d options, want 4: %#v", len(options), options)
	}
	learned := map[string]bool{}
	for _, id := range learnedIDs {
		learned[id] = true
	}
	foundDictionaryDistractor := false
	for _, option := range options {
		if option.ID == correct.ID {
			continue
		}
		if learned[option.ID] {
			t.Fatalf("expected distractors outside the user's fixed learned pool, got %s in %#v", option.ID, options)
		}
		if !vocabWordInOptionBand(option, correct) {
			t.Fatalf("expected A-band distractor, got %s (%s): %#v", option.English, option.Level, options)
		}
		foundDictionaryDistractor = true
	}
	if !foundDictionaryDistractor {
		t.Fatalf("expected at least one dictionary distractor: %#v", options)
	}
}

func TestSQLiteVocabularyCreatesPerformanceIndexes(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		testVocabWord("alpha", "A1", 1),
		testVocabWord("bravo", "A1", 2),
		testVocabWord("charlie", "A2", 3),
		testVocabWord("delta", "B1", 4),
	})

	db := currentSQLiteVocabularyDB()
	if db == nil {
		t.Fatal("expected SQLite vocabulary database")
	}
	wantWordIndexes := map[string]bool{
		"idx_vocabulary_words_language_position":            false,
		"idx_vocabulary_words_language_level_position":      false,
		"idx_vocabulary_words_language_rank":                false,
		"idx_vocabulary_words_language_level_rank_position": false,
		"idx_vocabulary_words_language_position_id":         false,
	}
	rows, err := db.Query(`PRAGMA index_list(vocabulary_words)`)
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var seq int
		var name string
		var unique int
		var origin string
		var partial int
		if err := rows.Scan(&seq, &name, &unique, &origin, &partial); err != nil {
			t.Fatal(err)
		}
		if _, ok := wantWordIndexes[name]; ok {
			wantWordIndexes[name] = true
		}
	}
	if err := rows.Close(); err != nil {
		t.Fatal(err)
	}
	for name, found := range wantWordIndexes {
		if !found {
			t.Fatalf("missing SQLite vocabulary performance index %s", name)
		}
	}
}

func TestSQLiteVocabularyUsesFileBackedTempStore(t *testing.T) {
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "vocabulary.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
	})
	if err := initSQLiteVocabularyTables(db); err != nil {
		t.Fatal(err)
	}

	var tempStore int
	if err := db.QueryRow(`PRAGMA temp_store`).Scan(&tempStore); err != nil {
		t.Fatal(err)
	}
	const sqliteTempStoreFile = 1
	if tempStore != sqliteTempStoreFile {
		t.Fatalf("PRAGMA temp_store = %d, want FILE (%d)", tempStore, sqliteTempStoreFile)
	}
}

func TestSQLiteVocabularyAICachePersistsHints(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		testVocabWord("station", "A1", 1),
		testVocabWord("ticket", "A1", 2),
		testVocabWord("hotel", "A1", 3),
		testVocabWord("train", "A1", 4),
	})
	word, ok := findVocabWord("en:station")
	if !ok {
		t.Fatal("expected station word")
	}
	key := "test-model|ru|en|learn word|en:station|станция"
	if value, ok, err := sqliteVocabularyAICacheGet(key); err != nil || ok || value != "" {
		t.Fatalf("empty cache lookup = %q, %v, %v; want miss", value, ok, err)
	}
	if err := sqliteVocabularyAICacheSet(key, "hint", word, "ru", "learn word", "test-model", "станция", "place where trains stop"); err != nil {
		t.Fatal(err)
	}
	value, ok, err := sqliteVocabularyAICacheGet(key)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || value != "place where trains stop" {
		t.Fatalf("cached hint = %q, %v; want persisted hint", value, ok)
	}
}

func TestSQLiteVocabularyAITranslationPersistsAndJoins(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:       "en:station",
			Language: "en",
			English:  "station",
			Russian:  "станция",
			Level:    "A1",
			Translations: map[string]string{
				"ru": "станция",
			},
		},
	})
	word, ok := findVocabWord("en:station")
	if !ok {
		t.Fatal("expected station in test vocabulary")
	}
	if got := wordTranslation(word, "es"); got != "" {
		t.Fatalf("wordTranslation before AI cache = %q, want empty", got)
	}
	if err := sqliteVocabularyAITranslationSet(word, "es", "test-model", "prompt", "estación; parada"); err != nil {
		t.Fatalf("sqliteVocabularyAITranslationSet: %v", err)
	}
	value, ok, err := sqliteVocabularyAITranslationGet(word.ID, "es")
	if err != nil || !ok || value != "estación; parada" {
		t.Fatalf("sqliteVocabularyAITranslationGet = %q, %v, %v; want persisted translation", value, ok, err)
	}
	refreshed, ok := findVocabWord("en:station")
	if !ok {
		t.Fatal("expected station after AI translation insert")
	}
	if got := wordTranslation(refreshed, "es"); got != "estación" {
		t.Fatalf("wordTranslation after AI cache = %q, want first persisted value", got)
	}
}

func TestSQLiteVocabularyAIWordSetPersistsAndFindsWord(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		testVocabWord("ticket", "A1", 1),
	})

	inserted, err := sqliteVocabularyAIWordSet(vocabWord{
		Language:      "en",
		English:       "umbrella",
		Russian:       "umbrella ru",
		Translations:  map[string]string{"ru": "umbrella ru", "es": "paraguas"},
		Level:         "A1",
		Topic:         "travel",
		PartOfSpeech:  "noun",
		Source:        "ai",
		FrequencyRank: 0,
	}, "test-model", "prompt")
	if err != nil {
		t.Fatalf("sqliteVocabularyAIWordSet: %v", err)
	}
	if !inserted {
		t.Fatal("sqliteVocabularyAIWordSet inserted = false, want true")
	}

	word, ok := findVocabWord("en:umbrella")
	if !ok {
		t.Fatal("expected AI vocabulary word to be queryable through normal lookup")
	}
	if word.ID != "en:umbrella" || word.Language != "en" || word.Level != "A1" {
		t.Fatalf("unexpected AI word identity: %#v", word)
	}
	if got := wordTranslation(word, "es"); got != "paraguas" {
		t.Fatalf("wordTranslation(AI word, es) = %q, want paraguas", got)
	}
	exists, err := sqliteVocabularyWordExists("en", "umbrella")
	if err != nil {
		t.Fatalf("sqliteVocabularyWordExists: %v", err)
	}
	if !exists {
		t.Fatal("sqliteVocabularyWordExists returned false for persisted AI word")
	}
}

func TestSQLiteVocabularyAIWordSetRejectsDuplicateWord(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		testVocabWord("umbrella", "A1", 1),
	})

	inserted, err := sqliteVocabularyAIWordSet(vocabWord{
		Language:     "en",
		English:      "umbrella",
		Russian:      "umbrella ru",
		Translations: map[string]string{"ru": "umbrella ru"},
		Level:        "A1",
	}, "test-model", "prompt")
	if err != nil {
		t.Fatalf("sqliteVocabularyAIWordSet duplicate: %v", err)
	}
	if inserted {
		t.Fatal("sqliteVocabularyAIWordSet inserted duplicate word, want rejection")
	}
}

func TestSQLiteVocabularyImportSkipsStaleAITranslations(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("VOCABULARY_DIR", dir)
	path := filepath.Join(dir, "vocabulary_words.json")
	if err := os.WriteFile(path, []byte(`[
		{"english":"station","russian":"station ru","translations":{"ru":"station ru"},"level":"A1"},
		{"english":"ticket","russian":"ticket ru","translations":{"ru":"ticket ru"},"level":"A1"}
	]`), 0o644); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "vocabulary.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		configureSQLiteVocabulary(nil)
		_ = db.Close()
	})
	if err := syncSQLiteVocabularyFromJSON(db); err != nil {
		t.Fatal(err)
	}
	configureSQLiteVocabulary(db)

	station, ok := findVocabWord("en:station")
	if !ok {
		t.Fatal("expected station before pruning")
	}
	if err := sqliteVocabularyAITranslationSet(station, "es", "test-model", "prompt", "estacion"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`[
		{"english":"ticket","russian":"ticket ru","translations":{"ru":"ticket ru"},"level":"A1"}
	]`), 0o644); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := importSQLiteVocabularyLanguage(db, "en", path, info); err != nil {
		t.Fatalf("importSQLiteVocabularyLanguage with stale AI translation: %v", err)
	}
	if _, ok := findVocabWord("en:station"); ok {
		t.Fatal("expected pruned seed word to stay absent")
	}
}

func TestSQLiteVocabularyAIWordReplayAfterJSONImport(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("VOCABULARY_DIR", dir)
	path := filepath.Join(dir, "vocabulary_words.json")
	if err := os.WriteFile(path, []byte(`[
		{"english":"ticket","russian":"ticket ru","translations":{"ru":"ticket ru"},"level":"A1"}
	]`), 0o644); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "vocabulary.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		configureSQLiteVocabulary(nil)
		_ = db.Close()
	})
	if err := syncSQLiteVocabularyFromJSON(db); err != nil {
		t.Fatal(err)
	}
	configureSQLiteVocabulary(db)

	inserted, err := sqliteVocabularyAIWordSet(vocabWord{
		Language:     "en",
		English:      "umbrella",
		Russian:      "umbrella ru",
		Translations: map[string]string{"ru": "umbrella ru", "es": "paraguas"},
		Level:        "A1",
	}, "test-model", "prompt")
	if err != nil || !inserted {
		t.Fatalf("sqliteVocabularyAIWordSet = %v, %v; want inserted AI word", inserted, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := importSQLiteVocabularyLanguage(db, "en", path, info); err != nil {
		t.Fatal(err)
	}
	word, ok := findVocabWord("en:umbrella")
	if !ok {
		t.Fatal("expected AI word to replay after JSON import")
	}
	if got := wordTranslation(word, "es"); got != "paraguas" {
		t.Fatalf("replayed AI word translation = %q, want paraguas", got)
	}
}

func TestVocabularyPromptUsesCachedInterfaceTranslationBeforeFallback(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:       "en:station",
			Language: "en",
			English:  "station",
			Russian:  "станция",
			Level:    "A1",
			Translations: map[string]string{
				"ru": "станция",
			},
		},
	})
	word, ok := findVocabWord("en:station")
	if !ok {
		t.Fatal("expected station in test vocabulary")
	}
	if fallback := vocabularyFallbackPrompt(word, "es"); fallback != "станция" {
		t.Fatalf("expected Russian fallback before AI cache, got %q", fallback)
	}
	if err := sqliteVocabularyAITranslationSet(word, "es", "test-model", "prompt", "estación; parada"); err != nil {
		t.Fatalf("sqliteVocabularyAITranslationSet: %v", err)
	}
	got := (*bot)(nil).vocabularyPrompt(context.Background(), userState{InterfaceLanguage: "es", LearningLanguage: "en"}, word, "learn word")
	if got != "estación" {
		t.Fatalf("vocabularyPrompt used %q, want cached Spanish interface prompt", got)
	}
}

func TestVocabularyPromptRejectsDictionaryValueThatLeaksAnswer(t *testing.T) {
	word := vocabWord{
		ID:       "en:creating",
		Language: "en",
		English:  "creating",
		Russian:  "creating",
		Translations: map[string]string{
			"ru": "creating",
			"en": "creating",
			"es": "creador",
		},
		Level:        "A2",
		PartOfSpeech: "adj",
	}
	configureSQLiteVocabularyForTest(t, []vocabWord{word})
	if got := wordTranslation(word, "ru"); got != "" {
		t.Fatalf("wordTranslation leaked target answer as prompt: %q", got)
	}
	if fallback := vocabularyFallbackPrompt(word, "ru"); fallback == "creating" {
		t.Fatalf("vocabularyFallbackPrompt leaked target answer")
	}
	if err := sqliteVocabularyAITranslationSet(word, "ru", "test-model", "prompt", "создание; сотворение"); err != nil {
		t.Fatal(err)
	}
	got := (*bot)(nil).vocabularyPrompt(context.Background(), userState{InterfaceLanguage: "ru", LearningLanguage: "en"}, word, "learn word")
	if got != "создание" {
		t.Fatalf("vocabularyPrompt = %q, want cached Russian translation", got)
	}
}

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

func TestSQLiteNextUnlearnedWordDrawsFromWideRandomPool(t *testing.T) {
	seedVocabularyRandomForTest(t, 20260527)
	words := make([]vocabWord, 0, 180)
	for index := 0; index < 180; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("poolword%03d", index), "A1", index+1))
	}
	configureSQLiteVocabularyForTest(t, words)

	user := userState{LearningLanguage: "en", InterfaceLanguage: "ru", Level: "A1"}
	seen := map[string]bool{}
	for round := 0; round < 80; round++ {
		word, ok := nextUnlearnedWord(user)
		if !ok {
			t.Fatal("expected unlearned word")
		}
		seen[word.ID] = true
	}
	if len(seen) < 36 {
		t.Fatalf("nextUnlearnedWord drew from too narrow a pool: %d unique words in 80 rounds", len(seen))
	}
}

func TestSQLiteNextUnlearnedWordRandomizesSparseEarlyPool(t *testing.T) {
	seedVocabularyRandomForTest(t, 2026052701)
	words := make([]vocabWord, 0, 640)
	for index := 0; index < 24; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("earlya%03d", index), "A1", index+1))
	}
	for index := 0; index < 616; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("lateb%03d", index), "B2", index+1000))
	}
	configureSQLiteVocabularyForTest(t, words)

	user := userState{LearningLanguage: "en", InterfaceLanguage: "ru", Level: "A1"}
	seen := map[string]bool{}
	for round := 0; round < 72; round++ {
		word, ok := nextUnlearnedWord(user)
		if !ok {
			t.Fatal("expected unlearned word")
		}
		if normalizeCEFRLevel(word.Level) != "A1" {
			t.Fatalf("expected A1 word, got %s (%s)", word.English, word.Level)
		}
		seen[word.ID] = true
	}
	if len(seen) < 16 {
		t.Fatalf("sparse SQLite pool fell back to repeated early words: %d unique words in 72 rounds", len(seen))
	}
}

func TestSQLiteNextUnlearnedWordFallbackStaysInsideCEFRBand(t *testing.T) {
	seedVocabularyRandomForTest(t, 2026052801)
	words := []vocabWord{
		testVocabWord("learneda1000", "A1", 1),
		testVocabWord("fallbacka2000", "A2", 2),
	}
	for index := 0; index < 80; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("wrongbandb%03d", index), "B2", index+1000))
		words = append(words, testVocabWord(fmt.Sprintf("wrongbandc%03d", index), "C1", index+2000))
	}
	configureSQLiteVocabularyForTest(t, words)

	user := userState{
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		Level:             "A1",
		LearnedWords: []learnedWordEntry{
			{ID: "en:learneda1000", Language: "en"},
		},
	}
	for round := 0; round < 24; round++ {
		word, ok := nextUnlearnedWord(user)
		if !ok {
			t.Fatal("expected A2 fallback word")
		}
		if normalizeCEFRLevel(word.Level) != "A2" {
			t.Fatalf("expected fallback inside A1/A2 band, got %s (%s)", word.Level, word.English)
		}
	}
}

func TestSQLiteNextUnlearnedWordDoesNotFallForwardWhenBandExhausted(t *testing.T) {
	seedVocabularyRandomForTest(t, 2026052802)
	words := []vocabWord{testVocabWord("learneda1000", "A1", 1)}
	for index := 0; index < 80; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("futureb%03d", index), "B1", index+1000))
		words = append(words, testVocabWord(fmt.Sprintf("futurec%03d", index), "C2", index+2000))
	}
	configureSQLiteVocabularyForTest(t, words)

	user := userState{
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		Level:             "A1",
		LearnedWords: []learnedWordEntry{
			{ID: "en:learneda1000", Language: "en"},
		},
	}
	if word, ok := nextUnlearnedWord(user); ok {
		t.Fatalf("expected no A-band word after A1/A2 exhaustion, got %s (%s)", word.English, word.Level)
	}
}

func TestSQLiteLearnedWordOptionsVaryAcrossReviewRounds(t *testing.T) {
	seedVocabularyRandomForTest(t, 8051)
	words := make([]vocabWord, 0, 96)
	for index := 0; index < 96; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("reviewword%03d", index), "A1", index+1))
	}
	configureSQLiteVocabularyForTest(t, words)

	correct, ok := findVocabWord("en:reviewword000")
	if !ok {
		t.Fatal("expected correct review word")
	}
	learnedIDs := []string{"en:reviewword000", "en:reviewword001", "en:reviewword002", "en:reviewword003"}
	learned := map[string]bool{}
	for _, id := range learnedIDs {
		learned[id] = true
	}
	distractors := map[string]bool{}
	for round := 0; round < 30; round++ {
		options := learnedWordOptions(correct, learnedIDs)
		if len(options) != 4 {
			t.Fatalf("learnedWordOptions returned %d options, want 4: %#v", len(options), options)
		}
		for _, option := range options {
			if option.ID == correct.ID {
				continue
			}
			if learned[option.ID] {
				t.Fatalf("review option reused fixed learned pool item %s: %#v", option.ID, options)
			}
			distractors[option.ID] = true
		}
	}
	if len(distractors) < 18 {
		t.Fatalf("learnedWordOptions used too narrow a distractor pool: %d unique distractors", len(distractors))
	}
}

func TestSQLiteWordOptionsRandomizeSparseEarlyDistractors(t *testing.T) {
	seedVocabularyRandomForTest(t, 2026052702)
	words := make([]vocabWord, 0, 640)
	for index := 0; index < 32; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("optiona%03d", index), "A1", index+1))
	}
	for index := 0; index < 608; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("optionb%03d", index), "B2", index+1000))
	}
	configureSQLiteVocabularyForTest(t, words)

	correct, ok := findVocabWord("en:optiona000")
	if !ok {
		t.Fatal("expected correct word")
	}
	distractors := map[string]bool{}
	for round := 0; round < 36; round++ {
		options := wordOptions(correct)
		if len(options) != 4 {
			t.Fatalf("wordOptions returned %d options, want 4: %#v", len(options), options)
		}
		for _, option := range options {
			if option.ID == correct.ID {
				continue
			}
			if normalizeCEFRLevel(option.Level) != "A1" {
				t.Fatalf("expected A1 distractor, got %s (%s): %#v", option.English, option.Level, options)
			}
			distractors[option.ID] = true
		}
	}
	if len(distractors) < 18 {
		t.Fatalf("sparse SQLite option pool repeated too narrowly: %d unique distractors", len(distractors))
	}
}

func TestSQLiteWordOptionsStayInsideCEFRBandWhenBandHasEnoughDistractors(t *testing.T) {
	seedVocabularyRandomForTest(t, 2026052803)
	words := []vocabWord{testVocabWord("correctb2000", "B2", 1)}
	for index := 0; index < 40; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("bandb%03d", index), "B1", index+10))
		words = append(words, testVocabWord(fmt.Sprintf("wronga%03d", index), "A1", index+1000))
		words = append(words, testVocabWord(fmt.Sprintf("wrongc%03d", index), "C1", index+2000))
	}
	configureSQLiteVocabularyForTest(t, words)

	correct, ok := findVocabWord("en:correctb2000")
	if !ok {
		t.Fatal("expected correct word")
	}
	for round := 0; round < 36; round++ {
		options := wordOptions(correct)
		if len(options) != 4 {
			t.Fatalf("wordOptions returned %d options, want 4: %#v", len(options), options)
		}
		for _, option := range options {
			if option.ID == correct.ID {
				continue
			}
			level := normalizeCEFRLevel(option.Level)
			if level != "B1" && level != "B2" {
				t.Fatalf("expected B-band distractor, got %s (%s): %#v", option.English, option.Level, options)
			}
		}
	}
}

func TestSQLiteWordOptionsDoNotBackfillFromWrongCEFRBand(t *testing.T) {
	seedVocabularyRandomForTest(t, 2026052804)
	words := []vocabWord{
		testVocabWord("correcta1000", "A1", 1),
		testVocabWord("onlya2000", "A2", 2),
	}
	for index := 0; index < 40; index++ {
		words = append(words, testVocabWord(fmt.Sprintf("wrongb%03d", index), "B1", index+100))
		words = append(words, testVocabWord(fmt.Sprintf("wrongc%03d", index), "C1", index+200))
	}
	configureSQLiteVocabularyForTest(t, words)

	correct, ok := findVocabWord("en:correcta1000")
	if !ok {
		t.Fatal("expected correct word")
	}
	options := wordOptions(correct)
	if len(options) != 2 {
		t.Fatalf("wordOptions returned %d options, want only correct plus A-band distractor: %#v", len(options), options)
	}
	for _, option := range options {
		level := normalizeCEFRLevel(option.Level)
		if level != "A1" && level != "A2" {
			t.Fatalf("expected no wrong-band backfill, got %s (%s): %#v", option.English, option.Level, options)
		}
	}
}

func TestLearnedWordsForLanguageNewestFirst(t *testing.T) {
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)
	user := userState{
		LearningLanguage: "en",
		LearnedWords: []learnedWordEntry{
			{ID: "en:old", Language: "en", English: "old", Russian: "\u0441\u0442\u0430\u0440\u043e\u0435", ReviewCorrectCount: learnedWordMasteryThreshold, LearnedAt: now},
			{ID: "en:new", Language: "en", English: "new", Russian: "\u043d\u043e\u0432\u043e\u0435", ReviewCorrectCount: learnedWordMasteryThreshold, LearnedAt: now.Add(time.Hour)},
		},
	}
	words := learnedWordsForLanguage(user)
	if len(words) != 2 {
		t.Fatalf("learnedWordsForLanguage() returned %d words, want 2", len(words))
	}
	if words[0].ID != "en:new" || words[1].ID != "en:old" {
		t.Fatalf("learned words order = %#v, want newest first", words)
	}
}

func TestNextUnlearnedWordStartsAtUserLevel(t *testing.T) {
	user := userState{LearningLanguage: "en", Level: "B2"}
	word, ok := nextUnlearnedWord(user)
	if !ok {
		t.Fatal("expected a word")
	}
	if cefrRank(word.Level) < cefrRank("B2") {
		t.Fatalf("expected B2+ word, got %s (%s)", word.Level, word.English)
	}
}

func TestNextUnlearnedWordSkipsUnsuitableDictionaryEntries(t *testing.T) {
	for _, word := range []string{"pm", "ice-free"} {
		if vocabularyWordSuitableForLearning(vocabWord{English: word, Russian: "перевод"}) {
			t.Fatalf("%q should not be suitable for learning rounds", word)
		}
	}
	if !vocabularyWordSuitableForLearning(vocabWord{English: "apple", Russian: "яблоко"}) {
		t.Fatal("apple should be suitable for learning rounds")
	}
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{English: "pm", Russian: "после полудня", Translations: map[string]string{"ru": "после полудня"}, Level: "A1", FrequencyRank: 1},
		{English: "ice-free", Russian: "свободный ото льда", Translations: map[string]string{"ru": "свободный ото льда"}, Level: "A1", FrequencyRank: 2},
		{English: "apple", Russian: "яблоко", Translations: map[string]string{"ru": "яблоко"}, Level: "A1", FrequencyRank: 3},
	})

	word, ok := nextUnlearnedWord(userState{LearningLanguage: "en", InterfaceLanguage: "ru", Level: "A1"})
	if !ok {
		t.Fatal("expected a suitable word")
	}
	if word.English != "apple" {
		t.Fatalf("nextUnlearnedWord() = %q, want apple", word.English)
	}
}

func TestNextUnlearnedWordUsesSelectedLanguage(t *testing.T) {
	user := userState{LearningLanguage: "es", Level: "A1"}
	word, ok := nextUnlearnedWord(user)
	if !ok {
		t.Fatal("expected a Spanish word")
	}
	if word.Language != "es" {
		t.Fatalf("expected Spanish word, got language %s", word.Language)
	}
}

func TestNextUnlearnedWordSupportsItalian(t *testing.T) {
	user := userState{LearningLanguage: "it", Level: "A1"}
	word, ok := nextUnlearnedWord(user)
	if !ok {
		t.Fatal("expected an Italian word")
	}
	if word.Language != "it" {
		t.Fatalf("expected Italian word, got language %s", word.Language)
	}
	if strings.Contains(word.English, " ") {
		t.Fatalf("expected a single learning word, got %q", word.English)
	}
}

func TestItalianVocabularyContainsOnlySingleLearningWords(t *testing.T) {
	for _, word := range vocabularyForLanguage("it") {
		if strings.TrimSpace(word.English) == "" || strings.TrimSpace(word.Russian) == "" {
			t.Fatalf("expected non-empty Italian vocabulary entry: %#v", word)
		}
		for _, char := range word.English {
			if !unicode.IsLetter(char) && char != '-' && char != '\'' && char != '’' {
				t.Fatalf("expected Italian vocabulary to contain single words only, got %q", word.English)
			}
		}
	}
}

func TestVocabularyUsesInterfaceTranslations(t *testing.T) {
	word, ok := findVocabWord("de:haus")
	if !ok {
		t.Fatal("expected German word Haus in vocabulary")
	}
	if got := wordTranslation(word, "es"); got != "casa" {
		t.Fatalf("wordTranslation(Haus, es) = %q, want casa", got)
	}
	if got := wordTranslation(word, "en"); got != "house" {
		t.Fatalf("wordTranslation(Haus, en) = %q, want house", got)
	}
	if got := wordTranslation(word, "de"); got == "Haus" || got == "" {
		t.Fatalf("same-language interface should use a non-answer prompt, got %q", got)
	}
}

func TestVocabularyKeepsNativeContextVariants(t *testing.T) {
	word := vocabWord{
		ID:       "en:train",
		Language: "en",
		English:  "train",
		Translations: map[string]string{
			"ru": "\u043f\u043e\u0435\u0437\u0434; \u0441\u043e\u0441\u0442\u0430\u0432; \u0442\u0440\u0435\u043d\u0438\u0440\u043e\u0432\u043a\u0430; \u043b\u0438\u0448\u043d\u0435\u0435",
		},
	}
	if prompt := wordTranslation(word, "ru"); prompt != "\u043f\u043e\u0435\u0437\u0434" {
		t.Fatalf("wordTranslation() = %q, want primary prompt", prompt)
	}
	context := wordContext(word, "ru")
	for _, want := range []string{"\u043f\u043e\u0435\u0437\u0434", "\u0441\u043e\u0441\u0442\u0430\u0432", "\u0442\u0440\u0435\u043d\u0438\u0440\u043e\u0432\u043a\u0430"} {
		if want == "\u043f\u043e\u0435\u0437\u0434" {
			continue
		}
		if !strings.Contains(context, want) {
			t.Fatalf("wordContext() = %q, want variant %q", context, want)
		}
	}
	if strings.Contains(context, "\u043b\u0438\u0448\u043d\u0435\u0435") {
		t.Fatalf("wordContext() kept more than three context variants: %q", context)
	}
}

func TestVocabularyUsesDictionaryPromptWords(t *testing.T) {
	word, ok := findVocabWord("de:sein")
	if !ok {
		t.Fatal("expected German word sein in vocabulary")
	}
	prompt := wordTranslation(word, "es")
	if prompt != "ser" {
		t.Fatalf("wordTranslation(sein, es) = %q, want first Spanish prompt word", prompt)
	}
}

func TestVocabularyPromptDoesNotShowEnglishPivotForOtherInterfaces(t *testing.T) {
	word, ok := findVocabWord("de:sein")
	if !ok {
		t.Fatal("expected German word sein in vocabulary")
	}
	if got := wordTranslation(word, "es"); got != "ser" {
		t.Fatalf("wordTranslation(sein, es) = %q, want Spanish prompt", got)
	}
	if block := wordStudyMarkdownBlock(word, "ru"); strings.Contains(strings.ToLower(block), "be") {
		t.Fatalf("Russian study block leaked English pivot: %q", block)
	}
	if context := wordContext(word, "es"); strings.Contains(strings.ToLower(context), "be") {
		t.Fatalf("Spanish context leaked English pivot: %q", context)
	}
}

func TestWordContextShowsNativeAlternatesWithoutRevealingTarget(t *testing.T) {
	word := vocabWord{
		ID:       "en:train",
		Language: "en",
		English:  "train",
		Translations: map[string]string{
			"ru": "поезд; состав; тренировка",
		},
	}
	if got := wordTranslation(word, "ru"); got != "поезд" {
		t.Fatalf("wordTranslation() = %q, want primary prompt", got)
	}
	context := wordContext(word, "ru")
	if !strings.Contains(context, "состав") || !strings.Contains(context, "тренировка") {
		t.Fatalf("wordContext() = %q, want native alternate values", context)
	}
	if strings.Contains(strings.ToLower(context), "train") {
		t.Fatalf("wordContext revealed target answer: %q", context)
	}
}

func TestVocabularyMissingInterfacePromptDoesNotFallbackToRussianOrAnswer(t *testing.T) {
	word := vocabWord{
		ID:       "en:apple",
		Language: "en",
		Russian:  "\u044f\u0431\u043b\u043e\u043a\u043e",
		English:  "apple",
		Translations: map[string]string{
			"en": "apple",
			"ru": "\u044f\u0431\u043b\u043e\u043a\u043e",
		},
	}
	if got := wordTranslation(word, "es"); got != "" {
		t.Fatalf("missing Spanish prompt fell back to %q", got)
	}
}

func TestGeneratedVocabularyPromptsAreLimitedValuesAndScriptSafe(t *testing.T) {
	cyrillicAllowed := map[string]bool{
		"ru": true,
		"uk": true,
		"tg": true,
		"kk": true,
		"ky": true,
		"tt": true,
	}
	for language := range vocabularyFiles {
		language := language
		t.Run(language, func(t *testing.T) {
			if err := forEachVocabularyWord(language, func(word vocabWord) bool {
				if !cyrillicAllowed[language] && containsCyrillic(word.English) {
					t.Fatalf("%s target %s contains Cyrillic: %q", language, word.ID, word.English)
				}
				for code, value := range word.Translations {
					values := dictionaryContextValues(value, "", word.English)
					if len(values) > 3 {
						t.Fatalf("%s translation %s for %s has too many values: %q", language, code, word.ID, value)
					}
					if !cyrillicAllowed[code] && containsCyrillic(value) {
						t.Fatalf("%s translation %s for %s contains Cyrillic: %q", language, code, word.ID, value)
					}
				}
				return true
			}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAllInterfaceLanguagesHaveLearningVocabulary(t *testing.T) {
	if got := len(learningLanguages); got != 35 {
		t.Fatalf("learningLanguages count = %d, want 35", got)
	}
	seen := map[string]bool{}
	for _, language := range learningLanguages {
		code := normalizeLearningLanguage(language.Code)
		seen[code] = true
		if _, ok := vocabularyFiles[code]; !ok {
			t.Fatalf("missing vocabulary file mapping for %s", code)
		}
		count := 0
		if err := forEachVocabularyWord(code, func(word vocabWord) bool {
			count++
			if normalizeLearningLanguage(word.Language) != code {
				t.Fatalf("%s dictionary returned word with language %q", code, word.Language)
			}
			if strings.TrimSpace(word.English) == "" || strings.TrimSpace(word.Russian) == "" {
				t.Fatalf("%s dictionary has unusable entry: %#v", code, word)
			}
			return true
		}); err != nil {
			t.Fatalf("%s dictionary failed to load: %v", code, err)
		}
		if count < 1000 {
			t.Fatalf("%s dictionary has only %d entries", code, count)
		}
	}
	for _, language := range interfaceLanguages() {
		if !seen[normalizeInterfaceLanguage(language.Code)] {
			t.Fatalf("interface language %s is not available as a learning language", language.Code)
		}
	}
}

func TestNewScriptDictionariesDoNotStartWithLatinHeadwords(t *testing.T) {
	scriptChecks := map[string]func(string) bool{
		"ar": func(value string) bool { return regexp.MustCompile(`[\x{0600}-\x{06ff}]`).MatchString(value) },
		"bn": func(value string) bool { return regexp.MustCompile(`[\x{0980}-\x{09ff}]`).MatchString(value) },
		"el": func(value string) bool { return regexp.MustCompile(`[\x{0370}-\x{03ff}]`).MatchString(value) },
		"hi": func(value string) bool { return regexp.MustCompile(`[\x{0900}-\x{097f}]`).MatchString(value) },
		"ta": func(value string) bool { return regexp.MustCompile(`[\x{0b80}-\x{0bff}]`).MatchString(value) },
		"te": func(value string) bool { return regexp.MustCompile(`[\x{0c00}-\x{0c7f}]`).MatchString(value) },
		"th": func(value string) bool { return regexp.MustCompile(`[\x{0e00}-\x{0e7f}]`).MatchString(value) },
	}
	for code, hasScript := range scriptChecks {
		code := code
		hasScript := hasScript
		t.Run(code, func(t *testing.T) {
			checked := 0
			if err := forEachVocabularyWord(code, func(word vocabWord) bool {
				if !hasScript(word.English) {
					t.Fatalf("%s headword lacks expected script: %#v", code, word)
				}
				checked++
				return checked < 100
			}); err != nil {
				t.Fatal(err)
			}
			if checked < 100 {
				t.Fatalf("%s checked only %d words", code, checked)
			}
		})
	}
}

func TestVocabularyContextDoesNotRevealEnglishAnswer(t *testing.T) {
	word, ok := findVocabWord("be")
	if !ok {
		t.Fatal("expected English word be in vocabulary")
	}
	if context := wordContext(word, "en"); dictionaryTextContainsTerm(context, "be") {
		t.Fatalf("English context revealed the answer: %q", context)
	}
}

func TestCommonCoreUsesCuratedLearningWords(t *testing.T) {
	seedVocabularyRandomForTest(t, 20260608)
	for _, language := range []string{"de", "es", "fr", "it", "pl", "pt", "ro", "uk"} {
		word, ok := nextUnlearnedWord(userState{LearningLanguage: language, Level: "A1"})
		if !ok {
			t.Fatalf("expected A1 word for %s", language)
		}
		if normalizeLearningLanguage(word.Language) != language || normalizeCEFRLevel(word.Level) != "A1" || !wordPromptAvailableForUser(word, userState{LearningLanguage: language, Level: "A1"}) {
			t.Fatalf("A1 word for %s = %#v; want usable curated A1 word", language, word)
		}
	}
}

func TestPracticePoolUsesCombinedMasteryCount(t *testing.T) {
	user := userState{
		LearningLanguage: "en",
		LearnedWords: []learnedWordEntry{
			{ID: "en:be", Language: "en", ReviewCorrectCount: 6, SpellingCorrectCount: 3},
			{ID: "en:have", Language: "en", ReviewCorrectCount: 6, SpellingCorrectCount: 4},
		},
	}

	if ids := reviewWordIDs(user); len(ids) != 1 || ids[0] != "en:be" {
		t.Fatalf("reviewWordIDs() = %#v, want only en:be", ids)
	}
	if ids := spellingWordIDs(user); len(ids) != 1 || ids[0] != "en:be" {
		t.Fatalf("spellingWordIDs() = %#v, want only en:be", ids)
	}
	if got := masteredWordCountForLanguage(user, "en"); got != 1 {
		t.Fatalf("masteredWordCountForLanguage() = %d, want 1", got)
	}
}

func TestLevelTestCallbackDoesNotCarryScore(t *testing.T) {
	index, answer, ok := parseLevelTestCallback("lt|3|2")
	if !ok {
		t.Fatal("expected valid callback")
	}
	if index != 3 || answer != 2 {
		t.Fatalf("parseLevelTestCallback() = %d, %d; want 3, 2", index, answer)
	}
	if _, _, ok := parseLevelTestCallback("lt|3|99|2"); ok {
		t.Fatal("expected old score-carrying callback to be rejected")
	}
	index, answer, ok = parseLevelTestCallback("lt|3|-1")
	if !ok {
		t.Fatal("expected unknown answer callback to be valid")
	}
	if index != 3 || answer != -1 {
		t.Fatalf("parseLevelTestCallback() = %d, %d; want 3, -1", index, answer)
	}

	mode := levelTestMode(4, 11)
	gotIndex, gotScore, ok := parseLevelTestMode(mode)
	if !ok {
		t.Fatal("expected valid level test mode")
	}
	if gotIndex != 4 || gotScore != 11 {
		t.Fatalf("parseLevelTestMode() = %d, %d; want 4, 11", gotIndex, gotScore)
	}

	modeWithOrder := levelTestModeWithOrder(2, 7, []int{3, 1, 0, 2})
	gotIndex, gotScore, order, ok := parseLevelTestModeState(modeWithOrder)
	if !ok {
		t.Fatal("expected valid level test mode with order")
	}
	if gotIndex != 2 || gotScore != 7 || !validLevelTestOrder(order, 4) {
		t.Fatalf("parseLevelTestModeState() = %d, %d, %#v; want valid ordered state", gotIndex, gotScore, order)
	}
}

func TestActiveLevelAssessmentQuestionsUseStoredOrder(t *testing.T) {
	baseUser := userState{InterfaceLanguage: "ru", LearningLanguage: "en"}
	base := levelAssessmentQuestionsForUser(baseUser)
	order := make([]int, len(base))
	for index := range order {
		order[index] = index
	}
	order[0], order[1] = 1, 0
	user := userState{
		InterfaceLanguage: "ru",
		LearningLanguage:  "en",
		Mode:              levelTestModeWithOrder(0, 0, order),
	}
	ordered := activeLevelAssessmentQuestionsForUser(user)
	if ordered[0].Question != base[1].Question || ordered[1].Question != base[0].Question {
		t.Fatalf("stored level order was not applied")
	}
}

func TestParseManualLevelCallback(t *testing.T) {
	level, ok := parseManualLevelCallback("level|B2")
	if !ok {
		t.Fatal("expected valid manual level callback")
	}
	if level != "B2" {
		t.Fatalf("parseManualLevelCallback() = %s, want B2", level)
	}
	if _, ok := parseManualLevelCallback("level|D1"); ok {
		t.Fatal("expected invalid CEFR level to be rejected")
	}
}

func TestParseLanguageCallback(t *testing.T) {
	language, ok := parseLanguageCallback("lang|it")
	if !ok {
		t.Fatal("expected valid language callback")
	}
	if language != "it" {
		t.Fatalf("parseLanguageCallback() = %s, want it", language)
	}
	language, ok = parseLanguageCallback("lang|pt")
	if !ok {
		t.Fatal("expected Portuguese language callback to be valid")
	}
	if language != "pt" {
		t.Fatalf("parseLanguageCallback() = %s, want pt", language)
	}
	language, ok = parseLanguageCallback("lang|ru")
	if !ok {
		t.Fatal("expected Russian language callback to be valid")
	}
	if language != "ru" {
		t.Fatalf("parseLanguageCallback() = %s, want ru", language)
	}
}

func TestLevelAssessmentUsesSelectedLanguage(t *testing.T) {
	for _, language := range learningLanguages {
		if got := len(levelAssessmentQuestionsForLanguage(language.Code)); got != 36 {
			t.Fatalf("expected %s assessment to have 36 questions, got %d", language.Code, got)
		}
	}
	if level := levelFromAssessmentScore(levelAssessmentMaxScore("fr"), levelAssessmentMaxScore("fr")); level != "C2" {
		t.Fatalf("expected max French score to map to C2, got %s", level)
	}
}

func TestLevelAssessmentForNonRussianInterfaceAvoidsRussianScaffolding(t *testing.T) {
	for _, language := range learningLanguages {
		user := userState{InterfaceLanguage: "en", LearningLanguage: language.Code}
		questions := levelAssessmentQuestionsForUser(user)
		if len(questions) == 0 {
			t.Fatalf("expected exam-style level assessment for %s with non-Russian interface", language.Code)
		}
		copy := systemUI(user)
		for index, question := range questions {
			text := fmt.Sprintf(copy.LevelQuestionText, index+1, len(questions), question.Question) + "\n\n" + fmt.Sprintf(copy.LevelCurrentScore, 0)
			if containsRussianScaffolding(text) {
				t.Fatalf("expected localized level question chrome for %s at %d, got %q", language.Code, index, text)
			}
			if !isCyrillicLearningLanguage(language.Code) && containsCyrillic(text) {
				t.Fatalf("expected no Cyrillic payload for %s at %d, got %q", language.Code, index, text)
			}
			for _, option := range question.Options {
				if !isCyrillicLearningLanguage(language.Code) && containsCyrillic(option) {
					t.Fatalf("expected non-Cyrillic level option for %s at %d, got %q", language.Code, index, option)
				}
			}
		}
	}
}

func TestLevelAssessmentUsesInterfaceLanguagePrompt(t *testing.T) {
	user := userState{InterfaceLanguage: "es", LearningLanguage: "de"}
	text := levelQuestionTextForUser(user, 0, 0)
	if !strings.Contains(text, "Pregunta 1/36") || !strings.Contains(text, "Elige la forma correcta") {
		t.Fatalf("expected Spanish level-test chrome, got %q", text)
	}
}

func TestMeaningQuestionsDoNotEchoCorrectAnswer(t *testing.T) {
	for _, language := range learningLanguages {
		for index, question := range levelAssessmentQuestionsForLanguage(language.Code) {
			if !isMeaningAssessmentPrompt(question.Question) {
				continue
			}
			payload := normalizeAssessmentEcho(assessmentPromptPayload(question.Question))
			correct := normalizeAssessmentEcho(question.Options[question.CorrectIndex])
			if payload != "" && strings.Contains(correct, payload) {
				t.Fatalf("expected meaning question for %s at %d not to echo %q in correct option %q", language.Code, index, question.Question, question.Options[question.CorrectIndex])
			}
		}
	}
}

func TestSanitizeGeneratedVocabularyWordRejectsTargetAsTranslation(t *testing.T) {
	word, ok := sanitizeGeneratedVocabularyWord("en", generatedVocabularyWord{
		Word:        "doctor",
		Translation: "doctor",
		Russian:     "doctor",
		Level:       "A1",
	}, "ru", "A1", nil)
	if ok {
		t.Fatalf("sanitizeGeneratedVocabularyWord accepted target as translation: %#v", word)
	}
}

func TestSanitizeGeneratedVocabularyWordRejectsUnsuitableTarget(t *testing.T) {
	word, ok := sanitizeGeneratedVocabularyWord("en", generatedVocabularyWord{
		Word:        "train station",
		Translation: "estacion",
		Russian:     "station ru",
		Level:       "A1",
	}, "es", "A1", nil)
	if ok {
		t.Fatalf("sanitizeGeneratedVocabularyWord accepted phrase target: %#v", word)
	}
}

func TestSanitizeGeneratedVocabularyWordRejectsRussianThatLeaksTarget(t *testing.T) {
	word, ok := sanitizeGeneratedVocabularyWord("en", generatedVocabularyWord{
		Word:        "doctor",
		Translation: "medico",
		Russian:     "doctor",
		Level:       "A1",
	}, "es", "A1", nil)
	if ok {
		t.Fatalf("sanitizeGeneratedVocabularyWord accepted Russian target leak: %#v", word)
	}
}

func TestSanitizeGeneratedVocabularyWordNormalizesLevelAndID(t *testing.T) {
	word, ok := sanitizeGeneratedVocabularyWord("en", generatedVocabularyWord{
		Word:         " Umbrella ",
		Translation:  "paraguas",
		Russian:      "umbrella ru",
		Level:        "b1",
		Topic:        " Travel ",
		PartOfSpeech: " Noun ",
	}, "es", "A2", map[string]bool{})
	if !ok {
		t.Fatal("sanitizeGeneratedVocabularyWord rejected valid generated card")
	}
	if word.ID != "en:umbrella" || word.Language != "en" || word.English != "Umbrella" {
		t.Fatalf("unexpected generated word identity: %#v", word)
	}
	if word.Level != "B1" {
		t.Fatalf("generated word level = %q, want B1", word.Level)
	}
	if got := wordTranslation(word, "es"); got != "paraguas" {
		t.Fatalf("generated word Spanish translation = %q, want paraguas", got)
	}
}

func TestVocabularyPromptInLanguageDoesNotFallbackToRussian(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:       "en:station",
			Language: "en",
			English:  "station",
			Russian:  "station ru",
			Level:    "A1",
			Translations: map[string]string{
				"ru": "station ru",
			},
		},
	})
	word, ok := findVocabWord("en:station")
	if !ok {
		t.Fatal("expected station word")
	}
	got := (*bot)(nil).vocabularyPromptInLanguage(context.Background(), userState{LearningLanguage: "en", InterfaceLanguage: "es"}, word, "es", "learn word")
	if got != "" {
		t.Fatalf("vocabularyPromptInLanguage fallback = %q, want empty Spanish prompt", got)
	}
}

func TestVocabularyGenerationPromptForbidsDuplicatesAndRequiresJSON(t *testing.T) {
	messages := vocabularyGenerationPrompt(
		userState{LearningLanguage: "en", InterfaceLanguage: "ru", Level: "A1"},
		interfaceLanguageByCode("ru"),
		[]string{"doctor"},
	)
	joined := ""
	for _, message := range messages {
		joined += message.Content + "\n"
	}
	for _, want := range []string{"Return strict JSON only", "doctor", `"word"`, `"translation"`, `"russian"`} {
		if !strings.Contains(joined, want) {
			t.Fatalf("generation prompt missing %q:\n%s", want, joined)
		}
	}
}

func TestVocabularyGenerationPromptLanguageFallsBackWhenInterfaceMatchesLearning(t *testing.T) {
	if got := vocabularyGenerationPromptLanguage(userState{LearningLanguage: "en", InterfaceLanguage: "en"}); got.Code != "ru" {
		t.Fatalf("English learner with English UI prompt language = %s, want ru", got.Code)
	}
	if got := vocabularyGenerationPromptLanguage(userState{LearningLanguage: "ru", InterfaceLanguage: "ru"}); got.Code != "en" {
		t.Fatalf("Russian learner with Russian UI prompt language = %s, want en", got.Code)
	}
	if got := vocabularyGenerationPromptLanguage(userState{LearningLanguage: "de", InterfaceLanguage: "es"}); got.Code != "es" {
		t.Fatalf("German learner with Spanish UI prompt language = %s, want es", got.Code)
	}
}

func TestVocabularyGenerationForbiddenWordsIncludesLearnedAndExtra(t *testing.T) {
	forbidden, visible := vocabularyGenerationForbiddenWords(userState{
		LearningLanguage: "en",
		LearnedWords: []learnedWordEntry{
			{ID: "en:doctor", Language: "en", English: "doctor"},
			{ID: "de:haus", Language: "de", English: "Haus"},
		},
	}, []string{"umbrella"})
	for _, key := range []string{"en:doctor", "doctor", "en:umbrella", "umbrella"} {
		if !forbidden[key] {
			t.Fatalf("forbidden[%q] = false, want true; visible=%#v", key, visible)
		}
	}
	if forbidden["de:haus"] || forbidden["haus"] {
		t.Fatalf("forbidden includes learned word from another language: %#v", forbidden)
	}
}

func TestSanitizeVocabularyTranslationRejectsTargetWord(t *testing.T) {
	word := vocabWord{ID: "en:okay", Language: "en", English: "okay", Russian: "хорошо"}
	got := sanitizeVocabularyTranslation(`{"translation":"okay; fine; all right"}`, word)
	if got != "fine; all right" {
		t.Fatalf("sanitizeVocabularyTranslation kept target word: got %q", got)
	}
}

func TestVocabularyTranslationPromptForbidsSameSurfaceForm(t *testing.T) {
	word := vocabWord{ID: "en:okay", Language: "en", English: "okay", Russian: "хорошо", Level: "A1"}
	messages := vocabularyTranslationPrompt(word, learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "learn words")
	joined := ""
	for _, message := range messages {
		joined += message.Content + "\n"
	}
	for _, want := range []string{"Never return the hidden target-language word itself", "must not equal or repeat the hidden target word"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("translation prompt missing %q:\n%s", want, joined)
		}
	}
}

func containsCyrillic(text string) bool {
	for _, r := range text {
		if unicode.In(r, unicode.Cyrillic) {
			return true
		}
	}
	return false
}

func containsRussianScaffolding(text string) bool {
	for _, marker := range []string{"Выбери", "Заполни", "Что значит", "Как сказать"} {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return false
}

func isMeaningAssessmentPrompt(question string) bool {
	return strings.HasPrefix(question, "Что значит ") || strings.HasPrefix(question, "What does ")
}

func normalizeAssessmentEcho(text string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
