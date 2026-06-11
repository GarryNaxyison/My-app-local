package main

import (
	"strings"
	"testing"
)

func validAITutorLessonPayloadForTest() aiTutorLessonPayload {
	return aiTutorLessonPayload{
		Title:             "A Morning Visit",
		Level:             "A1",
		LevelBand:         "A1-A2",
		TargetLanguage:    "en",
		InterfaceLanguage: "ru",
		Theme:             "daily life",
		LessonGoal:        "Understand and retell a short everyday story.",
		GrammarFocus: aiTutorGrammarFocus{
			Name:                      "Simple present",
			ShortExplanationInterface: "Use simple present for routines.",
			ModelSentenceTarget:       "I go to the shop.",
		},
		Story: aiTutorStory{
			TextTarget:      "Mia wakes up early. She drinks water. She walks to a small shop. She buys bread. She says thank you. She goes home.",
			AudioTextTarget: "Mia wakes up early. She drinks water. She walks to a small shop. She buys bread. She says thank you. She goes home.",
			SentenceCount:   6,
		},
		Words: []aiTutorWord{
			{ID: "w1", Target: "wake up", InterfaceTranslation: "prosypatsya", PartOfSpeech: "verb", ExampleSentenceTarget: "I wake up early.", ExampleTranslationInterface: "Ya prosypayus rano.", AudioTextTarget: "wake up", ExampleAudioTextTarget: "I wake up early."},
			{ID: "w2", Target: "water", InterfaceTranslation: "voda", PartOfSpeech: "noun", ExampleSentenceTarget: "I drink water.", ExampleTranslationInterface: "Ya p'yu vodu.", AudioTextTarget: "water", ExampleAudioTextTarget: "I drink water."},
			{ID: "w3", Target: "shop", InterfaceTranslation: "magazin", PartOfSpeech: "noun", ExampleSentenceTarget: "The shop is small.", ExampleTranslationInterface: "Magazin malenkiy.", AudioTextTarget: "shop", ExampleAudioTextTarget: "The shop is small."},
			{ID: "w4", Target: "bread", InterfaceTranslation: "hleb", PartOfSpeech: "noun", ExampleSentenceTarget: "She buys bread.", ExampleTranslationInterface: "Ona pokupaet hleb.", AudioTextTarget: "bread", ExampleAudioTextTarget: "She buys bread."},
			{ID: "w5", Target: "thank you", InterfaceTranslation: "spasibo", PartOfSpeech: "chunk", ExampleSentenceTarget: "Thank you for the bread.", ExampleTranslationInterface: "Spasibo za hleb.", AudioTextTarget: "thank you", ExampleAudioTextTarget: "Thank you for the bread."},
			{ID: "w6", Target: "home", InterfaceTranslation: "dom", PartOfSpeech: "noun", ExampleSentenceTarget: "She goes home.", ExampleTranslationInterface: "Ona idet domoy.", AudioTextTarget: "home", ExampleAudioTextTarget: "She goes home."},
		},
		RetellTask: aiTutorRetellTask{
			InstructionInterface: "Retell the story in your own words.",
			MinSentences:         2,
			FeedbackRubric:       []string{"mentions Mia", "mentions shop", "uses target language"},
		},
		ComprehensionQuestions: []aiTutorQuestion{
			{ID: "q1", QuestionTarget: "When does Mia wake up?", ExpectedPoints: []string{"early"}, FeedbackRuleInterface: "Check time detail."},
			{ID: "q2", QuestionTarget: "What does Mia buy?", ExpectedPoints: []string{"bread"}, FeedbackRuleInterface: "Check object detail."},
			{ID: "q3", QuestionTarget: "Where does Mia go at the end?", ExpectedPoints: []string{"home"}, FeedbackRuleInterface: "Check final place."},
		},
		WordLearning: []aiTutorWordTask{{WordID: "w1"}, {WordID: "w2"}, {WordID: "w3"}, {WordID: "w4"}, {WordID: "w5"}, {WordID: "w6"}},
		WordRecall:   []aiTutorWordTask{{WordID: "w1"}, {WordID: "w2"}, {WordID: "w3"}, {WordID: "w4"}, {WordID: "w5"}, {WordID: "w6"}},
		ProductionTask: aiTutorProductionTask{
			InstructionInterface:     "Write 2-3 sentences with at least three lesson words.",
			RequiredWordCount:        3,
			SentenceCount:            "2-3",
			EvaluationCriteria:       []string{"2-3 sentences", "uses at least three words", "target language"},
			RecommendationsInterface: []string{"Repeat thank you as one polite chunk.", "Keep A1 sentences short."},
		},
		ReviewOptions: []string{"tomorrow", "3_days", "1_week", "no_review"},
		QualitySelfCheck: aiTutorQualitySelfCheck{
			CEFRReason:    "Short A1 sentences and high-frequency words.",
			WhyReusable:   "Everyday routine topic.",
			DuplicateRisk: "low",
		},
	}
}

func TestValidateAITutorLessonRequiresAudioAndFinalRecommendations(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()
	lesson.Story.AudioTextTarget = ""
	lesson.Words[0].AudioTextTarget = ""
	lesson.Words[1].ExampleAudioTextTarget = ""
	lesson.ProductionTask.RecommendationsInterface = nil
	issues := strings.Join(validateAITutorLessonPayload(lesson), "\n")
	for _, want := range []string{"story audio", "word 1 audio", "word 2 example audio", "final recommendations"} {
		if !strings.Contains(issues, want) {
			t.Fatalf("issues miss %q:\n%s", want, issues)
		}
	}
}

func TestAITutorLevelBand(t *testing.T) {
	cases := map[string]string{
		"A1":    "A1-A2",
		"A2":    "A1-A2",
		"B1":    "B1-B2",
		"B2":    "B1-B2",
		"C1":    "C1-C2",
		"C2":    "C1-C2",
		"B1-B2": "B1-B2",
		"C1-C2": "C1-C2",
		"":      "A1-A2",
	}
	for level, want := range cases {
		if got := aiTutorLevelBand(level); got != want {
			t.Fatalf("aiTutorLevelBand(%q) = %q, want %q", level, got, want)
		}
	}
}

func TestAITutorRecallStepBuildsMultipleChoiceOptions(t *testing.T) {
	step := aiTutorBuildStep(validAITutorLessonPayloadForTest(), aiTutorWordRecallStage(1))
	if step.Kind != "word_recall" {
		t.Fatalf("kind = %q, want word_recall", step.Kind)
	}
	if len(step.Options) < 4 {
		t.Fatalf("recall options len = %d, want at least 4: %#v", len(step.Options), step.Options)
	}
	correct := 0
	for _, option := range step.Options {
		if strings.TrimSpace(option.ID) == "" || strings.TrimSpace(option.Text) == "" {
			t.Fatalf("recall option must include id and text: %#v", option)
		}
		if option.Correct {
			correct++
		}
	}
	if correct != 1 {
		t.Fatalf("recall options correct count = %d, want 1: %#v", correct, step.Options)
	}
}

func TestValidateAITutorLessonAcceptsValidPayload(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()
	if issues := validateAITutorLessonPayload(lesson); len(issues) != 0 {
		t.Fatalf("valid lesson rejected: %#v", issues)
	}
}

func TestValidateAITutorLessonRejectsWrongCounts(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()
	lesson.Words = lesson.Words[:5]
	lesson.ComprehensionQuestions = lesson.ComprehensionQuestions[:2]
	lesson.Story.TextTarget = "One. Two. Three. Four."
	lesson.Story.SentenceCount = 4
	issues := strings.Join(validateAITutorLessonPayload(lesson), "\n")
	for _, want := range []string{"exactly 6 words", "exactly 3 questions", "5-7 story sentences"} {
		if !strings.Contains(issues, want) {
			t.Fatalf("issues miss %q:\n%s", want, issues)
		}
	}
}

func TestAITutorCanonicalStages(t *testing.T) {
	stages := aiTutorCanonicalStages()
	if stages[0] != aiTutorStageStoryIntro {
		t.Fatalf("first stage = %q", stages[0])
	}
	if stages[len(stages)-1] != aiTutorStageComplete {
		t.Fatalf("last stage = %q", stages[len(stages)-1])
	}
	if len(stages) != 21 {
		t.Fatalf("stage count = %d, want 21", len(stages))
	}
}
