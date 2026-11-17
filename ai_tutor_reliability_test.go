package main

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
)

type unavailableTutorAI struct{}

func (unavailableTutorAI) complete(context.Context, []chatMessage, float64, int) (string, error) {
	return "", errors.New("model unavailable")
}

func TestTutorSQLiteCompletionAtomicAndReplaySafe(t *testing.T) {
	s, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer s.db.Close()
	_, err = s.getOrCreateUser(42, "test")
	if err != nil {
		t.Fatal(err)
	}
	err = s.createAITutorSession(aiTutorSessionRecord{ID: "atomic", TelegramID: 42, LessonID: "lesson", CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive})
	if err != nil {
		t.Fatal(err)
	}
	if err = s.scheduleAITutorReview(42, "lesson", "2026-11-13T12:00:00Z", "tomorrow"); err != nil {
		t.Fatal(err)
	}
	_, err = s.db.Exec(`CREATE TRIGGER reject_reward BEFORE UPDATE OF xp ON users BEGIN SELECT RAISE(ABORT,'simulated reward failure'); END`)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.completeAITutorSession("atomic", 42, 40, "2026-11-14T12:00:00Z", "no_review"); err == nil {
		t.Fatal("expected reward failure")
	}
	session, _, _ := s.getAITutorSession("atomic")
	if session.Status == aiTutorSessionComplete {
		t.Fatal("completion committed without reward")
	}
	if reviews, err := s.listScheduledAITutorReviews(42); err != nil || len(reviews) != 1 {
		t.Fatalf("failed completion changed the queue: %+v, %v", reviews, err)
	}
	_, err = s.db.Exec(`DROP TRIGGER reject_reward`)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err = s.completeAITutorSession("atomic", 42, 40, "2026-11-14T12:00:00Z", "no_review"); err != nil {
			t.Fatal(err)
		}
	}
	user, _, err := s.getUser(42)
	if err != nil {
		t.Fatal(err)
	}
	if user.XP != 40 || user.LessonCount != 1 {
		t.Fatalf("duplicate reward: %+v", user)
	}
	if reviews, err := s.listScheduledAITutorReviews(42); err != nil || len(reviews) != 0 {
		t.Fatalf("completed review is still scheduled: %+v, %v", reviews, err)
	}
}

func TestTutorStaleStageDoesNotAdvance(t *testing.T) {
	s := newTestJSONStore(t)
	_ = s.saveAITutorLesson(aiTutorLessonRecord{ID: "l", Payload: validAITutorLessonPayloadForTest(), Status: aiTutorStatusApproved})
	_ = s.createAITutorSession(aiTutorSessionRecord{ID: "s", LessonID: "l", TelegramID: 9, CurrentStage: aiTutorWordLearnStage(2), Status: aiTutorSessionActive})
	e := newAITutorEngine(s, nil)
	_, err := e.Submit(context.Background(), userState{TelegramID: 9}, "s", aiTutorSubmitInput{Stage: aiTutorWordLearnStage(1), Choice: "continue"})
	if err == nil {
		t.Fatal("stale answer accepted")
	}
	session, _, _ := s.getAITutorSession("s")
	if session.CurrentStage != aiTutorWordLearnStage(2) {
		t.Fatal("stale answer advanced lesson")
	}
}

func TestTutorCompletionReplacesScheduleOnce(t *testing.T) {
	sqlite, err := newSQLiteStore(filepath.Join(t.TempDir(), "reviews.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer sqlite.db.Close()
	for name, storage := range map[string]store{"sqlite": sqlite, "json": newTestJSONStore(t)} {
		t.Run(name, func(t *testing.T) {
			const completionTime = "2026-11-14T12:00:00Z"
			if err := storage.createAITutorSession(aiTutorSessionRecord{
				ID: "review", LessonID: "lesson", TelegramID: 42, Status: aiTutorSessionActive,
			}); err != nil {
				t.Fatal(err)
			}
			for _, seed := range []struct {
				owner            int64
				lesson, interval string
			}{
				{42, "lesson", "tomorrow"},
				{42, "lesson", "1_week"},
				{42, "other", "tomorrow"},
				{43, "lesson", "tomorrow"},
			} {
				if err := storage.scheduleAITutorReview(seed.owner, seed.lesson, "2026-11-13T12:00:00Z", seed.interval); err != nil {
					t.Fatal(err)
				}
			}
			if err := storage.completeAITutorSession("review", 42, 40, completionTime, "tomorrow"); err != nil {
				t.Fatal(err)
			}
			if err := storage.completeAITutorSession("review", 42, 40, completionTime, "3_days"); err != nil {
				t.Fatal(err)
			}
			reviews, err := storage.listScheduledAITutorReviews(42)
			if err != nil || len(reviews) != 2 {
				t.Fatalf("wrong schedule: %+v, %v", reviews, err)
			}
			for _, item := range reviews {
				if item.LessonID == "lesson" && (item.IntervalCode != "tomorrow" || item.DueAt != "2026-11-15T12:00:00Z") {
					t.Fatalf("replay changed the new schedule: %+v", item)
				}
			}
			if other, err := storage.listScheduledAITutorReviews(43); err != nil || len(other) != 1 {
				t.Fatalf("another account's schedule changed: %+v, %v", other, err)
			}
		})
	}
}

func TestTutorCompletionSurvivesModelOutage(t *testing.T) {
	s := newTestJSONStore(t)
	_ = s.saveAITutorLesson(aiTutorLessonRecord{ID: "l", Payload: validAITutorLessonPayloadForTest(), Status: aiTutorStatusInTrial})
	_ = s.createAITutorSession(aiTutorSessionRecord{ID: "s", LessonID: "l", TelegramID: 9, CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive})
	e := newAITutorEngine(s, unavailableTutorAI{})
	result, err := e.Submit(context.Background(), userState{TelegramID: 9}, "s", aiTutorSubmitInput{Choice: "no_review"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Session.Status != aiTutorSessionComplete {
		t.Fatal("model outage blocked completion")
	}
	_, err = e.Submit(context.Background(), userState{TelegramID: 9}, "s", aiTutorSubmitInput{Choice: "no_review"})
	if err != nil {
		t.Fatal(err)
	}
	user, _ := s.getOrCreateUser(9, "test")
	if user.XP != aiTutorCompletionXP {
		t.Fatalf("reward=%d", user.XP)
	}
}

func TestTutorPostCheckIncludesLearnerAnswers(t *testing.T) {
	s := newTestJSONStore(t)
	_ = s.saveAITutorLesson(aiTutorLessonRecord{ID: "l", Payload: validAITutorLessonPayloadForTest(), Status: aiTutorStatusInTrial})
	_ = s.saveAITutorAnswer(aiTutorAnswerRecord{SessionID: "s", Stage: aiTutorStageRetell, AnswerText: "My actual learner answer"})
	ai := &fakeAITutorClient{}
	e := newAITutorEngine(s, ai)
	lesson, _, _ := s.getAITutorLesson("l")
	if err := e.finishAITutorLesson(context.Background(), userState{TelegramID: 9}, lesson, "s"); err != nil {
		t.Fatal(err)
	}
	if len(ai.messages) != 1 || !strings.Contains(ai.messages[0][1].Content, "My actual learner answer") {
		t.Fatal("post-check omitted actual learner answers")
	}
}
