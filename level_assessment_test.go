package main

import (
	"math"
	"strings"
	"testing"
)

func TestLevelAssessmentScoringUsesQuestionWeights(t *testing.T) {
	user := userState{InterfaceLanguage: "en", LearningLanguage: "en"}
	questions := levelAssessmentQuestionsForUser(user)
	if len(questions) == 0 {
		t.Fatal("expected English level assessment questions")
	}

	score := 0
	first := questions[0]
	if first.Points <= 0 {
		t.Fatalf("first question points = %d, want positive", first.Points)
	}
	if first.CorrectIndex == 0 {
		score += first.Points
	}
	if score != 0 {
		t.Fatalf("wrong answer changed score: got %d", score)
	}
	score += first.Points
	if score != first.Points {
		t.Fatalf("correct answer score = %d, want %d", score, first.Points)
	}
}

func TestLevelAssessmentThresholdsUseMaxScoreRatio(t *testing.T) {
	maxScore := levelAssessmentMaxScoreForUser(userState{InterfaceLanguage: "en", LearningLanguage: "en"})
	if maxScore <= 0 {
		t.Fatal("expected positive max score")
	}
	cases := map[int]string{
		0:                                        "A1",
		int(math.Ceil(float64(maxScore) * 0.15)): "A2",
		int(math.Ceil(float64(maxScore) * 0.33)): "B1",
		int(math.Ceil(float64(maxScore) * 0.53)): "B2",
		int(math.Ceil(float64(maxScore) * 0.71)): "C1",
		int(math.Ceil(float64(maxScore) * 0.86)): "C2",
	}
	for score, want := range cases {
		if got := levelFromAssessmentScore(score, maxScore); got != want {
			t.Fatalf("levelFromAssessmentScore(%d, %d) = %s, want %s", score, maxScore, got, want)
		}
	}
}

func TestLevelAssessmentQuestionsExistForEveryLearningLanguage(t *testing.T) {
	for _, language := range learningLanguages {
		if _, ok := levelAssessmentQuestionsByLanguage[normalizeLearningLanguage(language.Code)]; !ok {
			t.Fatalf("level assessment for %s falls back to vocabulary guessing instead of a static exam", language.Code)
		}
		questions := levelAssessmentQuestionsForUser(userState{InterfaceLanguage: "ru", LearningLanguage: language.Code})
		if len(questions) != 36 {
			t.Fatalf("level assessment questions for %s = %d, want 36", language.Code, len(questions))
		}
		for index, question := range questions {
			if strings.HasPrefix(question.Question, userLearningLanguage(userState{LearningLanguage: language.Code}).InterfaceName+": ") {
				t.Fatalf("level assessment for %s at %d fell back to bare vocabulary guessing: %q", language.Code, index, question.Question)
			}
			if len(question.Options) != 4 {
				t.Fatalf("level assessment for %s at %d has %d options, want 4", language.Code, index, len(question.Options))
			}
		}
	}
}

func TestLevelAssessmentRussianPromptKeepsTaskType(t *testing.T) {
	user := userState{InterfaceLanguage: "ru", LearningLanguage: "en"}
	questions := levelAssessmentQuestionsForUser(user)
	if len(questions) < 4 {
		t.Fatal("expected English level assessment questions")
	}
	if got := questions[1].Question; !strings.HasPrefix(got, "Выбери правильный перевод:") {
		t.Fatalf("translation question lost task type: %q", got)
	}
	if got := questions[2].Question; !strings.HasPrefix(got, "Заполни пропуск:") || !strings.Contains(got, "She _____ to London last year.") {
		t.Fatalf("gap question lost task type or example: %q", got)
	}
	if got := questions[0].Question; got == strings.TrimSpace(ui(user).ChooseAnswer) {
		t.Fatalf("question collapsed to generic choose-answer prompt: %q", got)
	}
}

func TestLevelAssessmentTaskTypesLocalizedForEveryInterfaceLanguage(t *testing.T) {
	kinds := []string{"fill_gap", "correct_sentence", "correct_translation", "word_meaning", "natural_option", "how_to_say", "correct_variant"}
	for _, language := range interfaceLanguages() {
		t.Run(language.Code, func(t *testing.T) {
			generic := strings.TrimSpace(ui(userState{InterfaceLanguage: language.Code}).ChooseAnswer)
			for _, kind := range kinds {
				label := localizedAssessmentTaskLabel(kind, language.Code)
				if strings.TrimSpace(label) == "" {
					t.Fatalf("empty localized assessment label for %s/%s", language.Code, kind)
				}
				if generic != "" && label == generic && kind != "choose_answer" {
					t.Fatalf("assessment label for %s/%s collapsed to generic choose-answer label %q", language.Code, kind, label)
				}
			}
		})
	}
}

func TestLevelAssessmentDeterminesLevelsForNewLearningLanguages(t *testing.T) {
	newLearningLanguages := []string{
		"ar", "bn", "cs", "el", "hi",
		"hu", "id", "nl", "sv", "ta",
		"te", "th", "tl", "tr", "vi",
	}
	cases := []struct {
		name  string
		ratio float64
		want  string
	}{
		{name: "all wrong", ratio: 0, want: "A1"},
		{name: "A2 threshold", ratio: 0.15, want: "A2"},
		{name: "B1 threshold", ratio: 0.33, want: "B1"},
		{name: "B2 threshold", ratio: 0.53, want: "B2"},
		{name: "C1 threshold", ratio: 0.71, want: "C1"},
		{name: "C2 threshold", ratio: 0.86, want: "C2"},
		{name: "all correct", ratio: 1, want: "C2"},
	}

	for _, code := range newLearningLanguages {
		t.Run(code, func(t *testing.T) {
			user := userState{InterfaceLanguage: "ru", LearningLanguage: code}
			maxScore := levelAssessmentMaxScoreForUser(user)
			if maxScore <= 0 {
				t.Fatalf("levelAssessmentMaxScoreForUser(%s) = %d, want positive", code, maxScore)
			}
			for _, tc := range cases {
				t.Run(tc.name, func(t *testing.T) {
					score := int(math.Ceil(float64(maxScore) * tc.ratio))
					if got := levelFromAssessmentScore(score, maxScore); got != tc.want {
						t.Fatalf("levelFromAssessmentScore(%d, %d) for %s = %s, want %s", score, maxScore, code, got, tc.want)
					}
				})
			}
		})
	}
}

func TestNewLearningLanguagesUseInterfaceQuestionAndTargetOptions(t *testing.T) {
	newLearningLanguages := []string{
		"ar", "bn", "cs", "el", "hi",
		"hu", "id", "nl", "sv", "ta",
		"te", "th", "tl", "tr", "vi",
	}
	for _, code := range newLearningLanguages {
		t.Run(code, func(t *testing.T) {
			for _, interfaceLanguage := range []string{"ru", code} {
				t.Run(interfaceLanguage, func(t *testing.T) {
					user := userState{InterfaceLanguage: interfaceLanguage, LearningLanguage: code}
					questions := levelAssessmentQuestionsForUser(user)
					if len(questions) == 0 {
						t.Fatal("expected localized level assessment questions")
					}
					expectedPrompt := strings.TrimSpace(ui(user).ChooseAnswer)
					if expectedPrompt == "" {
						t.Fatal("expected localized ChooseAnswer copy")
					}
					questionText := strings.TrimSpace(questions[0].Question)
					if questionText == expectedPrompt || strings.HasPrefix(questionText, expectedPrompt+"\n") {
						t.Fatalf("first question for interface %s learning %s lost its task logic: %q", interfaceLanguage, code, questionText)
					}
					if strings.HasPrefix(questionText, userLearningLanguage(user).InterfaceName+": ") {
						t.Fatalf("first question for interface %s learning %s fell back to a bare language label: %q", interfaceLanguage, code, questionText)
					}
					for _, forbidden := range []string{"Fill the gap", "Choose the", "What does", "Выбери", "Заполни"} {
						if strings.Contains(questionText, forbidden) && !strings.Contains(expectedPrompt, forbidden) {
							t.Fatalf("first question for interface %s learning %s leaks raw instruction %q: %q", interfaceLanguage, code, forbidden, questionText)
						}
					}
					for _, option := range questions[0].Options {
						if strings.TrimSpace(option) == "" {
							t.Fatalf("empty level assessment option for %s", code)
						}
					}
					if scriptCheck, ok := targetOptionScriptChecks[code]; ok && !anyOptionMatches(questions[0].Options, scriptCheck) {
						t.Fatalf("options for %s do not look like target-language answer options: %v", code, questions[0].Options)
					}
				})
			}
		})
	}
}

var targetOptionScriptChecks = map[string]func(rune) bool{
	"ar": func(r rune) bool { return r >= 0x0600 && r <= 0x06FF },
	"bn": func(r rune) bool { return r >= 0x0980 && r <= 0x09FF },
	"el": func(r rune) bool { return r >= 0x0370 && r <= 0x03FF },
	"hi": func(r rune) bool { return r >= 0x0900 && r <= 0x097F },
	"ta": func(r rune) bool { return r >= 0x0B80 && r <= 0x0BFF },
	"te": func(r rune) bool { return r >= 0x0C00 && r <= 0x0C7F },
	"th": func(r rune) bool { return r >= 0x0E00 && r <= 0x0E7F },
}

func anyOptionMatches(options []string, match func(rune) bool) bool {
	for _, option := range options {
		for _, r := range option {
			if match(r) {
				return true
			}
		}
	}
	return false
}
