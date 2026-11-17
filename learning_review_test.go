package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestWebOfflineReviewSyncDeduplicates(t *testing.T) {
	api, _, cookie := newTestWebAPI(t)
	for _, ids := range [][]string{{"card-a", "card-a"}, {"card-a", "card-b"}} {
		result := requestJSON(t, api, cookie, http.MethodPost, "/api/offline-reviews", map[string]any{"card_ids": ids})
		if result["synced"] != float64(1) {
			t.Fatalf("expected one new card per batch: %#v", result)
		}
	}
}

func TestLearningReviewWebURL(t *testing.T) {
	got, err := learningReviewWebURL("https://neriva.ru/app?lang=ru", "ai_tutor:lesson/1")
	if err != nil || got != "https://neriva.ru/app?lang=ru&review=ai_tutor%3Alesson%2F1&view=home" {
		t.Fatalf("review URL = %q, error = %v", got, err)
	}
}

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
	if _, err := storage.db.Exec(`CREATE TRIGGER reject_review_xp BEFORE UPDATE OF xp ON users
		WHEN NEW.xp > OLD.xp BEGIN SELECT RAISE(ABORT,'simulated review reward failure'); END`); err != nil {
		t.Fatal(err)
	}
	status, _ := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": entry.Correction})
	if status != http.StatusInternalServerError {
		t.Fatalf("forced failure status = %d", status)
	}
	if _, err := storage.db.Exec("DROP TRIGGER reject_review_xp"); err != nil {
		t.Fatal(err)
	}
	queue = requestJSON(t, api, cookie, http.MethodGet, "/api/learning-reviews", nil)
	if items, ok := queue["items"].([]any); !ok || len(items) != 1 {
		t.Fatalf("failed answer consumed its review: %#v", queue)
	}
	result := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": entry.Correction})
	if result["correct"] != true {
		t.Fatalf("correction rejected: %#v", result)
	}
	queue = requestJSON(t, api, cookie, http.MethodGet, "/api/learning-reviews", nil)
	if items, ok := queue["items"].([]any); !ok || len(items) != 0 {
		t.Fatalf("corrected mistake remained in the queue: %#v", queue)
	}
	status, _ = requestJSONRaw(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": entry.Correction})
	user, err := storage.getOrCreateUser(-42, "")
	if status != http.StatusConflict || err != nil || user.XP != 8 {
		t.Fatalf("replay status = %d, XP = %d, error = %v", status, user.XP, err)
	}
	next := mistakeEntry{Word: "I went yesterday?", Correction: entry.Correction, AddedAt: time.Now()}
	if err := storage.addMistakes(-42, "en", []mistakeEntry{next}); err != nil {
		t.Fatal(err)
	}
	requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/start", map[string]any{"index": 0})
	status, _ = requestJSONRaw(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer",
		map[string]any{"text": entry.Correction, "review_id": item["id"]})
	if status != http.StatusConflict {
		t.Fatalf("stale answer after opening another exercise returned %d", status)
	}
}
func TestWebMistakeTransferResultIsSavedOnce(t *testing.T) {
	api, storage, cookie := newTestWebAPI(t)
	body := map[string]any{"review_id": "mistake:" + strings.Repeat("a", 64)}
	status, _ := requestJSONRaw(t, api, cookie, http.MethodPost, "/api/mistakes/transfer", body)
	if status != http.StatusNotFound {
		t.Fatalf("unknown exercise: %d", status)
	}
	entry := mistakeEntry{Language: "en", Word: "I go yesterday", Correction: "I went yesterday", AddedAt: time.Now().UTC()}
	if err := storage.addMistakes(-42, "en", []mistakeEntry{entry}); err != nil {
		t.Fatal(err)
	}
	body["review_id"] = mistakeLearningReviewID(entry)
	requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/start", body)
	requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/practice/answer", map[string]any{"text": entry.Correction, "review_id": body["review_id"]})
	for i := 0; i < 2; i++ {
		result := requestJSON(t, api, cookie, http.MethodPost, "/api/mistakes/transfer", body)
		if result["completed"] != float64(1) || result["completed_last_7_days"] != float64(1) {
			t.Fatalf("completion count after write %d: %#v", i, result)
		}
	}
	result := requestJSON(t, api, cookie, http.MethodGet, "/api/mistakes/transfer", nil)
	if result["completed"] != float64(1) {
		t.Fatalf("stored result: %#v", result)
	}
}

func TestMistakeTransferSevenDayBoundary(t *testing.T) {
	_, sqlite, _ := newTestWebAPI(t)
	cutoff := time.Date(2026, 11, 10, 9, 30, 0, 0, time.UTC)
	jsonStorage := newTestJSONStore(t)
	user := newUserState(-42, cutoff)
	user.MistakeTransfers = map[string]mistakeTransferRecord{}
	for id, item := range map[string]mistakeTransferRecord{
		"old":            {Language: "en", CompletedAt: formatDBTime(cutoff.Add(-time.Nanosecond))},
		"edge":           {Language: "en", CompletedAt: formatDBTime(cutoff)},
		"recent":         {Language: "en", CompletedAt: formatDBTime(cutoff.Add(time.Nanosecond))},
		"pending":        {Language: "en"},
		"other-language": {Language: "de", CompletedAt: formatDBTime(cutoff)},
	} {
		user.MistakeTransfers[id] = item
		if _, err := sqlite.db.Exec("INSERT INTO mistake_transfers VALUES (?, ?, ?, ?)", -42, id, item.Language, item.CompletedAt); err != nil {
			t.Fatal(err)
		}
	}
	jsonStorage.users[-42] = user
	for name, storage := range map[string]store{"sqlite": sqlite, "json": jsonStorage} {
		t.Run(name, func(t *testing.T) {
			progress, err := storage.mistakeTransferProgress(-42, "en", cutoff)
			if err != nil || progress.Completed != 3 || progress.CompletedLast7Days != 2 {
				t.Fatalf("unexpected progress %+v: %v", progress, err)
			}
		})
	}
}
