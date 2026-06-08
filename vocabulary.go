package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

const minimumVocabularySize = 57000
const learnedWordMasteryThreshold = 10

type vocabWord struct {
	ID            string            `json:"id,omitempty"`
	Language      string            `json:"language,omitempty"`
	Russian       string            `json:"russian"`
	English       string            `json:"english"`
	Context       string            `json:"context,omitempty"`
	Translations  map[string]string `json:"translations,omitempty"`
	Level         string            `json:"level,omitempty"`
	Topic         string            `json:"topic,omitempty"`
	PartOfSpeech  string            `json:"part_of_speech,omitempty"`
	Source        string            `json:"source,omitempty"`
	FrequencyRank int               `json:"frequency_rank,omitempty"`
}

var vocabularyFiles = map[string]string{
	"en": "vocabulary_words.json",
	"ru": "vocabulary_words_ru.json",
	"es": "vocabulary_words_es.json",
	"de": "vocabulary_words_de.json",
	"fr": "vocabulary_words_fr.json",
	"it": "vocabulary_words_it.json",
	"zh": "vocabulary_words_zh.json",
	"ja": "vocabulary_words_ja.json",
	"ko": "vocabulary_words_ko.json",
	"tg": "vocabulary_words_tg.json",
	"uz": "vocabulary_words_uz.json",
	"tt": "vocabulary_words_tt.json",
	"hy": "vocabulary_words_hy.json",
	"kk": "vocabulary_words_kk.json",
	"ky": "vocabulary_words_ky.json",
	"ka": "vocabulary_words_ka.json",
	"uk": "vocabulary_words_uk.json",
	"pl": "vocabulary_words_pl.json",
	"ro": "vocabulary_words_ro.json",
	"pt": "vocabulary_words_pt.json",
	"ar": "vocabulary_words_ar.json",
	"bn": "vocabulary_words_bn.json",
	"cs": "vocabulary_words_cs.json",
	"el": "vocabulary_words_el.json",
	"hi": "vocabulary_words_hi.json",
	"hu": "vocabulary_words_hu.json",
	"id": "vocabulary_words_id.json",
	"nl": "vocabulary_words_nl.json",
	"sv": "vocabulary_words_sv.json",
	"ta": "vocabulary_words_ta.json",
	"te": "vocabulary_words_te.json",
	"th": "vocabulary_words_th.json",
	"tl": "vocabulary_words_tl.json",
	"tr": "vocabulary_words_tr.json",
	"vi": "vocabulary_words_vi.json",
}

var wikiLinkPattern = regexp.MustCompile(`\[\[(?:[^\]|]+\|)?([^\]]+)\]\]`)

var vocabularyRandom = struct {
	mu  sync.Mutex
	rng *rand.Rand
}{
	rng: rand.New(rand.NewSource(time.Now().UnixNano())),
}

func vocabularyDir() string {
	if dir := strings.TrimSpace(os.Getenv("VOCABULARY_DIR")); dir != "" {
		return dir
	}
	if _, err := os.Stat(filepath.Join("data", "vocabulary")); err == nil {
		return filepath.Join("data", "vocabulary")
	}
	return "."
}

func vocabularyPath(language string) (string, bool) {
	file, ok := vocabularyFiles[normalizeLearningLanguage(language)]
	if !ok {
		return "", false
	}
	return filepath.Join(vocabularyDir(), file), true
}

func loadVocabularyFromDisk(language string) ([]vocabWord, error) {
	path, ok := vocabularyPath(language)
	if !ok {
		return nil, fmt.Errorf("unknown vocabulary language: %s", language)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return loadVocabulary(language, data)
}

func loadVocabulary(language string, data []byte) ([]vocabWord, error) {
	var raw []vocabWord
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse vocabulary_words.json: %w", err)
	}

	byID := map[string]vocabWord{}
	order := make([]string, 0, len(raw))
	for _, word := range raw {
		word, ok := normalizeVocabularyWord(language, word)
		if !ok {
			continue
		}
		if _, exists := byID[word.ID]; exists {
			continue
		}
		byID[word.ID] = word
		order = append(order, word.ID)
	}

	words := make([]vocabWord, 0, len(order))
	for _, id := range order {
		words = append(words, byID[id])
	}
	return words, nil
}

func normalizeVocabularyWord(language string, word vocabWord) (vocabWord, bool) {
	word.Language = normalizeLearningLanguage(language)
	word.English = strings.TrimSpace(word.English)
	word.Russian = cleanDictionaryDisplay(word.Russian)
	word.Context = normalizeVocabularyContext(word.Context)
	if word.English == "" || word.Russian == "" {
		return vocabWord{}, false
	}
	if word.Translations == nil {
		word.Translations = map[string]string{}
	}
	if strings.TrimSpace(word.Russian) != "" {
		word.Translations["ru"] = word.Russian
	}
	for code, value := range word.Translations {
		value = cleanDictionaryDisplay(value)
		if value == "" {
			delete(word.Translations, code)
			continue
		}
		delete(word.Translations, code)
		word.Translations[normalizeInterfaceLanguage(code)] = value
	}
	word.ID = makeVocabID(word.Language, word.English)
	return word, true
}

var vocabularyLearningWordPattern = regexp.MustCompile(`^[\p{L}][\p{L}\p{N}'’]*$`)

type generatedVocabularyWord struct {
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	Russian      string `json:"russian"`
	Context      string `json:"context"`
	Level        string `json:"level"`
	Topic        string `json:"topic"`
	PartOfSpeech string `json:"part_of_speech"`
}

func parseGeneratedVocabularyWord(raw string, language string, promptLanguage string, fallbackLevel string, forbidden map[string]bool) (vocabWord, bool) {
	raw = stripJSONCodeFence(cleanModelReply(raw))
	if raw == "" {
		return vocabWord{}, false
	}
	var payload generatedVocabularyWord
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return vocabWord{}, false
	}
	return sanitizeGeneratedVocabularyWord(language, payload, promptLanguage, fallbackLevel, forbidden)
}

func sanitizeGeneratedVocabularyWord(language string, generated generatedVocabularyWord, promptLanguage string, fallbackLevel string, forbidden map[string]bool) (vocabWord, bool) {
	language = normalizeLearningLanguage(language)
	promptLanguage = normalizeInterfaceLanguage(promptLanguage)
	if promptLanguage == language {
		if language == "ru" {
			promptLanguage = "en"
		} else {
			promptLanguage = "ru"
		}
	}
	target := cleanDictionaryDisplay(generated.Word)
	if target == "" {
		return vocabWord{}, false
	}
	id := makeVocabID(language, target)
	if vocabularyWordForbidden(id, target, forbidden) {
		return vocabWord{}, false
	}

	russian := cleanDictionaryDisplay(generated.Russian)
	if language == "ru" && russian == "" {
		russian = target
	}
	if language != "ru" && russian != "" && dictionaryTextContainsTerm(russian, target) {
		return vocabWord{}, false
	}
	context := normalizeVocabularyContext(generated.Context)
	if dictionaryTextContainsTerm(context, target) {
		context = ""
	}
	candidate := vocabWord{
		ID:           id,
		Language:     language,
		English:      target,
		Russian:      russian,
		Context:      context,
		Translations: map[string]string{},
		Level:        normalizeCEFRLevel(firstNonEmpty(generated.Level, fallbackLevel)),
		Topic:        cleanDictionaryDisplay(generated.Topic),
		PartOfSpeech: cleanDictionaryDisplay(generated.PartOfSpeech),
		Source:       "ai",
	}
	translation := sanitizeVocabularyTranslation(generated.Translation, candidate)
	if translation == "" {
		return vocabWord{}, false
	}
	if promptLanguage == "ru" {
		candidate.Russian = translation
	} else if candidate.Russian == "" {
		return vocabWord{}, false
	}
	candidate.Translations[promptLanguage] = translation

	word, ok := normalizeVocabularyWord(language, candidate)
	if !ok {
		return vocabWord{}, false
	}
	word.Level = normalizeCEFRLevel(candidate.Level)
	word.Topic = cleanDictionaryDisplay(candidate.Topic)
	word.PartOfSpeech = cleanDictionaryDisplay(candidate.PartOfSpeech)
	word.Source = "ai"
	if !vocabularyWordSuitableForLearning(word) {
		return vocabWord{}, false
	}
	return word, true
}

func vocabularyWordForbidden(id string, word string, forbidden map[string]bool) bool {
	if len(forbidden) == 0 {
		return false
	}
	id = strings.TrimSpace(id)
	word = cleanDictionaryDisplay(word)
	normalizedWord := normalizeAnswer(word)
	for _, candidate := range []string{id, legacyVocabID(id), word, strings.ToLower(word), normalizedWord} {
		if strings.TrimSpace(candidate) != "" && forbidden[candidate] {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func vocabularyGenerationPromptLanguage(user userState) learningLanguage {
	learningLanguage := normalizeLearningLanguage(user.LearningLanguage)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	if interfaceLanguage != learningLanguage {
		return interfaceLanguageByCode(interfaceLanguage)
	}
	if learningLanguage == "ru" {
		return interfaceLanguageByCode("en")
	}
	return interfaceLanguageByCode("ru")
}

func vocabularyGenerationForbiddenWords(user userState, extra []string) (map[string]bool, []string) {
	language := normalizeLearningLanguage(user.LearningLanguage)
	forbidden := map[string]bool{}
	visible := []string{}
	add := func(value string) {
		value = cleanDictionaryDisplay(value)
		if value == "" {
			return
		}
		id := makeVocabID(language, value)
		forbidden[id] = true
		forbidden[legacyVocabID(id)] = true
		forbidden[strings.ToLower(value)] = true
		forbidden[normalizeAnswer(value)] = true
		for _, existing := range visible {
			if sameDictionaryText(existing, value) {
				return
			}
		}
		visible = append(visible, value)
	}
	for _, learned := range user.LearnedWords {
		if normalizeLearningLanguage(learned.Language) != language {
			continue
		}
		if strings.TrimSpace(learned.ID) != "" {
			id := learned.ID
			if !strings.Contains(id, ":") {
				id = makeVocabID(language, id)
			}
			forbidden[id] = true
			forbidden[legacyVocabID(id)] = true
		}
		add(learned.English)
	}
	for _, item := range extra {
		add(item)
	}
	return forbidden, visible
}

func vocabularyWordSuitableForLearning(word vocabWord) bool {
	text := strings.TrimSpace(word.English)
	if text == "" || strings.TrimSpace(word.Russian) == "" {
		return false
	}
	if strings.ContainsAny(text, "-‐‑‒–—―_/\\") {
		return false
	}
	if len([]rune(text)) < 3 {
		return false
	}
	return vocabularyLearningWordPattern.MatchString(text)
}

func cleanDictionaryDisplay(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = wikiLinkPattern.ReplaceAllString(value, "$1")
	value = strings.ReplaceAll(value, " ;", ";")
	value = strings.ReplaceAll(value, "; ", "; ")
	value = strings.Join(strings.Fields(value), " ")
	return strings.Trim(value, " \t\r\n;,")
}

func normalizeVocabularyContext(value string) string {
	return cleanDictionaryDisplay(value)
}

func forEachVocabularyWord(language string, visit func(vocabWord) bool) error {
	if used, err := sqliteForEachVocabularyWord(language, visit); used || err != nil {
		return err
	}
	language = normalizeLearningLanguage(language)
	path, ok := vocabularyPath(language)
	if !ok {
		return fmt.Errorf("unknown vocabulary language: %s", language)
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
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

	for decoder.More() {
		var word vocabWord
		if err := decoder.Decode(&word); err != nil {
			return fmt.Errorf("decode %s: %w", path, err)
		}
		word, ok := normalizeVocabularyWord(language, word)
		if !ok {
			continue
		}
		if !visit(word) {
			return nil
		}
	}
	return nil
}

func wordTranslationForLanguage(word vocabWord, language string) string {
	language = normalizeInterfaceLanguage(language)
	learningLanguage := normalizeLearningLanguage(word.Language)
	if language == learningLanguage {
		return ""
	}
	if word.Translations != nil {
		if translation := dictionaryDisplayValues(word.Translations[language], word.English, 1); translation != "" {
			return firstDictionaryValue(translation)
		}
	}
	if language == "ru" {
		return firstDictionaryValue(dictionaryDisplayValues(word.Russian, word.English, 1))
	}
	return ""
}

func wordPromptLanguage(word vocabWord, interfaceLanguage string) string {
	interfaceLanguage = normalizeInterfaceLanguage(interfaceLanguage)
	learningLanguage := normalizeLearningLanguage(word.Language)
	if interfaceLanguage != learningLanguage {
		if translation := wordTranslationForLanguage(word, interfaceLanguage); translation != "" {
			return interfaceLanguage
		}
	}
	if interfaceLanguage == learningLanguage {
		for _, fallbackLanguage := range []string{"en", "ru"} {
			if fallbackLanguage == learningLanguage {
				continue
			}
			if translation := wordTranslationForLanguage(word, fallbackLanguage); translation != "" {
				return fallbackLanguage
			}
		}
		if word.Translations != nil {
			codes := make([]string, 0, len(word.Translations))
			for code := range word.Translations {
				code = normalizeInterfaceLanguage(code)
				if code != learningLanguage {
					codes = append(codes, code)
				}
			}
			sort.Strings(codes)
			for _, code := range codes {
				if translation := wordTranslationForLanguage(word, code); translation != "" {
					return code
				}
			}
		}
		return interfaceLanguage
	}
	if interfaceLanguage == "ru" {
		if translation := wordTranslationForLanguage(word, "ru"); translation != "" {
			return "ru"
		}
	}
	return interfaceLanguage
}

func wordTranslation(word vocabWord, interfaceLanguage string) string {
	learningLanguage := normalizeLearningLanguage(word.Language)
	interfaceLanguage = normalizeInterfaceLanguage(interfaceLanguage)
	promptLanguage := wordPromptLanguage(word, interfaceLanguage)
	if promptLanguage != learningLanguage {
		if translation := wordTranslationForLanguage(word, promptLanguage); translation != "" {
			return translation
		}
		return ""
	}
	if interfaceLanguage != learningLanguage {
		return ""
	}
	if interfaceLanguage == learningLanguage && word.Translations != nil {
		if interfaceLanguage != "en" {
			if translation := cleanDictionaryDisplay(word.Translations["en"]); translation != "" {
				return firstDictionaryValue(translation)
			}
		}
		if interfaceLanguage != "ru" {
			if translation := cleanDictionaryDisplay(word.Translations["ru"]); translation != "" {
				return firstDictionaryValue(translation)
			}
		}
	}
	if word.Translations != nil {
		if translation := cleanDictionaryDisplay(word.Translations[interfaceLanguage]); translation != "" && interfaceLanguage != learningLanguage {
			return firstDictionaryValue(translation)
		}
	}
	if interfaceLanguage == "ru" && cleanDictionaryDisplay(word.Russian) != "" && learningLanguage != "ru" {
		return firstDictionaryValue(word.Russian)
	}
	return ""
}

func wordPromptAvailableForUser(word vocabWord, user userState) bool {
	return wordTranslation(word, user.InterfaceLanguage) != ""
}

func wordContext(word vocabWord, interfaceLanguage string) string {
	interfaceLanguage = normalizeInterfaceLanguage(interfaceLanguage)
	learningLanguage := normalizeLearningLanguage(word.Language)
	if interfaceLanguage == learningLanguage {
		return ""
	}
	primary := wordTranslation(word, interfaceLanguage)
	target := cleanDictionaryDisplay(word.English)
	values := dictionaryContextValues("", primary, target)
	if word.Translations != nil {
		values = dictionaryContextValues(word.Translations[interfaceLanguage], primary, target)
	}
	if len(values) == 0 && interfaceLanguage == "ru" && learningLanguage != "ru" {
		values = dictionaryContextValues(word.Context, primary, target)
		if len(values) == 0 {
			values = dictionaryContextValues(word.Russian, primary, target)
		}
	}
	if len(values) == 0 {
		return ""
	}
	if len(values) > 2 {
		values = values[:2]
	}
	return strings.Join(values, " / ")
}

func learnedWordContext(entry learnedWordEntry, interfaceLanguage string) string {
	if word, ok := findVocabWord(entry.ID); ok {
		return wordContext(word, interfaceLanguage)
	}
	return wordContext(vocabWord{
		ID:       entry.ID,
		Language: entry.Language,
		Russian:  entry.Russian,
		English:  entry.English,
		Context:  entry.Context,
	}, interfaceLanguage)
}

func hasDictionaryContext(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	if strings.ContainsAny(value, ";/()") {
		return true
	}
	return len(strings.Fields(value)) > 3
}

func sameDictionaryText(left string, right string) bool {
	return normalizeAnswer(cleanDictionaryDisplay(left)) == normalizeAnswer(cleanDictionaryDisplay(right))
}

func dictionaryTextContainsTerm(text string, term string) bool {
	text = cleanDictionaryDisplay(text)
	term = cleanDictionaryDisplay(term)
	if text == "" || term == "" {
		return false
	}
	if sameDictionaryText(text, term) {
		return true
	}
	normalizedTerm := normalizeAnswer(term)
	for _, piece := range strings.FieldsFunc(text, func(r rune) bool {
		return r == ';' || r == ',' || r == '/' || r == '(' || r == ')' || r == '[' || r == ']'
	}) {
		if normalizeAnswer(piece) == normalizedTerm {
			return true
		}
	}
	return false
}

func firstDictionaryValue(value string) string {
	value = cleanDictionaryDisplay(value)
	if value == "" {
		return ""
	}
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ';' || r == '/' || r == ',' || r == '(' || r == ')'
	})
	for _, part := range parts {
		part = cleanDictionaryDisplay(part)
		if part != "" {
			return part
		}
	}
	return value
}

func dictionaryDisplayValues(value string, target string, max int) string {
	value = cleanDictionaryDisplay(value)
	target = cleanDictionaryDisplay(target)
	if value == "" || max <= 0 {
		return ""
	}
	values := make([]string, 0, max)
	for _, part := range strings.FieldsFunc(value, func(r rune) bool {
		return r == ';' || r == '/' || r == ',' || r == '(' || r == ')' || r == '[' || r == ']'
	}) {
		part = cleanDictionaryDisplay(part)
		if part == "" || dictionaryTextContainsTerm(part, target) {
			continue
		}
		duplicate := false
		for _, existing := range values {
			if sameDictionaryText(existing, part) {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}
		values = append(values, part)
		if len(values) >= max {
			break
		}
	}
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, "; ")
}

func dictionaryContextValues(value string, primary string, target string) []string {
	value = cleanDictionaryDisplay(value)
	primary = cleanDictionaryDisplay(primary)
	target = cleanDictionaryDisplay(target)
	if value == "" {
		return nil
	}
	values := []string{}
	for _, part := range strings.FieldsFunc(value, func(r rune) bool {
		return r == ';' || r == '/' || r == ',' || r == '(' || r == ')' || r == '[' || r == ']'
	}) {
		part = cleanDictionaryDisplay(part)
		if part == "" || sameDictionaryText(part, primary) || dictionaryTextContainsTerm(part, target) {
			continue
		}
		duplicate := false
		for _, existing := range values {
			if sameDictionaryText(existing, part) {
				duplicate = true
				break
			}
		}
		if !duplicate {
			values = append(values, part)
		}
	}
	return values
}

func learnedWordTranslation(entry learnedWordEntry, interfaceLanguage string) string {
	return learnedWordTranslationWithLookup(entry, interfaceLanguage, nil)
}

func learnedWordTranslationWithLookup(entry learnedWordEntry, interfaceLanguage string, lookup map[string]vocabWord) string {
	if word, ok := lookup[entry.ID]; ok {
		return wordTranslation(word, interfaceLanguage)
	}
	if word, ok := findVocabWord(entry.ID); ok {
		return wordTranslation(word, interfaceLanguage)
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return cleanDictionaryDisplay(entry.Russian)
	}
	return ""
}

func learnedWordContextWithLookup(entry learnedWordEntry, interfaceLanguage string, lookup map[string]vocabWord) string {
	if word, ok := lookup[entry.ID]; ok {
		return wordContext(word, interfaceLanguage)
	}
	return learnedWordContext(entry, interfaceLanguage)
}

func learnedEntryVocabWord(entry learnedWordEntry) vocabWord {
	if word, ok := findVocabWord(entry.ID); ok {
		return word
	}
	return vocabWord{
		ID:       entry.ID,
		Language: entry.Language,
		Russian:  entry.Russian,
		English:  entry.English,
		Context:  entry.Context,
		Translations: map[string]string{
			"ru": cleanDictionaryDisplay(entry.Russian),
		},
	}
}

func makeVocabID(language string, english string) string {
	id := strings.ToLower(strings.TrimSpace(english))
	replacer := strings.NewReplacer(" ", "-", "'", "", "/", "-", ".", "", ",", "")
	return normalizeLearningLanguage(language) + ":" + replacer.Replace(id)
}

func languageFromVocabID(id string) string {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		return normalizeLearningLanguage(parts[0])
	}
	return "en"
}

func vocabularySize() int {
	count := 0
	if err := forEachVocabularyWord("en", func(vocabWord) bool {
		count++
		return true
	}); err != nil {
		panic(err)
	}
	return count
}

func vocabularyLanguageHasAny(language string) bool {
	found := false
	if err := forEachVocabularyWord(language, func(vocabWord) bool {
		found = true
		return false
	}); err != nil {
		panic(err)
	}
	return found
}

func findVocabWord(id string) (vocabWord, bool) {
	id = strings.TrimSpace(id)
	words := findVocabWords([]string{id})
	word, ok := words[id]
	return word, ok
}

func findVocabWords(ids []string) map[string]vocabWord {
	if result, used, err := sqliteFindVocabWords(ids); used {
		if err != nil {
			panic(err)
		}
		return result
	}
	result := map[string]vocabWord{}
	targetsByLanguage := map[string]map[string][]string{}
	remainingByLanguage := map[string]int{}
	requested := map[string]bool{}
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || requested[id] {
			continue
		}
		requested[id] = true
		language := languageFromVocabID(id)
		if targetsByLanguage[language] == nil {
			targetsByLanguage[language] = map[string][]string{}
		}
		targetsByLanguage[language][id] = append(targetsByLanguage[language][id], id)
		if legacy := legacyVocabID(id); legacy != id {
			targetsByLanguage[language][legacy] = append(targetsByLanguage[language][legacy], id)
		}
		remainingByLanguage[language]++
	}
	for language, targets := range targetsByLanguage {
		remaining := remainingByLanguage[language]
		if err := forEachVocabularyWord(language, func(word vocabWord) bool {
			for _, key := range []string{word.ID, legacyVocabID(word.ID)} {
				for _, originalID := range targets[key] {
					if _, ok := result[originalID]; ok {
						continue
					}
					result[originalID] = word
					remaining--
				}
			}
			return remaining > 0
		}); err != nil {
			panic(err)
		}
	}
	return result
}

func vocabularyLookupForLearnedWords(words []learnedWordEntry) map[string]vocabWord {
	ids := make([]string, 0, len(words))
	for _, word := range words {
		ids = append(ids, word.ID)
	}
	return findVocabWords(ids)
}

func legacyVocabID(id string) string {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return id
}

var cefrLevels = []string{"A1", "A2", "B1", "B2", "C1", "C2"}

func normalizeCEFRLevel(level string) string {
	level = strings.ToUpper(strings.TrimSpace(level))
	for _, item := range cefrLevels {
		if level == item {
			return item
		}
	}
	return "A2"
}

func cefrOptionBand(level string) []string {
	switch normalizeCEFRLevel(level) {
	case "A1", "A2":
		return []string{"A1", "A2"}
	case "B1", "B2":
		return []string{"B1", "B2"}
	case "C1", "C2":
		return []string{"C1", "C2"}
	default:
		return nil
	}
}

func vocabWordInOptionBand(word vocabWord, correct vocabWord) bool {
	if strings.TrimSpace(word.Level) == "" {
		return false
	}
	for _, level := range cefrOptionBand(correct.Level) {
		if normalizeCEFRLevel(word.Level) == level {
			return true
		}
	}
	return false
}

func cefrRank(level string) int {
	level = normalizeCEFRLevel(level)
	for index, item := range cefrLevels {
		if level == item {
			return index
		}
	}
	return 1
}

func nextCEFRLevel(level string) (string, bool) {
	rank := cefrRank(level)
	if rank+1 >= len(cefrLevels) {
		return normalizeCEFRLevel(level), false
	}
	return cefrLevels[rank+1], true
}

func wordMatchesLearningLevel(word vocabWord, level string) bool {
	if strings.TrimSpace(word.Level) == "" {
		return false
	}
	for _, candidate := range cefrOptionBand(level) {
		if normalizeCEFRLevel(word.Level) == candidate {
			return true
		}
	}
	return false
}

func preferredLearningWord(word vocabWord) bool {
	topic := strings.ToLower(strings.TrimSpace(word.Topic))
	source := strings.ToLower(strings.TrimSpace(word.Source))
	return topic == "multilingual-core" || strings.Contains(source, "multilingual core")
}

func appendVocabularyPoolCandidate(pool []vocabWord, seen int, limit int, word vocabWord) []vocabWord {
	if len(pool) < limit {
		return append(pool, word)
	}
	if index := vocabularyRandomIntn(seen); index < limit {
		pool[index] = word
	}
	return pool
}

func nextUnlearnedWord(user userState) (vocabWord, bool) {
	if word, ok, used, err := sqliteNextUnlearnedWord(user); used {
		if err != nil {
			panic(err)
		}
		return word, ok
	}
	language := normalizeLearningLanguage(user.LearningLanguage)
	userRank := cefrRank(user.Level)
	learned := map[string]bool{}
	for _, word := range user.LearnedWords {
		if normalizeLearningLanguage(word.Language) == language {
			learned[word.ID] = true
			learned[legacyVocabID(word.ID)] = true
		}
	}
	const poolLimit = 768
	exactPool := make([]vocabWord, 0, 96)
	preferredExactPool := make([]vocabWord, 0, 96)
	fallbackPool := make([]vocabWord, 0, 96)
	preferredFallbackPool := make([]vocabWord, 0, 96)
	exactSeen := 0
	preferredExactSeen := 0
	fallbackSeen := 0
	preferredFallbackSeen := 0
	if err := forEachVocabularyWord(language, func(word vocabWord) bool {
		if learned[word.ID] || learned[legacyVocabID(word.ID)] {
			return true
		}
		if !vocabularyWordSuitableForLearning(word) {
			return true
		}
		if cefrRank(word.Level) != userRank {
			if wordMatchesLearningLevel(word, user.Level) {
				fallbackSeen++
				fallbackPool = appendVocabularyPoolCandidate(fallbackPool, fallbackSeen, poolLimit, word)
				if preferredLearningWord(word) {
					preferredFallbackSeen++
					preferredFallbackPool = appendVocabularyPoolCandidate(preferredFallbackPool, preferredFallbackSeen, poolLimit, word)
				}
			}
			return true
		}
		exactSeen++
		exactPool = appendVocabularyPoolCandidate(exactPool, exactSeen, poolLimit, word)
		if preferredLearningWord(word) {
			preferredExactSeen++
			preferredExactPool = appendVocabularyPoolCandidate(preferredExactPool, preferredExactSeen, poolLimit, word)
		}
		return true
	}); err != nil {
		panic(err)
	}
	if len(preferredExactPool) > 0 {
		return preferredExactPool[vocabularyRandomIntn(len(preferredExactPool))], true
	}
	if len(exactPool) > 0 {
		return exactPool[vocabularyRandomIntn(len(exactPool))], true
	}
	if len(preferredFallbackPool) > 0 {
		return preferredFallbackPool[vocabularyRandomIntn(len(preferredFallbackPool))], true
	}
	if len(fallbackPool) > 0 {
		return fallbackPool[vocabularyRandomIntn(len(fallbackPool))], true
	}
	return vocabWord{}, false
}

func learnedWordIDs(user userState) []string {
	return learnedWordIDsByMode(user, func(learnedWordEntry) bool { return true })
}

func reviewWordIDs(user userState) []string {
	ids := learnedWordIDsByMode(user, func(word learnedWordEntry) bool {
		return !learnedWordMastered(word)
	})
	randomStrings(ids)
	return ids
}

func spellingWordIDs(user userState) []string {
	ids := learnedWordIDsByMode(user, func(word learnedWordEntry) bool {
		return !learnedWordMastered(word)
	})
	randomStrings(ids)
	return ids
}

func learnedWordIDsByMode(user userState, accept func(learnedWordEntry) bool) []string {
	language := normalizeLearningLanguage(user.LearningLanguage)
	ids := make([]string, 0, len(user.LearnedWords))
	for _, word := range user.LearnedWords {
		if normalizeLearningLanguage(word.Language) == language && accept(word) {
			ids = append(ids, word.ID)
		}
	}
	return ids
}

func learnedWordTrainingCount(word learnedWordEntry) int {
	return word.ReviewCorrectCount + word.SpellingCorrectCount
}

func learnedWordMastered(word learnedWordEntry) bool {
	return learnedWordTrainingCount(word) >= learnedWordMasteryThreshold
}

func masteredWordCount(user userState) int {
	count := 0
	for _, word := range user.LearnedWords {
		if learnedWordMastered(word) {
			count++
		}
	}
	return count
}

func masteredWordCountForLanguage(user userState, languageCode string) int {
	count := 0
	language := normalizeLearningLanguage(languageCode)
	for _, word := range user.LearnedWords {
		if normalizeLearningLanguage(word.Language) == language && learnedWordMastered(word) {
			count++
		}
	}
	return count
}

func vocabularyForUser(user userState) []vocabWord {
	return vocabularyForLanguage(user.LearningLanguage)
}

func vocabularyForLanguage(language string) []vocabWord {
	language = normalizeLearningLanguage(language)
	if words, used, err := sqliteVocabularyForLanguage(language); used {
		if err != nil {
			panic(err)
		}
		if language == "en" && len(words) < minimumVocabularySize {
			panic(fmt.Sprintf("basic vocabulary must contain at least %d unique words", minimumVocabularySize))
		}
		if len(words) == 0 {
			panic(fmt.Sprintf("%s vocabulary must not be empty", language))
		}
		return words
	}
	words, err := loadVocabularyFromDisk(language)
	if err != nil {
		panic(err)
	}
	if language == "en" && len(words) < minimumVocabularySize {
		panic(fmt.Sprintf("basic vocabulary must contain at least %d unique words", minimumVocabularySize))
	}
	if len(words) == 0 {
		panic(fmt.Sprintf("%s vocabulary must not be empty", language))
	}
	return words
}

func wordOptions(correct vocabWord) []vocabWord {
	if options, used, err := sqliteWordOptions(correct); used {
		if err != nil {
			panic(err)
		}
		if len(options) > 1 {
			return options
		}
	}
	options := []vocabWord{correct}
	candidates := sampleVocabularyWords(correct.Language, 3, func(word vocabWord) bool {
		if word.ID == correct.ID || legacyVocabID(word.ID) == legacyVocabID(correct.ID) {
			return false
		}
		if !vocabWordInOptionBand(word, correct) {
			return false
		}
		return true
	})
	options = append(options, firstWords(candidates, 3)...)
	randomWords(options)
	return options
}

func learnedWordOptions(correct vocabWord, learnedIDs []string) []vocabWord {
	if options, used, err := sqliteLearnedWordOptions(correct, learnedIDs); used {
		if err != nil {
			panic(err)
		}
		if len(options) > 1 {
			return options
		}
	}
	options := []vocabWord{correct}
	learnedLookup := map[string]bool{}
	for _, id := range learnedIDs {
		learnedLookup[id] = true
		learnedLookup[legacyVocabID(id)] = true
	}
	seenCandidates := map[string]bool{correct.ID: true, legacyVocabID(correct.ID): true}
	candidates := sampleVocabularyWords(correct.Language, 3, func(word vocabWord) bool {
		legacyID := legacyVocabID(word.ID)
		if seenCandidates[word.ID] || seenCandidates[legacyID] || learnedLookup[word.ID] || learnedLookup[legacyID] {
			return false
		}
		return vocabWordInOptionBand(word, correct)
	})
	for _, word := range candidates {
		seenCandidates[word.ID] = true
		seenCandidates[legacyVocabID(word.ID)] = true
	}
	options = append(options, firstWords(candidates, 3)...)

	randomWords(options)
	return options
}

func sampleVocabularyWords(language string, count int, accept func(vocabWord) bool) []vocabWord {
	if count <= 0 {
		return nil
	}
	samples := make([]vocabWord, 0, count)
	seen := 0
	if err := forEachVocabularyWord(language, func(word vocabWord) bool {
		if !accept(word) {
			return true
		}
		seen++
		if len(samples) < count {
			samples = append(samples, word)
			return true
		}
		index := vocabularyRandomIntn(seen)
		if index < count {
			samples[index] = word
		}
		return true
	}); err != nil {
		panic(err)
	}
	return samples
}

func vocabularyRandomIntn(max int) int {
	vocabularyRandom.mu.Lock()
	defer vocabularyRandom.mu.Unlock()
	return vocabularyRandom.rng.Intn(max)
}

func randomWords(words []vocabWord) {
	vocabularyRandom.mu.Lock()
	defer vocabularyRandom.mu.Unlock()
	vocabularyRandom.rng.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
}

func randomStrings(values []string) {
	vocabularyRandom.mu.Lock()
	defer vocabularyRandom.mu.Unlock()
	vocabularyRandom.rng.Shuffle(len(values), func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})
}

func firstWords(words []vocabWord, count int) []vocabWord {
	if count > len(words) {
		count = len(words)
	}
	return words[:count]
}

func containsWordID(words []vocabWord, id string) bool {
	for _, word := range words {
		if word.ID == id {
			return true
		}
	}
	return false
}

func containsLegacyWordID(words []vocabWord, id string) bool {
	for _, word := range words {
		if legacyVocabID(word.ID) == id {
			return true
		}
	}
	return false
}

func knowledgeLevel(xp int) (int, string, int, int) {
	if xp < 0 {
		xp = 0
	}
	names := []string{
		"Новичок алфавита", "Собиратель слов", "Охотник за фразами", "Разговорный искатель", "Ученик полиглота",
		"Хранитель словаря", "Мастер small talk", "Грамматический разведчик", "Фразовый алхимик", "Уверенный спикер",
		"Путешественник речи", "Строитель диалогов", "Лингвистический навигатор", "Посол английского", "Полиглот-стратег",
		"Ловец акцентов", "Мастер живой речи", "Архитектор языка", "Легенда перевода", "Верховный полиглот",
	}
	thresholds := []int{
		0, 120, 280, 500, 800,
		1200, 1750, 2500, 3500, 4800,
		6500, 8700, 11500, 15000, 19500,
		25500, 33500, 44500, 60000, 80000,
	}

	level := 1
	for i := len(thresholds) - 1; i >= 0; i-- {
		if xp >= thresholds[i] {
			level = i + 1
			break
		}
	}
	name := names[level-1]

	current := 100
	need := 100
	if level == 20 {
		return level, name, current, need
	}
	current = xp - thresholds[level-1]
	need = thresholds[level] - thresholds[level-1]
	return level, name, current, need
}

func vocabularySourceSummary() string {
	sources := map[string]bool{}
	if err := forEachVocabularyWord("en", func(word vocabWord) bool {
		if word.Source != "" {
			sources[word.Source] = true
		}
		return true
	}); err != nil {
		panic(err)
	}
	list := make([]string, 0, len(sources))
	for source := range sources {
		list = append(list, source)
	}
	sort.Strings(list)
	return strings.Join(list, "; ")
}
