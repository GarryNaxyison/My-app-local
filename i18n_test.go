package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestEveryInterfaceLanguageHasUICopy(t *testing.T) {
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code})
		if copy.MainMenuTitle == "" || copy.ChooseTimezone == "" || copy.ChooseLearnLang == "" || copy.Tool.VoiceToText == "" || copy.Tool.ImageTranslate == "" || copy.Tool.VoicePrompt == "" || copy.Tool.ImagePrompt == "" {
			t.Fatalf("missing UI copy for %s: %#v", language.Code, copy)
		}
	}
}

func TestInterfaceLanguagesDoNotDuplicateRussian(t *testing.T) {
	count := 0
	for _, language := range interfaceLanguages() {
		if language.Code == "ru" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected one Russian interface language entry, got %d", count)
	}
}

func TestPracticeStartedFormattingDoesNotLeakFmtExtra(t *testing.T) {
	learning := learningLanguageByCode("en")
	for _, language := range interfaceLanguages() {
		user := userState{InterfaceLanguage: language.Code, LearningLanguage: learning.Code}
		copy := ui(user)
		text := fmt.Sprintf(systemUI(user).PracticeStarted, learning.PracticeDirection, language.InterfaceName, copy.StopButton)
		if strings.Contains(text, "%!(") {
			t.Fatalf("practice start format leaked fmt marker for %s: %q", language.Code, text)
		}
		if !strings.Contains(text, copy.StopButton) {
			t.Fatalf("practice start format for %s does not include stop button %q: %q", language.Code, copy.StopButton, text)
		}
	}
}

func TestToolUICopyUsesLocalizedPromptsWhenAvailable(t *testing.T) {
	copy := ui(userState{InterfaceLanguage: "ky"})
	if copy.Tool.VoicePrompt == englishToolUICopy().VoicePrompt {
		t.Fatal("expected Kyrgyz tools menu prompt to be localized")
	}
	if copy.Tool.ImagePrompt == englishToolUICopy().ImagePrompt {
		t.Fatal("expected Kyrgyz image tool prompt to be localized")
	}
}

func TestEveryNonEnglishInterfaceLanguageAvoidsEnglishTelegramFallback(t *testing.T) {
	english := englishUICopy()
	englishTool := englishToolUICopy()

	for _, language := range interfaceLanguages() {
		if language.Code == "en" {
			continue
		}
		copy := ui(userState{InterfaceLanguage: language.Code})
		checks := map[string][2]string{
			"MainMenuTitle":       {copy.MainMenuTitle, english.MainMenuTitle},
			"MainMenuBody":        {copy.MainMenuBody, english.MainMenuBody},
			"LearnWords":          {copy.LearnWords, english.LearnWords},
			"WordGame":            {copy.WordGame, english.WordGame},
			"Spelling":            {copy.Spelling, english.Spelling},
			"Vocabulary":          {copy.Vocabulary, english.Vocabulary},
			"Mistakes":            {copy.Mistakes, english.Mistakes},
			"ChooseBotLang":       {copy.ChooseBotLang, english.ChooseBotLang},
			"ChooseLearnLang":     {copy.ChooseLearnLang, english.ChooseLearnLang},
			"Tool.VoiceToText":    {copy.Tool.VoiceToText, englishTool.VoiceToText},
			"Tool.ImageTranslate": {copy.Tool.ImageTranslate, englishTool.ImageTranslate},
			"Tool.VoicePrompt":    {copy.Tool.VoicePrompt, englishTool.VoicePrompt},
			"Tool.ImagePrompt":    {copy.Tool.ImagePrompt, englishTool.ImagePrompt},
		}
		for field, values := range checks {
			if values[0] == "" {
				t.Fatalf("missing %s for %s", field, language.Code)
			}
			if values[0] == values[1] {
				t.Fatalf("expected %s %s to avoid English fallback, got %q", language.Code, field, values[0])
			}
		}
	}
}

func TestCJKInterfaceLanguagesHaveLocalizedCopies(t *testing.T) {
	for _, code := range []string{"zh", "ja", "ko"} {
		copy := ui(userState{InterfaceLanguage: code})
		if copy.MainMenuTitle == englishUICopy().MainMenuTitle {
			t.Fatalf("expected %s main UI copy to be localized", code)
		}
		if copy.Tool.VoicePrompt == englishToolUICopy().VoicePrompt {
			t.Fatalf("expected %s tool UI copy to be localized", code)
		}
		systemCopy := systemUI(userState{InterfaceLanguage: code})
		if systemCopy.LevelStartText == systemUICopies["en"].LevelStartText {
			t.Fatalf("expected %s system UI copy to be localized", code)
		}
		premiumCopy := premiumUI(userState{InterfaceLanguage: code})
		if premiumCopy.MonthPrice == englishPremiumUICopy().MonthPrice {
			t.Fatalf("expected %s premium UI copy to be localized", code)
		}
	}
}

func TestOnboardingWelcomeContainsEveryInterfaceLanguage(t *testing.T) {
	text := onboardingWelcomeText(userState{})
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code, InterfaceSelected: true})
		if !strings.Contains(text, language.InterfaceName) || !strings.Contains(text, copy.ChooseBotLang) || !strings.Contains(text, copy.ChooseTimezone) {
			t.Fatalf("onboarding welcome is missing localized setup text for %s: %q", language.Code, text)
		}
	}
}
