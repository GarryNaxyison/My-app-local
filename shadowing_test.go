package main

import (
	"strings"
	"testing"
)

func TestShadowingModeRoundTrip(t *testing.T) {
	target := "Could you say that a little slower, please?"
	mode := shadowingMode(target)
	got, ok := parseShadowingMode(mode)
	if !ok {
		t.Fatalf("expected mode to parse")
	}
	if got != target {
		t.Fatalf("target = %q, want %q", got, target)
	}
}

func TestPronunciationPracticeModeRoundTrip(t *testing.T) {
	target := "Could you say that clearly?"
	mode := pronunciationPracticeMode(target)
	got, ok := parsePronunciationMode(mode)
	if !ok {
		t.Fatalf("expected mode to parse")
	}
	if got != target {
		t.Fatalf("target = %q, want %q", got, target)
	}
	if _, ok := parseShadowingMode(mode); ok {
		t.Fatalf("pronunciation mode must not parse as shadowing: %q", mode)
	}
}

func TestShadowingScoreRewardsCloseRepeat(t *testing.T) {
	closeScore := shadowingScore("Could you say that a little slower please", "could you say that a little slower please")
	if closeScore < 95 {
		t.Fatalf("close repeat score = %d, want >= 95", closeScore)
	}
	weakScore := shadowingScore("Could you say that a little slower please", "hello tomorrow")
	if weakScore >= closeScore {
		t.Fatalf("weak score = %d must be below close score %d", weakScore, closeScore)
	}
}

func TestShadowingPhrasePromptIsTargetOnly(t *testing.T) {
	messages := shadowingPhrasePrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "A2", shadowingDeckSpec(42), "")
	if len(messages) != 2 {
		t.Fatalf("expected 2 chat messages, got %d", len(messages))
	}
	if !containsAll(messages[1].Content, []string{"Create exactly one natural English phrase", "10,000-card phrase deck", "Return only the target-language phrase"}) {
		t.Fatalf("unexpected shadowing prompt: %q", messages[1].Content)
	}
}

func TestShadowingPhrasePromptAvoidsPreviousPhrase(t *testing.T) {
	previous := "Could you say that a little slower, please?"
	messages := shadowingPhrasePrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "A2", shadowingDeckSpec(42), previous)
	if !strings.Contains(messages[1].Content, "Do not repeat or lightly paraphrase this previous phrase: "+previous) {
		t.Fatalf("prompt does not avoid previous phrase: %q", messages[1].Content)
	}
}

func TestBuildShadowingPhraseFallbackAvoidsPreviousPhrase(t *testing.T) {
	user := userState{TelegramID: 0, Level: "A2", LearningLanguage: "en"}
	previous := shadowingFallbackPhrase(learningLanguageByCode("en"), shadowingDeckIndex(user))
	got := (&bot{}).buildShadowingPhrase(t.Context(), user, previous)
	if got == "" {
		t.Fatal("expected fallback phrase")
	}
	if strings.EqualFold(got, previous) {
		t.Fatalf("fallback repeated previous phrase %q", previous)
	}
}

func TestTelegramListeningStartMessageHidesTargetPhrase(t *testing.T) {
	user := userState{InterfaceLanguage: "ru"}
	target := "So, just to clarify, did I understand that correctly?"

	got := shadowingStartMessage(user, target)

	if strings.Contains(got, target) {
		t.Fatalf("listening start message revealed target phrase:\n%s", got)
	}
	for _, unexpected := range []string{"\u0424\u0440\u0430\u0437\u0430:", "\u043e\u0442\u043f\u0440\u0430\u0432\u044c \u0433\u043e\u043b\u043e\u0441"} {
		if strings.Contains(got, unexpected) {
			t.Fatalf("listening start message contains %q:\n%s", unexpected, got)
		}
	}
	if !strings.Contains(got, "\u043d\u0430\u043f\u0438\u0448\u0438") && !strings.Contains(got, "\u041d\u0430\u043f\u0438\u0448\u0438") {
		t.Fatalf("listening start message should ask for typed text:\n%s", got)
	}
}

func TestShadowingUICopiesCoverInterfaceLanguagesAndStayUTF8(t *testing.T) {
	for _, language := range interfaceLanguages() {
		copy, ok := shadowingUICopies[language.Code]
		if !ok {
			t.Fatalf("shadowing copy missing for %s", language.Code)
		}
		fields := map[string]string{
			"Unavailable": copy.Unavailable,
			"Title":       copy.Title,
			"Start":       copy.Start,
			"Phrase":      copy.Phrase,
			"TextHint":    copy.TextHint,
			"NoGaps":      copy.NoGaps,
			"MatchGood":   copy.MatchGood,
			"MatchWeak":   copy.MatchWeak,
			"TipVoice":    copy.TipVoice,
			"TipText":     copy.TipText,
			"Repeat":      copy.Repeat,
			"Target":      copy.Target,
			"Heard":       copy.Heard,
			"Accuracy":    copy.Accuracy,
			"Bonus":       copy.Bonus,
			"Next":        copy.Next,
		}
		for field, value := range fields {
			if strings.TrimSpace(value) == "" {
				t.Fatalf("shadowing copy %s.%s is empty", language.Code, field)
			}
			if strings.ContainsAny(value, "?�") || strings.Contains(value, "Ã") || strings.Contains(value, "Â") || strings.Contains(value, "â€") {
				t.Fatalf("shadowing copy %s.%s looks misencoded: %q", language.Code, field, value)
			}
		}
	}
}

func TestShadowingFallbackPhrasesCoverLearningLanguagesAndStayUTF8(t *testing.T) {
	for _, language := range learningLanguages {
		phrases := shadowingFallbackPhrases[language.Code]
		if len(phrases) != 3 {
			t.Fatalf("shadowing fallback phrases for %s = %d, want 3", language.Code, len(phrases))
		}
		for index, phrase := range phrases {
			if strings.TrimSpace(phrase) == "" {
				t.Fatalf("shadowing fallback phrase %s[%d] is empty", language.Code, index)
			}
			if strings.ContainsAny(phrase, "?�") || strings.Contains(phrase, "Ã") || strings.Contains(phrase, "Â") || strings.Contains(phrase, "â€") {
				t.Fatalf("shadowing fallback phrase %s[%d] looks misencoded: %q", language.Code, index, phrase)
			}
		}
	}
}

func TestRussianShadowingResultUsesReadableLabels(t *testing.T) {
	user := userState{InterfaceLanguage: "ru"}
	got := shadowingResultMessage(user, "Could I have a coffee with milk, please?", "Закончил тут", 1, shadowingFallbackFeedback(user, "Could I have a coffee with milk, please?", "Закончил тут", 1, "", false), 2, nil)
	for _, want := range []string{"Фраза:", "В ответе:", "Совпадение: 1/100", "Опыт: +2 XP", "К сожалению"} {
		if !strings.Contains(got, want) {
			t.Fatalf("Russian listening result misses %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "??") || strings.Contains(got, "�") {
		t.Fatalf("Russian listening result contains replacement question marks:\n%s", got)
	}
}

func TestShadowingDeckAndTextFallback(t *testing.T) {
	if shadowingDeckSize != len(shadowingContexts)*len(shadowingMoves)*len(shadowingDetails) {
		t.Fatalf("shadowing deck size = %d, want generated combination count", shadowingDeckSize)
	}
	voiceScore := shadowingScoreForSource("Could you say that slower", "Could you say that slower", true)
	textScore := shadowingScoreForSource("Could you say that slower", "Could you say that slower", false)
	if voiceScore < 95 {
		t.Fatalf("voice score = %d, want close to perfect", voiceScore)
	}
	if textScore > 60 {
		t.Fatalf("text fallback score = %d, want capped at 60", textScore)
	}
}

func containsAll(text string, parts []string) bool {
	for _, part := range parts {
		if !strings.Contains(text, part) {
			return false
		}
	}
	return true
}
