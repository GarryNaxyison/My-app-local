package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebLearningReviewsReadExistingTutorSchedule(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	now := time.Now().UTC()
	for _, lesson := range []aiTutorLessonRecord{
		{ID: "older", LearningLanguage: "en", Theme: "Travel"},
		{ID: "next", LearningLanguage: "en", Theme: "Work"},
		{ID: "foreign", LearningLanguage: "de", Theme: "Deutsch"},
		{ID: "private", LearningLanguage: "en", Theme: "Private"},
	} {
		if err := store.saveAITutorLesson(lesson); err != nil {
			t.Fatal(err)
		}
	}
	for _, entry := range []struct {
		owner  int64
		id     string
		offset time.Duration
	}{
		{-42, "older", -time.Hour},
		{-42, "next", time.Hour},
		{-42, "foreign", -time.Hour},
		{-43, "private", -time.Hour},
	} {
		if err := store.scheduleAITutorReview(entry.owner, entry.id, formatDBTime(now.Add(entry.offset)), "tomorrow"); err != nil {
			t.Fatal(err)
		}
	}
	if err := store.scheduleAITutorReview(-42, "older", formatDBTime(now.Add(3*time.Hour)), "next_week"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec("UPDATE ai_tutor_reviews SET updated_at = ? WHERE lesson_id = ? AND interval_code = ?",
		formatDBTime(now.Add(-time.Hour)), "older", "next_week"); err != nil {
		t.Fatal(err)
	}
	response := requestJSON(t, api, cookie, http.MethodGet, "/api/learning-reviews", nil)
	items, ok := response["items"].([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("expected exactly two cards for this account and language: %#v", response)
	}
	first := items[0].(map[string]any)
	second := items[1].(map[string]any)
	if first["lesson_id"] != "older" || first["status"] != "due" || second["lesson_id"] != "next" || second["status"] != "upcoming" {
		t.Fatalf("unexpected queue order or status: %#v", items)
	}
}

func TestWebLearningReviewsRequireAuthentication(t *testing.T) {
	api, _, _ := newTestWebAPI(t)
	response := httptest.NewRecorder()
	api.handleLearningReviews(response, httptest.NewRequest(http.MethodGet, "/api/learning-reviews", nil))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", response.Code)
	}
}

func TestWebLearningReviewsIncludeMistakes(t *testing.T) {
	api, storage, cookie := newTestWebAPI(t)
	entry := mistakeEntry{Word: "I go yesterday", Correction: "I went yesterday", AddedAt: time.Now().Add(-time.Minute)}
	if err := storage.addMistakes(-42, "en", []mistakeEntry{entry}); err != nil {
		t.Fatal(err)
	}
	queue := requestJSON(t, api, cookie, http.MethodGet, "/api/learning-reviews", nil)
	items, ok := queue["items"].([]any)
	if !ok || len(items) != 1 {
		t.Fatalf("expected a review for the saved mistake: %#v", queue)
	}
	item := items[0].(map[string]any)
	if item["source"] != "mistake" || item["status"] != "due" || item["title"] != entry.Word {
		t.Fatalf("unexpected mistake review: %#v", item)
	}
	started := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/start", map[string]any{"review_id": item["id"]})
	if mistake, ok := started["mistake"].(map[string]any); !ok || mistake["word"] != entry.Word {
		t.Fatalf("wrong exercise: %#v", started)
	}
	result := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": entry.Correction})
	if result["correct"] != true {
		t.Fatalf("correction rejected: %#v", result)
	}
	queue = requestJSON(t, api, cookie, http.MethodGet, "/api/learning-reviews", nil)
	if items, ok := queue["items"].([]any); !ok || len(items) != 0 {
		t.Fatalf("corrected mistake remained in the queue: %#v", queue)
	}
}
