package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
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

func TestOnboardingWelcomeStaysShortAndLocalized(t *testing.T) {
	t.Run("unselected is short", func(t *testing.T) {
		text := onboardingWelcomeText(userState{})
		if strings.Contains(text, "->") {
			t.Fatalf("unselected onboarding welcome should not enumerate every interface language: %q", text)
		}
		if len(strings.Split(strings.TrimSpace(text), "\n")) > 4 {
			t.Fatalf("unselected onboarding welcome should stay short: %q", text)
		}
	})

	t.Run("selected uses localized copy only", func(t *testing.T) {
		user := userState{InterfaceLanguage: "ru", InterfaceSelected: true}
		copy := ui(user)
		want := "Poliglot AI\n\n" + copy.ChooseBotLang + "\n" + copy.ChooseTimezone + "\n" + copy.ChooseLearnLang
		text := onboardingWelcomeText(user)
		if text != want {
			t.Fatalf("selected onboarding welcome = %q, want %q", text, want)
		}
		for _, language := range interfaceLanguages() {
			if strings.Contains(text, language.InterfaceName) && language.Code != user.InterfaceLanguage {
				t.Fatalf("selected onboarding welcome should not mention other interface languages: %q", text)
			}
		}
	})
}

func TestAllInterfaceLanguagesUseExplicitLiteraryBotCopy(t *testing.T) {
	languages := interfaceLanguages()
	if len(languages) != 35 {
		t.Fatalf("interface language count = %d, want 35", len(languages))
	}

	englishUI := stringFields(englishUICopy())
	russianUI := stringFields(ui(userState{InterfaceLanguage: "ru"}))
	englishSystem := stringFields(systemUICopies["en"])
	russianSystem := stringFields(systemUICopies["ru"])
	englishPremium := stringFields(englishPremiumUICopy())
	russianPremium := stringFields(premiumUI(userState{InterfaceLanguage: "ru"}))

	for _, language := range languages {
		code := language.Code
		user := userState{InterfaceLanguage: code, InterfaceSelected: true, LearningLanguage: "en"}

		assertNoExactFallbackCopy(t, code, "ui", stringFields(ui(user)), englishUI, russianUI)
		assertNoExactFallbackCopy(t, code, "system", stringFields(systemUI(user)), englishSystem, russianSystem)
		assertNoExactFallbackCopy(t, code, "premium", stringFields(premiumUI(user)), englishPremium, russianPremium)

		for _, learning := range languages {
			name := strings.TrimSpace(localizedLanguageNameForInterface(learning, code))
			if name == "" {
				t.Fatalf("%s language name for %s is empty", code, learning.Code)
			}
			if code != "en" && name == strings.TrimSpace(localizedLanguageNameForInterface(learning, "en")) {
				t.Fatalf("%s language name for %s falls back to English: %q", code, learning.Code, name)
			}
			if code != "ru" && name == strings.TrimSpace(localizedLanguageNameForInterface(learning, "ru")) {
				t.Fatalf("%s language name for %s falls back to Russian: %q", code, learning.Code, name)
			}
		}
	}
}

func TestSystemUICopyDoesNotUseCompactGeneratedFallbacks(t *testing.T) {
	compactMarkers := map[string]string{
		"LevelDontKnow":      "?",
		"ReminderStatus":     "вЂў",
		"ReminderEnabled":    "+",
		"ReminderDisabled":   "-",
		"ReminderDefault":    "UTC",
		"ReminderHour":       "19:00",
		"ReminderFooter":     "%s",
		"VocabLevelQuestion": "%s: %s",
	}

	for _, language := range interfaceLanguages() {
		if language.Code == "en" || language.Code == "ru" {
			continue
		}
		fields := stringFields(systemUI(userState{InterfaceLanguage: language.Code}))
		for field, marker := range compactMarkers {
			if fields[field] == marker {
				t.Fatalf("%s system %s still uses compact generated fallback %q", language.Code, field, marker)
			}
		}
		for field, value := range fields {
			if strings.TrimSpace(value) == "" {
				t.Fatalf("%s system %s is empty", language.Code, field)
			}
		}
	}
}

func TestBotRuntimeAndPaymentCopyAvoidEnglishRussianFallback(t *testing.T) {
	englishRuntime := runtimeTextFields(userState{InterfaceLanguage: "en"})
	russianRuntime := runtimeTextFields(userState{InterfaceLanguage: "ru"})

	plan := premiumPlan{Title: "Premium 30", DaysLabel: "30 days", RubPrice: 300, StarsPrice: 150}
	payment := cryptoPayment{
		ID:        "pay-test",
		Amount:    "1.23",
		Currency:  "TON",
		Network:   "TON",
		Address:   "EQTEST",
		Memo:      "memo-test",
		Status:    cryptoStatusPending,
		ExpiresAt: time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC),
	}
	englishPayment := paymentCopyFields(userState{InterfaceLanguage: "en"}, plan, payment)
	russianPayment := paymentCopyFields(userState{InterfaceLanguage: "ru"}, plan, payment)

	for _, language := range interfaceLanguages() {
		code := language.Code
		user := userState{InterfaceLanguage: code, InterfaceSelected: true}
		assertNoExactFallbackCopy(t, code, "bot runtime", runtimeTextFields(user), englishRuntime, russianRuntime)
		assertNoExactFallbackCopy(t, code, "payment", paymentCopyFields(user, plan, payment), englishPayment, russianPayment)
	}
}

func TestVocabularyFallbackPromptUsesInterfaceLanguageMessage(t *testing.T) {
	word := vocabWord{
		ID:       "en:station",
		Language: "en",
		English:  "station",
		Russian:  "СЃС‚Р°РЅС†РёСЏ",
		Translations: map[string]string{
			"ru": "СЃС‚Р°РЅС†РёСЏ",
		},
	}

	for _, language := range interfaceLanguages() {
		code := language.Code
		if code == "en" || code == "ru" {
			continue
		}
		got := vocabularyFallbackPrompt(word, code)
		if strings.TrimSpace(got) == "" {
			t.Fatalf("%s vocabulary fallback is empty", code)
		}
		if got == firstDictionaryValue(word.Russian) || got == firstDictionaryValue(word.English) {
			t.Fatalf("%s vocabulary fallback leaked Russian/English dictionary value: %q", code, got)
		}
	}
}

func assertNoExactFallbackCopy(t *testing.T, code, area string, got, english, russian map[string]string) {
	t.Helper()
	for field, value := range got {
		value = strings.TrimSpace(value)
		if value == "" {
			t.Fatalf("%s %s %s is empty", code, area, field)
		}
		if isAllowedInvariantCopyValue(value) {
			continue
		}
		if code != "en" && value == strings.TrimSpace(english[field]) && shouldFailExactFallbackValue(value) {
			t.Fatalf("%s %s %s falls back to English: %q", code, area, field, value)
		}
		if code != "ru" && value == strings.TrimSpace(russian[field]) && shouldFailExactFallbackValue(value) {
			t.Fatalf("%s %s %s falls back to Russian: %q", code, area, field, value)
		}
	}
}

func shouldFailExactFallbackValue(value string) bool {
	value = strings.TrimSpace(value)
	if strings.Contains(value, "\n") {
		return true
	}
	if len([]rune(value)) > 18 {
		return true
	}
	if strings.Contains(value, "%") {
		return false
	}
	return strings.Contains(value, " ") && strings.ContainsAny(value, ".!?:;")
}

func runtimeTextFields(user userState) map[string]string {
	out := map[string]string{}
	for key := range botRuntimeTexts["en"] {
		out[key] = botRuntimeText(user, key)
	}
	return out
}

func paymentCopyFields(user userState, plan premiumPlan, payment cryptoPayment) map[string]string {
	return map[string]string{
		"premium_payment_options": premiumPaymentOptionsText(user, plan),
		"crypto_text":             cryptoPaymentText(user, plan, payment, "https://pay.example/invoice"),
		"crypto_pending":          cryptoPaymentPendingText(user, plan, payment),
		"crypto_text_v2":          cryptoPaymentTextV2(user, plan, payment, "https://pay.example/invoice"),
		"crypto_pending_v2":       cryptoPaymentPendingTextV2(user, plan, payment),
	}
}

func stringFields(value any) map[string]string {
	out := map[string]string{}
	collectStringFields(reflect.ValueOf(value), "", out)
	return out
}

func collectStringFields(value reflect.Value, prefix string, out map[string]string) {
	if !value.IsValid() {
		return
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	switch value.Kind() {
	case reflect.Struct:
		valueType := value.Type()
		for index := 0; index < value.NumField(); index++ {
			fieldName := valueType.Field(index).Name
			if prefix != "" {
				fieldName = prefix + "." + fieldName
			}
			collectStringFields(value.Field(index), fieldName, out)
		}
	case reflect.String:
		out[prefix] = value.String()
	case reflect.Slice:
		if value.Type().Elem().Kind() != reflect.String {
			return
		}
		for index := 0; index < value.Len(); index++ {
			out[fmt.Sprintf("%s[%d]", prefix, index)] = value.Index(index).String()
		}
	}
}

func isAllowedInvariantCopyValue(value string) bool {
	value = strings.TrimSpace(value)
	switch value {
	case "Premium", "Platinum", "AI Tutor", "GPT Agent", "GPT Агент", "GPT РђРіРµРЅС‚", "SBP", "YooKassa", "USDT", "Stars", "XP", "A1", "C2", "Plan", "Status", "Auto", "Menu", "Меню", "Тариф", "Прогресс", "Практика", "Стоп", "Слова", "📋 Menu", "📋 Меню", "⏹ Stop", "⏹ Стоп", "◀ Back", "◀ Menu", "◀ Меню", "Question %d/%d\n\n%s", "%s: %d Stars":
		return true
	}
	return false
}
