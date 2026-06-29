package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
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

func TestAITutorFindApprovedIgnoresOldFingerprintVersion(t *testing.T) {
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatalf("newJSONStore() error = %v", err)
	}
	payload := validAITutorLessonPayloadForTest()
	oldLesson := aiTutorLessonRecord{
		ID:                "old-approved",
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		ExactLevel:        "A1",
		LevelBand:         "A1-A2",
		Theme:             "daily life",
		Status:            aiTutorStatusApproved,
		Payload:           payload,
		Fingerprint:       "legacy-fingerprint",
		PostScore:         99,
	}
	currentLesson := oldLesson
	currentLesson.ID = "current-approved"
	currentLesson.Fingerprint = ""
	currentLesson.PostScore = 1
	if err := store.saveAITutorLesson(oldLesson); err != nil {
		t.Fatalf("save old lesson: %v", err)
	}
	if err := store.saveAITutorLesson(currentLesson); err != nil {
		t.Fatalf("save current lesson: %v", err)
	}
	found, ok, err := store.findApprovedAITutorLesson("en", "ru", "A1-A2", 42)
	if err != nil {
		t.Fatalf("findApprovedAITutorLesson() error = %v", err)
	}
	if !ok || found.ID != currentLesson.ID {
		t.Fatalf("found lesson = %#v ok=%v, want current approved lesson", found, ok)
	}
}

func TestSQLiteAITutorQualityPromotionWritesSeparateLessonBank(t *testing.T) {
	dir := t.TempDir()
	bankPath := filepath.Join(dir, "ai_tutor_lessons.sqlite")
	store, err := newSQLiteStore(filepath.Join(dir, "app.sqlite"), "", bankPath)
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()
	defer store.aiTutorDB.Close()

	lesson := aiTutorLessonRecord{
		ID:                "lesson-bank-1",
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		ExactLevel:        "A1",
		LevelBand:         "A1-A2",
		Theme:             "daily life",
		Status:            aiTutorStatusInTrial,
		Payload:           validAITutorLessonPayloadForTest(),
		PreflightScore:    92,
	}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if _, err := os.Stat(bankPath); err != nil {
		t.Fatalf("AI Tutor bank database was not created at %s: %v", bankPath, err)
	}
	var count int
	if err := store.aiTutorDB.QueryRow(`SELECT COUNT(*) FROM ai_tutor_lessons WHERE id = ?`, lesson.ID).Scan(&count); err != nil {
		t.Fatalf("count bank before promotion: %v", err)
	}
	if count != 0 {
		t.Fatalf("in-trial lesson should not be in approved bank, count=%d", count)
	}

	if err := store.updateAITutorLessonQuality(lesson.ID, aiTutorStatusApproved, 94); err != nil {
		t.Fatalf("updateAITutorLessonQuality() error = %v", err)
	}
	var status string
	var postScore int
	if err := store.aiTutorDB.QueryRow(`SELECT status, post_score FROM ai_tutor_lessons WHERE id = ?`, lesson.ID).Scan(&status, &postScore); err != nil {
		t.Fatalf("read promoted bank lesson: %v", err)
	}
	if status != aiTutorStatusApproved || postScore != 94 {
		t.Fatalf("bank lesson status=%q postScore=%d", status, postScore)
	}
}

func TestSQLiteAITutorFindApprovedUsesSeparateLessonBank(t *testing.T) {
	dir := t.TempDir()
	store, err := newSQLiteStore(filepath.Join(dir, "app.sqlite"), "", filepath.Join(dir, "ai_tutor_lessons.sqlite"))
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()
	defer store.aiTutorDB.Close()

	lesson := aiTutorLessonRecord{
		ID:                "lesson-bank-2",
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		ExactLevel:        "A1",
		LevelBand:         "A1-A2",
		Theme:             "daily life",
		Status:            aiTutorStatusApproved,
		Payload:           validAITutorLessonPayloadForTest(),
		PreflightScore:    92,
		PostScore:         95,
	}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if _, err := store.db.Exec(`DELETE FROM ai_tutor_lessons WHERE id = ?`, lesson.ID); err != nil {
		t.Fatalf("delete main lesson row: %v", err)
	}

	found, ok, err := store.findApprovedAITutorLesson("en", "ru", "A1-A2", 42)
	if err != nil {
		t.Fatalf("findApprovedAITutorLesson() error = %v", err)
	}
	if !ok || found.ID != lesson.ID {
		t.Fatalf("found lesson = %#v ok=%v", found, ok)
	}
	var mainCount int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM ai_tutor_lessons WHERE id = ?`, lesson.ID).Scan(&mainCount); err != nil {
		t.Fatalf("count restored main lesson: %v", err)
	}
	if mainCount != 1 {
		t.Fatalf("bank lesson should be copied into main DB for active sessions, mainCount=%d", mainCount)
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

func TestSQLiteStoreAITutorWordReportLifecycle(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()

	lesson := aiTutorLessonRecord{
		ID:                "lesson-word-report-1",
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		ExactLevel:        "A1",
		LevelBand:         "A1-A2",
		Status:            aiTutorStatusApproved,
		Payload:           validAITutorLessonPayloadForTest(),
	}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{
		ID:           "session-word-report-1",
		TelegramID:   42,
		LessonID:     lesson.ID,
		Surface:      "web",
		CurrentStage: aiTutorWordLearnStage(1),
		Status:       aiTutorSessionActive,
	}); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}

	created, duplicate, err := store.createAITutorWordReport(aiTutorWordReportRecord{
		ID:                  "aitwr-1",
		TelegramID:          42,
		SessionID:           "session-word-report-1",
		LessonID:            lesson.ID,
		Stage:               aiTutorWordLearnStage(1),
		WordIndex:           0,
		OriginalWord:        "wake up",
		OriginalTranslation: "prosypatsya",
		ProposedWord:        "get up",
		ProposedTranslation: "vstavat",
		Comment:             "better phrase for this card",
		Status:              aiTutorWordReportPending,
	})
	if err != nil {
		t.Fatalf("createAITutorWordReport() error = %v", err)
	}
	if duplicate {
		t.Fatal("first report should not be marked duplicate")
	}
	if created.ID != "aitwr-1" || created.Status != aiTutorWordReportPending || created.CreatedAt == "" || created.UpdatedAt == "" {
		t.Fatalf("created report mismatch: %#v", created)
	}

	if err := store.setAITutorWordReportAdminMessage(created.ID, "1001", 77); err != nil {
		t.Fatalf("setAITutorWordReportAdminMessage() error = %v", err)
	}
	if err := store.setAITutorWordReportFixPrompt(created.ID, "1001", 88); err != nil {
		t.Fatalf("setAITutorWordReportFixPrompt() error = %v", err)
	}
	byPrompt, ok, err := store.findPendingAITutorWordReportByFixPrompt("1001", 88)
	if err != nil || !ok || byPrompt.ID != created.ID {
		t.Fatalf("findPendingAITutorWordReportByFixPrompt() = %#v ok=%v err=%v", byPrompt, ok, err)
	}

	userCount, sessionCount, err := store.countAITutorWordReports(42, "session-word-report-1", time.Now().UTC().Add(-10*time.Minute))
	if err != nil {
		t.Fatalf("countAITutorWordReports() error = %v", err)
	}
	if userCount != 1 || sessionCount != 1 {
		t.Fatalf("report counts user=%d session=%d, want 1/1", userCount, sessionCount)
	}

	resolved, applied, err := store.resolveAITutorWordReport(created.ID, aiTutorWordReportAccepted, "get up", "vstavat")
	if err != nil {
		t.Fatalf("resolveAITutorWordReport() error = %v", err)
	}
	if !applied || resolved.Status != aiTutorWordReportAccepted || resolved.FinalWord != "get up" || resolved.FinalTranslation != "vstavat" || resolved.ResolvedAt == "" {
		t.Fatalf("resolved report mismatch applied=%v report=%#v", applied, resolved)
	}
	resolvedAgain, appliedAgain, err := store.resolveAITutorWordReport(created.ID, aiTutorWordReportRejected, "ignored", "ignored")
	if err != nil {
		t.Fatalf("resolveAITutorWordReport(second) error = %v", err)
	}
	if appliedAgain || resolvedAgain.Status != aiTutorWordReportAccepted || resolvedAgain.FinalWord != "get up" {
		t.Fatalf("resolve should be idempotent applied=%v report=%#v", appliedAgain, resolvedAgain)
	}
}

func TestSQLiteStoreAITutorWordReportDeduplicatesPendingStage(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()

	report := aiTutorWordReportRecord{
		ID:                  "aitwr-dup-1",
		TelegramID:          42,
		SessionID:           "session-dup",
		LessonID:            "lesson-dup",
		Stage:               aiTutorWordLearnStage(2),
		WordIndex:           1,
		OriginalWord:        "water",
		OriginalTranslation: "voda",
		ProposedWord:        "still water",
		ProposedTranslation: "voda bez gaza",
		Status:              aiTutorWordReportPending,
	}
	first, duplicate, err := store.createAITutorWordReport(report)
	if err != nil {
		t.Fatalf("createAITutorWordReport(first) error = %v", err)
	}
	if duplicate {
		t.Fatal("first insert should not be duplicate")
	}
	report.ID = "aitwr-dup-2"
	report.ProposedWord = "sparkling water"
	second, duplicate, err := store.createAITutorWordReport(report)
	if err != nil {
		t.Fatalf("createAITutorWordReport(second) error = %v", err)
	}
	if !duplicate || second.ID != first.ID || second.ProposedWord != first.ProposedWord {
		t.Fatalf("duplicate result = %#v duplicate=%v, want original %#v", second, duplicate, first)
	}
	globalCount, err := store.countAITutorWordReportsSince(time.Now().UTC().Add(-10 * time.Minute))
	if err != nil {
		t.Fatalf("countAITutorWordReportsSince() error = %v", err)
	}
	if globalCount != 1 {
		t.Fatalf("global report count = %d, want 1", globalCount)
	}
}

func TestApplyAITutorWordReportCorrectionUpdatesLessonAndVocabulary(t *testing.T) {
	configureSQLiteVocabularyForTest(t, []vocabWord{
		{
			ID:           "w1",
			Language:     "en",
			English:      "wake up",
			Russian:      "prosypatsya",
			Translations: map[string]string{"ru": "prosypatsya"},
			Level:        "A1",
		},
	})
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	defer store.db.Close()
	lessonPayload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{
		ID:                "lesson-apply-report",
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		ExactLevel:        "A1",
		LevelBand:         "A1-A2",
		Status:            aiTutorStatusApproved,
		Payload:           lessonPayload,
	}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	before, ok, err := store.getAITutorLesson(lesson.ID)
	if err != nil || !ok {
		t.Fatalf("getAITutorLesson(before) ok=%v err=%v", ok, err)
	}
	report := aiTutorWordReportRecord{
		ID:                  "aitwr-apply-1",
		TelegramID:          42,
		SessionID:           "session-apply-report",
		LessonID:            lesson.ID,
		Stage:               aiTutorWordLearnStage(1),
		WordIndex:           0,
		OriginalWord:        "wake up",
		OriginalTranslation: "prosypatsya",
		ProposedWord:        "get up",
		ProposedTranslation: "vstavat",
		Status:              aiTutorWordReportPending,
	}

	if err := applyAITutorWordReportCorrection(store, report, "get up", "vstavat"); err != nil {
		t.Fatalf("applyAITutorWordReportCorrection() error = %v", err)
	}
	after, ok, err := store.getAITutorLesson(lesson.ID)
	if err != nil || !ok {
		t.Fatalf("getAITutorLesson(after) ok=%v err=%v", ok, err)
	}
	if after.Payload.Words[0].ID != "w1" {
		t.Fatalf("word id changed: %#v", after.Payload.Words[0])
	}
	if after.Payload.Words[0].Target != "get up" || after.Payload.Words[0].InterfaceTranslation != "vstavat" || after.Payload.Words[0].AudioTextTarget != "get up" {
		t.Fatalf("lesson word was not corrected: %#v", after.Payload.Words[0])
	}
	if before.Fingerprint == after.Fingerprint {
		t.Fatalf("fingerprint did not change after word correction: %q", after.Fingerprint)
	}

	db := currentSQLiteVocabularyDB()
	if db == nil {
		t.Fatal("expected configured SQLite vocabulary db")
	}
	var word, russian string
	if err := db.QueryRow(`SELECT word, russian FROM vocabulary_words WHERE id = ?`, "w1").Scan(&word, &russian); err != nil {
		t.Fatalf("read vocabulary_words: %v", err)
	}
	if word != "get up" || russian != "vstavat" {
		t.Fatalf("vocabulary_words = %q/%q, want get up/vstavat", word, russian)
	}
	var translated string
	if err := db.QueryRow(`SELECT text FROM vocabulary_translations WHERE word_id = ? AND language = ?`, "w1", "ru").Scan(&translated); err != nil {
		t.Fatalf("read vocabulary_translations: %v", err)
	}
	if translated != "vstavat" {
		t.Fatalf("vocabulary translation = %q, want vstavat", translated)
	}
	if err := db.QueryRow(`SELECT word, russian FROM vocabulary_ai_words WHERE id = ?`, "w1").Scan(&word, &russian); err != nil {
		t.Fatalf("read vocabulary_ai_words: %v", err)
	}
	if word != "get up" || russian != "vstavat" {
		t.Fatalf("vocabulary_ai_words = %q/%q, want get up/vstavat", word, russian)
	}
	var sourceWord string
	if err := db.QueryRow(`SELECT source_word, translation FROM vocabulary_ai_translations WHERE word_id = ? AND target_language = ?`, "w1", "ru").Scan(&sourceWord, &translated); err != nil {
		t.Fatalf("read vocabulary_ai_translations: %v", err)
	}
	if sourceWord != "get up" || translated != "vstavat" {
		t.Fatalf("vocabulary_ai_translations = %q/%q, want get up/vstavat", sourceWord, translated)
	}
}
