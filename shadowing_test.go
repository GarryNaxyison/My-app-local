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
	messages := shadowingPhrasePrompt(learningLanguageByCode("en"), interfaceLanguageByCode("ru"), "A2", shadowingDeckSpec(42))
	if len(messages) != 2 {
		t.Fatalf("expected 2 chat messages, got %d", len(messages))
	}
	if !containsAll(messages[1].Content, []string{"Create exactly one natural English phrase", "10,000-card phrase deck", "Return only the target-language phrase"}) {
		t.Fatalf("unexpected shadowing prompt: %q", messages[1].Content)
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
