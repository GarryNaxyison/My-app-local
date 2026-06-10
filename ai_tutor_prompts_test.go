package main

import (
	"strings"
	"testing"
)

func TestAITutorGenerationPromptContract(t *testing.T) {
	language := learningLanguageByCode("en")
	interfaceLanguage := learningLanguageByCode("ru")
	messages := aiTutorLessonGenerationPrompt(language, interfaceLanguage, "A1", "food", []string{"old-fp"})
	if len(messages) != 2 {
		t.Fatalf("message count = %d", len(messages))
	}
	content := messages[1].Content
	for _, want := range []string{
		"Return strict JSON only",
		"Exactly 5-7 sentences",
		"Select exactly 6 useful target words",
		"Create exactly 3 comprehension questions",
		"old-fp",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("generation prompt misses %q:\n%s", want, content)
		}
	}
}

func TestAITutorCheckerPromptsReturnJSONOnly(t *testing.T) {
	language := learningLanguageByCode("en")
	interfaceLanguage := learningLanguageByCode("ru")
	lesson := validAITutorLessonPayloadForTest()
	prompts := [][]chatMessage{
		aiTutorRetellCheckPrompt(language, interfaceLanguage, "A1", lesson, "Mia goes shop."),
		aiTutorQuestionCheckPrompt(language, interfaceLanguage, "A1", lesson, lesson.ComprehensionQuestions[0], "Early."),
		aiTutorProductionCheckPrompt(language, interfaceLanguage, "A1", lesson, "I drink water. I buy bread."),
		aiTutorQualityPrompt(language, interfaceLanguage, lesson, nil, "preflight"),
	}
	for index, prompt := range prompts {
		if len(prompt) != 2 {
			t.Fatalf("prompt %d message count = %d", index, len(prompt))
		}
		if !strings.Contains(prompt[1].Content, "Return strict JSON only") {
			t.Fatalf("prompt %d does not force JSON:\n%s", index, prompt[1].Content)
		}
	}
}
