package main

import (
	"strings"
	"testing"
)

func variationSeedLine(t *testing.T, content string) string {
	t.Helper()
	marker := "Variation seed:"
	index := strings.Index(content, marker)
	if index < 0 {
		t.Fatalf("expected prompt to include %q, got %q", marker, content)
	}
	line := content[index:]
	if end := strings.Index(line, "."); end >= 0 {
		line = line[:end]
	}
	return line
}

func TestTrimPracticeHistoryKeepsLastFiveMessages(t *testing.T) {
	history := trimPracticeHistory([]string{"one", "two", "three", "four", "five", "six"}, practiceMemoryLimit)
	if len(history) != 5 {
		t.Fatalf("expected 5 messages, got %d", len(history))
	}
	if history[0] != "two" || history[4] != "six" {
		t.Fatalf("unexpected history: %#v", history)
	}
}

func TestTrimLessonHistoryKeepsLastTenLessons(t *testing.T) {
	history := trimLessonHistory([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"})
	if len(history) != 10 {
		t.Fatalf("expected 10 lessons, got %d", len(history))
	}
	if history[0] != "2" || history[9] != "11" {
		t.Fatalf("unexpected lesson history: %#v", history)
	}
}

func TestPracticePromptIncludesRecentMessages(t *testing.T) {
	messages := practicePrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "A2", []string{"Hello", "I goed home", "What next?"}, "", 4)
	if len(messages) != 2 {
		t.Fatalf("expected 2 chat messages, got %d", len(messages))
	}
	content := messages[1].Content
	for _, want := range []string{"1) Hello", "2) I goed home", "3) What next?", "Latest learner message: What next?"} {
		if !strings.Contains(content, want) {
			t.Fatalf("expected prompt to contain %q, got %q", want, content)
		}
	}
	if !strings.Contains(content, "own final line") {
		t.Fatalf("expected prompt to keep the final question extractable, got %q", content)
	}
}

func TestLessonPromptUsesSelectedLanguage(t *testing.T) {
	messages := lessonPrompt(learningLanguageByCode("it"), interfaceLanguageByCode("pl"), "A1", "", 2, []string{"A previous cafe lesson"})
	if len(messages) != 2 {
		t.Fatalf("expected 2 chat messages, got %d", len(messages))
	}
	if !strings.Contains(messages[0].Content, "Italian coach") {
		t.Fatalf("expected Italian system prompt, got %q", messages[0].Content)
	}
	if !strings.Contains(messages[1].Content, "Italian practice task") {
		t.Fatalf("expected Italian lesson prompt, got %q", messages[1].Content)
	}
	if !strings.Contains(messages[1].Content, "Polish") {
		t.Fatalf("expected prompt to use selected interface language, got %q", messages[1].Content)
	}
	if !strings.Contains(messages[1].Content, "Sytuacja:") || !strings.Contains(messages[1].Content, "Twoje zadanie:") {
		t.Fatalf("expected localized lesson section labels, got %q", messages[1].Content)
	}
	if strings.Contains(messages[1].Content, "Situation:, Pattern:, Useful chunks:, Your task:") {
		t.Fatalf("lesson prompt must not hard-code English labels for Polish UI: %q", messages[1].Content)
	}
	if !strings.Contains(messages[1].Content, "never repeat any of the last 10 lessons") || !strings.Contains(messages[1].Content, "A previous cafe lesson") {
		t.Fatalf("expected lesson prompt to include anti-repeat context, got %q", messages[1].Content)
	}
}

func TestLessonPromptIncludesDistinctVariationSeed(t *testing.T) {
	first := lessonPrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "A2", "business calls", 0, []string{"Cafe order with polite requests"})[1].Content
	second := lessonPrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "A2", "business calls", 1, []string{"Cafe order with polite requests"})[1].Content
	firstSeed := variationSeedLine(t, first)
	secondSeed := variationSeedLine(t, second)
	if firstSeed == secondSeed {
		t.Fatalf("expected different lesson variation seeds for different counters, got %q", firstSeed)
	}
	for _, content := range []string{first, second} {
		if !strings.Contains(content, "Do not reveal the variation seed") {
			t.Fatalf("expected lesson prompt to keep seed private, got %q", content)
		}
	}
}

func TestPracticePromptUsesLocalizedSectionLabels(t *testing.T) {
	messages := practicePrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "A2", []string{"I goed home"}, "", 7)
	content := messages[1].Content
	for _, want := range []string{"Исправление:", "Пример фразы:", "Твоя очередь:"} {
		if !strings.Contains(content, want) {
			t.Fatalf("expected practice prompt to contain %q, got %q", want, content)
		}
	}
	if strings.Contains(content, "Correction:, Model phrase:, Your turn:") {
		t.Fatalf("practice prompt must not hard-code English labels for Russian UI: %q", content)
	}
	if !strings.Contains(content, "Practice novelty rule") || !strings.Contains(content, "Practice turn number: 8") {
		t.Fatalf("expected practice prompt to include anti-repeat rotation, got %q", content)
	}
}

func TestPracticePromptIncludesDistinctVariationSeed(t *testing.T) {
	first := practicePrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "B1", []string{"I need explain delay"}, "work updates", 2)[1].Content
	second := practicePrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "B1", []string{"I need explain delay"}, "work updates", 3)[1].Content
	firstSeed := variationSeedLine(t, first)
	secondSeed := variationSeedLine(t, second)
	if firstSeed == secondSeed {
		t.Fatalf("expected different practice variation seeds for different counters, got %q", firstSeed)
	}
	for _, content := range []string{first, second} {
		if !strings.Contains(content, "Do not reveal the variation seed") {
			t.Fatalf("expected practice prompt to keep seed private, got %q", content)
		}
	}
}

func TestRoleplayPromptIncludesPrivateVariationSeed(t *testing.T) {
	content := roleplayPrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "Scenario: short cafe roleplay", "small talk")[1].Content
	_ = variationSeedLine(t, content)
	if !strings.Contains(content, "Do not reveal the variation seed") {
		t.Fatalf("expected roleplay prompt to keep seed private, got %q", content)
	}
}

func TestToolTranslationPromptUsesInterfaceLanguage(t *testing.T) {
	messages := translationPrompt("Hello", interfaceLanguageByCode("ka"))
	if !strings.Contains(messages[1].Content, "Georgian") {
		t.Fatalf("expected translation prompt to target interface language, got %q", messages[1].Content)
	}
	if strings.Contains(messages[1].Content, "natural Russian") {
		t.Fatalf("translation prompt must not be hard-coded to Russian: %q", messages[1].Content)
	}
}

func TestVocabularyExampleTranslationPromptUsesInterfaceLanguage(t *testing.T) {
	messages := vocabularyExampleTranslationPrompt(vocabWord{English: "size"}, "What size shoes do you wear?", learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "learn word")
	if len(messages) != 2 {
		t.Fatalf("expected 2 chat messages, got %d", len(messages))
	}
	content := messages[1].Content
	for _, want := range []string{"Russian", "What size shoes do you wear?"} {
		if !strings.Contains(content, want) {
			t.Fatalf("expected translation prompt to contain %q, got %q", want, content)
		}
	}
	if !strings.Contains(strings.ToLower(content), "translate") {
		t.Fatalf("expected translation prompt to ask for translation, got %q", content)
	}
	if strings.Contains(content, "write one short natural example") {
		t.Fatalf("translation prompt must not ask to generate a new example: %q", content)
	}
}

func TestCleanModelReplyRemovesPromptLeak(t *testing.T) {
	raw := "You are an encouraging English coach for learners.\nAct like an expert English teacher.\n\nГотово: перевод ниже."
	got := cleanModelReply(raw)
	if strings.Contains(got, "You are an encouraging") || strings.Contains(got, "Act like an expert") {
		t.Fatalf("expected prompt leak to be removed, got %q", got)
	}
	if !strings.Contains(got, "Готово") {
		t.Fatalf("expected real reply to remain, got %q", got)
	}
}
