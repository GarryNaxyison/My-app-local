package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestJSONStoreExtendPremiumAddsToExistingPremium(t *testing.T) {
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.getOrCreateUser(123, "Test")
	if err != nil {
		t.Fatal(err)
	}

	firstUntil, err := store.extendPremium(user.TelegramID, "payment-1", 30*24*time.Hour, "premium")
	if err != nil {
		t.Fatal(err)
	}
	secondUntil, err := store.extendPremium(user.TelegramID, "payment-2", 365*24*time.Hour, "premium")
	if err != nil {
		t.Fatal(err)
	}

	want := firstUntil.Add(365 * 24 * time.Hour)
	if !secondUntil.Equal(want) {
		t.Fatalf("secondUntil = %s, want %s", secondUntil, want)
	}
}

func TestJSONStoreExtendPremiumIgnoresSamePayment(t *testing.T) {
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatal(err)
	}

	user, err := store.getOrCreateUser(123, "Test")
	if err != nil {
		t.Fatal(err)
	}

	firstUntil, err := store.extendPremium(user.TelegramID, "payment-1", 30*24*time.Hour, "premium")
	if err != nil {
		t.Fatal(err)
	}
	secondUntil, err := store.extendPremium(user.TelegramID, "payment-1", 365*24*time.Hour, "premium")
	if err != nil {
		t.Fatal(err)
	}

	if !secondUntil.Equal(firstUntil) {
		t.Fatalf("duplicate payment changed premium until: got %s, want %s", secondUntil, firstUntil)
	}
}

func TestPremiumOnlyLimits(t *testing.T) {
	free := userState{Plan: "free"}
	lessonLimit, practiceLimit := limitsFor(free)
	if lessonLimit != freeLessonLimit || practiceLimit != freePracticeLimit {
		t.Fatalf("free limits = %d, %d; want %d, %d", lessonLimit, practiceLimit, freeLessonLimit, freePracticeLimit)
	}
	if voiceLimitFor(free) != 0 {
		t.Fatalf("free voice limit = %d; want 0", voiceLimitFor(free))
	}

	premium := userState{Plan: "premium", PremiumUntil: time.Now().Add(24 * time.Hour)}
	lessonLimit, practiceLimit = limitsFor(premium)
	if lessonLimit != premiumLessonLimit || practiceLimit != premiumPracticeLimit {
		t.Fatalf("premium limits = %d, %d; want %d, %d", lessonLimit, practiceLimit, premiumLessonLimit, premiumPracticeLimit)
	}
	if voiceLimitFor(premium) != premiumVoiceLimit {
		t.Fatalf("premium voice limit = %d; want %d", voiceLimitFor(premium), premiumVoiceLimit)
	}
}

func TestPremiumInvoiceDescriptionIsComplete(t *testing.T) {
	maxVoiceSeconds := 30
	b := &bot{cfg: config{MaxVoiceSeconds: maxVoiceSeconds}}
	description := b.premiumInvoiceDescription(userState{InterfaceLanguage: "ru"}, premiumPlan{Tier: "premium"})
	for _, want := range []string{
		fmt.Sprintf("%d уроков", premiumLessonLimit),
		fmt.Sprintf("%d сообщений", premiumPracticeLimit),
		fmt.Sprintf("%d голосовых", premiumVoiceLimit),
		fmt.Sprintf("%d секунд", maxVoiceSeconds),
		"голос в текст",
		"картинки",
	} {
		if !strings.Contains(description, want) {
			t.Fatalf("premium invoice description is missing %q: %q", want, description)
		}
	}
	if len([]rune(description)) > 255 {
		t.Fatalf("telegram invoice description is too long: %d chars", len([]rune(description)))
	}
}

func TestPremiumTextUsesInterfaceLanguage(t *testing.T) {
	b := &bot{cfg: config{MaxVoiceSeconds: 30, PremiumRubPrice: 300, PremiumStarsPrice: 150, PremiumYearRubPrice: 3000, PremiumYearStarsPrice: 1500, PlatinumRubPrice: 590, PlatinumStarsPrice: 300, PlatinumYearRubPrice: 5900, PlatinumYearStarsPrice: 3000}}
	text := plainTextFromMarkdownV2(b.premiumText(userState{InterfaceLanguage: "en"}))
	if strings.Contains(text, "уроков") || !strings.Contains(text, "lessons per day") {
		t.Fatalf("expected English premium text, got %q", text)
	}

	keyboard := premiumInlineKeyboard(true, nil, localizedPremiumPlans(userState{InterfaceLanguage: "es"}, b.premiumPlanList()), userState{InterfaceLanguage: "es"})
	rows := keyboard["inline_keyboard"].([][]map[string]any)
	if rows[0][0]["callback_data"].(string) != "premium_plan|"+premiumMonthlyProduct ||
		rows[1][0]["callback_data"].(string) != "premium_plan|"+premiumYearlyProduct ||
		!strings.Contains(rows[1][0]["text"].(string), "año") ||
		rows[len(rows)-3][0]["text"].(string) != "Invitar a un amigo" {
		t.Fatalf("expected Spanish premium keyboard, got %#v", rows)
	}
}

func TestPremiumTextUsesLandingPlanDescriptionsForRussian(t *testing.T) {
	b := &bot{cfg: config{MaxVoiceSeconds: 30, PremiumRubPrice: 300, PremiumStarsPrice: 150, PremiumYearRubPrice: 3000, PremiumYearStarsPrice: 1500, PlatinumRubPrice: 590, PlatinumStarsPrice: 300, PlatinumYearRubPrice: 5900, PlatinumYearStarsPrice: 3000}}
	text := plainTextFromMarkdownV2(b.premiumText(userState{InterfaceLanguage: "ru"}))

	for _, want := range []string{
		"Базовое текстовое обучение, тренировка слов, Phrasebook и обзор прогресса. AI Tutor, аудирование, произношение и голосовые проверки открываются в Premium.",
		"Основной режим для ежедневной практики: AI Tutor, listening, pronunciation, guided AI lessons, voice checks, photo tools и расширенные дневные лимиты.",
		"AI Tutor с максимальными дневными лимитами, глубиной roleplay, интенсивным review и максимальной voice/pronunciation практикой.",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("premium text missing landing description %q:\n%s", want, text)
		}
	}
}

func TestPremiumPaymentOptionsKeyboardShowsMethodsForSelectedPlan(t *testing.T) {
	user := userState{InterfaceLanguage: "en"}
	plan := premiumPlan{
		Product:    premiumYearlyProduct,
		Title:      "Premium for a year",
		DaysLabel:  "365 days",
		RubPrice:   3000,
		StarsPrice: 1500,
	}
	keyboard := premiumPaymentOptionsKeyboard(true, []cryptoPaymentMethod{{ID: cryptoMethodTON, Label: "TON"}}, plan, user)
	callbacks := collectCallbackData(t, keyboard)
	for _, want := range []string{
		"buy_stars|" + premiumYearlyProduct,
		"buy_yookassa|" + premiumYearlyProduct,
		"buy_crypto|" + premiumYearlyProduct + "|" + cryptoMethodTON,
		"menu_premium",
	} {
		if !callbacks[want] {
			t.Fatalf("payment options keyboard is missing %q; got %#v", want, callbacks)
		}
	}
}

func TestPremiumPaymentPayloadCarriesUserAndTimestamp(t *testing.T) {
	b := &bot{cfg: config{PremiumYearRubPrice: 3000, PremiumYearStarsPrice: 1500}}
	payload := premiumPaymentPayload(premiumYearlyProduct, 123, time.Unix(1700000000, 0))

	plan, userID, timestamp, ok := b.premiumPlanPaymentFromPayload(payload)
	if !ok {
		t.Fatalf("expected valid payload")
	}
	if plan.Product != premiumYearlyProduct || userID != 123 || timestamp != 1700000000 {
		t.Fatalf("unexpected parsed payload: plan=%s user=%d timestamp=%d", plan.Product, userID, timestamp)
	}

	for _, bad := range []string{
		premiumYearlyProduct,
		premiumYearlyProduct + "_0_1700000000",
		premiumYearlyProduct + "_123_0",
		premiumYearlyProduct + "_abc_1700000000",
		premiumYearlyProduct + "_123_abc",
	} {
		if _, _, _, ok := b.premiumPlanPaymentFromPayload(bad); ok {
			t.Fatalf("expected invalid payload %q", bad)
		}
	}
}

func TestPremiumStarsStartPayloadBuildsTelegramDeepLink(t *testing.T) {
	payload := premiumStarsStartPayload(platinumYearlyProduct)
	product, ok := parsePremiumStarsStartPayload(payload)
	if !ok || product != platinumYearlyProduct {
		t.Fatalf("unexpected start payload parse: product=%q ok=%v", product, ok)
	}
	link := telegramBotStartURL("@Poliglot_AI_bot", payload)
	if link != "https://t.me/Poliglot_AI_bot?start=buy_platinum_365d" {
		t.Fatalf("telegram deep link = %q", link)
	}
}

func TestCJKReferralShareTextIncludesWebAndPremium(t *testing.T) {
	for _, code := range []string{"zh", "ja", "ko"} {
		text := fmt.Sprintf(premiumUI(userState{InterfaceLanguage: code, InterfaceSelected: true}).ReferralShareText, "https://example.test/app?ref=CODE")
		if !strings.Contains(text, "https://example.test/app?ref=CODE") {
			t.Fatalf("%s referral text does not include link: %q", code, text)
		}
		if !strings.Contains(text, "Poliglot AI") || !strings.Contains(text, "Premium") || !(strings.Contains(strings.ToLower(text), "web") || strings.Contains(text, "웹")) {
			t.Fatalf("%s referral text should mention Poliglot AI, web, and Premium: %q", code, text)
		}
	}
}
