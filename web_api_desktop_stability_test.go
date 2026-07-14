package main

import (
	"strings"
	"testing"
)

func TestPracticePromptForWebInputRoutesRoleplayEnvelope(t *testing.T) {
	language := learningLanguageByCode("en")
	interfaceLanguage := learningLanguageByCode("ru")
	messages := []string{"ROLEPLAY_TOOL_V2\nScenario: job interview\nLearner line: I am applying for the manager position."}

	prompt := practicePromptForWebInput(language, interfaceLanguage, "A2", messages, "", 0)
	if len(prompt) < 2 {
		t.Fatalf("prompt has %d messages, want at least 2", len(prompt))
	}
	if !strings.Contains(prompt[0].Content, "dedicated AI roleplay tool") {
		t.Fatalf("roleplay envelope must use roleplay prompt, got system prompt: %q", prompt[0].Content)
	}
	if strings.Contains(prompt[1].Content, "Continue a friendly") {
		t.Fatalf("roleplay envelope must not use the generic practice-chat prompt: %q", prompt[1].Content)
	}
}

func TestLessonAnswerMatchesOfferedPhraseIgnoringCosmeticDifferences(t *testing.T) {
	task := "Example phrase: That sounds great! I would love to go. What time should we meet?"
	answer := "that sounds great, I would love to go, what time should we meet"

	if !lessonAnswerMatchesOfferedPhrase(task, answer) {
		t.Fatalf("an offered model phrase with case and terminal punctuation changes must match")
	}
	if lessonAnswerMatchesOfferedPhrase(task, "I would like to go tomorrow") {
		t.Fatalf("an unrelated answer must not match the lesson task")
	}
}
