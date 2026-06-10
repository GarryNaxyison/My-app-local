package main

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

type fakeAITutorClient struct {
	responses []string
	calls     int
}

func (f *fakeAITutorClient) complete(ctx context.Context, messages []chatMessage, temperature float64, maxTokens int) (string, error) {
	f.calls++
	if len(f.responses) == 0 {
		return `{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`, nil
	}
	out := f.responses[0]
	f.responses = f.responses[1:]
	return out, nil
}

func newTestJSONStore(t *testing.T) *jsonStore {
	t.Helper()
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatalf("newJSONStore() error = %v", err)
	}
	return store
}

func TestAITutorStartGeneratesValidatedSession(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	body, _ := json.Marshal(payload)
	ai := &fakeAITutorClient{responses: []string{
		string(body),
		`{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`,
	}}
	engine := newAITutorEngine(store, ai)
	engine.now = func() time.Time { return time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC) }
	user := userState{TelegramID: 77, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}
	result, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageStoryIntro {
		t.Fatalf("stage = %q", result.Session.CurrentStage)
	}
	if result.NextStep.Stage != aiTutorStageStoryIntro {
		t.Fatalf("next step = %#v", result.NextStep)
	}
	if ai.calls != 2 {
		t.Fatalf("AI calls = %d, want generation + preflight", ai.calls)
	}
}

func TestAITutorStartRepairsDerivedWordTasksAndReviewOptions(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	payload.WordLearning = []aiTutorWordTask{{WordID: "w1"}}
	payload.WordRecall = []aiTutorWordTask{{WordID: "w1"}}
	payload.ReviewOptions = nil
	body, _ := json.Marshal(payload)
	ai := &fakeAITutorClient{responses: []string{
		string(body),
		`{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`,
	}}
	engine := newAITutorEngine(store, ai)
	user := userState{TelegramID: 78, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}
	result, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("Start() should repair deterministic lesson scaffolding, got error = %v", err)
	}
	if got := len(result.Lesson.Payload.WordLearning); got != 6 {
		t.Fatalf("word learning tasks = %d, want 6", got)
	}
	if got := len(result.Lesson.Payload.WordRecall); got != 6 {
		t.Fatalf("word recall tasks = %d, want 6", got)
	}
	if !hasAITutorReviewOptions(result.Lesson.Payload.ReviewOptions) {
		t.Fatalf("review options were not repaired: %#v", result.Lesson.Payload.ReviewOptions)
	}
}

func TestAITutorStartAcceptsApprovedTenPointPreflightScore(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	body, _ := json.Marshal(payload)
	ai := &fakeAITutorClient{responses: []string{
		string(body),
		`{"approved":true,"score":10,"critical_issues":[],"fix_suggestions":[],"reasons":["10/10 lesson"]}`,
	}}
	engine := newAITutorEngine(store, ai)
	user := userState{TelegramID: 79, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}
	result, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("Start() should accept approved 10/10 preflight score, got error = %v", err)
	}
	if result.Lesson.PreflightScore < aiTutorPromotionThreshold {
		t.Fatalf("preflight score = %d, want normalized score above threshold", result.Lesson.PreflightScore)
	}
}

func TestAITutorSubmitAdvancesThroughDeterministicStages(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatal(err)
	}
	session := aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageStoryIntro, Status: aiTutorSessionActive}
	if err := store.createAITutorSession(session); err != nil {
		t.Fatal(err)
	}
	engine := newAITutorEngine(store, &fakeAITutorClient{})
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Text: "continue"})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageRetell {
		t.Fatalf("stage = %q, want retell", result.Session.CurrentStage)
	}
}

func TestAITutorReviewScheduleCompletesSession(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	session := aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive}
	_ = store.createAITutorSession(session)
	engine := newAITutorEngine(store, &fakeAITutorClient{})
	engine.now = func() time.Time { return time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC) }
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Choice: "3_days"})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.Status != aiTutorSessionComplete {
		t.Fatalf("status = %q", result.Session.Status)
	}
	if result.Session.CurrentStage != aiTutorStageComplete {
		t.Fatalf("stage = %q", result.Session.CurrentStage)
	}
}

func TestAITutorSubmitRetellUsesAIChecker(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageRetell, Status: aiTutorSessionActive})
	ai := &fakeAITutorClient{responses: []string{`{"correct":true,"comprehension_score":88,"corrected_answer_target":"Mia goes to the shop.","feedback_interface":"Good.","mistakes":[]}`}}
	engine := newAITutorEngine(store, ai)
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Text: "Mia shop."})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageQuestion1 {
		t.Fatalf("stage = %q", result.Session.CurrentStage)
	}
	if result.Feedback.JSON == "" {
		t.Fatalf("missing checker JSON feedback")
	}
	if ai.calls != 1 {
		t.Fatalf("AI calls = %d", ai.calls)
	}
}

func TestAITutorCompletionPromotesApprovedTrialLesson(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusInTrial, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive})
	ai := &fakeAITutorClient{responses: []string{`{"approved":true,"score":93,"critical_issues":[],"fix_suggestions":[],"reasons":["clear"]}`}}
	engine := newAITutorEngine(store, ai)
	_, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Choice: "tomorrow"})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	updated, _, _ := store.getAITutorLesson("lesson-1")
	if updated.Status != aiTutorStatusApproved || updated.PostScore != 93 {
		t.Fatalf("lesson not promoted: %#v", updated)
	}
}
