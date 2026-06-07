package main

import (
	"strings"
	"testing"
)

func TestPronunciationTechnicalReportExactRepeatFlagsMissingWords(t *testing.T) {
	user := userState{InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A2"}
	transcription := audioTranscription{
		Text:            "Could you tell me where station is",
		DurationSeconds: 3.2,
		Tokens: []transcriptionToken{
			{Token: "Could", LogProb: -0.05},
			{Token: "you", LogProb: -0.05},
			{Token: "tell", LogProb: -0.05},
			{Token: "me", LogProb: -0.05},
			{Token: "where", LogProb: -0.05},
			{Token: "station", LogProb: -0.35},
			{Token: "is", LogProb: -0.05},
		},
	}
	report := buildPronunciationTechnicalReport(user, "Could you tell me where the nearest station is?", transcription, pronunciationModeExact)
	if report.Similarity >= 100 {
		t.Fatalf("expected imperfect similarity, got %d", report.Similarity)
	}
	if len(report.MissingWords) == 0 {
		t.Fatalf("expected missing words in report")
	}
	if len(report.ProblemWords) == 0 {
		t.Fatalf("expected problem words in report")
	}
	assessment := localPronunciationAssessment(user, report, 80)
	if assessment.Score >= 100 {
		t.Fatalf("expected score below perfect, got %d", assessment.Score)
	}
}

func TestPronunciationTechnicalReportFreeSpeechUsesConfidenceAndFluency(t *testing.T) {
	user := userState{InterfaceLanguage: "en", LearningLanguage: "en", Level: "B1"}
	transcription := audioTranscription{
		Text:            "I usually read in the evening",
		DurationSeconds: 2.8,
		Tokens: []transcriptionToken{
			{Token: "I", LogProb: -0.02},
			{Token: "usually", LogProb: -1.2},
			{Token: "read", LogProb: -0.05},
			{Token: "in", LogProb: -0.05},
			{Token: "the", LogProb: -0.05},
			{Token: "evening", LogProb: -0.05},
		},
	}
	report := buildPronunciationTechnicalReport(user, "", transcription, pronunciationModeFree)
	if report.Mode != string(pronunciationModeFree) {
		t.Fatalf("expected free speech mode, got %q", report.Mode)
	}
	if report.Fluency <= 0 || report.AverageConfidence <= 0 {
		t.Fatalf("expected positive fluency and confidence, got fluency=%d confidence=%f", report.Fluency, report.AverageConfidence)
	}
	if len(report.ProblemWords) == 0 || report.ProblemWords[0].Word != "usually" {
		t.Fatalf("expected low-confidence word usually, got %#v", report.ProblemWords)
	}
}

func TestParsePronunciationCoachJSONStripsCodeFence(t *testing.T) {
	raw := "```json\n{\"model_correction\":84,\"accent_strength\":22,\"fluency\":91,\"problem_words\":[{\"word\":\"nearest\",\"confidence\":0.61}],\"tips\":[\"Open your vowels\"],\"overall_feedback\":\"Good repeat.\"}\n```"
	parsed, ok := parsePronunciationCoachJSON(raw)
	if !ok {
		t.Fatalf("expected coach JSON to parse")
	}
	if parsed.ModelCorrection != 84 || parsed.AccentStrength != 22 || len(parsed.ProblemWords) != 1 {
		t.Fatalf("unexpected parsed coach JSON: %#v", parsed)
	}
}

func TestPronunciationEstimatedConfidenceDoesNotAllowPerfectRepeat(t *testing.T) {
	user := userState{InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A2"}
	transcription := audioTranscription{
		Text:            "Could you tell me where the nearest station is",
		DurationSeconds: 3.1,
	}
	report := buildPronunciationTechnicalReport(user, "Could you tell me where the nearest station is?", transcription, pronunciationModeExact)
	if report.ConfidenceSource != "estimated" {
		t.Fatalf("expected estimated confidence without logprobs, got %q", report.ConfidenceSource)
	}
	assessment := localPronunciationAssessment(user, report, 100)
	if assessment.Score > 76 {
		t.Fatalf("expected estimated confidence exact repeat to be capped, got %d", assessment.Score)
	}
	if !strings.Contains(assessment.Feedback, "Балл ограничен") {
		t.Fatalf("expected capped-confidence feedback, got %q", assessment.Feedback)
	}
}
