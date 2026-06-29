package main

import (
	"path/filepath"
	"testing"
)

func assertClassicLessonCountsOnCompletion(t *testing.T, store store, telegramID int64) {
	t.Helper()
	if _, err := store.getOrCreateUser(telegramID, "tester"); err != nil {
		t.Fatalf("getOrCreateUser() error = %v", err)
	}
	if err := store.saveLesson(telegramID, "Answer this lesson once."); err != nil {
		t.Fatalf("saveLesson() error = %v", err)
	}
	started, err := store.getOrCreateUser(telegramID, "tester")
	if err != nil {
		t.Fatalf("get started user: %v", err)
	}
	if started.LessonCount != 0 || started.LessonsToday != 0 {
		t.Fatalf("lesson counters after start = count:%d today:%d, want unchanged", started.LessonCount, started.LessonsToday)
	}
	if started.Mode != "lesson" || started.LastLessonPrompt == "" {
		t.Fatalf("active lesson was not stored after start: mode=%q prompt=%q", started.Mode, started.LastLessonPrompt)
	}
	if err := store.completeLesson(telegramID); err != nil {
		t.Fatalf("completeLesson() error = %v", err)
	}
	completed, err := store.getOrCreateUser(telegramID, "tester")
	if err != nil {
		t.Fatalf("get completed user: %v", err)
	}
	if completed.LessonCount != 1 || completed.LessonsToday != 1 {
		t.Fatalf("lesson counters after completion = count:%d today:%d, want 1/1", completed.LessonCount, completed.LessonsToday)
	}
	if completed.Mode != "idle" || completed.LastLessonPrompt != "" {
		t.Fatalf("active lesson was not cleared after completion: mode=%q prompt=%q", completed.Mode, completed.LastLessonPrompt)
	}
}

func TestJSONStoreClassicLessonCountsOnlyOnCompletion(t *testing.T) {
	assertClassicLessonCountsOnCompletion(t, newTestJSONStore(t), 101)
}

func TestSQLiteStoreClassicLessonCountsOnlyOnCompletion(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()
	assertClassicLessonCountsOnCompletion(t, store, 102)
}
