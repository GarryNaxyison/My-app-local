package main

import (
	"errors"
	"log"
	"time"
)

func (s *jsonStore) aiTutorSessionAnswers(id string) ([]aiTutorAnswerRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	answers := []aiTutorAnswerRecord{}
	for _, answer := range s.aiTutorAnswers {
		if answer.SessionID == id {
			answers = append(answers, answer)
		}
	}
	return answers, nil
}

func (s *sqliteStore) aiTutorSessionAnswers(id string) ([]aiTutorAnswerRecord, error) {
	rows, err := s.db.Query(`SELECT session_id,stage,answer_text,answer_json,feedback_json,correct,created_at FROM ai_tutor_answers WHERE session_id=? ORDER BY created_at,stage`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	answers := []aiTutorAnswerRecord{}
	for rows.Next() {
		var answer aiTutorAnswerRecord
		if err := rows.Scan(&answer.SessionID, &answer.Stage, &answer.AnswerText, &answer.AnswerJSON, &answer.FeedbackJSON, &answer.Correct, &answer.CreatedAt); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}
	return answers, rows.Err()
}

// Completion and its reward share a commit boundary. Replays are no-ops.
func (s *sqliteStore) completeAITutorSession(id string, userID int64, amount int, stamp, reviewChoice string) error {
	user, exists, err := s.getUser(userID)
	if err == nil && !exists {
		user, err = s.getOrCreateUser(userID, "")
	}
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.Exec(`UPDATE ai_tutor_sessions SET status='complete', current_stage='complete', completed_at=?, updated_at=? WHERE id=? AND telegram_id=? AND status<>'complete'`, stamp, stamp, id, userID)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return tx.Commit()
	}
	reviewDueAt := aiTutorReviewDueAt(parseDBTime(stamp), reviewChoice)
	var lessonID string
	if err := tx.QueryRow(`SELECT lesson_id FROM ai_tutor_sessions WHERE id=?`, id).Scan(&lessonID); err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE ai_tutor_reviews SET status='complete', updated_at=?
		WHERE telegram_id=? AND lesson_id=? AND status='scheduled' AND
		(? <> '' OR julianday(due_at) <= julianday(?))`, stamp, userID, lessonID, reviewDueAt, stamp)
	if err != nil {
		return err
	}
	if reviewDueAt != "" {
		if err := upsertAITutorReview(tx, userID, lessonID, reviewDueAt, reviewChoice, stamp); err != nil {
			return err
		}
	}
	now := time.Now().UTC()
	day := todayUTC(now)
	habitDate := localDateForUser(user, now)
	_, err = tx.Exec(`UPDATE users SET xp=xp+?, lesson_count=lesson_count+1, lessons_today=CASE WHEN daily_date=? THEN lessons_today+1 ELSE 1 END, practice_today=CASE WHEN daily_date=? THEN practice_today ELSE 0 END, voice_today=CASE WHEN daily_date=? THEN voice_today ELSE 0 END, daily_date=?, updated_at=? WHERE telegram_id=?`, amount, day, day, day, day, stamp, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`INSERT INTO habit_days (telegram_id,habit_date,login,complete,claimed,claimed_at,lessons,practice,voice,updated_at) SELECT telegram_id,?,1,CASE WHEN lessons_today>=1 AND practice_today>=1 AND voice_today>=1 THEN 1 ELSE 0 END,0,'',lessons_today,practice_today,voice_today,? FROM users WHERE telegram_id=? ON CONFLICT(telegram_id,habit_date) DO UPDATE SET login=1, complete=MAX(habit_days.complete,excluded.complete), lessons=excluded.lessons, practice=excluded.practice, voice=excluded.voice, updated_at=excluded.updated_at`, habitDate, stamp, userID)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	if _, err := s.rewardReferralLevel(userID); err != nil {
		log.Printf("AI Tutor referral reward pending for user %d: %v", userID, err)
	}
	return nil
}

func (s *jsonStore) completeAITutorSession(id string, userID int64, amount int, stamp, reviewChoice string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.aiTutorSessions[id]
	if !ok || session.TelegramID != userID {
		return errors.New("AI Tutor session not found")
	}
	if aiTutorSessionIsCompleted(session) {
		return nil
	}
	previous, exists := s.users[userID]
	user := previous
	user.HabitLog = make(map[string]habitDayEntry, len(previous.HabitLog))
	for key, day := range previous.HabitLog {
		user.HabitLog[key] = day
	}
	if !exists {
		user = newUserState(userID, time.Now().UTC())
	}
	normalizeUser(&user, time.Now().UTC())
	user.XP += amount
	user.LessonCount++
	user.LessonsToday++
	recordHabitDay(&user, time.Now().UTC(), true)
	s.users[userID] = user
	inviter, inviterExists := s.users[user.InvitedBy]
	s.rewardReferralLevelLocked(userID, time.Now().UTC())
	if err := s.saveLocked(); err != nil {
		if inviterExists {
			s.users[user.InvitedBy] = inviter
		}
		if exists {
			s.users[userID] = previous
		} else {
			delete(s.users, userID)
		}
		return err
	}
	session.Status = aiTutorSessionComplete
	session.CurrentStage = aiTutorStageComplete
	session.CompletedAt = stamp
	session.UpdatedAt = stamp
	s.aiTutorSessions[id] = session
	completedAt := parseDBTime(stamp)
	reviewDueAt := aiTutorReviewDueAt(completedAt, reviewChoice)
	for key, review := range s.aiTutorReviews {
		if review.TelegramID != userID || review.LessonID != session.LessonID || review.Status != "scheduled" {
			continue
		}
		due := parseDBTime(review.DueAt)
		if reviewDueAt != "" || (!due.IsZero() && !due.After(completedAt)) {
			review.Status = "complete"
			review.UpdatedAt = stamp
			s.aiTutorReviews[key] = review
		}
	}
	if reviewDueAt != "" {
		s.scheduleAITutorReviewLocked(userID, session.LessonID, reviewDueAt, reviewChoice, stamp)
	}
	return nil
}
