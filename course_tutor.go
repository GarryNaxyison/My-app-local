package main

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

type tutorLessonWord struct {
	ID           string `json:"id"`
	Word         string `json:"word"`
	Translation  string `json:"translation"`
	Context      string `json:"context,omitempty"`
	Example      string `json:"example"`
	ExampleRu    string `json:"example_translation,omitempty"`
	PartOfSpeech string `json:"part_of_speech,omitempty"`
	Topic        string `json:"topic"`
	Level        string `json:"level"`
	Frequency    int    `json:"frequency,omitempty"`
	AudioText    string `json:"audio_text"`
}

type tutorLessonChoice struct {
	Prompt          string              `json:"prompt"`
	Options         []tutorChoiceOption `json:"options"`
	CorrectAnswerID string              `json:"correct_answer_id"`
	Feedback        string              `json:"feedback,omitempty"`
}

type tutorChoiceOption struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Label   string `json:"label,omitempty"`
	Why     string `json:"why,omitempty"`
	Skill   string `json:"skill,omitempty"`
	Quality string `json:"quality,omitempty"`
	Avoid   bool   `json:"avoid,omitempty"`
}

type tutorAnswerVariant struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Text    string `json:"text"`
	Why     string `json:"why"`
	Level   string `json:"level,omitempty"`
	UseCase string `json:"use_case,omitempty"`
	Avoid   bool   `json:"avoid,omitempty"`
}

type tutorLessonSRS struct {
	Again string `json:"again"`
	Hard  string `json:"hard"`
	Good  string `json:"good"`
	Easy  string `json:"easy"`
}

type tutorLessonStep struct {
	Code    string `json:"code"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type tutorLesson struct {
	ID                string               `json:"id"`
	Title             string               `json:"title"`
	Level             string               `json:"level"`
	Topic             string               `json:"topic"`
	VariantCode       string               `json:"variant_code,omitempty"`
	VariantTitle      string               `json:"variant_title,omitempty"`
	FocusSkill        string               `json:"focus_skill,omitempty"`
	SuccessCriteria   []string             `json:"success_criteria,omitempty"`
	LessonNumber      int                  `json:"lesson_number"`
	CourseSize        int                  `json:"course_size"`
	LearningLanguage  string               `json:"learning_language"`
	InterfaceLanguage string               `json:"interface_language"`
	DurationMinutes   int                  `json:"duration_minutes"`
	Goal              string               `json:"goal"`
	CanDo             string               `json:"can_do"`
	Scenario          string               `json:"scenario"`
	Steps             []tutorLessonStep    `json:"steps"`
	Words             []tutorLessonWord    `json:"words"`
	GrammarTitle      string               `json:"grammar_title"`
	Grammar           string               `json:"grammar"`
	MiniExplanation   string               `json:"mini_explanation"`
	Choice            tutorLessonChoice    `json:"choice"`
	Checks            []tutorLessonChoice  `json:"checks,omitempty"`
	WritingTask       string               `json:"writing_task"`
	WritingExpected   []string             `json:"writing_expected,omitempty"`
	AnswerVariants    []tutorAnswerVariant `json:"answer_variants,omitempty"`
	ListeningTask     string               `json:"listening_task"`
	ListeningText     string               `json:"listening_text"`
	ListeningQuestion string               `json:"listening_question,omitempty"`
	ListeningExpected []string             `json:"listening_expected,omitempty"`
	PronunciationText string               `json:"pronunciation_text,omitempty"`
	Dialogue          []string             `json:"dialogue"`
	DialoguePrompt    string               `json:"dialogue_prompt,omitempty"`
	DialogueGoal      string               `json:"dialogue_goal,omitempty"`
	DialogueVariants  []tutorAnswerVariant `json:"dialogue_variants,omitempty"`
	Review            []tutorLessonWord    `json:"review"`
	ReviewSummary     []string             `json:"review_summary,omitempty"`
	SRS               tutorLessonSRS       `json:"srs"`
	NextActions       []string             `json:"next_actions"`
	Source            string               `json:"source"`
}

type tutorTopic struct {
	Code       string
	TitleRU    string
	TitleEN    string
	ScenarioRU string
	ScenarioEN string
	Words      []string
}

type tutorCourseFunction struct {
	Code    string
	TitleRU string
	TitleEN string
	CanDoRU string
	CanDoEN string
}

type tutorCourseVariant struct {
	Code         string
	TitleRU      string
	TitleEN      string
	FocusRU      string
	FocusEN      string
	ConstraintRU string
	ConstraintEN string
}

const tutorCourseSize = 300

var tutorScenarioCourseTopics = []tutorTopic{
	{Code: "hotel", TitleRU: "Отель и поездка", TitleEN: "Hotel and travel", ScenarioRU: "Вы на ресепшене отеля: нужно подтвердить бронь, уточнить номер и показать документ.", ScenarioEN: "You are at a hotel reception: confirm a reservation, ask about the room, and show a document.", Words: []string{"hotel", "room", "reservation", "passport", "ticket", "station", "airport", "street", "address"}},
	{Code: "food", TitleRU: "Еда и заказ", TitleEN: "Food and ordering", ScenarioRU: "Вы в кафе: нужно попросить меню, заказать напиток и уточнить завтрак или ужин.", ScenarioEN: "You are in a cafe: ask for the menu, order a drink, and clarify breakfast or dinner.", Words: []string{"food", "water", "coffee", "tea", "bread", "menu", "restaurant", "breakfast", "dinner"}},
	{Code: "city", TitleRU: "Город и маршруты", TitleEN: "City and directions", ScenarioRU: "Вы в городе: нужно спросить дорогу, назвать адрес и понять, где находится нужное место.", ScenarioEN: "You are in the city: ask for directions, give an address, and understand where a place is.", Words: []string{"city", "street", "house", "shop", "market", "bank", "school", "office", "park"}},
	{Code: "work", TitleRU: "Работа и встреча", TitleEN: "Work and meetings", ScenarioRU: "Вы на работе: нужно подтвердить встречу, время, имя и контакт.", ScenarioEN: "You are at work: confirm a meeting, time, name, and contact.", Words: []string{"work", "job", "office", "meeting", "email", "phone", "name", "number", "time"}},
	{Code: "doctor", TitleRU: "Врач и самочувствие", TitleEN: "Doctor and health", ScenarioRU: "Вы в клинике: нужно объяснить симптом, попросить врача и договориться о времени.", ScenarioEN: "You are at a clinic: explain a symptom, ask for a doctor, and agree on a time.", Words: []string{"doctor", "hospital", "pain", "medicine", "water", "help", "today", "tomorrow", "morning"}},
	{Code: "shopping", TitleRU: "Магазин и покупка", TitleEN: "Shopping", ScenarioRU: "Вы в магазине: нужно спросить цену, размер, оплату и чек.", ScenarioEN: "You are in a shop: ask about price, size, payment, and a receipt.", Words: []string{"shop", "market", "price", "money", "card", "bag", "size", "new", "good"}},
	{Code: "home", TitleRU: "Дом и быт", TitleEN: "Home and daily life", ScenarioRU: "Вы говорите о доме: нужно назвать адрес, комнату, ключ и бытовую проблему.", ScenarioEN: "You are talking about home: give an address, mention a room, a key, and a daily problem.", Words: []string{"home", "house", "room", "key", "address", "water", "door", "family", "today"}},
	{Code: "study", TitleRU: "Учёба и экзамен", TitleEN: "Study and exams", ScenarioRU: "Вы на занятии: нужно спросить задание, уточнить время и подготовиться к экзамену.", ScenarioEN: "You are in class: ask about a task, clarify time, and prepare for an exam.", Words: []string{"school", "lesson", "exam", "question", "answer", "book", "time", "help", "today"}},
	{Code: "travel", TitleRU: "Транспорт и билеты", TitleEN: "Transport and tickets", ScenarioRU: "Вы на станции или в аэропорту: нужно купить билет, уточнить время и направление.", ScenarioEN: "You are at a station or airport: buy a ticket, clarify time, and confirm direction.", Words: []string{"ticket", "station", "airport", "train", "bus", "time", "street", "city", "tomorrow"}},
	{Code: "services", TitleRU: "Сервисы и документы", TitleEN: "Services and documents", ScenarioRU: "Вы обращаетесь в сервис: нужно назвать номер, документ, проблему и следующий шаг.", ScenarioEN: "You are at a service desk: give a number, show a document, explain a problem, and agree on the next step.", Words: []string{"bank", "phone", "number", "name", "passport", "help", "problem", "address", "today"}},
}

var tutorCourseFunctions = []tutorCourseFunction{
	{Code: "ask-help", TitleRU: "Попросить помощь", TitleEN: "Ask for help", CanDoRU: "попросить помощь в конкретной ситуации", CanDoEN: "ask for help in a specific situation"},
	{Code: "give-details", TitleRU: "Дать детали", TitleEN: "Give details", CanDoRU: "дать короткие детали: имя, время, место или предмет", CanDoEN: "give short details: name, time, place, or item"},
	{Code: "confirm-correct", TitleRU: "Подтвердить и исправить", TitleEN: "Confirm and correct", CanDoRU: "подтвердить информацию и вежливо исправить ошибку", CanDoEN: "confirm information and politely correct a mistake"},
	{Code: "solve-problem", TitleRU: "Решить проблему", TitleEN: "Solve a problem", CanDoRU: "объяснить простую проблему и попросить следующий шаг", CanDoEN: "explain a simple problem and ask for the next step"},
	{Code: "close-next", TitleRU: "Закрыть диалог", TitleEN: "Close and move on", CanDoRU: "закрыть короткий диалог и договориться о следующем действии", CanDoEN: "close a short dialogue and agree on the next action"},
}

var tutorCourseVariants = []tutorCourseVariant{
	{Code: "polite-minimum", TitleRU: "Polite minimum", TitleEN: "Polite minimum", FocusRU: "short polite line", FocusEN: "short polite line", ConstraintRU: "Say only the useful minimum: opener, key word, one detail.", ConstraintEN: "Say only the useful minimum: opener, key word, one detail."},
	{Code: "detail-upgrade", TitleRU: "Detail upgrade", TitleEN: "Detail upgrade", FocusRU: "add one concrete detail", FocusEN: "add one concrete detail", ConstraintRU: "Add one detail that helps the other person act.", ConstraintEN: "Add one detail that helps the other person act."},
	{Code: "clarify-correct", TitleRU: "Clarify and correct", TitleEN: "Clarify and correct", FocusRU: "clarify or repair information", FocusEN: "clarify or repair information", ConstraintRU: "Clarify one unclear point or politely correct one detail.", ConstraintEN: "Clarify one unclear point or politely correct one detail."},
	{Code: "time-pressure", TitleRU: "Time pressure", TitleEN: "Time pressure", FocusRU: "answer fast with time words", FocusEN: "answer fast with time words", ConstraintRU: "Use a time word and keep the line direct.", ConstraintEN: "Use a time word and keep the line direct."},
	{Code: "compare-options", TitleRU: "Compare options", TitleEN: "Compare options", FocusRU: "choose between two options", FocusEN: "choose between two options", ConstraintRU: "Compare two options and choose one.", ConstraintEN: "Compare two options and choose one."},
	{Code: "repair-followup", TitleRU: "Repair follow-up", TitleEN: "Repair follow-up", FocusRU: "ask the next repair question", FocusEN: "ask the next repair question", ConstraintRU: "End with a follow-up question or next step.", ConstraintEN: "End with a follow-up question or next step."},
}

var tutorCourseTopics = []tutorTopic{
	{Code: "hotel", TitleRU: "Отель и поездка", TitleEN: "Hotel and travel", Words: []string{"hotel", "room", "reservation", "passport", "ticket", "station", "airport", "street", "address"}},
	{Code: "food", TitleRU: "Еда и заказ", TitleEN: "Food and ordering", Words: []string{"food", "water", "coffee", "tea", "bread", "menu", "restaurant", "breakfast", "dinner"}},
	{Code: "city", TitleRU: "Город", TitleEN: "City", Words: []string{"city", "street", "house", "shop", "market", "bank", "school", "office", "park"}},
	{Code: "work", TitleRU: "Работа", TitleEN: "Work", Words: []string{"work", "job", "office", "meeting", "email", "phone", "name", "number", "time"}},
	{Code: "doctor", TitleRU: "Врач и самочувствие", TitleEN: "Doctor and health", Words: []string{"doctor", "hospital", "pain", "medicine", "water", "help", "today", "tomorrow", "morning"}},
}

func buildTutorLesson(user userState) (tutorLesson, error) {
	return buildTutorLessonForSequence(user, tutorLessonSequence(user))
}

func buildTutorLessonForSequence(user userState, lessonSequence int) (tutorLesson, error) {
	language := normalizeLearningLanguage(user.LearningLanguage)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	level := normalizeTutorLevel(user.Level)
	if lessonSequence < 0 {
		lessonSequence = 0
	}
	topicIndex := lessonSequence % len(tutorScenarioCourseTopics)
	functionIndex := (lessonSequence / len(tutorScenarioCourseTopics)) % len(tutorCourseFunctions)
	variantIndex := (lessonSequence / (len(tutorScenarioCourseTopics) * len(tutorCourseFunctions))) % len(tutorCourseVariants)
	topic := tutorScenarioCourseTopics[topicIndex]
	courseFunction := tutorCourseFunctions[functionIndex]
	variant := tutorCourseVariants[variantIndex]
	lessonNumber := tutorCourseLessonNumber(lessonSequence)
	words, err := selectTutorLessonWords(user, topic, level, 10, lessonSequence, variant.Code)
	if err != nil {
		return tutorLesson{}, err
	}
	if len(words) < 4 {
		return tutorLesson{}, fmt.Errorf("not enough local A1-A2 vocabulary for %s", language)
	}
	if len(words) > 8 {
		words = words[:8]
	}

	review := selectTutorReviewWords(user, interfaceLanguage, language, 3)
	checks := tutorScenarioChecks(words, topic, courseFunction, interfaceLanguage)
	choice := checks[0]
	languageName := learningLanguageByCode(language).ExplanationName
	topicTitle := tutorTopicTitle(topic, interfaceLanguage)
	first := words[0]
	second := words[tutorMinInt(1, len(words)-1)]
	listeningText, listeningQuestion, listeningExpected := tutorListeningBlock(words, topic, courseFunction, language, interfaceLanguage)
	dialogue, dialoguePrompt, dialogueGoal := tutorScenarioDialogue(words, topic, courseFunction, language, interfaceLanguage)
	title := tutorLocalized(interfaceLanguage, "AI Репетитор", "AI Tutor")

	if interfaceLanguage == "ru" {
		title = "AI Репетитор"
	}

	return tutorLesson{
		ID:                tutorLessonID(language, interfaceLanguage, level, lessonSequence, topic, courseFunction, variant),
		Title:             title,
		Level:             level,
		Topic:             topicTitle,
		VariantCode:       variant.Code,
		VariantTitle:      tutorVariantTitle(variant, interfaceLanguage),
		FocusSkill:        tutorVariantFocus(variant, interfaceLanguage),
		SuccessCriteria:   tutorSuccessCriteria(words, courseFunction, variant, interfaceLanguage),
		LessonNumber:      lessonNumber,
		CourseSize:        tutorCourseSize,
		LearningLanguage:  language,
		InterfaceLanguage: interfaceLanguage,
		DurationMinutes:   8,
		Goal:              tutorLessonGoal(topic, courseFunction, interfaceLanguage) + " " + tutorVariantFocusLine(variant, interfaceLanguage),
		CanDo:             tutorCanDo(level, topic, courseFunction, interfaceLanguage),
		Scenario:          tutorTopicScenario(topic, interfaceLanguage) + " " + tutorVariantConstraint(variant, interfaceLanguage),
		Steps:             tutorLessonSteps(interfaceLanguage),
		Words:             words,
		GrammarTitle:      tutorScenarioGrammarTitle(level, courseFunction, interfaceLanguage),
		Grammar:           tutorScenarioGrammarBody(level, languageName, topic, courseFunction, interfaceLanguage),
		MiniExplanation:   tutorMiniExplanation(level, topic, courseFunction, words, interfaceLanguage),
		Choice:            choice,
		Checks:            checks,
		WritingTask:       tutorWritingTask(first, second, topic, courseFunction, interfaceLanguage),
		WritingExpected:   []string{first.Word, second.Word},
		AnswerVariants:    tutorAnswerVariants(words, topic, courseFunction, language, interfaceLanguage, false),
		ListeningTask:     tutorListeningTask(topic, courseFunction, interfaceLanguage),
		ListeningText:     listeningText,
		ListeningQuestion: listeningQuestion,
		ListeningExpected: listeningExpected,
		PronunciationText: tutorPronunciationLine(words, topic, language),
		Dialogue:          dialogue,
		DialoguePrompt:    dialoguePrompt,
		DialogueGoal:      dialogueGoal,
		DialogueVariants:  tutorAnswerVariants(words, topic, courseFunction, language, interfaceLanguage, true),
		Review:            review,
		ReviewSummary:     tutorReviewSummary(words, topic, courseFunction, interfaceLanguage),
		SRS: tutorLessonSRS{
			Again: tutorLocalized(interfaceLanguage, "ошибка: повторить сегодня", "missed: repeat today"),
			Hard:  tutorLocalized(interfaceLanguage, "трудно: завтра", "hard: tomorrow"),
			Good:  tutorLocalized(interfaceLanguage, "верно: через 2 дня", "correct: in 2 days"),
			Easy:  tutorLocalized(interfaceLanguage, "легко: через неделю", "easy: in one week"),
		},
		NextActions: []string{
			tutorLocalized(interfaceLanguage, "Повторить новые слова в Review", "Review the new words"),
			tutorLocalized(interfaceLanguage, "Проиграть аудио ещё раз", "Replay the listening model"),
			tutorLocalized(interfaceLanguage, "Сохранить лучшую фразу в разговорник", "Save the best phrase to the phrasebook"),
		},
		Source: "local-a1-a2-course-bank",
	}, nil
}

func tutorLessonSteps(interfaceLanguage string) []tutorLessonStep {
	ru := normalizeInterfaceLanguage(interfaceLanguage) == "ru"
	if ru {
		return []tutorLessonStep{
			{Code: "words", Title: "Новые слова", Summary: "слова темы"},
			{Code: "explain", Title: "Мини-объяснение", Summary: "одно правило"},
			{Code: "choice", Title: "Выбор ответа", Summary: "узнать перевод"},
			{Code: "writing", Title: "Письмо", Summary: "собрать фразу"},
			{Code: "listening", Title: "Аудирование", Summary: "повторить вслух"},
			{Code: "dialogue", Title: "Мини-диалог", Summary: "контекст"},
			{Code: "review", Title: "SRS", Summary: "расписание памяти"},
		}
	}
	return []tutorLessonStep{
		{Code: "words", Title: "New words", Summary: "topic vocabulary"},
		{Code: "explain", Title: "Mini explanation", Summary: "one pattern"},
		{Code: "choice", Title: "Choice", Summary: "recognize meaning"},
		{Code: "writing", Title: "Writing", Summary: "build a sentence"},
		{Code: "listening", Title: "Listening", Summary: "catch details"},
		{Code: "pronunciation", Title: "Pronunciation", Summary: "repeat aloud"},
		{Code: "dialogue", Title: "Mini dialogue", Summary: "context"},
		{Code: "review", Title: "SRS", Summary: "memory schedule"},
	}
}

func normalizeTutorLevel(level string) string {
	level = normalizeCEFRLevel(level)
	if level == "" {
		return "A1"
	}
	return level
}

func selectTutorLessonWords(user userState, topic tutorTopic, level string, count int, lessonSequence int, variantCode string) ([]tutorLessonWord, error) {
	language := normalizeLearningLanguage(user.LearningLanguage)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	learned := map[string]bool{}
	for _, learnedWord := range user.LearnedWords {
		if normalizeLearningLanguage(learnedWord.Language) == language {
			learned[learnedWord.ID] = true
			learned[legacyVocabID(learnedWord.ID)] = true
		}
	}
	candidates := make([]vocabWord, 0, 64)
	if err := forEachVocabularyWord(language, func(word vocabWord) bool {
		if learned[word.ID] || learned[legacyVocabID(word.ID)] {
			return true
		}
		if !tutorAcceptWord(word, user, level) {
			return true
		}
		candidates = append(candidates, word)
		return len(candidates) < 1600
	}); err != nil {
		return nil, err
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		left := tutorWordScore(candidates[i], topic, level)
		right := tutorWordScore(candidates[j], topic, level)
		if left != right {
			return left > right
		}
		return tutorStableWordKey(candidates[i]) < tutorStableWordKey(candidates[j])
	})
	if len(candidates) > 1 {
		window := tutorMinInt(len(candidates), tutorMaxInt(48, count*8))
		offset := tutorStableIndex(fmt.Sprintf("%s:%s:%s:%s:%d", language, level, topic.Code, variantCode, lessonSequence), window)
		if offset > 0 {
			rotated := make([]vocabWord, 0, len(candidates))
			rotated = append(rotated, candidates[offset:window]...)
			rotated = append(rotated, candidates[:offset]...)
			rotated = append(rotated, candidates[window:]...)
			candidates = rotated
		}
	}
	words := make([]tutorLessonWord, 0, count)
	seen := map[string]bool{}
	for _, candidate := range candidates {
		item := tutorLessonWordFromVocab(candidate, interfaceLanguage, topic)
		if item.Word == "" || item.Translation == "" || seen[strings.ToLower(item.Word)] {
			continue
		}
		seen[strings.ToLower(item.Word)] = true
		words = append(words, item)
		if len(words) >= count {
			break
		}
	}
	return words, nil
}

func selectTutorReviewWords(user userState, interfaceLanguage string, language string, count int) []tutorLessonWord {
	review := make([]tutorLessonWord, 0, count)
	for _, learned := range user.LearnedWords {
		if len(review) >= count {
			break
		}
		if normalizeLearningLanguage(learned.Language) != language || learnedWordMastered(learned) {
			continue
		}
		word := learnedEntryVocabWord(learned)
		if !tutorLooksClean(word.English) || !tutorLooksClean(tutorTranslationForUser(word, interfaceLanguage)) {
			continue
		}
		item := tutorLessonWordFromVocab(word, interfaceLanguage, tutorCourseTopics[0])
		if item.Word != "" && item.Translation != "" {
			review = append(review, item)
		}
	}
	return review
}

func tutorAcceptWord(word vocabWord, user userState, level string) bool {
	wordLevel := normalizeCEFRLevel(word.Level)
	if wordLevel != "A1" && wordLevel != "A2" {
		return false
	}
	if cefrRank(wordLevel) > cefrRank(level)+1 {
		return false
	}
	if !tutorLooksClean(word.English) || !tutorLooksClean(tutorTranslationForUser(word, user.InterfaceLanguage)) {
		return false
	}
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" && normalizeLearningLanguage(word.Language) != "ru" && looksLikeEnglishOnly(tutorTranslationForUser(word, "ru")) {
		return false
	}
	if strings.Count(cleanDictionaryDisplay(word.English), " ") > 3 {
		return false
	}
	return true
}

func tutorWordScore(word vocabWord, topic tutorTopic, level string) int {
	score := 0
	if normalizeCEFRLevel(word.Level) == level {
		score += 50
	}
	if strings.EqualFold(word.Topic, "curated-core") || strings.HasPrefix(strings.ToLower(word.Source), "curated") {
		score += 25
	}
	if word.FrequencyRank > 0 {
		score += tutorMaxInt(0, 20-word.FrequencyRank/100)
	}
	needle := strings.ToLower(cleanDictionaryDisplay(word.English))
	for _, expected := range topic.Words {
		if needle == expected {
			score += 80
			break
		}
		if strings.Contains(needle, expected) || strings.Contains(strings.ToLower(word.Topic), expected) {
			score += 20
		}
	}
	return score
}

func tutorLessonWordFromVocab(word vocabWord, interfaceLanguage string, topic tutorTopic) tutorLessonWord {
	target := firstDictionaryValue(word.English)
	translation := tutorTranslationForUser(word, interfaceLanguage)
	context := wordContext(word, interfaceLanguage)
	if context == "" {
		context = firstDictionaryValue(word.Russian)
	}
	example, exampleTranslation := tutorScenarioExample(target, translation, interfaceLanguage, topic, normalizeLearningLanguage(word.Language))
	return tutorLessonWord{
		ID:           word.ID,
		Word:         target,
		Translation:  translation,
		Context:      context,
		Example:      example,
		ExampleRu:    exampleTranslation,
		PartOfSpeech: strings.TrimSpace(word.PartOfSpeech),
		Topic:        tutorTopicTitle(topic, interfaceLanguage),
		Level:        normalizeTutorLevel(word.Level),
		Frequency:    word.FrequencyRank,
		AudioText:    target,
	}
}

func tutorChoiceFromWords(words []tutorLessonWord, interfaceLanguage string) tutorLessonChoice {
	correct := words[0]
	options := make([]tutorChoiceOption, 0, tutorMinInt(4, len(words)))
	for index, word := range words {
		if index >= 4 {
			break
		}
		quality := "distractor"
		label := tutorLocalized(interfaceLanguage, "Похоже, но не здесь", "Close, but not here")
		why := tutorLocalized(interfaceLanguage, "Смотрите на слово в задании, а не на знакомый перевод рядом.", "Match the target word, not just a familiar meaning nearby.")
		if word.ID == correct.ID {
			quality = "correct"
			label = tutorLocalized(interfaceLanguage, "Точный перевод", "Exact meaning")
			why = tutorLocalized(interfaceLanguage, "Это значение нужно узнать перед тем, как собрать реплику.", "You need this meaning before building the line.")
		}
		options = append(options, tutorChoiceOption{ID: word.ID, Text: word.Translation, Label: label, Why: why, Skill: "meaning", Quality: quality, Avoid: word.ID != correct.ID})
	}
	return tutorLessonChoice{
		Prompt:          tutorLocalized(interfaceLanguage, "Выбери перевод: "+correct.Word, "Choose the translation: "+correct.Word),
		Options:         options,
		CorrectAnswerID: correct.ID,
	}
}

func tutorLessonSequence(user userState) int {
	sequence := user.WordLessonCount + user.LessonCount
	if sequence < 0 {
		return 0
	}
	return sequence
}

func tutorCourseLessonNumber(lessonSequence int) int {
	if lessonSequence < 0 {
		return 1
	}
	return lessonSequence%tutorCourseSize + 1
}

func tutorLessonID(language string, interfaceLanguage string, level string, lessonSequence int, topic tutorTopic, function tutorCourseFunction, variant tutorCourseVariant) string {
	return fmt.Sprintf(
		"tutor-%s-%s-%s-%04d-%s-%s-%s",
		normalizeLearningLanguage(language),
		normalizeInterfaceLanguage(interfaceLanguage),
		strings.ToLower(normalizeTutorLevel(level)),
		lessonSequence+1,
		topic.Code,
		function.Code,
		variant.Code,
	)
}

func tutorVariantTitle(variant tutorCourseVariant, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && strings.TrimSpace(variant.TitleRU) != "" {
		return variant.TitleRU
	}
	return variant.TitleEN
}

func tutorVariantFocus(variant tutorCourseVariant, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && strings.TrimSpace(variant.FocusRU) != "" {
		return variant.FocusRU
	}
	return variant.FocusEN
}

func tutorVariantFocusLine(variant tutorCourseVariant, interfaceLanguage string) string {
	focus := tutorVariantFocus(variant, interfaceLanguage)
	if focus == "" {
		return ""
	}
	return tutorLocalized(interfaceLanguage, "Focus: "+focus+".", "Focus: "+focus+".")
}

func tutorVariantConstraint(variant tutorCourseVariant, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && strings.TrimSpace(variant.ConstraintRU) != "" {
		return variant.ConstraintRU
	}
	return variant.ConstraintEN
}

func tutorSuccessCriteria(words []tutorLessonWord, function tutorCourseFunction, variant tutorCourseVariant, interfaceLanguage string) []string {
	first := tutorWordOrFallback(words, 0, "key word").Word
	second := tutorWordOrFallback(words, 1, "detail").Word
	return []string{
		"Use the lesson goal: " + function.CanDoEN,
		"Include " + first + " and one useful detail such as " + second,
		"Variant focus: " + tutorVariantFocus(variant, interfaceLanguage),
	}
}

func tutorLessonStorageContext(user userState) (string, string, string) {
	return normalizeLearningLanguage(user.LearningLanguage), normalizeInterfaceLanguage(user.InterfaceLanguage), normalizeTutorLevel(user.Level)
}

func tutorLessonMatchesContext(lesson tutorLesson, language string, interfaceLanguage string, level string) bool {
	return normalizeLearningLanguage(lesson.LearningLanguage) == normalizeLearningLanguage(language) &&
		normalizeInterfaceLanguage(lesson.InterfaceLanguage) == normalizeInterfaceLanguage(interfaceLanguage) &&
		normalizeTutorLevel(lesson.Level) == normalizeTutorLevel(level)
}

func tutorReusableLessonUser(user userState) userState {
	core := user
	core.TelegramID = 0
	core.LearnedWords = nil
	core.Mistakes = nil
	core.Phrasebook = nil
	core.LessonCount = 0
	core.WordLessonCount = 0
	return core
}

func tutorPersonalizeLessonForUser(lesson tutorLesson, user userState) tutorLesson {
	language := normalizeLearningLanguage(lesson.LearningLanguage)
	if language == "" {
		language = normalizeLearningLanguage(user.LearningLanguage)
	}
	interfaceLanguage := normalizeInterfaceLanguage(lesson.InterfaceLanguage)
	if interfaceLanguage == "" {
		interfaceLanguage = normalizeInterfaceLanguage(user.InterfaceLanguage)
	}
	review := selectTutorReviewWords(user, interfaceLanguage, language, 3)
	if len(review) > 0 {
		lesson.Review = review
	}
	return lesson
}

func tutorLessonGenerationAttemptLimit() int {
	return tutorMaxInt(64, tutorCourseSize)
}

func tutorLessonGoal(topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("Сценарная задача: %s в теме «%s».", strings.ToLower(function.TitleRU), tutorTopicTitle(topic, interfaceLanguage))
	}
	return fmt.Sprintf("Scenario task: %s in %s.", function.TitleEN, tutorTopicTitle(topic, interfaceLanguage))
}

func tutorCanDo(level string, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("После урока вы сможете на уровне %s %s.", level, function.CanDoRU)
	}
	return fmt.Sprintf("After this %s lesson you can %s.", level, function.CanDoEN)
}

func tutorTopicScenario(topic tutorTopic, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" && topic.ScenarioRU != "" {
		return topic.ScenarioRU
	}
	if topic.ScenarioEN != "" {
		return topic.ScenarioEN
	}
	return tutorTopicTitle(topic, interfaceLanguage)
}

func tutorScenarioChecks(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) []tutorLessonChoice {
	if len(words) < 4 {
		return []tutorLessonChoice{tutorChoiceFromWords(words, interfaceLanguage)}
	}
	first := words[0]
	second := words[tutorMinInt(1, len(words)-1)]
	third := words[tutorMinInt(2, len(words)-1)]
	fourth := words[tutorMinInt(3, len(words)-1)]
	ru := normalizeInterfaceLanguage(interfaceLanguage) == "ru"
	promptBestLine := tutorLocalized(interfaceLanguage, "Какая реплика лучше подходит к сценарию?", "Which line best fits the scenario?")
	promptDetail := tutorLocalized(interfaceLanguage, "Какую деталь нужно назвать в этом диалоге?", "Which detail do you need to give in this dialogue?")
	promptNext := tutorLocalized(interfaceLanguage, "Какой следующий шаг звучит естественно?", "Which next step sounds natural?")
	feedback := tutorLocalized(interfaceLanguage, "Верно: это реплика из сценария, а не просто перевод слова.", "Correct: this is a scenario line, not just a word translation.")

	bestLine := tutorCheckLine(topic, first, second, ru)
	detailLine := tutorCheckDetailLine(topic, second, third, ru)
	nextLine := tutorCheckNextLine(topic, third, fourth, ru)
	return []tutorLessonChoice{
		{
			Prompt:          promptBestLine,
			Options:         tutorCheckOptions("scenario", bestLine, tutorDistractorLines(topic, first, second, ru), interfaceLanguage),
			CorrectAnswerID: "scenario:correct",
			Feedback:        feedback,
		},
		{
			Prompt:          promptDetail,
			Options:         tutorCheckOptions("detail", detailLine, []string{first.Translation, fourth.Translation, tutorLocalized(interfaceLanguage, "Не отвечать и ждать", "Say nothing and wait")}, interfaceLanguage),
			CorrectAnswerID: "detail:correct",
			Feedback:        tutorLocalized(interfaceLanguage, "Верно: вы выбрали деталь, которая двигает диалог дальше.", "Correct: you picked the detail that moves the dialogue forward."),
		},
		{
			Prompt:          promptNext,
			Options:         tutorCheckOptions("next", nextLine, tutorNextDistractors(topic, ru), interfaceLanguage),
			CorrectAnswerID: "next:correct",
			Feedback:        tutorLocalized(interfaceLanguage, "Верно: это закрывает шаг и готовит ваш ответ.", "Correct: this closes the step and prepares your answer."),
		},
	}
}

func tutorCheckOptions(prefix string, correct string, distractors []string, interfaceLanguage string) []tutorChoiceOption {
	options := []tutorChoiceOption{{
		ID:      prefix + ":correct",
		Text:    correct,
		Label:   tutorLocalized(interfaceLanguage, "Рабочая реплика", "Usable line"),
		Why:     tutorLocalized(interfaceLanguage, "Есть действие, деталь и понятная цель. Такую фразу можно сказать в реальном диалоге.", "It has an action, a detail, and a clear goal. You can say it in a real dialogue."),
		Skill:   prefix,
		Quality: "correct",
	}}
	for index, distractor := range distractors {
		distractor = strings.TrimSpace(distractor)
		if distractor == "" || strings.EqualFold(distractor, correct) {
			continue
		}
		options = append(options, tutorChoiceOption{
			ID:      fmt.Sprintf("%s:wrong:%d", prefix, index),
			Text:    distractor,
			Label:   tutorWeakOptionLabel(index, interfaceLanguage),
			Why:     tutorWeakOptionWhy(index, interfaceLanguage),
			Skill:   prefix,
			Quality: "weak",
			Avoid:   true,
		})
		if len(options) >= 4 {
			break
		}
	}
	return options
}

func tutorWeakOptionLabel(index int, interfaceLanguage string) string {
	labelsRU := []string{"Просто перевод", "Повтор слов", "Пустая реакция"}
	labelsEN := []string{"Translation only", "Word repeat", "Empty reaction"}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return labelsRU[tutorMinInt(index, len(labelsRU)-1)]
	}
	return labelsEN[tutorMinInt(index, len(labelsEN)-1)]
}

func tutorWeakOptionWhy(index int, interfaceLanguage string) string {
	whyRU := []string{
		"Перевод слова не решает задачу собеседника.",
		"Повторить слова легче, но это не звучит как ответ.",
		"В реальной сцене собеседник ждёт действие или деталь.",
	}
	whyEN := []string{
		"A word translation does not solve the other person's task.",
		"Repeating words is easier, but it does not sound like an answer.",
		"In a real scene, the other person expects an action or detail.",
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return whyRU[tutorMinInt(index, len(whyRU)-1)]
	}
	return whyEN[tutorMinInt(index, len(whyEN)-1)]
}

func tutorCheckLine(topic tutorTopic, first tutorLessonWord, second tutorLessonWord, ru bool) string {
	if ru {
		switch topic.Code {
		case "hotel":
			return fmt.Sprintf("I have a reservation. Can I have a %s?", second.Word)
		case "doctor":
			return fmt.Sprintf("I have %s. Can I see a %s today?", first.Word, second.Word)
		case "food":
			return fmt.Sprintf("Can I see the %s and have %s?", first.Word, second.Word)
		case "travel":
			return fmt.Sprintf("I need a %s for tomorrow. What time is it?", first.Word)
		}
	}
	switch topic.Code {
	case "hotel":
		return fmt.Sprintf("I have a reservation. Can I have a %s?", second.Word)
	case "doctor":
		return fmt.Sprintf("I have %s. Can I see a %s today?", first.Word, second.Word)
	case "food":
		return fmt.Sprintf("Can I see the %s and have %s?", first.Word, second.Word)
	case "travel":
		return fmt.Sprintf("I need a %s for tomorrow. What time is it?", first.Word)
	}
	return fmt.Sprintf("Can you help me with %s and %s?", first.Word, second.Word)
}

func tutorCheckDetailLine(topic tutorTopic, second tutorLessonWord, third tutorLessonWord, ru bool) string {
	switch topic.Code {
	case "hotel", "services":
		return fmt.Sprintf("%s: %s", third.Word, third.Translation)
	case "work", "study", "travel":
		return fmt.Sprintf("%s: %s", second.Word, second.Translation)
	case "doctor":
		return fmt.Sprintf("%s: %s", second.Word, second.Translation)
	}
	return fmt.Sprintf("%s: %s", second.Word, second.Translation)
}

func tutorCheckNextLine(topic tutorTopic, third tutorLessonWord, fourth tutorLessonWord, ru bool) string {
	switch topic.Code {
	case "hotel":
		return fmt.Sprintf("Here is my %s. Thank you.", third.Word)
	case "doctor":
		return fmt.Sprintf("Can I come this %s?", fourth.Word)
	case "food":
		return fmt.Sprintf("Thank you. Can I pay by %s?", fourth.Word)
	case "travel":
		return fmt.Sprintf("Great. I will take this %s.", third.Word)
	}
	return fmt.Sprintf("Thank you. The next step is %s.", fourth.Word)
}

func tutorDistractorLines(topic tutorTopic, first tutorLessonWord, second tutorLessonWord, ru bool) []string {
	return []string{
		fmt.Sprintf("%s / %s", first.Translation, second.Translation),
		fmt.Sprintf("I repeat: %s, %s.", first.Word, second.Word),
		"Yes. No. Maybe.",
	}
}

func tutorNextDistractors(topic tutorTopic, ru bool) []string {
	if ru {
		return []string{"Повторить список слов без ответа", "Выбрать случайный перевод", "Закрыть приложение"}
	}
	return []string{"Repeat the word list only", "Pick a random translation", "Close the app"}
}

func tutorScenarioGrammarTitle(level string, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("Паттерн %s: %s", level, strings.ToLower(function.TitleRU))
	}
	return fmt.Sprintf("%s pattern: %s", level, function.TitleEN)
}

func tutorScenarioGrammarBody(level string, languageName string, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		switch level {
		case "A1":
			return "Соберите короткую реплику: вежливое начало + ключевое слово + одна деталь. В этом уроке важна не форма слова сама по себе, а то, какую задачу она решает в сцене."
		case "A2":
			return "Добавьте уточнение: время, место, имя или документ. Ответ должен звучать как реплика в диалоге, а не как список слов."
		default:
			return "Сначала выберите коммуникативную цель, затем переформулируйте мысль естественно: причина, деталь, следующий шаг. Это тренирует контекст и нюанс, а не механический перевод."
		}
	}
	switch level {
	case "A1":
		return "Build a short " + languageName + " line: polite opener + key word + one detail. The word matters because it solves the scene task."
	case "A2":
		return "Add one clarification: time, place, name, or document. Make it a dialogue line, not a word list."
	default:
		return "Start from the communicative goal, then rephrase naturally: reason, detail, and next step."
	}
}

func tutorMiniExplanation(level string, topic tutorTopic, function tutorCourseFunction, words []tutorLessonWord, interfaceLanguage string) string {
	first := ""
	second := ""
	if len(words) > 0 {
		first = words[0].Word
	}
	if len(words) > 1 {
		second = words[1].Word
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("Мини-объяснение: сцена «%s» учит %s. Слова «%s» и «%s» нужны не для угадывания перевода, а чтобы собрать рабочую реплику: запрос, деталь и следующий шаг.", tutorTopicTitle(topic, interfaceLanguage), function.CanDoRU, first, second)
	}
	return fmt.Sprintf("Mini explanation: the %s scene trains you to %s. Use \"%s\" and \"%s\" to build a usable line: request, detail, and next step.", tutorTopicTitle(topic, interfaceLanguage), function.CanDoEN, first, second)
}

func tutorWritingTask(first tutorLessonWord, second tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("Напишите одну рабочую реплику для сценария «%s». Используйте «%s» и по возможности «%s».", tutorTopicTitle(topic, interfaceLanguage), first.Word, second.Word)
	}
	return fmt.Sprintf("Write one usable line for the \"%s\" scenario. Use \"%s\" and, if possible, \"%s\".", tutorTopicTitle(topic, interfaceLanguage), first.Word, second.Word)
}

func tutorAnswerVariants(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string, dialogue bool) []tutorAnswerVariant {
	first := tutorWordOrFallback(words, 0, "reservation")
	second := tutorWordOrFallback(words, 1, "room")
	third := tutorWordOrFallback(words, 2, "passport")
	fourth := tutorWordOrFallback(words, 3, "today")
	ru := normalizeInterfaceLanguage(interfaceLanguage) == "ru"
	canDo := function.CanDoEN
	if ru {
		canDo = function.CanDoRU
	}

	lines := tutorVariantLines(topic, language, first, second, third, fourth, dialogue)
	labels := []string{"Base answer", "Natural answer", "Stronger answer", "Common weak answer"}
	why := []string{
		"Short, clear, and enough for the task.",
		"Sounds like a real line because it adds one useful detail.",
		"Best model: polite opener, task, detail, and next step.",
		"Too thin: it names words but does not solve the conversation.",
	}
	useCases := []string{"safe minimum", "daily conversation", "best practice model", "avoid"}
	if ru {
		labels = []string{"Базовый ответ", "Естественно", "Сильный вариант", "Частая слабая ошибка"}
		why = []string{
			"Коротко и достаточно для задачи.",
			"Звучит как живая реплика, потому что добавлена полезная деталь.",
			"Лучший образец: вежливое начало, задача, деталь и следующий шаг.",
			"Слишком тонко: слова названы, но разговорная задача не решена.",
		}
		useCases = []string{"минимум без риска", "живой разговор", "лучший образец", "не повторять"}
	}

	variants := make([]tutorAnswerVariant, 0, len(lines))
	seen := map[string]bool{}
	for index, line := range lines {
		line = strings.TrimSpace(line)
		key := strings.ToLower(line)
		if line == "" || seen[key] {
			continue
		}
		seen[key] = true
		variants = append(variants, tutorAnswerVariant{
			ID:      fmt.Sprintf("%s:%d", map[bool]string{true: "dialogue", false: "writing"}[dialogue], index+1),
			Label:   labels[tutorMinInt(index, len(labels)-1)],
			Text:    line,
			Why:     why[tutorMinInt(index, len(why)-1)] + " " + tutorLocalized(interfaceLanguage, "Цель урока: "+canDo+".", "Lesson goal: "+canDo+"."),
			Level:   []string{"A1", "A2", "B1", "avoid"}[tutorMinInt(index, 3)],
			UseCase: useCases[tutorMinInt(index, len(useCases)-1)],
			Avoid:   index == len(lines)-1,
		})
	}
	return variants
}

func tutorVariantLines(topic tutorTopic, language string, first tutorLessonWord, second tutorLessonWord, third tutorLessonWord, fourth tutorLessonWord, dialogue bool) []string {
	if normalizeLearningLanguage(language) != "en" {
		base := strings.TrimSpace(first.Word + " " + second.Word)
		natural := strings.TrimSpace(first.Word + " " + second.Word + " " + third.Word)
		strong := strings.TrimSpace(first.Word + " " + second.Word + " " + third.Word + " " + fourth.Word)
		weak := strings.TrimSpace(first.Word + " / " + second.Word)
		return []string{base, natural, strong, weak}
	}
	if dialogue {
		switch topic.Code {
		case "hotel":
			return []string{
				fmt.Sprintf("Yes, here is my %s.", third.Word),
				fmt.Sprintf("Yes, here is my %s. The %s is under my name.", third.Word, first.Word),
				fmt.Sprintf("Yes, here is my %s. The %s is under my name, and I need a quiet %s.", third.Word, first.Word, second.Word),
				fmt.Sprintf("%s. %s. Yes.", first.Word, second.Word),
			}
		case "doctor":
			return []string{
				fmt.Sprintf("Yes, I can come %s.", fourth.Word),
				fmt.Sprintf("Yes, I can come %s. I have %s.", fourth.Word, first.Word),
				fmt.Sprintf("Yes, I can come %s. I have %s and need a %s.", fourth.Word, first.Word, second.Word),
				fmt.Sprintf("%s, %s.", first.Word, second.Word),
			}
		case "food":
			return []string{
				fmt.Sprintf("%s, please.", second.Word),
				fmt.Sprintf("I would like %s, please.", second.Word),
				fmt.Sprintf("I would like %s, please. Could I also see the %s?", second.Word, first.Word),
				fmt.Sprintf("%s / %s.", first.Word, second.Word),
			}
		case "travel":
			return []string{
				fmt.Sprintf("I need a %s.", first.Word),
				fmt.Sprintf("I need a %s for %s.", first.Word, fourth.Word),
				fmt.Sprintf("I need a %s for %s. What time does it leave?", first.Word, fourth.Word),
				fmt.Sprintf("%s. Time?", first.Word),
			}
		}
		return []string{
			fmt.Sprintf("Yes, the detail is %s.", second.Word),
			fmt.Sprintf("Yes, the detail is %s. I need help with %s.", second.Word, first.Word),
			fmt.Sprintf("Yes, the detail is %s. Could you help me with the next step for %s?", second.Word, first.Word),
			fmt.Sprintf("%s, %s, okay.", first.Word, second.Word),
		}
	}
	switch topic.Code {
	case "hotel":
		return []string{
			fmt.Sprintf("I have a %s.", first.Word),
			fmt.Sprintf("I have a %s and need a quiet %s.", first.Word, second.Word),
			fmt.Sprintf("Good evening. I have a %s and need a quiet %s for %s.", first.Word, second.Word, fourth.Word),
			fmt.Sprintf("%s, %s, %s.", first.Word, second.Word, third.Word),
		}
	case "doctor":
		return []string{
			fmt.Sprintf("I have %s.", first.Word),
			fmt.Sprintf("I have %s and need a %s.", first.Word, second.Word),
			fmt.Sprintf("Hello. I have %s and need a %s %s.", first.Word, second.Word, fourth.Word),
			fmt.Sprintf("%s. %s. Help.", first.Word, second.Word),
		}
	case "food":
		return []string{
			fmt.Sprintf("Can I see the %s?", first.Word),
			fmt.Sprintf("Can I see the %s and have %s?", first.Word, second.Word),
			fmt.Sprintf("Hello. Can I see the %s and have %s, please?", first.Word, second.Word),
			fmt.Sprintf("%s / %s / please.", first.Word, second.Word),
		}
	case "travel":
		return []string{
			fmt.Sprintf("I need a %s.", first.Word),
			fmt.Sprintf("I need a %s for %s.", first.Word, fourth.Word),
			fmt.Sprintf("Hello. I need a %s for %s. What time does it leave?", first.Word, fourth.Word),
			fmt.Sprintf("%s. %s.", first.Word, fourth.Word),
		}
	}
	return []string{
		fmt.Sprintf("Can you help me with %s?", first.Word),
		fmt.Sprintf("Can you help me with %s? The detail is %s.", first.Word, second.Word),
		fmt.Sprintf("Could you help me with %s? The detail is %s, and the next step is %s.", first.Word, second.Word, third.Word),
		fmt.Sprintf("%s, %s, maybe.", first.Word, second.Word),
	}
}

func tutorListeningTask(topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return "Прослушайте модель и напишите, какие 2-3 ключевые детали вы услышали. Репетитор проверит, есть ли смысловые слова."
	}
	return "Listen to the model and type 2-3 key details you heard. The tutor checks the meaning words."
}

func tutorListeningBlock(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string) (string, string, []string) {
	first := tutorWordOrFallback(words, 0, "reservation")
	second := tutorWordOrFallback(words, 1, "room")
	third := tutorWordOrFallback(words, 2, "passport")
	fourth := tutorWordOrFallback(words, 3, "today")
	if normalizeLearningLanguage(language) != "en" {
		text := strings.Join([]string{first.Word, second.Word, third.Word, fourth.Word}, ". ")
		question := tutorLocalized(interfaceLanguage, "Какие ключевые слова вы услышали?", "Which key words did you hear?")
		return text + ".", question, []string{first.Word, second.Word, third.Word}
	}
	var text string
	switch topic.Code {
	case "hotel":
		text = fmt.Sprintf("Good evening. I have a %s. Can I have a quiet %s for one night? Here is my %s.", first.Word, second.Word, third.Word)
	case "doctor":
		text = fmt.Sprintf("Hello. I have %s in the morning. Can I see a %s %s?", first.Word, second.Word, fourth.Word)
	case "food":
		text = fmt.Sprintf("Hello. Can I see the %s? I would like %s and %s, please.", first.Word, second.Word, third.Word)
	case "travel":
		text = fmt.Sprintf("Hello. I need a %s to the %s for %s. What time does it leave?", first.Word, second.Word, fourth.Word)
	default:
		text = fmt.Sprintf("Hello. I need help with %s. My detail is %s, and the next step is %s.", first.Word, second.Word, third.Word)
	}
	question := tutorLocalized(interfaceLanguage, "Какие 2-3 детали важны в аудио?", "Which 2-3 details matter in the audio?")
	return text, question, []string{first.Word, second.Word, third.Word}
}

func tutorScenarioDialogue(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string) ([]string, string, string) {
	first := tutorWordOrFallback(words, 0, "reservation")
	second := tutorWordOrFallback(words, 1, "room")
	third := tutorWordOrFallback(words, 2, "passport")
	if normalizeLearningLanguage(language) != "en" {
		lines := []string{
			fmt.Sprintf("Tutor: %s?", tutorTopicTitle(topic, interfaceLanguage)),
			fmt.Sprintf("Learner: %s / %s.", first.Word, second.Word),
			fmt.Sprintf("Tutor: %s?", third.Word),
		}
		prompt := tutorLocalized(interfaceLanguage, "Ответьте одной короткой репликой с 1-2 словами урока.", "Answer with one short line using 1-2 lesson words.")
		return lines, prompt, tutorCanDo("A1", topic, function, interfaceLanguage)
	}
	var lines []string
	var prompt string
	switch topic.Code {
	case "hotel":
		lines = []string{
			"Receptionist: Good evening. Do you have a reservation?",
			fmt.Sprintf("Guest: Yes, I have a %s under Ivan Petrov.", first.Word),
			fmt.Sprintf("Receptionist: May I see your %s, please?", third.Word),
		}
		prompt = fmt.Sprintf("Answer as the guest: confirm the %s and offer the %s.", first.Word, third.Word)
	case "doctor":
		lines = []string{
			"Receptionist: Hello. How can I help?",
			fmt.Sprintf("Patient: I have %s and I need a %s.", first.Word, second.Word),
			"Receptionist: Can you come this morning?",
		}
		prompt = "Answer as the patient: say yes/no and give one detail."
	case "food":
		lines = []string{
			"Server: Hello. What would you like?",
			fmt.Sprintf("Guest: Can I see the %s, please?", first.Word),
			fmt.Sprintf("Server: Sure. Would you like %s or %s?", second.Word, third.Word),
		}
		prompt = "Answer as the guest: choose one item and say please."
	default:
		lines = []string{
			"Assistant: Hello. How can I help?",
			fmt.Sprintf("Learner: I need help with %s.", first.Word),
			fmt.Sprintf("Assistant: What detail should I write: %s or %s?", second.Word, third.Word),
		}
		prompt = "Answer with the correct detail and one next step."
	}
	return lines, tutorLocalized(interfaceLanguage, "Ответьте на последнюю реплику на изучаемом языке.", prompt), tutorCanDo("A1", topic, function, interfaceLanguage)
}

func tutorReviewSummary(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) []string {
	count := tutorMinInt(4, len(words))
	summary := make([]string, 0, count+2)
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		summary = append(summary, "Вы закрыли сценарий: "+tutorTopicTitle(topic, interfaceLanguage))
		summary = append(summary, "Цель: "+function.CanDoRU)
	} else {
		summary = append(summary, "Scenario completed: "+tutorTopicTitle(topic, interfaceLanguage))
		summary = append(summary, "Can-do: "+function.CanDoEN)
	}
	for index := 0; index < count; index++ {
		word := words[index]
		summary = append(summary, fmt.Sprintf("%s - %s", word.Word, word.Translation))
	}
	return summary
}

func tutorPronunciationLine(words []tutorLessonWord, topic tutorTopic, language string) string {
	if normalizeLearningLanguage(language) != "en" {
		for _, word := range words {
			example := strings.TrimSpace(word.Example)
			if example != "" {
				return example
			}
		}
		return strings.TrimSpace(strings.Join(tutorLessonWordTexts(words, 4), " "))
	}
	word := func(index int, fallback string) string {
		if index >= 0 && index < len(words) && strings.TrimSpace(words[index].Word) != "" {
			return strings.TrimSpace(words[index].Word)
		}
		return fallback
	}
	switch topic.Code {
	case "hotel":
		return "I have a reservation and I need a quiet room, please."
	case "doctor":
		return "I have pain here and I need to see a doctor today."
	case "food":
		return "Could I have the menu and a glass of water, please?"
	case "transport":
		return "I need a ticket for the next train, please."
	case "work":
		return "Can we check the meeting time and the main task?"
	case "study":
		return "Could you help me with this lesson question today?"
	default:
		return fmt.Sprintf("Please help me with %s and %s today.", word(0, "this"), word(1, "practice"))
	}
}

func tutorLessonWordTexts(words []tutorLessonWord, limit int) []string {
	out := make([]string, 0, tutorMinInt(limit, len(words)))
	for _, word := range words {
		text := strings.TrimSpace(word.Word)
		if text == "" {
			continue
		}
		out = append(out, text)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func tutorScenarioExample(word string, translation string, interfaceLanguage string, topic tutorTopic, language string) (string, string) {
	if normalizeLearningLanguage(language) != "en" {
		return word, translation
	}
	example := fmt.Sprintf("Can you help me with %s?", word)
	switch topic.Code {
	case "hotel":
		if strings.EqualFold(word, "hotel") {
			example = "The hotel is near the station."
		} else if strings.EqualFold(word, "room") {
			example = "Can I have a quiet room?"
		} else if strings.EqualFold(word, "passport") {
			example = "Here is my passport."
		} else if strings.EqualFold(word, "reservation") {
			example = "I have a reservation for tonight."
		}
	case "doctor":
		if strings.EqualFold(word, "pain") {
			example = "I have pain here."
		} else if strings.EqualFold(word, "doctor") {
			example = "Can I see a doctor today?"
		} else if strings.EqualFold(word, "medicine") {
			example = "Do I need medicine?"
		}
	case "food":
		if strings.EqualFold(word, "menu") {
			example = "Can I see the menu?"
		} else if strings.EqualFold(word, "water") {
			example = "A glass of water, please."
		} else if strings.EqualFold(word, "coffee") {
			example = "Can I have coffee with breakfast?"
		}
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return example, fmt.Sprintf("Смысл в уроке: %s.", translation)
	}
	return example, translation
}

func tutorWordOrFallback(words []tutorLessonWord, index int, fallback string) tutorLessonWord {
	if index >= 0 && index < len(words) {
		return words[index]
	}
	return tutorLessonWord{ID: fallback, Word: fallback, Translation: fallback}
}

func tutorDialogue(words []tutorLessonWord, interfaceLanguage string) []string {
	if len(words) < 3 {
		return nil
	}
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return []string{
			"А: " + words[0].Example + " " + words[0].ExampleRu,
			"Б: " + words[1].Example + " " + words[1].ExampleRu,
			"А: " + words[2].Example + " " + words[2].ExampleRu,
		}
	}
	return []string{
		"A: " + words[0].Example,
		"B: " + words[1].Example,
		"A: " + words[2].Example,
	}
}

func tutorExample(word string, translation string, interfaceLanguage string) (string, string) {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return fmt.Sprintf("I need %s.", word), fmt.Sprintf("Мне нужно: %s.", translation)
	}
	return fmt.Sprintf("I need %s.", word), fmt.Sprintf("%s.", translation)
}

func tutorGrammarTitle(level string, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		if level == "A1" {
			return "Паттерн A1: короткое утверждение"
		}
		return "Паттерн A2: просьба и уточнение"
	}
	if level == "A1" {
		return "A1 pattern: short statement"
	}
	return "A2 pattern: request and clarification"
}

func tutorGrammarBody(level string, languageName string, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		if level == "A1" {
			return "Собери фразу по схеме: кто/что + действие или состояние + одно полезное слово. Держи предложение коротким и понятным."
		}
		return "Собери просьбу: вежливое начало + нужное слово + уточнение времени, места или количества."
	}
	if level == "A1" {
		return "Build a short " + languageName + " sentence: subject + action or state + one useful word."
	}
	return "Build a polite request: opener + target word + time, place, or quantity detail."
}

func tutorTopicTitle(topic tutorTopic, interfaceLanguage string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return topic.TitleRU
	}
	return topic.TitleEN
}

func tutorLocalized(interfaceLanguage string, ru string, en string) string {
	if normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		return ru
	}
	return en
}

func tutorTranslationForUser(word vocabWord, interfaceLanguage string) string {
	translation := wordTranslation(word, interfaceLanguage)
	if translation != "" {
		return translation
	}
	if normalizeLearningLanguage(word.Language) == "ru" && normalizeInterfaceLanguage(interfaceLanguage) == "ru" {
		if translation := wordTranslationForLanguage(word, "en"); translation != "" {
			return translation
		}
		if word.English != "" {
			return firstDictionaryValue(word.English)
		}
	}
	return ""
}

func tutorLooksClean(value string) bool {
	value = cleanDictionaryDisplay(value)
	if value == "" || !utf8.ValidString(value) {
		return false
	}
	lower := strings.ToLower(value)
	badFragments := []string{"пїЅ", "Гђ", "Г‘", "вЂ", "В«", "В»", "\\u", "<", ">", "\ufffd"}
	for _, fragment := range badFragments {
		if strings.Contains(lower, strings.ToLower(fragment)) {
			return false
		}
	}
	for _, r := range value {
		switch r {
		case 'Ђ', 'Ѓ', 'Љ', 'Њ', 'Ќ', 'Ћ', 'Џ', 'ђ', 'ѓ', 'љ', 'њ', 'ќ', 'ћ', 'џ', 'ў', 'ґ':
			return false
		}
	}
	if len([]rune(value)) > 220 {
		return false
	}
	return true
}

func looksLikeEnglishOnly(value string) bool {
	value = cleanDictionaryDisplay(value)
	if value == "" {
		return false
	}
	letters := 0
	latin := 0
	cyrillic := 0
	for _, r := range value {
		if !unicode.IsLetter(r) {
			continue
		}
		letters++
		switch {
		case unicode.In(r, unicode.Latin):
			latin++
		case unicode.In(r, unicode.Cyrillic):
			cyrillic++
		}
	}
	return letters > 0 && latin == letters && cyrillic == 0
}

func tutorStableWordKey(word vocabWord) string {
	if word.FrequencyRank > 0 {
		return fmt.Sprintf("%08d:%s", word.FrequencyRank, word.ID)
	}
	return "99999999:" + word.ID
}

func tutorStableIndex(seed string, max int) int {
	if max <= 0 {
		return 0
	}
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(seed))
	return int(hash.Sum32() % uint32(max))
}

func tutorMinInt(left int, right int) int {
	if left < right {
		return left
	}
	return right
}

func tutorMaxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
