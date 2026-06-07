package main

import (
	"strings"
	"testing"
)

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

func TestToolTranslationPromptUsesInterfaceLanguage(t *testing.T) {
	messages := translationPrompt("Hello", interfaceLanguageByCode("ka"))
	if !strings.Contains(messages[1].Content, "Georgian") {
		t.Fatalf("expected translation prompt to target interface language, got %q", messages[1].Content)
	}
	if strings.Contains(messages[1].Content, "natural Russian") {
		t.Fatalf("translation prompt must not be hard-coded to Russian: %q", messages[1].Content)
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
