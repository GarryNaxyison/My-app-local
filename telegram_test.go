package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestSendInlineMessageEditsCallbackMessage(t *testing.T) {
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

	client := &telegramClient{baseURL: server.URL, http: server.Client()}
	ctx := contextWithTelegramEditTarget(context.Background(), 42, 77)

	if err := client.sendInlineMessage(ctx, 42, "menu", mainMenuInlineKeyboard()); err != nil {
		t.Fatalf("sendInlineMessage() error = %v", err)
	}
	if method != "/editMessageText" {
		t.Fatalf("expected editMessageText, got %s", method)
	}
	if payload["chat_id"].(float64) != 42 || payload["message_id"].(float64) != 77 {
		t.Fatalf("unexpected edit target payload: %#v", payload)
	}
}

func TestMainMenuContainsCoreBotFunctions(t *testing.T) {
	keyboard := mainMenuInlineKeyboard(ui(userState{InterfaceLanguage: "ru"}))
	callbacks := collectCallbackData(t, keyboard)
	for _, want := range []string{
		"menu_lesson",
		"menu_tutor",
		"menu_practice",
		"menu_word_lesson",
		"menu_word_game",
		"menu_spelling",
		"menu_vocabulary",
		"menu_phrasebook",
		"menu_mistakes",
		"menu_level_test",
		"menu_progress",
		"menu_leaderboard",
		"menu_limits",
		"menu_tools",
		"menu_reminders",
		"menu_settings",
		"menu_premium",
		"menu_referral",
	} {
		if !callbacks[want] {
			t.Fatalf("main menu is missing callback %q; got %#v", want, callbacks)
		}
	}
}

func TestAITutorTelegramRecallKeyboardUsesStructuredOptions(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()
	step := aiTutorBuildStep(lesson, aiTutorWordRecallStage(1))
	keyboard := aiTutorTelegramKeyboard("session-1", step, ui(userState{InterfaceLanguage: "ru"}))
	callbacks := collectCallbackData(t, keyboard)

	for _, want := range []string{"ait|session-1|choice|w1", "ait|session-1|choice|w2", "back_menu"} {
		if !callbacks[want] {
			t.Fatalf("recall keyboard is missing callback %q; got %#v", want, callbacks)
		}
	}
	if callbacks["ait|session-1|choice|continue"] {
		t.Fatalf("recall keyboard must use answer choices, got continue callback: %#v", callbacks)
	}
}

func TestAITutorTelegramAudioClipsUseGeneratedAudioText(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()

	storyClips := aiTutorTelegramAudioClips(aiTutorBuildStep(lesson, aiTutorStageStoryIntro))
	if len(storyClips) != 1 || storyClips[0].Text != lesson.Story.AudioTextTarget {
		t.Fatalf("story clips = %#v", storyClips)
	}

	wordClips := aiTutorTelegramAudioClips(aiTutorBuildStep(lesson, aiTutorWordLearnStage(1)))
	if len(wordClips) != 2 {
		t.Fatalf("word clips count = %d, want 2: %#v", len(wordClips), wordClips)
	}
	if wordClips[0].Text != lesson.Words[0].AudioTextTarget || wordClips[1].Text != lesson.Words[0].ExampleAudioTextTarget {
		t.Fatalf("word clips use wrong audio text: %#v", wordClips)
	}

	recallClips := aiTutorTelegramAudioClips(aiTutorBuildStep(lesson, aiTutorWordRecallStage(1)))
	if len(recallClips) != 0 {
		t.Fatalf("recall clips should not reveal the answer: %#v", recallClips)
	}
}

func TestTelegramAITutorCallbackDeletesPreviousAudioMessages(t *testing.T) {
	var deleted []int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if r.URL.Path == "/deleteMessage" {
			id, _ := payload["message_id"].(float64)
			deleted = append(deleted, int64(id))
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	store := &jsonStore{path: filepath.Join(t.TempDir(), "users.json"), users: map[int64]userState{}, aiTutorLessons: map[string]aiTutorLessonRecord{}, aiTutorSessions: map[string]aiTutorSessionRecord{}, aiTutorAnswers: map[string]aiTutorAnswerRecord{}}
	user := userState{TelegramID: 123, FirstName: "Test", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1", Mode: "ai_tutor:session-1"}
	store.users[user.TelegramID] = user
	_ = store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-audio-cleanup", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()})
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: user.TelegramID, LessonID: "lesson-tg-audio-cleanup", Surface: "telegram", CurrentStage: aiTutorStageStoryIntro, Status: aiTutorSessionActive})
	b := &bot{
		store:                 store,
		telegram:              &telegramClient{baseURL: server.URL, http: server.Client()},
		pronunciationMessages: map[int64][]int64{user.TelegramID: []int64{901, 902}},
	}

	err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:      "cb-1",
		From:    telegramUser{ID: user.TelegramID, FirstName: user.FirstName},
		Message: &telegramMessage{MessageID: 77, Chat: telegramChat{ID: user.TelegramID}},
		Data:    "ait|session-1|choice|continue",
	})
	if err != nil {
		t.Fatalf("handleCallbackQuery() error = %v", err)
	}
	if len(deleted) != 2 || deleted[0] != 901 || deleted[1] != 902 {
		t.Fatalf("deleted audio message IDs = %#v, want [901 902]", deleted)
	}
	if remaining := b.pronunciationMessages[user.TelegramID]; len(remaining) != 0 {
		t.Fatalf("pronunciationMessages still has chat entries: %#v", b.pronunciationMessages)
	}
}

func TestMistakesActionsKeyboardUsesDistinctActionLabels(t *testing.T) {
	english := englishUICopy()
	keyboard := mistakesActionsKeyboard(0, 1, english)

	practiceText := buttonTextByCallback(t, keyboard, "practice_mistakes")
	if !strings.Contains(practiceText, "Fix") || strings.Contains(practiceText, english.Mistakes) {
		t.Fatalf("practice_mistakes text = %q, want a fix-action label distinct from %q", practiceText, english.Mistakes)
	}
	clearText := buttonTextByCallback(t, keyboard, "clear_mistakes")
	if !strings.Contains(clearText, "Clear") || strings.Contains(clearText, english.Mistakes) {
		t.Fatalf("clear_mistakes text = %q, want a clear-action label distinct from %q", clearText, english.Mistakes)
	}

	russian := ui(userState{InterfaceLanguage: "ru"})
	russianKeyboard := mistakesActionsKeyboard(0, 1, russian)
	if got := buttonTextByCallback(t, russianKeyboard, "practice_mistakes"); strings.Contains(got, russian.Mistakes) {
		t.Fatalf("Russian practice_mistakes text still repeats the dictionary title: %q", got)
	}
	if got := buttonTextByCallback(t, russianKeyboard, "clear_mistakes"); strings.Contains(got, russian.Mistakes) {
		t.Fatalf("Russian clear_mistakes text still repeats the dictionary title: %q", got)
	}
}

func TestTelegramAITutorMenuStartsInteractiveSession(t *testing.T) {
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

	store := &jsonStore{
		path:            filepath.Join(t.TempDir(), "users.json"),
		users:           map[int64]userState{},
		aiTutorLessons:  map[string]aiTutorLessonRecord{},
		aiTutorSessions: map[string]aiTutorSessionRecord{},
		aiTutorAnswers:  map[string]aiTutorAnswerRecord{},
	}
	user := userState{
		TelegramID:        123,
		FirstName:         "Test",
		InterfaceLanguage: "ru",
		InterfaceSelected: true,
		TimezoneSelected:  true,
		LanguageSelected:  true,
		LearningLanguage:  "en",
		Level:             "A1",
	}
	store.users[user.TelegramID] = user
	if err := store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}

	err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:   "cb-1",
		From: telegramUser{ID: user.TelegramID, FirstName: user.FirstName},
		Data: "menu_tutor",
	})
	if err != nil {
		t.Fatalf("handleCallbackQuery(menu_tutor) error = %v", err)
	}
	refreshed := store.users[user.TelegramID]
	if !strings.HasPrefix(refreshed.Mode, "ai_tutor:") {
		t.Fatalf("mode = %q, want ai_tutor session", refreshed.Mode)
	}
	text, _ := payloads[len(payloads)-1]["text"].(string)
	if !strings.Contains(text, "A Morning Visit") {
		t.Fatalf("expected story step title, got %q", text)
	}
}

func TestMenuTutorDoesNotUseLocalCourseBank(t *testing.T) {
	var sentText string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		sentText, _ = payload["text"].(string)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	store := &jsonStore{
		path:            filepath.Join(t.TempDir(), "users.json"),
		users:           map[int64]userState{},
		aiTutorLessons:  map[string]aiTutorLessonRecord{},
		aiTutorSessions: map[string]aiTutorSessionRecord{},
		aiTutorAnswers:  map[string]aiTutorAnswerRecord{},
	}
	user := userState{
		TelegramID:        123,
		FirstName:         "Test",
		InterfaceLanguage: "ru",
		InterfaceSelected: true,
		TimezoneSelected:  true,
		LanguageSelected:  true,
		LearningLanguage:  "en",
		Level:             "A1",
	}
	store.users[user.TelegramID] = user
	if err := store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-local-regression", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}

	err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:   "cb-1",
		From: telegramUser{ID: user.TelegramID, FirstName: user.FirstName},
		Data: "menu_tutor",
	})
	if err != nil {
		t.Fatalf("handleCallbackQuery(menu_tutor) error = %v", err)
	}

	if !strings.Contains(sentText, "A Morning Visit") {
		t.Fatalf("menu_tutor did not render AI tutor story step: %q", sentText)
	}
	for _, oldMarker := range []string{"local-a1-a2-course-core", "Final word check", "Pattern:"} {
		if strings.Contains(sentText, oldMarker) {
			t.Fatalf("menu_tutor leaked local course bank marker %q in text: %q", oldMarker, sentText)
		}
	}
}

func TestTelegramAITutorTextRoutesToActiveSession(t *testing.T) {
	store := &jsonStore{path: filepath.Join(t.TempDir(), "users.json"), users: map[int64]userState{}, aiTutorLessons: map[string]aiTutorLessonRecord{}, aiTutorSessions: map[string]aiTutorSessionRecord{}, aiTutorAnswers: map[string]aiTutorAnswerRecord{}}
	user := userState{TelegramID: 123, FirstName: "Test", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1", Mode: "ai_tutor:session-1"}
	store.users[user.TelegramID] = user
	_ = store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-2", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()})
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: user.TelegramID, LessonID: "lesson-tg-2", Surface: "telegram", CurrentStage: aiTutorStageStoryIntro, Status: aiTutorSessionActive})
	var sent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		_ = json.NewDecoder(r.Body).Decode(&payload)
		sent, _ = payload["text"].(string)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	if err := b.handleAITutorText(context.Background(), user.TelegramID, user, "continue"); err != nil {
		t.Fatalf("handleAITutorText() error = %v", err)
	}
	session, _, _ := store.getAITutorSession("session-1")
	if session.CurrentStage != aiTutorStageRetell {
		t.Fatalf("stage = %q, want retell", session.CurrentStage)
	}
	if !strings.Contains(sent, "Retell") {
		t.Fatalf("sent text = %q", sent)
	}
}

func TestTelegramAITutorCallbackRoutesToActiveSession(t *testing.T) {
	store := &jsonStore{path: filepath.Join(t.TempDir(), "users.json"), users: map[int64]userState{}, aiTutorLessons: map[string]aiTutorLessonRecord{}, aiTutorSessions: map[string]aiTutorSessionRecord{}, aiTutorAnswers: map[string]aiTutorAnswerRecord{}}
	user := userState{TelegramID: 123, FirstName: "Test", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1", Mode: "ai_tutor:session-1"}
	store.users[user.TelegramID] = user
	_ = store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-3", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()})
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: user.TelegramID, LessonID: "lesson-tg-3", Surface: "telegram", CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`{"ok":true}`)) }))
	defer server.Close()
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	err := b.handleCallbackQuery(context.Background(), callbackQuery{ID: "cb-1", From: telegramUser{ID: user.TelegramID, FirstName: user.FirstName}, Data: "ait|session-1|choice|no_review"})
	if err != nil {
		t.Fatalf("handleCallbackQuery() error = %v", err)
	}
	session, _, _ := store.getAITutorSession("session-1")
	if session.Status != aiTutorSessionComplete {
		t.Fatalf("status = %q, want complete", session.Status)
	}
}

func TestTelegramAITutorWordReportAcceptAppliesCorrection(t *testing.T) {
	store, report := newAITutorWordReportModerationFixture(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":91}}`))
	}))
	defer server.Close()
	b := &bot{
		cfg:      config{TelegramOpsRecipients: []telegramOpsRecipient{{ChatID: "999"}}},
		store:    store,
		telegram: &telegramClient{baseURL: server.URL, http: server.Client()},
	}

	if err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:      "cb-report-accept",
		From:    telegramUser{ID: 999, FirstName: "Admin"},
		Message: &telegramMessage{MessageID: 77, Chat: telegramChat{ID: 999}},
		Data:    "ait_word_report|" + report.ID + "|accept",
	}); err != nil {
		t.Fatalf("handleCallbackQuery(accept) error = %v", err)
	}

	resolved, ok, err := store.getAITutorWordReport(report.ID)
	if err != nil || !ok {
		t.Fatalf("getAITutorWordReport() ok=%v err=%v", ok, err)
	}
	if resolved.Status != aiTutorWordReportAccepted || resolved.FinalWord != report.ProposedWord || resolved.FinalTranslation != report.ProposedTranslation {
		t.Fatalf("resolved report mismatch: %#v", resolved)
	}
	lesson, ok, err := store.getAITutorLesson(report.LessonID)
	if err != nil || !ok {
		t.Fatalf("getAITutorLesson() ok=%v err=%v", ok, err)
	}
	if got := lesson.Payload.Words[0].Target; got != report.ProposedWord {
		t.Fatalf("lesson word target = %q, want %q", got, report.ProposedWord)
	}
	if got := lesson.Payload.Words[0].InterfaceTranslation; got != report.ProposedTranslation {
		t.Fatalf("lesson word translation = %q, want %q", got, report.ProposedTranslation)
	}
}

func TestTelegramAITutorWordReportRejectDoesNotMutateLesson(t *testing.T) {
	store, report := newAITutorWordReportModerationFixture(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	b := &bot{
		cfg:      config{TelegramOpsRecipients: []telegramOpsRecipient{{ChatID: "999"}}},
		store:    store,
		telegram: &telegramClient{baseURL: server.URL, http: server.Client()},
	}

	if err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:      "cb-report-reject",
		From:    telegramUser{ID: 999, FirstName: "Admin"},
		Message: &telegramMessage{MessageID: 77, Chat: telegramChat{ID: 999}},
		Data:    "ait_word_report|" + report.ID + "|reject",
	}); err != nil {
		t.Fatalf("handleCallbackQuery(reject) error = %v", err)
	}

	resolved, _, _ := store.getAITutorWordReport(report.ID)
	if resolved.Status != aiTutorWordReportRejected {
		t.Fatalf("status = %q, want rejected", resolved.Status)
	}
	lesson, _, _ := store.getAITutorLesson(report.LessonID)
	if got := lesson.Payload.Words[0].Target; got != report.OriginalWord {
		t.Fatalf("reject mutated lesson word target = %q, want %q", got, report.OriginalWord)
	}
}

func TestTelegramAITutorWordReportFixReplyAppliesCorrection(t *testing.T) {
	store, report := newAITutorWordReportModerationFixture(t)
	var promptMessageID int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		_ = json.NewDecoder(r.Body).Decode(&payload)
		if text, _ := payload["text"].(string); strings.Contains(text, "word - translation") {
			promptMessageID = 91
		}
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":91}}`))
	}))
	defer server.Close()
	b := &bot{
		cfg:      config{TelegramOpsRecipients: []telegramOpsRecipient{{ChatID: "999"}}},
		store:    store,
		telegram: &telegramClient{baseURL: server.URL, http: server.Client()},
	}

	if err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:      "cb-report-fix",
		From:    telegramUser{ID: 999, FirstName: "Admin"},
		Message: &telegramMessage{MessageID: 77, Chat: telegramChat{ID: 999}},
		Data:    "ait_word_report|" + report.ID + "|fix",
	}); err != nil {
		t.Fatalf("handleCallbackQuery(fix) error = %v", err)
	}
	if promptMessageID == 0 {
		t.Fatal("expected fix prompt message")
	}

	if err := b.handleUpdate(context.Background(), telegramUpdate{Message: &telegramMessage{
		MessageID: 92,
		From:      telegramUser{ID: 999, FirstName: "Admin"},
		Chat:      telegramChat{ID: 999},
		Text:      "rise - podnimatsya",
		ReplyToMessage: &telegramMessage{
			MessageID: promptMessageID,
			Chat:      telegramChat{ID: 999},
		},
	}}); err != nil {
		t.Fatalf("handleUpdate(fix reply) error = %v", err)
	}

	resolved, _, _ := store.getAITutorWordReport(report.ID)
	if resolved.Status != aiTutorWordReportFixed || resolved.FinalWord != "rise" || resolved.FinalTranslation != "podnimatsya" {
		t.Fatalf("fixed report mismatch: %#v", resolved)
	}
	lesson, _, _ := store.getAITutorLesson(report.LessonID)
	if got := lesson.Payload.Words[0].Target; got != "rise" {
		t.Fatalf("fixed lesson word target = %q, want rise", got)
	}
}

func TestTelegramAITutorWordReportFixPlainAdminMessageAppliesLatestPrompt(t *testing.T) {
	store, report := newAITutorWordReportModerationFixture(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true,"result":{"message_id":91}}`))
	}))
	defer server.Close()
	b := &bot{
		cfg:      config{TelegramOpsRecipients: []telegramOpsRecipient{{ChatID: "999"}}},
		store:    store,
		telegram: &telegramClient{baseURL: server.URL, http: server.Client()},
	}

	if err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:      "cb-report-fix",
		From:    telegramUser{ID: 999, FirstName: "Admin"},
		Message: &telegramMessage{MessageID: 77, Chat: telegramChat{ID: 999}},
		Data:    "ait_word_report|" + report.ID + "|fix",
	}); err != nil {
		t.Fatalf("handleCallbackQuery(fix) error = %v", err)
	}

	if err := b.handleUpdate(context.Background(), telegramUpdate{Message: &telegramMessage{
		MessageID: 92,
		From:      telegramUser{ID: 999, FirstName: "Admin"},
		Chat:      telegramChat{ID: 999},
		Text:      "local park - местный парк",
	}}); err != nil {
		t.Fatalf("handleUpdate(plain fix message) error = %v", err)
	}

	resolved, _, _ := store.getAITutorWordReport(report.ID)
	if resolved.Status != aiTutorWordReportFixed || resolved.FinalWord != "local park" || resolved.FinalTranslation != "местный парк" {
		t.Fatalf("fixed report mismatch: %#v", resolved)
	}
	lesson, _, _ := store.getAITutorLesson(report.LessonID)
	if got := lesson.Payload.Words[0].Target; got != "local park" {
		t.Fatalf("fixed lesson word target = %q, want local park", got)
	}
}

func TestTelegramAITutorWordReportAdminOnly(t *testing.T) {
	store, report := newAITutorWordReportModerationFixture(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	b := &bot{
		cfg:      config{TelegramOpsRecipients: []telegramOpsRecipient{{ChatID: "999"}}},
		store:    store,
		telegram: &telegramClient{baseURL: server.URL, http: server.Client()},
	}

	if err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:      "cb-report-intruder",
		From:    telegramUser{ID: 123, FirstName: "Not Admin"},
		Message: &telegramMessage{MessageID: 77, Chat: telegramChat{ID: 123}},
		Data:    "ait_word_report|" + report.ID + "|accept",
	}); err != nil {
		t.Fatalf("handleCallbackQuery(non-admin) error = %v", err)
	}
	stored, _, _ := store.getAITutorWordReport(report.ID)
	if stored.Status != aiTutorWordReportPending {
		t.Fatalf("non-admin changed report status to %q", stored.Status)
	}
}

func newAITutorWordReportModerationFixture(t *testing.T) (*sqliteStore, aiTutorWordReportRecord) {
	t.Helper()
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "telegram-report.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	t.Cleanup(func() { _ = store.db.Close() })
	lesson := aiTutorLessonRecord{ID: "lesson-report-telegram", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{ID: "session-report-telegram", TelegramID: 42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorWordLearnStage(1), Status: aiTutorSessionActive}); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}
	report, _, err := store.createAITutorWordReport(aiTutorWordReportRecord{
		ID:                  "aitwr-telegram-1",
		TelegramID:          42,
		SessionID:           "session-report-telegram",
		LessonID:            lesson.ID,
		Stage:               aiTutorWordLearnStage(1),
		WordIndex:           0,
		OriginalWord:        "wake up",
		OriginalTranslation: "prosypatsya",
		ProposedWord:        "get up",
		ProposedTranslation: "vstavat",
		Status:              aiTutorWordReportPending,
	})
	if err != nil {
		t.Fatalf("createAITutorWordReport() error = %v", err)
	}
	return store, report
}

func TestFormatTelegramTutorLessonIncludesTeachingPointAndFinalCheck(t *testing.T) {
	user := userState{InterfaceLanguage: "en", LearningLanguage: "en", Level: "A1"}
	lesson, err := buildTutorLessonForSequence(user, 4)
	if err != nil {
		t.Fatalf("buildTutorLessonForSequence(doctor): %v", err)
	}
	text := formatTelegramTutorLesson(lesson, user)
	for _, want := range []string{
		"Hello\\. I need to see a doctor today, please\\.",
		"Final word check",
		"doctor",
		"today",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("telegram tutor lesson missing %q:\n%s", want, text)
		}
	}
}

func TestAITutorMenuLabelIsLocalizedForEveryInterfaceLanguage(t *testing.T) {
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code, InterfaceSelected: true})
		if strings.TrimSpace(copy.AITutor) == "" {
			t.Fatalf("missing AI Tutor label for %s", language.Code)
		}
		keyboard := mainMenuInlineKeyboard(copy)
		rows := keyboard["inline_keyboard"].([][]map[string]any)
		found := false
		for _, row := range rows {
			for _, button := range row {
				if button["callback_data"] == "menu_tutor" {
					found = true
					text, _ := button["text"].(string)
					if !strings.Contains(text, copy.AITutor) {
						t.Fatalf("menu_tutor label for %s = %q, want copy label %q", language.Code, text, copy.AITutor)
					}
					if language.Code != "ru" && strings.Contains(text, "AI Репетитор") {
						t.Fatalf("menu_tutor label for %s leaked Russian hardcode: %q", language.Code, text)
					}
				}
			}
		}
		if !found {
			t.Fatalf("missing menu_tutor for %s", language.Code)
		}
	}
}

func TestPhrasebookMenuLabelIsLocalizedForEveryInterfaceLanguage(t *testing.T) {
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code, InterfaceSelected: true})
		if strings.TrimSpace(copy.Phrasebook) == "" {
			t.Fatalf("missing phrasebook label for %s", language.Code)
		}
		keyboard := mainMenuInlineKeyboard(copy)
		rows := keyboard["inline_keyboard"].([][]map[string]any)
		found := false
		for _, row := range rows {
			for _, button := range row {
				if button["callback_data"] == "menu_phrasebook" {
					found = true
					text, _ := button["text"].(string)
					if !strings.Contains(text, copy.Phrasebook) {
						t.Fatalf("menu_phrasebook label for %s = %q, want copy label %q", language.Code, text, copy.Phrasebook)
					}
					if language.Code != "en" && strings.Contains(text, "Phrasebook") {
						t.Fatalf("menu_phrasebook label for %s leaked English fallback: %q", language.Code, text)
					}
				}
			}
		}
		if !found {
			t.Fatalf("missing menu_phrasebook for %s", language.Code)
		}
	}
}

func TestBottomReplyKeyboardUsesLocalizedButtons(t *testing.T) {
	copy := ui(userState{InterfaceLanguage: "zh", InterfaceSelected: true})
	keyboard := bottomReplyKeyboard(copy)
	rows := keyboard["keyboard"].([][]map[string]string)
	if rows[0][0]["text"] != copy.MenuButton {
		t.Fatalf("expected menu button %q, got %q", copy.MenuButton, rows[0][0]["text"])
	}
	if rows[0][1]["text"] != copy.StopButton {
		t.Fatalf("expected stop button %q, got %q", copy.StopButton, rows[0][1]["text"])
	}
}

func TestRuntimeCopyMentionsLocalizedMenuButton(t *testing.T) {
	copy := ui(userState{InterfaceLanguage: "zh", InterfaceSelected: true})
	if !strings.Contains(copy.UnknownButton, copy.MenuButton) {
		t.Fatalf("unknown-button text should reference localized menu button: %q", copy.UnknownButton)
	}
	if !strings.Contains(copy.Stopped, copy.MenuButton) {
		t.Fatalf("stopped text should reference localized menu button: %q", copy.Stopped)
	}
	if !strings.Contains(copy.Tool.VoicePremiumRequired, copy.MenuButton) {
		t.Fatalf("tool text should reference localized menu button: %q", copy.Tool.VoicePremiumRequired)
	}
}

func TestCJKToolInlineLabelsAreLocalized(t *testing.T) {
	user := userState{InterfaceLanguage: "ja", InterfaceSelected: true}
	copy := ui(user)

	translatorKeyboard := translatorLanguageKeyboard(user, false)
	rows := translatorKeyboard["inline_keyboard"].([][]map[string]any)
	if rows[0][0]["text"] != "↤ "+copy.Tool.AutoDetect {
		t.Fatalf("expected auto-detect label %q, got %q", copy.Tool.AutoDetect, rows[0][0]["text"])
	}

	mainKeyboard := mainMenuInlineKeyboard(copy)
	mainRows := mainKeyboard["inline_keyboard"].([][]map[string]any)
	gotWebApp := false
	for _, row := range mainRows {
		for _, button := range row {
			if text, _ := button["text"].(string); strings.Contains(text, copy.Tool.WebApp) {
				gotWebApp = true
			}
		}
	}
	if !gotWebApp {
		t.Fatalf("main keyboard does not contain localized web app label %q: %#v", copy.Tool.WebApp, mainRows)
	}

	toolsKeyboard := toolsInlineKeyboard("https://example.test/app", copy)
	toolRows := toolsKeyboard["inline_keyboard"].([][]map[string]any)
	for _, row := range toolRows {
		for _, button := range row {
			if text, _ := button["text"].(string); strings.Contains(text, copy.Tool.WebApp) {
				t.Fatalf("tools keyboard should not contain web app label %q: %#v", copy.Tool.WebApp, toolRows)
			}
		}
	}
}

func TestTranslatorLanguageKeyboardUsesInterfaceLanguageNames(t *testing.T) {
	user := userState{InterfaceLanguage: "ru", InterfaceSelected: true}
	keyboard := translatorLanguageKeyboard(user, true)
	rows := keyboard["inline_keyboard"].([][]map[string]any)

	foundJapanese := false
	for _, row := range rows {
		if len(row) != 1 {
			t.Fatalf("translator language buttons should be one per row: %#v", rows)
		}
		if row[0]["text"] == "↦ 日本語 (японский)" {
			foundJapanese = true
		}
	}
	if !foundJapanese {
		t.Fatalf("expected Japanese button to include Russian helper name: %#v", rows)
	}

	englishUser := userState{InterfaceLanguage: "en", InterfaceSelected: true}
	if got := translatorLanguageDisplayName("ja", englishUser); got != "日本語 (Japanese)" {
		t.Fatalf("expected English helper name for Japanese, got %q", got)
	}
}

func TestMenuAndStopButtonsAreRecognizedForEveryInterfaceLanguage(t *testing.T) {
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code})
		if !isMenuCommandText(copy.MenuButton) {
			t.Fatalf("menu button for %s is not recognized: %q", language.Code, copy.MenuButton)
		}
		if !isStopCommandText(copy.StopButton) {
			t.Fatalf("stop button for %s is not recognized: %q", language.Code, copy.StopButton)
		}
	}
}

func TestBackButtonsReturnToMainMenu(t *testing.T) {
	for name, keyboard := range map[string]map[string]any{
		"leaderboard": leaderboardMenuKeyboard(ui(userState{InterfaceLanguage: "ru"})),
		"manualLevel": manualLevelSelectionKeyboard(ui(userState{InterfaceLanguage: "ru"})),
	} {
		callbacks := collectCallbackData(t, keyboard)
		if !callbacks["back_menu"] {
			t.Fatalf("%s keyboard is missing back_menu callback; got %#v", name, callbacks)
		}
		if callbacks["menu_stats"] || callbacks["menu_level_test"] {
			t.Fatalf("%s keyboard still points to stale submenu; got %#v", name, callbacks)
		}
	}
}

func TestToolModeBackReturnsToToolsMenu(t *testing.T) {
	callbacks := collectCallbackData(t, backToToolsKeyboard(ui(userState{InterfaceLanguage: "ru"})))
	if !callbacks["menu_tools"] {
		t.Fatalf("tool mode back keyboard should return to tools menu; got %#v", callbacks)
	}
	if callbacks["back_menu"] {
		t.Fatalf("tool mode back keyboard should not jump to the main menu; got %#v", callbacks)
	}
}

func TestChangingBotLanguageReturnsCompletedUserToMainMenu(t *testing.T) {
	var methods []string
	var payloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		methods = append(methods, r.URL.Path)
		payloads = append(payloads, payload)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	store := &jsonStore{
		path:  filepath.Join(t.TempDir(), "users.json"),
		users: map[int64]userState{},
	}
	user := userState{
		TelegramID:        123,
		FirstName:         "Test",
		InterfaceLanguage: "ru",
		InterfaceSelected: true,
		TimezoneSelected:  true,
		LanguageSelected:  true,
		LearningLanguage:  "en",
	}
	store.users[user.TelegramID] = user
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}

	ctx := contextWithTelegramEditTarget(context.Background(), 42, 77)
	if err := b.handleInterfaceLanguageChoice(ctx, 42, user, "ui|en"); err != nil {
		t.Fatalf("handleInterfaceLanguageChoice() error = %v", err)
	}
	if len(methods) != 2 || methods[0] != "/sendMessage" || methods[1] != "/editMessageText" {
		t.Fatalf("expected language confirmation and edited main menu, got methods=%v payloads=%#v", methods, payloads)
	}
	if !strings.Contains(payloads[1]["text"].(string), "Main Menu") {
		t.Fatalf("expected edited message to be the localized main menu, got %#v", payloads[1]["text"])
	}
}

func TestInlineMarkdownFallsBackToPlainTextOnParseError(t *testing.T) {
	var calls []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		calls = append(calls, payload)
		if _, ok := payload["parse_mode"]; ok {
			_, _ = w.Write([]byte(`{"ok":false,"description":"Bad Request: can't parse entities: Character '-' is reserved"}`))
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := &telegramClient{baseURL: server.URL, http: server.Client()}
	if err := client.sendInlineMarkdownMessage(context.Background(), 42, "*Top* \\- ok", mainMenuInlineKeyboard()); err != nil {
		t.Fatalf("sendInlineMarkdownMessage() error = %v", err)
	}
	if len(calls) != 2 {
		t.Fatalf("expected markdown attempt and plain fallback, got %d calls", len(calls))
	}
	if _, ok := calls[1]["parse_mode"]; ok {
		t.Fatalf("fallback should not use MarkdownV2: %#v", calls[1])
	}
	if calls[1]["text"] != "Top - ok" {
		t.Fatalf("unexpected fallback text: %#v", calls[1]["text"])
	}
}

func TestBuyCommandOpensPremiumPlanMenu(t *testing.T) {
	var methods []string
	var payloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		methods = append(methods, r.URL.Path)
		payloads = append(payloads, payload)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	b := &bot{
		cfg: config{
			PremiumRubPrice:       300,
			PremiumStarsPrice:     150,
			PremiumYearRubPrice:   3000,
			PremiumYearStarsPrice: 1500,
		},
		telegram: &telegramClient{baseURL: server.URL, http: server.Client()},
	}

	if err := b.handleCommand(context.Background(), 42, "/buy", userState{InterfaceLanguage: "en"}); err != nil {
		t.Fatalf("handleCommand(/buy) error = %v", err)
	}
	if len(methods) != 1 || methods[0] != "/sendMessage" {
		t.Fatalf("expected /buy to open a tariff menu, got methods=%v payloads=%#v", methods, payloads)
	}
	replyMarkup, ok := payloads[0]["reply_markup"].(map[string]any)
	if !ok {
		t.Fatalf("expected reply markup, got %#v", payloads[0]["reply_markup"])
	}
	rows, ok := replyMarkup["inline_keyboard"].([]any)
	if !ok || len(rows) < 2 {
		t.Fatalf("expected premium plan rows, got %#v", replyMarkup["inline_keyboard"])
	}
	firstRow, ok := rows[0].([]any)
	if !ok || len(firstRow) == 0 {
		t.Fatalf("expected first premium plan row, got %#v", rows[0])
	}
	firstButton, ok := firstRow[0].(map[string]any)
	if !ok || firstButton["callback_data"] != "premium_plan|"+premiumMonthlyProduct {
		t.Fatalf("expected monthly plan callback, got %#v", firstRow[0])
	}
}

func TestStartCommandShowsPrivacyPolicyBeforeOnboarding(t *testing.T) {
	var methods []string
	var payloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		methods = append(methods, r.URL.Path)
		payloads = append(payloads, payload)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	b := &bot{telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	if err := b.handleCommand(context.Background(), 42, "/start", userState{TelegramID: 42}); err != nil {
		t.Fatalf("handleCommand(/start) error = %v", err)
	}
	if len(methods) != 1 || methods[0] != "/sendMessage" {
		t.Fatalf("expected a single privacy prompt, got methods=%v payloads=%#v", methods, payloads)
	}
	if !strings.Contains(payloads[0]["text"].(string), "Политику обработки персональных данных") {
		t.Fatalf("expected privacy prompt, got %q", payloads[0]["text"])
	}
	replyMarkup, ok := payloads[0]["reply_markup"].(map[string]any)
	if !ok {
		t.Fatalf("expected reply markup, got %#v", payloads[0]["reply_markup"])
	}
	rows, ok := replyMarkup["inline_keyboard"].([]any)
	if !ok || len(rows) != 2 {
		t.Fatalf("expected continue and policy rows, got %#v", replyMarkup["inline_keyboard"])
	}
	continueRow := rows[0].([]any)
	continueButton := continueRow[0].(map[string]any)
	if continueButton["callback_data"] != "privacy_continue" {
		t.Fatalf("expected privacy continue callback, got %#v", continueButton)
	}
	policyRow := rows[1].([]any)
	policyButton := policyRow[0].(map[string]any)
	if policyButton["url"] != privacyPolicyURL {
		t.Fatalf("expected privacy URL %q, got %#v", privacyPolicyURL, policyButton)
	}
}

func collectCallbackData(t *testing.T, keyboard map[string]any) map[string]bool {
	t.Helper()
	result := map[string]bool{}
	rows, ok := keyboard["inline_keyboard"].([][]map[string]any)
	if !ok {
		t.Fatalf("unexpected keyboard shape: %#v", keyboard)
	}
	for _, row := range rows {
		for _, button := range row {
			if callback, ok := button["callback_data"].(string); ok {
				result[callback] = true
			}
		}
	}
	return result
}

func buttonTextByCallback(t *testing.T, keyboard map[string]any, callback string) string {
	t.Helper()
	rows, ok := keyboard["inline_keyboard"].([][]map[string]any)
	if !ok {
		t.Fatalf("unexpected keyboard shape: %#v", keyboard)
	}
	for _, row := range rows {
		for _, button := range row {
			if button["callback_data"] == callback {
				text, _ := button["text"].(string)
				return text
			}
		}
	}
	t.Fatalf("missing callback %q in keyboard %#v", callback, keyboard)
	return ""
}
