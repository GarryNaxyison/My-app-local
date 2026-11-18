package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const sampleContextLesson = `{"questions":[
	{"context":"Нужен ___ для поездки.","options":["ticket","hotel","train"],"correct_index":0,"explanation":"Для поездки нужен билет."},
	{"context":"Поезд уходит с ___ номер три.","options":["train","platform","ticket"],"correct_index":1,"explanation":"Платформа — место посадки."},
	{"context":"В восемь утра отправляется ___.","options":["platform","ticket","train"],"correct_index":2,"explanation":"Train означает поезд."}]}`

func TestContextGameValidatesAndSeparatesSharedAndPrivateCache(t *testing.T) {
	api, store, first := newTestWebAPI(t)
	if _, err := store.getOrCreateUser(-43, "second"); err != nil {
		t.Fatal(err)
	}
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(-43, "second", hash); err != nil {
		t.Fatal(err)
	}
	if err := store.setLearningLanguage(-43, "en"); err != nil {
		t.Fatal(err)
	}
	if err := store.setInterfaceLanguage(-43, "ru"); err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	api.setSessionCookie(recorder, -43)
	second := recorder.Result().Cookies()[0]
	calls := 0
	api.bot.openrouter = newOpenRouterClient("test-key", "test-model", "", "", &http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			calls++
			body, _ := json.Marshal(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"content": sampleContextLesson}}}})
			return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(string(body)))}, nil
		}),
	})
	for _, cookie := range []*http.Cookie{first, second} {
		for _, body := range []map[string]string{{"preset": "travel"}, {"text": "The traveler needs a ticket."}, {"text": "The traveler needs a ticket."}} {
			result := requestJSON(t, api, cookie, http.MethodPost, "/api/games/context", body)
			if questions, ok := result["questions"].([]any); !ok || len(questions) != 3 {
				t.Fatalf("expected three questions: %v", result)
			}
		}
	}
	if calls != 3 {
		t.Fatalf("model calls = %d, want one shared and two private", calls)
	}
}

func TestContextGameRejectsDuplicateOptionsAndMissingAnswer(t *testing.T) {
	for _, raw := range []string{
		strings.Replace(sampleContextLesson, `["ticket","hotel","train"]`, `["ticket","Ticket","train"]`, 1),
		strings.Replace(sampleContextLesson, `"correct_index":0,`, "", 1),
	} {
		var lesson contextGameLesson
		if err := json.Unmarshal([]byte(raw), &lesson); err != nil {
			t.Fatal(err)
		}
		if validContextGame(lesson) {
			t.Fatal("accepted invalid exercise")
		}
	}
}
