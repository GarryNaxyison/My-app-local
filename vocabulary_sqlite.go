package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type sqliteVocabularyRuntime struct {
	mu      sync.RWMutex
	db      *sql.DB
	cacheMu sync.RWMutex
	idCache map[string][]string
}

var vocabularySQLiteRuntime sqliteVocabularyRuntime
var vocabularyJSONPersistMu sync.Mutex

func initializeVocabularyRuntime(cfg config) (*sql.DB, error) {
	db, err := openSQLiteDatabase("VOCABULARY_DATABASE_PATH", cfg.VocabularyDatabasePath)
	if err != nil {
		return nil, err
	}
	if err := initSQLiteVocabularyTables(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	configureSQLiteVocabulary(db)
	log.Printf("SQLite vocabulary opened: %s", cfg.VocabularyDatabasePath)
	go func() {
		if err := syncSQLiteVocabularyFromJSON(db); err != nil {
			log.Printf("failed to sync SQLite vocabulary in background: %v", err)
			return
		}
		log.Printf("SQLite vocabulary ready: %s", cfg.VocabularyDatabasePath)
	}()
	return db, nil
}

func configureSQLiteVocabulary(db *sql.DB) {
	vocabularySQLiteRuntime.mu.Lock()
	defer vocabularySQLiteRuntime.mu.Unlock()
	vocabularySQLiteRuntime.db = db
	vocabularySQLiteRuntime.cacheMu.Lock()
	vocabularySQLiteRuntime.idCache = map[string][]string{}
	vocabularySQLiteRuntime.cacheMu.Unlock()
}

func currentSQLiteVocabularyDB() *sql.DB {
	vocabularySQLiteRuntime.mu.RLock()
	defer vocabularySQLiteRuntime.mu.RUnlock()
	return vocabularySQLiteRuntime.db
}

func invalidateSQLiteVocabularyIDCache() {
	vocabularySQLiteRuntime.cacheMu.Lock()
	defer vocabularySQLiteRuntime.cacheMu.Unlock()
	vocabularySQLiteRuntime.idCache = map[string][]string{}
}

func initSQLiteVocabularyTables(db *sql.DB) error {
	if db == nil {
		return nil
	}
	statements := []string{
		`PRAGMA journal_mode=WAL`,
		`PRAGMA synchronous=NORMAL`,
		`PRAGMA busy_timeout=5000`,
		`PRAGMA temp_store=MEMORY`,
		`PRAGMA foreign_keys=ON`,
		`CREATE TABLE IF NOT EXISTS vocabulary_words (
			id TEXT PRIMARY KEY,
			language TEXT NOT NULL,
			word TEXT NOT NULL,
			russian TEXT NOT NULL,
			context TEXT NOT NULL DEFAULT '',
			level TEXT NOT NULL DEFAULT '',
			topic TEXT NOT NULL DEFAULT '',
			part_of_speech TEXT NOT NULL DEFAULT '',
			source TEXT NOT NULL DEFAULT '',
			frequency_rank INTEGER NOT NULL DEFAULT 0,
			position INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS vocabulary_translations (
			word_id TEXT NOT NULL,
			language TEXT NOT NULL,
			text TEXT NOT NULL,
			PRIMARY KEY (word_id, language),
			FOREIGN KEY (word_id) REFERENCES vocabulary_words(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS vocabulary_sources (
			language TEXT PRIMARY KEY,
			path TEXT NOT NULL,
			size_bytes INTEGER NOT NULL,
			mod_time_unix INTEGER NOT NULL,
			word_count INTEGER NOT NULL,
			imported_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS vocabulary_ai_cache (
			cache_key TEXT PRIMARY KEY,
			kind TEXT NOT NULL,
			word_id TEXT NOT NULL,
			learning_language TEXT NOT NULL,
			interface_language TEXT NOT NULL,
			mode TEXT NOT NULL,
			model TEXT NOT NULL,
			prompt TEXT NOT NULL,
			value TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS vocabulary_ai_translations (
			word_id TEXT NOT NULL,
			learning_language TEXT NOT NULL,
			target_language TEXT NOT NULL,
			source_word TEXT NOT NULL,
			translation TEXT NOT NULL,
			model TEXT NOT NULL,
			prompt TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			PRIMARY KEY (word_id, target_language)
		)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func ensureSQLiteVocabularyIndexes(db *sql.DB) error {
	if db == nil {
		return nil
	}
	statements := []string{
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_words_language_position ON vocabulary_words(language, position)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_words_language_level_position ON vocabulary_words(language, level, position)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_words_language_rank ON vocabulary_words(language, frequency_rank)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_words_language_level_rank_position ON vocabulary_words(language, level, frequency_rank, position)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_words_language_position_id ON vocabulary_words(language, position, id)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_translations_language ON vocabulary_translations(language)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_translations_language_word ON vocabulary_translations(language, word_id)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_ai_cache_word_kind ON vocabulary_ai_cache(word_id, kind, interface_language)`,
		`CREATE INDEX IF NOT EXISTS idx_vocabulary_ai_translations_language ON vocabulary_ai_translations(target_language, learning_language)`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func syncSQLiteVocabularyFromJSON(db *sql.DB) error {
	if db == nil {
		return nil
	}
	if err := initSQLiteVocabularyTables(db); err != nil {
		return err
	}
	languages := make([]string, 0, len(vocabularyFiles))
	for language := range vocabularyFiles {
		languages = append(languages, language)
	}
	sort.Strings(languages)
	for _, language := range languages {
		if err := syncSQLiteVocabularyLanguage(db, language); err != nil {
			return err
		}
	}
	if err := ensureSQLiteVocabularyIndexes(db); err != nil {
		return err
	}
	invalidateSQLiteVocabularyIDCache()
	return nil
}

func syncSQLiteVocabularyLanguage(db *sql.DB, language string) error {
	language = normalizeLearningLanguage(language)
	path, ok := vocabularyPath(language)
	if !ok {
		return nil
	}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if fresh, err := sqliteVocabularySourceFresh(db, language, path, info); err != nil {
		return err
	} else if fresh {
		return nil
	}
	log.Printf("Importing %s vocabulary into SQLite from %s", language, path)
	return importSQLiteVocabularyLanguage(db, language, path, info)
}

func sqliteVocabularySourceFresh(db *sql.DB, language string, path string, info os.FileInfo) (bool, error) {
	var sizeBytes int64
	var modTimeUnix int64
	var wordCount int
	err := db.QueryRow(
		`SELECT size_bytes, mod_time_unix, word_count FROM vocabulary_sources WHERE language = ?`,
		language,
	).Scan(&sizeBytes, &modTimeUnix, &wordCount)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if sizeBytes != info.Size() || modTimeUnix != info.ModTime().Unix() || wordCount <= 0 {
		return false, nil
	}
	var storedCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM vocabulary_words WHERE language = ?`, language).Scan(&storedCount); err != nil {
		return false, err
	}
	return storedCount == wordCount, nil
}

func importSQLiteVocabularyLanguage(db *sql.DB, language string, path string, info os.FileInfo) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if delimiter, ok := token.(json.Delim); !ok || delimiter != '[' {
		return fmt.Errorf("read %s: expected JSON array", path)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`DELETE FROM vocabulary_translations
		WHERE word_id IN (SELECT id FROM vocabulary_words WHERE language = ?)`,
		language,
	); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM vocabulary_words WHERE language = ?`, language); err != nil {
		return err
	}

	wordStmt, err := tx.Prepare(`INSERT INTO vocabulary_words
		(id, language, word, russian, context, level, topic, part_of_speech, source, frequency_rank, position)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer wordStmt.Close()

	translationStmt, err := tx.Prepare(`INSERT OR REPLACE INTO vocabulary_translations
		(word_id, language, text) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer translationStmt.Close()

	seen := map[string]bool{}
	position := 0
	for decoder.More() {
		var raw vocabWord
		if err := decoder.Decode(&raw); err != nil {
			return fmt.Errorf("decode %s: %w", path, err)
		}
		word, ok := normalizeVocabularyWord(language, raw)
		if !ok || seen[word.ID] {
			continue
		}
		seen[word.ID] = true
		if _, err := wordStmt.Exec(
			word.ID,
			word.Language,
			word.English,
			word.Russian,
			word.Context,
			strings.TrimSpace(word.Level),
			strings.TrimSpace(word.Topic),
			strings.TrimSpace(word.PartOfSpeech),
			strings.TrimSpace(word.Source),
			word.FrequencyRank,
			position,
		); err != nil {
			return err
		}
		for code, value := range word.Translations {
			code = normalizeInterfaceLanguage(code)
			value = cleanDictionaryDisplay(value)
			if code == "" || value == "" {
				continue
			}
			if _, err := translationStmt.Exec(word.ID, code, value); err != nil {
				return err
			}
		}
		position++
	}
	if _, err := tx.Exec(
		`INSERT OR REPLACE INTO vocabulary_translations (word_id, language, text)
		SELECT word_id, target_language, translation
		FROM vocabulary_ai_translations
		WHERE learning_language = ? AND translation <> ''`,
		language,
	); err != nil {
		return err
	}

	if _, err := tx.Exec(
		`INSERT INTO vocabulary_sources (language, path, size_bytes, mod_time_unix, word_count, imported_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(language) DO UPDATE SET
			path = excluded.path,
			size_bytes = excluded.size_bytes,
			mod_time_unix = excluded.mod_time_unix,
			word_count = excluded.word_count,
			imported_at = excluded.imported_at`,
		language,
		path,
		info.Size(),
		info.ModTime().Unix(),
		position,
		time.Now().UTC().Format(time.RFC3339),
	); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Imported %s vocabulary into SQLite: %d words", language, position)
	return nil
}

func sqliteVocabularyHasLanguage(db *sql.DB, language string) (bool, error) {
	var exists int
	err := db.QueryRow(`SELECT 1 FROM vocabulary_words WHERE language = ? LIMIT 1`, normalizeLearningLanguage(language)).Scan(&exists)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func sqliteForEachVocabularyWord(language string, visit func(vocabWord) bool) (bool, error) {
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return false, nil
	}
	language = normalizeLearningLanguage(language)
	hasLanguage, err := sqliteVocabularyHasLanguage(db, language)
	if err != nil {
		return true, err
	}
	if !hasLanguage {
		return false, nil
	}
	rows, err := db.Query(sqliteVocabularySelectSQL(`
		WHERE w.language = ?
		ORDER BY w.position, t.language`), language)
	if err != nil {
		return true, err
	}
	defer rows.Close()
	return true, scanSQLiteVocabularyRows(rows, visit)
}

func sqliteFindVocabWords(ids []string) (map[string]vocabWord, bool, error) {
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return nil, false, nil
	}

	requested := map[string]bool{}
	languages := map[string]bool{}
	candidateToOriginals := map[string][]string{}
	candidatesByLanguage := map[string][]string{}
	for _, originalID := range ids {
		originalID = strings.TrimSpace(originalID)
		if originalID == "" || requested[originalID] {
			continue
		}
		requested[originalID] = true
		language := languageFromVocabID(originalID)
		languages[language] = true
		legacy := legacyVocabID(originalID)
		candidates := []string{originalID}
		prefixed := makeVocabID(language, legacy)
		if prefixed != originalID {
			candidates = append(candidates, prefixed)
		}
		for _, candidate := range candidates {
			if candidate == "" {
				continue
			}
			candidateToOriginals[candidate] = append(candidateToOriginals[candidate], originalID)
			candidatesByLanguage[language] = append(candidatesByLanguage[language], candidate)
		}
	}
	if len(requested) == 0 {
		return map[string]vocabWord{}, true, nil
	}

	for language := range languages {
		hasLanguage, err := sqliteVocabularyHasLanguage(db, language)
		if err != nil {
			return nil, true, err
		}
		if !hasLanguage {
			return nil, false, nil
		}
	}

	result := map[string]vocabWord{}
	for _, candidates := range candidatesByLanguage {
		uniqueCandidates := uniqueStrings(candidates)
		if len(uniqueCandidates) == 0 {
			continue
		}
		words, err := sqliteVocabularyWordsByIDs(db, uniqueCandidates)
		if err != nil {
			return nil, true, err
		}
		for wordID, word := range words {
			for _, originalID := range candidateToOriginals[wordID] {
				if _, exists := result[originalID]; !exists {
					result[originalID] = word
				}
			}
		}
	}
	return result, true, nil
}

func sqliteVocabularyWordsByIDs(db *sql.DB, ids []string) (map[string]vocabWord, error) {
	result := map[string]vocabWord{}
	const batchSize = 400
	for start := 0; start < len(ids); start += batchSize {
		end := start + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		batch := ids[start:end]
		args := make([]any, 0, len(batch))
		for _, id := range batch {
			args = append(args, id)
		}
		rows, err := db.Query(sqliteVocabularySelectSQL(`
			WHERE w.id IN (`+sqlitePlaceholders(len(batch))+`)
			ORDER BY w.position, t.language`), args...)
		if err != nil {
			return nil, err
		}
		if err := scanSQLiteVocabularyRows(rows, func(word vocabWord) bool {
			result[word.ID] = word
			return true
		}); err != nil {
			rows.Close()
			return nil, err
		}
		if err := rows.Close(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func sqliteVocabularyForLanguage(language string) ([]vocabWord, bool, error) {
	words := []vocabWord{}
	used, err := sqliteForEachVocabularyWord(language, func(word vocabWord) bool {
		words = append(words, word)
		return true
	})
	if !used || err != nil {
		return nil, used, err
	}
	return words, true, nil
}

func sqliteVocabularyAICacheGet(cacheKey string) (string, bool, error) {
	cacheKey = strings.TrimSpace(cacheKey)
	if cacheKey == "" {
		return "", false, nil
	}
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return "", false, nil
	}
	var value string
	err := db.QueryRow(`SELECT value FROM vocabulary_ai_cache WHERE cache_key = ?`, cacheKey).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	value = cleanDictionaryDisplay(value)
	if value == "" {
		return "", false, nil
	}
	return value, true, nil
}

func sqliteVocabularyAICacheSet(cacheKey string, kind string, word vocabWord, interfaceLanguage string, mode string, model string, prompt string, value string) error {
	cacheKey = strings.TrimSpace(cacheKey)
	value = cleanDictionaryDisplay(value)
	if cacheKey == "" || value == "" {
		return nil
	}
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.Exec(
		`INSERT INTO vocabulary_ai_cache
			(cache_key, kind, word_id, learning_language, interface_language, mode, model, prompt, value, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(cache_key) DO UPDATE SET
			value = excluded.value,
			prompt = excluded.prompt,
			model = excluded.model,
			updated_at = excluded.updated_at`,
		cacheKey,
		strings.TrimSpace(kind),
		strings.TrimSpace(word.ID),
		normalizeLearningLanguage(word.Language),
		normalizeInterfaceLanguage(interfaceLanguage),
		strings.TrimSpace(mode),
		strings.TrimSpace(model),
		cleanDictionaryDisplay(prompt),
		value,
		now,
		now,
	)
	return err
}

func sqliteVocabularyAITranslationGet(wordID string, targetLanguage string) (string, bool, error) {
	wordID = strings.TrimSpace(wordID)
	targetLanguage = normalizeInterfaceLanguage(targetLanguage)
	if wordID == "" || targetLanguage == "" {
		return "", false, nil
	}
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return "", false, nil
	}
	var value string
	err := db.QueryRow(
		`SELECT translation FROM vocabulary_ai_translations WHERE word_id = ? AND target_language = ?`,
		wordID,
		targetLanguage,
	).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	value = cleanDictionaryDisplay(value)
	if value == "" {
		return "", false, nil
	}
	return value, true, nil
}

func sqliteVocabularyAITranslationSet(word vocabWord, targetLanguage string, model string, prompt string, value string) error {
	targetLanguage = normalizeInterfaceLanguage(targetLanguage)
	value = cleanDictionaryDisplay(value)
	if strings.TrimSpace(word.ID) == "" || targetLanguage == "" || value == "" {
		return nil
	}
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(
		`INSERT INTO vocabulary_ai_translations
			(word_id, learning_language, target_language, source_word, translation, model, prompt, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(word_id, target_language) DO UPDATE SET
			translation = excluded.translation,
			model = excluded.model,
			prompt = excluded.prompt,
			updated_at = excluded.updated_at`,
		strings.TrimSpace(word.ID),
		normalizeLearningLanguage(word.Language),
		targetLanguage,
		cleanDictionaryDisplay(word.English),
		value,
		strings.TrimSpace(model),
		cleanDictionaryDisplay(prompt),
		now,
		now,
	); err != nil {
		return err
	}
	if _, err := tx.Exec(
		`INSERT OR REPLACE INTO vocabulary_translations (word_id, language, text)
		SELECT id, ?, ? FROM vocabulary_words WHERE id = ?`,
		targetLanguage,
		value,
		strings.TrimSpace(word.ID),
	); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	invalidateSQLiteVocabularyIDCache()
	if err := persistVocabularyTranslationToJSON(word, targetLanguage, value); err != nil {
		return err
	}
	return nil
}

func persistVocabularyTranslationToJSON(word vocabWord, targetLanguage string, value string) error {
	language := normalizeLearningLanguage(word.Language)
	targetLanguage = normalizeInterfaceLanguage(targetLanguage)
	value = cleanDictionaryDisplay(value)
	if language == "" || targetLanguage == "" || value == "" || targetLanguage == language {
		return nil
	}
	path, ok := vocabularyPath(language)
	if !ok {
		return nil
	}

	vocabularyJSONPersistMu.Lock()
	defer vocabularyJSONPersistMu.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read vocabulary json %s: %w", path, err)
	}
	var rows []vocabWord
	if err := json.Unmarshal(data, &rows); err != nil {
		return fmt.Errorf("parse vocabulary json %s: %w", path, err)
	}
	targetID := strings.TrimSpace(word.ID)
	if targetID == "" {
		targetID = makeVocabID(language, word.English)
	}
	changed := false
	for index := range rows {
		rowID := strings.TrimSpace(rows[index].ID)
		if rowID == "" {
			rowID = makeVocabID(language, rows[index].English)
		}
		if rowID != targetID && !sameDictionaryText(rows[index].English, word.English) {
			continue
		}
		if rows[index].Translations == nil {
			rows[index].Translations = map[string]string{}
		}
		if cleanDictionaryDisplay(rows[index].Translations[targetLanguage]) != value {
			rows[index].Translations[targetLanguage] = value
			changed = true
		}
		if targetLanguage == "ru" && cleanDictionaryDisplay(rows[index].Russian) != value {
			rows[index].Russian = value
			changed = true
		}
		if rows[index].ID == "" && strings.TrimSpace(word.ID) != "" {
			rows[index].ID = targetID
			changed = true
		}
		if rows[index].Language == "" && language != "" {
			rows[index].Language = language
			changed = true
		}
		break
	}
	if !changed {
		return nil
	}
	payload, err := json.Marshal(rows)
	if err != nil {
		return fmt.Errorf("encode vocabulary json %s: %w", path, err)
	}
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return fmt.Errorf("write vocabulary json %s: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("replace vocabulary json %s: %w", path, err)
	}
	return nil
}

func sqliteNextUnlearnedWord(user userState) (vocabWord, bool, bool, error) {
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return vocabWord{}, false, false, nil
	}
	language := normalizeLearningLanguage(user.LearningLanguage)
	hasLanguage, err := sqliteVocabularyHasLanguage(db, language)
	if err != nil {
		return vocabWord{}, false, true, err
	}
	if !hasLanguage {
		return vocabWord{}, false, false, nil
	}

	learnedIDs := sqliteLearnedVocabularyIDs(user, language)
	id, ok, err := sqliteRandomUnlearnedWordID(db, language, user.Level, learnedIDs, true)
	if err != nil {
		return vocabWord{}, false, true, err
	}
	if !ok {
		id, ok, err = sqliteRandomUnlearnedWordID(db, language, user.Level, learnedIDs, false)
		if err != nil {
			return vocabWord{}, false, true, err
		}
	}
	if !ok {
		return vocabWord{}, false, true, nil
	}
	words, err := sqliteVocabularyWordsByIDs(db, []string{id})
	if err != nil {
		return vocabWord{}, false, true, err
	}
	word, ok := words[id]
	return word, ok, true, nil
}

func sqliteWordOptions(correct vocabWord) ([]vocabWord, bool, error) {
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return nil, false, nil
	}
	language := normalizeLearningLanguage(correct.Language)
	hasLanguage, err := sqliteVocabularyHasLanguage(db, language)
	if err != nil {
		return nil, true, err
	}
	if !hasLanguage {
		return nil, false, nil
	}

	options := []vocabWord{correct}
	exclude := map[string]bool{correct.ID: true}
	candidates, err := sqliteSampleVocabularyOptionCandidates(db, language, 3, correct.Level, exclude)
	if err != nil {
		return nil, true, err
	}
	options = append(options, candidates...)
	for _, candidate := range candidates {
		exclude[candidate.ID] = true
	}
	randomWords(options)
	return options, true, nil
}

func sqliteLearnedWordOptions(correct vocabWord, learnedIDs []string) ([]vocabWord, bool, error) {
	db := currentSQLiteVocabularyDB()
	if db == nil {
		return nil, false, nil
	}
	language := normalizeLearningLanguage(correct.Language)
	hasLanguage, err := sqliteVocabularyHasLanguage(db, language)
	if err != nil {
		return nil, true, err
	}
	if !hasLanguage {
		return nil, false, nil
	}

	options := []vocabWord{correct}
	exclude := map[string]bool{correct.ID: true, legacyVocabID(correct.ID): true}
	for _, id := range learnedIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if !strings.Contains(id, ":") {
			id = makeVocabID(language, id)
		}
		exclude[id] = true
		exclude[legacyVocabID(id)] = true
	}

	candidates, err := sqliteSampleVocabularyOptionCandidates(db, language, 3, correct.Level, exclude)
	if err != nil {
		return nil, true, err
	}
	options = append(options, candidates...)
	for _, candidate := range candidates {
		exclude[candidate.ID] = true
		exclude[legacyVocabID(candidate.ID)] = true
	}
	randomWords(options)
	return options, true, nil
}

func sqliteVocabularySelectSQL(suffix string) string {
	return `SELECT
			w.id,
			w.language,
			w.word,
			w.russian,
			w.context,
			w.level,
			w.topic,
			w.part_of_speech,
			w.source,
			w.frequency_rank,
			t.language,
			t.text
		FROM vocabulary_words w
		LEFT JOIN vocabulary_translations t ON t.word_id = w.id
	` + suffix
}

func scanSQLiteVocabularyRows(rows *sql.Rows, visit func(vocabWord) bool) error {
	var current vocabWord
	haveCurrent := false
	keepGoing := true
	flush := func() bool {
		if !haveCurrent {
			return true
		}
		return visit(current)
	}
	for rows.Next() {
		var word vocabWord
		var translationLanguage sql.NullString
		var translationText sql.NullString
		if err := rows.Scan(
			&word.ID,
			&word.Language,
			&word.English,
			&word.Russian,
			&word.Context,
			&word.Level,
			&word.Topic,
			&word.PartOfSpeech,
			&word.Source,
			&word.FrequencyRank,
			&translationLanguage,
			&translationText,
		); err != nil {
			return err
		}
		if !haveCurrent || word.ID != current.ID {
			if !flush() {
				keepGoing = false
				break
			}
			word.Translations = map[string]string{}
			current = word
			haveCurrent = true
		}
		if translationLanguage.Valid && translationText.Valid {
			current.Translations[normalizeInterfaceLanguage(translationLanguage.String)] = cleanDictionaryDisplay(translationText.String)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if keepGoing {
		flush()
	}
	return nil
}

func sqliteRandomUnlearnedWordID(db *sql.DB, language string, level string, learnedIDs []string, exactLevel bool) (string, bool, error) {
	where := []string{"w.language = ?"}
	args := []any{language}

	levelSQL, levelArgs := sqliteLearningLevelSQL(level, exactLevel)
	where = append(where, levelSQL)
	args = append(args, levelArgs...)

	whereSQL := strings.Join(where, " AND ")
	cacheKey := strings.Join([]string{
		"next",
		language,
		normalizeCEFRLevel(level),
		fmt.Sprintf("exact=%t", exactLevel),
	}, "|")
	exclude := map[string]bool{}
	for _, id := range learnedIDs {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		if !strings.Contains(id, ":") {
			id = makeVocabID(language, id)
		}
		exclude[id] = true
		exclude[legacyVocabID(id)] = true
	}
	return sqliteRandomVocabularyIDFromCachedWhere(db, cacheKey, whereSQL, args, exclude)
}

func sqliteFirstMatchingVocabularyID(db *sql.DB, whereSQL string, args []any, startPosition int) (string, bool, error) {
	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, startPosition)
	query := `SELECT w.id FROM vocabulary_words w WHERE ` + whereSQL + ` AND w.position >= ? ORDER BY w.position LIMIT 1`
	var id string
	err := db.QueryRow(query, queryArgs...).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func sqliteLearnedVocabularyIDs(user userState, language string) []string {
	seen := map[string]bool{}
	ids := make([]string, 0, len(user.LearnedWords))
	for _, word := range user.LearnedWords {
		if normalizeLearningLanguage(word.Language) != language {
			continue
		}
		id := strings.TrimSpace(word.ID)
		if id == "" {
			continue
		}
		if !strings.Contains(id, ":") {
			id = makeVocabID(language, id)
		}
		if seen[id] {
			continue
		}
		seen[id] = true
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func sqlitePromptAvailableSQL(learningLanguage string, interfaceLanguage string) (string, []any) {
	learningLanguage = normalizeLearningLanguage(learningLanguage)
	interfaceLanguage = normalizeInterfaceLanguage(interfaceLanguage)
	if interfaceLanguage != learningLanguage {
		if interfaceLanguage == "ru" {
			return "w.russian <> ''", nil
		}
		return `EXISTS (
			SELECT 1 FROM vocabulary_translations prompt
			WHERE prompt.word_id = w.id AND prompt.language = ? AND prompt.text <> ''
		)`, []any{interfaceLanguage}
	}
	return `EXISTS (
		SELECT 1 FROM vocabulary_translations prompt
		WHERE prompt.word_id = w.id AND prompt.language <> ? AND prompt.text <> ''
	)`, []any{learningLanguage}
}

func sqliteLearningLevelSQL(level string, exact bool) (string, []any) {
	allowed := []string{normalizeCEFRLevel(level)}
	if !exact {
		allowed = cefrOptionBand(level)
	}
	if len(allowed) == 0 {
		return "1 = 0", nil
	}
	args := make([]any, 0, len(allowed))
	for _, item := range allowed {
		args = append(args, item)
	}
	return "upper(w.level) IN (" + sqlitePlaceholders(len(allowed)) + ")", args
}

func sqliteSampleVocabularyOptionCandidates(db *sql.DB, language string, count int, level string, exclude map[string]bool) ([]vocabWord, error) {
	if count <= 0 {
		return nil, nil
	}
	where, args := sqliteVocabularyOptionWhere(language, level, nil)
	cacheKey := "options|" + language + "|" + sqliteVocabularyOptionBandKey(level)
	allIDs, err := sqliteVocabularyIDsForCachedWhere(db, cacheKey, strings.Join(where, " AND "), args)
	if err != nil {
		return nil, err
	}
	ids := sqliteSampleVocabularyIDs(allIDs, count, exclude)

	wordsByID, err := sqliteVocabularyWordsByIDs(db, ids)
	if err != nil {
		return nil, err
	}
	words := make([]vocabWord, 0, len(ids))
	for _, id := range ids {
		if word, ok := wordsByID[id]; ok {
			words = append(words, word)
		}
	}
	return words, nil
}

func sqliteVocabularyIDsForCachedWhere(db *sql.DB, cacheKey string, whereSQL string, args []any) ([]string, error) {
	vocabularySQLiteRuntime.cacheMu.RLock()
	if vocabularySQLiteRuntime.idCache != nil {
		if ids, ok := vocabularySQLiteRuntime.idCache[cacheKey]; ok {
			vocabularySQLiteRuntime.cacheMu.RUnlock()
			return ids, nil
		}
	}
	vocabularySQLiteRuntime.cacheMu.RUnlock()

	query := `SELECT w.id FROM vocabulary_words w WHERE ` + whereSQL + ` ORDER BY w.position`
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	vocabularySQLiteRuntime.cacheMu.Lock()
	defer vocabularySQLiteRuntime.cacheMu.Unlock()
	if vocabularySQLiteRuntime.idCache == nil {
		vocabularySQLiteRuntime.idCache = map[string][]string{}
	}
	if cached, ok := vocabularySQLiteRuntime.idCache[cacheKey]; ok {
		return cached, nil
	}
	vocabularySQLiteRuntime.idCache[cacheKey] = ids
	return ids, nil
}

func sqliteRandomVocabularyIDFromCachedWhere(db *sql.DB, cacheKey string, whereSQL string, args []any, exclude map[string]bool) (string, bool, error) {
	ids, err := sqliteVocabularyIDsForCachedWhere(db, cacheKey, whereSQL, args)
	if err != nil {
		return "", false, err
	}
	id, ok := sqliteRandomVocabularyID(ids, exclude)
	return id, ok, nil
}

func sqliteRandomVocabularyID(ids []string, exclude map[string]bool) (string, bool) {
	eligible := 0
	for _, id := range ids {
		if sqliteVocabularyIDExcluded(id, exclude) {
			continue
		}
		eligible++
	}
	if eligible == 0 {
		return "", false
	}
	target := vocabularyRandomIntn(eligible)
	for _, id := range ids {
		if sqliteVocabularyIDExcluded(id, exclude) {
			continue
		}
		if target == 0 {
			return id, true
		}
		target--
	}
	return "", false
}

func sqliteSampleVocabularyIDs(ids []string, count int, exclude map[string]bool) []string {
	if count <= 0 || len(ids) == 0 {
		return nil
	}
	result := make([]string, 0, count)
	seen := map[string]bool{}
	for id, value := range exclude {
		if value {
			seen[id] = true
		}
	}

	attempts := count*32 + 32
	for attempt := 0; attempt < attempts && len(result) < count; attempt++ {
		id := ids[vocabularyRandomIntn(len(ids))]
		if sqliteVocabularyIDExcluded(id, seen) {
			continue
		}
		seen[id] = true
		seen[legacyVocabID(id)] = true
		result = append(result, id)
	}
	if len(result) >= count {
		return result
	}

	needed := count - len(result)
	eligibleSeen := 0
	reservoir := make([]string, 0, needed)
	for _, id := range ids {
		if sqliteVocabularyIDExcluded(id, seen) {
			continue
		}
		eligibleSeen++
		if len(reservoir) < needed {
			reservoir = append(reservoir, id)
			continue
		}
		if index := vocabularyRandomIntn(eligibleSeen); index < needed {
			reservoir[index] = id
		}
	}
	result = append(result, reservoir...)
	return result
}

func sqliteVocabularyIDExcluded(id string, exclude map[string]bool) bool {
	if len(exclude) == 0 {
		return false
	}
	return exclude[id] || exclude[legacyVocabID(id)]
}

func sqliteVocabularyOptionBandKey(level string) string {
	if strings.TrimSpace(level) == "" {
		return "*"
	}
	band := cefrOptionBand(level)
	if len(band) == 0 {
		return normalizeCEFRLevel(level)
	}
	return strings.Join(band, ",")
}

func sqliteVocabularyOptionIDAtOrAfter(db *sql.DB, language string, level string, startPosition int, exclude map[string]bool) (string, bool, error) {
	id, ok, err := sqliteFirstVocabularyOptionID(db, language, level, startPosition, exclude)
	if err != nil || ok || startPosition == 0 {
		return id, ok, err
	}
	return sqliteFirstVocabularyOptionID(db, language, level, 0, exclude)
}

func sqliteFirstVocabularyOptionID(db *sql.DB, language string, level string, startPosition int, exclude map[string]bool) (string, bool, error) {
	where, args := sqliteVocabularyOptionWhere(language, level, exclude)
	where = append(where, "w.position >= ?")
	args = append(args, startPosition)
	query := `SELECT w.id FROM vocabulary_words w WHERE ` + strings.Join(where, " AND ") + ` ORDER BY w.position LIMIT 1`
	var id string
	err := db.QueryRow(query, args...).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return id, true, nil
}

func sqliteFirstVocabularyOptionIDs(db *sql.DB, language string, level string, count int, exclude map[string]bool) ([]string, error) {
	if count <= 0 {
		return nil, nil
	}
	where, args := sqliteVocabularyOptionWhere(language, level, exclude)
	args = append(args, count)
	query := `SELECT w.id FROM vocabulary_words w WHERE ` + strings.Join(where, " AND ") + ` ORDER BY w.position LIMIT ?`
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func sqliteVocabularyOptionWhere(language string, level string, exclude map[string]bool) ([]string, []any) {
	where := []string{"w.language = ?"}
	args := []any{language}
	if strings.TrimSpace(level) != "" {
		band := cefrOptionBand(level)
		if len(band) > 0 {
			where = append(where, "upper(w.level) IN ("+sqlitePlaceholders(len(band))+")")
			for _, item := range band {
				args = append(args, item)
			}
		}
	}
	excludeIDs := stringSetValues(exclude)
	notInSQL, notInArgs := sqliteNotInSQL("w.id", excludeIDs)
	if notInSQL != "" {
		where = append(where, notInSQL)
		args = append(args, notInArgs...)
	}
	return where, args
}

func sqliteNotInSQL(column string, values []string) (string, []any) {
	values = uniqueStrings(values)
	if len(values) == 0 {
		return "", nil
	}
	args := make([]any, 0, len(values))
	for _, value := range values {
		args = append(args, value)
	}
	return column + " NOT IN (" + sqlitePlaceholders(len(values)) + ")", args
}

func sqlitePlaceholders(count int) string {
	if count <= 0 {
		return ""
	}
	return strings.TrimRight(strings.Repeat("?,", count), ",")
}

func uniqueStrings(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}

func stringSetValues(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value, ok := range values {
		if ok {
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}
