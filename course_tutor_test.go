package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildTutorLessonUsesLocalA1A2CourseCore(t *testing.T) {
	user := userState{
		TelegramID:        42,
		FirstName:         "demo",
		InterfaceLanguage: "ru",
		LearningLanguage:  "en",
		Level:             "A1",
	}
	lesson, err := buildTutorLesson(user)
	if err != nil {
		t.Fatalf("buildTutorLesson() error = %v", err)
	}
	if lesson.Source != "local-a1-a2-course-bank" {
		t.Fatalf("lesson source = %q, want local course bank", lesson.Source)
	}
	if lesson.VariantCode == "" || lesson.VariantTitle == "" || lesson.FocusSkill == "" || len(lesson.SuccessCriteria) < 3 {
		t.Fatalf("lesson variant metadata is incomplete: variant=%q title=%q focus=%q criteria=%#v", lesson.VariantCode, lesson.VariantTitle, lesson.FocusSkill, lesson.SuccessCriteria)
	}
	if lesson.Level != "A1" {
		t.Fatalf("lesson level = %q, want A1", lesson.Level)
	}
	if len(lesson.Words) < 4 {
		t.Fatalf("lesson words = %d, want at least 4", len(lesson.Words))
	}
	if lesson.Grammar == "" || lesson.WritingTask == "" || lesson.ListeningText == "" || lesson.PronunciationText == "" || len(lesson.Dialogue) == 0 {
		t.Fatalf("lesson is missing required sections: %#v", lesson)
	}
	if lesson.PronunciationText == lesson.ListeningText {
		t.Fatalf("pronunciation text should be a standalone target, got same text as listening: %q", lesson.PronunciationText)
	}
	if len(lesson.Steps) < 7 {
		t.Fatalf("lesson steps = %d, want full guided flow", len(lesson.Steps))
	}
	if lesson.CourseSize < 300 || lesson.LessonNumber <= 0 {
		t.Fatalf("lesson course markers are incomplete: number=%d size=%d", lesson.LessonNumber, lesson.CourseSize)
	}
	if lesson.CanDo == "" || lesson.Scenario == "" {
		t.Fatalf("lesson is missing scenario/can-do: can_do=%q scenario=%q", lesson.CanDo, lesson.Scenario)
	}
	if lesson.Choice.CorrectAnswerID == "" || len(lesson.Choice.Options) < 4 {
		t.Fatalf("choice task is incomplete: %#v", lesson.Choice)
	}
	for _, option := range lesson.Choice.Options {
		if strings.TrimSpace(option.Label) == "" || strings.TrimSpace(option.Why) == "" {
			t.Fatalf("choice option lacks teaching metadata: %#v", option)
		}
	}
	if len(lesson.Checks) < 3 {
		t.Fatalf("checks = %d, want multi-question scenario check", len(lesson.Checks))
	}
	for _, check := range lesson.Checks {
		for _, option := range check.Options {
			if strings.TrimSpace(option.Label) == "" || strings.TrimSpace(option.Why) == "" {
				t.Fatalf("scenario option lacks teaching metadata: %#v", option)
			}
		}
	}
	if len(lesson.AnswerVariants) < 4 || len(lesson.DialogueVariants) < 4 {
		t.Fatalf("answer variants are too thin: writing=%#v dialogue=%#v", lesson.AnswerVariants, lesson.DialogueVariants)
	}
	if !hasAvoidVariant(lesson.AnswerVariants) || !hasAvoidVariant(lesson.DialogueVariants) {
		t.Fatalf("answer variants must include a common weak answer to avoid: writing=%#v dialogue=%#v", lesson.AnswerVariants, lesson.DialogueVariants)
	}
	if !hasStrongTutorVariant(lesson.AnswerVariants) || !hasStrongTutorVariant(lesson.DialogueVariants) {
		t.Fatalf("answer variants do not contain a useful model answer: writing=%#v dialogue=%#v", lesson.AnswerVariants, lesson.DialogueVariants)
	}
	if strings.Count(strings.ToLower(lesson.ListeningText), "i need") > 1 {
		t.Fatalf("listening text repeats weak I need pattern: %q", lesson.ListeningText)
	}
	if strings.Contains(strings.Join(lesson.Dialogue, "\n"), "Мне нужно:") {
		t.Fatalf("dialogue still looks like word examples, not a dialogue: %#v", lesson.Dialogue)
	}
	if len(lesson.ReviewSummary) < 3 {
		t.Fatalf("review summary is too thin: %#v", lesson.ReviewSummary)
	}
	for _, text := range []string{lesson.Title, lesson.Topic, lesson.Goal, lesson.CanDo, lesson.Scenario, lesson.GrammarTitle, lesson.Grammar, lesson.MiniExplanation, lesson.WritingTask, lesson.ListeningTask, lesson.ListeningQuestion, lesson.PronunciationText, lesson.DialoguePrompt} {
		if text == "" || !tutorLooksClean(text) {
			t.Fatalf("visible tutor text is not clean: %q", text)
		}
	}
	for _, variant := range append(append([]tutorAnswerVariant{}, lesson.AnswerVariants...), lesson.DialogueVariants...) {
		if strings.TrimSpace(variant.Text) == "" || strings.TrimSpace(variant.Label) == "" || strings.TrimSpace(variant.Why) == "" {
			t.Fatalf("answer variant is incomplete: %#v", variant)
		}
		if !tutorLooksClean(variant.Text) || !tutorLooksClean(variant.Label) || !tutorLooksClean(variant.Why) {
			t.Fatalf("answer variant is not clean: %#v", variant)
		}
	}
	for _, word := range lesson.Words {
		if word.Word == "" || word.Translation == "" {
			t.Fatalf("word has empty prompt/translation: %#v", word)
		}
		if word.Level != "A1" && word.Level != "A2" {
			t.Fatalf("word %q has level %q, want A1/A2", word.Word, word.Level)
		}
		if !tutorLooksClean(word.Word) || !tutorLooksClean(word.Translation) {
			t.Fatalf("word is not clean enough for tutor: %#v", word)
		}
		if looksLikeEnglishOnly(word.Translation) {
			t.Fatalf("word leaks English fallback in ru interface: %#v", word)
		}
	}
}

func TestTutorCafeScenarioUsesSlotsNotRandomWords(t *testing.T) {
	topic := tutorTestTopicByCode(t, "food")
	function := tutorCourseFunctions[0]
	words := []tutorLessonWord{
		{ID: "en:menu", Word: "menu", Translation: "меню", Level: "A1"},
		{ID: "en:football", Word: "football", Translation: "футбол", Level: "A1"},
		{ID: "en:coffee", Word: "coffee", Translation: "кофе", Level: "A1"},
		{ID: "en:breakfast", Word: "breakfast", Translation: "завтрак", Level: "A1"},
	}

	checks := tutorScenarioChecks(words, topic, function, "ru")
	variants := tutorAnswerVariants(words, topic, function, "en", "ru", false)
	dialogueVariants := tutorAnswerVariants(words, topic, function, "en", "ru", true)

	combined := strings.ToLower(strings.Join(append(tutorChoiceTexts(checks), tutorVariantTexts(append(variants, dialogueVariants...))...), "\n"))
	for _, want := range []string{"coffee", "for breakfast", "could i see the menu"} {
		if !strings.Contains(combined, want) {
			t.Fatalf("cafe scenario content misses %q:\n%s", want, combined)
		}
	}
	for _, bad := range []string{"have football", "like football", "football for breakfast"} {
		if strings.Contains(combined, bad) {
			t.Fatalf("cafe scenario used football as an order/detail via %q:\n%s", bad, combined)
		}
	}
}

func TestTutorRussianMiniExplanationDoesNotLeakServiceCriteria(t *testing.T) {
	user := userState{TelegramID: 91, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}
	lesson, err := buildTutorLessonForSequence(user, 1)
	if err != nil {
		t.Fatalf("buildTutorLessonForSequence(food) error = %v", err)
	}
	visible := strings.Join(append([]string{lesson.MiniExplanation}, lesson.SuccessCriteria...), "\n")
	for _, bad := range []string{"Use the lesson goal", "Include ", "Variant focus"} {
		if strings.Contains(visible, bad) {
			t.Fatalf("visible RU tutor text leaked service criteria %q:\n%s", bad, visible)
		}
	}
	if !strings.Contains(visible, "попросить меню") || !strings.Contains(visible, "coffee") || !strings.Contains(visible, "for breakfast") {
		t.Fatalf("RU cafe explanation does not expose the scenario pattern:\n%s", visible)
	}
}

func TestTutorScenarioChoicesAreDialogueReplies(t *testing.T) {
	topic := tutorTestTopicByCode(t, "food")
	function := tutorCourseFunctions[0]
	words := []tutorLessonWord{
		{ID: "en:menu", Word: "menu", Translation: "меню", Level: "A1"},
		{ID: "en:football", Word: "football", Translation: "футбол", Level: "A1"},
		{ID: "en:coffee", Word: "coffee", Translation: "кофе", Level: "A1"},
		{ID: "en:breakfast", Word: "breakfast", Translation: "завтрак", Level: "A1"},
	}

	checks := tutorScenarioChecks(words, topic, function, "ru")
	if len(checks) < 2 {
		t.Fatalf("checks = %d, want at least 2 scenario checks", len(checks))
	}
	best := tutorCorrectChoiceText(t, checks[0])
	detail := tutorCorrectChoiceText(t, checks[1])
	if best != "Could I see the menu and have coffee for breakfast, please?" {
		t.Fatalf("best cafe reply = %q", best)
	}
	if detail != "for breakfast" {
		t.Fatalf("detail check = %q, want for breakfast", detail)
	}
	if strings.Contains(best, " / ") || strings.Contains(detail, ":") {
		t.Fatalf("correct choices should be dialogue replies/details, got best=%q detail=%q", best, detail)
	}
}

func TestTutorWorkScenarioUsesMeetingSlotsNotRandomVocabulary(t *testing.T) {
	topic := tutorTestTopicByCode(t, "work")
	function := tutorCourseFunctions[0]
	words := []tutorLessonWord{
		{ID: "en:tax", Word: "tax", Translation: "налог", Level: "A2"},
		{ID: "en:whatever", Word: "whatever", Translation: "что угодно", Level: "A2"},
		{ID: "en:meeting", Word: "meeting", Translation: "встреча", Level: "A1"},
		{ID: "en:contact", Word: "contact", Translation: "контакт", Level: "A2"},
	}

	mini := tutorMiniExplanation("A2", topic, function, words, "ru")
	checks := tutorScenarioChecks(words, topic, function, "ru")
	variants := tutorAnswerVariants(words, topic, function, "en", "ru", false)
	dialogueVariants := tutorAnswerVariants(words, topic, function, "en", "ru", true)

	combined := strings.ToLower(strings.Join(append([]string{mini}, append(tutorChoiceTexts(checks), tutorVariantTexts(append(variants, dialogueVariants...))...)...), "\n"))
	for _, want := range []string{"meeting", "alex", "10", "contact"} {
		if !strings.Contains(combined, want) {
			t.Fatalf("work tutor scenario misses meeting slot %q:\n%s", want, combined)
		}
	}
	for _, bad := range []string{"tax", "whatever", "i repeat:", "налог /", " / что угодно", "use the lesson goal", "variant focus"} {
		if strings.Contains(combined, bad) {
			t.Fatalf("work tutor scenario leaked random vocabulary/service text %q:\n%s", bad, combined)
		}
	}
}

func tutorTestTopicByCode(t *testing.T, code string) tutorTopic {
	t.Helper()
	for _, topic := range tutorScenarioCourseTopics {
		if topic.Code == code {
			return topic
		}
	}
	t.Fatalf("topic %q not found", code)
	return tutorTopic{}
}

func tutorChoiceTexts(checks []tutorLessonChoice) []string {
	var texts []string
	for _, check := range checks {
		texts = append(texts, check.Prompt)
		for _, option := range check.Options {
			texts = append(texts, option.Text, option.Label, option.Why)
		}
	}
	return texts
}

func tutorVariantTexts(variants []tutorAnswerVariant) []string {
	texts := make([]string, 0, len(variants)*3)
	for _, variant := range variants {
		texts = append(texts, variant.Text, variant.Label, variant.Why)
	}
	return texts
}

func tutorCorrectChoiceText(t *testing.T, choice tutorLessonChoice) string {
	t.Helper()
	for _, option := range choice.Options {
		if option.ID == choice.CorrectAnswerID {
			return option.Text
		}
	}
	t.Fatalf("correct option %q not found in %#v", choice.CorrectAnswerID, choice.Options)
	return ""
}

func hasAvoidVariant(variants []tutorAnswerVariant) bool {
	for _, variant := range variants {
		if variant.Avoid {
			return true
		}
	}
	return false
}

func hasStrongTutorVariant(variants []tutorAnswerVariant) bool {
	for _, variant := range variants {
		if variant.Avoid {
			continue
		}
		text := strings.ToLower(variant.Text)
		if strings.Contains(text, "please") || strings.Contains(text, "could") || strings.Count(text, " ") >= 7 {
			return true
		}
	}
	return false
}

func TestTutorLessonAvailableForAllLearningLanguages(t *testing.T) {
	for _, language := range learningLanguages {
		user := userState{
			TelegramID:        100,
			FirstName:         "demo",
			InterfaceLanguage: "ru",
			LearningLanguage:  language.Code,
			Level:             "A2",
		}
		lesson, err := buildTutorLesson(user)
		if err != nil {
			t.Fatalf("buildTutorLesson(%s) error = %v", language.Code, err)
		}
		if lesson.LearningLanguage != normalizeLearningLanguage(language.Code) {
			t.Fatalf("lesson language = %q, want %q", lesson.LearningLanguage, language.Code)
		}
		if len(lesson.Words) < 4 {
			t.Fatalf("buildTutorLesson(%s) returned %d words, want at least 4", language.Code, len(lesson.Words))
		}
	}
}

func TestTutorLessonCourseHas300StableUniqueLessons(t *testing.T) {
	user := userState{
		TelegramID:        77,
		FirstName:         "demo",
		InterfaceLanguage: "en",
		LearningLanguage:  "en",
		Level:             "A1",
	}
	seen := map[string]bool{}
	variantSeen := map[string]bool{}
	for sequence := 0; sequence < tutorCourseSize; sequence++ {
		lesson, err := buildTutorLessonForSequence(user, sequence)
		if err != nil {
			t.Fatalf("buildTutorLessonForSequence(%d) error = %v", sequence, err)
		}
		if seen[lesson.ID] {
			t.Fatalf("duplicate lesson id at sequence %d: %s", sequence, lesson.ID)
		}
		seen[lesson.ID] = true
		variantSeen[lesson.VariantCode] = true
		if lesson.LessonNumber < 1 || lesson.LessonNumber > tutorCourseSize {
			t.Fatalf("lesson number out of course bounds at sequence %d: %d", sequence, lesson.LessonNumber)
		}
		if lesson.VariantCode == "" || lesson.FocusSkill == "" || len(lesson.SuccessCriteria) < 3 {
			t.Fatalf("lesson %s lacks variant teaching metadata: %#v", lesson.ID, lesson)
		}
	}
	if len(seen) != tutorCourseSize {
		t.Fatalf("unique lesson count = %d, want %d", len(seen), tutorCourseSize)
	}
	if len(variantSeen) != len(tutorCourseVariants) {
		t.Fatalf("variant coverage = %d, want %d", len(variantSeen), len(tutorCourseVariants))
	}
	next, err := buildTutorLessonForSequence(user, tutorCourseSize)
	if err != nil {
		t.Fatalf("buildTutorLessonForSequence(%d) error = %v", tutorCourseSize, err)
	}
	if seen[next.ID] {
		t.Fatalf("lesson after course cycle reused an id: %s", next.ID)
	}
}

func TestSQLiteStoreNextTutorLessonReusesBankWithoutRepeatingForUser(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore: %v", err)
	}
	defer store.db.Close()
	userOne := userState{TelegramID: 501, FirstName: "one", InterfaceLanguage: "en", LearningLanguage: "en", Level: "A1"}
	userTwo := userState{TelegramID: 502, FirstName: "two", InterfaceLanguage: "en", LearningLanguage: "en", Level: "A1"}
	factoryOne := func(sequence int) (tutorLesson, error) {
		return buildTutorLessonForSequence(tutorReusableLessonUser(userOne), sequence)
	}
	first, err := store.nextTutorLesson(userOne, factoryOne)
	if err != nil {
		t.Fatalf("first nextTutorLesson: %v", err)
	}
	second, err := store.nextTutorLesson(userOne, factoryOne)
	if err != nil {
		t.Fatalf("second nextTutorLesson: %v", err)
	}
	if first.ID == second.ID {
		t.Fatalf("same user received duplicate tutor lesson: %s", first.ID)
	}
	other, err := store.nextTutorLesson(userTwo, func(sequence int) (tutorLesson, error) {
		return buildTutorLessonForSequence(tutorReusableLessonUser(userTwo), sequence)
	})
	if err != nil {
		t.Fatalf("other user nextTutorLesson: %v", err)
	}
	if other.ID != first.ID {
		t.Fatalf("second user did not reuse the first banked lesson: got %s want %s", other.ID, first.ID)
	}
}
