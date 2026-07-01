package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestPremiumTextShowsCurrentPlanStatus(t *testing.T) {
	b := &bot{cfg: config{MaxVoiceSeconds: 30, PremiumRubPrice: 300, PremiumStarsPrice: 150, PremiumYearRubPrice: 3000, PremiumYearStarsPrice: 1500, PlatinumRubPrice: 590, PlatinumStarsPrice: 300, PlatinumYearRubPrice: 5900, PlatinumYearStarsPrice: 3000}}
	user := userState{
		InterfaceLanguage: "en",
		Plan:              "platinum",
		PremiumUntil:      time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC),
	}
	text := plainTextFromMarkdownV2(b.premiumText(user))
	for _, want := range []string{"Current plan", "Platinum", "Premium until", "2026-07-12"} {
		if !strings.Contains(text, want) {
			t.Fatalf("premium text missing %q:\n%s", want, text)
		}
	}
}

func TestPremiumTextUsesLandingPlanDescriptionsForRussian(t *testing.T) {
	b := &bot{cfg: config{MaxVoiceSeconds: 30, PremiumRubPrice: 300, PremiumStarsPrice: 150, PremiumYearRubPrice: 3000, PremiumYearStarsPrice: 1500, PlatinumRubPrice: 590, PlatinumStarsPrice: 300, PlatinumYearRubPrice: 5900, PlatinumYearStarsPrice: 3000}}
	text := plainTextFromMarkdownV2(b.premiumText(userState{InterfaceLanguage: "ru"}))

	for _, want := range []string{
		"Базовое текстовое обучение, тренировка слов, Phrasebook и обзор прогресса. AI Tutor, аудирование, произношение и голосовые проверки открываются в Premium.",
		"Основной режим для ежедневной практики: AI Tutor, аудирование, произношение, AI-уроки, проверки голоса, фото-инструменты и расширенные дневные лимиты.",
		"AI Tutor с максимальными дневными лимитами, глубокой ролевой практикой, интенсивным повторением и максимальной практикой голоса и произношения.",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("premium text missing landing description %q:\n%s", want, text)
		}
	}
}

func TestPremiumPlanLandingCopyMatchesLandingPricingCards(t *testing.T) {
	free := premiumPlanLandingCopy(userState{InterfaceLanguage: "ru"}, "free")
	if free.Label != "Начальный" || free.Description != "Базовое текстовое обучение, тренировка слов, Phrasebook и обзор прогресса. AI Tutor, аудирование, произношение и голосовые проверки открываются в Premium." {
		t.Fatalf("free landing copy mismatch: %+v", free)
	}
	assertStringSlicesEqual(t, free.Included, []string{"ежедневная привычка и стартовые уроки", "базовая тренировка слов", "заметки, phrasebook и обзор прогресса"})
	assertStringSlicesEqual(t, free.Locked, []string{"Уроки с AI Tutor", "Аудирование и произношение", "Проверка голоса и фото-инструменты"})
	if free.Note != "Подходит для знакомства с продуктом без оплаты." {
		t.Fatalf("free note mismatch: %q", free.Note)
	}

	premium := premiumPlanLandingCopy(userState{InterfaceLanguage: "ru"}, "premium")
	if premium.Label != "Регулярная учеба" || premium.Description != "Основной режим для ежедневной практики: AI Tutor, аудирование, произношение, AI-уроки, проверки голоса, фото-инструменты и расширенные дневные лимиты." {
		t.Fatalf("premium landing copy mismatch: %+v", premium)
	}
	assertStringSlicesEqual(t, premium.Included, []string{"голос в текст и перевод услышанного", "перевод текста с картинки", "практика по контексту голоса или фото", "словарь ошибок, notes, XP и streak"})
	if premium.Note != "Лучший выбор для стабильного ежедневного обучения." {
		t.Fatalf("premium note mismatch: %q", premium.Note)
	}

	platinum := premiumPlanLandingCopy(userState{InterfaceLanguage: "ru"}, "platinum")
	if platinum.Label != "Интенсив" || platinum.Description != "AI Tutor с максимальными дневными лимитами, глубокой ролевой практикой, интенсивным повторением и максимальной практикой голоса и произношения." {
		t.Fatalf("platinum landing copy mismatch: %+v", platinum)
	}
	assertStringSlicesEqual(t, platinum.Included, []string{"максимальные дневные лимиты", "больше практики по голосу или фото", "интенсивное повторение слабых мест", "лучший режим для плотной ежедневной учебы"})
	if platinum.Note != "Для поездки, работы, экзамена или очень плотного темпа." {
		t.Fatalf("platinum note mismatch: %q", platinum.Note)
	}
}

func TestPremiumLandingCopyLocalizesPremiumFeatureTerms(t *testing.T) {
	for _, language := range interfaceLanguages() {
		for _, tier := range []string{"free", "premium", "platinum"} {
			copy := premiumPlanLandingCopy(userState{InterfaceLanguage: language.Code}, tier)
			joined := strings.Join(append(append([]string{copy.Label, copy.Description, copy.Note}, copy.Included...), copy.Locked...), "\n")
			for _, forbidden := range []string{
				"Try the path",
				"Попробовать маршрут",
				"guided AI lessons",
				"guided lessons",
				"Listening and pronunciation",
				"Listening и pronunciation",
				"voice checks",
				"photo tools",
				"voice/pronunciation",
				"voice/photo-context",
				"heavy daily learning",
			} {
				if strings.Contains(joined, forbidden) {
					t.Fatalf("%s %s landing copy contains untranslated term %q:\n%s", language.Code, tier, forbidden, joined)
				}
			}
		}
	}
}

func TestPremiumPlanLandingCopyCoversAllInterfaceLanguages(t *testing.T) {
	languages := interfaceLanguages()
	if len(languages) != 35 {
		t.Fatalf("interface language count = %d, want 35", len(languages))
	}
	russian := map[string]string{
		"free":     premiumPlanLandingCopy(userState{InterfaceLanguage: "ru"}, "free").Description,
		"premium":  premiumPlanLandingCopy(userState{InterfaceLanguage: "ru"}, "premium").Description,
		"platinum": premiumPlanLandingCopy(userState{InterfaceLanguage: "ru"}, "platinum").Description,
	}

	for _, language := range languages {
		for _, tier := range []string{"free", "premium", "platinum"} {
			copy := premiumPlanLandingCopy(userState{InterfaceLanguage: language.Code}, tier)
			if copy.Label == "" || copy.Description == "" || copy.Note == "" {
				t.Fatalf("%s %s landing copy has empty headline fields: %+v", language.Code, tier, copy)
			}
			if len(copy.Included) == 0 {
				t.Fatalf("%s %s landing copy has no included features", language.Code, tier)
			}
			if tier == "free" && len(copy.Locked) == 0 {
				t.Fatalf("%s free landing copy has no locked features", language.Code)
			}
			if language.Code != "ru" && copy.Description == russian[tier] {
				t.Fatalf("%s %s landing copy still uses Russian description", language.Code, tier)
			}
		}
	}
}

func assertStringSlicesEqual(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("slice[%d] = %q, want %q; full=%#v", i, got[i], want[i], got)
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
	keyboard := premiumPaymentOptionsKeyboard(true, true, []cryptoPaymentMethod{{ID: cryptoMethodTON, Label: "TON"}}, plan, user)
	callbacks := collectCallbackData(t, keyboard)
	for _, want := range []string{
		"buy_stars|" + premiumYearlyProduct,
		"buy_rollypay|" + premiumYearlyProduct,
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

func TestRollyPayConfigRequiresSeparateCashboxes(t *testing.T) {
	cfg := config{
		RollyPayBotCashboxID: "bot-cashbox",
		RollyPayBotAPIKey:    "bot-key",
		RollyPayWebCashboxID: "web-cashbox",
		RollyPayWebAPIKey:    "web-key",
	}
	if !cfg.rollyPayBotEnabled() {
		t.Fatal("expected Telegram RollyPay to be enabled")
	}
	if !cfg.rollyPayWebEnabled() {
		t.Fatal("expected web RollyPay to be enabled")
	}
	cfg.RollyPayWebAPIKey = ""
	if cfg.rollyPayWebEnabled() {
		t.Fatal("web RollyPay must require web API key")
	}
	cfg = config{YooKassaReturnURL: "https://api.neriva.ru/payment/success"}
	if got := cfg.rollyPayWebReturnURL(); got != "https://neriva.ru/app?payment=success&provider=rollypay" {
		t.Fatalf("RollyPay web return URL = %q, want RollyPay app success URL", got)
	}
	cfg.WebPaymentReturnURL = "https://neriva.ru/app?payment=success"
	if got := cfg.rollyPayWebReturnURL(); got != cfg.WebPaymentReturnURL {
		t.Fatalf("RollyPay web return URL should honor WEB_PAYMENT_RETURN_URL: got %q", got)
	}
}

func TestRollyPayClientUsesDocumentedPaymentContract(t *testing.T) {
	var gotPath string
	var gotAPIKey string
	var gotNonce string
	var gotPayload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAPIKey = r.Header.Get("X-API-Key")
		gotNonce = r.Header.Get("X-Nonce")
		if r.Header.Get("Authorization") != "" {
			t.Fatalf("unexpected Authorization header %q", r.Header.Get("Authorization"))
		}
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatal(err)
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"id":      "pay_123",
			"pay_url": "https://pay.rollypay.io/pay/tok_123",
			"status":  "created",
		})
	}))
	defer server.Close()

	client := newRollyPayClient("terminal-123", "api-key-123", server.URL, "", server.Client())
	payment, err := client.createPremiumPayment(context.Background(), userState{TelegramID: 42, FirstName: "Ada", InterfaceLanguage: "en"}, premiumPlan{
		Product:  premiumMonthlyProduct,
		Title:    "Premium",
		RubPrice: 300,
	}, rollyPayChannelTelegram, "https://example.test/success", "https://example.test/fail", time.Unix(1700000000, 0))
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/v1/payments" {
		t.Fatalf("path = %q, want /api/v1/payments", gotPath)
	}
	if gotAPIKey != "api-key-123" {
		t.Fatalf("X-API-Key = %q", gotAPIKey)
	}
	if gotNonce == "" {
		t.Fatal("expected X-Nonce header")
	}
	for key, want := range map[string]any{
		"amount":               "300.00",
		"payment_currency":     "RUB",
		"order_id":             "premium_30d_42_1700000000",
		"terminal_id":          "terminal-123",
		"customer_id":          "42",
		"success_redirect_url": "https://example.test/success",
		"fail_redirect_url":    "https://example.test/fail",
	} {
		if gotPayload[key] != want {
			t.Fatalf("payload[%s] = %#v, want %#v; full=%#v", key, gotPayload[key], want, gotPayload)
		}
	}
	if _, ok := gotPayload["cashbox_id"]; ok {
		t.Fatalf("unexpected cashbox_id in payload: %#v", gotPayload)
	}
	if _, ok := gotPayload["payment_method"]; ok {
		t.Fatalf("unexpected fixed payment_method in payload: %#v", gotPayload)
	}
	if payment.ID != "pay_123" || payment.URL != "https://pay.rollypay.io/pay/tok_123" || payment.Status != "created" {
		t.Fatalf("unexpected payment: %+v", payment)
	}
}

func TestYooKassaClientUsesStableIdempotenceKeyForRetry(t *testing.T) {
	var gotKeys []string
	client := newYooKassaClient("shop", "secret", &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		gotKeys = append(gotKeys, req.Header.Get("Idempotence-Key"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"id":"pay_123","confirmation":{"confirmation_url":"https://pay.example/123"}}`)),
		}, nil
	})})
	plan := premiumPlan{Product: premiumMonthlyProduct, Title: "Premium", RubPrice: 300}

	for i := 0; i < 2; i++ {
		if _, err := client.createPremiumPayment(context.Background(), 42, "Ada", "en", plan, "https://example.test/success", "web", "checkout-retry-1"); err != nil {
			t.Fatal(err)
		}
	}
	if len(gotKeys) != 2 {
		t.Fatalf("captured %d idempotence keys, want 2", len(gotKeys))
	}
	if gotKeys[0] == "" || gotKeys[0] != gotKeys[1] {
		t.Fatalf("Idempotence-Key values = %#v, want stable non-empty key", gotKeys)
	}
}

func TestRollyPayClientUsesStableOrderIDForRetry(t *testing.T) {
	var orderIDs []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatal(err)
		}
		orderIDs = append(orderIDs, firstString(payload, "order_id"))
		writeJSON(w, http.StatusOK, map[string]any{
			"id":      "pay_123",
			"pay_url": "https://pay.rollypay.io/pay/tok_123",
		})
	}))
	defer server.Close()
	client := newRollyPayClient("terminal-123", "api-key-123", server.URL, "", server.Client())
	plan := premiumPlan{Product: premiumMonthlyProduct, Title: "Premium", RubPrice: 300}

	for _, now := range []time.Time{time.Unix(1700000000, 0), time.Unix(1700000030, 0)} {
		if _, err := client.createPremiumPayment(context.Background(), userState{TelegramID: 42, FirstName: "Ada"}, plan, rollyPayChannelWeb, "https://example.test/success", "https://example.test/fail", now, "checkout-retry-1"); err != nil {
			t.Fatal(err)
		}
	}
	if len(orderIDs) != 2 {
		t.Fatalf("captured %d order IDs, want 2", len(orderIDs))
	}
	if orderIDs[0] == "" || orderIDs[0] != orderIDs[1] {
		t.Fatalf("order IDs = %#v, want stable non-empty order ID", orderIDs)
	}
}

func TestRollyPayClientAvoidsDuplicateAPIPrefix(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		writeJSON(w, http.StatusOK, map[string]any{
			"payment_id": "pay_123",
			"pay_url":    "https://pay.rollypay.io/pay/tok_123",
		})
	}))
	defer server.Close()

	client := newRollyPayClient("terminal-123", "api-key-123", server.URL+"/api/v1", "", server.Client())
	if _, err := client.createPremiumPayment(context.Background(), userState{TelegramID: 42}, premiumPlan{
		Product:  premiumMonthlyProduct,
		Title:    "Premium",
		RubPrice: 300,
	}, rollyPayChannelWeb, "", "", time.Unix(1700000000, 0)); err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/v1/payments" {
		t.Fatalf("path = %q, want /api/v1/payments without duplicate prefix", gotPath)
	}
}

func TestRollyPayWebhookRequiresValidSignatureWhenSecretConfigured(t *testing.T) {
	body := []byte(`{"status":"paid","order_id":"premium_30d_123_1700000000","payment_id":"rp_123","amount":"300.00","currency":"RUB"}`)
	timestamp := "1700000100"
	signature := rollyPaySignature(body, timestamp, "signing-secret")

	req := httptest.NewRequest(http.MethodPost, "/rollypay/webhook/bot", strings.NewReader(string(body)))
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", signature)
	if !rollyPayWebhookSignatureAllowed(req, body, "signing-secret") {
		t.Fatal("expected valid signature")
	}
	req.Header.Set("X-Signature", "bad")
	if rollyPayWebhookSignatureAllowed(req, body, "signing-secret") {
		t.Fatal("expected invalid signature to be rejected")
	}
}

func TestRollyPayWebhookActivatesPremiumFromOrderID(t *testing.T) {
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.getOrCreateUser(123, "Test")
	if err != nil {
		t.Fatal(err)
	}
	b := &bot{
		cfg: config{
			MaxVoiceSeconds:        30,
			PremiumRubPrice:        300,
			PremiumStarsPrice:      150,
			PremiumYearRubPrice:    3000,
			PremiumYearStarsPrice:  1500,
			PlatinumRubPrice:       590,
			PlatinumStarsPrice:     300,
			PlatinumYearRubPrice:   5900,
			PlatinumYearStarsPrice: 3000,
			RollyPayWebhookSecret:  "secret",
		},
		store: store,
	}
	orderID := rollyPayOrderID(premiumMonthlyProduct, user.TelegramID, time.Unix(1700000000, 0))
	body := []byte(`{"status":"paid","order_id":"` + orderID + `","payment_id":"rp_123","amount":"300.00","currency":"RUB"}`)
	timestamp := "1700000100"
	req := httptest.NewRequest(http.MethodPost, "/rollypay/webhook/bot", strings.NewReader(string(body)))
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", rollyPaySignature(body, timestamp, "secret"))
	rec := httptest.NewRecorder()

	b.handleRollyPayWebhook(rec, req, rollyPayChannelTelegram)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	refreshed, err := store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		t.Fatal(err)
	}
	if !refreshed.isPremium(time.Now()) || refreshed.Plan != "premium" {
		t.Fatalf("premium not activated: %+v", refreshed)
	}
	if refreshed.LastPaymentChargeID != "rollypay:rp_123" {
		t.Fatalf("charge id = %q", refreshed.LastPaymentChargeID)
	}
}

func TestRollyPayWebhookActivatesPremiumFromMetadataWhenOrderIDMissing(t *testing.T) {
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatal(err)
	}
	user, err := store.getOrCreateUser(123, "Test")
	if err != nil {
		t.Fatal(err)
	}
	b := &bot{
		cfg: config{
			MaxVoiceSeconds:        30,
			PremiumRubPrice:        300,
			PremiumStarsPrice:      150,
			PremiumYearRubPrice:    3000,
			PremiumYearStarsPrice:  1500,
			PlatinumRubPrice:       590,
			PlatinumStarsPrice:     300,
			PlatinumYearRubPrice:   5900,
			PlatinumYearStarsPrice: 3000,
			RollyPayWebhookSecret:  "secret",
		},
		store: store,
	}
	body := []byte(`{"status":"paid","payment_id":"rp_meta","amount":"300.00","metadata":{"telegram_id":"123","product":"premium_30d","channel":"site"}}`)
	timestamp := "1700000100"
	req := httptest.NewRequest(http.MethodPost, "/rollypay/webhook/site", strings.NewReader(string(body)))
	req.Header.Set("X-Timestamp", timestamp)
	req.Header.Set("X-Signature", rollyPaySignature(body, timestamp, "secret"))
	rec := httptest.NewRecorder()

	b.handleRollyPayWebhook(rec, req, rollyPayChannelWeb)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	refreshed, err := store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		t.Fatal(err)
	}
	if !refreshed.isPremium(time.Now()) || refreshed.Plan != "premium" {
		t.Fatalf("premium not activated from metadata: %+v", refreshed)
	}
	if refreshed.LastPaymentChargeID != "rollypay:rp_meta" {
		t.Fatalf("charge id = %q", refreshed.LastPaymentChargeID)
	}
}

func rollyPaySignature(body []byte, timestamp string, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp))
	mac.Write([]byte("."))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
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
