package main

import (
	"path/filepath"
	"testing"
)

func TestSQLiteAITutorLessonAndSessionRoundTrip(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()
	lesson := aiTutorLessonRecord{
		ID:                "lesson-1",
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		ExactLevel:        "A1",
		LevelBand:         "A1-A2",
		Theme:             "daily life",
		Status:            aiTutorStatusApproved,
		Payload:           validAITutorLessonPayloadForTest(),
		Fingerprint:       "fp-1",
		PreflightScore:    92,
	}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	found, ok, err := store.findApprovedAITutorLesson("en", "ru", "A1-A2", 42)
	if err != nil {
		t.Fatalf("findApprovedAITutorLesson() error = %v", err)
	}
	if !ok || found.ID != lesson.ID {
		t.Fatalf("found lesson = %#v ok=%v", found, ok)
	}
	session := aiTutorSessionRecord{
		ID:           "session-1",
		TelegramID:   42,
		LessonID:     lesson.ID,
		Surface:      "web",
		CurrentStage: aiTutorStageStoryIntro,
		Status:       aiTutorSessionActive,
	}
	if err := store.createAITutorSession(session); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}
	loaded, ok, err := store.getAITutorSession("session-1")
	if err != nil || !ok {
		t.Fatalf("getAITutorSession() loaded=%#v ok=%v err=%v", loaded, ok, err)
	}
	if loaded.CurrentStage != aiTutorStageStoryIntro || loaded.TelegramID != 42 {
		t.Fatalf("loaded session mismatch: %#v", loaded)
	}
}

func TestSQLiteAITutorMigrationClearsLegacyTutorBase(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()
	lesson, err := buildTutorLessonForSequence(tutorReusableLessonUser(userState{TelegramID: 1, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}), 0)
	if err != nil {
		t.Fatalf("legacy lesson build error = %v", err)
	}
	tx, err := store.db.Begin()
	if err != nil {
		t.Fatalf("begin error = %v", err)
	}
	if err := store.saveTutorLessonTx(tx, lesson); err != nil {
		t.Fatalf("save legacy lesson error = %v", err)
	}
	if _, err := tx.Exec(`INSERT INTO tutor_user_lessons (telegram_id, lesson_id, assigned_at, lesson_order) VALUES (1, ?, 'now', 1)`, lesson.ID); err != nil {
		t.Fatalf("save legacy assignment error = %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit error = %v", err)
	}
	if err := store.clearLegacyTutorLessonBase(); err != nil {
		t.Fatalf("clearLegacyTutorLessonBase() error = %v", err)
	}
	var lessons, assignments int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM tutor_lessons`).Scan(&lessons); err != nil {
		t.Fatalf("count tutor_lessons error = %v", err)
	}
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM tutor_user_lessons`).Scan(&assignments); err != nil {
		t.Fatalf("count tutor_user_lessons error = %v", err)
	}
	if lessons != 0 || assignments != 0 {
		t.Fatalf("legacy rows remain: lessons=%d assignments=%d", lessons, assignments)
	}
}
