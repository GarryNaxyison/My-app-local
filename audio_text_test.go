package main

import "testing"

func TestLessonCorrectionPronunciationTextUsesCorrectedSection(t *testing.T) {
	feedback := "1) Corrected version:\nI went home yesterday.\n\n2) Improve:\nUse went.\n\n3) Summary:\nGood."
	got := lessonCorrectionPronunciationText(feedback, []mistakeEntry{{Correction: "I went home."}})
	if got != "I went home yesterday." {
		t.Fatalf("lesson correction audio = %q", got)
	}
}

func TestLessonCorrectionPronunciationTextFallsBackToMistakes(t *testing.T) {
	got := lessonCorrectionPronunciationText("Good work.", []mistakeEntry{{Correction: "I went home."}})
	if got != "I went home." {
		t.Fatalf("lesson fallback audio = %q", got)
	}
}

func TestPracticeCorrectionPronunciationTextPrefersModelPhrase(t *testing.T) {
	reply := "Correction:\nUse went, not goed.\n\nModel phrase:\nI went home yesterday.\n\nYour turn:\nWhere did you go yesterday?"
	got := practiceCorrectionPronunciationText(reply, []mistakeEntry{{Correction: "went"}})
	if got != "I went home yesterday." {
		t.Fatalf("practice correction audio = %q", got)
	}
}

func TestPracticeCorrectionPronunciationTextSupportsSameLineModelPhrase(t *testing.T) {
	reply := "Model phrase: I usually drink coffee in the morning.\nYour turn: What do you drink?"
	got := practiceCorrectionPronunciationText(reply, []mistakeEntry{{Correction: "drink coffee"}})
	if got != "I usually drink coffee in the morning." {
		t.Fatalf("practice same-line correction audio = %q", got)
	}
}

func TestPracticeCorrectionPronunciationTextStripsInterfaceTranslation(t *testing.T) {
	reply := "Model phrase: I would like to order some coffee, please. (\u042f \u0431\u044b \u0445\u043e\u0442\u0435\u043b \u0437\u0430\u043a\u0430\u0437\u0430\u0442\u044c \u043a\u043e\u0444\u0435, \u043f\u043e\u0436\u0430\u043b\u0443\u0439\u0441\u0442\u0430.)"
	got := practiceCorrectionPronunciationText(reply, nil)
	if got != "I would like to order some coffee, please." {
		t.Fatalf("practice translated correction audio = %q", got)
	}
}

func TestPracticeCorrectionPronunciationTextStripsLatinTranslation(t *testing.T) {
	target := "\u042f \u0445\u043e\u0447\u0443 \u0437\u0430\u043a\u0430\u0437\u0430\u0442\u044c \u043a\u043e\u0444\u0435, \u043f\u043e\u0436\u0430\u043b\u0443\u0439\u0441\u0442\u0430."
	reply := "Model phrase: " + target + " (I would like to order coffee, please.)"
	got := practiceCorrectionPronunciationText(reply, nil)
	if got != target {
		t.Fatalf("practice latin translated correction audio = %q", got)
	}
}

func TestSanitizeAudioTextKeepsShortClarifyingParenthetical(t *testing.T) {
	got := sanitizeAudioText("I live in Washington (state)")
	if got != "I live in Washington (state)" {
		t.Fatalf("clarifying parenthetical should stay, got %q", got)
	}
}

func TestPracticeQuestionPronunciationTextUsesFinalQuestion(t *testing.T) {
	reply := "Good, corrected version: I went home.\n\nWhere did you go yesterday?"
	got := practiceQuestionPronunciationText(reply)
	if got != "Where did you go yesterday?" {
		t.Fatalf("practice question audio = %q", got)
	}
}

func TestPracticeQuestionPronunciationTextStripsQuestionLabel(t *testing.T) {
	reply := "Nice.\nQuestion: Where do you live?"
	got := practiceQuestionPronunciationText(reply)
	if got != "Where do you live?" {
		t.Fatalf("practice labeled question audio = %q", got)
	}
}

func TestLessonQuestionPronunciationTextPrefersTargetExample(t *testing.T) {
	task := "\u0421\u0438\u0442\u0443\u0430\u0446\u0438\u044f: \u0442\u044b \u0432 \u043a\u0430\u0444\u0435.\n\n" +
		"\u041f\u0440\u0438\u043c\u0435\u0440 \u0444\u0440\u0430\u0437\u044b:\nI would like a coffee, please.\n\n" +
		"\u0422\u0432\u043e\u044f \u043e\u0447\u0435\u0440\u0435\u0434\u044c:\n\u041d\u0430\u043f\u0438\u0448\u0438 \u043e\u0434\u0438\u043d \u043a\u043e\u0440\u043e\u0442\u043a\u0438\u0439 \u043e\u0442\u0432\u0435\u0442 \u043d\u0430 \u0430\u043d\u0433\u043b\u0438\u0439\u0441\u043a\u043e\u043c."
	got := lessonQuestionPronunciationText(task)
	if got != "I would like a coffee, please." {
		t.Fatalf("lesson question audio = %q", got)
	}
}

func TestLessonQuestionPronunciationTextDoesNotReadInterfaceInstruction(t *testing.T) {
	task := "\u0421\u0438\u0442\u0443\u0430\u0446\u0438\u044f: \u0442\u044b \u0432 \u043a\u0430\u0444\u0435.\n\n" +
		"\u0422\u0432\u043e\u044f \u043e\u0447\u0435\u0440\u0435\u0434\u044c:\n\u041d\u0430\u043f\u0438\u0448\u0438 \u043e\u0434\u0438\u043d \u043a\u043e\u0440\u043e\u0442\u043a\u0438\u0439 \u043e\u0442\u0432\u0435\u0442 \u043d\u0430 \u0430\u043d\u0433\u043b\u0438\u0439\u0441\u043a\u043e\u043c."
	got := lessonQuestionPronunciationText(task)
	if got != "" {
		t.Fatalf("lesson instruction should not be voiced, got %q", got)
	}
}

func TestRoleplayQuestionPronunciationTextFallsBackToLastDialogueLine(t *testing.T) {
	reply := "Poliglot AI: Welcome to the station. Tell me where you need to go.\n\nCorrection: keep it short."
	got := roleplayQuestionPronunciationText(reply)
	if got != "Welcome to the station. Tell me where you need to go." {
		t.Fatalf("roleplay fallback audio = %q", got)
	}
}
