package main

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type fakeAITutorClient struct {
	responses []string
	calls     int
	messages  [][]chatMessage
}

func (f *fakeAITutorClient) complete(ctx context.Context, messages []chatMessage, temperature float64, maxTokens int) (string, error) {
	f.calls++
	f.messages = append(f.messages, messages)
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

func TestAITutorStartReusesActiveSessionForRepeatedStarts(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	body, _ := json.Marshal(payload)
	ai := &fakeAITutorClient{responses: []string{
		string(body),
		`{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`,
	}}
	engine := newAITutorEngine(store, ai)
	engine.now = func() time.Time { return time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC) }
	user := userState{TelegramID: 81, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}

	first, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("first Start() error = %v", err)
	}
	second, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("second Start() error = %v", err)
	}
	if first.Session.ID != second.Session.ID {
		t.Fatalf("second start created a new active session: first=%q second=%q", first.Session.ID, second.Session.ID)
	}
	if first.Lesson.ID == second.Lesson.ID {
		return
	}
	t.Fatalf("second start changed active lesson: first=%q second=%q", first.Lesson.ID, second.Lesson.ID)
}

func TestAITutorStartUsesUnusedGeneratedSequenceAfterCompletion(t *testing.T) {
	store := newTestJSONStore(t)
	firstPayload := validAITutorLessonPayloadForTest()
	firstBody, _ := json.Marshal(firstPayload)
	secondPayload := validAITutorLessonPayloadForTest()
	secondPayload.Title = "A Park Visit"
	secondPayload.Theme = "sports routines"
	secondBody, _ := json.Marshal(secondPayload)
	ai := &fakeAITutorClient{responses: []string{
		string(firstBody),
		`{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`,
		string(secondBody),
		`{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`,
	}}
	engine := newAITutorEngine(store, ai)
	engine.now = func() time.Time { return time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC) }
	user := userState{TelegramID: 82, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}

	first, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("first Start() error = %v", err)
	}
	if err := store.updateAITutorSessionStage(first.Session.ID, aiTutorStageComplete, aiTutorSessionComplete, formatDBTime(time.Date(2026, 6, 10, 12, 5, 0, 0, time.UTC))); err != nil {
		t.Fatalf("complete first session: %v", err)
	}
	second, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("second Start() error = %v", err)
	}
	if first.Lesson.ID == second.Lesson.ID {
		t.Fatalf("second start reused completed lesson %q", first.Lesson.ID)
	}
	if len(ai.messages) < 4 {
		t.Fatalf("captured AI messages = %d, want two generation calls and two preflight calls", len(ai.messages))
	}
	firstPrompt := ai.messages[0][1].Content
	secondPrompt := ai.messages[2][1].Content
	if firstPrompt == secondPrompt {
		t.Fatalf("second generation prompt reused the same topic seed:\n%s", secondPrompt)
	}
}

func TestAITutorStartUsesLearningFocusAsTopicSeed(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	body, _ := json.Marshal(payload)
	ai := &fakeAITutorClient{responses: []string{
		string(body),
		`{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`,
	}}
	engine := newAITutorEngine(store, ai)
	user := userState{TelegramID: 80, FirstName: "demo", InterfaceLanguage: "de", LearningLanguage: "en", Level: "A2", LearningFocus: "work meetings and business travel"}
	if _, err := engine.Start(context.Background(), user, "web"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if len(ai.messages) == 0 || len(ai.messages[0]) < 2 {
		t.Fatalf("generation prompt was not captured: %#v", ai.messages)
	}
	prompt := ai.messages[0][1].Content
	for _, want := range []string{"Topic/theme seed: work meetings and business travel", "Interface language: German", "Target learning language: English"} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("generation prompt misses %q:\n%s", want, prompt)
		}
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

func TestAITutorWordRecallRejectsWrongChoice(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorWordRecallStage(1), Status: aiTutorSessionActive})
	engine := newAITutorEngine(store, &fakeAITutorClient{})
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Choice: "w2"})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorWordRecallStage(1) {
		t.Fatalf("wrong recall choice advanced to %q", result.Session.CurrentStage)
	}
	if result.Feedback.OK || !strings.Contains(result.Feedback.Message, "Try again") {
		t.Fatalf("wrong recall choice feedback = %#v", result.Feedback)
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
	user, err := store.getOrCreateUser(9, "demo")
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if user.XP != aiTutorCompletionXP || user.LessonCount != 1 || user.LessonsToday != 1 {
		t.Fatalf("AI Tutor completion reward not applied: xp=%d lessons=%d today=%d", user.XP, user.LessonCount, user.LessonsToday)
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

func TestAITutorSubmitCheckerRejectsExplicitFalseWithoutAdvancing(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageRetell, Status: aiTutorSessionActive})
	ai := &fakeAITutorClient{responses: []string{`{"ok":false,"feedback_interface":"Add one more story detail before moving on.","mistakes":["missing detail"]}`}}
	engine := newAITutorEngine(store, ai)
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Text: "Mia shop."})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageRetell {
		t.Fatalf("explicit checker rejection advanced to %q", result.Session.CurrentStage)
	}
	if result.Feedback.OK || !strings.Contains(result.Feedback.Message, "Add one more story detail") {
		t.Fatalf("feedback = %#v", result.Feedback)
	}
	answer := store.aiTutorAnswers["session-1|"+aiTutorStageRetell]
	if answer.Correct {
		t.Fatalf("rejected checker answer was stored as correct: %#v", answer)
	}
}

func TestAITutorSubmitProductionRequiresAtLeastTwoSentences(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageProduction, Status: aiTutorSessionActive})
	ai := &fakeAITutorClient{responses: []string{`{"ok":true,"corrected_version_target":"Mia buys bread. She pays at the shop.","recommendations_interface":["Good."]}`}}
	engine := newAITutorEngine(store, ai)
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Text: "Mia buys bread at the shop."})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageProduction {
		t.Fatalf("one-sentence production answer advanced to %q", result.Session.CurrentStage)
	}
	feedback := strings.ToLower(result.Feedback.Message)
	if result.Feedback.OK || (!strings.Contains(feedback, "two") && !strings.Contains(feedback, "два")) {
		t.Fatalf("feedback = %#v", result.Feedback)
	}
	if ai.calls != 0 {
		t.Fatalf("AI checker should not run for structurally invalid production answer, calls=%d", ai.calls)
	}
}

func TestAITutorSubmitProductionShowsFinalRecommendations(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageProduction, Status: aiTutorSessionActive})
	ai := &fakeAITutorClient{responses: []string{`{"ok":true,"used_words":["shop","bread"],"missing_requirement":"","corrected_version_target":"Mia buys bread at the shop.","recommendations_interface":["Repeat shop and bread tomorrow.","Keep sentences short at A1."],"mistakes":[]}`}}
	engine := newAITutorEngine(store, ai)
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Text: "Mia buys bread in the shop. She pays for bread."})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if !strings.Contains(result.Feedback.Message, "Repeat shop and bread tomorrow.") {
		t.Fatalf("production feedback should include final recommendations, got %#v", result.Feedback)
	}
}

func TestAITutorRecallOptionsAreShuffledWithoutForcingCorrectAwayFromFirst(t *testing.T) {
	words := []aiTutorWord{
		{ID: "a", Target: "alpha"},
		{ID: "b", Target: "beta"},
		{ID: "c", Target: "gamma"},
		{ID: "d", Target: "delta"},
		{ID: "e", Target: "epsilon"},
		{ID: "f", Target: "zeta"},
	}
	changedOrder := false
	correctFirst := false
	for index := range words {
		options := aiTutorRecallOptions(words, index)
		if len(options) != 4 {
			t.Fatalf("options for word %d = %d, want 4", index+1, len(options))
		}
		if options[0].ID == words[index].ID {
			correctFirst = true
		}
		sourceOrder := []string{words[index].ID}
		for sourceIndex, word := range words {
			if sourceIndex != index {
				sourceOrder = append(sourceOrder, word.ID)
			}
			if len(sourceOrder) == len(options) {
				break
			}
		}
		for optionIndex := range options {
			if options[optionIndex].ID != sourceOrder[optionIndex] {
				changedOrder = true
				break
			}
		}
	}
	if !changedOrder {
		t.Fatalf("recall options kept source order for every word")
	}
	if !correctFirst {
		t.Fatalf("recall shuffle appears to force the correct answer away from the first slot")
	}
}

func TestAITutorStepOptionsDoNotExposeCorrectFlag(t *testing.T) {
	step := aiTutorBuildStep(validAITutorLessonPayloadForTest(), aiTutorWordRecallStage(1))
	body, err := json.Marshal(step)
	if err != nil {
		t.Fatalf("marshal step: %v", err)
	}
	if strings.Contains(string(body), `"correct"`) {
		t.Fatalf("step JSON exposes correct answer metadata: %s", string(body))
	}
}

func TestAITutorBuildStepUsesSpecificLocalizedInstructions(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()
	lesson.Title = "A Beach Evening"
	lesson.Theme = "a short travel story"
	lesson.Story.TitleInterface = "Вечер на пляже"
	lesson.RetellTask.InstructionInterface = ""
	lesson.ProductionTask.InstructionInterface = ""

	story := aiTutorBuildStep(lesson, aiTutorStageStoryIntro)
	if story.Title != "Вечер на пляже" {
		t.Fatalf("story title = %q, want story_title_interface", story.Title)
	}
	if strings.Contains(story.Instruction, "Дальше") || strings.Contains(story.Instruction, "?") || story.Instruction == "" {
		t.Fatalf("story instruction is not specific: %q", story.Instruction)
	}

	retell := aiTutorBuildStep(lesson, aiTutorStageRetell)
	if !strings.Contains(retell.Instruction, "2-3") || !strings.Contains(strings.ToLower(retell.Instruction), "перескаж") {
		t.Fatalf("retell instruction is not specific: %q", retell.Instruction)
	}

	production := aiTutorBuildStep(lesson, aiTutorStageProduction)
	if !strings.Contains(production.Instruction, "2-3") || !strings.Contains(production.Instruction, "3") {
		t.Fatalf("production instruction is not specific: %q", production.Instruction)
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
