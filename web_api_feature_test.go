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

func TestWebAppV2ShellAndPWAEntrypointsAreNoCache(t *testing.T) {
	api := newWebAPI(config{WebAPISessionSecret: "test-session-secret"}, nil)
	mux := http.NewServeMux()
	api.register(mux)

	for _, target := range []string{"/app/v2/", "/app/v2/manifest.webmanifest", "/app/v2/offline-deck-sw.js"} {
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
	if total, _ := vocabulary["total"].(float64); total != 1 {
		t.Fatalf("expected one mastered vocabulary word, got %#v", vocabulary)
	}
	if items, _ := vocabulary["items"].([]any); len(items) != 1 {
		t.Fatalf("expected one vocabulary item, got %#v", vocabulary)
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

	var telegramPayload map[string]any
	telegramServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sendMessage" {
			t.Fatalf("unexpected telegram method %s", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&telegramPayload); err != nil {
			t.Fatalf("decode telegram payload: %v", err)
		}
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
	if telegramPayload["chat_id"] != float64(telegramOpsRecipientID) {
		t.Fatalf("telegram chat_id = %#v, want %d", telegramPayload["chat_id"], telegramOpsRecipientID)
	}
	text, _ := telegramPayload["text"].(string)
	for _, want := range []string{"Bug report Poliglot AI", "From: tester (-42)", "View: home", "Меню лагает"} {
		if !strings.Contains(text, want) {
			t.Fatalf("telegram text %q does not contain %q", text, want)
		}
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
	if len(methods) != 2 || methods[0] != "/sendPhoto" || methods[1] != "/sendPhoto" {
		t.Fatalf("expected two Telegram sendPhoto calls, got %#v", methods)
	}
	multipartBody := strings.Join(multipartBodies, "\n---NEXT---\n")
	for _, want := range []string{`name="photo"; filename="`, `name="caption"`, "Bug report Poliglot AI", "Скрин показывает баг"} {
		if !strings.Contains(multipartBody, want) {
			t.Fatalf("sendPhoto body does not contain %q: %s", want, multipartBody)
		}
	}
	if !strings.Contains(multipartBody, "Screenshot: 2/2") {
		t.Fatalf("second screenshot caption was not sent: %s", multipartBody)
	}
}

func TestWebAppRoutesServeV2AndRedirectLegacyV1(t *testing.T) {
	api := newWebAPI(config{WebAPISessionSecret: "test-session-secret"}, nil)
	mux := http.NewServeMux()
	api.register(mux)

	request := httptest.NewRequest(http.MethodGet, "/app", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("GET /app returned %d", recorder.Code)
	}
	if body := recorder.Body.String(); !strings.Contains(body, `<div id="root">`) || !strings.Contains(body, `poliglot-boot`) || !strings.Contains(body, `/app/v2/assets/`) {
		sample := body
		if len(sample) > 220 {
			sample = sample[:220]
		}
		t.Fatalf("/app should serve the React V2 shell, got: %s", sample)
	}

	request = httptest.NewRequest(http.MethodGet, "/app/v1", nil)
	recorder = httptest.NewRecorder()
	mux.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusMovedPermanently {
		t.Fatalf("GET /app/v1 returned %d", recorder.Code)
	}
	if location := recorder.Header().Get("Location"); location != "/app/v2/" {
		t.Fatalf("/app/v1 should redirect to /app/v2/, got %q", location)
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
		WebTelegramLoginBot: "Poliglot_AI_bot",
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
	if botURL, _ := response["bot_url"].(string); botURL != "https://t.me/Poliglot_AI_bot?start=buy_premium_30d" {
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
	for _, want := range []string{"Payment Poliglot AI", "Source: Telegram Stars", "User: Paid User (123)", "stars-test-payment"} {
		if !strings.Contains(text, want) {
			t.Fatalf("ops payment text %q does not contain %q", text, want)
		}
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
