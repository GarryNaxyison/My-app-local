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
