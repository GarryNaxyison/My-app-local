package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestWebGeneratedVisualAssets(t *testing.T) {
	api := newWebAPI(config{WebAPISessionSecret: "test-session-secret"}, nil)
	mux := http.NewServeMux()
	api.register(mux)

	for _, target := range []string{"/assets/header-home.png", "/app/assets/header-home.png", "/assets/icon-lesson.png", "/app/assets/icon-lesson.png", "/assets/brand-assets-manifest.json"} {
		request := httptest.NewRequest(http.MethodGet, target, nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("GET %s returned %d", target, recorder.Code)
		}
		if recorder.Body.Len() == 0 {
			t.Fatalf("GET %s returned empty body", target)
		}
	}

	request := httptest.NewRequest(http.MethodHead, "/assets/header-home.png", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("HEAD /assets/header-home.png returned %d", recorder.Code)
	}
	if recorder.Body.Len() != 0 {
		t.Fatal("HEAD /assets/header-home.png should not write a body")
	}

	request = httptest.NewRequest(http.MethodGet, "/assets/../web_api.go", nil)
	recorder = httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code == http.StatusOK {
		body, _ := io.ReadAll(recorder.Body)
		t.Fatalf("expected traversal request to be rejected, got %d: %s", recorder.Code, string(body))
	}
}

func TestClientRateKeyIgnoresForwardedForFromUntrustedRemote(t *testing.T) {
	api := newWebAPI(config{WebAPISessionSecret: "test-session-secret"}, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	request.RemoteAddr = "203.0.113.10:43124"
	request.Header.Set("X-Forwarded-For", "198.51.100.99")
	request.Header.Set("X-Real-IP", "198.51.100.100")

	if got := api.clientRateKey(request); got != "203.0.113.10" {
		t.Fatalf("clientRateKey = %q, want remote address", got)
	}
}

func TestClientRateKeyUsesForwardedForFromTrustedProxy(t *testing.T) {
	api := newWebAPI(config{
		WebAPISessionSecret: "test-session-secret",
		TrustedProxyCIDRs:   []string{"10.0.0.0/8"},
	}, nil)
	request := httptest.NewRequest(http.MethodGet, "/api/session", nil)
	request.RemoteAddr = "10.20.30.40:43124"
	request.Header.Set("X-Forwarded-For", "198.51.100.99, 10.20.30.40")

	if got := api.clientRateKey(request); got != "198.51.100.99" {
		t.Fatalf("clientRateKey = %q, want forwarded client IP", got)
	}
}

func TestWebhookHTTPServerTimeouts(t *testing.T) {
	cfg := config{WebhookListenAddr: ":0", WebAPISessionSecret: "test-session-secret"}
	server := newWebhookHTTPServer(cfg, &bot{cfg: cfg})

	if server.ReadHeaderTimeout <= 0 {
		t.Fatal("ReadHeaderTimeout must be configured")
	}
	if server.ReadTimeout <= 0 {
		t.Fatal("ReadTimeout must be configured")
	}
	if server.WriteTimeout <= 0 {
		t.Fatal("WriteTimeout must be configured")
	}
	if server.IdleTimeout <= 0 {
		t.Fatal("IdleTimeout must be configured")
	}
}

func TestWebAppShellAndPWAEntrypointsAreNoCache(t *testing.T) {
	api := newWebAPI(config{WebAPISessionSecret: "test-session-secret"}, nil)
	mux := http.NewServeMux()
	api.register(mux)

	for _, target := range []string{"/app", "/app/manifest.webmanifest", "/app/offline-deck-sw.js"} {
		request := httptest.NewRequest(http.MethodGet, target, nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("GET %s returned %d", target, recorder.Code)
		}
		cacheControl := recorder.Header().Get("Cache-Control")
		if !strings.Contains(cacheControl, "no-store") {
			t.Fatalf("GET %s Cache-Control = %q, want no-store", target, cacheControl)
		}
	}
}

func TestWebUserDTOIncludesTelegramAccountIdentity(t *testing.T) {
	api := newWebAPI(config{}, nil)
	dto := api.userDTO(userState{
		TelegramID:        123456789,
		FirstName:         "demo_user",
		InterfaceLanguage: "ru",
		LearningLanguage:  "en",
		Level:             "A2",
	})
	if linked, _ := dto["telegram_linked"].(bool); !linked {
		t.Fatalf("telegram_linked = %#v, want true", dto["telegram_linked"])
	}
	account, ok := dto["telegram_account"].(map[string]any)
	if !ok {
		t.Fatalf("telegram_account = %#v, want object", dto["telegram_account"])
	}
	if account["id"] != int64(123456789) {
		t.Fatalf("telegram account id = %#v", account["id"])
	}
	if account["name"] != "demo_user" {
		t.Fatalf("telegram account name = %#v", account["name"])
	}
	if label, _ := account["label"].(string); !strings.Contains(label, "demo_user") || !strings.Contains(label, "123456789") {
		t.Fatalf("telegram account label = %#v", account["label"])
	}
}

func TestWebUserDTOIncludesActiveLessonForReload(t *testing.T) {
	api := newWebAPI(config{}, nil)
	dto := api.userDTO(userState{
		TelegramID:        123456789,
		FirstName:         "demo_user",
		InterfaceLanguage: "ru",
		LearningLanguage:  "en",
		Level:             "A2",
		Mode:              "lesson",
		LastLessonPrompt:  "Answer this active lesson.",
	})
	if dto["active_lesson_prompt"] != "Answer this active lesson." {
		t.Fatalf("active_lesson_prompt = %#v", dto["active_lesson_prompt"])
	}
	if instruction, _ := dto["active_lesson_instruction"].(string); strings.TrimSpace(instruction) == "" {
		t.Fatalf("active_lesson_instruction missing in dto: %#v", dto)
	}
}

func TestWebSettingsSavesLearningFocus(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	focus := "business calls and meetings"

	body := requestJSON(t, api, cookie, http.MethodPost, "/api/settings", map[string]any{
		"interface_language": "ru",
		"learning_language":  "en",
		"level":              "A2",
		"learning_focus":     focus,
	})
	user, ok := body["user"].(map[string]any)
	if !ok {
		t.Fatalf("missing user in settings response: %#v", body)
	}
	if user["learning_focus"] != focus {
		t.Fatalf("settings response learning_focus = %#v, want %q", user["learning_focus"], focus)
	}
	stored, _, err := store.getUser(-42)
	if err != nil {
		t.Fatal(err)
	}
	if stored.LearningFocus != focus {
		t.Fatalf("stored LearningFocus = %q, want %q", stored.LearningFocus, focus)
	}
}

func TestWebTelegramCodeFlowStartsBotLinkAndVerifiesSiteCode(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	start := requestJSON(t, api, cookie, http.MethodPost, "/api/auth/telegram/start", map[string]any{})
	token, _ := start["token"].(string)
	botURL, _ := start["bot_url"].(string)
	if token == "" || !strings.Contains(botURL, "start=web_") {
		t.Fatalf("unexpected telegram start response: %#v", start)
	}
	request, err := api.bot.webAuth.beginTelegram(token, telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "demo"})
	if err != nil {
		t.Fatalf("begin telegram: %v", err)
	}
	linked := requestJSON(t, api, cookie, http.MethodPost, "/api/auth/telegram/status", map[string]any{
		"token": token,
		"code":  request.Code,
	})
	if authenticated, _ := linked["authenticated"].(bool); !authenticated {
		t.Fatalf("expected verified telegram code to return session, got %#v", linked)
	}
	user, ok := linked["user"].(map[string]any)
	if !ok {
		t.Fatalf("expected linked user in session: %#v", linked)
	}
	if tg, _ := user["telegram_linked"].(bool); !tg {
		t.Fatalf("expected telegram_linked in session, got %#v", user)
	}
}

func TestWebTelegramCodeFlowReturnsLocalizedCodeWhenTelegramAlreadyLinked(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	linkedAccountHash, err := hashWebPassword("linked-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-99, "linked_user", linkedAccountHash); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.authenticateTelegramWebAccount(-99, telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "demo"}); err != nil {
		t.Fatal(err)
	}

	start := requestJSON(t, api, cookie, http.MethodPost, "/api/auth/telegram/start", map[string]any{})
	token, _ := start["token"].(string)
	request, err := api.bot.webAuth.beginTelegram(token, telegramWebProfile{ID: 1001, FirstName: "Demo", Username: "demo"})
	if err != nil {
		t.Fatalf("begin telegram: %v", err)
	}

	status, payload := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/auth/telegram/status", map[string]any{
		"token": token,
		"code":  request.Code,
	})
	if status != http.StatusConflict {
		t.Fatalf("telegram status = %d payload %#v, want conflict", status, payload)
	}
	apiError, _ := payload["error"].(map[string]any)
	if code, _ := apiError["code"].(string); code != "telegram_already_linked" {
		t.Fatalf("error code = %#v payload %#v, want telegram_already_linked", apiError["code"], payload)
	}
}

func TestWebAuthRegisterDotsAndLocalizedErrorCodes(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)

	status, invalid := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/auth/register", map[string]any{
		"login":            ".friend",
		"password":         "strong-password",
		"password_confirm": "strong-password",
	})
	if status != http.StatusBadRequest {
		t.Fatalf("invalid dotted login returned %d: %#v", status, invalid)
	}
	apiError, ok := invalid["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object, got %#v", invalid)
	}
	if apiError["code"] != "invalid_login_start" {
		t.Fatalf("error code = %#v, want invalid_login_start", apiError["code"])
	}
	if message, _ := apiError["message"].(string); strings.Contains(message, "Раздел") || strings.Contains(message, "Section") || message == "" {
		t.Fatalf("unexpected auth error message: %q", message)
	}

	status, valid := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/auth/register", map[string]any{
		"login":              "friend.name",
		"password":           "strong-password",
		"password_confirm":   "strong-password",
		"interface_language": "ru",
	})
	if status != http.StatusOK {
		t.Fatalf("valid dotted login returned %d: %#v", status, valid)
	}
	account, ok := valid["account"].(map[string]any)
	if !ok || account["login"] != "friend.name" {
		t.Fatalf("expected registered dotted login in account payload, got %#v", valid["account"])
	}
}

func TestWebLearningFeatureFlows(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)

	addLearnedTestWord(t, store, -42, "be")
	addLearnedTestWord(t, store, -42, "have")
	addLearnedTestWord(t, store, -42, "go")
	markTestWordMastered(t, store, -42, "be")
	if err := store.addMistakes(-42, "en", []mistakeEntry{{
		Language:    "en",
		Word:        "She go home",
		Correction:  "She goes home",
		Explanation: "Present Simple needs -s with she.",
		AddedAt:     time.Now(),
	}}); err != nil {
		t.Fatalf("addMistakes: %v", err)
	}

	vocabulary := requestJSON(t, api, cookie, http.MethodGet, "/api/vocabulary", nil)
	if total, _ := vocabulary["total"].(float64); total != 3 {
		t.Fatalf("expected all learned vocabulary words, got %#v", vocabulary)
	}
	if items, _ := vocabulary["items"].([]any); len(items) != 3 {
		t.Fatalf("expected all learned vocabulary items, got %#v", vocabulary)
	}

	wordGame := requestJSON(t, api, cookie, http.MethodPost, "/api/word-game/next", map[string]any{})
	if empty, _ := wordGame["empty"].(bool); empty {
		t.Fatal("word game should not be empty after learned words were added")
	}
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	correctWordID := strings.TrimPrefix(user.Mode, modeWebWordGamePrefix)
	if correctWordID == "" || correctWordID == user.Mode {
		t.Fatalf("expected active word game mode, got %q", user.Mode)
	}
	wordGameAnswer := requestJSON(t, api, cookie, http.MethodPost, "/api/word-game/answer", map[string]any{"answer_id": correctWordID})
	if correct, _ := wordGameAnswer["correct"].(bool); !correct {
		t.Fatalf("expected correct word game answer, got %#v", wordGameAnswer)
	}

	spelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/start", map[string]any{})
	if empty, _ := spelling["empty"].(bool); empty {
		t.Fatal("spelling should not be empty after learned words were added")
	}
	var wrongSpelling map[string]any
	for attempt := 1; attempt <= 3; attempt++ {
		wrongSpelling = requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/answer", map[string]any{"text": "wrong"})
		if correct, _ := wrongSpelling["correct"].(bool); correct {
			t.Fatalf("expected wrong spelling answer to stay incorrect, got %#v", wrongSpelling)
		}
		if hint, _ := wrongSpelling["hint"].(string); hint == "" {
			t.Fatalf("expected spelling hint, got %#v", wrongSpelling)
		}
		if attempts, _ := wrongSpelling["attempts"].(float64); int(attempts) != attempt {
			t.Fatalf("expected spelling attempt %d, got %#v", attempt, wrongSpelling)
		}
	}
	if canGiveUp, _ := wrongSpelling["can_give_up"].(bool); !canGiveUp {
		t.Fatalf("expected give-up option after 3 wrong spelling answers, got %#v", wrongSpelling)
	}
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	spellingWordID, _, ok := parseWebSpellingMode(user.Mode)
	if !ok {
		t.Fatalf("expected active spelling mode, got %q", user.Mode)
	}
	beforeGiveUpXP := user.XP
	beforeGiveUpSpelling := learnedWordSpellingCount(user, spellingWordID)
	giveUpSpelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/answer", map[string]any{"give_up": true})
	if gaveUp, _ := giveUpSpelling["gave_up"].(bool); !gaveUp {
		t.Fatalf("expected spelling give-up response, got %#v", giveUpSpelling)
	}
	if answer, _ := giveUpSpelling["correct_answer"].(string); answer == "" {
		t.Fatalf("expected give-up spelling response to include correct_answer, got %#v", giveUpSpelling)
	}
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if user.XP != beforeGiveUpXP {
		t.Fatalf("give-up changed XP: got %d want %d", user.XP, beforeGiveUpXP)
	}
	if count := learnedWordSpellingCount(user, spellingWordID); count != beforeGiveUpSpelling {
		t.Fatalf("give-up changed spelling count: got %d want %d", count, beforeGiveUpSpelling)
	}
	nextSpelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/start", map[string]any{})
	if empty, _ := nextSpelling["empty"].(bool); empty {
		t.Fatal("spelling should still have selectable words after give-up")
	}
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	spellingWordID, _, ok = parseWebSpellingMode(user.Mode)
	if !ok {
		t.Fatalf("expected active spelling mode after give-up restart, got %q", user.Mode)
	}
	giveUpWordID, _ := giveUpSpelling["word_id"].(string)
	if len(spellingWordIDs(user)) > 1 && spellingWordID == giveUpWordID {
		t.Fatalf("expected next spelling word after give-up, got the same word %q", spellingWordID)
	}
	word, ok := findVocabWord(spellingWordID)
	if !ok {
		t.Fatalf("active spelling word not found: %q", spellingWordID)
	}
	correctSpelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/answer", map[string]any{"text": word.English})
	if correct, _ := correctSpelling["correct"].(bool); !correct {
		t.Fatalf("expected correct spelling answer, got %#v", correctSpelling)
	}
	if answer, _ := correctSpelling["correct_answer"].(string); answer != word.English {
		t.Fatalf("expected correct_answer %q, got %#v", word.English, correctSpelling)
	}

	mistakes := requestJSON(t, api, cookie, http.MethodGet, "/api/mistakes", nil)
	if total, _ := mistakes["total"].(float64); total != 2 {
		t.Fatalf("expected initial mistake plus spelling mistake, got %#v", mistakes)
	}
	if items, _ := mistakes["items"].([]any); len(items) != 2 {
		t.Fatalf("expected /api/mistakes to include all persisted mistakes, got %#v", mistakes)
	}
	startMistake := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/start", map[string]any{"index": 0})
	if empty, _ := startMistake["empty"].(bool); empty {
		t.Fatalf("expected mistake practice item, got %#v", startMistake)
	}
	wrongMistake := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": "She went home"})
	if correct, _ := wrongMistake["correct"].(bool); correct {
		t.Fatalf("expected wrong mistake practice answer, got %#v", wrongMistake)
	}
	correctMistake := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": "She goes home"})
	if correct, _ := correctMistake["correct"].(bool); !correct {
		t.Fatalf("expected correct mistake practice answer, got %#v", correctMistake)
	}
	afterMistakes := requestJSON(t, api, cookie, http.MethodGet, "/api/mistakes", nil)
	if total, _ := afterMistakes["total"].(float64); total != 1 {
		t.Fatalf("expected repaired mistake removed while spelling mistake remains, got %#v", afterMistakes)
	}
}

func TestWebLearningActionsAwardXPAndReturnUpdatedUser(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	responses := []string{
		"Good correction.\n" + mistakesSentinel + "\n[]",
		"Good practice.\n" + mistakesSentinel + "\n[]",
	}
	call := 0
	api.bot.openrouter = newOpenRouterClient("test-key", "xp-model", "http://localhost", "test", &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://openrouter.ai/api/v1/chat/completions" {
				t.Fatalf("unexpected OpenRouter URL: %s", req.URL.String())
			}
			if call >= len(responses) {
				t.Fatalf("unexpected extra OpenRouter call %d", call+1)
			}
			body, err := json.Marshal(map[string]any{
				"choices": []map[string]any{{"message": map[string]any{"content": responses[call]}}},
			})
			if err != nil {
				t.Fatal(err)
			}
			call++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		}),
	})

	currentXP := func() int {
		t.Helper()
		user, err := store.getOrCreateUser(-42, "tester")
		if err != nil {
			t.Fatal(err)
		}
		return user.XP
	}
	assertXPDelta := func(label string, before int, delta int, payload map[string]any) {
		t.Helper()
		want := before + delta
		if got := currentXP(); got != want {
			t.Fatalf("%s stored XP = %d, want %d", label, got, want)
		}
		userPayload, ok := payload["user"].(map[string]any)
		if !ok {
			t.Fatalf("%s response missing user payload: %#v", label, payload)
		}
		if got := int(userPayload["xp"].(float64)); got != want {
			t.Fatalf("%s response user.xp = %d, want %d", label, got, want)
		}
	}

	if err := store.saveLesson(-42, "Say that you need help."); err != nil {
		t.Fatalf("save lesson: %v", err)
	}
	before := currentXP()
	lesson := requestJSON(t, api, cookie, http.MethodPost, "/api/lesson/answer", map[string]any{"text": "I need help"})
	assertXPDelta("lesson", before, 20, lesson)

	before = currentXP()
	practice := requestJSON(t, api, cookie, http.MethodPost, "/api/practice", map[string]any{"text": "Can you repeat?"})
	assertXPDelta("practice", before, 5, practice)

	_ = requestJSON(t, api, cookie, http.MethodPost, "/api/words/next", map[string]any{})
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	wordID := strings.TrimPrefix(user.Mode, modeWebWordPrefix)
	if wordID == "" || wordID == user.Mode {
		t.Fatalf("expected active word mode, got %q", user.Mode)
	}
	before = currentXP()
	word := requestJSON(t, api, cookie, http.MethodPost, "/api/words/answer", map[string]any{"answer_id": wordID})
	assertXPDelta("word", before, 15, word)

	_ = requestJSON(t, api, cookie, http.MethodPost, "/api/word-game/next", map[string]any{})
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	reviewWordID := strings.TrimPrefix(user.Mode, modeWebWordGamePrefix)
	if reviewWordID == "" || reviewWordID == user.Mode {
		t.Fatalf("expected active review mode, got %q", user.Mode)
	}
	before = currentXP()
	review := requestJSON(t, api, cookie, http.MethodPost, "/api/word-game/answer", map[string]any{"answer_id": reviewWordID})
	assertXPDelta("review", before, 10, review)

	_ = requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/start", map[string]any{})
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	spellingWordID, _, ok := parseWebSpellingMode(user.Mode)
	if !ok {
		t.Fatalf("expected active spelling mode, got %q", user.Mode)
	}
	spellingWord, ok := learnedWordByID(user, spellingWordID)
	if !ok {
		t.Fatalf("active spelling word not found: %q", spellingWordID)
	}
	before = currentXP()
	spelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/answer", map[string]any{"text": spellingWord.English})
	assertXPDelta("spelling", before, 7, spelling)

	if err := store.addMistakes(-42, "en", []mistakeEntry{{
		Language:    "en",
		Word:        "She go home",
		Correction:  "She goes home",
		Explanation: "Present Simple needs -s with she.",
		AddedAt:     time.Now(),
	}}); err != nil {
		t.Fatalf("addMistakes: %v", err)
	}
	_ = requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/start", map[string]any{"index": 0})
	before = currentXP()
	mistake := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": "She goes home"})
	assertXPDelta("mistake", before, 8, mistake)
}

func TestWebWordGameNextAutoTranslatesMissingPrompt(t *testing.T) {
	configureWebMissingPromptVocabularyForTest(t)
	seedVocabularyRandomForTest(t, 1)
	api, store, cookie := newTestWebAPI(t)
	configureWebVocabularyOpenRouterForTest(t, api, "вращающийся в две стороны", nil)
	addLearnedTestWord(t, store, -42, "en:birotate")

	body := requestJSON(t, api, cookie, http.MethodPost, "/api/word-game/next", map[string]any{})
	if body["prompt"] != "вращающийся в две стороны" {
		t.Fatalf("expected OpenRouter translation prompt, got %#v", body)
	}
	if body["message"] == ui(userState{InterfaceLanguage: "ru"}).LearnWords {
		t.Fatalf("word-game prompt leaked empty-state copy: %#v", body)
	}
	if stored, ok, err := sqliteVocabularyAITranslationGet("en:birotate", "ru"); err != nil || !ok || stored != "вращающийся в две стороны" {
		t.Fatalf("stored AI translation = %q ok=%v err=%v", stored, ok, err)
	}
}

func TestWebReviewAndSpellingUsePersistedAITranslationPrompts(t *testing.T) {
	configureWebMissingPromptVocabularyForTest(t)
	seedVocabularyRandomForTest(t, 1)
	api, store, cookie := newTestWebAPI(t)
	openRouterCalls := 0
	configureWebVocabularyOpenRouterForTest(t, api, "вращающийся в две стороны", &openRouterCalls)
	addLearnedTestWord(t, store, -42, "en:birotate")

	next := requestJSON(t, api, cookie, http.MethodPost, "/api/word-game/next", map[string]any{})
	if next["prompt"] != "вращающийся в две стороны" {
		t.Fatalf("word-game next prompt = %#v", next)
	}
	wrongReview := requestJSON(t, api, cookie, http.MethodPost, "/api/word-game/answer", map[string]any{"answer_id": "en:not-the-answer"})
	if wrongReview["prompt"] != "вращающийся в две стороны" {
		t.Fatalf("word-game wrong prompt = %#v", wrongReview)
	}
	spelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/start", map[string]any{})
	if spelling["prompt"] != "вращающийся в две стороны" {
		t.Fatalf("spelling start prompt = %#v", spelling)
	}
	wrongSpelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/answer", map[string]any{"text": "wrong"})
	if wrongSpelling["prompt"] != "вращающийся в две стороны" {
		t.Fatalf("spelling wrong prompt = %#v", wrongSpelling)
	}
	giveUpSpelling := requestJSON(t, api, cookie, http.MethodPost, "/api/spelling/answer", map[string]any{"give_up": true})
	if giveUpSpelling["translation"] != "вращающийся в две стороны" {
		t.Fatalf("spelling give-up translation = %#v", giveUpSpelling)
	}
	if openRouterCalls != 1 {
		t.Fatalf("OpenRouter calls = %d, want one persisted translation reused by follow-up paths", openRouterCalls)
	}
}

func TestWebWordGameNextFailsWhenAITranslationCannotPersist(t *testing.T) {
	configureWebMissingPromptVocabularyForTest(t)
	api, store, cookie := newTestWebAPI(t)
	openRouterCalls := 0
	configureWebVocabularyOpenRouterForTest(t, api, "вращающийся в две стороны", &openRouterCalls)
	addLearnedTestWord(t, store, -42, "en:birotate")
	db := currentSQLiteVocabularyDB()
	if db == nil {
		t.Fatal("expected configured SQLite vocabulary DB")
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA query_only = ON`); err != nil {
		t.Fatal(err)
	}

	status, body := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/word-game/next", map[string]any{})
	if status != http.StatusInternalServerError {
		t.Fatalf("word-game next status = %d body=%#v, want 500", status, body)
	}
	if openRouterCalls != 1 {
		t.Fatalf("OpenRouter calls = %d, want attempted translation before persistence failure", openRouterCalls)
	}
	if stored, ok, err := sqliteVocabularyAITranslationGet("en:birotate", "ru"); err != nil || ok || stored != "" {
		t.Fatalf("AI translation should not be partially persisted: value=%q ok=%v err=%v", stored, ok, err)
	}
}

func configureWebMissingPromptVocabularyForTest(t *testing.T) {
	t.Helper()
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:            "en:birotate",
			Language:      "en",
			English:       "birotate",
			Russian:       "birotate",
			Translations:  map[string]string{"en": "birotate"},
			Level:         "A1",
			FrequencyRank: 1,
		},
		{
			ID:            "en:user",
			Language:      "en",
			English:       "user",
			Russian:       "пользователь",
			Translations:  map[string]string{"ru": "пользователь"},
			Level:         "A1",
			FrequencyRank: 2,
		},
		{
			ID:            "en:birthday",
			Language:      "en",
			English:       "birthday",
			Russian:       "день рождения",
			Translations:  map[string]string{"ru": "день рождения"},
			Level:         "A1",
			FrequencyRank: 3,
		},
		{
			ID:            "en:women's",
			Language:      "en",
			English:       "women's",
			Russian:       "женский",
			Translations:  map[string]string{"ru": "женский"},
			Level:         "A1",
			FrequencyRank: 4,
		},
	})
}

func configureWebVocabularyOpenRouterForTest(t *testing.T, api *webAPI, translation string, calls *int) {
	t.Helper()
	api.cfg.OpenRouterVocabularyModel = "vocab-test-model"
	api.bot.cfg.OpenRouterVocabularyModel = "vocab-test-model"
	api.bot.openrouter = newOpenRouterClient("test-key", "fallback-model", "http://localhost", "test", &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://openrouter.ai/api/v1/chat/completions" {
				t.Fatalf("unexpected OpenRouter URL %s", req.URL.String())
			}
			var payload map[string]any
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Fatalf("decode OpenRouter payload: %v", err)
			}
			if payload["model"] != "vocab-test-model" {
				t.Fatalf("unexpected vocabulary model %#v", payload["model"])
			}
			if calls != nil && openRouterPayloadIsVocabularyTranslation(payload) {
				(*calls)++
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(`{"choices":[{"message":{"content":` + strconv.Quote(translation) + `}}]}`)),
			}, nil
		}),
	})
}

func openRouterPayloadIsVocabularyTranslation(payload map[string]any) bool {
	switch value := payload["max_tokens"].(type) {
	case float64:
		return int(value) == 160
	case int:
		return value == 160
	}
	return false
}

func TestWebLessonAnswerCompletesActiveLessonAndRejectsReplay(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	responses := []string{
		"First correction.\n" + mistakesSentinel + "\n[]",
	}
	call := 0
	api.bot.openrouter = newOpenRouterClient("test-key", "xp-once-model", "http://localhost", "test", &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if call >= len(responses) {
				t.Fatalf("unexpected extra OpenRouter call %d", call+1)
			}
			body, err := json.Marshal(map[string]any{
				"choices": []map[string]any{{"message": map[string]any{"content": responses[call]}}},
			})
			if err != nil {
				t.Fatal(err)
			}
			call++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		}),
	})

	if err := store.saveLesson(-42, "Answer this lesson once."); err != nil {
		t.Fatalf("save lesson: %v", err)
	}
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	before := user.XP

	first := requestJSON(t, api, cookie, http.MethodPost, "/api/lesson/answer", map[string]any{"text": "My first answer"})
	if got := first["feedback"]; got != "First correction." {
		t.Fatalf("first feedback = %#v", got)
	}
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := user.XP, before+20; got != want {
		t.Fatalf("XP after first answer = %d, want %d", got, want)
	}
	if strings.TrimSpace(user.LastLessonPrompt) != "" {
		t.Fatalf("LastLessonPrompt after first answer = %q, want cleared completed lesson", user.LastLessonPrompt)
	}

	status, second := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/lesson/answer", map[string]any{"text": "My second answer"})
	if status != http.StatusConflict {
		t.Fatalf("second answer status = %d body=%#v, want 409", status, second)
	}
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := user.XP, before+20; got != want {
		t.Fatalf("XP after second answer = %d, want %d", got, want)
	}
	if call != 1 {
		t.Fatalf("expected 1 OpenRouter call, got %d", call)
	}
}

func TestWebDailyClaimInsideTwentyFourHoursDoesNotAwardXP(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	beforeUser, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}

	first := requestJSON(t, api, cookie, http.MethodPost, "/api/daily/claim", map[string]any{"date": "2026-06-12"})
	if first["claimed"] != true {
		t.Fatalf("first daily claim response = %#v, want claimed=true", first)
	}
	afterFirst, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if afterFirst.XP != beforeUser.XP+25 {
		t.Fatalf("XP after first claim = %d, want %d", afterFirst.XP, beforeUser.XP+25)
	}

	second := requestJSON(t, api, cookie, http.MethodPost, "/api/daily/claim", map[string]any{"date": "2026-06-12"})
	if second["claimed"] != false {
		t.Fatalf("second daily claim response = %#v, want claimed=false", second)
	}
	afterSecond, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if afterSecond.XP != afterFirst.XP {
		t.Fatalf("XP after rejected claim = %d, want unchanged %d", afterSecond.XP, afterFirst.XP)
	}
}

func TestWebTranslatorSpeechRequiresPaidAudioBudget(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	api.cfg.OpenRouterTTSModel = "google/gemini-3.1-flash-tts-preview"
	api.bot.cfg = api.cfg
	api.bot.openrouter = newOpenRouterClient("test-key", "google/gemini-3.1-flash", "http://localhost", "test", &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"audio/mpeg"}},
			Body:       io.NopCloser(bytes.NewBufferString("audio")),
		}, nil
	})})

	status, response := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/tools/translator-speech", map[string]any{"text": "Hello", "target_language": "en"})
	if status != http.StatusPaymentRequired {
		t.Fatalf("free translator speech status = %d response=%#v; want %d", status, response, http.StatusPaymentRequired)
	}

	if _, err := api.bot.store.extendPremium(-42, "test-premium-audio", 24*time.Hour, "premium"); err != nil {
		t.Fatal(err)
	}
	body, code := requestRaw(t, api, cookie, http.MethodPost, "/api/tools/translator-speech", map[string]any{"text": "Hello", "target_language": "en"})
	if code != http.StatusOK || !bytes.Contains(body, []byte("audio")) {
		t.Fatalf("premium translator speech = status %d body %q", code, string(body))
	}
	user, err := api.bot.store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if user.VoiceToday != 1 {
		t.Fatalf("premium audio budget used = %d; want 1", user.VoiceToday)
	}
}

func TestWebLearnedWordPronunciationStaysFree(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	api.cfg.OpenRouterTTSModel = "google/gemini-3.1-flash-tts-preview"
	api.bot.cfg = api.cfg
	api.bot.openrouter = newOpenRouterClient("test-key", "google/gemini-3.1-flash", "http://localhost", "test", &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"audio/mpeg"}},
			Body:       io.NopCloser(bytes.NewBufferString("word-audio")),
		}, nil
	})})
	addLearnedTestWord(t, store, -42, "be")
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if len(user.LearnedWords) == 0 {
		t.Fatal("expected learned word")
	}

	body, code := requestRaw(t, api, cookie, http.MethodPost, "/api/words/pronunciation", map[string]any{"word_id": user.LearnedWords[0].ID})
	if code != http.StatusOK || !bytes.Contains(body, []byte("word-audio")) {
		t.Fatalf("free learned word pronunciation = status %d body %q", code, string(body))
	}
	user, err = store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if user.VoiceToday != 0 {
		t.Fatalf("learned word pronunciation consumed audio budget: %d", user.VoiceToday)
	}
}

func TestWebLearnedWordPronunciationAcceptsLegacyWordID(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	api.cfg.OpenRouterTTSModel = "google/gemini-3.1-flash-tts-preview"
	api.bot.cfg = api.cfg
	api.bot.openrouter = newOpenRouterClient("test-key", "google/gemini-3.1-flash", "http://localhost", "test", &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"audio/mpeg"}},
			Body:       io.NopCloser(bytes.NewBufferString("legacy-word-audio")),
		}, nil
	})})
	addLearnedTestWord(t, store, -42, "be")
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if len(user.LearnedWords) == 0 {
		t.Fatal("expected learned word")
	}
	legacyID := legacyVocabID(user.LearnedWords[0].ID)
	if legacyID == user.LearnedWords[0].ID {
		t.Fatalf("test word %q does not have a legacy id variant", user.LearnedWords[0].ID)
	}

	body, code := requestRaw(t, api, cookie, http.MethodPost, "/api/words/pronunciation", map[string]any{"word_id": legacyID})
	if code != http.StatusOK || !bytes.Contains(body, []byte("legacy-word-audio")) {
		t.Fatalf("legacy learned word pronunciation = status %d body %q", code, string(body))
	}
}

func TestWriteAudioResponsesUseDetectedWAVMetadata(t *testing.T) {
	wav := wrapPCM16LEAsWAV([]byte{0x01, 0x02}, 24000, 1)
	wordRecorder := httptest.NewRecorder()
	writePronunciationAudio(wordRecorder, vocabWord{ID: "en:water"}, wav)
	if got := wordRecorder.Header().Get("Content-Type"); got != "audio/wav" {
		t.Fatalf("word pronunciation content type = %q, want audio/wav", got)
	}
	if got := wordRecorder.Header().Get("Content-Disposition"); !strings.Contains(got, `filename="water.wav"`) {
		t.Fatalf("word pronunciation filename = %q, want water.wav", got)
	}

	toolRecorder := httptest.NewRecorder()
	writeToolAudio(toolRecorder, "translation.mp3", wav)
	if got := toolRecorder.Header().Get("Content-Type"); got != "audio/wav" {
		t.Fatalf("tool audio content type = %q, want audio/wav", got)
	}
	if got := toolRecorder.Header().Get("Content-Disposition"); !strings.Contains(got, `filename="translation.wav"`) {
		t.Fatalf("tool audio filename = %q, want translation.wav", got)
	}
}

func TestPronunciationCheckRequiresPremiumBeforeProviderCall(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	api.cfg.OpenRouterSTTModel = "openai/gpt-4o-transcribe"
	api.cfg.OpenRouterTTSModel = "google/gemini-3.1-flash-tts-preview"
	api.cfg.OpenRouterPronunciationAudioModel = "google/gemini-2.5-pro"
	api.bot.cfg = api.cfg
	providerCalled := false
	api.bot.openrouter = newOpenRouterClient("test-key", "google/gemini-3.1-flash", "http://localhost", "test", &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		providerCalled = true
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBufferString(`{"choices":[{"message":{"content":"{}"}}]}`)),
		}, nil
	})})

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("target", "Good morning"); err != nil {
		t.Fatal(err)
	}
	part, err := writer.CreateFormFile("voice", "voice.webm")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte("fake webm audio")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/pronunciation/check", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.AddCookie(cookie)
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusPaymentRequired {
		t.Fatalf("pronunciation check status = %d body=%s, want 402", recorder.Code, recorder.Body.String())
	}
	if providerCalled {
		t.Fatal("pronunciation provider was called for a free user")
	}
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if user.VoiceToday != 0 {
		t.Fatalf("free pronunciation check consumed voice budget: %d", user.VoiceToday)
	}
}

func TestWebPhrasebookPersistsInStoreAndSession(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	saved := requestJSON(t, api, cookie, http.MethodPost, "/api/phrasebook", map[string]any{
		"id":          "phrase-test-1",
		"phrase":      "Could you say that again, please?",
		"translation": "Ask politely to repeat.",
		"source":      "lesson",
		"language":    "en",
	})
	items, _ := saved["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected saved phrasebook item, got %#v", saved)
	}
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if len(user.Phrasebook) != 1 || user.Phrasebook[0].Phrase != "Could you say that again, please?" {
		t.Fatalf("phrasebook was not persisted in store: %#v", user.Phrasebook)
	}
	session := requestJSON(t, api, cookie, http.MethodGet, "/api/session", nil)
	sessionUser, _ := session["user"].(map[string]any)
	sessionPhrases, _ := sessionUser["phrasebook"].([]any)
	if len(sessionPhrases) != 1 {
		t.Fatalf("session did not include phrasebook entries: %#v", sessionUser)
	}
	removed := requestJSON(t, api, cookie, http.MethodDelete, "/api/phrasebook?id=phrase-test-1", nil)
	removedItems, _ := removed["items"].([]any)
	if len(removedItems) != 0 {
		t.Fatalf("expected phrasebook item to be deleted, got %#v", removed)
	}
}

func TestWebSessionPersistsHabitLoginInDatabase(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	session := requestJSON(t, api, cookie, http.MethodGet, "/api/session", nil)
	user, _ := session["user"].(map[string]any)
	habit, _ := user["habit_log"].(map[string]any)
	today := localDateForUser(userState{ReminderUTCOffset: 180}, time.Now().UTC())
	day, _ := habit[today].(map[string]any)
	if day == nil || day["login"] != true {
		t.Fatalf("expected session to include persisted login day %s, got %#v", today, habit)
	}
	stored, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if stored.HabitLog[today].Login != true {
		t.Fatalf("expected login day persisted in sqlite store, got %#v", stored.HabitLog)
	}
}

func TestWebNavigationLayoutPersistsInDatabaseAndSession(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	saved := requestJSON(t, api, cookie, http.MethodPost, "/api/navigation-layout", map[string]any{
		"function_ribbon": []string{"lesson", "home", "practice", "lesson", "bad-view"},
		"mobile_pinned":   []string{"home", "mistakes", "spelling"},
		"mobile_more":     []string{"settings", "premium"},
		"mobile_rail":     []string{"lesson", "practice", "mistakes"},
	})
	user, _ := saved["user"].(map[string]any)
	layout, _ := user["navigation_layout"].(map[string]any)
	ribbon, _ := layout["function_ribbon"].([]any)
	if got := strings.Join(anyStrings(ribbon), ","); got != "lesson,home,practice" {
		t.Fatalf("unexpected normalized function ribbon: %#v", layout)
	}
	stored, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(stored.NavigationLayout.MobileRail, ","); got != "lesson,practice,mistakes" {
		t.Fatalf("navigation layout was not persisted: %#v", stored.NavigationLayout)
	}
	session := requestJSON(t, api, cookie, http.MethodGet, "/api/session", nil)
	sessionUser, _ := session["user"].(map[string]any)
	sessionLayout, _ := sessionUser["navigation_layout"].(map[string]any)
	rail, _ := sessionLayout["mobile_rail"].([]any)
	if got := strings.Join(anyStrings(rail), ","); got != "lesson,practice,mistakes" {
		t.Fatalf("session did not include persisted navigation layout: %#v", sessionLayout)
	}
}

func TestWebNavigationLayoutPatchPreservesSavedMobileOrder(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	_ = requestJSON(t, api, cookie, http.MethodPost, "/api/navigation-layout", map[string]any{
		"function_ribbon": []string{"home", "lesson"},
		"mobile_rail":     []string{"settings", "home", "tutor"},
	})
	_ = requestJSON(t, api, cookie, http.MethodPost, "/api/navigation-layout", map[string]any{
		"function_ribbon": []string{"lesson", "home"},
	})
	stored, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(stored.NavigationLayout.FunctionRibbon, ","); got != "lesson,home" {
		t.Fatalf("function ribbon patch was not saved: %#v", stored.NavigationLayout)
	}
	if got := strings.Join(stored.NavigationLayout.MobileRail, ","); got != "settings,home,tutor" {
		t.Fatalf("mobile rail order was lost after partial patch: %#v", stored.NavigationLayout)
	}
}

func TestWebMistakesReturnsAllItemsForClientPagination(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	entries := make([]mistakeEntry, 0, 12)
	base := time.Now().UTC().Add(-time.Hour)
	for i := 0; i < 12; i++ {
		entries = append(entries, mistakeEntry{
			Language:    "en",
			Word:        "wrong " + strconv.Itoa(i+1),
			Correction:  "correct " + strconv.Itoa(i+1),
			Explanation: "test",
			AddedAt:     base.Add(time.Duration(i) * time.Minute),
		})
	}
	if err := store.addMistakes(-42, "en", entries); err != nil {
		t.Fatalf("add mistakes: %v", err)
	}
	response := requestJSON(t, api, cookie, http.MethodGet, "/api/mistakes", nil)
	if total, _ := response["total"].(float64); total != 12 {
		t.Fatalf("expected total 12, got %#v", response)
	}
	items, _ := response["items"].([]any)
	if len(items) != 12 {
		t.Fatalf("expected all 12 items for client pagination, got %#v", response)
	}
	pageItems, _ := response["page_items"].([]any)
	if len(pageItems) != 10 {
		t.Fatalf("expected first page slice of 10, got %#v", response)
	}
	first, _ := items[0].(map[string]any)
	if first["word"] != "wrong 12" {
		t.Fatalf("expected newest mistake first, got %#v", items[0])
	}
}

func TestWebMistakeDeleteRemovesOneItem(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	entries := []mistakeEntry{
		{
			Language:    "en",
			Word:        "first wrong",
			Correction:  "first correct",
			Explanation: "grammar",
			AddedAt:     time.Now().UTC().Add(-2 * time.Minute),
		},
		{
			Language:    "en",
			Word:        "second wrong",
			Correction:  "second correct",
			Explanation: "grammar",
			AddedAt:     time.Now().UTC().Add(-time.Minute),
		},
	}
	if err := store.addMistakes(-42, "en", entries); err != nil {
		t.Fatalf("add mistakes: %v", err)
	}

	response := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/delete", map[string]any{"index": 0})
	if ok, _ := response["ok"].(bool); !ok {
		t.Fatalf("delete response = %#v, want ok", response)
	}
	if total, _ := response["total"].(float64); total != 1 {
		t.Fatalf("delete total = %#v, want 1; response %#v", response["total"], response)
	}
	items, _ := response["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("delete items = %#v, want one remaining item", response)
	}
	remaining, _ := items[0].(map[string]any)
	if remaining["word"] != "second wrong" {
		t.Fatalf("wrong item was removed or ordering changed: %#v", response)
	}

	after := requestJSON(t, api, cookie, http.MethodGet, "/api/mistakes", nil)
	if total, _ := after["total"].(float64); total != 1 {
		t.Fatalf("persisted total = %#v, want 1; response %#v", after["total"], after)
	}
}

func TestWebLessonAndPracticeFallbackPersistMistakes(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	responses := []string{
		"1) She goes home.\n\n2) Use -s with she.\n\n3) Repeat: she goes.",
		`[{"word":"She go home","correction":"She goes home","explanation":"Use -s with she in Present Simple."}]`,
		"Model phrase: I would like coffee.\n\nQuestion audio: Would you like coffee?",
		`[{"word":"I want coffee please","correction":"I would like coffee, please.","explanation":"Use a more natural polite request."}]`,
	}
	call := 0
	api.bot.openrouter = newOpenRouterClient("test-key", "fallback-model", "http://localhost", "test", &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://openrouter.ai/api/v1/chat/completions" {
				t.Fatalf("unexpected OpenRouter URL: %s", req.URL.String())
			}
			if call >= len(responses) {
				t.Fatalf("unexpected extra OpenRouter call %d", call+1)
			}
			body, err := json.Marshal(map[string]any{
				"choices": []map[string]any{{"message": map[string]any{"content": responses[call]}}},
			})
			if err != nil {
				t.Fatal(err)
			}
			call++
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(body)),
			}, nil
		}),
	})

	if err := store.saveLesson(-42, "Say that she goes home."); err != nil {
		t.Fatalf("save lesson: %v", err)
	}
	lesson := requestJSON(t, api, cookie, http.MethodPost, "/api/lesson/answer", map[string]any{"text": "She go home"})
	if got, _ := lesson["mistakes"].([]any); len(got) != 1 {
		t.Fatalf("expected lesson fallback mistake, got %#v", lesson)
	}
	practice := requestJSON(t, api, cookie, http.MethodPost, "/api/practice", map[string]any{"text": "I want coffee please"})
	if got, _ := practice["mistakes"].([]any); len(got) != 1 {
		t.Fatalf("expected practice fallback mistake, got %#v", practice)
	}
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	mistakes := mistakesForLanguage(user)
	if len(mistakes) != 2 {
		t.Fatalf("expected two persisted mistakes, got %#v", mistakes)
	}
	if call != len(responses) {
		t.Fatalf("expected %d OpenRouter calls, got %d", len(responses), call)
	}
}

func TestWebPhrasebookAutoTranslatesEmptyNote(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	api.cfg.OpenRouterTranslatorModel = "translator-test"
	api.bot.cfg.OpenRouterTranslatorModel = "translator-test"
	api.bot.openrouter = newOpenRouterClient("test-key", "fallback-model", "http://localhost", "test", &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://openrouter.ai/api/v1/chat/completions" {
				t.Fatalf("unexpected OpenRouter URL %s", req.URL.String())
			}
			var payload map[string]any
			if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
				t.Fatalf("decode OpenRouter payload: %v", err)
			}
			if payload["model"] != "translator-test" {
				t.Fatalf("unexpected translator model %#v", payload["model"])
			}
			body := `{"choices":[{"message":{"content":"Не могли бы вы повторить, пожалуйста?"}}]}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}),
	})
	saved := requestJSON(t, api, cookie, http.MethodPost, "/api/phrasebook", map[string]any{
		"id":       "phrase-auto-translate",
		"phrase":   "Could you say that again, please?",
		"source":   "lesson",
		"language": "en",
	})
	items, _ := saved["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("expected saved phrasebook item, got %#v", saved)
	}
	item, _ := items[0].(map[string]any)
	if item["translation"] != "Не могли бы вы повторить, пожалуйста?" {
		t.Fatalf("expected auto translation in phrasebook item, got %#v", item)
	}
	user, err := store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if len(user.Phrasebook) != 1 || user.Phrasebook[0].Translation != "Не могли бы вы повторить, пожалуйста?" {
		t.Fatalf("phrasebook auto translation was not persisted: %#v", user.Phrasebook)
	}
}

func TestWebBugReportNotifiesTelegramRecipient(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	var telegramPayloads []map[string]any
	telegramServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sendMessage" {
			t.Fatalf("unexpected telegram method %s", r.URL.Path)
		}
		var telegramPayload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&telegramPayload); err != nil {
			t.Fatalf("decode telegram payload: %v", err)
		}
		telegramPayloads = append(telegramPayloads, telegramPayload)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ok":     true,
			"result": map[string]any{"message_id": 1},
		})
	}))
	defer telegramServer.Close()
	api.bot.telegram = &telegramClient{baseURL: telegramServer.URL, http: telegramServer.Client()}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("message", "Меню лагает при переносе кнопок"); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("view", "home"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/bug-report", &body)
	request.AddCookie(cookie)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("bug report returned %d: %s", recorder.Code, recorder.Body.String())
	}
	if len(telegramPayloads) != 2 {
		t.Fatalf("sent bug report to %d chats, want 2 default operators: %#v", len(telegramPayloads), telegramPayloads)
	}
	for index, want := range []any{float64(185156683), float64(297284024)} {
		if got := telegramPayloads[index]["chat_id"]; got != want {
			t.Fatalf("telegram chat %d = %#v, want %#v; payloads=%#v", index, got, want, telegramPayloads)
		}
	}
	text, _ := telegramPayloads[0]["text"].(string)
	for _, want := range []string{"Bug report NERIVA", "From: tester (-42)", "View: home", "Меню лагает"} {
		if !strings.Contains(text, want) {
			t.Fatalf("telegram text %q does not contain %q", text, want)
		}
	}
}

func TestWebBugReportUsesConfiguredTelegramOpsRecipients(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	api.cfg.TelegramOpsRecipients = []telegramOpsRecipient{{ChatID: "12345"}, {ChatID: "@poliglot_owner"}}
	api.bot.cfg = api.cfg
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	var chatIDs []any
	telegramServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sendMessage" {
			t.Fatalf("unexpected telegram method %s", r.URL.Path)
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode telegram payload: %v", err)
		}
		chatIDs = append(chatIDs, payload["chat_id"])
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": map[string]any{"message_id": len(chatIDs)}})
	}))
	defer telegramServer.Close()
	api.bot.telegram = &telegramClient{baseURL: telegramServer.URL, http: telegramServer.Client()}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("message", "Кнопка отчета не доходит владельцу"); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("view", "tutor"); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/bug-report", &body)
	request.AddCookie(cookie)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("bug report returned %d: %s", recorder.Code, recorder.Body.String())
	}
	if len(chatIDs) != 4 {
		t.Fatalf("sent bug report to %d chats, want 4: %#v", len(chatIDs), chatIDs)
	}
	if chatIDs[0] != float64(12345) || chatIDs[1] != "@poliglot_owner" || chatIDs[2] != float64(185156683) || chatIDs[3] != float64(297284024) {
		t.Fatalf("chat IDs = %#v, want configured recipients plus default operators", chatIDs)
	}
}

func TestConfigTelegramOpsRecipientsAlwaysIncludesDefaultOperators(t *testing.T) {
	recipients := (config{TelegramOpsRecipients: []telegramOpsRecipient{
		{ChatID: "12345"},
		{ChatID: "185156683"},
		{ChatID: "@poliglot_owner"},
	}}).telegramOpsRecipients()

	if len(recipients) != 4 {
		t.Fatalf("telegramOpsRecipients() returned %d recipients, want 4: %#v", len(recipients), recipients)
	}
	for index, want := range []string{"12345", "185156683", "@poliglot_owner", "297284024"} {
		if recipients[index].ChatID != want {
			t.Fatalf("recipient %d = %q, want %q; all recipients: %#v", index, recipients[index].ChatID, want, recipients)
		}
	}
}

func TestConfigFromEnvReadsTelegramOpsRecipients(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "bot-token")
	t.Setenv("OPENROUTER_API_KEY", "openrouter-key")
	t.Setenv("TELEGRAM_OPS_RECIPIENTS", "12345, @poliglot_owner")

	cfg, err := configFromEnv()
	if err != nil {
		t.Fatalf("configFromEnv() error = %v", err)
	}
	if len(cfg.TelegramOpsRecipients) != 2 {
		t.Fatalf("TelegramOpsRecipients = %#v, want 2 recipients", cfg.TelegramOpsRecipients)
	}
	if cfg.TelegramOpsRecipients[0].ChatID != "12345" || cfg.TelegramOpsRecipients[1].ChatID != "@poliglot_owner" {
		t.Fatalf("TelegramOpsRecipients = %#v, want configured ID and username", cfg.TelegramOpsRecipients)
	}
}

func TestWebPhrasebookRoutePersistsNoteInSession(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	body := strings.NewReader(`{"id":"note-1","phrase":"Could you repeat?","note":"Polite fallback","source":"manual","language":"en"}`)
	request := httptest.NewRequest(http.MethodPost, "/api/phrasebook", body)
	request.AddCookie(cookie)
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("phrasebook returned %d: %s", recorder.Code, recorder.Body.String())
	}

	user, err := api.bot.store.getOrCreateUser(-42, "tester")
	if err != nil {
		t.Fatal(err)
	}
	if len(user.Phrasebook) != 1 || user.Phrasebook[0].Note != "Polite fallback" {
		t.Fatalf("phrasebook note was not persisted: %#v", user.Phrasebook)
	}
	session := api.sessionPayload(user)
	sessionUser, ok := session["user"].(map[string]any)
	if !ok {
		t.Fatalf("session user has unexpected shape: %#v", session["user"])
	}
	if _, ok := sessionUser["phrasebook"]; !ok {
		t.Fatalf("session user does not include phrasebook: %#v", sessionUser)
	}
}

func TestWebSessionIncludesReferralInviteesFromSQLite(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	invitee, err := store.getOrCreateUser(-77, "invited learner")
	if err != nil {
		t.Fatal(err)
	}
	invitee.InvitedBy = -42
	invitee.XP = 520
	invitee.Level = "A2"
	invitee.ReferralLevelRewarded = true
	if err := store.saveUser(invitee); err != nil {
		t.Fatal(err)
	}

	session := requestJSON(t, api, cookie, http.MethodGet, "/api/session", nil)
	sessionUser, _ := session["user"].(map[string]any)
	invitees, _ := sessionUser["referral_invitees"].([]any)
	if len(invitees) != 1 {
		t.Fatalf("expected one referral invitee from sqlite, got %#v", sessionUser)
	}
	first, _ := invitees[0].(map[string]any)
	if first["name"] != "invited learner" || first["reached_level_3"] != true {
		t.Fatalf("unexpected invitee payload: %#v", first)
	}
}

func TestWebBugReportAcceptsImageOnlyAttachment(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("view", "practice"); err != nil {
		t.Fatal(err)
	}
	part, err := writer.CreateFormFile("screenshot", "clipboard.heic")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/bug-report", &body)
	request.AddCookie(cookie)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("image-only bug report returned %d: %s", recorder.Code, recorder.Body.String())
	}
	if _, err := os.Stat(filepath.Join("tmp", "bug-reports", "bug-reports.jsonl")); err != nil {
		t.Fatalf("bug report file was not saved: %v", err)
	}
}

func TestWebBugReportSendsScreenshotAsTelegramPhoto(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })

	var methods []string
	var multipartBodies []string
	telegramServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.URL.Path)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read telegram body: %v", err)
		}
		multipartBodies = append(multipartBodies, string(body))
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true, "result": map[string]any{"message_id": 1}})
	}))
	defer telegramServer.Close()
	api.bot.telegram = &telegramClient{baseURL: telegramServer.URL, http: telegramServer.Client()}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("message", "Скрин показывает баг в ошибках"); err != nil {
		t.Fatal(err)
	}
	if err := writer.WriteField("view", "mistakes"); err != nil {
		t.Fatal(err)
	}
	part, err := writer.CreateFormFile("screenshot", "bug.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	part, err = writer.CreateFormFile("screenshot", "bug-2.png")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n', 4, 5, 6}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/bug-report", &body)
	request.AddCookie(cookie)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("bug report returned %d: %s", recorder.Code, recorder.Body.String())
	}
	if len(methods) != 4 {
		t.Fatalf("expected four Telegram sendPhoto calls for two screenshots and two default operators, got %#v", methods)
	}
	for _, method := range methods {
		if method != "/sendPhoto" {
			t.Fatalf("expected only Telegram sendPhoto calls, got %#v", methods)
		}
	}
	multipartBody := strings.Join(multipartBodies, "\n---NEXT---\n")
	for _, want := range []string{`name="photo"; filename="`, `name="caption"`, "Bug report NERIVA", "Скрин показывает баг"} {
		if !strings.Contains(multipartBody, want) {
			t.Fatalf("sendPhoto body does not contain %q: %s", want, multipartBody)
		}
	}
	if !strings.Contains(multipartBody, "Screenshot: 2/2") {
		t.Fatalf("second screenshot caption was not sent: %s", multipartBody)
	}
}

func TestWebAppRoutesServeWebAndRedirectVersionedPaths(t *testing.T) {
	api := newWebAPI(config{WebAPISessionSecret: "test-session-secret"}, nil)
	mux := http.NewServeMux()
	api.register(mux)

	request := httptest.NewRequest(http.MethodGet, "/app", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /app returned %d", recorder.Code)
	}
	if body := recorder.Body.String(); !strings.Contains(body, `<div id="root">`) || !strings.Contains(body, `poliglot-boot`) || !strings.Contains(body, `/app/assets/`) {
		sample := body
		if len(sample) > 220 {
			sample = sample[:220]
		}
		t.Fatalf("/app should serve the React web shell, got: %s", sample)
	}

	for _, legacyPath := range []string{"/app/v1", "/app/v2"} {
		request = httptest.NewRequest(http.MethodGet, legacyPath+"?view=home", nil)
		recorder = httptest.NewRecorder()
		mux.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusMovedPermanently {
			t.Fatalf("GET %s returned %d", legacyPath, recorder.Code)
		}
		if location := recorder.Header().Get("Location"); location != "/app?view=home" {
			t.Fatalf("%s should redirect to /app preserving query, got %q", legacyPath, location)
		}
	}
}

func TestPremiumPlansDTOIncludesLandingPricingCopy(t *testing.T) {
	cfg := config{
		PremiumRubPrice:        300,
		PremiumStarsPrice:      150,
		PremiumYearRubPrice:    3000,
		PremiumYearStarsPrice:  1500,
		PlatinumRubPrice:       590,
		PlatinumStarsPrice:     300,
		PlatinumYearRubPrice:   5900,
		PlatinumYearStarsPrice: 3000,
	}
	api := newWebAPI(cfg, &bot{cfg: cfg})

	plans := api.premiumPlansDTO(userState{InterfaceLanguage: "ru"})
	if len(plans) != 5 {
		t.Fatalf("plans len = %d, want free + 4 paid plans: %#v", len(plans), plans)
	}

	free := plans[0]
	if free["product"] != "free" || free["label"] != "Начальный" || free["body"] != "Базовое текстовое обучение, тренировка слов, Phrasebook и обзор прогресса. AI Tutor, аудирование, произношение и голосовые проверки открываются в Premium." {
		t.Fatalf("free plan copy mismatch: %#v", free)
	}
	assertStringSlicesEqual(t, free["features"].([]string), []string{"ежедневная привычка и стартовые уроки", "базовая тренировка слов", "заметки, phrasebook и обзор прогресса"})
	assertStringSlicesEqual(t, free["locked_features"].([]string), []string{"Уроки с AI Tutor", "Аудирование и произношение", "Проверка голоса и фото-инструменты"})
	if free["note"] != "Подходит для знакомства с продуктом без оплаты." {
		t.Fatalf("free plan note = %q", free["note"])
	}

	premium := plans[1]
	assertStringSlicesEqual(t, premium["features"].([]string), []string{"голос в текст и перевод услышанного", "перевод текста с картинки", "практика по контексту голоса или фото", "словарь ошибок, notes, XP и streak"})
	if premium["label"] != "Регулярная учеба" || premium["note"] != "Лучший выбор для стабильного ежедневного обучения." {
		t.Fatalf("premium plan copy mismatch: %#v", premium)
	}

	platinumYear := plans[4]
	if platinumYear["title"] != "Platinum на год" || platinumYear["label"] != "Интенсив" {
		t.Fatalf("platinum yearly title/label mismatch: %#v", platinumYear)
	}

	englishPlans := api.premiumPlansDTO(userState{InterfaceLanguage: "en"})
	if englishPlans[0]["label"] != "Starter" || englishPlans[0]["body"] == free["body"] {
		t.Fatalf("English free plan was not localized: %#v", englishPlans[0])
	}
	if englishPlans[1]["label"] != "Regular study" || englishPlans[1]["body"] == premium["body"] {
		t.Fatalf("English premium plan was not localized: %#v", englishPlans[1])
	}
}

func TestPremiumPlansDTOUsesLiveTonAPIRateForUSDTPrices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rates" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"rates":{"USDT_MASTER":{"prices":{"RUB":75}}}}`))
	}))
	defer server.Close()

	cfg := config{
		PremiumRubPrice:           300,
		PremiumStarsPrice:         150,
		CryptoUSDTRubRate:         72,
		CryptoTONAPIBaseURL:       server.URL,
		CryptoUSDTTONJettonMaster: "USDT_MASTER",
		CryptoUSDTTONWallet:       "EQ_USDT",
	}
	bot := &bot{cfg: cfg, cryptoRates: newCryptoRateProvider(cfg, server.Client())}
	api := newWebAPI(cfg, bot)

	plans := api.premiumPlansDTO(userState{InterfaceLanguage: "en"})
	if len(plans) < 2 {
		t.Fatalf("plans len = %d, want paid plans: %#v", len(plans), plans)
	}
	premium := plans[1]
	if premium["usdt_price"] != "4 USDT" {
		t.Fatalf("premium usdt_price = %#v, want live 75 RUB/USDT conversion", premium["usdt_price"])
	}
	methods, ok := premium["crypto_methods"].([]map[string]any)
	if !ok || len(methods) == 0 {
		t.Fatalf("premium crypto_methods missing: %#v", premium["crypto_methods"])
	}
	if methods[0]["currency"] != cryptoCurrencyUSDT || methods[0]["network"] != cryptoNetworkTON {
		t.Fatalf("unexpected crypto method: %#v", methods[0])
	}
}

func TestWebUserDTOUsesLiveTonAPIRateForReferralUSDT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rates" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"rates":{"USDT_MASTER":{"prices":{"RUB":75}}}}`))
	}))
	defer server.Close()

	cfg := config{
		CryptoUSDTRubRate:         72,
		CryptoTONAPIBaseURL:       server.URL,
		CryptoUSDTTONJettonMaster: "USDT_MASTER",
	}
	bot := &bot{cfg: cfg, cryptoRates: newCryptoRateProvider(cfg, server.Client())}
	api := newWebAPI(cfg, bot)

	dto := api.userDTO(userState{ReferralBalanceKopecks: 7500 * 100})
	if dto["usdt_rub_rate"] != "75" {
		t.Fatalf("usdt_rub_rate = %#v, want live rate", dto["usdt_rub_rate"])
	}
	if dto["referral_balance_usdt"] != "100.00 USDT" {
		t.Fatalf("referral_balance_usdt = %#v, want live 75 RUB/USDT conversion", dto["referral_balance_usdt"])
	}
	if dto["referral_withdraw_min_usdt"] != "13.33 USDT" {
		t.Fatalf("referral_withdraw_min_usdt = %#v, want live 75 RUB/USDT conversion", dto["referral_withdraw_min_usdt"])
	}
}

func TestWebDirectCryptoPaymentCreatesTONInvoice(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "tester", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-42, "tester"); err != nil {
		t.Fatal(err)
	}
	cfg := config{
		WebAPISessionSecret:     "test-session-secret",
		CryptoTONWallet:         "EQ_TEST",
		CryptoTONMonthAmount:    "1.25",
		CryptoTONYearAmount:     "10",
		CryptoTONCenterBaseURL:  "https://toncenter.com/api/v2",
		CryptoPaymentTTLMinutes: 60,
	}
	bot := &bot{cfg: cfg, store: store}
	api := newWebAPI(cfg, bot)
	recorder := httptest.NewRecorder()
	api.setSessionCookie(recorder, -42)
	cookies := recorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}
	payment := requestJSON(t, api, cookies[0], http.MethodPost, "/api/premium/crypto/payment", map[string]any{"product": premiumMonthlyProduct})
	if payment["network"] != cryptoNetworkTON || payment["currency"] != cryptoCurrencyTON {
		t.Fatalf("unexpected crypto payment response: %#v", payment)
	}
	if memo, _ := payment["memo"].(string); !strings.HasPrefix(memo, "POLIGLOT:ton_") {
		t.Fatalf("expected unique memo, got %#v", payment)
	}
	if invoiceURL, _ := payment["invoice_url"].(string); !strings.HasPrefix(invoiceURL, "https://app.tonkeeper.com/transfer/") {
		t.Fatalf("expected Tonkeeper invoice url, got %#v", payment)
	}
	var count int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM crypto_payments`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("crypto_payments count = %d, want 1", count)
	}
}

func TestWebStarsPaymentSendsTelegramInvoice(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(123, "telegramuser", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(123, "telegramuser"); err != nil {
		t.Fatal(err)
	}

	var method string
	var payload map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cfg := config{
		WebAPISessionSecret: "test-session-secret",
		WebTelegramLoginBot: "NERIVAapp_bot",
		MaxVoiceSeconds:     30,
		PremiumRubPrice:     300,
		PremiumStarsPrice:   150,
	}
	bot := &bot{cfg: cfg, store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	api := newWebAPI(cfg, bot)
	recorder := httptest.NewRecorder()
	api.setSessionCookie(recorder, 123)
	cookies := recorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}

	response := requestJSON(t, api, cookies[0], http.MethodPost, "/api/premium/stars", map[string]any{"product": premiumMonthlyProduct})
	if sent, _ := response["sent"].(bool); !sent {
		t.Fatalf("expected sent response, got %#v", response)
	}
	if botURL, _ := response["bot_url"].(string); botURL != "https://t.me/NERIVAapp_bot?start=buy_premium_30d" {
		t.Fatalf("unexpected bot url: %#v", response)
	}
	if method != "/sendInvoice" {
		t.Fatalf("expected sendInvoice call, got %s", method)
	}
	if payload["chat_id"].(float64) != 123 || payload["currency"].(string) != "XTR" {
		t.Fatalf("unexpected invoice payload: %#v", payload)
	}
	prices := payload["prices"].([]any)
	if len(prices) != 1 || int(prices[0].(map[string]any)["amount"].(float64)) != 150 {
		t.Fatalf("expected 150 Stars invoice, got %#v", payload)
	}
	if invoicePayload, _ := payload["payload"].(string); !strings.Contains(invoicePayload, premiumMonthlyProduct+"_123_") {
		t.Fatalf("unexpected invoice payload id: %#v", payload)
	}
}

func TestTelegramStarsSuccessfulPaymentNotifiesOpsRecipient(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()
	user, err := store.getOrCreateUser(123, "Paid User")
	if err != nil {
		t.Fatal(err)
	}

	var payloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		payloads = append(payloads, payload)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cfg := config{PremiumRubPrice: 300, PremiumStarsPrice: 150}
	bot := &bot{cfg: cfg, store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	plan, ok := bot.premiumPlan(premiumMonthlyProduct)
	if !ok {
		t.Fatal("expected monthly premium plan")
	}
	err = bot.handleSuccessfulPayment(context.Background(), user.TelegramID, user, successfulPayment{
		Currency:                "XTR",
		TotalAmount:             plan.StarsPrice,
		InvoicePayload:          premiumPaymentPayload(plan.Product, user.TelegramID, time.Now()),
		TelegramPaymentChargeID: "stars-test-payment",
	})
	if err != nil {
		t.Fatal(err)
	}

	var opsPayload map[string]any
	for _, payload := range payloads {
		if payload["chat_id"] == float64(telegramOpsRecipientID) {
			opsPayload = payload
			break
		}
	}
	if opsPayload == nil {
		t.Fatalf("expected ops notification payload, got %#v", payloads)
	}
	text, _ := opsPayload["text"].(string)
	for _, want := range []string{"Payment NERIVA", "Source: Telegram Stars", "User: Paid User (123)", "stars-test-payment"} {
		if !strings.Contains(text, want) {
			t.Fatalf("ops payment text %q does not contain %q", text, want)
		}
	}
}

func TestWebAITutorStartReturnsSessionStep(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	body := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/start", map[string]any{})
	if _, ok := body["session"].(map[string]any); !ok {
		t.Fatalf("missing session in response: %#v", body)
	}
	next, ok := body["next_step"].(map[string]any)
	if !ok || next["stage"] != aiTutorStageStoryIntro {
		t.Fatalf("next_step = %#v", body["next_step"])
	}
}

func TestWebAITutorStartRequiresPremium(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-free-blocked", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}

	status, body := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/ai-tutor/start", map[string]any{})
	if status != http.StatusPaymentRequired {
		t.Fatalf("start status = %d body=%#v, want 402", status, body)
	}
	apiError, _ := body["error"].(map[string]any)
	if apiError["code"] != "premium_required" {
		t.Fatalf("error = %#v, want premium_required", body["error"])
	}
}

func TestWebTutorStartUsesAITutorEngine(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-legacy-route", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}

	body := requestJSON(t, api, cookie, http.MethodPost, "/api/tutor/start", map[string]any{})
	if _, ok := body["next_step"].(map[string]any); !ok {
		t.Fatalf("legacy tutor route did not return AI tutor step: %#v", body)
	}
	if _, ok := body["tutor_lesson"]; ok {
		t.Fatalf("legacy tutor route returned local tutor lesson payload: %#v", body)
	}
}

func TestWebAITutorAnswerAdvancesSession(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-2", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	start := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/start", map[string]any{})
	session := start["session"].(map[string]any)
	sessionID, _ := session["ID"].(string)
	if sessionID == "" {
		sessionID, _ = session["id"].(string)
	}
	answer := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/answer", map[string]any{"session_id": sessionID, "text": "continue"})
	if answer["current_stage"] != aiTutorStageRetell {
		t.Fatalf("current_stage = %#v response=%#v", answer["current_stage"], answer)
	}
}

func TestWebAITutorReviewCompletesSession(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-3", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	start := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/start", map[string]any{})
	session := start["session"].(map[string]any)
	sessionID, _ := session["ID"].(string)
	if sessionID == "" {
		sessionID, _ = session["id"].(string)
	}
	if err := store.updateAITutorSessionStage(sessionID, aiTutorStageReviewSchedule, aiTutorSessionActive, ""); err != nil {
		t.Fatalf("updateAITutorSessionStage() error = %v", err)
	}
	done := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/review", map[string]any{"session_id": sessionID, "choice": "no_review"})
	if done["current_stage"] != aiTutorStageComplete {
		t.Fatalf("current_stage = %#v response=%#v", done["current_stage"], done)
	}
	user, ok := done["user"].(map[string]any)
	if !ok {
		t.Fatalf("missing refreshed user in response: %#v", done)
	}
	if xp, _ := user["xp"].(float64); int(xp) < aiTutorCompletionXP {
		t.Fatalf("AI Tutor completion XP missing from response user: %#v", user)
	}
}

func TestWebAITutorCompletedLessonsReturnsFinishedSessions(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	payload := validAITutorLessonPayloadForTest()
	payload.Title = "Sports Sunday"
	payload.Theme = "sports"
	lesson := aiTutorLessonRecord{ID: "lesson-completed-history", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	session := aiTutorSessionRecord{ID: "session-completed-history", TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageComplete, Status: aiTutorSessionComplete, CompletedAt: "2026-06-12T12:00:00Z"}
	if err := store.createAITutorSession(session); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}
	otherLesson := aiTutorLessonRecord{ID: "lesson-completed-other-user", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}
	if err := store.saveAITutorLesson(otherLesson); err != nil {
		t.Fatalf("saveAITutorLesson(other) error = %v", err)
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{ID: "session-completed-other-user", TelegramID: -7, LessonID: otherLesson.ID, Surface: "web", CurrentStage: aiTutorStageComplete, Status: aiTutorSessionComplete, CompletedAt: "2026-06-12T13:00:00Z"}); err != nil {
		t.Fatalf("createAITutorSession(other) error = %v", err)
	}

	body := requestJSON(t, api, cookie, http.MethodGet, "/api/ai-tutor/completed", nil)
	items, _ := body["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("items = %#v, want one completed lesson for current user", body["items"])
	}
	item, _ := items[0].(map[string]any)
	if item["lesson_id"] != lesson.ID || item["session_id"] != session.ID || item["title"] != "Sports Sunday" {
		t.Fatalf("completed lesson item mismatch: %#v", item)
	}
	if _, ok := item["lesson"].(map[string]any); !ok {
		t.Fatalf("completed lesson item missing lesson payload: %#v", item)
	}
}

func TestWebAITutorRestartCreatesNewSessionForCompletedLesson(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-restart-history", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{ID: "session-old-restart", TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageComplete, Status: aiTutorSessionComplete, CompletedAt: "2026-06-12T12:00:00Z"}); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}

	body := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/restart", map[string]any{"lesson_id": lesson.ID})
	if body["current_stage"] != aiTutorStageStoryIntro {
		t.Fatalf("current_stage = %#v response=%#v", body["current_stage"], body)
	}
	session, _ := body["session"].(map[string]any)
	if session["LessonID"] != lesson.ID && session["lesson_id"] != lesson.ID {
		t.Fatalf("restart returned wrong lesson session: %#v", session)
	}
	if session["ID"] == "session-old-restart" || session["id"] == "session-old-restart" {
		t.Fatalf("restart reused the old completed session: %#v", session)
	}
}

func TestWebAITutorRestartRejectsUnrelatedLesson(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	lesson := aiTutorLessonRecord{ID: "lesson-restart-forbidden", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}

	status, body := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/ai-tutor/restart", map[string]any{"lesson_id": lesson.ID})
	if status != http.StatusForbidden {
		t.Fatalf("restart status = %d body=%#v, want 403", status, body)
	}
}

func TestAITutorWordReportCreatesTelegramNotificationAndSkipsWord(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	var telegramPayloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode telegram payload: %v", err)
		}
		telegramPayloads = append(telegramPayloads, payload)
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":77}}`))
	}))
	defer server.Close()
	api.cfg.TelegramOpsRecipients = []telegramOpsRecipient{{ChatID: "12345"}}
	api.bot.cfg = api.cfg
	api.bot.telegram = &telegramClient{baseURL: server.URL, http: server.Client()}

	lesson := aiTutorLessonRecord{ID: "lesson-word-report-web", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{ID: "session-word-report-web", TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorWordLearnStage(1), Status: aiTutorSessionActive}); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}

	body := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/word-report", map[string]any{
		"session_id":           "session-word-report-web",
		"stage":                aiTutorWordLearnStage(1),
		"proposed_word":        "get up",
		"proposed_translation": "vstavat",
		"comment":              "more natural phrase",
	})
	if body["current_stage"] != aiTutorWordLearnStage(2) {
		t.Fatalf("current_stage = %#v response=%#v, want %s", body["current_stage"], body, aiTutorWordLearnStage(2))
	}
	next, _ := body["next_step"].(map[string]any)
	if next["stage"] != aiTutorWordLearnStage(2) || next["kind"] != "word_learn" {
		t.Fatalf("next_step = %#v, want next word learn step", next)
	}
	report, ok, err := store.findPendingAITutorWordReportByFixPrompt("", 0)
	if err != nil {
		t.Fatalf("findPendingAITutorWordReportByFixPrompt(empty) error = %v", err)
	}
	if ok {
		t.Fatalf("empty fix prompt should not match report: %#v", report)
	}
	userCount, sessionCount, err := store.countAITutorWordReports(-42, "session-word-report-web", time.Now().UTC().Add(-10*time.Minute))
	if err != nil {
		t.Fatalf("countAITutorWordReports() error = %v", err)
	}
	if userCount != 1 || sessionCount != 1 {
		t.Fatalf("report counts user=%d session=%d, want 1/1", userCount, sessionCount)
	}
	var stored aiTutorWordReportRecord
	rows, err := store.db.Query(`SELECT ` + aiTutorWordReportSelectColumns + ` FROM ai_tutor_word_reports`)
	if err != nil {
		t.Fatalf("query reports: %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Fatal("expected stored report row")
	}
	stored, err = scanAITutorWordReportRecord(rows)
	if err != nil {
		t.Fatalf("scan stored report: %v", err)
	}
	if stored.OriginalWord != "wake up" || stored.ProposedWord != "get up" || stored.ProposedTranslation != "vstavat" || stored.WordIndex != 0 {
		t.Fatalf("stored report mismatch: %#v", stored)
	}
	if len(telegramPayloads) == 0 {
		t.Fatal("expected telegram ops notification")
	}
	if len(telegramPayloads) != 3 {
		t.Fatalf("sent word report to %d chats, want configured recipient plus two default operators: %#v", len(telegramPayloads), telegramPayloads)
	}
	for index, want := range []any{float64(12345), float64(185156683), float64(297284024)} {
		if got := telegramPayloads[index]["chat_id"]; got != want {
			t.Fatalf("word report chat %d = %#v, want %#v; payloads=%#v", index, got, want, telegramPayloads)
		}
	}
	text, _ := telegramPayloads[0]["text"].(string)
	if !strings.Contains(text, stored.ID) ||
		!strings.Contains(text, "wake up - prosypatsya") ||
		!strings.Contains(text, "get up - vstavat") ||
		!strings.Contains(text, "vocabulary_ai_translations") {
		t.Fatalf("telegram text missing report details: %q", text)
	}
	keyboardJSON, _ := json.Marshal(telegramPayloads[0]["reply_markup"])
	for _, want := range []string{"ait_word_report|" + stored.ID + "|accept", "ait_word_report|" + stored.ID + "|reject", "ait_word_report|" + stored.ID + "|fix"} {
		if !strings.Contains(string(keyboardJSON), want) {
			t.Fatalf("telegram keyboard missing %q: %s", want, string(keyboardJSON))
		}
	}
}

func TestWebWordsNextIncludesReportMetadata(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:           "en:western",
			Language:     "en",
			English:      "western",
			Russian:      "западный",
			Translations: map[string]string{"ru": "западный"},
			Context:      "Side of the world where the sun sets.",
			Level:        "A1",
		},
	})
	api, _, cookie := newTestWebAPI(t)

	body := requestJSON(t, api, cookie, http.MethodPost, "/api/words/next", map[string]any{})
	if body["word_id"] != "en:western" {
		t.Fatalf("word_id = %#v body=%#v, want en:western", body["word_id"], body)
	}
	if body["word"] != "western" || body["translation"] != "западный" {
		t.Fatalf("word metadata mismatch: %#v", body)
	}
	if body["reportable"] != true {
		t.Fatalf("reportable = %#v body=%#v, want true", body["reportable"], body)
	}
}

func TestWebWordsReportCreatesTelegramNotification(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:           "en:western",
			Language:     "en",
			English:      "western",
			Russian:      "западный",
			Translations: map[string]string{"ru": "западный"},
			Context:      "Side of the world where the sun sets.",
			Level:        "A1",
		},
	})
	api, store, cookie := newTestWebAPI(t)
	var telegramPayloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode telegram payload: %v", err)
		}
		telegramPayloads = append(telegramPayloads, payload)
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":88}}`))
	}))
	defer server.Close()
	api.cfg.TelegramOpsRecipients = []telegramOpsRecipient{{ChatID: "12345"}}
	api.bot.cfg = api.cfg
	api.bot.telegram = &telegramClient{baseURL: server.URL, http: server.Client()}

	next := requestJSON(t, api, cookie, http.MethodPost, "/api/words/next", map[string]any{})
	body := requestJSON(t, api, cookie, http.MethodPost, "/api/words/report", map[string]any{
		"word_id":              next["word_id"],
		"proposed_word":        "westward",
		"proposed_translation": "на запад",
		"comment":              "current card uses adjective instead of direction",
	})
	if body["ok"] != true {
		t.Fatalf("report response = %#v, want ok true", body)
	}

	var stored aiTutorWordReportRecord
	rows, err := store.db.Query(`SELECT ` + aiTutorWordReportSelectColumns + ` FROM ai_tutor_word_reports`)
	if err != nil {
		t.Fatalf("query reports: %v", err)
	}
	defer rows.Close()
	if !rows.Next() {
		t.Fatal("expected stored report row")
	}
	stored, err = scanAITutorWordReportRecord(rows)
	if err != nil {
		t.Fatalf("scan stored report: %v", err)
	}
	if stored.Stage != "vocabulary_word" || stored.LessonID != "en:western" || stored.SessionID != "vocabulary|ru|en:western" {
		t.Fatalf("stored vocabulary report routing mismatch: %#v", stored)
	}
	if stored.OriginalWord != "western" || stored.OriginalTranslation != "западный" || stored.ProposedWord != "westward" || stored.ProposedTranslation != "на запад" {
		t.Fatalf("stored report content mismatch: %#v", stored)
	}
	if len(telegramPayloads) == 0 {
		t.Fatal("expected telegram ops notification")
	}
	text, _ := telegramPayloads[0]["text"].(string)
	for _, want := range []string{
		"Vocabulary word report",
		stored.ID,
		"western - западный",
		"westward - на запад",
		"vocabulary_words",
		"vocabulary_ai_translations",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("telegram text missing %q: %q", want, text)
		}
	}
	keyboardJSON, _ := json.Marshal(telegramPayloads[0]["reply_markup"])
	for _, want := range []string{"ait_word_report|" + stored.ID + "|accept", "ait_word_report|" + stored.ID + "|reject", "ait_word_report|" + stored.ID + "|fix"} {
		if !strings.Contains(string(keyboardJSON), want) {
			t.Fatalf("telegram keyboard missing %q: %s", want, string(keyboardJSON))
		}
	}
}

func TestTelegramVocabularyWordReportAcceptAppliesSQLiteCorrection(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:           "en:western",
			Language:     "en",
			English:      "western",
			Russian:      "западный",
			Translations: map[string]string{"ru": "западный"},
			Context:      "Side of the world where the sun sets.",
			Level:        "A1",
		},
	})
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })
	report := aiTutorWordReportRecord{
		ID:                  "report-vocabulary-accept",
		TelegramID:          -42,
		SessionID:           "vocabulary|ru|en:western",
		LessonID:            "en:western",
		Stage:               "vocabulary_word",
		WordIndex:           -1,
		OriginalWord:        "western",
		OriginalTranslation: "западный",
		ProposedWord:        "westward",
		ProposedTranslation: "на запад",
		Status:              aiTutorWordReportPending,
	}
	if _, _, err := store.createAITutorWordReport(report); err != nil {
		t.Fatalf("createAITutorWordReport() error = %v", err)
	}
	var telegramPayloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode telegram payload: %v", err)
		}
		telegramPayloads = append(telegramPayloads, payload)
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":89}}`))
	}))
	defer server.Close()
	b := &bot{
		cfg:      config{TelegramOpsRecipients: []telegramOpsRecipient{{ChatID: "12345"}}},
		store:    store,
		telegram: &telegramClient{baseURL: server.URL, http: server.Client()},
	}

	if err := b.acceptAITutorWordReport(context.Background(), 12345, report); err != nil {
		t.Fatalf("acceptAITutorWordReport() error = %v", err)
	}

	db := currentSQLiteVocabularyDB()
	if db == nil {
		t.Fatal("expected configured SQLite vocabulary DB")
	}
	var word, russian, translated, aiSource, aiTranslation string
	if err := db.QueryRow(`SELECT word, russian, source FROM vocabulary_words WHERE id = ?`, "en:western").Scan(&word, &russian, &aiSource); err != nil {
		t.Fatalf("query vocabulary_words: %v", err)
	}
	if err := db.QueryRow(`SELECT text FROM vocabulary_translations WHERE word_id = ? AND language = ?`, "en:western", "ru").Scan(&translated); err != nil {
		t.Fatalf("query vocabulary_translations: %v", err)
	}
	if err := db.QueryRow(`SELECT translation FROM vocabulary_ai_translations WHERE word_id = ? AND target_language = ?`, "en:western", "ru").Scan(&aiTranslation); err != nil {
		t.Fatalf("query vocabulary_ai_translations: %v", err)
	}
	if word != "westward" || russian != "на запад" || translated != "на запад" || aiTranslation != "на запад" {
		t.Fatalf("SQLite correction mismatch: word=%q russian=%q translated=%q aiTranslation=%q", word, russian, translated, aiTranslation)
	}
	if !strings.Contains(aiSource, "report") {
		t.Fatalf("vocabulary source = %q, want report marker", aiSource)
	}
	if len(telegramPayloads) == 0 {
		t.Fatal("expected resolved notice")
	}
}

func TestAITutorWordReportRejectsStaleStage(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	calledTelegram := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledTelegram = true
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	api.bot.telegram = &telegramClient{baseURL: server.URL, http: server.Client()}
	lesson := aiTutorLessonRecord{ID: "lesson-word-report-stale", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{ID: "session-word-report-stale", TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorWordLearnStage(2), Status: aiTutorSessionActive}); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}

	status, response := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/ai-tutor/word-report", map[string]any{
		"session_id":           "session-word-report-stale",
		"stage":                aiTutorWordLearnStage(1),
		"proposed_word":        "get up",
		"proposed_translation": "vstavat",
	})
	if status != http.StatusConflict {
		t.Fatalf("stale report status=%d response=%#v, want 409", status, response)
	}
	if calledTelegram {
		t.Fatal("stale report should not send telegram notification")
	}
}

func TestAITutorWordReportRateLimitsUser(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	grantTestPremium(t, store, -42)
	lesson := aiTutorLessonRecord{ID: "lesson-word-report-rate", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	for index := 1; index <= 3; index++ {
		sessionID := "session-word-report-rate-" + strconv.Itoa(index)
		if err := store.createAITutorSession(aiTutorSessionRecord{ID: sessionID, TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorWordLearnStage(index), Status: aiTutorSessionActive}); err != nil {
			t.Fatalf("createAITutorSession(%s) error = %v", sessionID, err)
		}
		requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/word-report", map[string]any{
			"session_id":           sessionID,
			"stage":                aiTutorWordLearnStage(index),
			"proposed_word":        "word " + strconv.Itoa(index),
			"proposed_translation": "translation " + strconv.Itoa(index),
		})
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{ID: "session-word-report-rate-4", TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorWordLearnStage(4), Status: aiTutorSessionActive}); err != nil {
		t.Fatalf("createAITutorSession(rate-4) error = %v", err)
	}
	status, response := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/ai-tutor/word-report", map[string]any{
		"session_id":           "session-word-report-rate-4",
		"stage":                aiTutorWordLearnStage(4),
		"proposed_word":        "word 4",
		"proposed_translation": "translation 4",
	})
	if status != http.StatusTooManyRequests {
		t.Fatalf("rate-limited report status=%d response=%#v, want 429", status, response)
	}
}

func newTestWebAPI(t *testing.T) (*webAPI, *sqliteStore, *http.Cookie) {
	t.Helper()
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })

	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-42, "tester", hash); err != nil {
		t.Fatal(err)
	}
	if _, err := store.getOrCreateUser(-42, "tester"); err != nil {
		t.Fatal(err)
	}
	if err := store.setInterfaceLanguage(-42, "ru"); err != nil {
		t.Fatal(err)
	}
	if err := store.setLearningLanguage(-42, "en"); err != nil {
		t.Fatal(err)
	}

	cfg := config{WebAPISessionSecret: "test-session-secret"}
	bot := &bot{cfg: cfg, store: store}
	api := newWebAPI(cfg, bot)
	recorder := httptest.NewRecorder()
	api.setSessionCookie(recorder, -42)
	cookies := recorder.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected session cookie")
	}
	return api, store, cookies[0]
}

func grantTestPremium(t *testing.T, store *sqliteStore, telegramID int64) {
	t.Helper()
	user, ok, err := store.getUser(telegramID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("test user %d not found", telegramID)
	}
	user.Plan = "premium"
	user.PremiumUntil = time.Now().UTC().Add(30 * 24 * time.Hour)
	if err := store.saveUser(user); err != nil {
		t.Fatalf("grant premium: %v", err)
	}
}

func addLearnedTestWord(t *testing.T, store *sqliteStore, userID int64, wordID string) {
	t.Helper()
	word, ok := findVocabWord(wordID)
	if !ok {
		t.Fatalf("test vocabulary word not found: %s", wordID)
	}
	if _, _, err := store.addLearnedWord(userID, word); err != nil {
		t.Fatalf("addLearnedWord(%s): %v", wordID, err)
	}
}

func learnedWordSpellingCount(user userState, wordID string) int {
	for _, word := range user.LearnedWords {
		if learnedWordIDMatches(word, wordID) {
			return word.SpellingCorrectCount
		}
	}
	return 0
}

func anyStrings(values []any) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if text, ok := value.(string); ok {
			result = append(result, text)
		}
	}
	return result
}

func requestJSON(t *testing.T, api *webAPI, cookie *http.Cookie, method string, target string, body any) map[string]any {
	t.Helper()
	status, response := requestJSONRaw(t, api, cookie, method, target, body)
	if status < 200 || status >= 300 {
		t.Fatalf("%s %s returned %d: %#v", method, target, status, response)
	}
	return response
}

func requestJSONRaw(t *testing.T, api *webAPI, cookie *http.Cookie, method string, target string, body any) (int, map[string]any) {
	t.Helper()
	var reader *bytes.Reader
	if body == nil {
		reader = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(payload)
	}
	request := httptest.NewRequest(method, target, reader)
	if cookie != nil {
		request.AddCookie(cookie)
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)
	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v\nbody: %s", err, recorder.Body.String())
	}
	return recorder.Code, response
}

func requestRaw(t *testing.T, api *webAPI, cookie *http.Cookie, method string, target string, body any) ([]byte, int) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(data)
	}
	request := httptest.NewRequest(method, target, reader)
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if cookie != nil {
		request.AddCookie(cookie)
	}
	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	api.register(mux)
	mux.ServeHTTP(recorder, request)
	return recorder.Body.Bytes(), recorder.Code
}
